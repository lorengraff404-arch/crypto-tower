package services

import (
	"encoding/json"
	"math/rand"
)

// StatusEffect represents an active buff or debuff
type StatusEffect struct {
	Effect    string `json:"effect"`     // "burn", "amped", "stun", etc.
	TurnsLeft int    `json:"turns_left"` // Decrements each turn
	Stacks    int    `json:"stacks"`     // For stacking effects
}

// StatusEffectManager handles application and processing of status effects
type StatusEffectManager struct {
	Effects []StatusEffect
}

// NewStatusEffectManager creates manager from JSON string
func NewStatusEffectManager(jsonData string) *StatusEffectManager {
	sem := &StatusEffectManager{Effects: []StatusEffect{}}
	if jsonData != "" {
		json.Unmarshal([]byte(jsonData), &sem.Effects)
	}
	return sem
}

// ToJSON serializes effects back to JSON string
func (sem *StatusEffectManager) ToJSON() string {
	data, _ := json.Marshal(sem.Effects)
	return string(data)
}

// AddEffect adds or refreshes a status effect
func (sem *StatusEffectManager) AddEffect(effect string, duration int) {
	// Check if effect already exists
	for i, e := range sem.Effects {
		if e.Effect == effect {
			// Refresh duration to max
			if duration > e.TurnsLeft {
				sem.Effects[i].TurnsLeft = duration
			}
			return
		}
	}
	// Add new effect
	sem.Effects = append(sem.Effects, StatusEffect{
		Effect:    effect,
		TurnsLeft: duration,
		Stacks:    1,
	})
}

// HasEffect checks if an effect is active
func (sem *StatusEffectManager) HasEffect(effect string) bool {
	for _, e := range sem.Effects {
		if e.Effect == effect && e.TurnsLeft > 0 {
			return true
		}
	}
	return false
}

// ProcessTurnEffects applies DoT, decrements durations, returns damage dealt
func (sem *StatusEffectManager) ProcessTurnEffects(maxHP int64) int64 {
	var dotDamage int64 = 0

	// Process Damage over Time
	if sem.HasEffect("burn") {
		dotDamage += int64(float64(maxHP) * 0.05) // 5% max HP
	}
	if sem.HasEffect("poison") {
		dotDamage += int64(float64(maxHP) * 0.08) // 8% max HP
	}
	if sem.HasEffect("bleed") {
		dotDamage += int64(float64(maxHP) * 0.03) // 3% max HP
	}

	// Decrement all durations
	newEffects := []StatusEffect{}
	for _, e := range sem.Effects {
		e.TurnsLeft--
		if e.TurnsLeft > 0 {
			newEffects = append(newEffects, e)
		}
	}
	sem.Effects = newEffects

	return dotDamage
}

// CanAct checks if unit can take action (stun/sleep/freeze check)
func (sem *StatusEffectManager) CanAct() bool {
	if sem.HasEffect("stun") {
		return false
	}
	if sem.HasEffect("sleep") {
		return false
	}
	if sem.HasEffect("freeze") {
		return false
	}
	if sem.HasEffect("paralyze") {
		// 25% chance to skip
		return rand.Float64() > 0.25
	}
	return true
}

// GetStatModifier returns multiplier for a stat (atk, def, spd)
func (sem *StatusEffectManager) GetStatModifier(stat string) float64 {
	modifier := 1.0

	switch stat {
	case "atk":
		if sem.HasEffect("amped") {
			modifier *= 1.3
		}
		if sem.HasEffect("feeble") {
			modifier *= 0.7
		}
	case "def":
		if sem.HasEffect("bulked") {
			modifier *= 1.3
		}
		if sem.HasEffect("fragile") {
			modifier *= 0.7
		}
	case "spd":
		if sem.HasEffect("haste") {
			modifier *= 1.3
		}
		if sem.HasEffect("slow") {
			modifier *= 0.7
		}
	}

	return modifier
}

// GetAccuracyModifier returns hit chance modifier
func (sem *StatusEffectManager) GetAccuracyModifier() float64 {
	if sem.HasEffect("blind") {
		return 0.5 // 50% hit chance
	}
	return 1.0
}

// WakeFromSleep removes sleep on damage taken
func (sem *StatusEffectManager) WakeFromSleep() {
	newEffects := []StatusEffect{}
	for _, e := range sem.Effects {
		if e.Effect != "sleep" {
			newEffects = append(newEffects, e)
		}
	}
	sem.Effects = newEffects
}

// GetStatusEffectMultipliers returns ATK, DEF, SPD, and CRIT multipliers from active effects (Phase 15.3)
func GetStatusEffectMultipliers(effects []StatusEffect) (atkMult, defMult, spdMult, critMult float64) {
	atkMult, defMult, spdMult, critMult = 1.0, 1.0, 1.0, 0.0 // crit is additive

	for _, effect := range effects {
		switch effect.Effect {
		// Buffs
		case "amped", "atk_up":
			atkMult *= 1.3
		case "bulked", "def_up":
			defMult *= 1.3
		case "haste", "spd_up":
			spdMult *= 1.5
		case "crit_boost":
			critMult += 0.3

		// Debuffs
		case "feeble", "atk_down":
			atkMult *= 0.7
		case "fragile", "def_down":
			defMult *= 0.7
		case "slow":
			spdMult *= 0.7
		}
	}

	return
}
