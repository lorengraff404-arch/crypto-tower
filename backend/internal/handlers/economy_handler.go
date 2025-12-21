package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// EconomyHandler handles token economy operations
type EconomyHandler struct {
	tokenService *services.TokenService
}

// NewEconomyHandler creates a new economy handler
func NewEconomyHandler() *EconomyHandler {
	return &EconomyHandler{
		tokenService: services.NewTokenService(),
	}
}

// ConvertTowerToGTK converts TOWER to GTK
// POST /api/v1/economy/convert/tower-to-gtk
func (h *EconomyHandler) ConvertTowerToGTK(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Amount int64 `json:"amount" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tokenService.ConvertTowerToGTK(userID, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Conversion successful",
	})
}

// ConvertGTKToTower converts GTK to TOWER
// POST /api/v1/economy/convert/gtk-to-tower
func (h *EconomyHandler) ConvertGTKToTower(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Amount int64 `json:"amount" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tokenService.ConvertGTKToTower(userID, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Conversion successful",
	})
}

// WithdrawTower initiates TOWER withdrawal
// POST /api/v1/economy/withdraw
func (h *EconomyHandler) WithdrawTower(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Amount        int64  `json:"amount" binding:"required,min=100"`
		WalletAddress string `json:"wallet_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tokenService.WithdrawTower(userID, req.Amount, req.WalletAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Withdrawal request submitted",
	})
}

// DepositTower processes TOWER deposit
// POST /api/v1/economy/deposit
func (h *EconomyHandler) DepositTower(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		TxHash string `json:"tx_hash" binding:"required"`
		Amount int64  `json:"amount" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tokenService.DepositTower(userID, req.TxHash, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deposit confirmed",
	})
}

// GetBalance returns user's token balances
// GET /api/v1/economy/balance
func (h *EconomyHandler) GetBalance(c *gin.Context) {
	userID := c.GetUint("user_id")

	balances, err := h.tokenService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balances": balances,
	})
}

// GetTransactionHistory returns user's transaction history
// GET /api/v1/economy/history
func (h *EconomyHandler) GetTransactionHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	if limit > 100 {
		limit = 100
	}

	transactions, err := h.tokenService.GetTransactionHistory(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
	})
}
