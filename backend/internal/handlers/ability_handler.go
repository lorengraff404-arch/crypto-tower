package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// AbilityHandler handles ability HTTP requests
type AbilityHandler struct {
	abilityService *services.AbilityService
	equipService   *services.AbilityEquipService
}

// NewAbilityHandler creates a new ability handler
func NewAbilityHandler() *AbilityHandler {
	return &AbilityHandler{
		abilityService: services.NewAbilityService(),
		equipService:   services.NewAbilityEquipService(),
	}
}

// GetAbilitiesByClass returns all abilities for a class
// GET /api/v1/abilities?class=Warrior
func (h *AbilityHandler) GetAbilitiesByClass(c *gin.Context) {
	class := c.Query("class")
	if class == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class parameter required"})
		return
	}

	abilities, err := h.abilityService.GetAbilitiesByClass(class)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch abilities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"class":     class,
		"abilities": abilities,
		"count":     len(abilities),
	})
}

// GetCharacterAbilities returns learned abilities for a character
// GET /api/v1/characters/:id/abilities
func (h *AbilityHandler) GetCharacterAbilities(c *gin.Context) {
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

	// Phase 10.5 Fix: Self-healing check
	// Ensure character has learned all abilities for their level
	if _, err := h.abilityService.AutoLearnAbilities(uint(characterID), char.Level); err != nil {
		// Log warning but proceed (non-critical)
		// log.Printf("AutoLearn warning for char %d: %v", characterID, err)
	}

	learnedAbilities, err := h.abilityService.GetLearnedAbilities(uint(characterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch learned abilities"})
		return
	}

	availableAbilities, err := h.abilityService.GetAvailableAbilities(uint(characterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available abilities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"character_id": characterID,
		"level":        char.Level,
		"class":        char.Class,
		"learned":      learnedAbilities,
		"available":    availableAbilities,
	})
}

// GetAbilityDetails returns detailed info about an ability with element bonuses
// GET /api/v1/abilities/:id?element=Fire
func (h *AbilityHandler) GetAbilityDetails(c *gin.Context) {
	abilityID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ability ID"})
		return
	}

	element := c.Query("element")
	if element == "" {
		element = "Fire" // Default
	}

	details, err := h.abilityService.GetAbilityDetails(uint(abilityID), element)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ability not found"})
		return
	}

	c.JSON(http.StatusOK, details)
}

// ==================== ABILITY LOADOUT ENDPOINTS ====================

// GetLearnedAbilities returns all abilities a character has learned (unlimited)
// GET /api/v1/abilities/learned?character_id=123
func (h *AbilityHandler) GetLearnedAbilities(c *gin.Context) {
	characterIDStr := c.Query("character_id")
	characterID, err := strconv.ParseUint(characterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character_id"})
		return
	}

	abilities, err := h.equipService.GetLearnedAbilities(uint(characterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load learned abilities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"character_id": characterID,
		"learned":      abilities,
		"count":        len(abilities),
	})
}

// GetEquippedAbilities returns currently equipped abilities for battle (slot-limited)
// GET /api/v1/abilities/equipped?character_id=123
func (h *AbilityHandler) GetEquippedAbilities(c *gin.Context) {
	characterIDStr := c.Query("character_id")
	characterID, err := strconv.ParseUint(characterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character_id"})
		return
	}

	abilities, err := h.equipService.GetEquippedAbilities(uint(characterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load equipped abilities"})
		return
	}

	// Get slot mapping for UI
	slotMap, _ := h.equipService.GetEquippedAbilitiesWithSlots(uint(characterID))

	c.JSON(http.StatusOK, gin.H{
		"character_id": characterID,
		"equipped":     abilities,
		"slots":        slotMap,
		"count":        len(abilities),
	})
}

// EquipAbility equips an ability to a specific slot
// POST /api/v1/abilities/equip
// Body: {"character_id": 123, "ability_id": 456, "slot": 1}
func (h *AbilityHandler) EquipAbility(c *gin.Context) {
	var req struct {
		CharacterID uint `json:"character_id" binding:"required"`
		AbilityID   uint `json:"ability_id" binding:"required"`
		Slot        int  `json:"slot" binding:"required,min=1,max=16"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	// SECURITY: Verify user owns character (Handled by Service)
	userID := c.GetUint("user_id")

	err := h.equipService.EquipAbility(userID, req.CharacterID, req.AbilityID, req.Slot)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "ability equipped successfully",
		"character_id": req.CharacterID,
		"ability_id":   req.AbilityID,
		"slot":         req.Slot,
	})
}

// UnequipAbility removes an ability from equipment
// DELETE /api/v1/abilities/unequip
// Body: {"character_id": 123, "ability_id": 456}
func (h *AbilityHandler) UnequipAbility(c *gin.Context) {
	var req struct {
		CharacterID uint `json:"character_id" binding:"required"`
		AbilityID   uint `json:"ability_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetUint("user_id")

	err := h.equipService.UnequipAbility(userID, req.CharacterID, req.AbilityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "ability unequipped successfully",
		"character_id": req.CharacterID,
		"ability_id":   req.AbilityID,
	})
}
