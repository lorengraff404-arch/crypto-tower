package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// TeamHandler handles team HTTP requests
type TeamHandler struct {
	teamService *services.TeamService
}

// NewTeamHandler creates a new handler
func NewTeamHandler() *TeamHandler {
	return &TeamHandler{
		teamService: services.NewTeamService(),
	}
}

// GetMyTeams returns all teams for the current user
// GET /api/v1/teams
func (h *TeamHandler) GetMyTeams(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	teams, err := h.teamService.GetUserTeams(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

// CreateTeam creates a new team
// POST /api/v1/teams
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.teamService.CreateTeam(userID.(uint), req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// UpdateTeamMember adds/updates a member in the team
// POST /api/v1/teams/:id/members
func (h *TeamHandler) UpdateTeamMember(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req struct {
		CharacterID uint `json:"character_id" binding:"required"`
		Slot        int  `json:"slot"` // 0-2 (checked in service)
		IsBackup    bool `json:"is_backup"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.teamService.AddMember(uint(teamID), req.CharacterID, req.Slot, req.IsBackup)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member updated successfully"})
}

// RemoveTeamMember removes a member
// DELETE /api/v1/teams/:id/members/:charId
func (h *TeamHandler) RemoveTeamMember(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	characterID, err := strconv.ParseUint(c.Param("charId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	err = h.teamService.RemoveMember(uint(teamID), uint(characterID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// GetActiveTeam returns the user's active team (currently the first one)
// GET /api/v1/teams/active
func (h *TeamHandler) GetActiveTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	teams, err := h.teamService.GetUserTeams(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	if len(teams) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active team found"})
		return
	}

	// For now, just return the first team as the active one
	c.JSON(http.StatusOK, teams[0])
}
