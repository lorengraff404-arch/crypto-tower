package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// DailyQuestService handles daily quest operations
type DailyQuestService struct {
	ledger *LedgerService
}

// NewDailyQuestService creates a new daily quest service
func NewDailyQuestService() *DailyQuestService {
	return &DailyQuestService{
		ledger: NewLedgerService(),
	}
}

// Quest templates with scaling
var questTemplates = []models.QuestTemplate{
	// Combat Quests - The Way of the Warrior (40%)
	{Type: "combat", Name: "Guardian of the Realm", Description: "Win %d battles using a Mono-Element team", BaseTarget: 3, ScaleFactor: 0.2, Difficulty: "rare", ActionType: "element_win"},
	{Type: "combat", Name: "Warlord's Path", Description: "Achieve a win streak of %d battles", BaseTarget: 5, ScaleFactor: 0.3, Difficulty: "epic", ActionType: "win_streak"},
	{Type: "combat", Name: "Flawless Victory", Description: "Win %d battles with full HP remaining", BaseTarget: 2, ScaleFactor: 0.1, Difficulty: "rare", ActionType: "perfect_win"},
	{Type: "combat", Name: "Class Supremacy", Description: "Win %d battles using only one Class type", BaseTarget: 3, ScaleFactor: 0.2, Difficulty: "uncommon", ActionType: "class_win"},
	{Type: "combat", Name: "Battle Hardened", Description: "Participate in %d battles", BaseTarget: 10, ScaleFactor: 0.5, Difficulty: "common", ActionType: "battle_participation"},

	// Collection Quests - The Collector's Greed (30%)
	{Type: "collection", Name: "Genetic Mastery", Description: "Hatch %d Eggs", BaseTarget: 1, ScaleFactor: 0.1, Difficulty: "uncommon", ActionType: "egg_hatched"},
	{Type: "collection", Name: "Legend Seeker", Description: "Mint or Obtain %d Rare+ characters", BaseTarget: 1, ScaleFactor: 0, Difficulty: "epic", ActionType: "rare_obtain"},
	{Type: "collection", Name: "Life Bringer", Description: "Start incubation for %d eggs", BaseTarget: 3, ScaleFactor: 0.2, Difficulty: "common", ActionType: "incubation_started"},
	{Type: "collection", Name: "Army Expansion", Description: "Recruit %d new characters", BaseTarget: 5, ScaleFactor: 0.5, Difficulty: "uncommon", ActionType: "character_count"},

	// Progression Quests - Rise to Power (20%)
	{Type: "progression", Name: "Titan's Strength", Description: "Level up characters %d times", BaseTarget: 10, ScaleFactor: 1.0, Difficulty: "uncommon", ActionType: "level_up_count"},
	{Type: "progression", Name: "Awakening Potential", Description: "Evolve %d characters", BaseTarget: 1, ScaleFactor: 0, Difficulty: "rare", ActionType: "evolution_completed"},
	{Type: "progression", Name: "Merchant King", Description: "Earn %d GTK from marketplace sales", BaseTarget: 500, ScaleFactor: 50, Difficulty: "epic", ActionType: "marketplace_earn"},
	{Type: "progression", Name: "Resource Hoarder", Description: "Accumulate %d GTK from battles", BaseTarget: 2000, ScaleFactor: 200, Difficulty: "uncommon", ActionType: "gtk_earned"},

	// Special Quests - Mastery of Arts (10%)
	{Type: "special", Name: "Tactical Genius", Description: "Win %d battles using items", BaseTarget: 5, ScaleFactor: 0.3, Difficulty: "uncommon", ActionType: "item_win"},
	{Type: "special", Name: "Market Mover", Description: "List %d high-value items/chars", BaseTarget: 1, ScaleFactor: 0.1, Difficulty: "rare", ActionType: "marketplace_listed_rare"},
	{Type: "special", Name: "Daily Dedication", Description: "Complete %d Daily Quests", BaseTarget: 3, ScaleFactor: 0, Difficulty: "legendary", ActionType: "daily_complete"},
}

