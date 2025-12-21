package handlers

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(cfg),
	}
}

// GetNonce returns a nonce for wallet signature
// POST /api/v1/auth/nonce
func (h *AuthHandler) GetNonce(c *gin.Context) {
	var req struct {
		WalletAddress string `json:"wallet_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Get or create user
	user, err := h.authService.GetOrCreateUser(req.WalletAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nonce":  user.Nonce,
		"message": "Sign this message to authenticate: " + user.Nonce,
	})
}

// VerifySignature verifies wallet signature and returns JWT
// POST /api/v1/auth/verify
func (h *AuthHandler) VerifySignature(c *gin.Context) {
	var req struct {
		WalletAddress string `json:"wallet_address" binding:"required"`
		Signature     string `json:"signature" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// CRITICAL: Normalize wallet address (same as GetOrCreateUser does)
	normalizedAddress := common.HexToAddress(req.WalletAddress).Hex()

	// Get user WITHOUT regenerating nonce (just fetch from DB)
	var user models.User
	result := db.DB.Where("wallet_address = ?", normalizedAddress).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found or nonce expired",
		})
		return
	}

	// Build message that was signed (must match the nonce from GET /nonce)
	message := "Sign this message to authenticate: " + user.Nonce

	// Verify signature
	valid, err := h.authService.VerifySignature(req.WalletAddress, req.Signature, message)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid signature",
		})
		return
	}

	// Generate JWT
	token, err := h.authService.GenerateJWT(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Update last login
	_ = h.authService.UpdateLastLogin(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":             user.ID,
			"wallet_address": user.WalletAddress,
			"level":          user.Level,
			"rank":           user.Rank,
			"gtk_balance":    user.GTKBalance,
			"tower_balance":  user.TOWERBalance,
		},
	})
}

// GetProfile returns authenticated user's profile
// GET /api/v1/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user": gin.H{
			"id":             user.ID,
			"wallet_address": user.WalletAddress,
			"level":          user.Level,
			"experience":     user.Experience,
			"rank":           user.Rank,
			"rank_tier":      user.RankTier,
			"elo":            user.ELO,
			"gtk_balance":    user.GTKBalance,
			"tower_balance":  user.TOWERBalance,
		},
	})
}
