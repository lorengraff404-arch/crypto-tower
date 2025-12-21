package formulas

import "math"

// GetXPForLevel returns the total XP required to reach a specific level
// Uses exponential curve: XP = 100 * level^2.5
func GetXPForLevel(level int) int {
	if level <= 1 {
		return 0
	}
	return int(math.Floor(100 * math.Pow(float64(level), 2.5)))
}

// GetLevelFromXP returns the level for a given total XP amount
func GetLevelFromXP(totalXP int) int {
	level := 1
	for level < 100 {
		requiredXP := GetXPForLevel(level + 1)
		if totalXP < requiredXP {
			break
		}
		level++
	}
	return level
}

// GetXPForNextLevel returns XP needed for next level from current XP
func GetXPForNextLevel(currentLevel int, totalXP int) int {
	if currentLevel >= 100 {
		return 0
	}
	nextLevelXP := GetXPForLevel(currentLevel + 1)
	return nextLevelXP - totalXP
}

// GetXPProgressPercent returns progress to next level as percentage
func GetXPProgressPercent(currentLevel int, totalXP int) float64 {
	if currentLevel >= 100 {
		return 100.0
	}
	
	currentLevelXP := GetXPForLevel(currentLevel)
	nextLevelXP := GetXPForLevel(currentLevel + 1)
	
	if nextLevelXP == currentLevelXP {
		return 100.0
	}
	
	progress := float64(totalXP-currentLevelXP) / float64(nextLevelXP-currentLevelXP) * 100
	return math.Min(100.0, math.Max(0.0, progress))
}

// ValidateXPGain checks if XP gain is reasonable to prevent exploits
func ValidateXPGain(source string, difficulty string, xpGained int) bool {
	maxAllowed := 0
	
	switch source {
	case "island_raid":
		switch difficulty {
		case "beginner":
			maxAllowed = 200
		case "advanced":
			maxAllowed = 500
		case "expert":
			maxAllowed = 1000
		case "legendary":
			maxAllowed = 2000
		}
	case "pvp_battle":
		maxAllowed = 1500 // Based on level difference
	case "quest":
		maxAllowed = 2500
	case "daily_bonus":
		maxAllowed = 500
	default:
		return false
	}
	
	return xpGained > 0 && xpGained <= maxAllowed
}
