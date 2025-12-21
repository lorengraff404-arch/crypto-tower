package services

import (
	"errors"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// CharacterService handles character-related business logic
type CharacterService struct{}

// NewCharacterService creates a new character service
func NewCharacterService() *CharacterService {
	return &CharacterService{}
}

// CreateCharacter creates a new character (from egg hatching)
func (s *CharacterService) CreateCharacter(ownerID uint, charType, element, rarity, class string) (*models.Character, error) {
	// Calculate base stats based on rarity
	baseStats := s.calculateBaseStats(rarity)

	// Initialize ManaService for mana calculations
	manaService := NewManaService()
	maxMana := manaService.CalculateMaxMana(rarity, 1) // Level 1
	manaRegen := manaService.GetManaRegenRate(rarity, 1)

	character := &models.Character{
		OwnerID:        ownerID,
		CharacterType:  charType,
		Element:        element,
		Rarity:         rarity,
		Class:          class,
		BaseAttack:     baseStats.Attack,
		BaseDefense:    baseStats.Defense,
		BaseHP:         baseStats.HP,
		BaseSpeed:      baseStats.Speed,
		CurrentAttack:  baseStats.Attack,
		CurrentDefense: baseStats.Defense,
		CurrentHP:      baseStats.HP,
		CurrentSpeed:   baseStats.Speed,
		Level:          1,
		Experience:     0,
		Durability:     100,
		Fatigue:        0,
		Abilities:      "[]", // Empty JSON array, will be populated based on type/element

		// Mana System - Calculated based on rarity
		MaxMana:       maxMana,
		CurrentMana:   maxMana, // Start with full mana
		ManaRegenRate: manaRegen,
	}

	if err := db.DB.Create(character).Error; err != nil {
		return nil, err
	}

	// Phase 10.2: Assign default moves based on element
	s.AssignDefaultMoves(character)

	// Phase 2: Auto-learn Level 1 abilities (Fix for Locked abilities)
	// Ensure Slash/Block etc are learned immediately
	abilityService := NewAbilityService()
	_, err := abilityService.AutoLearnAbilities(character.ID, 1)
	if err != nil {
		// Log error but don't fail creation, strictly speaking
		// fmt.Printf("Failed to auto-learn abilities: %v\n", err)
	}

	return character, nil
}

// AssignDefaultMoves gives a character 4 basic moves based on element (Phase 10.3 Fix)
func (s *CharacterService) AssignDefaultMoves(char *models.Character) {
	moves := s.getMovesForElement(char.Element)
	for i, tmpl := range moves {
		move := models.CharacterMove{
			CharacterID:     char.ID,
			MoveSlot:        i + 1,
			Name:            tmpl.Name,
			Type:            tmpl.Type,
			Category:        tmpl.Category,
			Power:           tmpl.Power,
			Accuracy:        tmpl.Accuracy,
			BasePP:          tmpl.PP,
			CurrentPP:       tmpl.PP,
			Priority:        tmpl.Priority,
			EffectChance:    tmpl.EffectChance,
			EffectType:      tmpl.EffectType,
			EffectMagnitude: tmpl.EffectMagnitude,
			Description:     tmpl.Description,
			Animation:       tmpl.Animation,
		}
		db.DB.Create(&move)
	}
}

func (s *CharacterService) getMovesForElement(element string) []models.MoveTemplate {
	// Simple mapping based on models.DefaultMoves
	switch element {
	case "FIRE":
		return []models.MoveTemplate{models.DefaultMoves[0], models.DefaultMoves[1], models.DefaultMoves[3], models.DefaultMoves[17]}
	case "WATER":
		return []models.MoveTemplate{models.DefaultMoves[4], models.DefaultMoves[5], models.DefaultMoves[7], models.DefaultMoves[18]}
	case "GRASS":
		return []models.MoveTemplate{models.DefaultMoves[8], models.DefaultMoves[9], models.DefaultMoves[10], models.DefaultMoves[11]}
	case "ELECTRIC":
		return []models.MoveTemplate{models.DefaultMoves[12], models.DefaultMoves[13], models.DefaultMoves[15], models.DefaultMoves[19]}
	case "ICE":
		return []models.MoveTemplate{models.DefaultMoves[4], models.DefaultMoves[14], models.DefaultMoves[15], models.DefaultMoves[18]}
	case "DRAGON":
		return []models.MoveTemplate{models.DefaultMoves[2], models.DefaultMoves[6], models.DefaultMoves[16], models.DefaultMoves[17]}
	default:
		return []models.MoveTemplate{models.DefaultMoves[14], models.DefaultMoves[15], models.DefaultMoves[17], models.DefaultMoves[20]}
	}
}

// GetUserCharacters retrieves all characters for a user
func (s *CharacterService) GetUserCharacters(userID uint, includeEggs bool) ([]models.Character, error) {
	var characters []models.Character
	query := db.DB.Where("owner_id = ?", userID)

	if !includeEggs {
		query = query.Where("is_egg = ?", false)
	}

	if err := query.Order("created_at DESC").Find(&characters).Error; err != nil {
		return nil, err
	}

	return characters, nil
}

// GetCharacterByID retrieves a character by ID
func (s *CharacterService) GetCharacterByID(characterID, ownerID uint) (*models.Character, error) {
	var character models.Character
	if err := db.DB.Where("id = ? AND owner_id = ?", characterID, ownerID).First(&character).Error; err != nil {
		return nil, err
	}
	return &character, nil
}

// HatchEgg converts an egg to a character
func (s *CharacterService) HatchEgg(characterID, ownerID uint) (*models.Character, error) {
	var character models.Character
	if err := db.DB.Where("id = ? AND owner_id = ?", characterID, ownerID).First(&character).Error; err != nil {
		return nil, errors.New("character not found")
	}

	if !character.IsEgg {
		return nil, errors.New("character is not an egg")
	}

	if character.HatchTime != nil && time.Now().Before(*character.HatchTime) {
		return nil, errors.New("egg is not ready to hatch")
	}

	// Apply care slot bonuses
	bonusMultiplier := 1.0
	if character.CareSlotCalibrate {
		bonusMultiplier += 0.05
	}
	if character.CareSlotNurture {
		bonusMultiplier += 0.05
	}
	if character.CareSlotStabilize {
		bonusMultiplier += 0.05
	}

	// Apply bonuses to base stats
	character.BaseAttack = int(float64(character.BaseAttack) * bonusMultiplier)
	character.BaseDefense = int(float64(character.BaseDefense) * bonusMultiplier)
	character.BaseHP = int(float64(character.BaseHP) * bonusMultiplier)
	character.BaseSpeed = int(float64(character.BaseSpeed) * bonusMultiplier)

	// Update current stats
	character.CurrentAttack = character.BaseAttack
	character.CurrentDefense = character.BaseDefense
	character.CurrentHP = character.BaseHP
	character.CurrentSpeed = character.BaseSpeed

	// Mark as hatched
	character.IsEgg = false
	now := time.Now()
	character.MintedAt = &now

	if err := db.DB.Save(&character).Error; err != nil {
		return nil, err
	}

	return &character, nil
}

// UpdateDurabilityAndFatigue updates character condition after battle
func (s *CharacterService) UpdateDurabilityAndFatigue(characterID uint, durabilityLoss, fatigueGain int) error {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return err
	}

	character.Durability -= durabilityLoss
	if character.Durability < 0 {
		character.Durability = 0
		character.IsDead = true
	}

	character.Fatigue += fatigueGain
	if character.Fatigue > 100 {
		character.Fatigue = 100
	}

	now := time.Now()
	character.LastBattleAt = &now

	return db.DB.Save(&character).Error
}

