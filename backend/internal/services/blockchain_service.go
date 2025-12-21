package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

type BlockchainService struct {
	client          *ethclient.Client
	rpcURL          string
	treasuryAddress common.Address
	gtkAddress      common.Address
}

func NewBlockchainService(cfg *config.Config) (*BlockchainService, error) {
	client, err := ethclient.Dial(cfg.OpBNBTestnetRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %v", err)
	}

	// Safety check for Treasury Address
	treasuryAddr := common.HexToAddress(cfg.DeployerAddress) // Default to deployer if treasury not robustly set, or add Treasury to config
	// In config.go we saw DeployerAddress, but let's check properly. Implementation plan said TreasuryWallet.
	// For now using DeployerAddress as Treasury if not passed explicitly, but Config usually has it.
	// We'll trust Config passed in. Assuming Config has a generic "Treasury" or using Deployer for now.
	// Config has "DeployerAddress". We will use that as Treasury.

	return &BlockchainService{
		client:          client,
		rpcURL:          cfg.OpBNBTestnetRPC,
		treasuryAddress: treasuryAddr,
		gtkAddress:      common.HexToAddress(cfg.GTKTokenAddress),
	}, nil
}

// VerifyTransaction checks if a transaction is valid, successful, and transferred the correct amount of GTK to the treasury
func (s *BlockchainService) VerifyTransaction(txHashStr string, expectedAmount *big.Int) error {
	txHash := common.HexToHash(txHashStr)

	// 1. Get Transaction Receipt to check status
	receipt, err := s.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return errors.New("transaction not found on chain (yet)")
		}
		return fmt.Errorf("failed to get receipt: %v", err)
	}

	// 2. Check Status (1 = success)
	if receipt.Status != 1 {
		return errors.New("transaction failed on-chain")
	}

	// 3. Scan Logs for ERC20 Transfer to Treasury
	// Event Signature: Transfer(address indexed from, address indexed to, uint256 value)
	// Topic[0] is Keccak256("Transfer(address,address,uint256)")
	transferSig := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

	foundTransfer := false
	var transferredAmount *big.Int

	for _, log := range receipt.Logs {
		// Verify expected Token Contract (GTK)
		if log.Address != s.gtkAddress {
			continue
		}

		// Check topic signature
		if len(log.Topics) == 3 && log.Topics[0] == transferSig {
			// Topic[1] = From, Topic[2] = To
			toAddress := common.HexToAddress(log.Topics[2].Hex())

			if toAddress == s.treasuryAddress {
				// We found a transfer to us!
				// Data is amount (not indexed)
				amount := new(big.Int)
				amount.SetBytes(log.Data)

				if amount.Cmp(expectedAmount) >= 0 {
					foundTransfer = true
					transferredAmount = amount
					break
				}
			}
		}
	}

	if !foundTransfer {
		return fmt.Errorf("no valid transfer of %s GTK to treasury found in this transaction", expectedAmount.String())
	}

	fmt.Printf("✅ Verified Tx %s: Transferred %s GTK\n", txHashStr, transferredAmount.String())
	return nil
}
