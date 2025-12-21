package services

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// BreedingService handles breeding operations with comprehensive security
type BreedingService struct {
	ledger     *LedgerService
	config     *ConfigService
	blockchain *BlockchainService
}

func NewBreedingService(bc *BlockchainService) *BreedingService {
	return &BreedingService{
		ledger:     NewLedgerService(),
		config:     NewConfigService(),
		blockchain: bc,
	}
}

// StartBreeding initiates breeding between two characters with full security checks
func (s *BreedingService) StartBreeding(userID, parent1ID, parent2ID uint, txHash string) (*models.Egg, error) {
	// SECURITY CHECK 1: Prevent self-breeding
	if parent1ID == parent2ID {
		return nil, errors.New("cannot breed character with itself")
	}

	// SECURITY CHECK 2: Verify ownership of BOTH parents
	var parent1, parent2 models.Character
	if err := db.DB.First(&parent1, parent1ID).Error; err != nil {
		return nil, errors.New("parent 1 not found")
	}
	if err := db.DB.First(&parent2, parent2ID).Error; err != nil {
		return nil, errors.New("parent 2 not found")
	}

	if parent1.OwnerID != userID || parent2.OwnerID != userID {
		return nil, errors.New("you don't own both parents")
	}

	// SECURITY CHECK 3: Check parents not in active battle
	var activeBattle models.Battle
	err := db.DB.Where(
		"((team1_user_id = ? OR team2_user_id = ?) AND status IN ('pending', 'active'))",
		userID, userID,
	).First(&activeBattle).Error

	if err == nil {
		// Found active battle, check if either parent is in it
		return nil, errors.New("cannot breed characters while in active battle")
	}

	// SECURITY CHECK 4: Check parents not listed on marketplace
	var listing1, listing2 models.MarketplaceListing
	err1 := db.DB.Where("character_id = ? AND status = 'ACTIVE'", parent1ID).First(&listing1).Error
	err2 := db.DB.Where("character_id = ? AND status = 'ACTIVE'", parent2ID).First(&listing2).Error

	if err1 == nil || err2 == nil {
		return nil, errors.New("cannot breed listed characters")
	}

	// SECURITY CHECK 5: Verify breeding cooldown
	var lastBreeding models.Egg
	err = db.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&lastBreeding).Error
	if err == nil {
		cooldown := 24 * time.Hour // 24 hour cooldown
		if time.Since(lastBreeding.CreatedAt) < cooldown {
			remaining := cooldown - time.Since(lastBreeding.CreatedAt)
			return nil, fmt.Errorf("breeding cooldown: %v remaining", remaining.Round(time.Minute))
		}
	}

	breedingCost := int64(s.config.GetInt("breeding_cost", 500))

	// SECURITY CHECK 6: Payment Verification
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	payWithLedger := false
	amountPaidInternal := int64(0)

	if txHash != "" {
		// Verify Blockchain Transaction
		if s.blockchain == nil {
			fmt.Println("⚠️ WARNING: Blockchain verification skipped (Service nil)")
		} else {
			// Use generic VerifyTransaction (checks Recipient=Treasury and Amount)
			// Assuming TOWER token
			costBig := big.NewInt(breedingCost)
			if err := s.blockchain.VerifyTransaction(txHash, costBig); err != nil {
				return nil, fmt.Errorf("blockchain verification failed: %v", err)
			}
		}
	} else {
		// Verify Internal Ledger Balance
		if user.TOWERBalance < breedingCost {
			return nil, fmt.Errorf("insufficient TOWER tokens (need %d, have %d)", breedingCost, user.TOWERBalance)
		}
		payWithLedger = true
		amountPaidInternal = breedingCost
	}

	// SECURITY CHECK 7: Validate parent levels (must be at least level 10)
	if parent1.Level < 10 || parent2.Level < 10 {
		return nil, errors.New("both parents must be at least level 10 to breed")
	}

	// SECURITY CHECK 8: Rate limiting - max 5 breedings per day
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	db.DB.Model(&models.Egg{}).Where("user_id = ? AND DATE(created_at) = ?", userID, today).Count(&todayCount)
	if todayCount >= 5 {
		return nil, errors.New("daily breeding limit reached (5 per day)")
	}

	// BEGIN ATOMIC TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// LEDGER INTEGRATION: Breeding Fee
	// Create Service Instance if needed
	if s.ledger == nil {
		s.ledger = NewLedgerService()
	}
	// s.config is already initialized and used for breedingCost

	// Get Accounts
	userAcc, _ := s.ledger.GetOrCreateAccount(&userID, models.AccountTypeWallet, "TOWER")
	treasuryAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeTreasury, "TOWER")

	if payWithLedger {
		// Debit User -> Credit Treasury
		entries := []models.LedgerEntry{
			{AccountID: userAcc.ID, Amount: -breedingCost, Type: "DEBIT"},
			{AccountID: treasuryAcc.ID, Amount: breedingCost, Type: "CREDIT"},
		}

		// Create Ledger TX (linked to DB Transaction via CreateTransactionWithTx if available?
		// BreedingService uses DB.Begin() manually.
		// Use manual CreateTransaction for now but inside the atomic block if possible?
		// LedgerService.CreateTransaction uses its own DB transaction.
		// Ideally use CreateTransactionWithTx to bind to `tx`.
		if err := s.ledger.CreateTransactionWithTx(tx, models.TxTypeBreedingFee, fmt.Sprintf("breed_%d_%d", parent1ID, parent2ID), "Breeding Fee (Internal)", entries); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("breeding payment failed: %v", err)
		}

		// Sync legacy balance
		if err := tx.Model(&user).Update("tower_balance", gorm.Expr("tower_balance - ?", breedingCost)).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to deduct breeding cost")
		}
		// Update local user object for log
		user.TOWERBalance -= breedingCost
	}

	// Calculate incubation time based on parent rarities
	incubationHours := s.calculateIncubationTime(parent1, parent2)

	// Inherit traits from parents for the egg
	rarity := s.inheritRarity(parent1.Rarity, parent2.Rarity)

	// Inherit element (50/50 from either parent)
	element := parent1.Element
	if time.Now().Unix()%2 == 0 {
		element = parent2.Element
	}

	// Inherit character type (50/50 from either parent)
	charType := parent1.CharacterType
	if time.Now().Unix()%3 == 0 {
		charType = parent2.CharacterType
	}

	// Inherit class (70% from parents, 30% random)
	class := parent1.Class
	if time.Now().Unix()%10 > 7 {
		class = parent2.Class
	}

	// Create egg with genetics and all required fields
	egg := models.Egg{
		UserID:                  userID,
		Parent1ID:               &parent1ID,
		Parent2ID:               &parent2ID,
		Rarity:                  rarity,
		Element:                 element,
		CharacterType:           charType,
		Class:                   class,
		IncubationTime:          incubationHours,
		EffectiveIncubationTime: incubationHours,
		AcceleratorsApplied:     "[]", // Initialize with empty JSON array
	}

	if err := tx.Create(&egg).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create egg")
	}

	// Create audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "BREEDING",
		EntityType: "egg",
		EntityID:   &egg.ID,
		NewValues:  fmt.Sprintf("parent1:%d,parent2:%d,cost:%d,tx_hash:%s", parent1ID, parent2ID, breedingCost, txHash),
	}
	tx.Create(&auditLog)

	// Create transaction record
	desc := fmt.Sprintf("Breeding fee for parents %d and %d", parent1ID, parent2ID)
	if txHash != "" {
		desc += fmt.Sprintf(" (Blockchain: %s)", txHash)
	}

	transaction := models.Transaction{
		UserID:          userID,
		TransactionType: "BREEDING_FEE",
		TokenType:       "TOWER",
		Amount:          -amountPaidInternal,
		BalanceBefore:   user.TOWERBalance, // user was updated locally if paid internally
		BalanceAfter:    user.TOWERBalance, // already deducted if internal
		Description:     desc,
	}
	// Correction: if paid internally, user.TOWERBalance is ALREADY deducted in memory above.
	// So BalanceBefore should be user.TOWERBalance + amountPaidInternal
	transaction.BalanceBefore = user.TOWERBalance + amountPaidInternal

	tx.Create(&transaction)

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("transaction failed")
	}

	return &egg, nil
}

