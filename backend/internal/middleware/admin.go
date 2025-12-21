package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

// AdminMiddleware verifies that the user is an admin
func AdminMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get wallet address from context (set by AuthMiddleware)
		walletAddress, exists := c.Get("wallet_address")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		// Get admin wallet from environment
		adminWallet := os.Getenv("ADMIN_WALLET")
		if adminWallet == "" {
			// Fallback to deployer wallet
			adminWallet = os.Getenv("DEPLOYER_WALLET")
		}
		if adminWallet == "" {
			adminWallet = "0xdCb8ca66Ae0809Eed5dB73E8e1c3787c8178327e" // Your deployer wallet
		}

		// Compare wallet addresses (case-insensitive)
		userWallet := walletAddress.(string)
		if !strings.EqualFold(userWallet, adminWallet) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Access denied - admin only",
				"required_wallet": adminWallet,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
