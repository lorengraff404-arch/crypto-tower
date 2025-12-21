package main

import (
	"fmt"
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
	"github.com/lorengraff/crypto-tower-defense/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()
	fmt.Println("🚀 Starting Marketplace Cleanup...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	if err := db.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("🧹 cleaning up ORPHANED listings (referencing missing characters)...")
	// Delete listings where character_id points to a non-existent (or soft-deleted) character
	result := db.DB.Exec(`
		DELETE FROM marketplace_listings 
		WHERE character_id IS NOT NULL 
		AND character_id NOT IN (SELECT id FROM characters WHERE deleted_at IS NULL)
	`)
	if result.Error != nil {
		log.Fatal("Error cleaning orphans:", result.Error)
	}
	fmt.Printf("✅ Deleted %d orphaned listing(s)\n", result.RowsAffected)

	fmt.Println("🧹 cleaning up DUPLICATE listings (keeping oldest)...")
	// Delete duplicates, keeping the one with MIN(id)
	result = db.DB.Exec(`
		DELETE FROM marketplace_listings
		WHERE status = 'ACTIVE'
		AND character_id IS NOT NULL
		AND id NOT IN (
			SELECT MIN(id)
			FROM marketplace_listings
			WHERE status = 'ACTIVE' 
			AND character_id IS NOT NULL
			GROUP BY character_id
		)
	`)
	if result.Error != nil {
		log.Fatal("Error cleaning duplicates:", result.Error)
	}
	fmt.Printf("✅ Deleted %d duplicate listing(s)\n", result.RowsAffected)

	fmt.Println("✨ Cleanup complete! The marketplace is now clean.")

	// DEBUG: Investigate why remaining listings fail to load character
	fmt.Println("🔍 Investigating remaining listings...")
	var listings []models.MarketplaceListing
	if err := db.DB.Preload("Character").Find(&listings).Error; err != nil {
		log.Println("Error finding listings:", err)
	}

	for _, l := range listings {
		fmt.Printf("Listing %d (CharID: %v)\n", l.ID, l.CharacterID)
		if l.Character == nil && l.CharacterID != nil {
			fmt.Println("  ⚠️ Character PRELOAD failed. Trying manual load...")
			var c models.Character
			err := db.DB.Unscoped().First(&c, *l.CharacterID).Error
			fmt.Printf("  -> Manual Load Result: %v (Error: %v)\n", c.ID, err)
			if err == nil {
				fmt.Printf("  -> Character exists but Preload failed. DeletedAt: %v\n", c.DeletedAt)
			}
		} else if l.Character != nil {
			c := l.Character
			fmt.Printf("  ✅ Loaded via Preload. Name='%s' Atk=%d Def=%d HP=%d Rarity=%s\n",
				c.Name, c.CurrentAttack, c.CurrentDefense, c.CurrentHP, c.Rarity)
		}
	}
}