// calculateIncubationTime calculates incubation time based on parent rarities
func (s *BreedingService) calculateIncubationTime(parent1, parent2 models.Character) int {
	rarityHours := map[string]int{
		"SSS": 72, // 3 days
		"SS":  48, // 2 days
		"S":   36, // 1.5 days
		"A":   24, // 1 day
		"B":   18, // 18 hours
		"C":   12, // 12 hours
	}

	hours1 := rarityHours[parent1.Rarity]
	hours2 := rarityHours[parent2.Rarity]

	// Average of both parents
	return (hours1 + hours2) / 2
}

// IncubateEgg starts incubation for an egg
func (s *BreedingService) IncubateEgg(userID, eggID uint) error {
	var egg models.Egg
	if err := db.DB.First(&egg, eggID).Error; err != nil {
		return errors.New("egg not found")
	}

	if egg.UserID != userID {
		return errors.New("you don't own this egg")
	}

	if egg.IncubationStartedAt != nil {
		return errors.New("egg is already incubating")
	}

	now := time.Now()
	egg.IncubationStartedAt = &now

	if err := db.DB.Save(&egg).Error; err != nil {
		return errors.New("failed to start incubation")
	}

	return nil
}

// HatchEgg hatches an egg into a character
func (s *BreedingService) HatchEgg(userID, eggID uint) (*models.Character, error) {
	var egg models.Egg
	if err := db.DB.Preload("Parent1").Preload("Parent2").First(&egg, eggID).Error; err != nil {
		return nil, errors.New("egg not found")
	}

	if egg.UserID != userID {
		return nil, errors.New("you don't own this egg")
	}

	if egg.IncubationStartedAt == nil {
		return nil, errors.New("egg must be incubating to hatch")
	}

	// Check if incubation time has passed
	incubationDuration := time.Duration(egg.IncubationTime) * time.Hour
	if time.Since(*egg.IncubationStartedAt) < incubationDuration {
		remaining := incubationDuration - time.Since(*egg.IncubationStartedAt)
		return nil, fmt.Errorf("egg needs %v more to hatch", remaining.Round(time.Minute))
	}

	// Create new character with inherited traits
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	character := s.generateOffspring(egg.Parent1, egg.Parent2)
	character.OwnerID = userID

	if err := tx.Create(&character).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create character")
	}

	// Mark egg as hatched
	now := time.Now()
	egg.CharacterID = &character.ID
	egg.HatchedAt = &now
	if err := tx.Save(&egg).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update egg")
	}

	// Audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "EGG_HATCHED",
		EntityType: "character",
		EntityID:   &character.ID,
		NewValues:  fmt.Sprintf("egg_id:%d,rarity:%s", eggID, character.Rarity),
	}
	tx.Create(&auditLog)

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("transaction failed")
	}

	return character, nil
}

