package services

import "math"

// PassiveAbility represents a class-specific passive effect
type PassiveAbility struct {
	Name        string
	Description string
	Class       string
}

// GetPassiveAbility returns the passive ability for a class
func GetPassiveAbility(class string) PassiveAbility {
	passives := map[string]PassiveAbility{
		"Warrior": {
			Name:        "Berserker",
			Description: "+10% ATK when HP < 50%",
			Class:       "Warrior",
		},
		"Mage": {
			Name:        "Mana Surge",
			Description: "+15% damage on critical hits",
			Class:       "Mage",
		},
		"Tank": {
			Name:        "Fortify",
			Description: "-20% damage taken when HP > 70%",
			Class:       "Tank",
		},
		"Archer": {
			Name:        "Precision",
			Description: "+20% critical hit chance",
			Class:       "Archer",
		},
		"Healer": {
			Name:        "Regeneration",
			Description: "Restore 5% HP at turn start",
			Class:       "Healer",
		},
	}

	if passive, ok := passives[class]; ok {
		return passive
	}

	return PassiveAbility{
		Name:        "None",
		Description: "No passive ability",
		Class:       "Unknown",
	}
}

// ApplyPassiveAbility applies passive effects to combat calculations
// Returns multiplier to apply to damage/stats
func ApplyPassiveAbility(class string, currentHP, maxHP int, isCrit bool, isAttacker bool) float64 {
	hpPercent := float64(currentHP) / float64(maxHP) * 100

	switch class {
	case "Warrior":
		// Berserker: +10% ATK when HP < 50%
		if isAttacker && hpPercent < 50 {
			return 1.10
		}

	case "Mage":
		// Mana Surge: +15% damage on critical hits
		if isAttacker && isCrit {
			return 1.15
		}

	case "Tank":
		// Fortify: -20% damage taken when HP > 70%
		if !isAttacker && hpPercent > 70 {
			return 0.80 // Reduce damage taken
		}

	case "Archer":
		// Precision: +20% crit chance (handled in calculateDamage)
		// This function doesn't apply multiplier for Archer
		return 1.0

	case "Healer":
		// Regeneration: Restore HP at turn start (handled separately)
		return 1.0
	}

	return 1.0 // No passive effect
}

// GetArcherCritBonus returns additional crit chance for Archers
func GetArcherCritBonus(class string) float64 {
	if class == "Archer" {
		return 0.20 // +20% crit chance
	}
	return 0.0
}

// ApplyHealerRegeneration restores HP for Healers at turn start
func ApplyHealerRegeneration(class string, currentHP, maxHP int) int {
	if class == "Healer" {
		regenAmount := int(math.Ceil(float64(maxHP) * 0.05)) // 5% max HP
		newHP := currentHP + regenAmount
		if newHP > maxHP {
			newHP = maxHP
		}
		return newHP
	}
	return currentHP
}
