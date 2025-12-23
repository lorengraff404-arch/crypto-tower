package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

type AdminHandler struct {
	adminService  *services.AdminService
	configService *services.ConfigService
}

func NewAdminHandler(adminService *services.AdminService, configService *services.ConfigService) *AdminHandler {
	return &AdminHandler{
		adminService:  adminService,
		configService: configService,
	}
}

// BanUser Endpoint
func (h *AdminHandler) BanUser(c *gin.Context) {
	var req struct {
		TargetID uint   `json:"target_id" binding:"required"`
		Reason   string `json:"reason" binding:"required"`
		Duration int    `json:"duration" binding:"required"` // Hours
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.BanUser(req.TargetID, req.Reason, req.Duration, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
}

// GetSettings returns all system settings
// GET /api/v1/admin/settings
func (h *AdminHandler) GetSettings(c *gin.Context) {
	settings, err := h.configService.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// UpdateSetting updates a single system setting
// PUT /api/v1/admin/settings
func (h *AdminHandler) UpdateSetting(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.UpdateSystemSetting(req.Key, req.Value, "", adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Setting updated successfully",
		"key":     req.Key,
		"value":   req.Value,
	})
}

// GetRevenueStats returns aggregated revenue statistics
func (h *AdminHandler) GetRevenueStats(c *gin.Context) {
	stats, err := h.adminService.GetRevenueStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// UnbanUser endpoint
func (h *AdminHandler) UnbanUser(c *gin.Context) {
	var req struct {
		TargetID uint `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.UnbanUser(req.TargetID, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unbanned successfully"})
}

// FreezeFunds endpoint
func (h *AdminHandler) FreezeFunds(c *gin.Context) {
	var req struct {
		TargetID uint `json:"target_id" binding:"required"`
		Freeze   bool `json:"freeze" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.FreezeFunds(req.TargetID, req.Freeze, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	action := "frozen"
	if !req.Freeze {
		action = "unfrozen"
	}
	c.JSON(http.StatusOK, gin.H{"message": "Funds " + action + " successfully"})
}

// ListUsers returns paginated user list
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.adminService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetAuditLogs returns recent admin actions
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	logs, err := h.adminService.GetAuditLogs(50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// CreateShopItem creates a new shop item
func (h *AdminHandler) CreateShopItem(c *gin.Context) {
	var item models.ShopItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.CreateShopItem(item, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateShopItem updates an existing shop item
func (h *AdminHandler) UpdateShopItem(c *gin.Context) {
	var item models.ShopItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.UpdateShopItem(item, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteShopItem deletes a shop item
func (h *AdminHandler) DeleteShopItem(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.DeleteShopItem(req.ID, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shop item deleted successfully"})
}

// GetAdminShopItems returns all shop items for admin
func (h *AdminHandler) GetAdminShopItems(c *gin.Context) {
	items, err := h.adminService.ListShopItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateQuestTemplate creates a new quest template
func (h *AdminHandler) CreateQuestTemplate(c *gin.Context) {
	var template models.QuestTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.CreateQuestTemplate(template, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// UpdateQuestTemplate updates an existing quest template
func (h *AdminHandler) UpdateQuestTemplate(c *gin.Context) {
	var template models.QuestTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.UpdateQuestTemplate(template, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteQuestTemplate deletes a quest template
func (h *AdminHandler) DeleteQuestTemplate(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.adminService.DeleteQuestTemplate(req.ID, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quest template deleted successfully"})
}

// GetAdminQuestTemplates returns all quest templates for admin
func (h *AdminHandler) GetAdminQuestTemplates(c *gin.Context) {
	templates, err := h.adminService.ListQuestTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// ==================== ABILITIES MANAGEMENT ====================

// GetAbilities returns all abilities for admin management
func (h *AdminHandler) GetAbilities(c *gin.Context) {
	abilities, err := h.adminService.GetAllAbilities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, abilities)
}

// CreateAbility creates a new ability
func (h *AdminHandler) CreateAbility(c *gin.Context) {
	var ability models.Ability
	if err := c.ShouldBindJSON(&ability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminService.CreateAbility(&ability); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ability)
}

// UpdateAbility updates an existing ability
func (h *AdminHandler) UpdateAbility(c *gin.Context) {
	var ability models.Ability
	if err := c.ShouldBindJSON(&ability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminService.UpdateAbility(&ability); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ability)
}

// DeleteAbility deletes an ability
func (h *AdminHandler) DeleteAbility(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminService.DeleteAbility(req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ability deleted successfully"})
}