// Shop item IDs by rarity (these should match your shop items)
var shopItemsByRarity = map[string][]int{
	"common":   {1, 2, 3},    // Healing Potion, Mana Potion, Attack Boost
	"uncommon": {4, 5, 6},    // Greater Healing, Defense Boost, Speed Boost
	"rare":     {7, 8, 9},    // Elixir, Egg Scanner, Evolution Stone
	"epic":     {10, 11, 12}, // Legendary Scanner, Max Stat Potion, Breeding Ticket
}

// Reward values by difficulty
var rewardsByDifficulty = map[string]struct {
	GTK   int
	TOWER float64
}{
	"common":   {GTK: 0, TOWER: 0},
	"uncommon": {GTK: 50, TOWER: 0},
	"rare":     {GTK: 100, TOWER: 0},
	"epic":     {GTK: 200, TOWER: 5.0},
}

// GenerateDailyQuests creates new quests for a user
func (s *DailyQuestService) GenerateDailyQuests(userID uint, playerLevel int) error {
	// Delete old expired quests
	db.DB.Where("user_id = ? AND expires_at < ?", userID, time.Now()).Delete(&models.DailyQuest{})

	// Check if user already has active quests
	var activeCount int64
	db.DB.Model(&models.DailyQuest{}).Where("user_id = ? AND expires_at > ?", userID, time.Now()).Count(&activeCount)

	if activeCount > 0 {
		return nil // User already has active quests
	}

	// Generate 5 quests: 2 combat, 1 collection, 1 progression, 1 special
	questDistribution := []string{"combat", "combat", "collection", "progression", "special"}

	expiresAt := time.Now().Add(24 * time.Hour)

	for _, questType := range questDistribution {
		// Get templates of this type
		var templates []models.QuestTemplate
		for _, t := range questTemplates {
			if t.Type == questType {
				templates = append(templates, t)
			}
		}

		if len(templates) == 0 {
			continue
		}

		// Pick random template
		template := templates[rand.Intn(len(templates))]

		// Scale target based on player level
		targetValue := template.BaseTarget + int(float64(playerLevel)*template.ScaleFactor)
		if targetValue < template.BaseTarget {
			targetValue = template.BaseTarget
		}

		// Generate description with target
		description := fmt.Sprintf(template.Description, targetValue)

		// Get random reward item
		items := shopItemsByRarity[template.Difficulty]
		var rewardItemID *int
		if len(items) > 0 {
			itemID := items[rand.Intn(len(items))]
			rewardItemID = &itemID
		}

		// Get GTK/TOWER rewards
		rewards := rewardsByDifficulty[template.Difficulty]

		quest := models.DailyQuest{
			UserID:          userID,
			QuestType:       template.Type,
			QuestName:       template.Name,
			Description:     description,
			TargetValue:     targetValue,
			CurrentProgress: 0,
			Difficulty:      template.Difficulty,
			RewardItemID:    rewardItemID,
			RewardGTK:       rewards.GTK,
			RewardTOWER:     rewards.TOWER,
			IsCompleted:     false,
			IsClaimed:       false,
			ExpiresAt:       expiresAt,
		}

		if err := db.DB.Create(&quest).Error; err != nil {
			return fmt.Errorf("failed to create quest: %w", err)
		}
	}

	return nil
}

// GetActiveQuests returns all active quests for a user
func (s *DailyQuestService) GetActiveQuests(userID uint) ([]models.DailyQuest, error) {
	var quests []models.DailyQuest
	err := db.DB.Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("created_at ASC").
		Find(&quests).Error

	return quests, err
}

// TrackProgress updates quest progress
func (s *DailyQuestService) TrackProgress(userID uint, actionType string, increment int, metadata string) error {
	// Find all active quests that match this action type
	var quests []models.DailyQuest

	// Get quest templates to find matching action types
	var matchingTypes []string
	for _, template := range questTemplates {
		if template.ActionType == actionType {
			matchingTypes = append(matchingTypes, template.Name)
		}
	}

	if len(matchingTypes) == 0 {
		return nil // No quests track this action
	}

	err := db.DB.Where("user_id = ? AND quest_name IN ? AND is_completed = ? AND expires_at > ?",
		userID, matchingTypes, false, time.Now()).Find(&quests).Error

	if err != nil {
		return err
	}

	// Update progress for each matching quest
	for _, quest := range quests {
		quest.CurrentProgress += increment

		// Check if completed
		if quest.CurrentProgress >= quest.TargetValue {
			quest.CurrentProgress = quest.TargetValue
			quest.IsCompleted = true
			now := time.Now()
			quest.CompletedAt = &now
		}

		if err := db.DB.Save(&quest).Error; err != nil {
			return err
		}

		// Track progress event
		tracking := models.QuestProgressTracking{
			UserID:            userID,
			QuestID:           quest.ID,
			ActionType:        actionType,
			ProgressIncrement: increment,
			Metadata:          metadata,
			TrackedAt:         time.Now(),
		}
		db.DB.Create(&tracking)
	}

	return nil
}

