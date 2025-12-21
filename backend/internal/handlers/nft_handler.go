package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

// NFTHandler handles NFT minting operations
type NFTHandler struct {
	nftService *services.NFTService
}

// NewNFTHandler creates a new NFT handler
func NewNFTHandler(cfg *config.Config) (*NFTHandler, error) {
	nftService, err := services.NewNFTService(cfg)
	if err != nil {
		return nil, err
	}

	return &NFTHandler{
		nftService: nftService,
	}, nil
}

// MintFirstCharacter mints the first character for free
// POST /api/v1/nft/mint-first
func (h *NFTHandler) MintFirstCharacter(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		CharacterID uint `json:"character_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.nftService.MintFirstCharacter(userID, req.CharacterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "First character minted successfully!",
	})
}

// MintCharacter mints a character as NFT (paid)
// POST /api/v1/nft/mint
func (h *NFTHandler) MintCharacter(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		CharacterID uint `json:"character_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.nftService.MintCharacter(userID, req.CharacterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Character minted successfully!",
	})
}

// GetGasEstimate returns gas estimate for minting
// GET /api/v1/nft/gas-estimate
func (h *NFTHandler) GetGasEstimate(c *gin.Context) {
	gasEstimate, err := h.nftService.EstimateGas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gas_estimate": gasEstimate.String(),
	})
}

// GetContractInfo returns NFT contract information
// GET /api/v1/nft/contract-info
func (h *NFTHandler) GetContractInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"contract_address": h.nftService.GetContractAddress(),
		"chain_id":         h.nftService.GetChainID().String(),
		"network":          "opBNB Testnet",
	})
}

// BuildMintTx builds a mint transaction for user to sign
// POST /api/v1/nft/build-mint-tx
func (h *NFTHandler) BuildMintTx(c *gin.Context) {
	var req struct {
		CharacterID uint   `json:"character_id" binding:"required"`
		UserAddress string `json:"user_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txData, err := h.nftService.BuildMintTransaction(req.UserAddress, req.CharacterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": txData,
	})
}

// VerifyOwnership verifies on-chain NFT ownership
// GET /api/v1/nft/verify/:tokenId
func (h *NFTHandler) VerifyOwnership(c *gin.Context) {
	tokenIDStr := c.Param("tokenId")
	expectedOwner := c.Query("owner")

	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
		return
	}

	isOwner := h.nftService.VerifyNFTOwnership(&tokenID, expectedOwner)

	c.JSON(http.StatusOK, gin.H{
		"is_owner": isOwner,
		"token_id": tokenID,
	})
}
