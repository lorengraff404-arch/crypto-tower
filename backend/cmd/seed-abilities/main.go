package main

import (
	"encoding/json"
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := db.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("🎯 Seeding Abilities Database...")

	// Seed all abilities
	seedWarriorAbilities()
	seedMageAbilities()
	seedTankAbilities()

	log.Println("✅ Ability seeding complete!")
}

func seedWarriorAbilities() {
	log.Println("  ⚔️  Seeding Warrior abilities...")

	elementBonusDefault, _ := json.Marshal(models.ElementBonus{
		Fire: 1.2, Water: 1.0, Ice: 1.0, Thunder: 1.0,
		Dark: 1.0, Plant: 0.9, Earth: 1.1, Wind: 1.0,
	})

	abilities := []models.Ability{
		// Level 1 - Starter
		{
			Name:          "Slash",
			Description:   "A swift sword strike dealing moderate physical damage to a single enemy.",
			Class:         "Warrior",
			UnlockLevel:   1,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      3,
			ManaCost:      15,
			BaseDamage:    80,
			IconURL:       "/assets/abilities/warrior/slash.png",
			AnimationName: "slash_effect",
			SoundEffect:   "sword_slash.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		{
			Name:          "Block",
			Description:   "Raise your shield to reduce incoming damage by 40% for 4 seconds.",
			Class:         "Warrior",
			UnlockLevel:   1,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      8,
			ManaCost:      20,
			DurationSecs:  4,
			EffectPower:   40,
			AppliesBuff:   "Defense Up",
			IconURL:       "/assets/abilities/warrior/block.png",
			AnimationName: "shield_aura",
			SoundEffect:   "shield_up.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 5
		{
			Name:          "Power Strike",
			Description:   "Channel your strength into a devastating blow dealing 150% weapon damage.",
			Class:         "Warrior",
			UnlockLevel:   5,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      6,
			ManaCost:      25,
			BaseDamage:    150,
			IconURL:       "/assets/abilities/warrior/power_strike.png",
			AnimationName: "power_slash",
			SoundEffect:   "heavy_impact.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 10
		{
			Name:          "Double Attack",
			Description:   "Strike twice in rapid succession, each hit dealing 70% damage.",
			Class:         "Warrior",
			UnlockLevel:   10,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      10,
			ManaCost:      30,
			BaseDamage:    140, // 70 x 2
			IconURL:       "/assets/abilities/warrior/double_attack.png",
			AnimationName: "dual_slash",
			SoundEffect:   "double_slash.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 15
		{
			Name:          "Iron Defense",
			Description:   "Fortify your armor, increasing defense by 50% for 8 seconds.",
			Class:         "Warrior",
			UnlockLevel:   15,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      15,
			ManaCost:      35,
			DurationSecs:  8,
			EffectPower:   50,
			AppliesBuff:   "Fortified",
			IconURL:       "/assets/abilities/warrior/iron_defense.png",
			AnimationName: "iron_skin",
			SoundEffect:   "armor_clank.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 20
		{
			Name:          "Charge",
			Description:   "Rush forward and slam into an enemy, dealing damage and stunning for 2 seconds.",
			Class:         "Warrior",
			UnlockLevel:   20,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      12,
			ManaCost:      40,
			BaseDamage:    100,
			DurationSecs:  2,
			AppliesDebuff: "Stunned",
			IconURL:       "/assets/abilities/warrior/charge.png",
			AnimationName: "charge_impact",
			SoundEffect:   "charge.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 25
		{
			Name:          "Rend",
			Description:   "Tear through armor, dealing damage and applying Bleed (2% HP/sec) for 6 seconds.",
			Class:         "Warrior",
			UnlockLevel:   25,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      14,
			ManaCost:      45,
			BaseDamage:    120,
			DurationSecs:  6,
			AppliesDebuff: "Bleed",
			IconURL:       "/assets/abilities/warrior/rend.png",
			AnimationName: "rend_effect",
			SoundEffect:   "rend.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 30
		{
			Name:          "Battle Cry",
			Description:   "Rally nearby allies, increasing their attack by 25% for 10 seconds.",
			Class:         "Warrior",
			UnlockLevel:   30,
			AbilityType:   "ACTIVE",
			TargetType:    "ALL_ALLIES",
			Cooldown:      20,
			ManaCost:      50,
			DurationSecs:  10,
			EffectPower:   25,
			AppliesBuff:   "Inspired",
			IconURL:       "/assets/abilities/warrior/battle_cry.png",
			AnimationName: "rally_aura",
			SoundEffect:   "war_cry.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 40
		{
			Name:          "Whirlwind",
			Description:   "Spin with your weapon, hitting all nearby enemies for 90% damage.",
			Class:         "Warrior",
			UnlockLevel:   40,
			AbilityType:   "ACTIVE",
			TargetType:    "AOE",
			Cooldown:      18,
			ManaCost:      60,
			BaseDamage:    180,
			IconURL:       "/assets/abilities/warrior/whirlwind.png",
			AnimationName: "whirlwind_spin",
			SoundEffect:   "whirlwind.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 50
		{
			Name:          "Titan Strike",
			Description:   "Deliver a crushing blow that breaks enemy defense, dealing massive damage and reducing their defense by 30% for 5 seconds.",
			Class:         "Warrior",
			UnlockLevel:   50,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      25,
			ManaCost:      70,
			BaseDamage:    250,
			DurationSecs:  5,
			AppliesDebuff: "Armor Broken",
			IconURL:       "/assets/abilities/warrior/titan_strike.png",
			AnimationName: "titan_impact",
			SoundEffect:   "titan_slam.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 60
		{
			Name:          "Unstoppable Force",
			Description:   "Become immune to crowd control and gain 40% movement speed for 6 seconds.",
			Class:         "Warrior",
			UnlockLevel:   60,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      30,
			ManaCost:      80,
			DurationSecs:  6,
			EffectPower:   40,
			AppliesBuff:   "Unstoppable",
			IconURL:       "/assets/abilities/warrior/unstoppable.png",
			AnimationName: "berserker_glow",
			SoundEffect:   "roar.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 80
		{
			Name:          "Warrior's Wrath",
			Description:   "Enter a state of pure fury, doubling attack speed and gaining lifesteal for 12 seconds.",
			Class:         "Warrior",
			UnlockLevel:   80,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      45,
			ManaCost:      100,
			DurationSecs:  12,
			EffectPower:   100,
			AppliesBuff:   "Wrath",
			IconURL:       "/assets/abilities/warrior/wrath.png",
			AnimationName: "rage_aura",
			SoundEffect:   "battle_rage.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
		// Level 100 - ULTIMATE
		{
			Name:          "Sword of Legends",
			Description:   "Summon the legendary blade of heroes, dealing devastating AOE damage and executing enemies below 20% HP instantly.",
			Class:         "Warrior",
			UnlockLevel:   100,
			AbilityType:   "ULTIMATE",
			TargetType:    "AOE",
			Cooldown:      120,
			ManaCost:      150,
			BaseDamage:    500,
			IconURL:       "/assets/abilities/warrior/legendary_sword.png",
			AnimationName: "excalibur",
			SoundEffect:   "legendary_strike.mp3",
			ElementBonuses: string(elementBonusDefault),
		},
	}

	for _, ability := range abilities {
		if err := db.DB.FirstOrCreate(&ability, models.Ability{Name: ability.Name, Class: ability.Class}).Error; err != nil {
			log.Printf("    ⚠️  Error seeding %s: %v", ability.Name, err)
		} else {
			log.Printf("    ✓ Seeded: Lvl %d - %s", ability.UnlockLevel, ability.Name)
		}
	}
}

func seedMageAbilities() {
	log.Println("  🔮 Seeding Mage abilities...")

	elementBonusMagic, _ := json.Marshal(models.ElementBonus{
		Fire: 1.3, Water: 1.2, Ice: 1.2, Thunder: 1.3,
		Dark: 1.4, Plant: 1.0, Earth: 0.9, Wind: 1.1,
	})

	abilities := []models.Ability{
		// Level 1
		{
			Name:          "Magic Missile",
			Description:   "Fire a bolt of pure magical energy that never misses.",
			Class:         "Mage",
			UnlockLevel:   1,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      2,
			ManaCost:      10,
			BaseDamage:    60,
			IconURL:       "/assets/abilities/mage/magic_missile.png",
			AnimationName: "arcane_missile",
			SoundEffect:   "magic_whoosh.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		{
			Name:          "Mana Shield",
			Description:   "Create a magical barrier absorbing 150 damage.",
			Class:         "Mage",
			UnlockLevel:   1,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      10,
			ManaCost:      25,
			EffectPower:   150,
			AppliesBuff:   "Mana Shield",
			IconURL:       "/assets/abilities/mage/mana_shield.png",
			AnimationName: "shield_bubble",
			SoundEffect:   "magic_shield.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 5
		{
			Name:          "Fire Ball",
			Description:   "Hurl a blazing fireball, dealing fire damage and burning enemies.",
			Class:         "Mage",
			UnlockLevel:   5,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      5,
			ManaCost:      30,
			BaseDamage:    120,
			DurationSecs:  4,
			AppliesDebuff: "Burn",
			IconURL:       "/assets/abilities/mage/fireball.png",
			AnimationName: "fireball_explosion",
			SoundEffect:   "fireball.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 10
		{
			Name:          "Ice Shard",
			Description:   "Launch sharp ice projectiles, dealing damage and slowing by 40% for 3 seconds.",
			Class:         "Mage",
			UnlockLevel:   10, 
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      6,
			ManaCost:      35,
			BaseDamage:    100,
			DurationSecs:  3,
			AppliesDebuff: "Slow",
			IconURL:       "/assets/abilities/mage/ice_shard.png",
			AnimationName: "ice_spike",
			SoundEffect:   "ice_shatter.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 15
		{
			Name:          "Lightning Bolt",
			Description:   "Call down lightning, dealing high electric damage with a chance to stun.",
			Class:         "Mage",
			UnlockLevel:   15,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      8,
			ManaCost:      40,
			BaseDamage:    140,
			IconURL:       "/assets/abilities/mage/lightning.png",
			AnimationName: "lightning_strike",
			SoundEffect:   "thunder.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 20
		{
			Name:          "Mana Burst",
			Description:   "Release a wave of mana, dealing damage to all nearby enemies.",
			Class:         "Mage",
			UnlockLevel:   20,
			AbilityType:   "ACTIVE",
			TargetType:    "AOE",
			Cooldown:      12,
			ManaCost:      50,
			BaseDamage:    160,
			IconURL:       "/assets/abilities/mage/mana_burst.png",
			AnimationName: "mana_explosion",
			SoundEffect:   "mana_blast.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 25
		{
			Name:          "Chain Lightning",
			Description:   "Unleash lightning that bounces to 3 additional targets, each dealing 80% damage.",
			Class:         "Mage",
			UnlockLevel:   25,
			AbilityType:   "ACTIVE",
			TargetType:    "CHAIN",
			Cooldown:      15,
			ManaCost:      60,
			BaseDamage:    200,
			IconURL:       "/assets/abilities/mage/chain_lightning.png",
			AnimationName: "chain_zap",
			SoundEffect:   "chain_lightning.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 30
		{
			Name:          "Meteor",
			Description:   "Summon a meteor from the sky, dealing massive AOE fire damage.",
			Class:         "Mage",
			UnlockLevel:   30,
			AbilityType:   "ACTIVE",
			TargetType:    "AOE",
			Cooldown:      20,
			ManaCost:      75,
			BaseDamage:    280,
			AppliesDebuff: "Burn",
			IconURL:       "/assets/abilities/mage/meteor.png",
			AnimationName: "meteor_impact",
			SoundEffect:   "meteor_crash.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 40
		{
			Name:          "Frost Nova",
			Description:   "Freeze all nearby enemies solid for 3 seconds.",
			Class:         "Mage",
			UnlockLevel:   40,
			AbilityType:   "ACTIVE",
			TargetType:    "AOE",
			Cooldown:      25,
			ManaCost:      80,
			DurationSecs:  3,
			AppliesDebuff: "Frozen",
			IconURL:       "/assets/abilities/mage/frost_nova.png",
			AnimationName: "freeze_explosion",
			SoundEffect:   "ice_burst.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 50
		{
			Name:          "Arcane Explosion",
			Description:   "Detonate pure arcane energy, dealing massive damage and silencing enemies for 4 seconds.",
			Class:         "Mage",
			UnlockLevel:   50,
			AbilityType:   "ACTIVE",
			TargetType:    "AOE",
			Cooldown:      30,
			ManaCost:      100,
			BaseDamage:    350,
			DurationSecs:  4,
			AppliesDebuff: "Silenced",
			IconURL:       "/assets/abilities/mage/arcane_explosion.png",
			AnimationName: "arcane_nova",
			SoundEffect:   "arcane_boom.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 60
		{
			Name:          "Time Warp",
			Description:   "Slow time itself, reducing all enemy attack and movement speed by 50% for 6 seconds.",
			Class:         "Mage",
			UnlockLevel:   60,
			AbilityType:   "ACTIVE",
			TargetType:    "AOE",
			Cooldown:      40,
			ManaCost:      120,
			DurationSecs:  6,
			EffectPower:   50,
			AppliesDebuff: "Time Slowed",
			IconURL:       "/assets/abilities/mage/time_warp.png",
			AnimationName: "chronosphere",
			SoundEffect:   "time_distortion.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 80
		{
			Name:          "Elemental Mastery",
			Description:   "Channel all elements at once, gaining 60% spell power and casting speed for 15 seconds.",
			Class:         "Mage",
			UnlockLevel:   80,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      60,
			ManaCost:      150,
			DurationSecs:  15,
			EffectPower:   60,
			AppliesBuff:   "Elemental Master",
			IconURL:       "/assets/abilities/mage/elemental_mastery.png",
			AnimationName: "prismatic_aura",
			SoundEffect:   "elemental_surge.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
		// Level 100 - ULTIMATE
		{
			Name:          "Reality Tear",
			Description:   "Rip open the fabric of reality, dealing catastrophic damage to all enemies and banishing the weakest.",
			Class:         "Mage",
			UnlockLevel:   100,
			AbilityType:   "ULTIMATE",
			TargetType:    "AOE",
			Cooldown:      150,
			ManaCost:      200,
			BaseDamage:    666,
			IconURL:       "/assets/abilities/mage/reality_tear.png",
			AnimationName: "void_rift",
			SoundEffect:   "reality_shatter.mp3",
			ElementBonuses: string(elementBonusMagic),
		},
	}

	for _, ability := range abilities {
		if err := db.DB.FirstOrCreate(&ability, models.Ability{Name: ability.Name, Class: ability.Class}).Error; err != nil {
			log.Printf("    ⚠️  Error seeding %s: %v", ability.Name, err)
		} else {
			log.Printf("    ✓ Seeded: Lvl %d - %s", ability.UnlockLevel, ability.Name)
		}
	}
}

func seedTankAbilities() {
	log.Println("  🛡️  Seeding Tank abilities...")

	elementBonusTank, _ := json.Marshal(models.ElementBonus{
		Fire: 0.9, Water: 1.1, Ice: 1.0, Thunder: 0.9,
		Dark: 1.0, Plant: 1.0, Earth: 1.4, Wind: 0.8,
	})

	abilities := []models.Ability{
		// Level 1
		{
			Name:          "Taunt",
			Description:   "Force all nearby enemies to attack you for 4 seconds.",
			Class:         "Tank",
			UnlockLevel:   1,
			AbilityType:   "ACTIVE",
			TargetType:    "AOE",
			Cooldown:      8,
			ManaCost:      20,
			DurationSecs:  4,
			AppliesDebuff: "Taunted",
			IconURL:       "/assets/abilities/tank/taunt.png",
			AnimationName: "threat_wave",
			SoundEffect:   "taunt.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		{
			Name:          "Shield Bash",
			Description:   "Bash with your shield, dealing damage and stunning for 1.5 seconds.",
			Class:         "Tank",
			UnlockLevel:   1,
			AbilityType:   "ACTIVE",
			TargetType:    "SINGLE_ENEMY",
			Cooldown:      6,
			ManaCost:      15,
			BaseDamage:    50,
			DurationSecs:  2,
			AppliesDebuff: "Stunned",
			IconURL:       "/assets/abilities/tank/shield_bash.png",
			AnimationName: "shield_slam",
			SoundEffect:   "bash.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 5
		{
			Name:          "Fortify",
			Description:   "Reinforce your defenses, increasing max HP by 30% for 10 seconds.",
			Class:         "Tank",
			UnlockLevel:   5,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      15,
			ManaCost:      30,
			DurationSecs:  10,
			EffectPower:   30,
			AppliesBuff:   "Fortified",
			IconURL:       "/assets/abilities/tank/fortify.png",
			AnimationName: "iron_aura",
			SoundEffect:   "iron_skin.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 10
		{
			Name:          "Reflect Damage",
			Description:   "Return 50% of damage taken back to attackers for 5 seconds.",
			Class:         "Tank",
			UnlockLevel:   10,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      18,
			ManaCost:      40,
			DurationSecs:  5,
			EffectPower:   50,
			AppliesBuff:   "Thorns",
			IconURL:       "/assets/abilities/tank/reflect.png",
			AnimationName: "spike_aura",
			SoundEffect:   "thorns.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 15
		{
			Name:          "Last Stand",
			Description:   "Cannot be reduced below 1 HP for 4 seconds. Increases defense by 100%.",
			Class:         "Tank",
			UnlockLevel:   15,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      45,
			ManaCost:      50,
			DurationSecs:  4,
			EffectPower:   100,
			AppliesBuff:   "Undying",
			IconURL:       "/assets/abilities/tank/last_stand.png",
			AnimationName: "golden_shield",
			SoundEffect:   "immortal.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 20
		{
			Name:          "Shield Wall",
			Description:   "Create an impenetrable wall, blocking all projectiles and reducing damage by 80% for allies behind you.",
			Class:         "Tank",
			UnlockLevel:   20,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      30,
			ManaCost:      60,
			DurationSecs:  6,
			EffectPower:   80,
			AppliesBuff:   "Protected",
			IconURL:       "/assets/abilities/tank/shield_wall.png",
			AnimationName: "barrier_wall",
			SoundEffect:   "shield_wall.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 25
		{
			Name:          "Counter Strike",
			Description:   "Enter a defensive stance. Next attack that hits you is blocked and countered for 200% damage.",
			Class:         "Tank",
			UnlockLevel:   25,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      20,
			ManaCost:      45,
			BaseDamage:    200,
			DurationSecs:  3,
			AppliesBuff:   "Counter Ready",
			IconURL:       "/assets/abilities/tank/counter.png",
			AnimationName: "parry_stance",
			SoundEffect:   "counter.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 30
		{
			Name:          "Guardian Aura",
			Description:   "Emit a protective aura, granting all allies 40% damage reduction for 8 seconds.",
			Class:         "Tank",
			UnlockLevel:   30,
			AbilityType:   "ACTIVE",
			TargetType:    "ALL_ALLIES",
			Cooldown:      35,
			ManaCost:      70,
			DurationSecs:  8,
			EffectPower:   40,
			AppliesBuff:   "Guardian's Blessing",
			IconURL:       "/assets/abilities/tank/guardian_aura.png",
			AnimationName: "holy_circle",
			SoundEffect:   "guardian.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 40
		{
			Name:          "Immovable Object",
			Description:   "Become immune to knockback, stuns, and slows. Gain 60% defense for 10 seconds.",
			Class:         "Tank",
			UnlockLevel:   40,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      40,
			ManaCost:      80,
			DurationSecs:  10,
			EffectPower:   60,
			AppliesBuff:   "Immovable",
			IconURL:       "/assets/abilities/tank/immovable.png",
			AnimationName: "stone_skin",
			SoundEffect:   "mountain.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 50
		{
			Name:          "Sacrifice",
			Description:   "Transfer all damage taken by allies to yourself for 5 seconds. Heal 10% max HP per second while active.",
			Class:         "Tank",
			UnlockLevel:   50,
			AbilityType:   "ACTIVE",
			TargetType:    "ALL_ALLIES",
			Cooldown:      60,
			ManaCost:      100,
			DurationSecs:  5,
			EffectPower:   10,
			AppliesBuff:   "Martyr",
			IconURL:       "/assets/abilities/tank/sacrifice.png",
			AnimationName: "divine_link",
			SoundEffect:   "sacrifice.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 60
		{
			Name:          "Fortress",
			Description:   "Transform into an unbreakable fortress, becoming stationary but gaining 90% damage reduction and taunting all enemies.",
			Class:         "Tank",
			UnlockLevel:   60,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      50,
			ManaCost:      120,
			DurationSecs:  8,
			EffectPower:   90,
			AppliesBuff:   "Fortress Mode",
			IconURL:       "/assets/abilities/tank/fortress.png",
			AnimationName: "castle_form",
			SoundEffect:   "fortress.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 80
		{
			Name:          "Titan's Endurance",
			Description:   "Regenerate 5% max HP per second and become immune to debuffs for 12 seconds.",
			Class:         "Tank",
			UnlockLevel:   80,
			AbilityType:   "ACTIVE",
			TargetType:    "SELF",
			Cooldown:      70,
			ManaCost:      150,
			DurationSecs:  12,
			EffectPower:   5,
			AppliesBuff:   "Titan Regen",
			IconURL:       "/assets/abilities/tank/titan_endurance.png",
			AnimationName: "titan_glow",
			SoundEffect:   "titan_roar.mp3",
			ElementBonuses: string(elementBonusTank),
		},
		// Level 100 - ULTIMATE
		{
			Name:          "Immortal Bastion",
			Description:   "Become completely invulnerable for 6 seconds. All allies gain 75% damage reduction. Cannot be dispelled.",
			Class:         "Tank",
			UnlockLevel:   100,
			AbilityType:   "ULTIMATE",
			TargetType:    "ALL_ALLIES",
			Cooldown:      180,
			ManaCost:      200,
			DurationSecs:  6,
			EffectPower:   75,
			AppliesBuff:   "Invulnerable",
			IconURL:       "/assets/abilities/tank/immortal_bastion.png",
			AnimationName: "divine_shield",
			SoundEffect:   "immortality.mp3",
			ElementBonuses: string(elementBonusTank),
		},
	}

	for _, ability := range abilities {
		if err := db.DB.FirstOrCreate(&ability, models.Ability{Name: ability.Name, Class: ability.Class}).Error; err != nil {
			log.Printf("    ⚠️  Error seeding %s: %v", ability.Name, err)
		} else {
			log.Printf("    ✓ Seeded: Lvl %d - %s", ability.UnlockLevel, ability.Name)
		}
	}
}
