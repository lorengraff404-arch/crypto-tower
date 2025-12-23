package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/middleware"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

// AuthService handles authentication logic
type AuthService struct {
	cfg *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

// GenerateNonce creates a random nonce for signature verification
func (s *AuthService) GenerateNonce() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetOrCreateUser retrieves existing user or creates new one
func (s *AuthService) GetOrCreateUser(walletAddress string) (*models.User, error) {
	walletAddress = common.HexToAddress(walletAddress).Hex() // Normalize address

	var user models.User
	result := db.DB.Where("wallet_address = ?", walletAddress).First(&user)

	if result.Error == nil {
		// User exists, generate new nonce
		nonce, err := s.GenerateNonce()
		if err != nil {
			return nil, err
		}
		user.Nonce = nonce
		db.DB.Save(&user)
		return &user, nil
	}

	// Create new user
	nonce, err := s.GenerateNonce()
	if err != nil {
		return nil, err
	}

	// SUPER_ADMIN Fallback: If this is the first user, make them SUPER_ADMIN
	var userCount int64
	db.DB.Model(&models.User{}).Count(&userCount)
	role := "PLAYER"
	if userCount == 0 {
		role = "SUPER_ADMIN"
	}

	user = models.User{
		WalletAddress: walletAddress,
		Nonce:         nonce,
		Role:          role, // Set assigned role
		Level:         1,
		Rank:          "Cadete",
		RankTier:      1,
		ELO:           1000,
		GTKBalance:    0,
		TOWERBalance:  100,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// VerifySignature validates wallet signature
func (s *AuthService) VerifySignature(walletAddress, signature, message string) (bool, error) {
	// Normalize address
	walletAddress = common.HexToAddress(walletAddress).Hex()

	// Decode signature
	sig, err := hex.DecodeString(signature[2:]) // Remove "0x" prefix
	if err != nil {
		return false, fmt.Errorf("invalid signature format: %w", err)
	}

	// Adjust signature format (Ethereum uses v = 27/28, we need 0/1)
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	// Hash message with Ethereum prefix
	messageHash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	// Recover public key
	pubKey, err := crypto.SigToPub(messageHash.Bytes(), sig)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	// Get address from public key
	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	return recoveredAddress == walletAddress, nil
}

// GenerateJWT creates a JWT token for authenticated user
func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
	claims := middleware.Claims{
		UserID:        user.ID,
		WalletAddress: user.WalletAddress,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "crypto-tower-defense",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// UpdateLastLogin updates user's last login timestamp
func (s *AuthService) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return db.DB.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", now).Error
}