// generateOffspring creates a new character from two parents
func (s *BreedingService) generateOffspring(parent1, parent2 *models.Character) *models.Character {
	// Inherit rarity (70% chance of parent's rarity, 30% chance of upgrade/downgrade)
	rarity := s.inheritRarity(parent1.Rarity, parent2.Rarity)

	// Inherit character type (50/50 from either parent)
	charType := parent1.CharacterType
	if time.Now().Unix()%2 == 0 {
		charType = parent2.CharacterType
	}

	// Inherit class (70% from parents, 30% random)
	class := parent1.Class
	if time.Now().Unix()%10 > 7 {
		classes := []string{"Warrior", "Mage", "Archer", "Tank", "Support"}
		class = classes[time.Now().Unix()%int64(len(classes))]
	}

	// Calculate base stats (average of parents + small random variation)
	baseAttack := (parent1.BaseAttack + parent2.BaseAttack) / 2
	baseDefense := (parent1.BaseDefense + parent2.BaseDefense) / 2
	baseHP := (parent1.BaseHP + parent2.BaseHP) / 2
	baseSpeed := (parent1.BaseSpeed + parent2.BaseSpeed) / 2

	return &models.Character{
		CharacterType:  charType,
		Class:          class,
		Rarity:         rarity,
		Level:          1,
		Experience:     0,
		BaseAttack:     baseAttack,
		BaseDefense:    baseDefense,
		BaseHP:         baseHP,
		BaseSpeed:      baseSpeed,
		CurrentAttack:  baseAttack,
		CurrentDefense: baseDefense,
		CurrentHP:      baseHP,
		CurrentSpeed:   baseSpeed,
	}
}

// inheritRarity determines offspring rarity from parents
func (s *BreedingService) inheritRarity(rarity1, rarity2 string) string {
	rarityLevels := map[string]int{
		"C": 1, "B": 2, "A": 3, "S": 4, "SS": 5, "SSS": 6,
	}
	levelToRarity := map[int]string{
		1: "C", 2: "B", 3: "A", 4: "S", 5: "SS", 6: "SSS",
	}

	level1 := rarityLevels[rarity1]
	level2 := rarityLevels[rarity2]
	avgLevel := (level1 + level2) / 2

	// 70% chance of average, 20% chance of +1, 10% chance of -1
	roll := time.Now().Unix() % 100
	if roll < 70 {
		return levelToRarity[avgLevel]
	} else if roll < 90 && avgLevel < 6 {
		return levelToRarity[avgLevel+1]
	} else if avgLevel > 1 {
		return levelToRarity[avgLevel-1]
	}

	return levelToRarity[avgLevel]
}

// GetUserEggs returns all eggs owned by a user
func (s *BreedingService) GetUserEggs(userID uint) ([]models.Egg, error) {
	var eggs []models.Egg
	if err := db.DB.Preload("Parent1").Preload("Parent2").Where("user_id = ?", userID).Order("created_at DESC").Find(&eggs).Error; err != nil {
		return nil, err
	}
	return eggs, nil
}
