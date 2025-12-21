package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// StatusEffectHandler handles status effect HTTP requests
type StatusEffectHandler struct {
	statusEffectService *services.StatusEffectService
}

// NewStatusEffectHandler creates a new handler
func NewStatusEffectHandler() *StatusEffectHandler {
	return &StatusEffectHandler{
		statusEffectService: services.NewStatusEffectService(),
	}
}

// GetCharacterEffects returns all active effects for a character
// GET /api/v1/characters/:id/effects
func (h *StatusEffectHandler) GetCharacterEffects(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	characterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	// Verify ownership
	characterService := services.NewCharacterService()
	char, err := characterService.GetCharacterByID(uint(characterID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	effects, err := h.statusEffectService.GetActiveEffects(uint(characterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch effects"})
		return
	}

	// Enhance with definitions
	enrichedEffects := make([]map[string]interface{}, 0)
	for _, effect := range effects {
		def, _ := models.GetEffectDefinition(effect.EffectName)
		enrichedEffects = append(enrichedEffects, map[string]interface{}{
			"id":              effect.ID,
			"effect_name":     effect.EffectName,
			"effect_type":     effect.EffectType,
			"stacks":          effect.Stacks,
			"turns_remaining": effect.TurnsRemaining,
			"icon":            def.Icon,
			"description":     def.Description,
			"stat_modifier":   effect.StatModifier,
			"damage_per_turn": effect.DamagePerTurn,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"character_id":   characterID,
		"character_name": char.Name,
		"effects":        enrichedEffects,
		"total_effects":  len(enrichedEffects),
	})
}

// GetAllEffectDefinitions returns all buff and debuff definitions
// GET /api/v1/effects/definitions
func (h *StatusEffectHandler) GetAllEffectDefinitions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"buffs":   models.BuffDefinitions,
		"debuffs": models.DebuffDefinitions,
	})
}
