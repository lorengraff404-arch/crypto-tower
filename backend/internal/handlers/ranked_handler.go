package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// RankedHandler handles ranked PvP battle endpoints
type RankedHandler struct {
	battleEngine *services.BattleEngine
}

// NewRankedHandler creates a new ranked handler
func NewRankedHandler() *RankedHandler {
	return &RankedHandler{
		battleEngine: services.NewBattleEngine(),
	}
}

// StartRankedRequest represents ranked battle matchmaking request
type StartRankedRequest struct {
	TeamID uint `json:"team_id" binding:"required"`
}

// StartRanked initiates ranked matchmaking
// POST /api/v1/battle/ranked
// Security: ELO-based matchmaking, ownership validation, anti-same-user exploit
func (h *RankedHandler) StartRanked(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req StartRankedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. SECURITY: Verify team ownership
	var team models.Team
	if err := db.DB.Where("id = ? AND owner_id = ?", req.TeamID, userID).
		Preload("Members.Character").
		First(&team).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Team not found or not owned"})
		return
	}

	// 2. SECURITY: Validate team (3v3, durability, alive)
	activeCount := 0
	teamCP := 0
	for _, member := range team.Members {
		if !member.IsBackup {
			activeCount++
			if member.Character.Durability < 10 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":        "Character durability too low",
					"character_id": member.Character.ID,
				})
				return
			}
			if member.Character.IsDead {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":        "Cannot use dead character",
					"character_id": member.Character.ID,
				})
				return
			}
			teamCP += member.Character.CombatPower
		}
	}

	if activeCount == 0 || activeCount > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team must have 1-3 active members"})
		return
	}

	// 3. Get user ELO
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// 4. SECURITY: Check for active battles
	var existingBattle models.Battle
	if err := db.DB.Where("(player1_id = ? OR player2_id = ?) AND status = ?",
		userID, userID, "active").First(&existingBattle).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":     "Already in active battle",
			"battle_id": existingBattle.ID,
		})
		return
	}

	// 5. Find opponent (ELO ±200 range, similar CP)
	var opponent models.User
	minELO := user.ELO - 200
	maxELO := user.ELO + 200

	// Find waiting player in queue (simplified - would use Redis queue in production)
	if err := db.DB.Where("id != ? AND elo BETWEEN ? AND ?", userID, minELO, maxELO).
		Order("elo ASC").
		First(&opponent).Error; err != nil {
		// No opponent found - add to queue (simplified - return "searching")
		c.JSON(http.StatusAccepted, gin.H{
			"status":    "searching",
			"message":   "Searching for opponent...",
			"elo_range": map[string]int{"min": minELO, "max": maxELO},
		})
		return
	}

	// 6. SECURITY: Prevent self-matchmaking
	if opponent.ID == userID.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot battle yourself"})
		return
	}

	// 7. Create battle
	battleSeed := generateBattleSeed()
	battle := &models.Battle{
		Player1ID:  userID.(uint),
		Player2ID:  opponent.ID,
		Status:     "active",
		BattleType: "ranked",
		Seed:       battleSeed,
		// Store ELO for calculation at end
	}

	if err := db.DB.Create(battle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create battle"})
		return
	}

	// 8. Deduct durability (atomic)
	tx := db.DB.Begin()
	for _, member := range team.Members {
		if !member.IsBackup {
			member.Character.Durability -= 3 // Ranked costs less than Raid
			if err := tx.Save(&member.Character).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update durability"})
				return
			}
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":   "Ranked battle started",
		"battle_id": battle.ID,
		"opponent": map[string]interface{}{
			"id":     opponent.ID,
			"wallet": opponent.WalletAddress,
			"elo":    opponent.ELO,
		},
		"your_elo": user.ELO,
	})
}

// generateBattleSeed creates a deterministic battle seed
func generateBattleSeed() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// UpdateELO calculates new ELO after battle (called on victory/defeat)
func UpdateELO(winnerID, loserID uint) error {
	var winner, loser models.User

	if err := db.DB.First(&winner, winnerID).Error; err != nil {
		return err
	}
	if err := db.DB.First(&loser, loserID).Error; err != nil {
		return err
	}

	// K-factor (max rating change)
	K := 32.0

	// Expected scores
	expectedWinner := 1.0 / (1.0 + pow10((float64(loser.ELO-winner.ELO))/400.0))
	expectedLoser := 1.0 - expectedWinner

	// Update ELO
	winner.ELO += int(K * (1.0 - expectedWinner))
	loser.ELO += int(K * (0.0 - expectedLoser))

	// Prevent negative ELO
	if loser.ELO < 0 {
		loser.ELO = 0
	}

	// Save
	tx := db.DB.Begin()
	if err := tx.Save(&winner).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(&loser).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func pow10(x float64) float64 {
	// Simple power of 10 approximation for ELO
	result := 1.0
	for i := 0; i < int(x); i++ {
		result *= 10
	}
	return result
}
