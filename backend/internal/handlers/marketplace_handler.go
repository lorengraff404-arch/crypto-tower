package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

type MarketplaceHandler struct {
	marketplaceService *services.MarketplaceService
}

func NewMarketplaceHandler() *MarketplaceHandler {
	return &MarketplaceHandler{
		marketplaceService: &services.MarketplaceService{},
	}
}

// CreateListing creates a new marketplace listing
func (h *MarketplaceHandler) CreateListing(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ItemType string `json:"item_type" binding:"required"`
		ItemID   uint   `json:"item_id" binding:"required"`
		Price    int    `json:"price" binding:"required"`
		Currency string `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Currency == "" {
		req.Currency = "TOWER"
	}

	listing, err := h.marketplaceService.CreateListing(userID, req.ItemType, req.ItemID, req.Price, req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"listing": listing,
	})
}

// GetListings returns active marketplace listings
func (h *MarketplaceHandler) GetListings(c *gin.Context) {
	itemType := c.Query("type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	listings, err := h.marketplaceService.GetActiveListings(itemType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"listings": listings,
	})
}

// BuyListing purchases an item
func (h *MarketplaceHandler) BuyListing(c *gin.Context) {
	userID := c.GetUint("user_id")
	listingID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	err := h.marketplaceService.BuyListing(userID, uint(listingID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// CancelListing cancels a listing
func (h *MarketplaceHandler) CancelListing(c *gin.Context) {
	userID := c.GetUint("user_id")
	listingID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	err := h.marketplaceService.CancelListing(userID, uint(listingID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
