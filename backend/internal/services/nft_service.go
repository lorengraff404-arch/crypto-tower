package services

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
	"gorm.io/gorm"
)

// NFTService handles NFT character minting and blockchain operations
type NFTService struct {
	client       *ethclient.Client
	contractAddr common.Address
	privateKey   *ecdsa.PrivateKey
	chainID      *big.Int
	cfg          *config.Config
}

// NewNFTService creates a new NFT service
func NewNFTService(cfg *config.Config) (*NFTService, error) {
	// Connect to opBNB testnet
	client, err := ethclient.Dial(cfg.OpBNBTestnetRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to opBNB: %v", err)
	}

	// Load private key (for server-side minting)
	privateKey, err := crypto.HexToECDSA(cfg.DeployerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	// Load contract address from config
	contractAddr := common.HexToAddress(cfg.CharacterNFTAddress)

	return &NFTService{
		client:       client,
		contractAddr: contractAddr,
		privateKey:   privateKey,
		chainID:      big.NewInt(5611), // opBNB testnet
		cfg:          cfg,
	}, nil
}

// MintFirstCharacter mints the first character for free (user pays gas)
func (s *NFTService) MintFirstCharacter(userID uint, characterID uint) error {
	// SECURITY CHECK 1: Verify user hasn't minted before
	var existingChar models.Character
	err := db.DB.Where("owner_id = ? AND on_chain_token_id IS NOT NULL", userID).First(&existingChar).Error
	if err == nil {
		return errors.New("user already minted first character")
	}

	// SECURITY CHECK 2: Get user and verify wallet
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	if user.WalletAddress == "" {
		return errors.New("wallet not connected - please connect your wallet first")
	}

	// SECURITY CHECK 3: Verify user has TOWER balance (any amount > 0)
	if user.TOWERBalance <= 0 {
		return errors.New("wallet must have TOWER balance to mint first character")
	}

	// SECURITY CHECK 4: Verify character ownership
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return errors.New("character not found")
	}

	if character.OwnerID != userID {
		return errors.New("you don't own this character")
	}

	if character.OnChainTokenID != nil {
		return errors.New("character already minted as NFT")
	}

	// SECURITY CHECK 5: Validate character is eligible for minting
	if character.Level < 1 {
		return errors.New("character must be at least level 1 to mint")
	}

	// Prepare metadata
	// metadata := s.PrepareMetadata(&character)

	// TODO: Upload metadata to IPFS
	metadataURI := fmt.Sprintf("ipfs://placeholder/%d", characterID)

	// BEGIN ATOMIC TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// For now, simulate minting (TODO: actual smart contract call)
	// In production, this would call the NFT contract's mint function
	tokenID := uint64(characterID)            // Simplified for now
	txHash := fmt.Sprintf("0x%064d", tokenID) // Placeholder

	// Update character with on-chain data
	character.OnChainTokenID = &tokenID
	character.MetadataURI = metadataURI
	character.MintTxHash = txHash
	character.IsMinted = true

	if err := tx.Save(&character).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update character")
	}

	// Mark user as having minted first character
	// user.HasMintedFirstChar = true // Field removed
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update user")
	}

	// Create audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "FIRST_CHARACTER_MINT",
		EntityType: "character",
		EntityID:   &characterID,
		NewValues:  fmt.Sprintf("token_id:%d,tx:%s", tokenID, txHash),
	}
	tx.Create(&auditLog)

	// Create transaction record
	transaction := models.Transaction{
		UserID:           userID,
		TransactionType:  "NFT_MINT",
		TokenType:        "TOWER",
		Amount:           0, // Free mint
		BalanceBefore:    user.TOWERBalance,
		BalanceAfter:     user.TOWERBalance,
		Description:      fmt.Sprintf("Minted first character #%d as NFT", characterID),
		BlockchainTxHash: &txHash,
		ChainID:          204, // opBNB testnet
		IsOnChain:        true,
		CharacterID:      &characterID,
	}
	tx.Create(&transaction)

	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction failed")
	}

	return nil
}

// MintCharacter mints a character as NFT (paid mint)
func (s *NFTService) MintCharacter(userID uint, characterID uint) error {
	mintCost := int64(1000) // 1000 TOWER to mint

	// Get user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// SECURITY CHECK 1: Verify wallet connected
	if user.WalletAddress == "" {
		return errors.New("wallet not connected")
	}

	// SECURITY CHECK 2: Verify sufficient TOWER balance
	if user.TOWERBalance < mintCost {
		return fmt.Errorf("insufficient TOWER balance (need %d, have %d)", mintCost, user.TOWERBalance)
	}

	// SECURITY CHECK 3: Verify character ownership
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return errors.New("character not found")
	}

	if character.OwnerID != userID {
		return errors.New("you don't own this character")
	}

	if character.OnChainTokenID != nil {
		return errors.New("character already minted as NFT")
	}

	// Prepare metadata
	// metadata := s.PrepareMetadata(&character)
	metadataURI := fmt.Sprintf("ipfs://placeholder/%d", characterID)

	// BEGIN ATOMIC TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Deduct mint cost
	if err := tx.Model(&user).Update("tower_balance", gorm.Expr("tower_balance - ?", mintCost)).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to deduct mint cost")
	}

	// Simulate minting
	tokenID := uint64(characterID + 1000000) // Offset for paid mints
	txHash := fmt.Sprintf("0x%064d", tokenID)

	// Update character
	character.OnChainTokenID = &tokenID
	character.MetadataURI = metadataURI
	character.MintTxHash = txHash
	character.IsMinted = true

	if err := tx.Save(&character).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update character")
	}

	// Create audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "CHARACTER_MINT",
		EntityType: "character",
		EntityID:   &characterID,
		NewValues:  fmt.Sprintf("token_id:%d,cost:%d,tx:%s", tokenID, mintCost, txHash),
	}
	tx.Create(&auditLog)

	// Create transaction record
	transaction := models.Transaction{
		UserID:           userID,
		TransactionType:  "NFT_MINT",
		TokenType:        "TOWER",
		Amount:           -mintCost,
		BalanceBefore:    user.TOWERBalance,
		BalanceAfter:     user.TOWERBalance - mintCost,
		Description:      fmt.Sprintf("Minted character #%d as NFT", characterID),
		BlockchainTxHash: &txHash,
		ChainID:          204,
		IsOnChain:        true,
		CharacterID:      &characterID,
	}
	tx.Create(&transaction)

	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction failed")
	}

	return nil
}

