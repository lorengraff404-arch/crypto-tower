package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// CharacterHandler handles character HTTP requests
type CharacterHandler struct {
	characterService *services.CharacterService
}

// NewCharacterHandler creates a new character handler
func NewCharacterHandler() *CharacterHandler {
	return &CharacterHandler{
		characterService: services.NewCharacterService(),
	}
}

// ListCharacters returns all characters for authenticated user
// GET /api/v1/characters
func (h *CharacterHandler) ListCharacters(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	includeEggs := c.Query("include_eggs") == "true"

	characters, err := h.characterService.GetUserCharacters(userID.(uint), includeEggs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": characters,
		"count":      len(characters),
	})
}

// GetCharacter returns a specific character
// GET /api/v1/characters/:id
func (h *CharacterHandler) GetCharacter(c *gin.Context) {
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

	character, err := h.characterService.GetCharacterByID(uint(characterID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	// Return character directly (backward compatible)
	// Progression info available at character.progression if needed
	c.JSON(http.StatusOK, character)
}

// CreateCharacter creates a new character for testing
// POST /api/v1/characters
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		CharacterType string `json:"character_type" binding:"required"`
		Element       string `json:"element" binding:"required"`
		Rarity        string `json:"rarity" binding:"required"`
		Class         string `json:"class" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	character, err := h.characterService.CreateCharacter(
		userID.(uint),
		req.CharacterType,
		req.Element,
		req.Rarity,
		req.Class,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
		return
	}

	// SKILL SYSTEM INTEGRATION: Initialize skills for new character
	skillInitService := services.NewSkillInitializationService()
	if err := skillInitService.InitializeCharacterSkills(character.ID); err != nil {
		// Log error but don't fail character creation
		c.JSON(http.StatusCreated, gin.H{
			"character": character,
			"warning":   "Character created but skill initialization failed",
		})
		return
	}

	c.JSON(http.StatusCreated, character)
}

// HatchEgg hatches an egg character
// POST /api/v1/characters/:id/hatch
func (h *CharacterHandler) HatchEgg(c *gin.Context) {
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

	character, err := h.characterService.HatchEgg(uint(characterID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Character hatched successfully",
		"character": character,
	})
}

// RecoverFatigue recovers character fatigue
// POST /api/v1/characters/:id/recover
func (h *CharacterHandler) RecoverFatigue(c *gin.Context) {
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
		Amount int `json:"amount" binding:"required,min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Verify ownership first
	_, err = h.characterService.GetCharacterByID(uint(characterID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	if err := h.characterService.RecoverFatigue(uint(characterID), req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recover fatigue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fatigue recovered successfully",
	})
}
