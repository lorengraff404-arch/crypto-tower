package db

import (
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// SeedAbilityLearning populates the ability_learning table
// Defines when characters can learn abilities based on Rank and Level
func SeedAbilityLearning(db *gorm.DB) error {
	// Check if already seeded
	var count int64
	db.Model(&models.AbilityLearning{}).Count(&count)
	if count > 0 {
		log.Println("AbilityLearning already seeded, skipping...")
		return nil
	}

	learningData := []models.AbilityLearning{
		// ==========================================
		// BASIC ABILITIES (All Ranks, Level 1)
		// ==========================================
		// Everyone gets Tackle (basic physical attack)
		{AbilityID: 1, MinRank: "C", LearnLevel: 1, IsStarting: true},

		// Everyone gets Defend (reduce damage)
		{AbilityID: 2, MinRank: "C", LearnLevel: 1, IsStarting: true},

		// ==========================================
		// INTERMEDIATE ABILITIES (Rank B+, Level 5)
		// ==========================================
		// Power Strike - Medium damage attack
		{AbilityID: 10, MinRank: "B", LearnLevel: 5},

		// Quick Attack - Low damage, high priority
		{AbilityID: 11, MinRank: "B", LearnLevel: 5},

		// Focus Energy - Increases critical hit ratio
		{AbilityID: 12, MinRank: "B", LearnLevel: 5},

		// ==========================================
		// ADVANCED ABILITIES (Rank A+, Level 10)
		// ==========================================
		// Mega Punch - High damage physical
		{AbilityID: 20, MinRank: "A", LearnLevel: 10},

		// Elemental Burst - Type-based attack
		{AbilityID: 21, MinRank: "A", LearnLevel: 10},

		// Battle Cry - Buff attack for 3 turns
		{AbilityID: 22, MinRank: "A", LearnLevel: 10},

		// Meditation - Restore mana
		{AbilityID: 23, MinRank: "A", LearnLevel: 10},

		// ==========================================
		// LEGENDARY ABILITIES (Rank S+, Level 15+)
		// ==========================================
		// Hyper Beam - Massive damage, recharge required
		{AbilityID: 30, MinRank: "S", LearnLevel: 15, IsUltimate: true},

		// Divine Shield - Immune to damage for 1 turn
		{AbilityID: 31, MinRank: "S", LearnLevel: 15},

		// Soul Drain - Damage + heal
		{AbilityID: 32, MinRank: "S", LearnLevel: 15},

		// ==========================================
		// SS TIER (Level 20+)
		// ==========================================
		// Meteor Strike - AOE damage
		{AbilityID: 40, MinRank: "SS", LearnLevel: 20, IsUltimate: true},

		// Time Warp - Extra turn
		{AbilityID: 41, MinRank: "SS", LearnLevel: 20},

		// ==========================================
		// SSS TIER (Level 25+)
		// ==========================================
		// Omega Destruction - Ultimate attack
		{AbilityID: 50, MinRank: "SSS", LearnLevel: 25, IsUltimate: true},

		// Immortality - Survive guaranteed KO once per battle
		{AbilityID: 51, MinRank: "SSS", LearnLevel: 25},
	}

	// Create in batches
	if err := db.CreateInBatches(&learningData, 100).Error; err != nil {
		return err
	}

	log.Printf("✅ Seeded %d ability learning entries", len(learningData))
	return nil
}
