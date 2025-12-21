package main

import (
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

func main() {
	log.Println("Loading config...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	log.Println("Connecting to database...")
	if err := db.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Seeding abilities...")
	if err := db.SeedAbilities(db.DB); err != nil {
		log.Fatalf("Failed to seed abilities: %v", err)
	}

	log.Println("Seeding ability learnings...")
	if err := db.SeedAbilityLearning(db.DB); err != nil {
		log.Fatalf("Failed to seed ability learnings: %v", err)
	}

	log.Println("✅ Database seeding completed successfully!")
}
