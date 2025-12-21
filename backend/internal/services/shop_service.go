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

// ShopService handles shop operations
type ShopService struct {
	ledger     *LedgerService
	blockchain *BlockchainService
}

// NewShopService creates a new shop service
func NewShopService(bc *BlockchainService) *ShopService {
	return &ShopService{
		ledger:     NewLedgerService(),
		blockchain: bc,
	}
}

// GetShopItems returns all available shop items
func (s *ShopService) GetShopItems(category string) ([]models.ShopItem, error) {
	var items []models.ShopItem
	query := db.DB.Where("is_available = true")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("category, gtk_cost").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// BuyItem purchases an item from the shop
func (s *ShopService) BuyItem(userID, itemID uint, quantity int, txHash string) error {
	// SECURITY CHECK 1: Validate quantity
	if quantity < 1 || quantity > 99 {
		return errors.New("quantity must be between 1 and 99")
	}

	// Get item
	var item models.ShopItem
	if err := db.DB.First(&item, itemID).Error; err != nil {
		return errors.New("item not found")
	}

	// SECURITY CHECK 2: Check availability
	if !item.IsAvailable {
		return errors.New("item not available")
	}

	// Calculate total cost
	totalCost := item.GTKCost * int64(quantity)

	// SECURITY CHECK 3: Verify Blockchain Transaction (TxHash)
	// Note: We use big.Int for blockchain calls
	if s.blockchain != nil {
		// Convert totalCost to Token Units (assuming 18 decimals or consistent unit)
		// Usually DB stores integer representation. Assuming DB stores raw units (Wei).
		costBig := big.NewInt(totalCost)
		// We can optionally verify the sender using "from" if we extracted it, but validating receipt of funds is primary
		if err := s.blockchain.VerifyTransaction(txHash, costBig); err != nil {
			return fmt.Errorf("blockchain verification failed: %v", err)
		}
	} else {
		// Log warning if blockchain service not active (DEV MODE)
		fmt.Println("⚠️ WARNING: Blockchain verification skipped (Service nil)")
	}

	// Get user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// NOTE: Balance verification happens on blockchain via tx_hash
	// Frontend transfers GTK before calling this endpoint
	// Backend will verify the transaction in the handler

	// LEDGER INTEGRATION: Execute Purchase Transaction
	// Debit: User Wallet, Credit: System Sink (or Treasury)
	// We use "sink" for shop purchases usually as it removes tokens from circulation

	userAcc, _ := s.ledger.GetOrCreateAccount(&userID, models.AccountTypeWallet, "GTK")
	sinkAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeSink, "GTK")

	ledgerEntries := []models.LedgerEntry{
		{AccountID: userAcc.ID, Amount: -totalCost, Type: "DEBIT"},
		{AccountID: sinkAcc.ID, Amount: totalCost, Type: "CREDIT"},
	}

	// Create Ledger Transaction
	if err := s.ledger.CreateTransaction(models.TxTypeShopBuy, fmt.Sprintf("shop_buy_%d_%d", userID, time.Now().Unix()), fmt.Sprintf("Bought %dx %s", quantity, item.Name), ledgerEntries); err != nil {
		return fmt.Errorf("purchase failed: %v", err)
	}

	// Inventory Management (Still using DB tx for atomic inventory update, but relying on Ledger for funds)
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Add to inventory (or increase quantity)
	var inventory models.UserInventory
	err := tx.Where("user_id = ? AND item_id = ?", userID, itemID).First(&inventory).Error

	if err == gorm.ErrRecordNotFound {
		// Create new inventory entry
		inventory = models.UserInventory{
			UserID:   userID,
			ItemID:   itemID,
			Quantity: quantity,
		}
		if err := tx.Create(&inventory).Error; err != nil {
			tx.Rollback()
			return errors.New("failed to add to inventory")
		}
	} else {
		// Update existing
		newQuantity := inventory.Quantity + quantity
		if newQuantity > item.MaxStack {
			tx.Rollback()
			return fmt.Errorf("max stack limit (%d) exceeded", item.MaxStack)
		}
		if err := tx.Model(&inventory).Update("quantity", newQuantity).Error; err != nil {
			tx.Rollback()
			return errors.New("failed to update inventory")
		}
	}

	// Update legacy balance just for UI sync (optional, or remove if fully relying on Ledger view)
	// For now, we sync it so legacy endpoints still work
	if err := tx.Model(&user).Update("gtk_balance", gorm.Expr("gtk_balance - ?", totalCost)).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to sync legacy balance")
	}

	// Create transaction record (Legacy)
	transaction := models.Transaction{
		UserID:          userID,
		TransactionType: "SHOP_PURCHASE",
		TokenType:       "GTK",
		Amount:          -totalCost,
		BalanceBefore:   user.GTKBalance, // This might be slightly stale if concurrent, but acceptable for legacy log
		BalanceAfter:    user.GTKBalance - totalCost,
		Description:     fmt.Sprintf("Bought %dx %s", quantity, item.Name),
	}
	tx.Create(&transaction)

	// Audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "SHOP_PURCHASE",
		EntityType: "shop_item",
		EntityID:   &itemID,
		NewValues:  fmt.Sprintf("quantity:%d,cost:%d", quantity, totalCost),
	}
	tx.Create(&auditLog)

	if err := tx.Commit().Error; err != nil {
		// If DB fails, we have a problem because Ledger succeeded.
		// Ideally wrap Ledger + DB in distributed tx or saga.
		// For this implementation, we assume DB commit reliability is high after Ledger check.
		// A rigorous solution would use a pending state.
		return errors.New("transaction commit failed")
	}

	return nil
}

