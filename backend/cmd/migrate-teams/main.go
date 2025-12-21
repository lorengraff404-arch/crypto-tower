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

	log.Println("🔄 Migrating Team tables...")

	// Auto-migrate
	if err := db.DB.AutoMigrate(&models.Team{}, &models.TeamMember{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	log.Println("✅ Team system migration complete!")
	log.Println("📊 Tables created: teams, team_members")
}
