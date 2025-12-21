package main

import (
	"fmt"
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := db.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	fmt.Println("🌱 Seeding game data...")

	seedSkills()
	// seedItems() // Deprecated: Using ShopItems as source of truth
	seedIslands()
	seedShopItems()

	fmt.Println("✅ Seeding complete!")
}

func seedSkills() {
	abilities := []models.Ability{
		// Warrior Skills
		{Name: "Slash", Description: "Basic physical attack", BaseDamage: 60, ManaCost: 0, Cooldown: 0, AbilityType: "active", Class: "Warrior", DamageType: "physical", AnimationName: "slash"},
		{Name: "Heavy Strike", Description: "Strong physical attack", BaseDamage: 120, ManaCost: 30, Cooldown: 3, AbilityType: "active", Class: "Warrior", DamageType: "physical", AnimationName: "heavy_strike"},
		{Name: "Iron Will", Description: "Boosts Defense", BuffDuration: 3, AppliesBuff: "defense_up", ManaCost: 20, Cooldown: 5, AbilityType: "active", Class: "Warrior", AnimationName: "buff_def"},

		// Mage Skills
		{Name: "Fireball", Description: "Ranged fire damage", BaseDamage: 80, ManaCost: 20, Cooldown: 2, AbilityType: "active", Class: "Mage", DamageType: "magical", AnimationName: "fireball", Element: "Fire"},
		{Name: "Ice Spike", Description: "Ice damage, chance to freeze", BaseDamage: 70, ManaCost: 25, Cooldown: 3, AbilityType: "active", Class: "Mage", DamageType: "magical", AnimationName: "ice_spike", Element: "Ice", AppliesDebuff: "freeze"},
		{Name: "Thunderbolt", Description: "High lightning damage", BaseDamage: 150, ManaCost: 50, Cooldown: 4, AbilityType: "active", Class: "Mage", DamageType: "magical", AnimationName: "thunderbolt", Element: "Electric"},

		// Archer Skills
		{Name: "Quick Shot", Description: "Fast arrow shot", BaseDamage: 50, ManaCost: 10, Cooldown: 1, AbilityType: "active", Class: "Archer", DamageType: "physical", AnimationName: "arrow"},
		{Name: "Poison Tip", Description: "Applies poison", BaseDamage: 40, ManaCost: 25, Cooldown: 3, AbilityType: "active", Class: "Archer", DamageType: "physical", AnimationName: "poison_arrow", AppliesDebuff: "poison"},
	}

	for _, a := range abilities {
		if err := db.DB.Where("name = ?", a.Name).FirstOrCreate(&a).Error; err != nil {
			log.Printf("Failed to seed ability %s: %v", a.Name, err)
		}
	}
}

func seedShopItems() {
	// Provide ID mapping to real items manually or verify name
	shopItems := []models.ShopItem{
		// Consumables
		{Name: "Small Potion", Description: "Heals 50 HP", GTKCost: 100, Category: "consumable", IsConsumable: true, EffectType: "heal_hp", EffectValue: 50, MaxStack: 99, IsAvailable: true},
		{Name: "Medium Potion", Description: "Heals 150 HP", GTKCost: 300, Category: "consumable", IsConsumable: true, EffectType: "heal_hp", EffectValue: 150, MaxStack: 99, IsAvailable: true},
		{Name: "Large Potion", Description: "Heals 500 HP", GTKCost: 800, Category: "consumable", IsConsumable: true, EffectType: "heal_hp", EffectValue: 500, MaxStack: 99, IsAvailable: true},
		{Name: "Revival Herb", Description: "Revives fainted character", GTKCost: 500, Category: "consumable", IsConsumable: true, EffectType: "revive", EffectValue: 50, MaxStack: 10, IsAvailable: true},
		{Name: "Elixir of Energy", Description: "Restores 50 Energy", GTKCost: 200, Category: "consumable", IsConsumable: true, EffectType: "restore_energy", EffectValue: 50, MaxStack: 99, IsAvailable: true},

		// Equipment - Weapons
		{Name: "Wooden Sword", Description: "Basic sword (+5 Atk)", GTKCost: 150, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 5, IsAvailable: true},
		{Name: "Iron Sword", Description: "Standard sword (+15 Atk)", GTKCost: 500, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 15, IsAvailable: true},
		{Name: "Steel Broadsword", Description: "Heavy sword (+30 Atk)", GTKCost: 1200, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 30, IsAvailable: true},
		{Name: "Mithril Blade", Description: "Rare sword (+50 Atk)", GTKCost: 5000, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 50, IsAvailable: true},

		// Equipment - Shields
		{Name: "Wooden Shield", Description: "Basic shield (+3 Def)", GTKCost: 100, Category: "armor", IsConsumable: false, EffectType: "equip_def", EffectValue: 3, IsAvailable: true},
		{Name: "Iron Buckler", Description: "Sturdy shield (+10 Def)", GTKCost: 400, Category: "armor", IsConsumable: false, EffectType: "equip_def", EffectValue: 10, IsAvailable: true},
		{Name: "Tower Shield", Description: "Massive shield (+25 Def)", GTKCost: 1500, Category: "armor", IsConsumable: false, EffectType: "equip_def", EffectValue: 25, IsAvailable: true},

		// Special
		{Name: "TOWER Token", Description: "Premium Currency", GTKCost: 1000, Category: "currency", IsConsumable: false, IsAvailable: true},
		{Name: "XP Scroll", Description: "Grants 1000 XP", GTKCost: 2000, Category: "consumable", IsConsumable: true, EffectType: "grant_xp", EffectValue: 1000, MaxStack: 99, IsAvailable: true},
	}
	for _, s := range shopItems {
		if err := db.DB.Where("name = ?", s.Name).FirstOrCreate(&s).Error; err != nil {
			log.Printf("Failed to seed shop item %s: %v", s.Name, err)
		}
	}
}

func seedIslands() {
	island := models.Island{
		Name:        "Tutorial Island",
		Description: "The beginning of your journey",
		Difficulty:  1,
	}
	db.DB.Where("name = ?", island.Name).FirstOrCreate(&island)

	missions := []models.IslandMission{
		{IslandID: island.ID, Sequence: 1, Name: "Slime Encounter", EnemyName: "Green Slime", EnemyHP: 100, EnemyAtk: 10, EnemyDef: 5, EnemySpeed: 5, EnemyType: "Grass", RewardsPool: "{}"},
		{IslandID: island.ID, Sequence: 2, Name: "Forest Wolf", EnemyName: "Dire Wolf", EnemyHP: 250, EnemyAtk: 25, EnemyDef: 10, EnemySpeed: 15, EnemyType: "Normal", RewardsPool: "{}"},
		{IslandID: island.ID, Sequence: 3, Name: "Goblin Camp", EnemyName: "Goblin Scout", EnemyHP: 400, EnemyAtk: 40, EnemyDef: 15, EnemySpeed: 20, EnemyType: "Earth", RewardsPool: "{}"},
		{IslandID: island.ID, Sequence: 4, Name: "Guardian Golem", EnemyName: "Stone Golem", EnemyHP: 1000, EnemyAtk: 80, EnemyDef: 50, EnemySpeed: 10, EnemyType: "Earth", RewardsPool: "{}"}, // Boss
	}

	for _, m := range missions {
		if err := db.DB.Where("island_id = ? AND sequence = ?", m.IslandID, m.Sequence).FirstOrCreate(&m).Error; err != nil {
			log.Printf("Failed to seed mission %s: %v", m.Name, err)
		}
	}
}
