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

	log.Println("🔄 Migrating Ability tables...")

	// Auto-migrate Ability and CharacterAbility models
	if err := db.DB.AutoMigrate(&models.Ability{}, &models.CharacterAbility{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	log.Println("✅ Ability tables created successfully!")
	log.Println("   📋 abilities - Stores all learnable abilities")
	log.Println("   📋 character_abilities - Tracks learned abilities per character")
}
