package services

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// BattleEngine handles core battle mechanics
type BattleEngine struct {
	config *ConfigService
}

// NewBattleEngine creates a new battle engine
func NewBattleEngine() *BattleEngine {
	return &BattleEngine{
		config: GetConfigService(),
	}
}

// CalculateTurnOrder determines action order based on Speed
func (e *BattleEngine) CalculateTurnOrder(participants []models.BattleParticipant) []models.TurnEntry {
	var queue []models.TurnEntry

	for _, p := range participants {
		if !p.IsFainted && p.IsActive {
			queue = append(queue, models.TurnEntry{
				ParticipantID: p.CharacterID,
				Speed:         p.Speed,
				Priority:      0, // Default priority
			})
		}
	}

	// Sort by priority (desc), then speed (desc), then random
	sort.Slice(queue, func(i, j int) bool {
		if queue[i].Priority != queue[j].Priority {
			return queue[i].Priority > queue[j].Priority
		}
		if queue[i].Speed != queue[j].Speed {
			return queue[i].Speed > queue[j].Speed
		}
		return rand.Intn(2) == 0 // Random tiebreaker
	})

	return queue
}

// CalculateDamage computes damage with type advantages, stats, and randomness
func (e *BattleEngine) CalculateDamage(
	attacker models.BattleParticipant,
	defender models.BattleParticipant,
	ability models.Ability,
) (damage int, critical bool, effectiveness string) {

	// Base damage from ability
	baseDamage := float64(ability.Damage)

	// Apply attacker's attack stat
	attackMultiplier := float64(attacker.Attack) / 100.0
	damage = int(baseDamage * attackMultiplier)

	// Apply defender's defense
	defCap := e.config.GetFloat("battle_def_reduction_cap", 0.75)
	defenseReduction := float64(defender.Defense) / 200.0
	if defenseReduction > defCap {
		defenseReduction = defCap // Use dynamic cap
	}
	damage = int(float64(damage) * (1.0 - defenseReduction))

	// Full Type Chart Implementation from Frontend (type-chart.html)
	// Types: BEAST, DRAGON, INSECT, MINERAL, SPIRIT, AVIAN, AQUA, FLORA
	// Elements: FIRE, WATER, ICE, THUNDER, DARK, PLANT, EARTH, WIND

	// Map Ability Element -> Map[Defender Type] -> Multiplier
	// Note: Defender has a "Type" (e.g. BEAST) but Ability has an "Element" (e.g. FIRE).
	// The frontend chart defines Resistances per Type against Elements.
	// We verify if the Defender.Element (which stores Type like BEAST) matches.

	typeEffectiveness := map[string]map[string]float64{
		"BEAST":   {"FIRE": 0.5, "DARK": 1.5, "PLANT": 1.5},
		"DRAGON":  {"FIRE": 1.5, "WATER": 0.5, "ICE": 0.5, "THUNDER": 0.5},
		"INSECT":  {"FIRE": 0.5, "ICE": 1.5, "PLANT": 2.0, "EARTH": 1.5, "WIND": 2.0},
		"MINERAL": {"FIRE": 1.5, "WATER": 2.0, "ICE": 1.5, "PLANT": 0.5, "EARTH": 1.5},
		"SPIRIT":  {"DARK": 0.5, "EARTH": 2.0, "WIND": 1.5},
		"AVIAN":   {"ICE": 0.5, "THUNDER": 0.5, "PLANT": 1.5, "EARTH": 2.0, "WIND": 2.0},
		"AQUA":    {"FIRE": 2.0, "WATER": 0.5, "ICE": 1.5, "THUNDER": 2.0, "PLANT": 2.0},
		"FLORA":   {"FIRE": 0.5, "WATER": 1.5, "ICE": 0.5, "PLANT": 1.5, "EARTH": 1.5, "WIND": 2.0},
	}

	effectiveness = "normal"
	effectivenessMultiplier := 1.0

	// Defender.Element actually stores the Character Type (BEAST, DRAGON, etc.) based on our models
	// Ability.Element stores the Attack Element (FIRE, WATER, etc.)

	if defenderTypeChart, ok := typeEffectiveness[defender.Element]; ok {
		// Check how this Type resists the incoming Element
		// The chart lists RESISTANCES/WEAKNESSES directly
		// > 1.0 = Weakness (Takes more damage) ? No, chart says "Resistant (1.5x)"... wait.

		// Let's re-read the HTML logic carefully.
		// value >= 1.5 ? 'Resistance' (Green) ... wait. usually Red is weak.
		// Frontend JS:
		// value >= 1.5 ? 'weak' (Green bg? No rgba(16, 185, 129) is Green/Success).
		// value <= 0.5 ? 'super-effective' (Red bg).

		// In HTML:
		// value >= 1.5 ? 'resistance-item weak' ... class name is confusing.
		// "Resistant to (Takes less damage)" is usually < 1.0
		// But the frontend text says: "Resistant (1.5x)" ?? That implies TAKING MORE damage?
		// OR does it mean "Effectiveness on Defender"?

		// Let's look at standard gaming logic vs this chart's text.
		// BEAST: FIRE 0.5. (Red badge "Weak (0.5x)" in HTML?)
		// HTML Code: value <= 0.5 ? 'resistance-item super-effective' -> class 'super-effective' usually means attack is strong.
		// Logic: If Attack is Super Effective, Multiplier should be > 1.0 (e.g. 2.0).
		// The HTML data says BEAST: FIRE = 0.5. And labels it "Weak (0.5x)" or "Super Effective".
		// This suggests the value stored is the DEFENSE MULTIPLIER? or RESISTANCE?
		// If BEAST takes 0.5x damage from FIRE, it is RESISTANT.
		// If BEAST takes 2.0x damage, it is WEAK.

		// Let's check a specific entry:
		// DRAGON: ICE 0.5. Dragons usually weak to Ice? Or Resist?
		// AQUA (Fish): THUNDER 2.0. Fish usually WEAK to Thunder.
		// So 2.0 = Weak (Takes Double Damage).
		// BEAST: FIRE 0.5. Beasts usually afraid of fire? Or resist?
		// FLORA (Plant): FIRE 0.5. Plants BURN. So they should take 2.0x.
		// BUT the data says FLORA: FIRE = 0.5.
		// Wait. Plants are WEAK to Fire. So Fire should do MORE damage.
		// If the value is 0.5 ... maybe it means "Resistance is 0.5" (Low resistance = High damage)?
		// OR maybe the chart keys are swapped?

		// Let's look at FLORA: WATER 1.5. Plants absorb water. Restistant?
		// If 1.5 means Resistant... then 1.5x damage is BAD for defender? No that's contradictory.

		// Let's assume standard Pokemon logic:
		// Grass (Flora) vs Fire: Weak (2x damage). Data says 0.5.
		// Grass (Flora) vs Water: Resist (0.5x damage). Data says 1.5.

		// HYPOTHESIS: The values in `typeData` are INVERTED or represent "Safety".
		// OR: My reading of the HTML text is crucial.
		// HTML: value >= 1.5 ? 'weak' class. Label: "Resistant (1.5x)".
		// This is confusing phrasing. "Resistant" usually means good for defender.
		// But "Weak" class usually means bad for defender.

		// Let's stick to the visual cues in `type-chart.html`:
		// .cell-super { background: red; color: light red; } -> typically "Super Effective" damage on me (OUCH).
		// .cell-resistant { background: green; color: light green; } -> typically "Resistant" (GOOD).

		// JS Logic for Class:
		// value >= 1.5 -> 'cell-resistant' (Green).
		// value <= 0.5 -> 'cell-super' (Red).

		// So:
		// 1.5+ = Green = Resistant (Good for Defender). Means Multiplier should be < 1.0 (e.g. 0.5).
		// 0.5- = Red = Super Effective (Bad for Defender). Means Multiplier should be > 1.0 (e.g. 2.0).

		// RE-MAPPING DATA TO IN-GAME MULTIPLIERS:
		// Data: 0.5 -> Real Multiplier: 2.0 (Super Effective)
		// Data: 1.5/2.0 -> Real Multiplier: 0.5 (Resistant)
		// Data: 1.0 -> Real Multiplier: 1.0 (Neutral)

		// Let's verify with FLORA (Plant).
		// Data: FIRE = 0.5. Real Multiplier = 2.0. Logic: Power of Fire vs Plant = 2x. CORRECT.
		// Data: WATER = 1.5. Real Multiplier = 0.5. Logic: Power of Water vs Plant = 0.5x. CORRECT.

		// Conclusion: The numbers in `typeData` are "Resistance Scores" where Lower = Worse for Defender (Higher Dmg).

		if val, found := defenderTypeChart[ability.Element]; found {
			if val <= 0.5 {
				effectiveness = "super_effective"
				effectivenessMultiplier = 2.0
			} else if val >= 1.5 {
				effectiveness = "not_very_effective"
				effectivenessMultiplier = 0.5
			}
		}
	}

	damage = int(float64(damage) * effectivenessMultiplier)

	// Apply buffs
	for _, buff := range attacker.Buffs {
		if buff.Stat == "attack" {
			damage = int(float64(damage) * buff.Modifier)
		}
	}

	for _, buff := range defender.Buffs {
		if buff.Stat == "defense" {
			damage = int(float64(damage) / buff.Modifier)
		}
	}

	// Critical hit chance (10% base default)
	critChance := e.config.GetFloat("battle_crit_chance", 0.10)
	if rand.Float64() < critChance {
		critical = true
		critMult := e.config.GetFloat("battle_crit_multiplier", 1.5)
		damage = int(float64(damage) * critMult)
	}

	// Randomness (±10% base default)
	randFactor := e.config.GetFloat("battle_randomness_factor", 0.10)
	randomFactor := (1.0 - randFactor) + (rand.Float64() * randFactor * 2.0)
	damage = int(float64(damage) * randomFactor)

	// Minimum damage
	if damage < 1 {
		damage = 1
	}

	return damage, critical, effectiveness
}

