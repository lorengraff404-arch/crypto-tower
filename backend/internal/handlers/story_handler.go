package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// StoryHandler handles narrative content HTTP requests
type StoryHandler struct{}

// NewStoryHandler creates a new story handler
func NewStoryHandler() *StoryHandler {
	return &StoryHandler{}
}

// GetMissionDialogues godoc
// @Summary Get dialogues for a mission
// @Description Get all story dialogues for a specific mission level
// @Tags story
// @Accept json
// @Produce json
// @Param level path int true "Mission Level"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/story/missions/{level}/dialogues [get]
func (h *StoryHandler) GetMissionDialogues(c *gin.Context) {
	levelStr := c.Param("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid mission level",
		})
		return
	}

	var dialogues []models.StoryDialogue
	err = db.DB.Where("mission_level = ?", level).Order("sort_order ASC").Find(&dialogues).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch dialogues",
		})
		return
	}

	// Group by dialogue type
	briefings := []models.StoryDialogue{}
	postMissions := []models.StoryDialogue{}
	cutscenes := []models.StoryDialogue{}

	for _, dialogue := range dialogues {
		switch dialogue.DialogueType {
		case "briefing":
			briefings = append(briefings, dialogue)
		case "post_mission":
			postMissions = append(postMissions, dialogue)
		case "cutscene":
			cutscenes = append(cutscenes, dialogue)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"mission_level": level,
		"briefings":     briefings,
		"post_missions": postMissions,
		"cutscenes":     cutscenes,
	})
}

// GetAvailableFragments godoc
// @Summary Get available story fragments
// @Description Get story fragments unlockable at user's level
// @Tags story
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/story/fragments [get]
func (h *StoryHandler) GetAvailableFragments(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Get user level (assuming it's in context or we fetch it)
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found",
		})
		return
	}

	// Get fragments user can unlock
	var fragments []models.StoryFragment
	err := db.DB.Where("unlock_level <= ?", user.Level).Order("unlock_level ASC").Find(&fragments).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch fragments",
		})
		return
	}

	// Get user's collected fragments
	var userProgress models.UserStoryProgress
	db.DB.Where("user_id = ?", userID).First(&userProgress)

	c.JSON(http.StatusOK, gin.H{
		"available_fragments": fragments,
		"user_progress":       userProgress,
	})
}

// GetStoryProgress godoc
// @Summary Get user's story progression
// @Description Get current story progress including act, corruption, relationships
// @Tags story
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/story/progress [get]
func (h *StoryHandler) GetStoryProgress(c *gin.Context) {
	userID := c.GetUint("user_id")

	var progress models.UserStoryProgress
	err := db.DB.Where("user_id = ?", userID).First(&progress).Error

	if err != nil {
		// Create initial progress if doesn't exist
		progress = models.UserStoryProgress{
			UserID:             userID,
			CurrentAct:         1,
			AriaCorruption:     0,
			KairosRelationship: "neutral",
			VoiceEncounters:    0,
			CollectedFragments: "[]",
			ViewedCutscenes:    "[]",
		}
		db.DB.Create(&progress)
	}

	c.JSON(http.StatusOK, gin.H{
		"progress": progress,
	})
}

// RecordChoice godoc
// @Summary Record a player narrative choice
// @Description Save player's choice for branching narrative
// @Tags story
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/story/choices [post]
func (h *StoryHandler) RecordChoice(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ChoiceID    string `json:"choice_id" binding:"required"`
		ChoiceValue string `json:"choice_value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	choice := models.PlayerChoice{
		UserID:      userID,
		ChoiceID:    req.ChoiceID,
		ChoiceValue: req.ChoiceValue,
	}

	if err := db.DB.Create(&choice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to record choice",
		})
		return
	}

	// Update story progress based on choice
	var progress models.UserStoryProgress
	db.DB.Where("user_id = ?", userID).First(&progress)

	// Update relationships based on choice
	switch req.ChoiceID {
	case "kairos_befriend":
		progress.KairosRelationship = "friend"
	case "kairos_antagonize":
		progress.KairosRelationship = "enemy"
	}

	db.DB.Save(&progress)

	c.JSON(http.StatusOK, gin.H{
		"choice":  choice,
		"message": "Choice recorded successfully",
	})
}

// UpdateStoryProgress godoc
// @Summary Update story progression
// @Description Update Aria corruption, act number, etc.
// @Tags story
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/story/progress [put]
func (h *StoryHandler) UpdateStoryProgress(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		CurrentAct         *int    `json:"current_act"`
		AriaCorruption     *int    `json:"aria_corruption"`
		VoiceEncounters    *int    `json:"voice_encounters"`
		CollectedFragments *string `json:"collected_fragments"`
		ViewedCutscenes    *string `json:"viewed_cutscenes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	var progress models.UserStoryProgress
	result := db.DB.Where("user_id = ?", userID).First(&progress)

	if result.Error != nil {
		// Create if doesn't exist
		progress = models.UserStoryProgress{
			UserID: userID,
		}
	}

	// Update only provided fields
	if req.CurrentAct != nil {
		progress.CurrentAct = *req.CurrentAct
	}
	if req.AriaCorruption != nil {
		progress.AriaCorruption = *req.AriaCorruption
	}
	if req.VoiceEncounters != nil {
		progress.VoiceEncounters = *req.VoiceEncounters
	}
	if req.CollectedFragments != nil {
		progress.CollectedFragments = *req.CollectedFragments
	}
	if req.ViewedCutscenes != nil {
		progress.ViewedCutscenes = *req.ViewedCutscenes
	}

	if result.Error != nil {
		db.DB.Create(&progress)
	} else {
		db.DB.Save(&progress)
	}

	c.JSON(http.StatusOK, gin.H{
		"progress": progress,
		"message":  "Story progress updated",
	})
}
