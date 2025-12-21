package main

import (
	"fmt"
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed config: %v", err)
	}

	if err := db.Connect(cfg); err != nil {
		log.Fatalf("Failed db: %v", err)
	}

	log.Println("🔄 Migrating Raid tables...")

	// DROP old tables to handle schema changes cleanly in dev
	db.DB.Migrator().DropTable(&models.IslandMission{}, &models.UserCampaignProgress{}, &models.RaidSession{}, &models.RaidBoss{}, &models.Island{})

	if err := db.DB.AutoMigrate(&models.Island{}, &models.RaidBoss{}, &models.IslandMission{}, &models.UserCampaignProgress{}, &models.RaidSession{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("🌱 Seeding Islands & Missions...")
	seedIslandsAndMissions()

	log.Println("✅ Raid system ready!")
}

func seedIslandsAndMissions() {
	islands := []models.Island{
		{Name: "Volcanic Wasteland", Description: "A scorching realm of fire and ash.", Difficulty: 1, MinLevelReq: 1, ImageURL: "img/islands/volcanic.jpg"},
		{Name: "Frozen Spire", Description: "An icy mountain peak shrouded in blizzards.", Difficulty: 2, MinLevelReq: 25, ImageURL: "img/islands/frozen.jpg"},
		{Name: "Emerald Jungle", Description: "Dense foliage hides lethal predators.", Difficulty: 3, MinLevelReq: 50, ImageURL: "img/islands/jungle.jpg"},
	}

	// Create Islands first
	for i := range islands {
		island := &islands[i]
		var existingIsland models.Island
		err := db.DB.Where("name = ?", island.Name).First(&existingIsland).Error

		if err != nil && err == gorm.ErrRecordNotFound {
			// Create island with boss
			boss := models.RaidBoss{
				Name:          fmt.Sprintf("%s Guardian", island.Name),
				Element:       "FIRE",
				CharacterType: "DRAGON",
				TotalHP:       int64(1500 * island.Difficulty), // 1500 HP = ~25 turns for team of 4 doing 60dmg/turn
				BaseAttack:    45 * island.Difficulty,          // Stronger than mobs
				BaseDefense:   30 * island.Difficulty,
				Speed:         80,
				ImageURL:      "",
			}
			island.Bosses = []models.RaidBoss{boss}

			if err := db.DB.Create(island).Error; err != nil {
				log.Printf("Failed to seed %s: %v", island.Name, err)
			} else {
				log.Printf("Created island: %s (ID: %d)", island.Name, island.ID)
			}
		} else if err == nil {
			islands[i].ID = existingIsland.ID
			log.Printf("Island %s already exists (ID: %d)", island.Name, existingIsland.ID)
		}
	}

	// Now create missions
	for _, island := range islands {
		if island.ID == 0 {
			continue
		}

		var count int64
		db.DB.Model(&models.IslandMission{}).Where("island_id = ?", island.ID).Count(&count)
		if count > 0 {
			log.Printf("Missions for %s already exist, skipping", island.Name)
			continue
		}

		for i := 1; i <= 5; i++ {
			// REBALANCED V3 (Phase 11): Stat Squish Match
			// Player Team: ~400 Total HP (4x100), ~15 Avg Dmg per hit (60 Dmg/Turn total)

			// Enemy HP: Target 4-6 turns for normal mobs
			// 4 turns * 60 dmg = 240 HP
			baseHP := 200 + (i-1)*50 // 200, 250, 300, 350, 400

			// Enemy ATK: Target ~10-15 dmg per hit to player (Player has ~25 Def)
			// (Atk * Pwr / Def) -> We want result ~15
			// (X * 40 / 25) / 15 = ~10 -> X ~= 25-30
			baseAtk := 25 + (i-1)*5 // 25, 30, 35, 40, 45

			// XP & Reward Tuning
			tokens := 20 + i*10
			xp := 50 + i*20
			// Enemy DEF: Team should still do decent damage
			baseDef := 5 + (i-1)*3 // 5, 8, 11, 14, 17

			// Determine enemy type based on island
			enemyType := "NORMAL"
			switch island.Name {
			case "Volcanic Wasteland":
				enemyType = "FIRE"
			case "Frozen Spire":
				enemyType = "ICE"
			case "Emerald Jungle":
				enemyType = "GRASS"
			}

			mission := models.IslandMission{
				IslandID:    island.ID,
				Sequence:    i,
				Name:        fmt.Sprintf("%s - Mission %d", island.Name, i),
				Description: fmt.Sprintf("Stage %d of the expedition.", i),
				EnemyName:   fmt.Sprintf("Minion %d", i),
				EnemyType:   enemyType,
				EnemyHP:     int64(baseHP * island.Difficulty),
				EnemyAtk:    baseAtk * island.Difficulty,
				EnemyDef:    baseDef * island.Difficulty,
				EnemySpeed:  10 * i,
				EnemyImage:  "",
				RewardsPool: fmt.Sprintf(`{"tokens": %d, "xp": %d}`, tokens, xp),
			}
			if i == 5 {
				// Boss: 3x stronger than mission 4 (real challenge)
				mission.Name = fmt.Sprintf("%s - BOSS BATTLE", island.Name)
				mission.EnemyName = fmt.Sprintf("%s Guardian", island.Name)
				mission.EnemyHP = int64(800 * island.Difficulty) // Boss HP (requires full team effort)
				mission.EnemyAtk = 70 * island.Difficulty        // Can hurt but not instant wipe
				mission.EnemyDef = 25 * island.Difficulty
			}
			if err := db.DB.Create(&mission).Error; err != nil {
				log.Printf("Failed to create mission %d for %s: %v", i, island.Name, err)
			} else {
				log.Printf("Created mission: %s (HP: %d, ATK: %d, DEF: %d, TYPE: %s)", mission.Name, mission.EnemyHP, mission.EnemyAtk, mission.EnemyDef, mission.EnemyType)
			}
		}
	}

	log.Println("✅ Missions seeded!")
}
