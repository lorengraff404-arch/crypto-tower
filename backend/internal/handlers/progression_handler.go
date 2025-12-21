package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// ProgressionHandler handles character progression HTTP requests
type ProgressionHandler struct {
	progressionService *services.ProgressionService
}

// NewProgressionHandler creates a new progression handler
func NewProgressionHandler() *ProgressionHandler {
	return &ProgressionHandler{
		progressionService: services.NewProgressionService(),
	}
}

// GainXP grants experience to a character
// POST /api/v1/characters/:id/gain-xp
func (h *ProgressionHandler) GainXP(c *gin.Context) {
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

	var req struct {
		XPGained   int    `json:"xp_gained" binding:"required,min=1"`
		Source     string `json:"source" binding:"required"`
		Difficulty string `json:"difficulty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Verify ownership
	characterService := services.NewCharacterService()
	char, err := characterService.GetCharacterByID(uint(characterID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	// Gain experience
	updatedChar, err := h.progressionService.GainExperience(char.ID, req.XPGained, req.Source, req.Difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Experience gained",
		"character": updatedChar,
		"leveled_up": updatedChar.Level > char.Level,
	})
}

// GetProgressionInfo returns character progression details
// GET /api/v1/characters/:id/progression
func (h *ProgressionHandler) GetProgressionInfo(c *gin.Context) {
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
	_, err = characterService.GetCharacterByID(uint(characterID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	info, err := h.progressionService.GetProgressionInfo(uint(characterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progression info"})
		return
	}

	c.JSON(http.StatusOK, info)
}

// ValidateIntegrity checks character integrity (anti-cheat)
// GET /api/v1/characters/:id/validate
func (h *ProgressionHandler) ValidateIntegrity(c *gin.Context) {
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
	_, err = characterService.GetCharacterByID(uint(characterID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	err = h.progressionService.ValidateCharacterIntegrity(uint(characterID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"message": "Character integrity validated",
	})
}
