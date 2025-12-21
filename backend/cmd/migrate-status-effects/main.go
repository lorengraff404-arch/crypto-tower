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

	log.Println("🔄 Migrating StatusEffect table...")

	// Auto-migrate
	if err := db.DB.AutoMigrate(&models.StatusEffect{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	log.Println("✅ StatusEffect migration complete!")
	log.Println("📊 Status effects ready:")
	log.Println("  - 6 Buffs: Amped, Bulked, Haste, Fleet, Warded, Regen")
	log.Println("  - 12 Debuffs: Burn, Poison, Bleed, Slow, Feeble, Fragile, Stun, Sleep, Freeze, Paralyze, Silence, Blind")
}