// ExecuteAbility executes an ability and returns the result
func (e *BattleEngine) ExecuteAbility(
	attacker *models.BattleParticipant,
	defender *models.BattleParticipant,
	ability models.Ability,
) (*models.TurnResult, error) {

	// Check mana cost
	if attacker.CurrentMana < ability.ManaCost {
		return nil, errors.New("not enough mana")
	}

	// Deduct mana
	attacker.CurrentMana -= ability.ManaCost

	// Calculate damage
	damage, critical, effectiveness := e.CalculateDamage(*attacker, *defender, ability)

	// Apply damage
	defender.CurrentHP -= damage
	if defender.CurrentHP < 0 {
		defender.CurrentHP = 0
	}

	// Check if fainted
	fainted := defender.CurrentHP == 0
	if fainted {
		defender.IsFainted = true
	}

	// Build message
	message := fmt.Sprintf("%s used %s! ", "Attacker", ability.Name)
	if critical {
		message += "Critical Hit! "
	}
	if effectiveness == "super_effective" {
		message += "It's super effective! "
	} else if effectiveness == "not_very_effective" {
		message += "It's not very effective... "
	}
	message += fmt.Sprintf("Dealt %d damage!", damage)

	return &models.TurnResult{
		AttackerID:      attacker.CharacterID,
		DefenderID:      defender.CharacterID,
		AbilityUsed:     ability.Name,
		Damage:          damage,
		Critical:        critical,
		Effectiveness:   effectiveness,
		DefenderHP:      defender.CurrentHP,
		DefenderFainted: fainted,
		Message:         message,
	}, nil
}

