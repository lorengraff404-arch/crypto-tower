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

	log.Println("🔄 Migrating Character Type field...")

	// Auto-migrate to add character_type column
	if err := db.DB.AutoMigrate(&models.Character{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	log.Println("✅ Character type migration complete!")

	// Set default types for existing characters
	log.Println("📝 Setting default types for existing characters...")

	result := db.DB.Exec("UPDATE characters SET character_type = 'BEAST' WHERE character_type IS NULL OR character_type = ''")
	if result.Error != nil {
		log.Printf("⚠️  Warning: Failed to set defaults: %v", result.Error)
	} else {
		log.Printf("✅ Updated %d characters with default type", result.RowsAffected)
	}

	log.Println("🎯 Migration complete! Character types ready.")
}
