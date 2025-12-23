package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// RaidHandler handles raid battle endpoints
type RaidHandler struct {
	raidService  *services.RaidService
	battleEngine *services.BattleEngine
}

// NewRaidHandler creates a new raid handler with dependencies injected
func NewRaidHandler(raidService *services.RaidService) *RaidHandler {
	return &RaidHandler{
		raidService:  raidService,
		battleEngine: services.NewBattleEngine(),
	}
}

// StartRaidRequest represents the request to start a raid
type StartRaidRequest struct {
	TeamID     uint   `json:"team_id" binding:"required"`
	MissionID  uint   `json:"mission_id" binding:"required"`
	Difficulty string `json:"difficulty"` // easy, normal, hard (optional)
}

// StartRaid initiates a new raid battle
// POST /api/v1/raids/start
// Security: Validates user ownership, team validity, character durability
func (h *RaidHandler) StartRaid(c *gin.Context) {
	// 1. SECURITY: Get authenticated user
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req StartRaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// 2. SECURITY: Verify team ownership
	var team models.Team
	if err := db.DB.Where("id = ? AND owner_id = ?", req.TeamID, userID).
		Preload("Members.Character").
		First(&team).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Team not found or not owned by user"})
		return
	}

	// 3. SECURITY: Validate team composition (max 3 active members for 3v3)
	activeCount := 0
	for _, member := range team.Members {
		if !member.IsBackup {
			activeCount++
		}
	}
	if activeCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team must have at least 1 active member"})
		return
	}
	if activeCount > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team cannot have more than 3 active members (3v3 limit)"})
		return
	}

	// 4. SECURITY: Check character durability (anti-exploit)
	for _, member := range team.Members {
		if !member.IsBackup && member.Character.Durability < 10 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":        "Character has too low durability",
				"character_id": member.Character.ID,
				"durability":   member.Character.Durability,
				"message":      "Characters need at least 10 durability to battle. Use revival items or rest.",
			})
			return
		}
		if member.Character.IsDead {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":        "Cannot use dead character",
				"character_id": member.Character.ID,
				"message":      "Revive character first using Sacred Ash or Revival Herb",
			})
			return
		}
	}

	// 5. SECURITY: Verify mission exists and is not locked
	var mission models.IslandMission
	if err := db.DB.First(&mission, req.MissionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	// 6. SECURITY: Rate limiting - Check for active raid session
	var existingSession models.RaidSession
	if err := db.DB.Where("user_id = ? AND status = ?", userID, "active").
		First(&existingSession).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":      "You already have an active raid session",
			"session_id": existingSession.ID,
			"message":    "Complete or abandon your current raid first",
		})
		return
	}

	// 7. Create raid session
	raidSession := &models.RaidSession{
		UserID:    userID.(uint),
		TeamID:    req.TeamID,
		MissionID: req.MissionID,
		Status:    "active",
	}

	if err := db.DB.Create(raidSession).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create raid session"})
		return
	}

	// 8. Deduct durability immediately (anti-cheat)
	tx := db.DB.Begin()
	for _, member := range team.Members {
		if !member.IsBackup {
			member.Character.Durability -= 5
			if err := tx.Save(&member.Character).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update durability"})
				return
			}
		}
	}
	tx.Commit()

	// 9. Return success
	c.JSON(http.StatusOK, gin.H{
		"message":        "Raid started successfully",
		"session_id":     raidSession.ID,
		"mission_id":     mission.ID,
		"team_id":        team.ID,
		"active_members": activeCount,
		"next_action":    "Use /raid/turn endpoint to execute turns",
	})
}

// GetRaidStatus returns current raid session status
// GET /api/v1/raids/:sessionId/status
func (h *RaidHandler) GetRaidStatus(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID, _ := c.Get("userID")

	var session models.RaidSession
	if err := db.DB.Where("id = ? AND user_id = ?", sessionID, userID).
		First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Raid session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.ID,
		"status":     session.Status,
	})
}
