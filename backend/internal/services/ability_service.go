package services

import (
	"encoding/json"
	"fmt"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// AbilityService handles ability-related business logic
type AbilityService struct{}

// NewAbilityService creates a new ability service
func NewAbilityService() *AbilityService {
	return &AbilityService{}
}

// GetAbilitiesByClass returns all abilities for a specific class
func (s *AbilityService) GetAbilitiesByClass(class string) ([]models.Ability, error) {
	var abilities []models.Ability
	if err := db.DB.Where("class = ?", class).Order("unlock_level ASC").Find(&abilities).Error; err != nil {
		return nil, err
	}
	return abilities, nil
}

// GetAvailableAbilities returns abilities a character can learn at their current level
// SECURITY: Now includes rank-based filtering to prevent learning abilities above character's rank
func (s *AbilityService) GetAvailableAbilities(characterID uint) ([]models.Ability, error) {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return nil, err
	}

	// Initialize restriction service
	restrictionService := NewAbilityRestrictionService(GetConfigService())
	learnableRarities := restrictionService.GetLearnableAbilityRarities(character.Rarity)

	// Get all abilities for this character's class up to their level and within their rank
	// SECURITY: Filter by class, level, AND rarity to prevent rank exploitation
	var abilities []models.Ability
	if err := db.DB.Where("class = ? AND unlock_level <= ? AND rarity IN ?",
		character.Class, character.Level, learnableRarities).
		Order("unlock_level ASC").Find(&abilities).Error; err != nil {
		return nil, err
	}

	return abilities, nil
}

// GetLearnedAbilities returns abilities a character has already learned
func (s *AbilityService) GetLearnedAbilities(characterID uint) ([]models.Ability, error) {
	var characterAbilities []models.CharacterAbility
	if err := db.DB.Where("character_id = ?", characterID).
		Preload("Ability").Find(&characterAbilities).Error; err != nil {
		return nil, err
	}

	abilities := make([]models.Ability, len(characterAbilities))
	for i, ca := range characterAbilities {
		abilities[i] = ca.Ability
	}

	return abilities, nil
}

// AutoLearnAbilities automatically learns abilities when character levels up
// SECURITY: Now enforces rank restrictions and slot limits for game balance
func (s *AbilityService) AutoLearnAbilities(characterID uint, characterLevel int) ([]models.Ability, error) {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return nil, err
	}

	// Initialize restriction service for validation
	restrictionService := NewAbilityRestrictionService(GetConfigService())

	// Get all abilities the character should have learned by now
	// SECURITY: Only fetch abilities for character's own class to prevent ability theft
	var availableAbilities []models.Ability
	if err := db.DB.Where("class = ? AND unlock_level <= ?", character.Class, characterLevel).
		Order("unlock_level ASC, rarity ASC").Find(&availableAbilities).Error; err != nil {
		return nil, err
	}

	// Get already learned abilities to check slots
	var learnedAbilityIDs []uint
	var characterAbilities []models.CharacterAbility
	if err := db.DB.Where("character_id = ?", characterID).
		Find(&characterAbilities).Error; err != nil {
		return nil, err
	}

	for _, ca := range characterAbilities {
		learnedAbilityIDs = append(learnedAbilityIDs, ca.AbilityID)
	}

	// SECURITY: Check slot limit to prevent ability hoarding
	maxSlots := restrictionService.GetMaxAbilitySlots(character.Rarity)
	currentSlots := len(learnedAbilityIDs)

	// Learn new abilities (with rank restrictions)
	newlyLearned := []models.Ability{}
	for _, ability := range availableAbilities {
		// Check if slot limit reached
		if currentSlots >= maxSlots {
			break // Stop learning if slots are full
		}

		// SECURITY: Validate character can learn this ability rarity
		if !restrictionService.CanLearnAbility(character.Rarity, ability.Rarity) {
			continue // Skip abilities above character's rank
		}

		// Check if already learned
		alreadyLearned := false
		for _, learnedID := range learnedAbilityIDs {
			if learnedID == ability.ID {
				alreadyLearned = true
				break
			}
		}

		if !alreadyLearned {
			// Learn this ability
			characterAbility := models.CharacterAbility{
				CharacterID: characterID,
				AbilityID:   ability.ID,
			}
			if err := db.DB.Create(&characterAbility).Error; err != nil {
				return nil, err
			}
			newlyLearned = append(newlyLearned, ability)
			currentSlots++ // Increment slot count
		}
	}

	return newlyLearned, nil
}

// CalculateAbilityDamage calculates damage with element bonuses
func (s *AbilityService) CalculateAbilityDamage(ability *models.Ability, characterElement string) int {
	if ability.BaseDamage == 0 {
		return 0
	}

	// Parse element bonuses
	var elementBonus models.ElementBonus
	if err := json.Unmarshal([]byte(ability.ElementBonuses), &elementBonus); err != nil {
		return ability.BaseDamage // Return base if can't parse
	}

	// Get multiplier for character's element
	multiplier := 1.0
	switch characterElement {
	case "Fire":
		multiplier = elementBonus.Fire
	case "Water":
		multiplier = elementBonus.Water
	case "Ice":
		multiplier = elementBonus.Ice
	case "Thunder":
		multiplier = elementBonus.Thunder
	case "Dark":
		multiplier = elementBonus.Dark
	case "Plant":
		multiplier = elementBonus.Plant
	case "Earth":
		multiplier = elementBonus.Earth
	case "Wind":
		multiplier = elementBonus.Wind
	}

	return int(float64(ability.BaseDamage) * multiplier)
}

// GetAbilityDetails returns detailed info about an ability including element-modified stats
func (s *AbilityService) GetAbilityDetails(abilityID uint, characterElement string) (map[string]interface{}, error) {
	var ability models.Ability
	if err := db.DB.First(&ability, abilityID).Error; err != nil {
		return nil, err
	}

	damage := s.CalculateAbilityDamage(&ability, characterElement)

	return map[string]interface{}{
		"ability":           ability,
		"calculated_damage": damage,
		"element_bonus":     fmt.Sprintf("%.0f%%", (float64(damage)/float64(ability.BaseDamage)-1)*100),
	}, nil
}
