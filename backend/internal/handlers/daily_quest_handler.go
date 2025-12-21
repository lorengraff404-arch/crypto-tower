package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

type DailyQuestHandler struct {
	questService *services.DailyQuestService
}

func NewDailyQuestHandler() *DailyQuestHandler {
	return &DailyQuestHandler{
		questService: services.NewDailyQuestService(),
	}
}

// GetDailyQuests returns user's active daily quests
// GET /api/v1/daily-quests
func (h *DailyQuestHandler) GetDailyQuests(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Get user to check level
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Generate quests if user doesn't have any
	h.questService.GenerateDailyQuests(userID, user.Level)

	// Get active quests
	quests, err := h.questService.GetActiveQuests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quests": quests,
		"count":  len(quests),
	})
}

// ClaimQuestReward claims a completed quest reward
// POST /api/v1/daily-quests/claim/:id
func (h *DailyQuestHandler) ClaimQuestReward(c *gin.Context) {
	userID := c.GetUint("user_id")
	questID := c.Param("id")

	var questIDUint uint
	if _, err := fmt.Sscanf(questID, "%d", &questIDUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	// Claim reward
	quest, err := h.questService.ClaimReward(userID, questIDUint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Distribute rewards
	var user models.User
	db.DB.First(&user, userID)

	// Add GTK (convert to int64)
	if quest.RewardGTK > 0 {
		user.GTKBalance += int64(quest.RewardGTK)
	}

	// Add TOWER (convert to int64)
	if quest.RewardTOWER > 0 {
		user.TOWERBalance += int64(quest.RewardTOWER * 1000000) // Convert to smallest unit
	}

	db.DB.Save(&user)

	// Add item to inventory if applicable
	var itemName string
	if quest.RewardItemID != nil {
		var item models.ShopItem
		if err := db.DB.First(&item, *quest.RewardItemID).Error; err == nil {
			// Add to user inventory
			inventory := models.UserInventory{
				UserID:   userID,
				ItemID:   uint(*quest.RewardItemID),
				Quantity: 1,
			}
			db.DB.Create(&inventory)
			itemName = item.Name
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Quest reward claimed!",
		"rewards": gin.H{
			"gtk":   quest.RewardGTK,
			"tower": quest.RewardTOWER,
			"item":  itemName,
		},
	})
}

// RefreshQuests manually refreshes quests (admin only)
// POST /api/v1/daily-quests/refresh
func (h *DailyQuestHandler) RefreshQuests(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Delete all current quests
	db.DB.Where("user_id = ?", userID).Delete(&models.DailyQuest{})

	// Generate new quests
	if err := h.questService.GenerateDailyQuests(userID, user.Level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Quests refreshed",
	})
}
