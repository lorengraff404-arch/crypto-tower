package services

import (
	"errors"
	"fmt"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/formulas"
)

// ProgressionService handles character progression logic
type ProgressionService struct{}

// NewProgressionService creates a new progression service
func NewProgressionService() *ProgressionService {
	return &ProgressionService{}
}

// GainExperience adds XP to character and handles level ups
func (s *ProgressionService) GainExperience(characterID uint, xpGained int, source string, difficulty string) (*models.Character, error) {
	// Anti-cheat: Validate XP gain
	if !formulas.ValidateXPGain(source, difficulty, xpGained) {
		return nil, fmt.Errorf("invalid XP gain: %d from %s/%s", xpGained, source, difficulty)
	}

	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return nil, err
	}

	// Add XP
	character.TotalXP += xpGained

	// Calculate new level based on total XP
	newLevel := formulas.GetLevelFromXP(character.TotalXP)
	leveledUp := newLevel > character.Level

	if leveledUp {
		// Update level
		oldLevel := character.Level
		character.Level = newLevel

		// Check for rarity upgrade
		newRarity := formulas.GetRarityForLevel(newLevel)
		rarityUpgraded := newRarity != character.Rarity
		if rarityUpgraded {
			character.Rarity = newRarity
		}

		// Check for evolution
		newEvolutionStage := formulas.GetEvolutionStage(newLevel)
		evolved := newEvolutionStage > character.EvolutionStage
		if evolved {
			character.EvolutionStage = newEvolutionStage
		}

		// Recalculate all stats
		att, def, hp, spd := formulas.RecalculateAllStats(
			int(character.BaseAttack),
			int(character.BaseDefense),
			int(character.BaseHP),
			int(character.BaseSpeed),
			character.Level,
			character.Rarity,
			character.EvolutionStage,
		)
		character.CurrentAttack = att
		character.CurrentDefense = def
		character.CurrentHP = hp
		character.CurrentSpeed = spd

		// SKILL SYSTEM INTEGRATION: Update mana scaling using ManaService
		manaService := NewManaService()
		character.MaxMana = manaService.CalculateMaxMana(character.Rarity, newLevel)
		character.ManaRegenRate = manaService.GetManaRegenRate(character.Rarity, newLevel)
		character.CurrentMana = character.MaxMana // Restore to full on level up

		// Auto-learn new abilities
		abilityService := NewAbilityService()
		newAbilities, err := abilityService.AutoLearnAbilities(characterID, newLevel)
		if err == nil && len(newAbilities) > 0 {
			fmt.Printf("Character %d learned %d new abilities!\n", characterID, len(newAbilities))
		}

		// SKILL SYSTEM INTEGRATION: Unlock skill slots
		skillInitService := NewSkillInitializationService()
		if err := skillInitService.CheckAndUnlockSlots(characterID, newLevel); err != nil {
			fmt.Printf("Warning: Failed to unlock skill slots: %v\n", err)
		}

		// Log the level up for auditing
		fmt.Printf("Character %d leveled up: %d → %d (Rarity: %s, Evolution: %d, Mana: %d/%d)\n",
			characterID, oldLevel, newLevel, character.Rarity, character.EvolutionStage, character.CurrentMana, character.MaxMana)
	}

	// Update current level XP for progress bar
	currentLevelXP := formulas.GetXPForLevel(character.Level)
	character.Experience = character.TotalXP - currentLevelXP

	// Save to database
	if err := db.DB.Save(&character).Error; err != nil {
		return nil, err
	}

	return &character, nil
}

