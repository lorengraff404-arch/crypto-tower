package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// RaidService handles raid-related business logic
type RaidService struct {
	skillService *SkillActivationService
	ledger       *LedgerService
}

// RaidSessionWithSprites contains raid session data with character sprites loaded
type RaidSessionWithSprites struct {
	Session   *models.RaidSession   `json:"session"`
	TeamChars []models.Character    `json:"team_characters"`
	Mission   *models.IslandMission `json:"mission"`
}

// NewRaidService creates a new raid service
func NewRaidService() *RaidService {
	return &RaidService{
		skillService: NewSkillActivationService(),
		ledger:       NewLedgerService(),
	}
}

// ListIslands returns all available islands
func (s *RaidService) ListIslands() ([]models.Island, error) {
	var islands []models.Island
	err := db.DB.Find(&islands).Error // No need to preload bosses/missions here, heavy
	return islands, err
}

// GetActiveSession returns the current in-progress session for the user if any
func (s *RaidService) GetActiveSession(userID uint) (*models.RaidSession, error) {
	var session models.RaidSession
	err := db.DB.Preload("Mission").Preload("Boss").Preload("Team.Members.Character").
		Where("user_id = ? AND status = 'IN_PROGRESS'", userID).
		Where("expires_at > ?", time.Now()). // Check expiry
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active session
		}
		return nil, err
	}
	return &session, nil
}

type MissionDTO struct {
	models.IslandMission
	IsLocked    bool `json:"is_locked"`
	IsCompleted bool `json:"is_completed"`
}

