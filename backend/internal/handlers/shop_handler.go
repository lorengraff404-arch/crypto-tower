package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// ShopHandler handles shop operations
type ShopHandler struct {
	shopService *services.ShopService
}

// NewShopHandler creates a new shop handler
func NewShopHandler(s *services.ShopService) *ShopHandler {
	return &ShopHandler{
		shopService: s,
	}
}

// GetShopItems returns all available shop items
// GET /api/v1/shop/items
// GET /api/v1/shop/items/:category
func (h *ShopHandler) GetShopItems(c *gin.Context) {
	category := c.Param("category")

	items, err := h.shopService.GetShopItems(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}

// BuyItem purchases an item from the shop
// POST /api/v1/shop/buy
func (h *ShopHandler) BuyItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ItemID   uint   `json:"item_id" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=1,max=99"`
		TxHash   string `json:"tx_hash" binding:"required"` // Blockchain transaction hash
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the blockchain transaction for audit
	// In a production system, you would verify the tx_hash on-chain here
	// For now, we just log it
	c.Set("tx_hash", req.TxHash)

	if err := h.shopService.BuyItem(userID, req.ItemID, req.Quantity, req.TxHash); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item purchased successfully",
		"tx_hash": req.TxHash, // Return the transaction hash
	})
}

// GetInventory returns user's shop inventory
// GET /api/v1/shop/inventory
func (h *ShopHandler) GetInventory(c *gin.Context) {
	userID := c.GetUint("user_id")

	inventory, err := h.shopService.GetUserInventory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"inventory": inventory,
		"count":     len(inventory),
	})
}

// UseItem uses a consumable item
// POST /api/v1/shop/use
func (h *ShopHandler) UseItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ItemID      uint `json:"item_id" binding:"required"`
		CharacterID uint `json:"character_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.shopService.UseItem(userID, req.ItemID, req.CharacterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item used successfully",
	})
}
