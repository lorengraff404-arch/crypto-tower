package services

import (
	"errors"
	"fmt"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// AbilityEquipService handles ability equipping logic with rank-based slot limits
type AbilityEquipService struct {
	restrictionService *AbilityRestrictionService
}

// NewAbilityEquipService creates a new ability equip service
func NewAbilityEquipService() *AbilityEquipService {
	return &AbilityEquipService{
		restrictionService: NewAbilityRestrictionService(GetConfigService()),
	}
}

// EquipAbility equips an ability to a specific slot with validation
// SECURITY: Validates character ownership, slot limits, and ability learning status
func (s *AbilityEquipService) EquipAbility(userID, characterID, abilityID uint, slot int) error {
	// Get character for validation
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	// SECURITY: Verify ownership
	if character.OwnerID != userID {
		return errors.New("you do not own this character")
	}

	// SECURITY: Validate slot number based on character's rank
	maxSlots := s.restrictionService.GetMaxAbilitySlots(character.Rarity)
	if slot < 1 || slot > maxSlots {
		return fmt.Errorf("invalid slot: must be 1-%d for rank %s", maxSlots, character.Rarity)
	}

	// SECURITY: Verify character has learned this ability
	var learned models.CharacterAbility
	if err := db.DB.Where("character_id = ? AND ability_id = ?", characterID, abilityID).
		First(&learned).Error; err != nil {
		return errors.New("ability not learned yet - cannot equip unlearned ability")
	}

	// Get ability to validate rank restriction
	var ability models.Ability
	if err := db.DB.First(&ability, abilityID).Error; err != nil {
		return fmt.Errorf("ability not found: %w", err)
	}

	// SECURITY: Validate character can use this ability rarity
	if !s.restrictionService.CanLearnAbility(character.Rarity, ability.Rarity) {
		return fmt.Errorf("rank %s cannot equip rank %s abilities", character.Rarity, ability.Rarity)
	}

	// Check if ability is already equipped somewhere
	var existing models.EquippedAbility
	err := db.DB.Where("character_id = ? AND ability_id = ?", characterID, abilityID).
		First(&existing).Error

	if err == nil {
		// Already equipped - update slot position
		existing.SlotPosition = slot
		return db.DB.Save(&existing).Error
	}

	// Check if slot is already occupied
	var occupiedSlot models.EquippedAbility
	err = db.DB.Where("character_id = ? AND slot_position = ?", characterID, slot).
		First(&occupiedSlot).Error

	if err == nil {
		// Slot occupied - delete old entry
		db.DB.Delete(&occupiedSlot)
	}

	// Equip to new slot
	equipped := models.EquippedAbility{
		CharacterID:  characterID,
		AbilityID:    abilityID,
		SlotPosition: slot,
	}

	if err := db.DB.Create(&equipped).Error; err != nil {
		return fmt.Errorf("failed to equip ability: %w", err)
	}

	return nil
}

// UnequipAbility removes an ability from equipment
func (s *AbilityEquipService) UnequipAbility(userID, characterID, abilityID uint) error {
	// SECURITY: Verify ownership first
	var char models.Character
	if err := db.DB.First(&char, characterID).Error; err != nil {
		return err
	}
	if char.OwnerID != userID {
		return errors.New("you do not own this character")
	}

	result := db.DB.Where("character_id = ? AND ability_id = ?", characterID, abilityID).
		Delete(&models.EquippedAbility{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("ability not equipped")
	}

	return nil
}

// UnequipSlot removes ability from a specific slot
func (s *AbilityEquipService) UnequipSlot(characterID uint, slot int) error {
	result := db.DB.Where("character_id = ? AND slot_position = ?", characterID, slot).
		Delete(&models.EquippedAbility{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("slot is empty")
	}

	return nil
}

// GetEquippedAbilities returns allequipped abilities for a character
func (s *AbilityEquipService) GetEquippedAbilities(characterID uint) ([]models.Ability, error) {
	var equipped []models.EquippedAbility
	if err := db.DB.Where("character_id = ?", characterID).
		Order("slot_position ASC").
		Preload("Ability").
		Find(&equipped).Error; err != nil {
		return nil, err
	}

	abilities := make([]models.Ability, len(equipped))
	for i, eq := range equipped {
		abilities[i] = eq.Ability
	}

	return abilities, nil
}

// GetEquippedAbilitiesWithSlots returns equipped abilities mapped to their slots
func (s *AbilityEquipService) GetEquippedAbilitiesWithSlots(characterID uint) (map[int]models.Ability, error) {
	var equipped []models.EquippedAbility
	if err := db.DB.Where("character_id = ?", characterID).
		Preload("Ability").
		Find(&equipped).Error; err != nil {
		return nil, err
	}

	slotMap := make(map[int]models.Ability)
	for _, eq := range equipped {
		slotMap[eq.SlotPosition] = eq.Ability
	}

	return slotMap, nil
}

// GetEquippedCount returns how many abilities are currently equipped
func (s *AbilityEquipService) GetEquippedCount(characterID uint) (int, error) {
	var count int64
	if err := db.DB.Model(&models.EquippedAbility{}).
		Where("character_id = ?", characterID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetLearnedAbilities returns ALL abilities a character has learned (unlimited)
func (s *AbilityEquipService) GetLearnedAbilities(characterID uint) ([]models.Ability, error) {
	var learned []models.CharacterAbility
	if err := db.DB.Where("character_id = ?", characterID).
		Preload("Ability").
		Order("id ASC").
		Find(&learned).Error; err != nil {
		return nil, err
	}

	abilities := make([]models.Ability, len(learned))
	for i, l := range learned {
		abilities[i] = l.Ability
	}

	return abilities, nil
}

// GetAvailableAbilitiesToLearn returns abilities character can learn but hasn't yet
func (s *AbilityEquipService) GetAvailableAbilitiesToLearn(characterID uint) ([]models.Ability, error) {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return nil, err
	}

	// Get learnable rarities for this character's rank
	learnableRarities := s.restrictionService.GetLearnableAbilityRarities(character.Rarity)

	// Get all abilities for character's class, level, and rank
	var allAbilities []models.Ability
	if err := db.DB.Where("class = ? AND unlock_level <= ? AND rarity IN ?",
		character.Class, character.Level, learnableRarities).
		Find(&allAbilities).Error; err != nil {
		return nil, err
	}

	// Get already learned abilities
	var learned []models.CharacterAbility
	db.DB.Where("character_id = ?", characterID).Find(&learned)
	learnedMap := make(map[uint]bool)
	for _, l := range learned {
		learnedMap[l.AbilityID] = true
	}

	// Filter out already learned
	var available []models.Ability
	for _, ability := range allAbilities {
		if !learnedMap[ability.ID] {
			available = append(available, ability)
		}
	}

	return available, nil
}
