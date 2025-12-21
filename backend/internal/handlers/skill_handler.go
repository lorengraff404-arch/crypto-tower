package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// SkillHandler handles skill-related HTTP requests
type SkillHandler struct {
	skillService *services.SkillActivationService
}

// NewSkillHandler creates a new skill handler
func NewSkillHandler() *SkillHandler {
	return &SkillHandler{
		skillService: services.NewSkillActivationService(),
	}
}

// ActivateSkill handles POST /api/v1/skills/activate
func (h *SkillHandler) ActivateSkill(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		CharacterID uint   `json:"character_id" binding:"required"`
		AbilityID   uint   `json:"ability_id" binding:"required"`
		TargetID    uint   `json:"target_id"`
		TargetIDs   []uint `json:"target_ids"`
		BattleID    uint   `json:"battle_id" binding:"required"`
		TurnNumber  int    `json:"turn_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify character ownership
	if !h.verifyCharacterOwnership(userID, req.CharacterID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your character"})
		return
	}

	// Activate skill
	activationReq := services.SkillActivationRequest{
		CharacterID: req.CharacterID,
		AbilityID:   req.AbilityID,
		TargetID:    req.TargetID,
		TargetIDs:   req.TargetIDs,
		BattleID:    req.BattleID,
		TurnNumber:  req.TurnNumber,
	}

	result, err := h.skillService.ActivateSkill(activationReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// GetActiveSkills handles GET /api/v1/characters/:id/active-skills
func (h *SkillHandler) GetActiveSkills(c *gin.Context) {
	characterID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	skills, err := h.skillService.GetActiveSkills(uint(characterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skills": skills,
	})
}

// SwapSkill handles POST /api/v1/characters/:id/swap-skill
func (h *SkillHandler) SwapSkill(c *gin.Context) {
	userID := c.GetUint("user_id")
	characterID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		OldAbilityID uint `json:"old_ability_id" binding:"required"`
		NewAbilityID uint `json:"new_ability_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership
	if !h.verifyCharacterOwnership(userID, uint(characterID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your character"})
		return
	}

	// Swap skill
	if err := h.skillService.SwapActiveSkill(uint(characterID), req.OldAbilityID, req.NewAbilityID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "skill swapped successfully",
	})
}

// GetSkillCooldowns handles GET /api/v1/characters/:id/cooldowns
func (h *SkillHandler) GetSkillCooldowns(c *gin.Context) {
	characterID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var cooldowns []struct {
		AbilityID         uint   `json:"ability_id"`
		AbilityName       string `json:"ability_name"`
		CooldownRemaining int    `json:"cooldown_remaining"`
	}

	// Query cooldowns from database
	// TODO: Implement actual query
	_ = characterID // Use the variable to avoid lint error for now

	c.JSON(http.StatusOK, gin.H{
		"cooldowns": cooldowns,
	})
}

// StartTurn handles POST /api/v1/battles/:id/start-turn
// Regenerates mana and reduces cooldowns
func (h *SkillHandler) StartTurn(c *gin.Context) {
	var req struct {
		CharacterID uint `json:"character_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Regenerate mana
	if err := h.skillService.RegenerateMana(req.CharacterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reduce cooldowns
	h.skillService.ReduceCooldowns(req.CharacterID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "turn started, mana regenerated, cooldowns reduced",
	})
}

// Helper function to verify character ownership
func (h *SkillHandler) verifyCharacterOwnership(userID, characterID uint) bool {
	// This would query the database
	// Simplified for now - always return true in development
	_ = userID      // TODO: Implement actual ownership check
	_ = characterID // TODO: Implement actual ownership check
	return true
}
