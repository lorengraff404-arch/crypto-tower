package services

// GetClassAdvantage returns damage multiplier for class matchup (Phase 15.1)
// Similar to type effectiveness but for character classes
func GetClassAdvantage(attackerClass, defenderClass string) float64 {
	// Class advantage chart (rock-paper-scissors style)
	advantages := map[string]map[string]float64{
		"Warrior": {
			"Mage":   1.3, // Warriors overwhelm low-defense mages
			"Tank":   0.7, // Warriors can't penetrate tank armor
			"Healer": 1.2, // Warriors pressure healers
			"Archer": 1.0, // Neutral
		},
		"Mage": {
			"Warrior": 0.7, // Mages have low physical defense
			"Tank":    1.4, // Magic penetrates armor
			"Archer":  1.1, // Slight advantage
			"Healer":  1.0, // Neutral
		},
		"Tank": {
			"Warrior": 1.3, // Armor absorbs physical damage
			"Mage":    0.6, // Vulnerable to magic
			"Archer":  0.8, // Arrows partially blocked
			"Healer":  1.0, // Neutral
		},
		"Archer": {
			"Warrior": 1.0, // Neutral
			"Mage":    1.1, // Slight advantage
			"Tank":    0.7, // Armor blocks arrows
			"Healer":  1.4, // Interrupts healing/casting
		},
		"Healer": {
			"Warrior": 0.8, // Vulnerable to aggression
			"Mage":    1.0, // Neutral
			"Tank":    1.0, // Neutral
			"Archer":  0.6, // Interrupted by arrows
		},
	}

	// Lookup advantage
	if classMap, ok := advantages[attackerClass]; ok {
		if mult, ok := classMap[defenderClass]; ok {
			return mult
		}
	}

	// Default: neutral (1.0x)
	return 1.0
}

// GetClassDescription returns flavor text for class
func GetClassDescription(class string) string {
	descriptions := map[string]string{
		"Warrior": "Physical DPS with high attack. Strong vs Mages, weak vs Tanks.",
		"Mage":    "Magic DPS with burst damage. Strong vs Tanks, weak vs Warriors.",
		"Tank":    "Defender with high HP/DEF. Strong vs Warriors, weak vs Mages.",
		"Archer":  "Ranged DPS with high speed. Strong vs Healers, weak vs Tanks.",
		"Healer":  "Support with sustain. Strong vs Warriors, weak vs Archers.",
	}

	if desc, ok := descriptions[class]; ok {
		return desc
	}
	return "Unknown class"
}

// GetClassEmoji returns emoji representation
func GetClassEmoji(class string) string {
	emojis := map[string]string{
		"Warrior": "⚔️",
		"Mage":    "🔮",
		"Tank":    "🛡️",
		"Archer":  "🏹",
		"Healer":  "✨",
	}

	if emoji, ok := emojis[class]; ok {
		return emoji
	}
	return "🎮"
}
