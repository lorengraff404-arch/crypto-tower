package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

type BattleHandler struct {
	battleService *services.BattleService
}

func NewBattleHandler() *BattleHandler {
	return &BattleHandler{
		battleService: services.NewBattleService(),
	}
}

// ProcessTurn
// POST /api/v1/battles/:id/turn
func (h *BattleHandler) ProcessTurn(c *gin.Context) {
	battleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := c.GetUint("userID")

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	battle, err := h.battleService.ProcessTurn(uint(battleID), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, battle)
}

// CompleteBattle
// POST /api/v1/battles/:id/complete
func (h *BattleHandler) CompleteBattle(c *gin.Context) {
	battleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		WinnerID uint `json:"winner_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Security: In production, verify signature from game server or authorized client
	// For now, we trust the client (DEV MODE ONLY WARNING)

	if err := h.battleService.CompleteBattle(uint(battleID), req.WinnerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Battle completed"})
}
