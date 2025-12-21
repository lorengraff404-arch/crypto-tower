package services

import (
	"errors"
	"fmt"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// AbilityEquipService handles ability equipping logic
type AbilityEquipService struct{}

// NewAbilityEquipService creates a new ability equip service
func NewAbilityEquipService() *AbilityEquipService {
	return &AbilityEquipService{}
}

// EquipAbility equips an ability to a specific slot (1-4)
func (s *AbilityEquipService) EquipAbility(characterID, abilityID uint, slot int) error {
	// Validate slot number
	if slot < 1 || slot > 4 {
		return errors.New("invalid slot: must be 1-4")
	}

	// Get character
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	// Verify character has learned this ability
	var learned models.CharacterLearnedAbility
	if err := db.DB.Where("character_id = ? AND ability_id = ?", characterID, abilityID).
		First(&learned).Error; err != nil {
		return errors.New("ability not learned yet")
	}

	// Equip to slot
	switch slot {
	case 1:
		character.EquippedAbility1 = &abilityID
	case 2:
		character.EquippedAbility2 = &abilityID
	case 3:
		character.EquippedAbility3 = &abilityID
	case 4:
		character.EquippedAbility4 = &abilityID
	}

	// Save
	if err := db.DB.Save(&character).Error; err != nil {
		return fmt.Errorf("failed to equip ability: %w", err)
	}

	return nil
}

// UnequipAbility removes an ability from a slot
func (s *AbilityEquipService) UnequipAbility(characterID uint, slot int) error {
	if slot < 1 || slot > 4 {
		return errors.New("invalid slot: must be 1-4")
	}

	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	// Unequip from slot
	switch slot {
	case 1:
		character.EquippedAbility1 = nil
	case 2:
		character.EquippedAbility2 = nil
	case 3:
		character.EquippedAbility3 = nil
	case 4:
		character.EquippedAbility4 = nil
	}

	if err := db.DB.Save(&character).Error; err != nil {
		return fmt.Errorf("failed to unequip ability: %w", err)
	}

	return nil
}

// GetLearnedAbilities returns ALL abilities a character has learned
func (s *AbilityEquipService) GetLearnedAbilities(characterID uint) ([]models.Ability, error) {
	var learned []models.CharacterLearnedAbility
	if err := db.DB.Where("character_id = ?", characterID).
		Preload("Ability").
		Find(&learned).Error; err != nil {
		return nil, err
	}

	abilities := make([]models.Ability, len(learned))
	for i, l := range learned {
		abilities[i] = l.Ability
	}

	return abilities, nil
}

// GetEquippedAbilities returns the 4 equipped abilities (some may be nil)
func (s *AbilityEquipService) GetEquippedAbilities(characterID uint) ([4]*models.Ability, error) {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return [4]*models.Ability{}, err
	}

	var equipped [4]*models.Ability

	// Load each equipped ability
	if character.EquippedAbility1 != nil {
		var ability models.Ability
		if err := db.DB.First(&ability, *character.EquippedAbility1).Error; err == nil {
			equipped[0] = &ability
		}
	}
	if character.EquippedAbility2 != nil {
		var ability models.Ability
		if err := db.DB.First(&ability, *character.EquippedAbility2).Error; err == nil {
			equipped[1] = &ability
		}
	}
	if character.EquippedAbility3 != nil {
		var ability models.Ability
		if err := db.DB.First(&ability, *character.EquippedAbility3).Error; err == nil {
			equipped[2] = &ability
		}
	}
	if character.EquippedAbility4 != nil {
		var ability models.Ability
		if err := db.DB.First(&ability, *character.EquippedAbility4).Error; err == nil {
			equipped[3] = &ability
		}
	}

	return equipped, nil
}

// LearnAbilityOnLevelUp automatically learns abilities when leveling up
func (s *AbilityEquipService) LearnAbilityOnLevelUp(characterID uint, newLevel int, rank string) ([]models.Ability, error) {
	// Find all abilities the character can learn at this level
	var learnable []models.AbilityLearning
	if err := db.DB.Where("learn_level <= ? AND min_rank <= ?", newLevel, rank).
		Preload("Ability").
		Find(&learnable).Error; err != nil {
		return nil, err
	}

	var newlyLearned []models.Ability

	// Check each ability
	for _, learning := range learnable {
		// Check if already learned
		var existing models.CharacterLearnedAbility
		err := db.DB.Where("character_id = ? AND ability_id = ?", characterID, learning.AbilityID).
			First(&existing).Error

		if err != nil {
			// Not learned yet, learn it now
			learned := models.CharacterLearnedAbility{
				CharacterID: characterID,
				AbilityID:   learning.AbilityID,
			}
			if err := db.DB.Create(&learned).Error; err != nil {
				continue // Skip if error
			}
			newlyLearned = append(newlyLearned, learning.Ability)
		}
	}

	return newlyLearned, nil
}

// GetAvailableAbilitiesToLearn returns abilities character can learn but hasn't yet
func (s *AbilityEquipService) GetAvailableAbilitiesToLearn(characterID uint) ([]models.Ability, error) {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return nil, err
	}

	// Get all abilities character CAN learn
	var learnable []models.AbilityLearning
	if err := db.DB.Where("learn_level <= ? AND min_rank <= ?", character.Level, character.Rarity).
		Preload("Ability").
		Find(&learnable).Error; err != nil {
		return nil, err
	}

	// Get already learned
	var learned []models.CharacterLearnedAbility
	db.DB.Where("character_id = ?", characterID).Find(&learned)
	learnedMap := make(map[uint]bool)
	for _, l := range learned {
		learnedMap[l.AbilityID] = true
	}

	// Filter out already learned
	var available []models.Ability
	for _, learning := range learnable {
		if !learnedMap[learning.AbilityID] {
			available = append(available, learning.Ability)
		}
	}

	return available, nil
}