// RecoverFatigue reduces fatigue (farming mechanic)
func (s *CharacterService) RecoverFatigue(characterID uint, amount int) error {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return err
	}

	character.Fatigue -= amount
	if character.Fatigue < 0 {
		character.Fatigue = 0
	}

	return db.DB.Save(&character).Error
}

// calculateBaseStats returns rebalanced base stats (Phase 11 Match Update)
func (s *CharacterService) calculateBaseStats(rarity string) struct{ Attack, Defense, HP, Speed int } {
	stats := struct{ Attack, Defense, HP, Speed int }{}

	// STAT SQUISH: Lower numbers for easier mental math and balance
	// Baseline (C Rarity Level 1): 20 Atk, 20 Def, 100 HP, 20 Speed
	switch rarity {
	case "SSS":
		stats.Attack = 60
		stats.Defense = 60
		stats.HP = 250
		stats.Speed = 50
	case "SS":
		stats.Attack = 50
		stats.Defense = 50
		stats.HP = 210
		stats.Speed = 45
	case "S":
		stats.Attack = 42
		stats.Defense = 42
		stats.HP = 180
		stats.Speed = 40
	case "A":
		stats.Attack = 35
		stats.Defense = 35
		stats.HP = 150
		stats.Speed = 35
	case "B":
		stats.Attack = 28
		stats.Defense = 28
		stats.HP = 125
		stats.Speed = 30
	case "C":
		stats.Attack = 20
		stats.Defense = 20
		stats.HP = 100
		stats.Speed = 25
	default:
		stats.Attack = 20
		stats.Defense = 20
		stats.HP = 100
		stats.Speed = 25
	}

	return stats
}
