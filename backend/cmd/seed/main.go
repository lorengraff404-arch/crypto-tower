package main

import (
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

func main() {
	log.Println("🌱 Starting Database Seeding...")

	// Load config and connect
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := db.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Seed in order
	seedAbilities()
	seedSystemSettings()
	seedShopItems()

	log.Println("✅ Database seeding completed successfully!")
}

func seedAbilities() {
	log.Println("📚 Seeding abilities...")

	abilities := []models.Ability{
		// ========== WARRIOR ABILITIES ==========
		// C Rank
		{Name: "Basic Strike", Description: "A simple melee attack", Class: "Warrior", Rarity: "C", UnlockLevel: 1, BaseDamage: 15, ManaCost: 5, Cooldown: 0, Element: "Neutral", TargetType: "single"},
		{Name: "Power Attack", Description: "A strong overhead slam", Class: "Warrior", Rarity: "C", UnlockLevel: 3, BaseDamage: 25, ManaCost: 10, Cooldown: 1, Element: "Neutral", TargetType: "single"},
		{Name: "Shield Bash", Description: "Bash with shield", Class: "Warrior", Rarity: "C", UnlockLevel: 5, BaseDamage: 20, ManaCost: 8, Cooldown: 2, Element: "Neutral", TargetType: "single"},
		{Name: "Rage", Description: "Increase attack for 3 turns", Class: "Warrior", Rarity: "C", UnlockLevel: 7, BaseDamage: 0, ManaCost: 15, Cooldown: 4, Element: "Neutral", TargetType: "self"},
		{Name: "Cleave", Description: "Hit 2 enemies", Class: "Warrior", Rarity: "C", UnlockLevel: 10, BaseDamage: 18, ManaCost: 12, Cooldown: 1, Element: "Neutral", TargetType: "aoe"},
		{Name: "Battle Cry", Description: "Taunt all enemies", Class: "Warrior", Rarity: "C", UnlockLevel: 12, BaseDamage: 0, ManaCost: 10, Cooldown: 3, Element: "Neutral", TargetType: "aoe"},

		// B Rank
		{Name: "Whirlwind", Description: "Spinning attack hitting all", Class: "Warrior", Rarity: "B", UnlockLevel: 15, BaseDamage: 30, ManaCost: 20, Cooldown: 3, Element: "Air", TargetType: "aoe"},
		{Name: "Crushing Blow", Description: "Massive single damage", Class: "Warrior", Rarity: "B", UnlockLevel: 20, BaseDamage: 60, ManaCost: 25, Cooldown: 4, Element: "Earth", TargetType: "single"},
		{Name: "Iron Will", Description: "Reduce damage taken", Class: "Warrior", Rarity: "B", UnlockLevel: 25, BaseDamage: 0, ManaCost: 20, Cooldown: 5, Element: "Neutral", TargetType: "self"},

		// A Rank
		{Name: "Earthquake", Description: "Earth damage to all foes", Class: "Warrior", Rarity: "A", UnlockLevel: 35, BaseDamage: 45, ManaCost: 35, Cooldown: 4, Element: "Earth", TargetType: "aoe"},
		{Name: "Berserker", Description: "Trade defense for attack", Class: "Warrior", Rarity: "A", UnlockLevel: 45, BaseDamage: 0, ManaCost: 30, Cooldown: 6, Element: "Neutral", TargetType: "self"},

		// S Rank
		{Name: "Titan's Wrath", Description: "Devastating AOE", Class: "Warrior", Rarity: "S", UnlockLevel: 60, BaseDamage: 80, ManaCost: 50, Cooldown: 6, Element: "Earth", TargetType: "aoe"},

		// SS Rank
		{Name: "Immortal Stance", Description: "Cannot die for 3 turns", Class: "Warrior", Rarity: "SS", UnlockLevel: 80, BaseDamage: 0, ManaCost: 60, Cooldown: 10, Element: "Light", TargetType: "self"},

		// SSS Rank
		{Name: "Armageddon Slash", Description: "Ultimate slash", Class: "Warrior", Rarity: "SSS", UnlockLevel: 100, BaseDamage: 200, ManaCost: 100, Cooldown: 12, Element: "Dark", TargetType: "aoe"},

		// ========== MAGE ABILITIES ==========
		// C Rank
		{Name: "Fire Bolt", Description: "Basic fire projectile", Class: "Mage", Rarity: "C", UnlockLevel: 1, BaseDamage: 20, ManaCost: 10, Cooldown: 0, Element: "Fire", TargetType: "single"},
		{Name: "Ice Shard", Description: "Ice projectile", Class: "Mage", Rarity: "C", UnlockLevel: 2, BaseDamage: 18, ManaCost: 12, Cooldown: 1, Element: "Water", TargetType: "single"},
		{Name: "Lightning Strike", Description: "Quick lightning", Class: "Mage", Rarity: "C", UnlockLevel: 4, BaseDamage: 25, ManaCost: 15, Cooldown: 2, Element: "Lightning", TargetType: "single"},
		{Name: "Mana Shield", Description: "Absorb damage", Class: "Mage", Rarity: "C", UnlockLevel: 6, BaseDamage: 0, ManaCost: 20, Cooldown: 4, Element: "Neutral", TargetType: "self"},
		{Name: "Arcane Missiles", Description: "3 magic missiles", Class: "Mage", Rarity: "C", UnlockLevel: 8, BaseDamage: 12, ManaCost: 18, Cooldown: 1, Element: "Neutral", TargetType: "single"},
		{Name: "Frost Nova", Description: "Freeze nearby enemies", Class: "Mage", Rarity: "C", UnlockLevel: 11, BaseDamage: 15, ManaCost: 22, Cooldown: 3, Element: "Water", TargetType: "aoe"},
		{Name: "Fireball", Description: "Explosive fire damage", Class: "Mage", Rarity: "C", UnlockLevel: 13, BaseDamage: 35, ManaCost: 25, Cooldown: 2, Element: "Fire", TargetType: "single"},

		// B Rank
		{Name: "Chain Lightning", Description: "Lightning bounces", Class: "Mage", Rarity: "B", UnlockLevel: 18, BaseDamage: 40, ManaCost: 35, Cooldown: 3, Element: "Lightning", TargetType: "aoe"},
		{Name: "Blizzard", Description: "Ice storm", Class: "Mage", Rarity: "B", UnlockLevel: 22, BaseDamage: 35, ManaCost: 40, Cooldown: 4, Element: "Water", TargetType: "aoe"},
		{Name: "Flame Wall", Description: "Fire barrier", Class: "Mage", Rarity: "B", UnlockLevel: 27, BaseDamage: 0, ManaCost: 30, Cooldown: 5, Element: "Fire", TargetType: "self"},
		{Name: "Time Warp", Description: "Gain extra turn", Class: "Mage", Rarity: "B", UnlockLevel: 30, BaseDamage: 0, ManaCost: 50, Cooldown: 8, Element: "Neutral", TargetType: "self"},

		// A Rank
		{Name: "Meteor", Description: "Summon fiery meteor", Class: "Mage", Rarity: "A", UnlockLevel: 40, BaseDamage: 70, ManaCost: 55, Cooldown: 5, Element: "Fire", TargetType: "aoe"},
		{Name: "Frozen Tomb", Description: "Freeze enemy", Class: "Mage", Rarity: "A", UnlockLevel: 48, BaseDamage: 30, ManaCost: 45, Cooldown: 6, Element: "Water", TargetType: "single"},
		{Name: "Arcane Explosion", Description: "Massive AoE magic", Class: "Mage", Rarity: "A", UnlockLevel: 52, BaseDamage: 65, ManaCost: 60, Cooldown: 5, Element: "Neutral", TargetType: "aoe"},

		// S Rank
		{Name: "Phoenix Fire", Description: "Revive on death", Class: "Mage", Rarity: "S", UnlockLevel: 65, BaseDamage: 0, ManaCost: 70, Cooldown: 15, Element: "Fire", TargetType: "self"},
		{Name: "Black Hole", Description: "Suck all enemies", Class: "Mage", Rarity: "S", UnlockLevel: 70, BaseDamage: 90, ManaCost: 80, Cooldown: 7, Element: "Dark", TargetType: "aoe"},

		// SS Rank
		{Name: "Dimensional Rift", Description: "Reality tear", Class: "Mage", Rarity: "SS", UnlockLevel: 85, BaseDamage: 100, ManaCost: 90, Cooldown: 10, Element: "Dark", TargetType: "single"},

		// SSS Rank
		{Name: "Apocalypse", Description: "End of days", Class: "Mage", Rarity: "SSS", UnlockLevel: 100, BaseDamage: 250, ManaCost: 150, Cooldown: 15, Element: "Dark", TargetType: "aoe"},

		// ========== TANK ABILITIES ==========
		// C Rank
		{Name: "Provoke", Description: "Force enemy to attack", Class: "Tank", Rarity: "C", UnlockLevel: 1, BaseDamage: 5, ManaCost: 8, Cooldown: 2, Element: "Neutral", TargetType: "single"},
		{Name: "Shield Block", Description: "Reduce damage", Class: "Tank", Rarity: "C", UnlockLevel: 3, BaseDamage: 0, ManaCost: 10, Cooldown: 3, Element: "Neutral", TargetType: "self"},
		{Name: "Counter", Description: "Return damage", Class: "Tank", Rarity: "C", UnlockLevel: 5, BaseDamage: 0, ManaCost: 12, Cooldown: 4, Element: "Neutral", TargetType: "self"},
		{Name: "Fortify", Description: "Increase defense", Class: "Tank", Rarity: "C", UnlockLevel: 8, BaseDamage: 0, ManaCost: 15, Cooldown: 5, Element: "Earth", TargetType: "self"},
		{Name: "Guardian's Light", Description: "Heal ally", Class: "Tank", Rarity: "C", UnlockLevel: 10, BaseDamage: -30, ManaCost: 18, Cooldown: 3, Element: "Light", TargetType: "ally"},
		{Name: "Wall of Stone", Description: "Party defense", Class: "Tank", Rarity: "C", UnlockLevel: 12, BaseDamage: 0, ManaCost: 25, Cooldown: 6, Element: "Earth", TargetType: "aoe"},

		// B Rank
		{Name: "Last Stand", Description: "Cannot go below 1 HP", Class: "Tank", Rarity: "B", UnlockLevel: 16, BaseDamage: 0, ManaCost: 30, Cooldown: 8, Element: "Light", TargetType: "self"},
		{Name: "Earthquake Stomp", Description: "Stun all enemies", Class: "Tank", Rarity: "B", UnlockLevel: 21, BaseDamage: 25, ManaCost: 28, Cooldown: 5, Element: "Earth", TargetType: "aoe"},
		{Name: "Regeneration", Description: "Heal over time", Class: "Tank", Rarity: "B", UnlockLevel: 26, BaseDamage: 0, ManaCost: 35, Cooldown: 7, Element: "Light", TargetType: "self"},

		// A Rank
		{Name: "Titan's Shield", Description: "Invulnerable", Class: "Tank", Rarity: "A", UnlockLevel: 38, BaseDamage: 0, ManaCost: 50, Cooldown: 10, Element: "Earth", TargetType: "self"},
		{Name: "Mass Heal", Description: "Heal all allies", Class: "Tank", Rarity: "A", UnlockLevel: 46, BaseDamage: -50, ManaCost: 60, Cooldown: 6, Element: "Light", TargetType: "aoe"},

		// S Rank
		{Name: "Divine Protection", Description: "Party damage reduction", Class: "Tank", Rarity: "S", UnlockLevel: 62, BaseDamage: 0, ManaCost: 70, Cooldown: 10, Element: "Light", TargetType: "aoe"},

		// SS Rank
		{Name: "Sacrifice", Description: "Revive all allies", Class: "Tank", Rarity: "SS", UnlockLevel: 82, BaseDamage: -100, ManaCost: 80, Cooldown: 20, Element: "Light", TargetType: "aoe"},

		// SSS Rank
		{Name: "Eternal Guardian", Description: "Party immortality", Class: "Tank", Rarity: "SSS", UnlockLevel: 100, BaseDamage: 0, ManaCost: 120, Cooldown: 25, Element: "Light", TargetType: "aoe"},
	}

	for _, ability := range abilities {
		var count int64
		db.DB.Model(&models.Ability{}).Where("name = ?", ability.Name).Count(&count)
		if count == 0 {
			if err := db.DB.Create(&ability).Error; err != nil {
				log.Printf("❌ Failed to create ability %s: %v", ability.Name, err)
			} else {
				log.Printf("✅ Created ability: %s (%s rank)", ability.Name, ability.Rarity)
			}
		} else {
			log.Printf("⏭️  Skipping %s (already exists)", ability.Name)
		}
	}
	log.Printf("📊 Total abilities: %d", len(abilities))
}

func seedSystemSettings() {
	log.Println("⚙️  Seeding system settings...")

	settings := []models.SystemSetting{
		{Key: "ability_slots_c", Value: "4", Description: "Max slots for C-rank"},
		{Key: "ability_slots_b", Value: "6", Description: "Max slots for B-rank"},
		{Key: "ability_slots_a", Value: "8", Description: "Max slots for A-rank"},
		{Key: "ability_slots_s", Value: "10", Description: "Max slots for S-rank"},
		{Key: "ability_slots_ss", Value: "12", Description: "Max slots for SS-rank"},
		{Key: "ability_slots_sss", Value: "16", Description: "Max slots for SSS-rank"},
	}

	for _, setting := range settings {
		var count int64
		db.DB.Model(&models.SystemSetting{}).Where("key = ?", setting.Key).Count(&count)
		if count == 0 {
			if err := db.DB.Create(&setting).Error; err != nil {
				log.Printf("❌ Failed to create setting %s: %v", setting.Key, err)
			} else {
				log.Printf("✅ Created: %s = %s", setting.Key, setting.Value)
			}
		}
	}
}

func seedShopItems() {
	log.Println("🛍️  Seeding shop items...")

	shopItems := []models.ShopItem{
		// Consumables
		{Name: "Health Potion", Description: "Restores 50 HP", Category: "healing", EffectType: "heal_hp", EffectValue: 50, GTKCost: 100, IsConsumable: true, MaxStack: 99, IsAvailable: true, IconURL: "assets/items/potion_hp.png"},
		{Name: "Mana Potion", Description: "Restores 20 MP", Category: "healing", EffectType: "restore_mp", EffectValue: 20, GTKCost: 150, IsConsumable: true, MaxStack: 99, IsAvailable: true, IconURL: "assets/items/potion_mp.png"},
		{Name: "Elixir", Description: "Restores 50% HP & MP", Category: "healing", EffectType: "restore_all", EffectValue: 50, GTKCost: 400, IsConsumable: true, MaxStack: 99, IsAvailable: true, IconURL: "assets/items/elixir.png"},

		// Status Cures
		{Name: "Antidote", Description: "Cures Poison", Category: "status", EffectType: "cure_status", EffectValue: 0, GTKCost: 50, IsConsumable: true, MaxStack: 99, IsAvailable: true, IconURL: "assets/items/antidote.png"},

		// Breeding/Hatching
		{Name: "Incubator Heat Lamp", Description: "Reduces remaining incubation time by 1 hour", Category: "egg", EffectType: "reduce_time", EffectValue: 60, GTKCost: 500, IsConsumable: true, MaxStack: 10, IsAvailable: true, IconURL: "assets/items/heat_lamp.png"},
		{Name: "Nutrient Injection", Description: "Increases chance of higher stats", Category: "egg", EffectType: "buff_stats", EffectValue: 10, GTKCost: 1000, IsConsumable: true, MaxStack: 5, IsAvailable: true, IconURL: "assets/items/nutrient.png"},
	}

	for _, item := range shopItems {
		var count int64
		db.DB.Model(&models.ShopItem{}).Where("name = ?", item.Name).Count(&count)
		if count == 0 {
			if err := db.DB.Create(&item).Error; err != nil {
				log.Printf("❌ Failed to create shop item %s: %v", item.Name, err)
			} else {
				log.Printf("✅ Created shop item: %s", item.Name)
			}
		} else {
			log.Printf("⏭️  Skipping %s (already exists)", item.Name)
		}
	}
}