// GetUserInventory returns user's shop inventory
func (s *ShopService) GetUserInventory(userID uint) ([]models.UserInventory, error) {
	var inventory []models.UserInventory
	if err := db.DB.Where("user_id = ?", userID).Preload("Item").Find(&inventory).Error; err != nil {
		return nil, err
	}

	return inventory, nil
}

// UseItem uses a consumable item
func (s *ShopService) UseItem(userID, itemID, characterID uint) error {
	// Get item
	var item models.ShopItem
	if err := db.DB.First(&item, itemID).Error; err != nil {
		return errors.New("item not found")
	}

	// Get inventory
	var inventory models.UserInventory
	err := db.DB.Where("user_id = ? AND item_id = ?", userID, itemID).First(&inventory).Error
	if err != nil {
		return errors.New("item not in inventory")
	}

	if inventory.Quantity < 1 {
		return errors.New("no items remaining")
	}

	// Get character
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return errors.New("character not found")
	}

	// SECURITY CHECK: Verify ownership
	if character.OwnerID != userID {
		return errors.New("you don't own this character")
	}

	// Apply item effect
	if err := s.applyItemEffect(&item, &character); err != nil {
		return err
	}

	// Consume item if consumable
	if item.IsConsumable {
		tx := db.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if err := tx.Model(&inventory).Update("quantity", gorm.Expr("quantity - 1")).Error; err != nil {
			tx.Rollback()
			return errors.New("failed to consume item")
		}

		// Delete if quantity = 0
		if inventory.Quantity == 1 {
			tx.Delete(&inventory)
		}

		if err := tx.Commit().Error; err != nil {
			return errors.New("transaction failed")
		}
	}

	return nil
}

// applyItemEffect applies the item's effect to character
func (s *ShopService) applyItemEffect(item *models.ShopItem, character *models.Character) error {
	switch item.EffectType {
	case "heal_hp":
		// Restore HP
		newHP := character.CurrentHP + item.EffectValue
		if newHP > character.BaseHP {
			newHP = character.BaseHP
		}
		character.CurrentHP = newHP
		db.DB.Save(character)

	case "heal_full_hp":
		// Full heal
		character.CurrentHP = character.BaseHP
		db.DB.Save(character)

	case "cure_poison":
		// Remove poison status
		db.DB.Where("character_id = ? AND status_type = ?", character.ID, "poison").
			Delete(&models.CharacterStatusEffect{})

	case "cure_burn":
		// Remove burn status
		db.DB.Where("character_id = ? AND status_type = ?", character.ID, "burn").
			Delete(&models.CharacterStatusEffect{})

	case "cure_freeze":
		// Remove freeze status
		db.DB.Where("character_id = ? AND status_type = ?", character.ID, "freeze").
			Delete(&models.CharacterStatusEffect{})

	case "cure_paralysis":
		// Remove paralysis status
		db.DB.Where("character_id = ? AND status_type = ?", character.ID, "paralysis").
			Delete(&models.CharacterStatusEffect{})

	case "cure_sleep":
		// Remove sleep status
		db.DB.Where("character_id = ? AND status_type = ?", character.ID, "sleep").
			Delete(&models.CharacterStatusEffect{})

	case "cure_all_status":
		// Remove all status effects
		db.DB.Where("character_id = ?", character.ID).
			Delete(&models.CharacterStatusEffect{})

	case "revive":
		// Revive fainted character
		if !character.IsFainted {
			return errors.New("character is not fainted")
		}
		character.IsFainted = false
		healAmount := int(float64(character.BaseHP) * float64(item.EffectValue) / 100.0)
		character.CurrentHP = healAmount
		db.DB.Save(character)

	case "buff_attack":
		// Apply attack buff
		buff := models.CharacterBuff{
			CharacterID:    character.ID,
			BuffType:       "x_attack",
			Multiplier:     float64(item.EffectValue) / 100.0,
			TurnsRemaining: 3, // Default 3 turns
		}
		db.DB.Create(&buff)

	case "buff_defense":
		// Apply defense buff
		buff := models.CharacterBuff{
			CharacterID:    character.ID,
			BuffType:       "x_defense",
			Multiplier:     float64(item.EffectValue) / 100.0,
			TurnsRemaining: 3,
		}
		db.DB.Create(&buff)

	case "buff_speed":
		// Apply speed buff
		buff := models.CharacterBuff{
			CharacterID:    character.ID,
			BuffType:       "x_speed",
			Multiplier:     float64(item.EffectValue) / 100.0,
			TurnsRemaining: 3,
		}
		db.DB.Create(&buff)

	case "buff_guard":
		// Apply guard buff (prevent stat reduction)
		buff := models.CharacterBuff{
			CharacterID:    character.ID,
			BuffType:       "guard_spec",
			Multiplier:     1.0,
			TurnsRemaining: 5,
		}
		db.DB.Create(&buff)

	case "buff_critical":
		// Apply critical buff
		buff := models.CharacterBuff{
			CharacterID:    character.ID,
			BuffType:       "dire_hit",
			Multiplier:     float64(item.EffectValue) / 100.0,
			TurnsRemaining: 3,
		}
		db.DB.Create(&buff)

	default:
		return fmt.Errorf("unknown effect type: %s", item.EffectType)
	}

	return nil
}
