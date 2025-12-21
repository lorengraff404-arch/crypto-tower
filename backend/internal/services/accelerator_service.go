package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// AcceleratorService handles egg accelerator items
type AcceleratorService struct{}

// NewAcceleratorService creates a new accelerator service
func NewAcceleratorService() *AcceleratorService {
	return &AcceleratorService{}
}

// ApplyAccelerator applies an accelerator item to an egg
func (s *AcceleratorService) ApplyAccelerator(userID, eggID, itemID uint) error {
	// Get egg
	var egg models.Egg
	if err := db.DB.First(&egg, eggID).Error; err != nil {
		return errors.New("egg not found")
	}

	// SECURITY CHECK 1: Verify ownership
	if egg.UserID != userID {
		return errors.New("you don't own this egg")
	}

	// SECURITY CHECK 2: Check egg not already hatched
	if egg.HatchedAt != nil {
		return errors.New("egg already hatched")
	}

	// Get item
	var item models.ShopItem
	if err := db.DB.First(&item, itemID).Error; err != nil {
		return errors.New("item not found")
	}

	// SECURITY CHECK 3: Verify item is accelerator
	if item.Category != "egg" {
		return errors.New("item is not an egg accelerator")
	}

	// Get user inventory
	var inventory models.UserInventory
	err := db.DB.Where("user_id = ? AND item_id = ?", userID, itemID).First(&inventory).Error
	if err != nil {
		return errors.New("you don't have this item")
	}

	// SECURITY CHECK 4: Check quantity
	if inventory.Quantity < 1 {
		return errors.New("no items remaining")
	}

	// Parse currently applied accelerators
	var appliedAccelerators []uint
	if egg.AcceleratorsApplied != "" && egg.AcceleratorsApplied != "[]" {
		json.Unmarshal([]byte(egg.AcceleratorsApplied), &appliedAccelerators)
	}

	// SECURITY CHECK 5: Check if item already applied
	for _, appliedID := range appliedAccelerators {
		if appliedID == itemID {
			return errors.New("this accelerator is already applied to this egg")
		}
	}

	// SECURITY CHECK 6: Check synergy/conflicts
	if err := s.validateSynergy(appliedAccelerators, itemID); err != nil {
		return err
	}

	// BEGIN TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Handle instant hatch
	if item.EffectType == "instant_hatch" {
		// Instant hatch - set incubation to 0
		egg.EffectiveIncubationTime = 0

		// If incubation already started, set it to ready to hatch
		if egg.IncubationStartedAt != nil {
			pastTime := time.Now().Add(-time.Duration(egg.IncubationTime+1) * time.Hour)
			egg.IncubationStartedAt = &pastTime
		}
	} else {
		// Calculate time reduction
		reduction := s.calculateTimeReduction(appliedAccelerators, itemID)
		newTime := int(float64(egg.IncubationTime) * (1.0 - reduction))

		// Ensure minimum 1 hour
		if newTime < 1 {
			newTime = 1
		}

		egg.EffectiveIncubationTime = newTime
	}

	// Add to applied accelerators
	appliedAccelerators = append(appliedAccelerators, itemID)
	acceleratorsJSON, _ := json.Marshal(appliedAccelerators)
	egg.AcceleratorsApplied = string(acceleratorsJSON)

	// Save egg
	if err := tx.Save(&egg).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update egg")
	}

	// Consume item
	if err := tx.Model(&inventory).Update("quantity", gorm.Expr("quantity - 1")).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to consume item")
	}

	// Delete inventory entry if quantity = 0
	if inventory.Quantity == 1 {
		tx.Delete(&inventory)
	}

	// Create audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "ACCELERATOR_APPLIED",
		EntityType: "egg",
		EntityID:   &eggID,
		NewValues:  fmt.Sprintf("item:%s,new_time:%d", item.Name, egg.EffectiveIncubationTime),
	}
	tx.Create(&auditLog)

	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction failed")
	}

	return nil
}

