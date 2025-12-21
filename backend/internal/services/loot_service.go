package services

import (
	"encoding/json"
	"math/rand"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// LootService handles loot generation and rewards
type LootService struct{}

// GenerateLoot generates loot drops for a completed mission
func (s *LootService) GenerateLoot(missionID uint, performanceGrade string) ([]models.LootDrop, error) {
	// Get loot table for mission
	var lootTable models.LootTable
	if err := db.DB.Preload("Entries.Item").
		Where("mission_id = ?", missionID).
		First(&lootTable).Error; err != nil {
		return nil, err
	}

	// Performance multiplier
	performanceMult := s.getPerformanceMultiplier(performanceGrade)

	drops := []models.LootDrop{}

	for _, entry := range lootTable.Entries {
		if entry.Item == nil {
			continue
		}

		// Roll for drop
		dropChance := entry.DropChance * performanceMult
		if rand.Float64()*100 < dropChance {
			// Determine quantity
			quantity := entry.MinQuantity
			if entry.MaxQuantity > entry.MinQuantity {
				quantity += rand.Intn(entry.MaxQuantity - entry.MinQuantity + 1)
			}

			drops = append(drops, models.LootDrop{
				ItemID:   entry.Item.ID,
				ItemName: entry.Item.Name,
				Quantity: quantity,
				Rarity:   entry.Item.Rarity,
			})
		}
	}

	return drops, nil
}

// CalculateRewards calculates tokens and XP for a battle
func (s *LootService) CalculateRewards(session *models.RaidSession, performanceGrade string) (tokens, xp int) {
	// Base rewards from mission
	baseTokens := 100 + (int(1) * 50)
	baseXP := 50 + (int(1) * 25)

	// Performance multiplier
	mult := s.getPerformanceMultiplier(performanceGrade)

	tokens = int(float64(baseTokens) * mult)
	xp = int(float64(baseXP) * mult)

	return
}

// CalculatePerformanceGrade determines grade based on battle performance
func (s *LootService) CalculatePerformanceGrade(session *models.RaidSession) string {
	// Factors:
	// - Turns taken (fewer = better)
	// - Damage taken (less = better)
	// - Characters fainted (fewer = better)

	maxTurns := 20
	turnScore := float64(maxTurns-session.TurnCount) / float64(maxTurns)
	if turnScore < 0 {
		turnScore = 0
	}

	// Damage score (less damage = higher score)
	damageScore := 1.0 - (float64(session.TotalDamageTaken) / float64(session.Mission.EnemyHP*3))
	if damageScore < 0 {
		damageScore = 0
	}

	// Average scores
	totalScore := (turnScore + damageScore) / 2.0

	// Convert to grade
	if totalScore >= 0.9 {
		return "S"
	} else if totalScore >= 0.75 {
		return "A"
	} else if totalScore >= 0.6 {
		return "B"
	} else if totalScore >= 0.4 {
		return "C"
	}
	return "D"
}

// SaveBattleReward saves reward to database
func (s *LootService) SaveBattleReward(sessionID, userID uint, tokens, xp int, drops []models.LootDrop, grade string) error {
	// Convert drops to JSON
	dropsJSON, _ := json.Marshal(drops)

	reward := models.BattleReward{
		RaidSessionID:    sessionID,
		UserID:           userID,
		TokensEarned:     tokens,
		XPEarned:         xp,
		ItemsJSON:        string(dropsJSON),
		PerformanceGrade: grade,
	}

	if err := db.DB.Create(&reward).Error; err != nil {
		return err
	}

	// Add items to inventory
	for _, drop := range drops {
		if err := s.AddToInventory(userID, drop.ItemID, drop.Quantity); err != nil {
			return err
		}
	}

	// Add tokens to user
	db.DB.Exec("UPDATE users SET tokens = tokens + ? WHERE id = ?", tokens, userID)

	return nil
}

// AddToInventory adds items to user inventory
func (s *LootService) AddToInventory(userID, itemID uint, quantity int) error {
	var inventory models.UserInventory

	// Check if item already exists
	err := db.DB.Where("user_id = ? AND item_id = ?", userID, itemID).First(&inventory).Error

	if err != nil {
		// Create new inventory entry
		inventory = models.UserInventory{
			UserID:   userID,
			ItemID:   itemID,
			Quantity: quantity,
		}
		return db.DB.Create(&inventory).Error
	}

	// Update existing quantity
	inventory.Quantity += quantity
	return db.DB.Save(&inventory).Error
}

// getPerformanceMultiplier returns reward multiplier based on grade
func (s *LootService) getPerformanceMultiplier(grade string) float64 {
	switch grade {
	case "S":
		return 2.0
	case "A":
		return 1.5
	case "B":
		return 1.2
	case "C":
		return 1.0
	case "D":
		return 0.8
	default:
		return 1.0
	}
}