// ValidateCharacterIntegrity checks if character stats are valid
func (s *ProgressionService) ValidateCharacterIntegrity(characterID uint) error {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return err
	}

	// Validate level matches total XP
	expectedLevel := formulas.GetLevelFromXP(character.TotalXP)
	if character.Level != expectedLevel {
		return fmt.Errorf("level mismatch: expected %d, got %d", expectedLevel, character.Level)
	}

	// Validate rarity matches level
	expectedRarity := formulas.GetRarityForLevel(character.Level)
	if character.Rarity != expectedRarity {
		return fmt.Errorf("rarity mismatch: expected %s, got %s", expectedRarity, character.Rarity)
	}

	// Validate evolution stage
	expectedEvolution := formulas.GetEvolutionStage(character.Level)
	if character.EvolutionStage != expectedEvolution {
		return fmt.Errorf("evolution mismatch: expected %d, got %d", expectedEvolution, character.EvolutionStage)
	}

	// Validate stats
	if !formulas.ValidateStats(character.BaseAttack, character.CurrentAttack, character.Level, character.Rarity, character.EvolutionStage) {
		return errors.New("attack stat tampering detected")
	}

	return nil
}

// GetProgressionInfo returns progression details for a character
func (s *ProgressionService) GetProgressionInfo(characterID uint) (map[string]interface{}, error) {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return nil, err
	}

	xpForNext := formulas.GetXPForNextLevel(character.Level, character.TotalXP)
	progress := formulas.GetXPProgressPercent(character.Level, character.TotalXP)

	return map[string]interface{}{
		"level":             character.Level,
		"total_xp":          character.TotalXP,
		"current_level_xp":  character.Experience,
		"xp_for_next_level": xpForNext,
		"progress_percent":  progress,
		"rarity":            character.Rarity,
		"evolution_stage":   character.EvolutionStage,
		"can_level_up":      character.Level < 100,
	}, nil
}

// RecalculateStats forcefully recalculates character stats (admin function)
func (s *ProgressionService) RecalculateStats(characterID uint) error {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return err
	}

	// Recalculate level from total XP
	character.Level = formulas.GetLevelFromXP(character.TotalXP)

	// Update rarity
	character.Rarity = formulas.GetRarityForLevel(character.Level)

	// Update evolution
	character.EvolutionStage = formulas.GetEvolutionStage(character.Level)

	// Recalculate all stats
	// Assuming formulas.RecalculateAllStats expects int for base stats and returns int
	// and character.Base/Current stats are int64
	att, def, hp, spd := formulas.RecalculateAllStats(
		int(character.BaseAttack),
		int(character.BaseDefense),
		int(character.BaseHP),
		int(character.BaseSpeed),
		character.Level,
		character.Rarity,
		character.EvolutionStage,
	)
	character.CurrentAttack = att
	character.CurrentDefense = def
	character.CurrentHP = hp
	character.CurrentSpeed = spd

	return db.DB.Save(&character).Error
}

// GetBaseManaForRarity returns base mana for a rarity
func (s *ProgressionService) GetBaseManaForRarity(rarity string) int {
	switch rarity {
	case "C":
		return 80
	case "B":
		return 100
	case "A":
		return 120
	case "S":
		return 150
	case "SS":
		return 200
	case "SSS":
		return 300
	default:
		return 100
	}
}

// GetBaseManaRegenForRarity returns base mana regen for a rarity
func (s *ProgressionService) GetBaseManaRegenForRarity(rarity string) int {
	switch rarity {
	case "C":
		return 8
	case "B":
		return 10
	case "A":
		return 12
	case "S":
		return 15
	case "SS":
		return 20
	case "SSS":
		return 30
	default:
		return 10
	}
}

// CalculateManaAtLevel calculates max mana at a given level
func (s *ProgressionService) CalculateManaAtLevel(baseMana int, level int) int {
	// Mana increases 5% per level
	return int(float64(baseMana) * (1.0 + float64(level-1)*0.05))
}

// CalculateManaRegenAtLevel calculates mana regen at a given level
func (s *ProgressionService) CalculateManaRegenAtLevel(baseRegen int, level int) int {
	// Mana regen increases 3% per level
	return int(float64(baseRegen) * (1.0 + float64(level-1)*0.03))
}
