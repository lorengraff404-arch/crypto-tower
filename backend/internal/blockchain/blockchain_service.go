package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockchainService handles blockchain interactions
type BlockchainService struct {
	client          *ethclient.Client
	treasuryAddress common.Address
	gtkAddress      common.Address
	towerAddress    common.Address
}

// NewBlockchainService creates a new blockchain service
func NewBlockchainService(rpcURL, treasury, gtk, tower string) (*BlockchainService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	return &BlockchainService{
		client:          client,
		treasuryAddress: common.HexToAddress(treasury),
		gtkAddress:      common.HexToAddress(gtk),
		towerAddress:    common.HexToAddress(tower),
	}, nil
}

// TransactionVerification holds transaction verification data
type TransactionVerification struct {
	TxHash      string
	From        string
	To          string
	Amount      float64
	TokenType   string
	IsValid     bool
	BlockNumber uint64
}

// VerifyTOWERTransfer verifies a TOWER token transfer to treasury
func (s *BlockchainService) VerifyTOWERTransfer(txHash string, expectedAmount float64, senderAddress string) (*TransactionVerification, error) {
	// Check if service is initialized
	if s == nil {
		return nil, fmt.Errorf("blockchain service not initialized")
	}
	if s.client == nil {
		return nil, fmt.Errorf("blockchain client not connected")
	}
	return s.verifyTokenTransfer(txHash, expectedAmount, senderAddress, s.towerAddress, "TOWER")
}

// VerifyGTKTransfer verifies a GTK token transfer to treasury
func (s *BlockchainService) VerifyGTKTransfer(txHash string, expectedAmount float64, senderAddress string) (*TransactionVerification, error) {
	// Check if service is initialized
	if s == nil {
		return nil, fmt.Errorf("blockchain service not initialized")
	}
	if s.client == nil {
		return nil, fmt.Errorf("blockchain client not connected")
	}
	return s.verifyTokenTransfer(txHash, expectedAmount, senderAddress, s.gtkAddress, "GTK")
}

// verifyTokenTransfer is the core verification logic
func (s *BlockchainService) verifyTokenTransfer(txHashStr string, expectedAmount float64, senderAddress string, tokenAddress common.Address, tokenType string) (*TransactionVerification, error) {
	ctx := context.Background()

	// Parse transaction hash
	txHash := common.HexToHash(txHashStr)

	// Get transaction receipt
	receipt, err := s.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Check transaction status
	if receipt.Status != types.ReceiptStatusSuccessful {
		return &TransactionVerification{
			TxHash:  txHashStr,
			IsValid: false,
		}, fmt.Errorf("transaction failed or reverted")
	}

	// Get transaction details
	tx, isPending, err := s.client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	if isPending {
		return &TransactionVerification{
			TxHash:  txHashStr,
			IsValid: false,
		}, fmt.Errorf("transaction still pending")
	}

	// Verify transaction is an ERC20 transfer
	if tx.To() == nil || *tx.To() != tokenAddress {
		return &TransactionVerification{
			TxHash:  txHashStr,
			IsValid: false,
		}, fmt.Errorf("transaction is not to the %s token contract", tokenType)
	}

	// Parse ERC20 transfer event from logs
	var transferEvent *types.Log
	for _, vLog := range receipt.Logs {
		// ERC20 Transfer event signature: Transfer(address,address,uint256)
		// Keccak256 hash: 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
		transferSig := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

		if vLog.Topics[0] == transferSig && len(vLog.Topics) == 3 {
			transferEvent = vLog
			break
		}
	}

	if transferEvent == nil {
		return &TransactionVerification{
			TxHash:  txHashStr,
			IsValid: false,
		}, fmt.Errorf("no Transfer event found in transaction")
	}

	// Parse sender, receiver, and amount from event
	from := common.BytesToAddress(transferEvent.Topics[1].Bytes())
	to := common.BytesToAddress(transferEvent.Topics[2].Bytes())
	amount := new(big.Int).SetBytes(transferEvent.Data)

	// Convert amount from wei to tokens (18 decimals)
	amountFloat := new(big.Float).Quo(
		new(big.Float).SetInt(amount),
		new(big.Float).SetInt(big.NewInt(1e18)),
	)
	actualAmount, _ := amountFloat.Float64()

	// Verify sender
	if !strings.EqualFold(from.Hex(), senderAddress) {
		return &TransactionVerification{
			TxHash:      txHashStr,
			From:        from.Hex(),
			To:          to.Hex(),
			Amount:      actualAmount,
			TokenType:   tokenType,
			IsValid:     false,
			BlockNumber: receipt.BlockNumber.Uint64(),
		}, fmt.Errorf("sender address mismatch: expected %s, got %s", senderAddress, from.Hex())
	}

	// Verify receiver (must be treasury)
	if !strings.EqualFold(to.Hex(), s.treasuryAddress.Hex()) {
		return &TransactionVerification{
			TxHash:      txHashStr,
			From:        from.Hex(),
			To:          to.Hex(),
			Amount:      actualAmount,
			TokenType:   tokenType,
			IsValid:     false,
			BlockNumber: receipt.BlockNumber.Uint64(),
		}, fmt.Errorf("receiver address mismatch: expected %s (treasury), got %s", s.treasuryAddress.Hex(), to.Hex())
	}

	// Verify amount (allow 0.01% tolerance for rounding)
	tolerance := expectedAmount * 0.0001
	if actualAmount < expectedAmount-tolerance || actualAmount > expectedAmount+tolerance {
		return &TransactionVerification{
			TxHash:      txHashStr,
			From:        from.Hex(),
			To:          to.Hex(),
			Amount:      actualAmount,
			TokenType:   tokenType,
			IsValid:     false,
			BlockNumber: receipt.BlockNumber.Uint64(),
		}, fmt.Errorf("amount mismatch: expected %.6f, got %.6f", expectedAmount, actualAmount)
	}

	log.Printf("✅ %s transfer verified: %.6f from %s (tx: %s, block: %d)",
		tokenType, actualAmount, from.Hex(), txHashStr, receipt.BlockNumber.Uint64())

	return &TransactionVerification{
		TxHash:      txHashStr,
		From:        from.Hex(),
		To:          to.Hex(),
		Amount:      actualAmount,
		TokenType:   tokenType,
		IsValid:     true,
		BlockNumber: receipt.BlockNumber.Uint64(),
	}, nil
}

// GetTransactionStatus returns the status of a transaction
func (s *BlockchainService) GetTransactionStatus(txHashStr string) (string, error) {
	ctx := context.Background()
	txHash := common.HexToHash(txHashStr)

	receipt, err := s.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return "not_found", err
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		return "success", nil
	}

	return "failed", nil
}

// Close closes the blockchain client connection
func (s *BlockchainService) Close() {
	s.client.Close()
}
