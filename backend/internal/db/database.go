package db

import (
	"fmt"
	"log"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes database connection
func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// Migrate runs all database migrations
func Migrate() error {
	// Auto-migrate models
	log.Println("Running database migrations...")
	err := DB.AutoMigrate(
		&models.User{},
		&models.Character{},
		&models.CharacterMove{},
		&models.Battle{},
		&models.BattleState{}, // Added migration
		&models.MarketplaceListing{},
		&models.Transaction{},
		&models.Egg{},
		&models.Equipment{},
		&models.DailyQuest{},
		&models.ShopItem{},
		&models.Item{},
		&models.UserInventory{}, // Added UserInventory
		&models.Mission{},
		&models.Team{},
		&models.TeamMember{},
		&models.Island{},
		&models.RaidBoss{},
		&models.RaidSession{},

		// Skill System (Phase 20)
		&models.Ability{},
		&models.AbilityLearning{},         // NEW: Ability learning requirements
		&models.CharacterLearnedAbility{}, // NEW: Track learned abilities
		&models.CharacterAbility{},
		&models.CharacterActiveSkill{},
		&models.CharacterSkillCooldown{},
		&models.AbilityUsage{},
		&models.CharacterStatusEffect{},
		&models.CharacterBuff{},

		// Sprite System
		&models.SpriteGenerationJob{},

		// Admin Panel
		&models.RevenueTransaction{},
		&models.AdminAction{},
		&models.AntiCheatFlag{},
		&models.GameConfig{},

		// System (Audit, Notifications, etc.)
		&models.AuditLog{},
	)

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migrations completed successfully")

	// Seed missions (tutorial + progression)
	if err := SeedMissions(DB); err != nil {
		log.Printf("Warning: Failed to seed missions: %v", err)
	} else {
		log.Println("Missions seeded successfully")
	}

	// Seed abilities FIRST (required for ability_learning foreign keys)
	if err := SeedAbilities(DB); err != nil {
		log.Printf("Warning: Failed to seed abilities: %v", err)
	} else {
		log.Println("Abilities seeded successfully")
	}

	// Seed ability learning requirements
	if err := SeedAbilityLearning(DB); err != nil {
		log.Printf("Warning: Failed to seed ability learning: %v", err)
	} else {
		log.Println("Ability learning seeded successfully")
	}

	// Seed story dialogues
	if err := SeedDialogues(DB); err != nil {
		log.Printf("Warning: Failed to seed dialogues: %v", err)
	} else {
		log.Println("Story dialogues seeded successfully")
	}

	// Seed story fragments
	if err := SeedStoryFragments(DB); err != nil {
		log.Printf("Warning: Failed to seed story fragments: %v", err)
	} else {
		log.Println("Story fragments seeded successfully")
	}

	// Add indexes for performance
	addIndexes()

	return nil
}

// addIndexes creates additional indexes for query optimization
func addIndexes() {
	// User indexes
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_level ON users(level)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_rank_tier ON users(rank_tier)")

	// Character indexes
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_characters_owner_type ON characters(owner_id, character_type)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_characters_rarity_level ON characters(rarity, level)")

	// Battle indexes
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_battles_created_at ON battles(created_at DESC)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_battles_type_status ON battles(battle_type, status)")

	// Transaction indexes
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_transactions_user_created ON transactions(user_id, created_at DESC)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_transactions_type_created ON transactions(transaction_type, created_at DESC)")

	// Marketplace indexes
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_marketplace_status_price ON marketplace_listings(status, price)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_marketplace_asset_type_status ON marketplace_listings(asset_type, status)")

	log.Println("Additional indexes created successfully")
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
