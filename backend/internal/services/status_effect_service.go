package services

import (
	"errors"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// StatusEffectService handles status effect logic
type StatusEffectService struct{}

// NewStatusEffectService creates a new service
func NewStatusEffectService() *StatusEffectService {
	return &StatusEffectService{}
}

// ApplyEffect applies a status effect to a character
func (s *StatusEffectService) ApplyEffect(characterID uint, effectName string, duration int, casterID *uint) error {
	// Get effect definition
	def, exists := models.GetEffectDefinition(effectName)
	if !exists {
		return errors.New("invalid effect name")
	}

	// Check for existing effect
	var existing models.StatusEffect
	err := db.DB.Where("character_id = ? AND effect_name = ?", characterID, effectName).
		First(&existing).Error

	if err == nil {
		// Effect already exists, increase stacks (max 3)
		if existing.Stacks < 3 {
			existing.Stacks++
			existing.Duration = duration // Refresh duration
			existing.TurnsRemaining = duration
			existing.ExpiresAt = time.Now().Add(time.Duration(duration) * time.Second)
			return db.DB.Save(&existing).Error
		}
		// Already at max stacks, just refresh duration
		existing.Duration = duration
		existing.TurnsRemaining = duration
		existing.ExpiresAt = time.Now().Add(time.Duration(duration) * time.Second)
		return db.DB.Save(&existing).Error
	}

	// Create new effect
	if duration == 0 {
		duration = def.DefaultDuration
	}

	effect := models.StatusEffect{
		EffectType:     def.Type,
		EffectName:     effectName,
		CharacterID:    characterID,
		Stacks:         1,
		Duration:       duration,
		MaxDuration:    duration,
		TurnsRemaining: duration,
		ExpiresAt:      time.Now().Add(time.Duration(duration) * time.Second),
		CasterID:       casterID,
		StatModifier:   def.Modifier,
		DamagePerTurn:  def.DamagePerTurn,
	}

	return db.DB.Create(&effect).Error
}

// RemoveEffect removes a status effect
func (s *StatusEffectService) RemoveEffect(characterID uint, effectName string) error {
	return db.DB.Where("character_id = ? AND effect_name = ?", characterID, effectName).
		Delete(&models.StatusEffect{}).Error
}

// GetActiveEffects returns all active effects for a character
func (s *StatusEffectService) GetActiveEffects(characterID uint) ([]models.StatusEffect, error) {
	var effects []models.StatusEffect
	err := db.DB.Where("character_id = ? AND expires_at > ?", characterID, time.Now()).
		Order("effect_type ASC, effect_name ASC").
		Find(&effects).Error
	return effects, err
}

// ProcessTurnEffects processes all turn-based effects (DoT, duration reduction)
func (s *StatusEffectService) ProcessTurnEffects(characterID uint) (int, []string, error) {
	effects, err := s.GetActiveEffects(characterID)
	if err != nil {
		return 0, nil, err
	}

	totalDamage := 0
	var expiredEffects []string

	// Get character for HP calculations
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return 0, nil, err
	}

	for i := range effects {
		effect := &effects[i]

		// Apply DoT damage
		if effect.DamagePerTurn > 0 {
			damage := int64(float64(character.CurrentHP) * (float64(effect.DamagePerTurn) / 100.0) * float64(effect.Stacks))
			totalDamage += int(damage)
		}

		// Reduce turn duration
		effect.TurnsRemaining--

		if effect.TurnsRemaining <= 0 {
			// Effect expired
			expiredEffects = append(expiredEffects, effect.EffectName)
			db.DB.Delete(effect)
		} else {
			// Save updated duration
			db.DB.Save(effect)
		}
	}

	return totalDamage, expiredEffects, nil
}

// HasEffect checks if character has a specific effect
func (s *StatusEffectService) HasEffect(characterID uint, effectName string) bool {
	var count int64
	db.DB.Model(&models.StatusEffect{}).
		Where("character_id = ? AND effect_name = ? AND expires_at > ?",
			characterID, effectName, time.Now()).
		Count(&count)
	return count > 0
}

// GetStatModifier returns the total stat modifier from all active effects
func (s *StatusEffectService) GetStatModifier(characterID uint, statType string) float64 {
	effects, _ := s.GetActiveEffects(characterID)

	totalModifier := 0.0
	for _, effect := range effects {
		switch statType {
		case "ATTACK":
			if effect.EffectName == "AMPED" || effect.EffectName == "FEEBLE" {
				totalModifier += effect.StatModifier * float64(effect.Stacks)
			}
		case "DEFENSE":
			if effect.EffectName == "BULKED" || effect.EffectName == "FRAGILE" {
				totalModifier += effect.StatModifier * float64(effect.Stacks)
			}
		case "SPEED":
			if effect.EffectName == "HASTE" || effect.EffectName == "SLOW" {
				totalModifier += effect.StatModifier * float64(effect.Stacks)
			}
		}
	}

	return totalModifier
}

// ClearAllEffects removes all effects from a character
func (s *StatusEffectService) ClearAllEffects(characterID uint) error {
	return db.DB.Where("character_id = ?", characterID).Delete(&models.StatusEffect{}).Error
}

// CleanupExpired removes all expired effects
func (s *StatusEffectService) CleanupExpired() error {
	return db.DB.Where("expires_at < ?", time.Now()).Delete(&models.StatusEffect{}).Error
}
