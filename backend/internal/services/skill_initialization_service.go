package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// SkillInitializationService handles initial skill assignment for new characters
type SkillInitializationService struct{}

// NewSkillInitializationService creates a new skill initialization service
func NewSkillInitializationService() *SkillInitializationService {
	return &SkillInitializationService{}
}

// InitializeCharacterSkills assigns initial skills to a newly created character
func (s *SkillInitializationService) InitializeCharacterSkills(characterID uint) error {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	// Get skill pool for this character
	skillPool, err := s.GetSkillPoolForCharacter(&character)
	if err != nil {
		return err
	}

	// Determine how many skills to assign based on rarity
	totalSkills := s.GetTotalSkillsByRarity(character.Rarity)
	activeSlots := s.GetActiveSlotsByRarity(character.Rarity)

	// Randomly select skills from pool
	selectedSkills := s.RandomSelectSkills(skillPool, totalSkills)

	// Assign skills to character
	for i, skill := range selectedSkills {
		// Create character_ability record
		charAbility := models.CharacterAbility{
			CharacterID: characterID,
			AbilityID:   skill.ID,
		}
		if err := db.DB.Create(&charAbility).Error; err != nil {
			return fmt.Errorf("failed to assign skill: %w", err)
		}

		// Create active skill slot if within active slot limit
		if i < activeSlots {
			activeSkill := models.CharacterActiveSkill{
				CharacterID:  characterID,
				AbilityID:    skill.ID,
				SlotPosition: i + 1,
				IsLocked:     i < 2, // First 2 slots unlocked by default
				UnlockLevel:  s.GetUnlockLevelForSlot(i, character.Rarity),
			}
			if err := db.DB.Create(&activeSkill).Error; err != nil {
				return fmt.Errorf("failed to create active slot: %w", err)
			}
		}
	}

	return nil
}

// GetSkillPoolForCharacter returns available skills for a character
func (s *SkillInitializationService) GetSkillPoolForCharacter(char *models.Character) ([]models.Ability, error) {
	var skills []models.Ability

	// Get skills that match character's class and rarity
	query := db.DB.Where("rarity <= ?", char.Rarity)

	// Filter by class if skill has class requirement
	query = query.Where("required_class IS NULL OR ? = ANY(required_class)", char.Class)

	// Filter by element if skill has element requirement
	query = query.Where("required_element IS NULL OR ? = ANY(required_element)", char.Element)

	// Get skills up to character's level
	query = query.Where("unlock_level <= ?", char.Level)

	if err := query.Find(&skills).Error; err != nil {
		return nil, err
	}

	return skills, nil
}

// GetTotalSkillsByRarity returns total skills a character can have
func (s *SkillInitializationService) GetTotalSkillsByRarity(rarity string) int {
	switch rarity {
	case "C":
		return 4
	case "B":
		return 6
	case "A":
		return 8
	case "S":
		return 12
	case "SS":
		return 20
	case "SSS":
		return 30
	default:
		return 4
	}
}

// GetActiveSlotsByRarity returns active skill slots for a rarity
func (s *SkillInitializationService) GetActiveSlotsByRarity(rarity string) int {
	switch rarity {
	case "C":
		return 2
	case "B":
		return 3
	case "A":
		return 4
	case "S":
		return 5
	case "SS":
		return 6
	case "SSS":
		return 7
	default:
		return 2
	}
}

// GetUnlockLevelForSlot returns the level required to unlock a slot
func (s *SkillInitializationService) GetUnlockLevelForSlot(slotIndex int, rarity string) int {
	// First 2 slots always unlocked at level 1
	if slotIndex < 2 {
		return 1
	}

	// Progressive unlock based on rarity
	switch rarity {
	case "C":
		return []int{1, 1, 5, 10}[slotIndex]
	case "B":
		return []int{1, 1, 5, 10, 15, 20}[slotIndex]
	case "A":
		return []int{1, 1, 8, 12, 16, 20, 25, 30}[slotIndex]
	case "S":
		return []int{1, 1, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55}[slotIndex]
	case "SS":
		// More slots, progressive unlocking
		if slotIndex < 6 {
			return 1 + (slotIndex * 5)
		}
		return 30 + ((slotIndex - 5) * 10)
	case "SSS":
		// Most slots, progressive unlocking
		if slotIndex < 7 {
			return 1 + (slotIndex * 5)
		}
		return 35 + ((slotIndex - 6) * 8)
	default:
		return 1
	}
}

// RandomSelectSkills randomly selects N skills from pool
func (s *SkillInitializationService) RandomSelectSkills(pool []models.Ability, count int) []models.Ability {
	if len(pool) <= count {
		return pool
	}

	// Shuffle pool
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]models.Ability, len(pool))
	copy(shuffled, pool)

	for i := range shuffled {
		j := rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	// Return first N skills
	return shuffled[:count]
}

// UnlockSkillSlot unlocks a skill slot when character levels up
func (s *SkillInitializationService) UnlockSkillSlot(characterID uint, slotPosition int) error {
	return db.DB.Model(&models.CharacterActiveSkill{}).
		Where("character_id = ? AND slot_position = ?", characterID, slotPosition).
		Update("is_locked", false).Error
}

// CheckAndUnlockSlots checks if any slots should be unlocked based on level
func (s *SkillInitializationService) CheckAndUnlockSlots(characterID uint, newLevel int) error {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return err
	}

	// Get all active skill slots
	var activeSlots []models.CharacterActiveSkill
	if err := db.DB.Where("character_id = ?", characterID).Find(&activeSlots).Error; err != nil {
		return err
	}

	// Unlock slots that meet level requirement
	for _, slot := range activeSlots {
		if slot.IsLocked && newLevel >= slot.UnlockLevel {
			if err := s.UnlockSkillSlot(characterID, slot.SlotPosition); err != nil {
				return err
			}
		}
	}

	return nil
}