// PrepareMetadata creates NFT metadata for a character
func (s *NFTService) PrepareMetadata(character *models.Character) map[string]interface{} {
	attributes := []map[string]interface{}{
		{"trait_type": "Rarity", "value": character.Rarity},
		{"trait_type": "Class", "value": character.Class},
		{"trait_type": "Element", "value": character.Element},
		{"trait_type": "Level", "value": character.Level},
		{"trait_type": "Attack", "value": character.CurrentAttack, "display_type": "number"},
		{"trait_type": "Defense", "value": character.CurrentDefense, "display_type": "number"},
		{"trait_type": "HP", "value": character.CurrentHP, "display_type": "number"},
		{"trait_type": "Speed", "value": character.CurrentSpeed, "display_type": "number"},
	}

	metadata := map[string]interface{}{
		"name":         fmt.Sprintf("%s #%d", character.CharacterType, character.ID),
		"description":  fmt.Sprintf("A %s %s character from Crypto Tower Defense", character.Rarity, character.Class),
		"image":        fmt.Sprintf("ipfs://placeholder/images/%s.png", character.CharacterType),
		"attributes":   attributes,
		"external_url": fmt.Sprintf("https://cryptotowerdefense.com/character/%d", character.ID),
	}

	return metadata
}

// VerifyNFTOwnership verifies on-chain ownership of an NFT
func (s *NFTService) VerifyNFTOwnership(tokenID *uint64, expectedOwner string) bool {
	if tokenID == nil {
		return false
	}

	// TODO: Implement actual on-chain verification
	// For now, return true (placeholder)
	// In production, this would call the NFT contract's ownerOf function

	return true
}

// GetNFTOwner returns the on-chain owner of an NFT
func (s *NFTService) GetNFTOwner(tokenID uint64) (string, error) {
	// TODO: Implement actual on-chain query
	// For now, return placeholder
	// In production, this would call the NFT contract's ownerOf function

	return "0x0000000000000000000000000000000000000000", nil
}

// GetTowerBalance returns TOWER balance from on-chain
func (s *NFTService) GetTowerBalance(walletAddress string) (int64, error) {
	// TODO: Implement actual on-chain balance check
	// For now, return 1 (placeholder to allow first mint)
	// In production, this would call the TOWER token contract's balanceOf function

	return 1, nil
}

// MintNFT calls the smart contract to mint an NFT (placeholder)
func (s *NFTService) MintNFT(recipient string, metadataURI string) (uint64, string, error) {
	// TODO: Implement actual smart contract call
	// This would use go-ethereum to call the CharacterNFT contract's mint function

	// For now, return placeholder values
	tokenID := uint64(1000000)
	txHash := "0x0000000000000000000000000000000000000000000000000000000000000000"

	return tokenID, txHash, nil
}

// UploadToIPFS uploads metadata to IPFS (placeholder)
func (s *NFTService) UploadToIPFS(metadata map[string]interface{}) (string, error) {
	// TODO: Implement actual IPFS upload
	// This would use an IPFS client or Pinata/NFT.Storage API

	return "ipfs://QmPlaceholder", nil
}

// EstimateGas estimates gas cost for minting
func (s *NFTService) EstimateGas() (*big.Int, error) {
	ctx := context.Background()

	// Get current gas price
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	// Estimate gas limit for mint (typically ~100k gas)
	gasLimit := big.NewInt(100000)

	// Calculate total cost
	totalCost := new(big.Int).Mul(gasPrice, gasLimit)

	return totalCost, nil
}

// GetTransactionReceipt gets receipt for a transaction
func (s *NFTService) GetTransactionReceipt(txHash string) (bool, error) {
	ctx := context.Background()

	hash := common.HexToHash(txHash)
	receipt, err := s.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return false, err
	}

	// Check if transaction was successful
	return receipt.Status == 1, nil
}

// BuildMintTransaction builds a mint transaction for user to sign
func (s *NFTService) BuildMintTransaction(userAddress string, characterID uint) (map[string]interface{}, error) {
	// TODO: Implement transaction building
	// This would create the transaction data for the user's wallet to sign

	txData := map[string]interface{}{
		"to":    s.contractAddr.Hex(),
		"from":  userAddress,
		"data":  "0x", // Placeholder for encoded function call
		"value": "0",
	}

	return txData, nil
}

// GetChainID returns the chain ID
func (s *NFTService) GetChainID() *big.Int {
	return s.chainID
}

// GetContractAddress returns the NFT contract address
func (s *NFTService) GetContractAddress() string {
	return s.contractAddr.Hex()
}