// validateSynergy checks if new item can be applied with existing accelerators
func (s *AcceleratorService) validateSynergy(appliedIDs []uint, newItemID uint) error {
	// Get all applied items
	var appliedItems []models.ShopItem
	if len(appliedIDs) > 0 {
		db.DB.Where("id IN ?", appliedIDs).Find(&appliedItems)
	}

	// Get new item
	var newItem models.ShopItem
	if err := db.DB.First(&newItem, newItemID).Error; err != nil {
		return errors.New("item not found")
	}

	// Instant hatch cannot be combined with anything
	if newItem.EffectType == "instant_hatch" {
		if len(appliedItems) > 0 {
			return errors.New("instant hatch cannot be combined with other accelerators")
		}
	}

	// Check if instant hatch already applied
	for _, item := range appliedItems {
		if item.EffectType == "instant_hatch" {
			return errors.New("instant hatch already applied - cannot add more accelerators")
		}
	}

	// Synergy rules based on item names
	itemNames := make(map[string]bool)
	for _, item := range appliedItems {
		itemNames[item.Name] = true
	}

	// Nest and Incubator conflict
	if newItem.Name == "Basic Nest" && itemNames["Advanced Incubator"] {
		return errors.New("nest cannot be combined with incubator")
	}
	if newItem.Name == "Advanced Incubator" && itemNames["Basic Nest"] {
		return errors.New("incubator cannot be combined with nest")
	}

	// Incubator and Solar Lamp conflict
	if newItem.Name == "Advanced Incubator" && itemNames["Solar Heat Lamp"] {
		return errors.New("incubator cannot be combined with solar lamp")
	}
	if newItem.Name == "Solar Heat Lamp" && itemNames["Advanced Incubator"] {
		return errors.New("solar lamp cannot be combined with incubator")
	}

	// Nest + Solar Lamp is OK (synergy)
	// No conflict check needed

	return nil
}

// calculateTimeReduction calculates total time reduction from accelerators
func (s *AcceleratorService) calculateTimeReduction(appliedIDs []uint, newItemID uint) float64 {
	// Get all items including new one
	allIDs := append(appliedIDs, newItemID)

	var items []models.ShopItem
	db.DB.Where("id IN ?", allIDs).Find(&items)

	totalReduction := 0.0

	// Check for synergy combinations
	hasNest := false
	hasSolarLamp := false

	for _, item := range items {
		if item.Name == "Basic Nest" {
			hasNest = true
		}
		if item.Name == "Solar Heat Lamp" {
			hasSolarLamp = true
		}
	}

	// Nest + Solar Lamp synergy = 87.5% reduction (not 100%)
	if hasNest && hasSolarLamp {
		// 25% + 75% = 100%, but we cap at 87.5% for balance
		totalReduction = 0.875
	} else {
		// No synergy, just add reductions
		for _, item := range items {
			reduction := float64(item.EffectValue) / 100.0
			totalReduction += reduction
		}

		// Cap at 90% reduction
		if totalReduction > 0.90 {
			totalReduction = 0.90
		}
	}

	return totalReduction
}

// GetAcceleratorInfo returns info about an accelerator item
func (s *AcceleratorService) GetAcceleratorInfo(itemID uint) (*models.ShopItem, error) {
	var item models.ShopItem
	if err := db.DB.First(&item, itemID).Error; err != nil {
		return nil, errors.New("item not found")
	}

	if item.Category != "egg" {
		return nil, errors.New("item is not an egg accelerator")
	}

	return &item, nil
}

// GetAppliedAccelerators returns list of accelerators applied to an egg
func (s *AcceleratorService) GetAppliedAccelerators(eggID uint) ([]models.ShopItem, error) {
	var egg models.Egg
	if err := db.DB.First(&egg, eggID).Error; err != nil {
		return nil, errors.New("egg not found")
	}

	if egg.AcceleratorsApplied == "" || egg.AcceleratorsApplied == "[]" {
		return []models.ShopItem{}, nil
	}

	var appliedIDs []uint
	json.Unmarshal([]byte(egg.AcceleratorsApplied), &appliedIDs)

	var items []models.ShopItem
	if len(appliedIDs) > 0 {
		db.DB.Where("id IN ?", appliedIDs).Find(&items)
	}

	return items, nil
}
