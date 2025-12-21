package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

type InventoryHandler struct {
	lootService *services.LootService
}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{
		lootService: &services.LootService{},
	}
}

// GetInventory returns user's inventory (Stackables + Unique Items)
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	type InventoryItem struct {
		ID       uint   `json:"id"` // Unique ID or ShopItem ID
		Name     string `json:"name"`
		ItemType string `json:"item_type"` // CONSUMABLE, WEAPON, etc.
		Rarity   string `json:"rarity"`
		Quantity int    `json:"quantity"`
		IconURL  string `json:"icon_url"`
		IsUnique bool   `json:"is_unique"` // True for Equipment, False for Stackables
		Effect   string `json:"effect,omitempty"`
	}

	var responseList []InventoryItem

	// 1. Fetch Stackables (UserInventory -> ShopItem)
	// We need a helper or direct DB query here.
	// To avoid circular deps or complex service logic, we'll query DB directly here for now,
	// or move this to ItemService later. Direct DB is acceptable for Handler in this context.
	// But services is better. Let's use LootService or ItemService if possible.
	// LootService has AddToInventory but not GetInventory.
	// Let's add GetStackableInventory to LootService? Or ItemService.
	// We'll trust ItemService to handle Unique items.
	// We'll add a specific struct for stackables query here.

	type StackableResult struct {
		models.UserInventory
		ShopItem models.ShopItem `gorm:"foreignKey:ItemID"`
	}
	var userStackables []models.UserInventory
	if err := db.DB.Preload("Item").Where("user_id = ?", userID).Find(&userStackables).Error; err == nil {
		for _, inv := range userStackables {
			if inv.Quantity > 0 {
				responseList = append(responseList, InventoryItem{
					ID:       inv.ItemID, // ShopItem ID
					Name:     inv.Item.Name,
					ItemType: "CONSUMABLE", // ShopItems are mostly consumables
					Rarity:   "C",          // Default
					Quantity: inv.Quantity,
					IconURL:  inv.Item.IconURL,
					IsUnique: false,
					Effect:   inv.Item.EffectType,
				})
			}
		}
	}

	// 2. Fetch Unique Items (models.Item)
	// We can use ItemService for this
	uniqueService := services.NewItemService()
	uniqueItems, err := uniqueService.GetUserItems(userID, "")
	if err == nil {
		for _, item := range uniqueItems {
			responseList = append(responseList, InventoryItem{
				ID:       item.ID,
				Name:     item.Name,
				ItemType: item.ItemType,
				Rarity:   item.Rarity,
				Quantity: item.Quantity,
				IconURL:  item.IconURL,
				IsUnique: true,
				Effect:   item.SpecialEffects,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"inventory": responseList,
	})
}

// UseItem uses a consumable item
func (h *InventoryHandler) UseItem(c *gin.Context) {

	var req struct {
		CharacterID uint `json:"character_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement item usage logic
	// - Check if user owns item
	// - Apply item effect
	// - Decrement quantity

	c.JSON(http.StatusOK, gin.H{
		"message": "Item used successfully",
	})
}

// DiscardItem removes item from inventory
func (h *InventoryHandler) DiscardItem(c *gin.Context) {

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement discard logic
	// - Verify ownership
	// - Decrement or remove item

	c.JSON(http.StatusOK, gin.H{
		"message": "Item discarded",
	})
}
