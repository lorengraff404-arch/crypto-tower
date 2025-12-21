package formulas

// TypeResistance represents the 8x8 Type vs Element resistance matrix
// Returns damage multiplier: 2.0 = Super Effective, 1.0 = Normal, 0.5 = Not Very Effective

var TypeElementMatrix = map[string]map[string]float64{
	// BEAST type resistances
	"BEAST": {
		"FIRE":    0.5, // Weak to fire
		"WATER":   1.0,
		"ICE":     1.0,
		"THUNDER": 1.0,
		"DARK":    1.5, // Resistant to dark
		"PLANT":   1.5, // Resistant to plant
		"EARTH":   1.0,
		"WIND":    1.0,
	},
	// DRAGON type resistances
	"DRAGON": {
		"FIRE":    1.5, // Resistant to fire
		"WATER":   0.5, // Weak to water
		"ICE":     0.5, // Weak to ice
		"THUNDER": 2.0, // Very weak to thunder
		"DARK":    1.0,
		"PLANT":   1.0,
		"EARTH":   1.0,
		"WIND":    1.0,
	},
	// INSECT type resistances
	"INSECT": {
		"FIRE":    2.0, // Very weak to fire
		"WATER":   1.0,
		"ICE":     1.5, // Resistant to ice
		"THUNDER": 1.0,
		"DARK":    1.0,
		"PLANT":   0.5, // Weak to plant
		"EARTH":   1.5, // Resistant to earth
		"WIND":    0.5, // Weak to wind
	},
	// MINERAL type resistances
	"MINERAL": {
		"FIRE":    1.5, // Resistant to fire
		"WATER":   0.5, // Weak to water
		"ICE":     1.5, // Resistant to ice
		"THUNDER": 1.0,
		"DARK":    1.0,
		"PLANT":   2.0, // Very weak to plant
		"EARTH":   1.5, // Resistant to earth
		"WIND":    1.0,
	},
	// SPIRIT type resistances
	"SPIRIT": {
		"FIRE":    1.0,
		"WATER":   1.0,
		"ICE":     1.0,
		"THUNDER": 1.0,
		"DARK":    2.0, // Very weak to dark
		"PLANT":   1.0,
		"EARTH":   0.5, // Weak to earth
		"WIND":    1.5, // Resistant to wind
	},
	// AVIAN type resistances
	"AVIAN": {
		"FIRE":    1.0,
		"WATER":   1.0,
		"ICE":     2.0, // Very weak to ice
		"THUNDER": 2.0, // Very weak to thunder
		"DARK":    1.0,
		"PLANT":   1.5, // Resistant to plant
		"EARTH":   0.5, // Weak to earth
		"WIND":    0.5, // Weak to wind
	},
	// AQUA type resistances
	"AQUA": {
		"FIRE":    0.5, // Weak to fire
		"WATER":   2.0, // Very resistant to water
		"ICE":     1.5, // Resistant to ice
		"THUNDER": 0.5, // Weak to thunder
		"DARK":    1.0,
		"PLANT":   0.5, // Weak to plant
		"EARTH":   1.0,
		"WIND":    1.0,
	},
	// FLORA type resistances
	"FLORA": {
		"FIRE":    2.0, // Very weak to fire
		"WATER":   1.5, // Resistant to water
		"ICE":     2.0, // Very weak to ice
		"THUNDER": 1.0,
		"DARK":    1.0,
		"PLANT":   1.5, // Resistant to plant
		"EARTH":   1.5, // Resistant to earth
		"WIND":    0.5, // Weak to wind
	},
}

// GetTypeResistance returns the damage multiplier for a Type vs Element matchup
func GetTypeResistance(characterType string, attackElement string) float64 {
	// Default to neutral if type or element not found
	if typeResistances, ok := TypeElementMatrix[characterType]; ok {
		if resistance, ok := typeResistances[attackElement]; ok {
			return resistance
		}
	}
	return 1.0 // Neutral damage
}

// GetEffectivenessText returns user-friendly text for the effectiveness
func GetEffectivenessText(multiplier float64) string {
	if multiplier >= 2.0 {
		return "Super Effective!"
	} else if multiplier >= 1.5 {
		return "It's effective!"
	} else if multiplier <= 0.5 {
		return "Not very effective..."
	}
	return ""
}

// ValidateCharacterType checks if a character type is valid
func ValidateCharacterType(characterType string) bool {
	validTypes := []string{"BEAST", "DRAGON", "INSECT", "MINERAL", "SPIRIT", "AVIAN", "AQUA", "FLORA"}
	for _, t := range validTypes {
		if t == characterType {
			return true
		}
	}
	return false
}

// ValidateElement checks if an element is valid
func ValidateElement(element string) bool {
	validElements := []string{"FIRE", "WATER", "ICE", "THUNDER", "DARK", "PLANT", "EARTH", "WIND"}
	for _, e := range validElements {
		if e == element {
			return true
		}
	}
	return false
}

// GetTypeEmoji returns emoji representation for character types
func GetTypeEmoji(characterType string) string {
	emojis := map[string]string{
		"BEAST":   "🦁",
		"DRAGON":  "🐉",
		"INSECT":  "🐛",
		"MINERAL": "💎",
		"SPIRIT":  "👻",
		"AVIAN":   "🦅",
		"AQUA":    "🐠",
		"FLORA":   "🌿",
	}
	if emoji, ok := emojis[characterType]; ok {
		return emoji
	}
	return "✨"
}

// GetElementEmoji returns emoji representation for elements
func GetElementEmoji(element string) string {
	emojis := map[string]string{
		"FIRE":    "🔥",
		"WATER":   "💧",
		"ICE":     "❄️",
		"THUNDER": "⚡",
		"DARK":    "🌑",
		"PLANT":   "🌱",
		"EARTH":   "🪨",
		"WIND":    "💨",
	}
	if emoji, ok := emojis[element]; ok {
		return emoji
	}
	return "✨"
}