func (s *RaidService) ListMissions(userID, islandID uint) ([]MissionDTO, error) {
	var missions []models.IslandMission
	if err := db.DB.Where("island_id = ?", islandID).Order("sequence asc").Find(&missions).Error; err != nil {
		return nil, err
	}

	// Get Progress
	var progress models.UserCampaignProgress
	db.DB.Where("user_id = ? AND island_id = ?", userID, islandID).First(&progress)
	// If not found, highest is 0

	var dtos []MissionDTO
	for _, m := range missions {
		dto := MissionDTO{IslandMission: m}
		if m.Sequence > progress.HighestSequence+1 {
			dto.IsLocked = true
		}
		if m.Sequence <= progress.HighestSequence {
			dto.IsCompleted = true
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

// Performance grading functions (Phase 9)
func (s *RaidService) calculatePerformanceGrade(session *models.RaidSession) string {
	// Calculate % damage taken
	damagePercent := float64(0)
	if session.InitialTeamHP > 0 {
		damagePercent = (float64(session.TotalDamageTaken) / float64(session.InitialTeamHP)) * 100
	}

	turns := session.TurnCount

	// Grading Logic (Pokemon-style)
	// S Rank: Fast victory (≤3 turns), minimal damage (<20%)
	if turns <= 3 && damagePercent < 20 {
		return "S"
	}

	// A Rank: Quick victory (≤5 turns), low damage (<40%)
	if turns <= 5 && damagePercent < 40 {
		return "A"
	}

	// B Rank: Decent victory (≤8 turns), moderate damage (<60%)
	if turns <= 8 && damagePercent < 60 {
		return "B"
	}

	// C Rank: Slow but safe
	if damagePercent < 80 {
		return "C"
	}

	// D Rank: Pyrrhic victory (barely survived)
	return "D"
}

func (s *RaidService) getGradeMultiplier(grade string) float64 {
	switch grade {
	case "S":
		return 3.0 // 3x rewards!
	case "A":
		return 2.0 // 2x rewards
	case "B":
		return 1.5 // 1.5x rewards
	case "C":
		return 1.0 // Normal rewards
	case "D":
		return 0.5 // Reduced rewards (you almost died!)
	default:
		return 1.0
	}
}

// StartRaidSession initiates a battle session
// NOW ACCEPTS MISSION ID instead of IslandID for the target
func (s *RaidService) StartRaidSession(userID, missionID, teamID uint) (*models.RaidSession, error) {
	// 0. Check for existing active session
	active, _ := s.GetActiveSession(userID)
	if active != nil {
		// Reload session with all required relationships
		var fullSession models.RaidSession
		if err := db.DB.
			Preload("Mission").
			Preload("Team").
			Preload("Team.Members").
			Preload("Team.Members.Character").
			Preload("Team.Members.Character.Moves").
			Preload("Boss").
			First(&fullSession, active.ID).Error; err != nil {
			return nil, fmt.Errorf("failed to reload session: %w", err)
		}

		return &fullSession, nil // Return existing if alive
	}

	// 1. Verify Team
	team, err := s.GetTeamIfValid(userID, teamID)
	if err != nil {
		return nil, err
	}

	// 2. Get Mission Info
	var mission models.IslandMission
	if err := db.DB.First(&mission, missionID).Error; err != nil {
		return nil, errors.New("mission not found")
	}

	// 2b. Verify Unlock Status
	var progress models.UserCampaignProgress
	db.DB.Where("user_id = ? AND island_id = ?", userID, mission.IslandID).First(&progress)
	if mission.Sequence > progress.HighestSequence+1 {
		return nil, errors.New("mission is locked")
	}

	// Calculate Initial Team HP
	var currentTeamHP int64 = 0
	for _, m := range team.Members {
		if !m.IsBackup {
			currentTeamHP += int64(m.Character.CurrentHP)
		}
	}

	// === AXIE-STYLE TURN SYSTEM (Phase 13) ===
	// Build Turn Queue (sorted by speed)
	turnQueue, err := s.buildTurnQueue(team, &mission)
	if err != nil {
		return nil, err
	}
	turnQueueJSON, _ := json.Marshal(turnQueue)

	// Initialize Character States (individual HP tracking)
	charStates, err := s.initializeCharacterStates(team)
	if err != nil {
		return nil, err
	}
	charStatesJSON, _ := json.Marshal(charStates)

	// Set first active character (first in turn queue)
	var activeCharID *uint
	if len(turnQueue) > 0 && turnQueue[0].Type == "player" {
		activeCharID = &turnQueue[0].CharID
	}

	// 3. Create Session
	// Note: We keep BossID as non-null foreign key for now, so we need a dummy or the real boss if mapped
	var dummyBoss models.RaidBoss
	db.DB.Where("island_id = ?", mission.IslandID).First(&dummyBoss)

	session := models.RaidSession{
		UserID:        userID,
		MissionID:     mission.ID,
		BossID:        &dummyBoss.ID, // Use pointer
		TeamID:        teamID,
		Status:        "IN_PROGRESS",
		CurrentStage:  1,
		TotalStages:   1,
		CurrentBossHP: mission.EnemyHP,
		TotalHP:       mission.EnemyHP, // Store max HP
		CurrentTeamHP: currentTeamHP,
		InitialTeamHP: currentTeamHP, // Track for performance %
		BattleSeed:    generateSeed(),
		ExpiresAt:     &[]time.Time{time.Now().Add(time.Minute * 15)}[0], // 15 min session timeout

		// Turn System Fields
		TurnQueue:         string(turnQueueJSON),
		CharacterStates:   string(charStatesJSON),
		ActiveCharacterID: activeCharID,
		CurrentTurnIndex:  0,
	}

	if err := db.DB.Create(&session).Error; err != nil {
		return nil, err
	}

	// Return full data for UI
	// Reload session to ensure all preloads are active
	if err := db.DB.Preload("Mission").
		Preload("Team").
		Preload("Team.Members").
		Preload("Team.Members.Character").
		Preload("Team.Members.Character.Moves"). // Critical for Move Selection
		Preload("Boss").
		First(&session, session.ID).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *RaidService) updateCampaignProgress(userID, islandID uint, sequence int) {
	var progress models.UserCampaignProgress
	// Find or Create
	err := db.DB.Where("user_id = ? AND island_id = ?", userID, islandID).First(&progress).Error
	if err != nil {
		// Create new
		progress = models.UserCampaignProgress{
			UserID:          userID,
			IslandID:        islandID,
			HighestSequence: sequence,
			UpdatedAt:       time.Now(),
		}
		db.DB.Create(&progress)
	} else {
		// Update only if higher
		if sequence > progress.HighestSequence {
			progress.HighestSequence = sequence
			progress.UpdatedAt = time.Now()
			db.DB.Save(&progress)
		}
	}
}

// ProcessTurn calculates damage and updates state securely
func (s *RaidService) ProcessTurn(sessionID uint) (*models.RaidSession, string, error) {
	var session models.RaidSession
	if err := db.DB.Preload("Mission").
		Preload("Boss").
		Preload("Team.Members.Character.StatusEffects").
		Preload("Team.Members.Character.Moves"). // Added for consistency
		First(&session, sessionID).Error; err != nil {
		return nil, "", errors.New("session not found")
	}

	if session.Status != "IN_PROGRESS" {
		return &session, "Battle already ended", nil
	}

	// Load status effects
	teamEffects := NewStatusEffectManager(session.ActiveStatusEffects)

	// Process DoT at start of turn
	dotDamage := teamEffects.ProcessTurnEffects(session.InitialTeamHP)
	if dotDamage > 0 {
		session.CurrentTeamHP -= dotDamage
		session.TotalDamageTaken += dotDamage
	}

	// Check if team can act (stun/sleep/freeze check)
	if !teamEffects.CanAct() {
		session.TurnCount++
		session.ActiveStatusEffects = teamEffects.ToJSON()
		db.DB.Save(&session)
		return &session, "Your team is unable to act this turn...", nil
	}

	// Reset team HP to max if needed (for first turn logic)
	if session.TurnCount == 0 {
		var maxHP int64 = 0
		for _, m := range session.Team.Members {
			if !m.IsBackup {
				maxHP += int64(m.Character.CurrentHP)
			}
		}
		session.CurrentTeamHP = maxHP
	}

	// 1. Player Attack with stat modifiers
	basePlayerDmg := int64(session.Team.TotalPower / 10)

	// Apply ATK buff/debuff
	atkModifier := teamEffects.GetStatModifier("atk")
	playerDmg := int64(float64(basePlayerDmg) * atkModifier)

	// Apply accuracy check
	if rand.Float64() > teamEffects.GetAccuracyModifier() {
		// Miss!
		session.TurnCount++
		session.ActiveStatusEffects = teamEffects.ToJSON()
		db.DB.Save(&session)
		return &session, "Attack missed!", nil
	}

	variance := float64(playerDmg) * 0.1
	playerDmg += int64(rand.Float64()*variance*2 - variance)
	if playerDmg < 1 {
		playerDmg = 1
	}

	session.CurrentBossHP -= playerDmg
	session.DamageDealt += playerDmg
	session.TurnCount++

	// Restore log generation
	logMsg := fmt.Sprintf("Dealt %d damage!", playerDmg)
	if atkModifier > 1.0 {
		logMsg += " (BOOSTED!)"
	}

	// Check Victory
	// Check Battle Status
	if session.CurrentBossHP <= 0 {
		session.CurrentBossHP = 0

		// If this was the last stage (Boss/Stage 5)
		if session.CurrentStage >= session.TotalStages {
			session.Status = "COMPLETED"
			now := time.Now()
			session.CompletedAt = &now

			// PHASE 9: Calculate Performance Grade
			grade := s.calculatePerformanceGrade(&session)
			session.PerformanceGrade = grade

			// logMsg := fmt.Sprintf("VICTORY! Performance: %s Rank!", grade)

			// Distribute Rewards with multipliers
			baseTokens := 50 + (session.Mission.Sequence * 25) // Scales with mission
			baseXP := 100 + (session.Mission.Sequence * 50)

			multiplier := s.getGradeMultiplier(grade)
			finalTokens := int(float64(baseTokens) * multiplier)
			finalXP := int(float64(baseXP) * multiplier)

			session.TokensEarned = finalTokens
			session.XPEarned = finalXP

			// Give rewards to user
			var user models.User
			if err := db.DB.First(&user, session.UserID).Error; err == nil {
				user.GTKBalance += int64(finalTokens)
				user.Experience += finalXP
				db.DB.Save(&user)
			}

			// ✅ PHASE 10.1: Distribute XP to team characters
			s.distributeExpToTeam(&session, finalXP)

			// Update Progress
			s.updateCampaignProgress(session.UserID, session.Mission.IslandID, session.Mission.Sequence)

			logMsg = fmt.Sprintf("VICTORY! %s Rank! +%d GTK, +%d XP", grade, finalTokens, finalXP)
		} else {
			// Shouldn't happen in mission mode (TotalStages = 1)
			session.Status = "COMPLETED"
			now := time.Now()
			session.CompletedAt = &now
			logMsg = "Mission Complete!"

			s.distributeRewards(session.UserID, 1)
			s.updateCampaignProgress(session.UserID, session.Mission.IslandID, session.Mission.Sequence)
		}
	} else {
		// Enemy Counter-Attack (Only if alive)
		// Enemy Damage = (EnemyAtk * 1.5) - (TeamAvDef / 2)
		teamDef := 0
		activeCount := 0
		for _, m := range session.Team.Members {
			if !m.IsBackup {
				teamDef += m.Character.CurrentDefense
				activeCount++
			}
		}
		avgDef := 0
		if activeCount > 0 {
			avgDef = teamDef / activeCount
		}

		baseAtk := session.Mission.EnemyAtk
		if session.Boss != nil { // Override if boss exists (rare)
			baseAtk = session.Boss.BaseAttack
		}

		bossDmg := int64(float64(baseAtk)*1.5) - int64(avgDef/2)

		// Apply team's DEF buff/debuff
		defModifier := teamEffects.GetStatModifier("def")
		if defModifier > 1.0 {
			// Higher DEF = less damage taken
			bossDmg = int64(float64(bossDmg) / defModifier)
		} else if defModifier < 1.0 {
			// Lower DEF = more damage taken
			bossDmg = int64(float64(bossDmg) * (1.0 / defModifier)) // Corrected calculation for debuff
		}

		if bossDmg < 10 {
			bossDmg = 10
		}

		session.CurrentTeamHP -= bossDmg
		session.TotalDamageTaken += bossDmg
		logMsg += fmt.Sprintf(" Enemy dealt %d damage!", bossDmg)

		// Wake from sleep if damaged
		teamEffects.WakeFromSleep()

		// Check Defeat
		if session.CurrentTeamHP <= 0 {
			session.CurrentTeamHP = 0
			session.Status = "FAILED"
			logMsg = "DEFEAT! Your team was wiped out."
		}
	}

	// Save status effects state
	session.ActiveStatusEffects = teamEffects.ToJSON()

	// Save session
	if err := db.DB.Save(&session).Error; err != nil {
		return nil, "", err
	}

	return &session, logMsg, nil
}

// distributeRewards handles atomic reward giving
// Now returns calculated rewards for session update
type RaidRewards struct {
	XP     int
	Tokens int64
}

func (s *RaidService) distributeRewards(userID uint, difficulty int) (RaidRewards, error) {
	var rewards RaidRewards

	// 1. Calculate Rewards based on difficulty (1-10)
	xpGain := 100 * difficulty
	tokenGain := int64(50*difficulty + rand.Intn(20*difficulty))

	rewards.XP = xpGain
	rewards.Tokens = tokenGain

	// 2. LEDGER INTEGRATION: Reward Tokens
	// Debit: Reward Pool (Inflation), Credit: User Wallet
	if s.ledger == nil {
		s.ledger = NewLedgerService()
	}

	rewardAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeReward, "GTK")
	userAcc, _ := s.ledger.GetOrCreateAccount(&userID, models.AccountTypeWallet, "GTK")

	entries := []models.LedgerEntry{
		{AccountID: rewardAcc.ID, Amount: -tokenGain, Type: "DEBIT"},
		{AccountID: userAcc.ID, Amount: tokenGain, Type: "CREDIT"},
	}

	// We use a transaction to ensure DB consistency with Ledger
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Log to Ledger
		if err := s.ledger.CreateTransactionWithTx(tx, models.TxTypeRaidReward, fmt.Sprintf("raid_reward_%d_%d", userID, time.Now().Unix()), fmt.Sprintf("Raid Reward (Diff %d)", difficulty), entries); err != nil {
			return err
		}

		// Update User Balance (Legacy Sync)
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		if err := tx.Model(&user).UpdateColumn("gtk_balance", gorm.Expr("gtk_balance + ?", tokenGain)).Error; err != nil {
			return err
		}

		return nil
	})

	return rewards, err
}

