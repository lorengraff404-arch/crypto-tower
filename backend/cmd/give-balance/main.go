package main

import (
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
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

	log.Println("🎁 Giving starter bonuses to all users...")

	// Give 1000 GTK and 500 TOWER to all users
	result := db.DB.Exec("UPDATE users SET gtk_balance = 1000, tower_balance = 500 WHERE deleted_at IS NULL")

	if result.Error != nil {
		log.Fatalf("Failed to update users: %v", result.Error)
	}

	log.Printf("✅ Updated %d users with starter bonuses!", result.RowsAffected)
	log.Println("   💰 +1000 GTK")
	log.Println("   🏆 +500 TOWER")
}
