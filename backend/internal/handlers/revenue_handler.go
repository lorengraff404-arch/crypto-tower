package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// RevenueHandler handles revenue-related requests
type RevenueHandler struct {
	revenueService *services.RevenueService
}

// NewRevenueHandler creates a new revenue handler
func NewRevenueHandler(revenueService *services.RevenueService) *RevenueHandler {
	return &RevenueHandler{
		revenueService: revenueService,
	}
}

// GetRevenueStats returns revenue distribution statistics
// GET /api/v1/revenue/stats
func (h *RevenueHandler) GetRevenueStats(c *gin.Context) {
	log.Println("Incoming request to GET /api/v1/revenue/stats")
	if h.revenueService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Revenue service not initialized"})
		return
	}

	stats, err := h.revenueService.GetRevenueStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
