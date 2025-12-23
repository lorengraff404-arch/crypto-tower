package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetActiveBattles returns all currently active battles
// GET /api/v1/admin/battles/active
func (h *AdminHandler) GetActiveBattles(c *gin.Context) {
	battles, err := h.adminService.GetActiveBattles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"battles": battles})
}

// GetBattleHistory returns global battle history
// GET /api/v1/admin/battles/history
func (h *AdminHandler) GetBattleHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	battles, err := h.adminService.GetBattleHistory(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"battles": battles})
}

// TerminateBattle forces a battle to end
// POST /api/v1/admin/battles/:id/terminate
func (h *AdminHandler) TerminateBattle(c *gin.Context) {
	battleIDStr := c.Param("id")
	battleID, err := strconv.ParseUint(battleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid battle ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reason is required"})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.TerminateBattle(uint(battleID), req.Reason, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Battle terminated successfully"})
}
