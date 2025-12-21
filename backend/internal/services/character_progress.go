package services

import (
	"math"

	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// CharacterProgressInfo contains XP and level progression data for frontend display
type CharacterProgressInfo struct {
	CurrentLevel   int     `json:"current_level"`
	CurrentXP      int     `json:"current_xp"`        // XP in current level
	XPForNextLevel int     `json:"xp_for_next_level"` // XP needed to reach next level
	XPProgressPct  float64 `json:"xp_progress_pct"`   // Percentage to next level (0-100)
	TotalXP        int     `json:"total_xp"`          // Lifetime XP earned
	NextLevel      int     `json:"next_level"`
}

// GetCharacterProgressInfo returns detailed XP progression info for a character
func GetCharacterProgressInfo(char *models.Character) CharacterProgressInfo {
	currentXP := char.Experience
	currentLevel := char.Level
	xpForNext := getExpForNextLevel(currentLevel)

	progress := float64(0)
	if xpForNext > 0 {
		progress = (float64(currentXP) / float64(xpForNext)) * 100.0
	}

	return CharacterProgressInfo{
		CurrentLevel:   currentLevel,
		CurrentXP:      currentXP,
		XPForNextLevel: xpForNext,
		XPProgressPct:  math.Round(progress*100) / 100, // Round to 2 decimals
		TotalXP:        char.TotalXP,
		NextLevel:      currentLevel + 1,
	}
}

// getExpForNextLevel calculates XP required for next level (Pokemon Medium-Fast curve)
// Made standalone to be reused by both raid_leveling.go and this file
func getExpForNextLevel(currentLevel int) int {
	// Medium-Fast growth: Level³
	// Level 1→2: 8 XP
	// Level 5→6: 216 XP
	// Level 10→11: 1331 XP
	nextLevel := currentLevel + 1
	return int(math.Pow(float64(nextLevel), 3))
}
