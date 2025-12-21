package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
	"gorm.io/gorm"
)

type MissionHandler struct {
	missionService *services.MissionService
}

func NewMissionHandler(db *gorm.DB) *MissionHandler {
	return &MissionHandler{
		missionService: services.NewMissionService(db),
	}
}

// ListMissions godoc
// @Summary Get available missions
// @Description Get all available missions for the authenticated user
// @Tags missions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/missions [get]
func (h *MissionHandler) ListMissions(c *gin.Context) {
	userID := c.GetUint("user_id")

	missions, err := h.missionService.GetAvailableMissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch missions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"missions": missions,
	})
}

// GetCurrentMission godoc
// @Summary Get current active mission
// @Description Get the user's current active mission
// @Tags missions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/missions/current [get]
func (h *MissionHandler) GetCurrentMission(c *gin.Context) {
	userID := c.GetUint("user_id")

	mission, err := h.missionService.GetCurrentMission(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch current mission",
		})
		return
	}

	if mission == nil {
		c.JSON(http.StatusOK, gin.H{
			"mission": nil,
			"message": "No active mission",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mission": mission,
	})
}

// StartMission godoc
// @Summary Start a mission
// @Description Start a specific mission by ID
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/missions/{id}/start [post]
func (h *MissionHandler) StartMission(c *gin.Context) {
	userID := c.GetUint("user_id")

	missionIDStr := c.Param("id")
	missionID, err := strconv.ParseUint(missionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid mission ID",
		})
		return
	}

	err = h.missionService.StartMission(userID, uint(missionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mission started successfully",
	})
}

// GetMissionProgress godoc
// @Summary Get mission progress
// @Description Get progress for a specific mission
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/missions/{id}/progress [get]
func (h *MissionHandler) GetMissionProgress(c *gin.Context) {
	userID := c.GetUint("user_id")

	missionIDStr := c.Param("id")
	missionID, err := strconv.ParseUint(missionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid mission ID",
		})
		return
	}

	// Get mission progress (placeholder - should query from DB)
	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"mission_id": missionID,
		"message":    "Progress fetched successfully",
	})
}

// CompleteMission godoc
// @Summary Complete a mission
// @Description Attempt to complete a mission and claim rewards
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/missions/{id}/complete [post]
func (h *MissionHandler) CompleteMission(c *gin.Context) {
	userID := c.GetUint("user_id")

	missionIDStr := c.Param("id")
	missionID, err := strconv.ParseUint(missionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid mission ID",
		})
		return
	}

	rewards, err := h.missionService.CompleteMission(userID, uint(missionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mission completed successfully",
		"rewards": rewards,
	})
}

// GetUnlockedFeatures godoc
// @Summary Get unlocked features
// @Description Get all features unlocked by the user
// @Tags missions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user/unlocks [get]
func (h *MissionHandler) GetUnlockedFeatures(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Check various unlocks
	breeding, _ := h.missionService.CheckFeatureUnlocked(userID, "breeding")
	crafting, _ := h.missionService.CheckFeatureUnlocked(userID, "crafting")
	raids, _ := h.missionService.CheckFeatureUnlocked(userID, "raids")
	ranked, _ := h.missionService.CheckFeatureUnlocked(userID, "ranked")
	advancedCrafting, _ := h.missionService.CheckFeatureUnlocked(userID, "advanced_crafting")

	c.JSON(http.StatusOK, gin.H{
		"unlocked_features": gin.H{
			"breeding":          breeding,
			"crafting":          crafting,
			"raids":             raids,
			"ranked":            ranked,
			"advanced_crafting": advancedCrafting,
		},
	})
}
