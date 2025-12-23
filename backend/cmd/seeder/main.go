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
	shopItems := []models.ShopItem{
		// Consumables - Healing
		{Name: "Small Potion", Description: "Restores 50 HP", GTKCost: 100, Category: "consumable", IsConsumable: true, EffectType: "heal_hp", EffectValue: 50, IconURL: "🧪", IsAvailable: true},
		{Name: "Medium Potion", Description: "Restores 150 HP", GTKCost: 300, Category: "consumable", IsConsumable: true, EffectType: "heal_hp", EffectValue: 150, IconURL: "🧴", IsAvailable: true},
		{Name: "Large Potion", Description: "Restores 400 HP", GTKCost: 700, Category: "consumable", IsConsumable: true, EffectType: "heal_hp", EffectValue: 400, IconURL: "⚗️", IsAvailable: true},
		{Name: "Mega Potion", Description: "Fully restores HP", GTKCost: 1500, Category: "consumable", IsConsumable: true, EffectType: "heal_hp", EffectValue: 9999, IconURL: "💊", IsAvailable: true},
		{Name: "Revival Herb", Description: "Revives KO'd character", GTKCost: 2000, Category: "consumable", IsConsumable: true, EffectType: "revive", IconURL: "🌿", IsAvailable: true},
		{Name: "Full Restore", Description: "Restores HP + all status", GTKCost: 2500, Category: "consumable", IsConsumable: true, EffectType: "full_heal", IconURL: "✨", IsAvailable: true},

		// Status Cures
		{Name: "Antidote", Description: "Cures Poison", GTKCost: 100, Category: "consumable", IsConsumable: true, EffectType: "cure_poison", IconURL: "💉", IsAvailable: true},
		{Name: "Burn Heal", Description: "Cures Burn", GTKCost: 150, Category: "consumable", IsConsumable: true, EffectType: "cure_burn", IconURL: "🧊", IsAvailable: true},
		{Name: "Ice Heal", Description: "Cures Freeze", GTKCost: 150, Category: "consumable", IsConsumable: true, EffectType: "cure_freeze", IconURL: "🔥", IsAvailable: true},
		{Name: "Paralyze Heal", Description: "Cures Paralysis", GTKCost: 150, Category: "consumable", IsConsumable: true, EffectType: "cure_paralysis", IconURL: "⚡", IsAvailable: true},

		// Equipment - Weapons
		{Name: "Bronze Dagger", Description: "Fast but weak (+3 Atk)", GTKCost: 80, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 3, IconURL: "🗡️", IsAvailable: true},
		{Name: "Wooden Sword", Description: "Basic sword (+5 Atk)", GTKCost: 150, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 5, IconURL: "⚔️", IsAvailable: true},
		{Name: "Iron Sword", Description: "Standard sword (+15 Atk)", GTKCost: 500, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 15, IconURL: "⚔️", IsAvailable: true},
		{Name: "Steel Broadsword", Description: "Heavy sword (+35 Atk)", GTKCost: 1200, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 35, IconURL: "🗡️", IsAvailable: true},
		{Name: "Mithril Blade", Description: "Magical sword (+60 Atk)", GTKCost: 5000, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 60, IconURL: "⚔️", IsAvailable: true},
		{Name: "Dragon Bone Blade", Description: "Legendary weapon (+120 Atk)", GTKCost: 25000, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 120, IconURL: "🐉", IsAvailable: true},
		{Name: "Crystal Staff", Description: "Focuses magical energy (+45 Atk)", GTKCost: 3500, Category: "weapon", IsConsumable: false, EffectType: "equip_atk", EffectValue: 45, IconURL: "🪄", IsAvailable: true},

		// Equipment - Armor/Defensive
		{Name: "Leather Armor", Description: "Light protection (+5 Def)", GTKCost: 200, Category: "armor", IsConsumable: false, EffectType: "equip_def", EffectValue: 5, IconURL: "🦺", IsAvailable: true},
		{Name: "Wooden Shield", Description: "Basic shield (+3 Def)", GTKCost: 100, Category: "armor", IsConsumable: false, EffectType: "equip_def", EffectValue: 3, IconURL: "🛡️", IsAvailable: true},
		{Name: "Iron Buckler", Description: "Sturdy shield (+10 Def)", GTKCost: 400, Category: "armor", IsConsumable: false, EffectType: "equip_def", EffectValue: 10, IconURL: "🛡️", IsAvailable: true},
		{Name: "Steel Plate", Description: "Heavy protection (+30 Def)", GTKCost: 2000, Category: "armor", IsConsumable: false, EffectType: "equip_def", EffectValue: 30, IconURL: "🛡️", IsAvailable: true},
		{Name: "Tower Shield", Description: "Massive defense (+50 Def)", GTKCost: 4500, Category: "armor", IsConsumable: false, EffectType: "equip_def", EffectValue: 50, IconURL: "🛡️", IsAvailable: true},
		{Name: "Runite Guard", Description: "Enchanted armor (+85 Def)", GTKCost: 15000, Category: "armor", IsConsumable: false, EffectType: "equip_def", EffectValue: 85, IconURL: "💎", IsAvailable: true},

		// Egg Items (All under 'egg' category)
		{Name: "Basic Nest", Description: "Simple incubation spot (-1h)", GTKCost: 200, Category: "egg", IsConsumable: true, EffectType: "accelerate", EffectValue: 60, IconURL: "🥚", IsAvailable: true},
		{Name: "Advanced Incubator", Description: "High-tech warmth (-6h)", GTKCost: 1200, Category: "egg", IsConsumable: true, EffectType: "accelerate", EffectValue: 360, IconURL: "🔬", IsAvailable: true},
		{Name: "Solar Heat Lamp", Description: "Extreme acceleration (-24h)", GTKCost: 4000, Category: "egg", IsConsumable: true, EffectType: "accelerate", EffectValue: 1440, IconURL: "☀️", IsAvailable: true},
		{Name: "Time Skip Device", Description: "Skip all remaining time", GTKCost: 8000, Category: "egg", IsConsumable: true, EffectType: "instant_hatch", IconURL: "⏰", IsAvailable: true},
		{Name: "Egg Scanner", Description: "View hidden egg stats", GTKCost: 1000, Category: "egg", IsConsumable: true, EffectType: "scan", IconURL: "🔍", IsAvailable: true},
		{Name: "Trait Revealer", Description: "Show concealed traits", GTKCost: 1500, Category: "egg", IsConsumable: true, EffectType: "scan", IconURL: "👁️", IsAvailable: true},

		// Special/Currency
		{Name: "XP Scroll", Description: "Instant 500 XP to character", GTKCost: 750, Category: "currency", IsConsumable: true, EffectType: "grant_xp", EffectValue: 500, IconURL: "📜", IsAvailable: true},
		{Name: "Master XP Scroll", Description: "Instant 2500 XP to character", GTKCost: 3000, Category: "currency", IsConsumable: true, EffectType: "grant_xp", EffectValue: 2500, IconURL: "📖", IsAvailable: true},
		{Name: "TOWER Voucher", Description: "Redeemable for 10 TOWER", GTKCost: 10000, Category: "currency", IsConsumable: true, IconURL: "🎫", IsAvailable: true},
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
