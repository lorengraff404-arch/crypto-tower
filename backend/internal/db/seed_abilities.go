package db

import (
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// SeedAbilities seeds basic abilities for the game
func SeedAbilities(db *gorm.DB) error {
	var count int64
	db.Model(&models.Ability{}).Count(&count)
	if count > 0 {
		log.Println("Abilities already seeded, skipping...")
		return nil
	}

	abilities := []models.Ability{
		// ID 1-2: BASIC (All ranks, Level 1)
		{
			ID: 1, Name: "Tackle", Description: "Basic physical attack",
			Class: "Warrior", UnlockLevel: 1, AbilityType: "ACTIVE",
			TargetType: "SINGLE_ENEMY", ManaCost: 0, Damage: 40,
			Element: "Normal", Category: "physical", Accuracy: 100,
			DamageType: "physical", Rarity: "C",
		},
		{
			ID: 2, Name: "Defend", Description: "Reduce incoming damage by 50%",
			Class: "Tank", UnlockLevel: 1, AbilityType: "ACTIVE",
			TargetType: "SELF", ManaCost: 0, Damage: 0,
			Element: "Normal", Category: "status", Accuracy: 100,
			AppliesBuff: "defense_up", BuffDuration: 2,
			DamageType: "physical", Rarity: "C",
		},

		// ID 10-12: INTERMEDIATE (Rank B+, Level 5)
		{
			ID: 10, Name: "Power Strike", Description: "Medium damage attack",
			Class: "Warrior", UnlockLevel: 5, AbilityType: "ACTIVE",
			TargetType: "SINGLE_ENEMY", ManaCost: 20, Damage: 70,
			Element: "Normal", Category: "physical", Accuracy: 95,
			DamageType: "physical", Rarity: "B",
		},
		{
			ID: 11, Name: "Quick Attack", Description: "Low damage, high priority",
			Class: "Warrior", UnlockLevel: 5, AbilityType: "ACTIVE",
			TargetType: "SINGLE_ENEMY", ManaCost: 15, Damage: 50,
			Element: "Normal", Category: "physical", Accuracy: 100,
			Priority: 1, DamageType: "physical", Rarity: "B",
		},
		{
			ID: 12, Name: "Focus Energy", Description: "Increases critical hit ratio",
			Class: "Mage", UnlockLevel: 5, AbilityType: "ACTIVE",
			TargetType: "SELF", ManaCost: 15, Damage: 0,
			Element: "Normal", Category: "status", Accuracy: 100,
			AppliesBuff: "crit_rate_up", BuffDuration: 3,
			DamageType: "physical", Rarity: "B",
		},

		// ID 20-23: ADVANCED (Rank A+, Level 10)
		{
			ID: 20, Name: "Mega Punch", Description: "High damage physical attack",
			Class: "Warrior", UnlockLevel: 10, AbilityType: "ACTIVE",
			TargetType: "SINGLE_ENEMY", ManaCost: 40, Damage: 120,
			Element: "Normal", Category: "physical", Accuracy: 90,
			DamageType: "physical", Rarity: "A",
		},
		{
			ID: 21, Name: "Fireball", Description: "Fire elemental attack",
			Class: "Mage", UnlockLevel: 10, AbilityType: "ACTIVE",
			TargetType: "SINGLE_ENEMY", ManaCost: 35, Damage: 100,
			Element: "Fire", Category: "special", Accuracy: 95,
			StatusEffectChance: 10, AppliesDebuff: "burn",
			DamageType: "magical", Rarity: "A",
		},
		{
			ID: 22, Name: "Battle Cry", Description: "Buff attack for 3 turns",
			Class: "Warrior", UnlockLevel: 10, AbilityType: "ACTIVE",
			TargetType: "SELF", ManaCost: 25, Damage: 0,
			Element: "Normal", Category: "status", Accuracy: 100,
			AppliesBuff: "attack_up", BuffDuration: 3,
			DamageType: "physical", Rarity: "A",
		},
		{
			ID: 23, Name: "Meditation", Description: "Restore 50 mana",
			Class: "Mage", UnlockLevel: 10, AbilityType: "ACTIVE",
			TargetType: "SELF", ManaCost: 0, Damage: 0,
			Element: "Normal", Category: "status", Accuracy: 100,
			BaseHeal: 50, DamageType: "physical", Rarity: "A",
		},

		// ID 30-32: LEGENDARY (Rank S+, Level 15)
		{
			ID: 30, Name: "Hyper Beam", Description: "Massive damage, recharge required",
			Class: "Mage", UnlockLevel: 15, AbilityType: "ACTIVE",
			TargetType: "SINGLE_ENEMY", ManaCost: 80, Damage: 200,
			Element: "Normal", Category: "special", Accuracy: 90,
			Cooldown: 3, IsUltimate: true,
			DamageType: "magical", Rarity: "S",
		},
		{
			ID: 31, Name: "Divine Shield", Description: "Immune to damage for 1 turn",
			Class: "Tank", UnlockLevel: 15, AbilityType: "ACTIVE",
			TargetType: "SELF", ManaCost: 60, Damage: 0,
			Element: "Light", Category: "status", Accuracy: 100,
			AppliesBuff: "invulnerable", BuffDuration: 1,
			DamageType: "physical", Rarity: "S",
		},
		{
			ID: 32, Name: "Soul Drain", Description: "Damage + heal 50%",
			Class: "Mage", UnlockLevel: 15, AbilityType: "ACTIVE",
			TargetType: "SINGLE_ENEMY", ManaCost: 50, Damage: 90,
			Element: "Dark", Category: "special", Accuracy: 95,
			BaseHeal: 45, DamageType: "magical", Rarity: "S",
		},

		// ID 40-41: SS TIER (Level 20)
		{
			ID: 40, Name: "Meteor Strike", Description: "AOE massive damage",
			Class: "Mage", UnlockLevel: 20, AbilityType: "ACTIVE",
			TargetType: "AOE", ManaCost: 100, Damage: 150,
			Element: "Fire", Category: "special", Accuracy: 85,
			AOERadius: 3, IsUltimate: true,
			DamageType: "magical", Rarity: "SS",
		},
		{
			ID: 41, Name: "Time Warp", Description: "Extra turn",
			Class: "Mage", UnlockLevel: 20, AbilityType: "ACTIVE",
			TargetType: "SELF", ManaCost: 80, Damage: 0,
			Element: "Time", Category: "status", Accuracy: 100,
			Cooldown: 5, DamageType: "physical", Rarity: "SS",
		},

		// ID 50-51: SSS TIER (Level 25)
		{
			ID: 50, Name: "Omega Destruction", Description: "Ultimate attack",
			Class: "Mage", UnlockLevel: 25, AbilityType: "ULTIMATE",
			TargetType: "AOE", ManaCost: 150, Damage: 250,
			Element: "Chaos", Category: "special", Accuracy: 90,
			AOERadius: 5, IsUltimate: true,
			DamageType: "true", Rarity: "SSS",
		},
		{
			ID: 51, Name: "Immortality", Description: "Survive guaranteed KO once",
			Class: "Tank", UnlockLevel: 25, AbilityType: "PASSIVE",
			TargetType: "SELF", ManaCost: 0, Damage: 0,
			Element: "Divine", Category: "status", Accuracy: 100,
			DamageType: "physical", Rarity: "SSS",
		},
	}

	if err := db.Create(&abilities).Error; err != nil {
		return err
	}

	log.Printf("✅ Seeded %d abilities", len(abilities))
	return nil
}
