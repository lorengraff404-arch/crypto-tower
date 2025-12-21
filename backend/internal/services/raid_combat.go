package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// BattleResult represents the outcome of a single turn action
type BattleResult struct {
	Attacker       string  `json:"attacker"`
	Defender       string  `json:"defender"`
	MoveName       string  `json:"move_name"`
	Damage         int64   `json:"damage"`
	Effectiveness  float64 `json:"effectiveness"`   // Type effectiveness
	ClassAdvantage float64 `json:"class_advantage"` // Class advantage (Phase 15.1)
	IsCritical     bool    `json:"is_critical"`
	Message        string  `json:"message"`
	DefenderHP     int64   `json:"defender_hp"`
	DefenderMaxHP  int64   `json:"defender_max_hp"`
	TargetCharID   uint    `json:"target_char_id"` // For frontend animation targeting
}

// ExecutePlayerTurn processes a player character's attack turn (Axie-style)
func (s *RaidService) ExecutePlayerTurn(sessionID, characterID uint, moveSlot int) (*models.RaidSession, *BattleResult, error) {
	// 1. Load session with team data
	var session models.RaidSession
	if err := db.DB.Preload("Mission").
		Preload("Team.Members.Character.Moves").
		First(&session, sessionID).Error; err != nil {
		return nil, nil, errors.New("session not found")
	}

	// 2. Verify it's this character's turn
	currentTurn, err := s.getCurrentTurn(&session)
	if err != nil {
		return nil, nil, err
	}

	if currentTurn.Type != "player" {
		return nil, nil, errors.New("not player's turn (enemy turn)")
	}

	if currentTurn.CharID != characterID {
		return nil, nil, fmt.Errorf("not this character's turn (expected char %d, got %d)", currentTurn.CharID, characterID)
	}

	// 3. Get character and validate move
	var character *models.Character
	for _, member := range session.Team.Members {
		if member.Character.ID == characterID {
			character = &member.Character
			break
		}
	}

	if character == nil {
		return nil, nil, errors.New("character not in team")
	}

	if moveSlot < 0 || moveSlot >= len(character.Moves) {
		return nil, nil, errors.New("invalid move slot")
	}

	move := character.Moves[moveSlot]

	// 4. Check PP
	if move.CurrentPP <= 0 {
		return nil, nil, errors.New("move has no PP remaining")
	}

	// 4.5. Process status effects at turn start (Phase 15.3)
	// Check if character can act (stun/freeze check)
	sem := NewStatusEffectManager(session.CharacterStates)
	if !sem.CanAct() {
		// Character is stunned/frozen, skip turn
		result := &BattleResult{
			Attacker: character.Class,
			Defender: session.Mission.EnemyName,
			MoveName: "Stunned",
			Message:  fmt.Sprintf("%s is stunned and cannot act!", character.Class),
		}
		s.advanceTurn(&session)
		db.DB.Save(&session)
		return &session, result, nil
	}

	// Process DoT effects (burn, poison)
	dotDamage := sem.ProcessTurnEffects(int64(character.CurrentHP))
	if dotDamage > 0 {
		// Apply DoT damage to character
		var charStates []CharacterState
		json.Unmarshal([]byte(session.CharacterStates), &charStates)
		for i, state := range charStates {
			if state.CharID == characterID {
				charStates[i].CurrentHP -= dotDamage
				if charStates[i].CurrentHP < 0 {
					charStates[i].CurrentHP = 0
				}
				break
			}
		}
		statesJSON, _ := json.Marshal(charStates)
		session.CharacterStates = string(statesJSON)
	}

	// 5. Use shared skill logic instead of custom damage calculation
	ability := s.convertMoveToAbility(move)

	// Calculate outcome using centralized logic
	skillResult := s.skillService.CalculateSkillOutcome(character, ability)

	// Apply damage to enemy
	damage := int64(skillResult.Damage)
	session.CurrentBossHP -= damage
	if session.CurrentBossHP < 0 {
		session.CurrentBossHP = 0
	}
	session.DamageDealt += damage

	// 7. Decrement PP
	move.CurrentPP--
	db.DB.Save(&move)

	// 7.5. Build result message with effect and effectiveness
	effectMsg := skillResult.Message
	if skillResult.CriticalHit {
		effectMsg += " Critical hit!"
	}

	result := &BattleResult{
		Attacker:       character.Class,
		Defender:       session.Mission.EnemyName,
		MoveName:       move.Name,
		Damage:         damage,
		Effectiveness:  1.0, // Can enhance later with element matching
		ClassAdvantage: 1.0, // Can enhance later
		IsCritical:     skillResult.CriticalHit,
		Message:        effectMsg,
		DefenderHP:     session.CurrentBossHP,
		DefenderMaxHP:  session.Mission.EnemyHP,
	}

	// 9. Check for victory
	if session.CurrentBossHP <= 0 {
		session.Status = "COMPLETED"
		now := time.Now()
		session.CompletedAt = &now

		// Calculate rewards
		s.calculateRewards(&session)
		s.distributeExpToTeam(&session, session.XPEarned)
		s.updateCampaignProgress(session.UserID, session.Mission.IslandID, session.Mission.Sequence)
	} else {
		// 10. Advance turn (only if battle continues)
		s.advanceTurn(&session)
	}

	// 11. Save session
	db.DB.Save(&session)

	// 12. Reload to get updated data
	db.DB.Preload("Mission").
		Preload("Team.Members.Character.Moves").
		First(&session, session.ID)

	return &session, result, nil
}