func (s *RaidService) GetTeamIfValid(userID, teamID uint) (*models.Team, error) {
	var team models.Team
	// Corrected line: Ensure Preload is correctly chained before First
	if err := db.DB.Preload("Members.Character").First(&team, teamID).Error; err != nil {
		return nil, errors.New("team not found")
	}
	if team.UserID != userID {
		return nil, errors.New("team belongs to another user")
	}

	// Check if team has at least 1 active member
	activeCount := 0
	for _, m := range team.Members {
		if !m.IsBackup {
			activeCount++
		}
	}
	if activeCount == 0 {
		return nil, errors.New("team must have at least 1 active member")
	}

	return &team, nil
}

func generateSeed() string {
	// Placeholder for simple seed generation
	return time.Now().String()
}

// AbandonSession marks a session as ABANDONED when user clicks RUN (Phase 12)
func (s *RaidService) AbandonSession(sessionID, userID uint) error {
	var session models.RaidSession
	err := db.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		return errors.New("session not found or unauthorized")
	}

	if session.Status != "IN_PROGRESS" {
		return errors.New("session is not in progress")
	}

	session.Status = "ABANDONED"
	now := time.Now()
	session.CompletedAt = &now

	return db.DB.Save(&session).Error
}

// GetRaidSessionWithSprites retrieves a raid session with team character sprites preloaded
func (s *RaidService) GetRaidSessionWithSprites(sessionID uint) (*RaidSessionWithSprites, error) {
	var session models.RaidSession
	if err := db.DB.
		Preload("Mission").
		Preload("Team").
		Preload("Team.Members").
		Preload("Team.Members.Character").
		Preload("Team.Members.Character.Moves").
		Preload("Boss").
		First(&session, sessionID).Error; err != nil {
		return nil, err
	}

	var teamChars []models.Character
	if session.Team != nil {
		for _, member := range session.Team.Members {
			teamChars = append(teamChars, member.Character)
		}
	}

	return &RaidSessionWithSprites{
		Session:   &session,
		TeamChars: teamChars,
		Mission:   &session.Mission,
	}, nil
}

// convertMoveToAbility converts a legacy CharacterMove to a standard Ability for calculation
func (s *RaidService) convertMoveToAbility(move models.CharacterMove) *models.Ability {
	dmgType := "physical"
	if move.Category == "SPECIAL" {
		dmgType = "magical"
	}

	return &models.Ability{
		Name:          move.Name,
		BaseDamage:    move.Power,
		DamageType:    dmgType,
		AnimationName: move.Animation,
		// Map other fields as needed
		AppliesBuff:        move.EffectType, // Simplification
		StatusEffectChance: move.EffectChance,
	}
}
