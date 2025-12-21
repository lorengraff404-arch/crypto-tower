package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// RaidTurnHandler struct to hold service
type RaidTurnHandler struct {
	raidService *services.RaidService
}

func NewRaidTurnHandler(service *services.RaidService) *RaidTurnHandler {
	return &RaidTurnHandler{raidService: service}
}

// ProcessRaidTurn handles turn execution in raid battles
func (h *RaidTurnHandler) ProcessRaidTurn(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid Session ID required"})
		return
	}

	var req struct {
		CharacterID uint   `json:"character_id"`
		ActionType  string `json:"action_type"`
		MoveSlot    int    `json:"move_slot"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var session *models.RaidSession
	var result *services.BattleResult

	if req.ActionType == "enemy" {
		session, result, err = h.raidService.ExecuteEnemyTurn(uint(sessionID))
	} else {
		session, result, err = h.raidService.ExecutePlayerTurn(uint(sessionID), req.CharacterID, req.MoveSlot)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session": session,
		"result":  result,
	})
}

// CompleteRaid checks status
func (h *RaidTurnHandler) CompleteRaid(c *gin.Context) {
	// Status check only
	c.JSON(http.StatusOK, gin.H{"message": "Raid completion is handled automatically at end of turn."})
}

// FleeRaid abandons raid
func (h *RaidTurnHandler) FleeRaid(c *gin.Context) {
	userID := c.GetUint("user_id")
	sessionIDStr := c.Param("sessionId")
	sessionID, _ := strconv.ParseUint(sessionIDStr, 10, 64)

	if err := h.raidService.AbandonSession(uint(sessionID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
