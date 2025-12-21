package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// ItemHandler handles item HTTP requests
type ItemHandler struct {
	itemService *services.ItemService
}

// NewItemHandler creates a new item handler
func NewItemHandler() *ItemHandler {
	return &ItemHandler{
		itemService: services.NewItemService(),
	}
}

// ListItems returns all items for authenticated user
// GET /api/v1/items?type=WEAPON
func (h *ItemHandler) ListItems(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemType := c.Query("type") // Optional filter

	items, err := h.itemService.GetUserItems(userID.(uint), itemType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}

// GetItem returns a specific item
// GET /api/v1/items/:id
func (h *ItemHandler) GetItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.itemService.GetItemByID(uint(itemID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateItem creates a new item for testing
// POST /api/v1/items
func (h *ItemHandler) CreateItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		ItemType           string `json:"item_type" binding:"required"`
		Name               string `json:"name" binding:"required"`
		Rarity             string `json:"rarity" binding:"required"`
		IsConsumable       bool   `json:"is_consumable"`
		IsCraftingMaterial bool   `json:"is_crafting_material"`
		IsStackable        bool   `json:"is_stackable"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	item, err := h.itemService.CreateItem(
		userID.(uint),
		req.ItemType,
		req.Name,
		req.Rarity,
		req.IsConsumable,
		req.IsCraftingMaterial,
		req.IsStackable,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// EquipItem equips an item to a character
// POST /api/v1/items/:id/equip
func (h *ItemHandler) EquipItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		CharacterID uint `json:"character_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.itemService.EquipItem(uint(itemID), req.CharacterID, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item equipped successfully",
	})
}

// UnequipItem unequips an item
// POST /api/v1/items/:id/unequip
func (h *ItemHandler) UnequipItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := h.itemService.UnequipItem(uint(itemID), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item unequipped successfully",
	})
}

// UseItem uses a consumable item
// POST /api/v1/items/:id/use
func (h *ItemHandler) UseItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		TargetCharacterID uint `json:"target_character_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.itemService.UseConsumable(uint(itemID), req.TargetCharacterID, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item used successfully",
	})
}
