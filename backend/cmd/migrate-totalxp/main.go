package main

import (
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

	log.Println("🔄 Running database migration for TotalXP column...")

	// Auto-migrate the Character model (will add total_xp column)
	if err := db.DB.AutoMigrate(&models.Character{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	log.Println("✅ Migration complete! TotalXP column added.")
	log.Println("📊 Syncing existing characters...")

	// Update all existing characters: set total_xp based on current level
	result := db.DB.Exec(`
		UPDATE characters 
		SET total_xp = FLOOR(100 * POWER(level, 2.5))
		WHERE total_xp = 0 AND level > 1
	`)

	if result.Error != nil {
		log.Printf("⚠️  Warning: Failed to sync total_xp: %v", result.Error)
	} else {
		log.Printf("✅ Synced %d characters with correct total_xp", result.RowsAffected)
	}
}