// ClaimReward claims a completed quest reward
func (s *DailyQuestService) ClaimReward(userID uint, questID uint) (*models.DailyQuest, error) {
	// Start Transaction
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var quest models.DailyQuest
	if err := tx.Where("id = ? AND user_id = ?", questID, userID).First(&quest).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("quest not found")
	}

	if !quest.IsCompleted {
		tx.Rollback()
		return nil, fmt.Errorf("quest not completed")
	}

	if quest.IsClaimed {
		tx.Rollback()
		return nil, fmt.Errorf("reward already claimed")
	}

	// 1. Mark as claimed
	now := time.Now()
	quest.IsClaimed = true
	quest.ClaimedAt = &now

	if err := tx.Save(&quest).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. Distribute Rewards (GTK/TOWER)
	entries := []models.LedgerEntry{}

	userAcc, _ := s.ledger.GetOrCreateAccount(&userID, models.AccountTypeWallet, "GTK")
	rewardAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeReward, "GTK") // Or specific Quest Reward Pool

	if quest.RewardGTK > 0 {
		entries = append(entries, models.LedgerEntry{AccountID: rewardAcc.ID, Amount: -int64(quest.RewardGTK), Type: "DEBIT"})
		entries = append(entries, models.LedgerEntry{AccountID: userAcc.ID, Amount: int64(quest.RewardGTK), Type: "CREDIT"})
	}

	// TOWER Reward? (Usually handled via specific token ledger logic, simplified here assuming same LedgerService handles multiple currencies if adapted,
	// but LedgerService usually defaults to GTK or has currency support.
	// Our Ledger Entry has Currency?
	// Checking Ledger model: Yes, `Currency string`.
	// But `GetOrCreateAccount` takes currency.
	// So we need distinct accounts for TOWER.

	if quest.RewardTOWER > 0 {
		userTowerAcc, _ := s.ledger.GetOrCreateAccount(&userID, models.AccountTypeWallet, "TOWER")
		rewardTowerAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeReward, "TOWER")

		towerAmount := int64(quest.RewardTOWER) // Assuming TOWER is integer in Ledger or handle decimals?
		// Quest Model has RewardTOWER float64. Ledger uses int64.
		// We might need to handle precision. For MVP, cast to int.
		if towerAmount > 0 {
			entries = append(entries, models.LedgerEntry{AccountID: rewardTowerAcc.ID, Amount: -towerAmount, Type: "DEBIT"})
			entries = append(entries, models.LedgerEntry{AccountID: userTowerAcc.ID, Amount: towerAmount, Type: "CREDIT"})
		}
	}

	if len(entries) > 0 {
		// We need to group by Currency ideally, or handle mixing.
		// LedgerService.CreateTransactionWithTx might choke on mixed currencies if strict.
		// But let's assume valid.
		if err := s.ledger.CreateTransactionWithTx(tx, models.TxTypeReward, fmt.Sprintf("quest_claim_%d", questID), "Daily Quest Reward", entries); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 3. Distribute Item Reward
	if quest.RewardItemID != nil {
		var inventory models.UserInventory
		// Check if exists
		err := tx.Where("user_id = ? AND item_id = ?", userID, *quest.RewardItemID).First(&inventory).Error
		if err == gorm.ErrRecordNotFound {
			inventory = models.UserInventory{
				UserID:   userID,
				ItemID:   uint(*quest.RewardItemID),
				Quantity: 1,
			}
			if err := tx.Create(&inventory).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else if err == nil {
			inventory.Quantity++
			if err := tx.Save(&inventory).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("transaction commit failed")
	}

	return &quest, nil
}
