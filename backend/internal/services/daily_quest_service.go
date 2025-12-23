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

// Quest templates are now stored in the database (QuestTemplate model)

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

	// Generate 5 random quests from DB templates
	var allTemplates []models.QuestTemplate
	if err := db.DB.Find(&allTemplates).Error; err != nil {
		return err
	}

	if len(allTemplates) == 0 {
		return errors.New("no quest templates found in database")
	}

	// Shuffle templates
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allTemplates), func(i, j int) {
		allTemplates[i], allTemplates[j] = allTemplates[j], allTemplates[i]
	})

	// Take up to 5
	count := 5
	if len(allTemplates) < 5 {
		count = len(allTemplates)
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	for i := 0; i < count; i++ {
		template := allTemplates[i]

		// Scale target based on player level
		targetValue := template.BaseTarget + int(float64(playerLevel)*template.ScaleFactor)
		if targetValue < template.BaseTarget {
			targetValue = template.BaseTarget
		}

		// Generate description with target
		description := fmt.Sprintf(template.Description, targetValue)

		// Get GTK/TOWER rewards from template
		rewardGTK := template.RewardGTK
		rewardTOWER := float64(template.RewardTOWER)

		// Get random reward item ID (if difficulty matches)
		items := shopItemsByRarity[template.Difficulty]
		var rewardItemID *int
		if len(items) > 0 {
			itemID := items[rand.Intn(len(items))]
			rewardItemID = &itemID
		}

		quest := models.DailyQuest{
			UserID:          userID,
			QuestType:       template.Type,
			QuestName:       template.Name,
			Description:     description,
			TargetValue:     targetValue,
			CurrentProgress: 0,
			Difficulty:      template.Difficulty,
			RewardItemID:    rewardItemID,
			RewardGTK:       rewardGTK,
			RewardTOWER:     rewardTOWER,
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
	// Find all templates matching this action type to get their names
	var matchingTemplates []models.QuestTemplate
	if err := db.DB.Where("action_type = ?", actionType).Find(&matchingTemplates).Error; err != nil {
		return err
	}

	if len(matchingTemplates) == 0 {
		return nil // No templates track this action
	}

	var matchingNames []string
	for _, t := range matchingTemplates {
		matchingNames = append(matchingNames, t.Name)
	}

	// Find all active quests that match these names
	var quests []models.DailyQuest
	err := db.DB.Where("user_id = ? AND quest_name IN ? AND is_completed = ? AND expires_at > ?",
		userID, matchingNames, false, time.Now()).Find(&quests).Error

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
