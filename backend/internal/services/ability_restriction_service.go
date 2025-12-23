package services

import "fmt"

// AbilityRestrictionService handles ability learning restrictions based on rank
type AbilityRestrictionService struct {
	configService *ConfigService
}

// NewAbilityRestrictionService creates a new ability restriction service
func NewAbilityRestrictionService(configService *ConfigService) *AbilityRestrictionService {
	return &AbilityRestrictionService{
		configService: configService,
	}
}

// GetMaxAbilitySlots returns the maximum number of abilities a character can have based on rank
func (s *AbilityRestrictionService) GetMaxAbilitySlots(characterRank string) int {
	// Check config first for customizable values
	configKey := fmt.Sprintf("ability_slots_%s", characterRank)
	slotCount := s.configService.GetInt(configKey, 0)

	if slotCount > 0 {
		return slotCount
	}

	// Default values if not in config
	slots := map[string]int{
		"C":   4,
		"B":   6,
		"A":   8,
		"S":   10,
		"SS":  12,
		"SSS": 16,
	}

	if count, ok := slots[characterRank]; ok {
		return count
	}

	return 4 // Default to C rank
}

// CanLearnAbility checks if a character of given rank can learn an ability of given rarity
// Higher ranks can learn lower rank abilities, but not vice versa
func (s *AbilityRestrictionService) CanLearnAbility(characterRank string, abilityRarity string) bool {
	rankOrder := []string{"C", "B", "A", "S", "SS", "SSS"}

	charIndex := indexOf(rankOrder, characterRank)
	abilityIndex := indexOf(rankOrder, abilityRarity)

	// If either rank not found, default to false for safety
	if charIndex == -1 || abilityIndex == -1 {
		return false
	}

	// Character rank must be >= ability rarity
	// E.g., rank S (index 3) can learn rank B (index 1) abilities
	return charIndex >= abilityIndex
}

// ValidateAbilityLearning performs comprehensive validation for ability learning
func (s *AbilityRestrictionService) ValidateAbilityLearning(characterRank string, characterLevel int, currentAbilityCount int, abilityRarity string, abilityUnlockLevel int) error {
	// 1. Check rank restriction
	if !s.CanLearnAbility(characterRank, abilityRarity) {
		return fmt.Errorf("rank %s characters cannot learn rank %s abilities", characterRank, abilityRarity)
	}

	// 2. Check level requirement
	if characterLevel < abilityUnlockLevel {
		return fmt.Errorf("requires character level %d (current: %d)", abilityUnlockLevel, characterLevel)
	}

	// 3. Check ability slot limit
	maxSlots := s.GetMaxAbilitySlots(characterRank)
	if currentAbilityCount >= maxSlots {
		return fmt.Errorf("ability slots full (%d/%d) - rank %s limit", currentAbilityCount, maxSlots, characterRank)
	}

	return nil
}

// GetLearnableAbilityRarities returns list of rarities a character can learn
func (s *AbilityRestrictionService) GetLearnableAbilityRarities(characterRank string) []string {
	rankOrder := []string{"C", "B", "A", "S", "SS", "SSS"}
	charIndex := indexOf(rankOrder, characterRank)

	if charIndex == -1 {
		return []string{"C"} // Default to C only
	}

	// Return all rarities up to and including character's rank
	return rankOrder[:charIndex+1]
}

// Helper function to find index in slice
func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}
