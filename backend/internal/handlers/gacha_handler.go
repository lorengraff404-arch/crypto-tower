package handlers

import (
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// GachaHandler handles gacha/minting operations
type GachaHandler struct {
	gachaService      *services.GachaService
	blockchainService *services.BlockchainService
	questService      *services.DailyQuestService
	// TODO: Add revenueService when RevenueService is created
}

// NewGachaHandler creates a new gacha handler
func NewGachaHandler(blockchainService *services.BlockchainService, db interface{}) *GachaHandler {
	return &GachaHandler{
		gachaService:      services.NewGachaService(blockchainService),
		blockchainService: blockchainService,
		questService:      services.NewDailyQuestService(),
	}
}

// MintEgg mints a new egg with gacha mechanics
// POST /api/v1/gacha/mint
func (h *GachaHandler) MintEgg(c *gin.Context) {
	userID := c.GetUint("user_id")
	// walletAddress := c.GetString("wallet_address") // Unused now

	var req struct {
		TowerAmount int64  `json:"tower_amount" binding:"min=0,max=10000"`
		TxHash      string `json:"tx_hash"` // Blockchain transaction hash (optional for free mint)
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If TOWER amount > 0, verify blockchain transaction
	if req.TowerAmount > 0 {
		// Skip blockchain verification if service is not initialized (development mode)
		if h.blockchainService == nil {
			log.Printf("⚠️ Blockchain service not initialized - skipping verification for %d TOWER (DEVELOPMENT MODE)", req.TowerAmount)
		} else {
			if req.TxHash == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction hash required for paid mint"})
				return
			}

			// Verify TOWER payment on blockchain using generic service
			// NOTE: VerifyTransaction checks strict Amount and Recipient (Treasury)
			// It does not currently verify SENDER against walletAddress, but for a purchase,
			// as long as the Treasury received money, we generally accept it.
			// (Use big.Int for amount)
			costBig := big.NewInt(req.TowerAmount)

			if err := h.blockchainService.VerifyTransaction(req.TxHash, costBig); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid or failed blockchain transaction",
					"details": err.Error(),
				})
				return
			}

			log.Printf("✅ Verified TOWER payment: %d (tx: %s)", req.TowerAmount, req.TxHash)
		}
	}

	// MintEgg with gacha service
	egg, err := h.gachaService.MintEgg(userID, req.TowerAmount, req.TxHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: If TOWER was paid, distribute revenue
	// if req.TowerAmount > 0 {
	//     err = h.revenueService.DistributeGTKRevenue("gacha_mint", float64(req.TowerAmount))
	//     if err != nil {
	//         log.Printf("⚠️ Revenue distribution failed: %v", err)
	//     }
	// }

	// Track quest progress for minting/incubation
	if err := h.questService.TrackProgress(userID, "incubation_started", 1, ""); err != nil {
		log.Printf("⚠️ Failed to track quest progress: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Egg minted successfully!",
		"egg":     egg,
		"tx_hash": req.TxHash,
	})
}

// GetOddsPreview returns probability preview for a given TOWER amount
// GET /api/v1/gacha/odds/:amount
func (h *GachaHandler) GetOddsPreview(c *gin.Context) {
	amountStr := c.Param("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	if amount < 0 || amount > 10000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be between 0 and 10,000"})
		return
	}

	odds := h.gachaService.GetOddsPreview(amount)

	c.JSON(http.StatusOK, gin.H{
		"amount": amount,
		"odds":   odds,
	})
}

// GetMyEggs returns user's eggs
// GET /api/v1/gacha/my-eggs
func (h *GachaHandler) GetMyEggs(c *gin.Context) {
	userID := c.GetUint("user_id")

	eggs, err := h.gachaService.GetUserEggs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"eggs":  eggs,
		"count": len(eggs),
	})
}

// StartIncubation starts egg incubation
// POST /api/v1/gacha/start-incubation/:id
func (h *GachaHandler) StartIncubation(c *gin.Context) {
	userID := c.GetUint("user_id")
	eggIDStr := c.Param("id")
	eggID, err := strconv.ParseUint(eggIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid egg ID"})
		return
	}

	if err := h.gachaService.StartIncubation(userID, uint(eggID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Incubation started",
	})
}

// HatchEgg hatches an egg into a character
// POST /api/v1/gacha/hatch/:id
func (h *GachaHandler) HatchEgg(c *gin.Context) {
	userID := c.GetUint("user_id")
	eggIDStr := c.Param("id")
	eggID, err := strconv.ParseUint(eggIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid egg ID"})
		return
	}

	character, err := h.gachaService.HatchEgg(userID, uint(eggID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Track quest progress for hatching
	if err := h.questService.TrackProgress(userID, "egg_hatched", 1, ""); err != nil {
		log.Printf("⚠️ Failed to track quest progress: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Egg hatched successfully!",
		"character": character,
	})
}

// ScanEgg reveals egg stats before hatching
// POST /api/v1/gacha/scan-egg/:id
func (h *GachaHandler) ScanEgg(c *gin.Context) {
	userID := c.GetUint("user_id")
	eggIDStr := c.Param("id")
	eggID, err := strconv.ParseUint(eggIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid egg ID"})
		return
	}

	stats, err := h.gachaService.ScanEgg(userID, uint(eggID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Egg scanned successfully",
		"stats":   stats,
	})
}

// ApplyAccelerator applies an accelerator item to an egg
// POST /api/v1/gacha/apply-accelerator
func (h *GachaHandler) ApplyAccelerator(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		EggID  uint `json:"egg_id" binding:"required"`
		ItemID uint `json:"item_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.gachaService.ApplyAccelerator(userID, req.EggID, req.ItemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Accelerator applied successfully",
	})
}
