package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

type BreedingHandler struct {
	breedingService *services.BreedingService
}

func NewBreedingHandler(service *services.BreedingService) *BreedingHandler {
	return &BreedingHandler{
		breedingService: service,
	}
}

// StartBreeding initiates breeding between two characters
func (h *BreedingHandler) StartBreeding(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Parent1ID uint   `json:"parent1_id" binding:"required"`
		Parent2ID uint   `json:"parent2_id" binding:"required"`
		TxHash    string `json:"tx_hash"` // Optional blockchain payment
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.breedingService.StartBreeding(userID, req.Parent1ID, req.Parent2ID, req.TxHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"session": session,
	})
}

// GetUserEggs returns user's eggs
func (h *BreedingHandler) GetUserEggs(c *gin.Context) {
	userID := c.GetUint("user_id")

	eggs, err := h.breedingService.GetUserEggs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"eggs":  eggs,
		"count": len(eggs),
	})
}

// StartIncubation starts egg incubation
func (h *BreedingHandler) StartIncubation(c *gin.Context) {
	userID := c.GetUint("user_id")
	eggIDStr := c.Param("id")
	eggID, _ := strconv.ParseUint(eggIDStr, 10, 32)

	if err := h.breedingService.IncubateEgg(userID, uint(eggID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// HatchEgg hatches an egg
func (h *BreedingHandler) HatchEgg(c *gin.Context) {
	userID := c.GetUint("user_id")
	eggID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	character, err := h.breedingService.HatchEgg(userID, uint(eggID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"character": character,
	})
}
