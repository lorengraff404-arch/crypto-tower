package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	adminID := c.GetUint("userID")
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

// UpdateSetting updates a system setting
// PUT /api/v1/admin/settings
func (h *AdminHandler) UpdateSetting(c *gin.Context) {
	var req struct {
		Key         string `json:"key" binding:"required"`
		Value       string `json:"value" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("userID")
	if err := h.configService.SetConfig(req.Key, req.Value, req.Description, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Setting updated"})
}

// GetRevenueStats returns financial metrics
func (h *AdminHandler) GetRevenueStats(c *gin.Context) {
	stats, err := h.adminService.GetRevenueStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
