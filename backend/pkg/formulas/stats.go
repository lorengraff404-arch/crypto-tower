package formulas

import "math"

// GetRarityMultiplier returns stat multiplier for rarity tier
func GetRarityMultiplier(rarity string) float64 {
	switch rarity {
	case "SSS":
		return 4.0
	case "SS":
		return 3.4
	case "S":
		return 2.8
	case "A":
		return 2.2
	case "B":
		return 1.6
	case "C":
		return 1.0
	default:
		return 1.0
	}
}

// GetRarityForLevel returns the rarity tier for a given level
func GetRarityForLevel(level int) string {
	switch {
	case level >= 95:
		return "SSS"
	case level >= 80:
		return "SS"
	case level >= 60:
		return "S"
	case level >= 40:
		return "A"
	case level >= 20:
		return "B"
	default:
		return "C"
	}
}

// GetEvolutionStage returns evolution stage for level
func GetEvolutionStage(level int) int {
	switch {
	case level >= 75:
		return 3 // Ultimate
	case level >= 50:
		return 2 // Mature
	case level >= 25:
		return 1 // Growth
	default:
		return 0 // Base
	}
}

// GetEvolutionBonus returns stat multiplier for evolution stage
func GetEvolutionBonus(stage int) float64 {
	return 1.0 + (float64(stage) * 0.1)
}

// CalculateStat computes final stat value with all bonuses
func CalculateStat(baseStat int, level int, rarity string, evolutionStage int) int {
	rarityMultiplier := GetRarityMultiplier(rarity)
	evolutionBonus := GetEvolutionBonus(evolutionStage)
	levelScaling := 1.0 + (float64(level) * 0.05) // +5% per level
	
	finalStat := float64(baseStat) * rarityMultiplier * evolutionBonus * levelScaling
	return int(math.Floor(finalStat))
}

// RecalculateAllStats updates all character stats based on current level/rarity/evolution
func RecalculateAllStats(baseAttack, baseDefense, baseHP, baseSpeed int, level int, rarity string, evolutionStage int) (int, int, int, int) {
	return CalculateStat(baseAttack, level, rarity, evolutionStage),
		CalculateStat(baseDefense, level, rarity, evolutionStage),
		CalculateStat(baseHP, level, rarity, evolutionStage),
		CalculateStat(baseSpeed, level, rarity, evolutionStage)
}

// ValidateStats checks if character stats match expected values
func ValidateStats(baseAttack, currentAttack int, level int, rarity string, evolutionStage int) bool {
	expectedAttack := CalculateStat(baseAttack, level, rarity, evolutionStage)
	// Allow 1 point of difference for rounding
	return abs(currentAttack-expectedAttack) <= 1
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
