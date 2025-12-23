package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

type LeaderboardHandler struct{}

func NewLeaderboardHandler() *LeaderboardHandler {
	return &LeaderboardHandler{}
}

// GetLeaderboard returns the top N players by ELO
// GET /api/v1/leaderboard
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	var users []models.User
	// Fetch ID, Username (Wallet?), ELO, Rank
	// Note: Username is not in User struct shown earlier?
	// Checking User model: WalletAddress is main ID. No Username field explicitly visible in `view_file` output earlier?
	// Let's check `User` model output again.
	// Line 24: WalletAddress string
	// Line 47: Characters
	// No "Username". We will use WalletAddress (truncated).

	if err := db.DB.Select("id, wallet_address, elo_rating, rank_tier, pvp_wins, pvp_losses").
		Order("elo_rating desc").
		Limit(limit).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	// Transform for Safe Output (Hide full wallet if sensitive? Public is ok usually)
	var response []gin.H
	for _, u := range users {
		response = append(response, gin.H{
			"id":        u.ID,
			"wallet":    u.WalletAddress, // Frontend can truncate
			"elo":       u.ELO,
			"rank_tier": u.RankTier,
			"wins":      u.PvPWins,
			"losses":    u.PvPLosses,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": response,
	})
}
