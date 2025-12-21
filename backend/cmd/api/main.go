package main

import (
	"fmt"
	"log"

	// Added by instruction
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin" // Added by instruction

	// 	"github.com/lorengraff/crypto-tower-defense/internal/blockchain"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/handlers"
	"github.com/lorengraff/crypto-tower-defense/internal/middleware"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
	"github.com/lorengraff/crypto-tower-defense/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()
	logger.Info("Starting Crypto Tower Defense API Server...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	if err := db.Connect(cfg); err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		log.Fatal(err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		logger.Error(fmt.Sprintf("Failed to run migrations: %v", err))
		log.Fatal(err)
	}

	// Initialize Gin router
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000", "http://localhost:8080",
		"http://127.0.0.1:3000", "http://127.0.0.1:8080",
		"http://127.0.0.1:3000", "http://127.0.0.1:8080",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Cache-Control", "Pragma", "Expires"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Serve admin dashboard static files (ABSOLUTE PATH)
	router.Static("/admin", "./admin-dashboard")

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"service":  "crypto-tower-defense-api",
			"version":  "0.1.0",
			"database": "connected",
		})
	})

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(cfg)
	characterHandler := handlers.NewCharacterHandler()
	missionHandler := handlers.NewMissionHandler(db.DB) // Fixed Signature
	storyHandler := handlers.NewStoryHandler()
	itemHandler := handlers.NewItemHandler()
	// battleHandler := handlers.NewBattleHandler() // Legacy, using battleEngine in other handlers
	progressionHandler := handlers.NewProgressionHandler()

	// Game Modes (Refactored)
	// gameModeHandler := handlers.NewGameModeHandler()

	// Initialize Revenue service
	sqlDB, _ := db.DB.DB()
	revenueService := services.NewRevenueService(sqlDB)
	revenueHandler := handlers.NewRevenueHandler(revenueService)

	// Enterprise Services
	ledgerService := services.NewLedgerService()
	configService := services.NewConfigService()
	adminService := services.NewAdminService()

	adminHandler := handlers.NewAdminHandler(adminService, configService)

	// Inject Ledger into Revenue (TODO: refactor RevenueService to accept Ledger)
	_ = ledgerService
	_ = configService
	_ = adminService
	_ = adminHandler // Prevent unused error

	// Blockchain Service (Fail safe)
	blockchainService, err := services.NewBlockchainService(cfg)
	if err != nil {
		log.Printf("⚠️ WARNING: Blockchain Service failed to initialize: %v (Shop verification will be disabled)", err)
		blockchainService = nil
	} else {
		log.Println("✅ Blockchain Service initialized")
	}

	// Initialize ShopService with Blockchain dependencies
	shopService := services.NewShopService(blockchainService)

	// Inject Ledger into Revenue (TODO: refactor RevenueService to accept Ledger)
	_ = ledgerService
	_ = configService
	_ = adminService
	_ = adminHandler // Prevent unused error

	// API v1 routes group
	v1 := router.Group("/api/v1")
	{
		// Apply security middleware
		v1.Use(middleware.SecurityHeaders())
		v1.Use(middleware.MaxBodySizeMiddleware(10 * 1024 * 1024)) // 10MB max
		v1.Use(middleware.ValidateJSONMiddleware())

		// Diagnostic endpoint
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "pong", "v1": "active"})
		})

		// Public auth routes with strict rate limiting
		authRoutes := v1.Group("/auth")
		authRoutes.Use(middleware.StrictRateLimiter()) // 10 req/min
		{
			authRoutes.POST("/nonce", authHandler.GetNonce)
			authRoutes.POST("/verify", authHandler.VerifySignature)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		protected.Use(middleware.StandardRateLimiter()) // 100 req/min for normal operations
		{
			// Auth
			protected.GET("/auth/profile", authHandler.GetProfile)

			// Character routes
			protected.GET("/characters", characterHandler.ListCharacters)
			protected.GET("/characters/:id", characterHandler.GetCharacter)
			protected.POST("/characters", characterHandler.CreateCharacter)
			protected.POST("/characters/:id/hatch", characterHandler.HatchEgg)
			protected.POST("/characters/:id/recover", characterHandler.RecoverFatigue)

			// Progression routes
			protected.POST("/characters/:id/gain-xp", progressionHandler.GainXP)
			protected.GET("/characters/:id/progression", progressionHandler.GetProgressionInfo)
			protected.GET("/characters/:id/validate", progressionHandler.ValidateIntegrity)

			// Ability routes
			abilityHandler := handlers.NewAbilityHandler()
			protected.GET("/abilities", abilityHandler.GetAbilitiesByClass)
			protected.GET("/abilities/:id", abilityHandler.GetAbilityDetails)
			protected.GET("/characters/:id/abilities", abilityHandler.GetCharacterAbilities)

			// Skill System routes (Phase 20)
			skillHandler := handlers.NewSkillHandler()
			protected.POST("/skills/activate", skillHandler.ActivateSkill)
			protected.GET("/characters/:id/active-skills", skillHandler.GetActiveSkills)
			protected.POST("/characters/:id/swap-skill", skillHandler.SwapSkill)
			protected.GET("/characters/:id/cooldowns", skillHandler.GetSkillCooldowns)
			// protected.POST("/battles/:id/start-turn", skillHandler.StartTurn)

			// Status Effect routes
			statusEffectHandler := handlers.NewStatusEffectHandler()
			protected.GET("/characters/:id/effects", statusEffectHandler.GetCharacterEffects)
			protected.GET("/effects/definitions", statusEffectHandler.GetAllEffectDefinitions)

			// Team routes
			teamHandler := handlers.NewTeamHandler()
			protected.GET("/teams", teamHandler.GetMyTeams)
			protected.GET("/teams/active", teamHandler.GetActiveTeam) // Added Active Team route
			protected.POST("/teams", teamHandler.CreateTeam)
			protected.POST("/teams/:id/members", teamHandler.UpdateTeamMember)
			protected.DELETE("/teams/:id/members/:charId", teamHandler.RemoveTeamMember)

			// Raid routes - NEW secure implementation
			raidHandler := handlers.NewRaidHandler()
			// We can reuse raidHandler if we put methods there, but currently split.
			// Let's create RaidTurnHandler sharing the same service?
			// Ideally refactor handlers to be one struct.
			// For minimal change, init turn handler.

			// Hack: RaidHandler in raid_handler.go creates a NEW service instance internally.
			// RaidTurnHandler needs one too.
			// Better: Expose the service from RaidHandler or create shared service.
			// services.NewRaidService()

			raidService := services.NewRaidService()
			// Recreate RaidHandler with injection?
			// Existing NewRaidHandler() creates generic. Let's rely on that for now or refactor.
			// But NewRaidHandler doesn't take args.
			// Creating both handlers separately means separate service instances.
			// RaidService is stateless except for DB/Ledger, so it's SAFE to have 2 instances.

			raidTurnHandler := handlers.NewRaidTurnHandler(raidService)

			raids := protected.Group("/raids")
			{
				raids.POST("/start", raidHandler.StartRaid)
				raids.GET("/:sessionId/status", raidHandler.GetRaidStatus)
				raids.POST("/:sessionId/turn", raidTurnHandler.ProcessRaidTurn)
				raids.POST("/:sessionId/complete", raidTurnHandler.CompleteRaid)
				raids.POST("/:sessionId/flee", raidTurnHandler.FleeRaid)
				raids.GET("/:sessionId/state", handlers.GetRaidBattleState) // Legacy? Check implementation
			}

			// Ranked Battle Endpoints (NEW - ELO matchmaking)
			rankedHandler := handlers.NewRankedHandler()
			ranked := protected.Group("/battle")
			{
				ranked.POST("/ranked", rankedHandler.StartRanked)
			}

			// Wager Battle Endpoints (NEW - Real GTK)
			wagerHandler := handlers.NewWagerHandler()
			wager := protected.Group("/battle")
			wager.GET("/wager-preview", handlers.GetWagerPreview)

			{
				wager.POST("/wager", wagerHandler.StartWager)
			}

			// Mission routes
			missionsGroup := protected.Group("/missions")
			{
				// missionsGroup.GET("", missionHandler.ListMissions)
				// missionsGroup.GET("/current", missionHandler.GetCurrentMission)
				// missionsGroup.POST("/:id/start", missionHandler.StartMission)
				missionsGroup.GET("/:id/progress", missionHandler.GetMissionProgress)
				missionsGroup.POST("/:id/complete", missionHandler.CompleteMission)
			}

			// User features
			protected.GET("/user/unlocks", missionHandler.GetUnlockedFeatures)

			// Story/Narrative endpoints
			protected.GET("/story/missions/:level/dialogues", storyHandler.GetMissionDialogues)
			protected.GET("/story/fragments", storyHandler.GetAvailableFragments)
			protected.GET("/story/progress", storyHandler.GetStoryProgress)
			protected.POST("/story/choices", storyHandler.RecordChoice)
			protected.PUT("/story/progress", storyHandler.UpdateStoryProgress)

			// Item routes
			protected.GET("/items", itemHandler.ListItems)
			protected.GET("/items/:id", itemHandler.GetItem)
			protected.POST("/items", itemHandler.CreateItem)
			protected.POST("/items/:id/equip", itemHandler.EquipItem)
			protected.POST("/items/:id/unequip", itemHandler.UnequipItem)
			protected.POST("/items/:id/use", itemHandler.UseItem)

			// Battle routes
			battleHandler := handlers.NewBattleHandler() // NEW
			protected.POST("/battles/:id/turn", battleHandler.ProcessTurn)
			protected.POST("/battles/:id/complete", battleHandler.CompleteBattle)

			// Game Modes (Refactored to Raids/Ranked/Wager below)
			// protected.GET("/game-modes", gameModeHandler.GetAllGameModes)
			// protected.GET("/game-modes/:mode/info", gameModeHandler.GetGameModeInfo)
			// protected.POST("/game-modes/free/create", gameModeHandler.CreateFreeBattle)
			// protected.POST("/game-modes/ranked/create", gameModeHandler.CreateRankedBattle)
			// protected.POST("/game-modes/ranked/:id/update-elo", gameModeHandler.UpdateRankedELO)
			// protected.POST("/game-modes/wager/find-opponent", gameModeHandler.FindWagerOpponent)
			// protected.POST("/game-modes/wager/create", gameModeHandler.CreateWagerBattle)

			// Breeding Routes
			breedingService := services.NewBreedingService(blockchainService)
			breedingHandler := handlers.NewBreedingHandler(breedingService)
			protected.POST("/breeding/start", breedingHandler.StartBreeding)
			protected.GET("/breeding/eggs", breedingHandler.GetUserEggs)
			protected.POST("/breeding/incubate/:id", breedingHandler.StartIncubation)
			protected.POST("/breeding/hatch/:id", breedingHandler.HatchEgg)

			// Marketplace (Phase 19)
			marketplaceHandler := handlers.NewMarketplaceHandler()
			protected.GET("/marketplace", marketplaceHandler.GetListings)
			protected.POST("/marketplace/list", marketplaceHandler.CreateListing)
			protected.POST("/marketplace/:id/buy", marketplaceHandler.BuyListing)
			protected.DELETE("/marketplace/:id", marketplaceHandler.CancelListing)

			// Inventory (Phase 16)
			inventoryHandler := handlers.NewInventoryHandler()
			protected.GET("/inventory", inventoryHandler.GetInventory)
			protected.POST("/inventory/:id/use", inventoryHandler.UseItem)
			protected.DELETE("/inventory/:id/discard", inventoryHandler.DiscardItem)

			// Economy routes (Token conversion & withdrawal)
			economyHandler := handlers.NewEconomyHandler()
			protected.POST("/economy/convert/tower-to-gtk", economyHandler.ConvertTowerToGTK)
			protected.POST("/economy/convert/gtk-to-tower", economyHandler.ConvertGTKToTower)
			protected.POST("/economy/withdraw", economyHandler.WithdrawTower)
			protected.POST("/economy/deposit", economyHandler.DepositTower)
			protected.GET("/economy/balance", economyHandler.GetBalance)
			protected.GET("/economy/history", economyHandler.GetTransactionHistory)

			// NFT routes (Character minting)
			nftHandler, err := handlers.NewNFTHandler(cfg)
			if err != nil {
				log.Printf("Warning: NFT handler initialization failed: %v", err)
			} else {
				protected.POST("/nft/mint-first", nftHandler.MintFirstCharacter)
				protected.POST("/nft/mint", nftHandler.MintCharacter)
				protected.GET("/nft/gas-estimate", nftHandler.GetGasEstimate)
				protected.GET("/nft/contract-info", nftHandler.GetContractInfo)
				protected.POST("/nft/build-mint-tx", nftHandler.BuildMintTx)
				protected.GET("/nft/verify/:tokenId", nftHandler.VerifyOwnership)
			}

			// Gacha routes
			// blockchainService is passed (may be nil if init failed, handled gracefully in service)
			gachaHandler := handlers.NewGachaHandler(blockchainService, nil)
			protected.POST("/gacha/mint", gachaHandler.MintEgg)
			protected.GET("/gacha/odds/:amount", gachaHandler.GetOddsPreview)
			protected.GET("/gacha/my-eggs", gachaHandler.GetMyEggs)
			protected.POST("/gacha/start-incubation/:id", gachaHandler.StartIncubation)
			protected.POST("/gacha/hatch/:id", gachaHandler.HatchEgg)
			protected.POST("/gacha/scan-egg/:id", gachaHandler.ScanEgg)
			protected.POST("/gacha/apply-accelerator", gachaHandler.ApplyAccelerator)

			// Shop routes (Item shop system)
			shopHandler := handlers.NewShopHandler(shopService) // Modified constructor
			protected.GET("/shop/items", shopHandler.GetShopItems)
			protected.GET("/shop/items/:category", shopHandler.GetShopItems)
			protected.POST("/shop/buy", shopHandler.BuyItem)
			protected.GET("/shop/inventory", shopHandler.GetInventory)
			protected.POST("/shop/use", shopHandler.UseItem)

			// Daily Quests
			questHandler := handlers.NewDailyQuestHandler()
			protected.GET("/daily-quests", questHandler.GetDailyQuests)
			protected.POST("/daily-quests/claim/:id", questHandler.ClaimQuestReward)
			protected.POST("/daily-quests/refresh", questHandler.RefreshQuests) // Admin only
		}

		// Revenue routes (admin only)
		v1.GET("/revenue/stats", middleware.AuthMiddleware(cfg), middleware.AdminMiddleware(cfg), revenueHandler.GetRevenueStats)

		// Admin Routes (Enterprise)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))
		admin.Use(middleware.AdminMiddleware(cfg))
		{
			admin.POST("/ban", adminHandler.BanUser)
			admin.GET("/settings", adminHandler.GetSettings)   // NEW
			admin.PUT("/settings", adminHandler.UpdateSetting) // NEW
			admin.GET("/revenue/stats", adminHandler.GetRevenueStats)
		}
	}

	// Start server
	address := fmt.Sprintf(":%s", cfg.Port)
	logger.Info(fmt.Sprintf("Server starting on %s", address))

	if err := router.Run(address); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		log.Fatal(err)
	}
}
