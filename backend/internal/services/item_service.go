package services

import (
	"errors"
	"fmt"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// ItemService handles item-related business logic
type ItemService struct{}

// NewItemService creates a new item service
func NewItemService() *ItemService {
	return &ItemService{}
}

// CreateItem creates a new item
func (s *ItemService) CreateItem(ownerID uint, itemType, name, rarity string, isConsumable, isCraftingMaterial, isStackable bool) (*models.Item, error) {
	// Calculate stat bonuses based on rarity and type
	bonuses := s.calculateItemBonuses(itemType, rarity)

	item := &models.Item{
		OwnerID:            ownerID,
		ItemType:           itemType,
		Name:               name,
		Rarity:             rarity,
		AttackBonus:        bonuses.Attack,
		DefenseBonus:       bonuses.Defense,
		HPBonus:            bonuses.HP,
		SpeedBonus:         bonuses.Speed,
		Durability:         100,
		IsConsumable:       isConsumable,
		IsCraftingMaterial: isCraftingMaterial,
		IsStackable:        isStackable,
		Quantity:           1,
	}

	if err := db.DB.Create(item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

// GetUserItems retrieves all items for a user
func (s *ItemService) GetUserItems(userID uint, itemType string) ([]models.Item, error) {
	var items []models.Item
	query := db.DB.Where("owner_id = ?", userID)

	if itemType != "" {
		query = query.Where("item_type = ?", itemType)
	}

	if err := query.Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// GetItemByID retrieves an item by ID
func (s *ItemService) GetItemByID(itemID, ownerID uint) (*models.Item, error) {
	var item models.Item
	if err := db.DB.Where("id = ? AND owner_id = ?", itemID, ownerID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// EquipItem equips an item to a character
func (s *ItemService) EquipItem(itemID, characterID, ownerID uint) error {
	// Get item
	var item models.Item
	if err := db.DB.Where("id = ? AND owner_id = ?", itemID, ownerID).First(&item).Error; err != nil {
		return errors.New("item not found")
	}

	// Get character
	var character models.Character
	if err := db.DB.Where("id = ? AND owner_id = ?", characterID, ownerID).First(&character).Error; err != nil {
		return errors.New("character not found")
	}

	// Check if item is equipment
	if item.ItemType != "WEAPON" && item.ItemType != "ARMOR" && item.ItemType != "ACCESSORY" && item.ItemType != "RUNE" {
		return errors.New("item is not equipment")
	}

	// Check if already equipped
	if item.IsEquipped {
		return errors.New("item already equipped")
	}

	// Update item
	item.IsEquipped = true
	item.EquippedByID = &characterID

	// Apply item effects
	character.CurrentAttack += item.AttackBonus
	character.CurrentDefense += item.DefenseBonus
	character.CurrentHP += item.HPBonus
	character.CurrentSpeed += item.SpeedBonus

	// Save changes
	if err := db.DB.Save(character).Error; err != nil {
		return fmt.Errorf("failed to update character stats: %w", err)
	}

	// Also save the item's equipped status
	if err := db.DB.Save(item).Error; err != nil {
		return fmt.Errorf("failed to update item equipped status: %w", err)
	}

	return nil
}

// UnequipItem removes an item from a character
func (s *ItemService) UnequipItem(userID, itemID uint) error {
	// Find the item
	var item models.Item
	if err := db.DB.First(&item, itemID).Error; err != nil {
		return errors.New("item not found")
	}

	if item.OwnerID != userID {
		return errors.New("unauthorized")
	}

	if !item.IsEquipped || item.EquippedByID == nil {
		return errors.New("item is not equipped")
	}

	characterID := *item.EquippedByID

	var character models.Character
	if err := db.DB.Preload("EquippedItems").First(&character, characterID).Error; err != nil {
		return errors.New("character not found")
	}

	// Remove item (simplified association removal)
	item.IsEquipped = false
	item.EquippedByID = nil
	// Update character stats
	character.CurrentAttack -= item.AttackBonus
	character.CurrentDefense -= item.DefenseBonus
	character.CurrentHP -= item.HPBonus
	character.CurrentSpeed -= item.SpeedBonus

	// Ensure stats don't go below base
	if character.CurrentAttack < character.BaseAttack {
		character.CurrentAttack = character.BaseAttack
	}
	if character.CurrentDefense < character.BaseDefense {
		character.CurrentDefense = character.BaseDefense
	}
	if character.CurrentHP < character.BaseHP {
		character.CurrentHP = character.BaseHP
	}
	if character.CurrentSpeed < character.BaseSpeed {
		character.CurrentSpeed = character.BaseSpeed
	}

	// Update item
	item.IsEquipped = false
	item.EquippedByID = nil

	// Save both in transaction
	tx := db.DB.Begin()
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(&character).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// UseConsumable uses a consumable item
func (s *ItemService) UseConsumable(itemID, targetID, ownerID uint) error {
	// Get item
	var item models.Item
	if err := db.DB.Where("id = ? AND owner_id = ?", itemID, ownerID).First(&item).Error; err != nil {
		return errors.New("item not found")
	}

	if !item.IsConsumable {
		return errors.New("item is not consumable")
	}

	if item.Quantity <= 0 {
		return errors.New("no items remaining")
	}

	// Get target character
	var character models.Character
	if err := db.DB.Where("id = ? AND owner_id = ?", targetID, ownerID).First(&character).Error; err != nil {
		return errors.New("character not found")
	}

	// Apply effect based on consume effect type
	switch item.ConsumeEffect {
	case "REVIVE":
		if !character.IsDead {
			return errors.New("character is not dead")
		}
		character.IsDead = false
		character.Durability = 50 // Revive to 50% durability

	case "REDUCE_FATIGUE":
		character.Fatigue -= 50
		if character.Fatigue < 0 {
			character.Fatigue = 0
		}

	case "XP_BOOST":
		character.Experience += 1000

	case "REPAIR":
		character.Durability = 100
		character.IsDead = false

	default:
		return errors.New("unknown consumable effect")
	}

	// Decrease item quantity
	item.Quantity--

	// Save both in transaction
	tx := db.DB.Begin()
	if err := tx.Save(&character).Error; err != nil {
		tx.Rollback()
		return err
	}

	if item.Quantity == 0 {
		// Delete item if quantity is 0
		if err := tx.Delete(&item).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// UpdateDurability updates item durability after battle
func (s *ItemService) UpdateDurability(itemID uint, durabilityLoss int) error {
	var item models.Item
	if err := db.DB.First(&item, itemID).Error; err != nil {
		return err
	}

	item.Durability -= durabilityLoss
	if item.Durability <= 0 {
		item.Durability = 0
		item.IsBroken = true
		// Auto-unequip broken items
		item.IsEquipped = false
		item.EquippedByID = nil
	}

	return db.DB.Save(&item).Error
}

// calculateItemBonuses returns stat bonuses based on type and rarity
func (s *ItemService) calculateItemBonuses(itemType, rarity string) struct{ Attack, Defense, HP, Speed int } {
	bonuses := struct{ Attack, Defense, HP, Speed int }{}

	// Rarity multiplier
	multiplier := 1.0
	switch rarity {
	case "SSS":
		multiplier = 3.0
	case "SS":
		multiplier = 2.5
	case "S":
		multiplier = 2.0
	case "A":
		multiplier = 1.5
	case "B":
		multiplier = 1.2
	case "C":
		multiplier = 1.0
	}

	// Base bonuses by type
	switch itemType {
	case "WEAPON":
		bonuses.Attack = int(float64(50) * multiplier)
		bonuses.Speed = int(float64(10) * multiplier)
	case "ARMOR":
		bonuses.Defense = int(float64(50) * multiplier)
		bonuses.HP = int(float64(100) * multiplier)
	case "ACCESSORY":
		bonuses.Attack = int(float64(20) * multiplier)
		bonuses.Defense = int(float64(20) * multiplier)
		bonuses.HP = int(float64(50) * multiplier)
		bonuses.Speed = int(float64(20) * multiplier)
	case "RUNE":
		bonuses.Attack = int(float64(30) * multiplier)
		bonuses.Defense = int(float64(30) * multiplier)
	}

	return bonuses
}
