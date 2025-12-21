package services

import (
	"math"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// XP & Leveling System (Phase 10.1)

// distributeExpToTeam splits XP among all active team members
func (s *RaidService) distributeExpToTeam(session *models.RaidSession, totalXP int) {
	// Get active members (non-backup)
	activeMembers := []models.TeamMember{}
	for _, m := range session.Team.Members {
		if !m.IsBackup {
			activeMembers = append(activeMembers, m)
		}
	}

	if len(activeMembers) == 0 {
		return
	}

	// Split XP evenly
	xpPerMember := totalXP / len(activeMembers)

	for _, member := range activeMembers {
		// Load full character data
		var char models.Character
		if err := db.DB.First(&char, member.CharacterID).Error; err != nil {
			continue
		}

		// Add XP
		char.Experience += xpPerMember

		// Check for level ups (can level up multiple times if enough XP)
		leveledUp := false
		for char.Experience >= s.getExpForNextLevel(char.Level) {
			s.levelUpCharacter(&char)
			leveledUp = true
		}

		// Save character
		db.DB.Save(&char)

		// Log level up (optional: could add to battle log)
		if leveledUp {
			// Could emit event or add to session log
			// For now, just save the updated character
		}
	}
}

// getExpForNextLevel calculates XP required for next level (Pokemon Medium-Fast curve)
func (s *RaidService) getExpForNextLevel(currentLevel int) int {
	// Medium-Fast growth: Level³
	// Level 1→2: 8 XP
	// Level 5→6: 216 XP
	// Level 10→11: 1331 XP
	nextLevel := currentLevel + 1
	return int(math.Pow(float64(nextLevel), 3))
}

// levelUpCharacter increases level and stats
func (s *RaidService) levelUpCharacter(char *models.Character) {
	// Deduct XP for level up
	char.Experience -= s.getExpForNextLevel(char.Level)
	char.Level++

	// Stat increases (balanced progression)
	// HP: +5 base + bonus every 5 levels
	hpGain := 5
	if char.Level%5 == 0 {
		hpGain += char.Level / 5
	}
	char.CurrentHP += hpGain // Increase current HP (also heals on level up)

	// Attack: +2 base + bonus every 10 levels
	atkGain := 2
	if char.Level%10 == 0 {
		atkGain += char.Level / 10
	}
	char.CurrentAttack += atkGain

	// Defense: +2 base + bonus every 10 levels
	defGain := 2
	if char.Level%10 == 0 {
		defGain += char.Level / 10
	}
	char.CurrentDefense += defGain

	// Speed: +1 every level
	char.CurrentSpeed += 1

	// Special stats (if they exist in your model)
	// char.SpecialAttack += 1
	// char.SpecialDefense += 1

	// Check for evolution at key levels (16, 32, 48)
	if char.Level == 16 || char.Level == 32 || char.Level == 48 {
		s.checkEvolution(char)
	}
}

// checkEvolution checks if character should evolve
func (s *RaidService) checkEvolution(char *models.Character) {
	// Evolution logic
	// For now, just a placeholder
	// In full implementation:
	// - Check if character has an evolution form
	// - Update CharacterType, stats, potentially moves
	// - Update image/appearance

	// Example:
	// if char.CharacterType == "BASIC_FIRE" && char.Level >= 16 {
	//     char.CharacterType = "EVOLVED_FIRE"
	//     char.MaxHP += 20
	//     char.CurrentAttack += 10
	//     // etc.
	// }
}