// CheckBattleEnd determines if battle is over
func (e *BattleEngine) CheckBattleEnd(team1, team2 []models.BattleParticipant) (ended bool, winnerTeam int) {
	team1Alive := false
	team2Alive := false

	for _, p := range team1 {
		if !p.IsFainted {
			team1Alive = true
			break
		}
	}

	for _, p := range team2 {
		if !p.IsFainted {
			team2Alive = true
			break
		}
	}

	if !team1Alive {
		return true, 2
	}
	if !team2Alive {
		return true, 1
	}

	return false, 0
}

// RegenerateMana applies mana regeneration
func (e *BattleEngine) RegenerateMana(participant *models.BattleParticipant, regenRate int) {
	participant.CurrentMana += regenRate
	if participant.CurrentMana > participant.MaxMana {
		participant.CurrentMana = participant.MaxMana
	}
}

// ApplyStatusDamage applies damage from status effects like burn/poison
func (e *BattleEngine) ApplyStatusDamage(participant *models.BattleParticipant) int {
	damage := 0
	switch participant.StatusEffect {
	case "burn":
		damage = participant.MaxHP / 16 // 1/16 max HP
	case "poison":
		damage = participant.MaxHP / 8 // 1/8 max HP
	}

	if damage > 0 {
		participant.CurrentHP -= damage
		if participant.CurrentHP < 0 {
			participant.CurrentHP = 0
			participant.IsFainted = true
		}
	}

	return damage
}

// UpdateBuffs decrements buff duration and removes expired buffs
func (e *BattleEngine) UpdateBuffs(participant *models.BattleParticipant) {
	// Update buffs
	var activeBuffs []models.Buff
	for _, buff := range participant.Buffs {
		buff.TurnsRemaining--
		if buff.TurnsRemaining > 0 {
			activeBuffs = append(activeBuffs, buff)
		}
	}
	participant.Buffs = activeBuffs

	// Update debuffs
	var activeDebuffs []models.Buff
	for _, debuff := range participant.Debuffs {
		debuff.TurnsRemaining--
		if debuff.TurnsRemaining > 0 {
			activeDebuffs = append(activeDebuffs, debuff)
		}
	}
	participant.Debuffs = activeDebuffs
}

// ExecuteRaidTurn executes a turn in a raid battle
func (e *BattleEngine) ExecuteRaidTurn(sessionID uint, characterID uint, actionType string, abilityID *uint, targetID *uint, itemID *uint) (map[string]interface{}, error) {
	// Placeholder for now - returns simple success
	return map[string]interface{}{
		"success":     true,
		"message":     "Turn  executed successfully",
		"action_type": actionType,
	}, nil
}
