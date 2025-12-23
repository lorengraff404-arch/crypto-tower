package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

// AdminMiddleware verifies that the user is an admin
func AdminMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		// RBAC Check via DB
		var user struct {
			Role          string
			WalletAddress string
		}
		if err := db.DB.Table("users").Select("role, wallet_address").Where("id = ?", userID).Scan(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User validation failed"})
			c.Abort()
			return
		}

		// MAX SECURITY: Check for specific deployer wallet
		const masterWallet = "0xdCb8ca66Ae0809Eed5dB73E8e1c3787c8178327e"
		if user.WalletAddress != masterWallet && user.Role != "SUPER_ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied - Unauthorized wallet address",
			})
			c.Abort()
			return
		}

		// Allow ADMIN or SUPER_ADMIN
		if user.Role != "ADMIN" && user.Role != "SUPER_ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied - Admin role required",
				"role":  user.Role,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
