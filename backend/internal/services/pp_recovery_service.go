package services

import (
	"errors"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// PPRecoveryService handles PP (Power Points) recovery for character abilities
type PPRecoveryService struct{}

// NewPPRecoveryService creates a new PP recovery service
func NewPPRecoveryService() *PPRecoveryService {
	return &PPRecoveryService{}
}

// UseAbility deducts PP when an ability is used
func (s *PPRecoveryService) UseAbility(characterID, abilityID uint) error {
	var usage models.AbilityUsage
	err := db.DB.Where("character_id = ? AND ability_id = ?", characterID, abilityID).First(&usage).Error

	if err == gorm.ErrRecordNotFound {
		// First time using this ability, create usage record
		var ability models.Ability
		if err := db.DB.First(&ability, abilityID).Error; err != nil {
			return errors.New("ability not found")
		}

		usage = models.AbilityUsage{
			CharacterID: characterID,
			AbilityID:   abilityID,
			CurrentPP:   ability.MaxPP - 1, // Use 1 PP
			MaxPP:       ability.MaxPP,
		}

		now := time.Now()
		usage.LastUsedAt = &now
		usage.LastRecoveredAt = now

		if err := db.DB.Create(&usage).Error; err != nil {
			return errors.New("failed to create ability usage")
		}

		return nil
	}

	// Check if PP available
	if usage.CurrentPP <= 0 {
		return errors.New("no PP remaining for this ability")
	}

	// Deduct 1 PP
	now := time.Now()
	usage.CurrentPP--
	usage.LastUsedAt = &now

	if err := db.DB.Save(&usage).Error; err != nil {
		return errors.New("failed to update PP")
	}

	return nil
}

// RecoverPPNaturally recovers PP naturally over time (1 PP per 30 minutes)
func (s *PPRecoveryService) RecoverPPNaturally(characterID uint) error {
	var usages []models.AbilityUsage
	db.DB.Where("character_id = ?", characterID).Find(&usages)

	now := time.Now()

	for _, usage := range usages {
		// Skip if already at max PP
		if usage.CurrentPP >= usage.MaxPP {
			continue
		}

		// Calculate time since last recovery
		timeSinceRecovery := now.Sub(usage.LastRecoveredAt)

		// 1 PP per 30 minutes
		ppToRecover := int(timeSinceRecovery.Minutes() / 30)

		if ppToRecover > 0 {
			// Recover PP
			usage.CurrentPP += ppToRecover
			if usage.CurrentPP > usage.MaxPP {
				usage.CurrentPP = usage.MaxPP
			}

			usage.LastRecoveredAt = now

			db.DB.Save(&usage)
		}
	}

	return nil
}

// UsePPPotion uses a PP recovery item
func (s *PPRecoveryService) UsePPPotion(userID, characterID, abilityID, itemID uint) error {
	// Get item
	var item models.ShopItem
	if err := db.DB.First(&item, itemID).Error; err != nil {
		return errors.New("item not found")
	}

	// Verify item is PP recovery type
	if item.Category != "pp" {
		return errors.New("item is not a PP recovery item")
	}

	// Get user inventory
	var inventory models.UserInventory
	err := db.DB.Where("user_id = ? AND item_id = ?", userID, itemID).First(&inventory).Error
	if err != nil {
		return errors.New("you don't have this item")
	}

	if inventory.Quantity < 1 {
		return errors.New("no items remaining")
	}

	// Get ability usage
	var usage models.AbilityUsage
	err = db.DB.Where("character_id = ? AND ability_id = ?", characterID, abilityID).First(&usage).Error
	if err == gorm.ErrRecordNotFound {
		return errors.New("ability not yet used")
	}

	// Check if already at max PP
	if usage.CurrentPP >= usage.MaxPP {
		return errors.New("ability already at max PP")
	}

	// BEGIN TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Apply PP recovery based on item type
	switch item.EffectType {
	case "restore_pp":
		// PP Potion: Restore 5 PP to one ability
		usage.CurrentPP += item.EffectValue
		if usage.CurrentPP > usage.MaxPP {
			usage.CurrentPP = usage.MaxPP
		}

	case "restore_pp_full":
		// PP Max: Fully restore PP to one ability
		usage.CurrentPP = usage.MaxPP

	case "restore_pp_all":
		// Elixir: Restore 5 PP to all abilities
		return s.usePPPotionAll(userID, characterID, itemID, item.EffectValue, false, tx)

	case "restore_pp_all_full":
		// Max Elixir: Fully restore PP to all abilities
		return s.usePPPotionAll(userID, characterID, itemID, 0, true, tx)
	}

	// Save usage
	if err := tx.Save(&usage).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to restore PP")
	}

	// Consume item
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

	return nil
}

// usePPPotionAll restores PP to all abilities
func (s *PPRecoveryService) usePPPotionAll(userID, characterID, itemID uint, amount int, fullRestore bool, tx *gorm.DB) error {
	// Get all ability usages for character
	var usages []models.AbilityUsage
	tx.Where("character_id = ?", characterID).Find(&usages)

	if len(usages) == 0 {
		tx.Rollback()
		return errors.New("no abilities used yet")
	}

	// Restore PP to all
	for _, usage := range usages {
		if fullRestore {
			usage.CurrentPP = usage.MaxPP
		} else {
			usage.CurrentPP += amount
			if usage.CurrentPP > usage.MaxPP {
				usage.CurrentPP = usage.MaxPP
			}
		}
		tx.Save(&usage)
	}

	// Consume item
	var inventory models.UserInventory
	err := tx.Where("user_id = ? AND item_id = ?", userID, itemID).First(&inventory).Error
	if err != nil {
		tx.Rollback()
		return errors.New("you don't have this item")
	}

	if inventory.Quantity < 1 {
		tx.Rollback()
		return errors.New("no items remaining")
	}

	if err := tx.Model(&inventory).Update("quantity", gorm.Expr("quantity - 1")).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to consume item")
	}

	if inventory.Quantity == 1 {
		tx.Delete(&inventory)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction failed")
	}

	return nil
}

// GetPPStatus returns PP status for all abilities of a character
func (s *PPRecoveryService) GetPPStatus(characterID uint) ([]models.AbilityUsage, error) {
	var usages []models.AbilityUsage
	if err := db.DB.Where("character_id = ?", characterID).Preload("Ability").Find(&usages).Error; err != nil {
		return nil, err
	}

	return usages, nil
}

// StartRest initiates natural PP recovery for a character
func (s *PPRecoveryService) StartRest(characterID uint) error {
	// This is called when a character starts resting
	// Natural recovery happens automatically via background job
	// For now, just trigger immediate recovery
	return s.RecoverPPNaturally(characterID)
}

// RecoverAllCharactersPP background job to recover PP for all characters
func (s *PPRecoveryService) RecoverAllCharactersPP() error {
	// Get all ability usages that need recovery
	var usages []models.AbilityUsage
	db.DB.Where("current_pp < max_pp").Find(&usages)

	now := time.Now()
	recovered := 0

	for _, usage := range usages {
		// Calculate time since last recovery
		timeSinceRecovery := now.Sub(usage.LastRecoveredAt)

		// 1 PP per 30 minutes
		ppToRecover := int(timeSinceRecovery.Minutes() / 30)

		if ppToRecover > 0 {
			usage.CurrentPP += ppToRecover
			if usage.CurrentPP > usage.MaxPP {
				usage.CurrentPP = usage.MaxPP
			}

			usage.LastRecoveredAt = now
			db.DB.Save(&usage)
			recovered++
		}
	}

	return nil
}