// ExecuteEnemyTurn processes the enemy's attack turn (AI-controlled)
func (s *RaidService) ExecuteEnemyTurn(sessionID uint) (*models.RaidSession, *BattleResult, error) {
	// 1. Load session
	var session models.RaidSession
	if err := db.DB.Preload("Mission").
		Preload("Team.Members.Character").
		First(&session, sessionID).Error; err != nil {
		return nil, nil, errors.New("session not found")
	}

	// 2. Verify it's enemy's turn
	currentTurn, err := s.getCurrentTurn(&session)
	if err != nil {
		return nil, nil, err
	}

	if currentTurn.Type != "enemy" {
		return nil, nil, errors.New("not enemy's turn")
	}

	// 3. Select random living player character as target
	var charStates []CharacterState
	json.Unmarshal([]byte(session.CharacterStates), &charStates)

	livingChars := []CharacterState{}
	for _, state := range charStates {
		if state.CurrentHP > 0 {
			livingChars = append(livingChars, state)
		}
	}

	if len(livingChars) == 0 {
		// All characters fainted - defeat
		session.Status = "FAILED"
		now := time.Now()
		session.CompletedAt = &now
		db.DB.Save(&session)

		return &session, &BattleResult{
			Message: "All characters fainted! Defeat!",
		}, nil
	}

	// Random target (Go 1.20+ auto-seeds, no need for rand.Seed)
	target := livingChars[rand.Intn(len(livingChars))]

	// 4. Get target character info
	var targetChar *models.Character
	for _, member := range session.Team.Members {
		if member.Character.ID == target.CharID {
			targetChar = &member.Character
			break
		}
	}

	// 5. Calculate damage
	damage, _, classAdv, isCrit := s.calculateDamage(
		session.Mission.EnemyAtk,
		targetChar.CurrentDefense,
		50, // Enemy base power
		session.Mission.EnemyType,
		targetChar.Element,
		"Enemy",                      // Attacker class
		targetChar.Class,             // Defender class
		int(session.CurrentBossHP),   // Attacker HP
		int(session.Mission.EnemyHP), // Attacker Max HP
		int(target.CurrentHP),        // Defender HP
		int(targetChar.CurrentHP),    // Defender Max HP
		session.ID,                   // Session ID (for status effects)
		nil,                          // Attacker char ID (boss has none)
		true,                         // Attacker is enemy
		&target.CharID,               // Defender char ID
		false,                        // Defender is not enemy
	)

	// 6. Apply damage to target character
	s.updateCharacterHP(&session, target.CharID, target.CurrentHP-damage)
	session.TotalDamageTaken += damage

	// 7. Build result
	result := &BattleResult{
		Attacker:       session.Mission.EnemyName,
		Defender:       targetChar.Class,
		MoveName:       "Attack",
		Damage:         damage,
		ClassAdvantage: classAdv,
		IsCritical:     isCrit,
		Message:        fmt.Sprintf("%s attacked %s for %d damage!", session.Mission.EnemyName, targetChar.Class, damage),
		DefenderHP:     target.CurrentHP - damage,
		DefenderMaxHP:  int64(targetChar.CurrentHP),
		TargetCharID:   target.CharID, // For frontend animation
	}

	// 8. Check for defeat
	if s.checkAllPlayerCharactersFainted(&session) {
		session.Status = "FAILED"
		now := time.Now()
		session.CompletedAt = &now
	} else {
		// Advance turn
		s.advanceTurn(&session)
	}

	// 9. Save
	db.DB.Save(&session)

	// Reload
	db.DB.Preload("Mission").
		Preload("Team.Members.Character.Moves").
		First(&session, session.ID)

	return &session, result, nil
}

