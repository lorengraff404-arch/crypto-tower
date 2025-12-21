package services

import (

	"errors"
	"fmt"
	"math/rand"
	"sort"


	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// BattleEngine handles core battle mechanics
type BattleEngine struct{}

// NewBattleEngine creates a new battle engine
func NewBattleEngine() *BattleEngine {
	return &BattleEngine{}
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
	defenseReduction := float64(defender.Defense) / 200.0
	if defenseReduction > 0.75 {
		defenseReduction = 0.75 // Cap at 75% reduction
	}
	damage = int(float64(damage) * (1.0 - defenseReduction))

	// Type effectiveness (simplified - just fire/water/grass triangle)
	effectiveness = "normal"
	effectivenessMultiplier := 1.0

	// Simplified type chart
	if ability.Element == "Fire" && defender.Attack > 0 { // Placeholder for element
		// Would check defender.Element here
		effectivenessMultiplier = 1.0 // Normal for now
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

	// Critical hit chance (10% base)
	if rand.Intn(100) < 10 {
		critical = true
		damage = int(float64(damage) * 1.5)
	}

	// Randomness ±10%
	randomFactor := 0.9 + (rand.Float64() * 0.2)
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
		"success": true,
		"message": "Turn  executed successfully",
		"action_type": actionType,
	}, nil
}
