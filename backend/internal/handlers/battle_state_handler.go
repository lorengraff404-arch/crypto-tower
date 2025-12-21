package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// GetRaidBattleState returns current battle state for real-time UI updates
// GET /api/v1/raids/:sessionId/state
func GetRaidBattleState(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID := c.Param("sessionId")

	// SECURITY: Verify ownership
	var session models.RaidSession
	if err := db.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Raid not found"})
		return
	}

	// Parse current battle state
	// CharacterStates contains HP/Mana/Status of all participants
	c.JSON(http.StatusOK, gin.H{
		"session_id":       session.ID,
		"status":           session.Status,
		"character_states": session.CharacterStates,
		"is_victory":       session.Status == "VICTORY",
		"is_defeat":        session.Status == "DEFEAT",
		"can_complete":     session.Status == "VICTORY",
	})
}

// GetWagerPreview shows potential rewards/risks before battle
// GET /api/v1/battle/wager-preview
func GetWagerPreview(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Load user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// SECURITY: Check minimum balance
	minBet := int64(100)
	maxBet := int64(100000)

	if user.GTKBalance < minBet {
		c.JSON(http.StatusOK, gin.H{
			"can_play":     false,
			"reason":       "Insufficient balance",
			"your_balance": user.GTKBalance,
			"minimum_bet":  minBet,
		})
		return
	}

	// Calculate potential rewards based on ELO
	c.JSON(http.StatusOK, gin.H{
		"can_play":     true,
		"your_balance": user.GTKBalance,
		"your_elo":     user.ELO,
		"minimum_bet":  minBet,
		"maximum_bet":  maxBet,
		"platform_fee": "5%",
		"risk_info":    "Bet amount escrowed. Winner gets (bet - 5% fee). Dynamic payout based on opponent difficulty.",
	})
}