// calculateDamage computes damage with type, class, passives, status effects, and crits (Phase 15)
func (s *RaidService) calculateDamage(atk, def, power int, attackerElement, defenderElement, attackerClass, defenderClass string, attackerHP, attackerMaxHP, defenderHP, defenderMaxHP int, sessionID uint, attackerCharID *uint, attackerIsEnemy bool, defenderCharID *uint, defenderIsEnemy bool) (int64, float64, float64, bool) {
	// Base damage formula (simplified Pokemon)
	baseDmg := ((2*5)/5 + 2) * power * (atk / def) / 50

	// Type effectiveness (existing)
	typeEff := s.getTypeEffectiveness(attackerElement, defenderElement)

	// Class advantage (Phase 15.1)
	classAdv := GetClassAdvantage(attackerClass, defenderClass)

	// Critical hit (10% base + Archer bonus)
	baseCritChance := 0.1
	archerBonus := GetArcherCritBonus(attackerClass)

	// Status effect crit boost (Phase 15.3)
	_, _, _, critBoost := s.GetStatusEffectMultipliers(sessionID, attackerCharID, attackerIsEnemy)

	isCrit := rand.Float64() < (baseCritChance + archerBonus + critBoost - 0.1) // Subtract base 1.0
	critMult := 1.0
	if isCrit {
		critMult = 1.5
	}

	// Passive abilities (Phase 15.2)
	attackerPassive := ApplyPassiveAbility(attackerClass, attackerHP, attackerMaxHP, isCrit, true)
	defenderPassive := ApplyPassiveAbility(defenderClass, defenderHP, defenderMaxHP, isCrit, false)

	// Status effects (Phase 15.3)
	atkBuff, _, _, _ := s.GetStatusEffectMultipliers(sessionID, attackerCharID, attackerIsEnemy)
	_, defDebuff, _, _ := s.GetStatusEffectMultipliers(sessionID, defenderCharID, defenderIsEnemy)

	// Final damage = base * type * class * crit * attacker passive * attacker buffs / (defender passive * defender buffs)
	finalDmg := float64(baseDmg) * typeEff * classAdv * critMult * attackerPassive * atkBuff / (defenderPassive * defDebuff)

	// Add randomness (85%-100%)
	randomFactor := 0.85 + rand.Float64()*0.15
	finalDmg *= randomFactor

	if finalDmg < 1 {
		finalDmg = 1 // Minimum 1 damage
	}

	return int64(finalDmg), typeEff, classAdv, isCrit
}

// GetStatusEffectMultipliers retrieves status effect multipliers for a character/enemy (Phase 15.3)
func (s *RaidService) GetStatusEffectMultipliers(sessionID uint, characterID *uint, isEnemy bool) (atkMult, defMult, spdMult, critMult float64) {
	// For now, return neutral values (1.0, 1.0, 1.0, 0.0)
	// TODO: Query active_status_effects table when database is set up
	return 1.0, 1.0, 1.0, 0.0
}

// getTypeEffectiveness returns type matchup multiplier
func (s *RaidService) getTypeEffectiveness(attacker, defender string) float64 {
	// Simple type chart (expand as needed)
	effectiveness := map[string]map[string]float64{
		"FIRE": {
			"GRASS": 2.0,
			"WATER": 0.5,
			"FIRE":  0.5,
			"ICE":   2.0,
		},
		"WATER": {
			"FIRE":     2.0,
			"GRASS":    0.5,
			"ELECTRIC": 0.5,
		},
		"GRASS": {
			"WATER": 2.0,
			"FIRE":  0.5,
			"GRASS": 0.5,
		},
		"ELECTRIC": {
			"WATER": 2.0,
			"GRASS": 0.5,
		},
		"ICE": {
			"GRASS": 2.0,
			"FIRE":  0.5,
		},
	}

	if typeMap, ok := effectiveness[attacker]; ok {
		if mult, ok := typeMap[defender]; ok {
			return mult
		}
	}

	return 1.0 // Neutral
}

// calculateRewards determines tokens and XP earned (already implemented but referenced here)
func (s *RaidService) calculateRewards(session *models.RaidSession) {
	// Grade based on HP remaining
	hpPercent := float64(session.CurrentTeamHP) / float64(session.InitialTeamHP) * 100

	grade := "D"
	mult := 1.0
	if hpPercent >= 90 {
		grade = "S"
		mult = 2.0
	} else if hpPercent >= 70 {
		grade = "A"
		mult = 1.5
	} else if hpPercent >= 50 {
		grade = "B"
		mult = 1.2
	} else if hpPercent >= 30 {
		grade = "C"
		mult = 1.0
	}

	// Parse rewards from mission
	var rewards map[string]int
	json.Unmarshal([]byte(session.Mission.RewardsPool), &rewards)

	baseTokens := rewards["tokens"]
	baseXP := rewards["xp"]

	session.TokensEarned = int(float64(baseTokens) * mult)
	session.XPEarned = int(float64(baseXP) * mult)
	session.PerformanceGrade = grade
}
