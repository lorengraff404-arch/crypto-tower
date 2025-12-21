#!/bin/bash

# Test Web3 Integration
echo "🧪 Testing Backend Web3 Integration"
echo "===================================="

# Check if Anvil is running
if ! lsof -i:8545 > /dev/null 2>&1; then
    echo "❌ Anvil not running on port 8545"
    echo "Start with: anvil --port 8545"
    exit 1
fi

echo "✅ Anvil running"

# Source env
cd backend
source .env

echo ""
echo "📋 Configuration:"
echo "  RPC: $OPBNB_RPC_URL"
echo "  GameToken: $GAME_TOKEN_ADDRESS"
echo "  TowerToken: $TOWER_TOKEN_ADDRESS"
echo "  CharacterNFT: $CHARACTER_NFT_ADDRESS"
echo "  ItemNFT: $ITEM_NFT_ADDRESS"

# Build test program
echo ""
echo "🔨 Building test program..."
cat > /tmp/test_web3.go << 'EOF'
package main

import (
	"context"
	"github.com/lorengraff/crypto-tower-defense/internal/blockchain"
	"fmt"
	"log"
	"math/big"
	"os"
	
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}
	
	// Create client
	rpcURL := os.Getenv("OPBNB_RPC_URL")
	client, err := blockchain.NewClient(rpcURL)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()
	
	fmt.Println("\n✅ Blockchain client created successfully!")
	
	// Test: Get GTK balance
	testAddress := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	balance, err := client.GetGTKBalance(context.Background(), testAddress)
	if err != nil {
		log.Fatal("Failed to get balance:", err)
	}
	
	fmt.Printf("\n📊 GTK Balance of %s: %s wei\n", testAddress.Hex(), balance.String())
	
	// Test: Mint 1000 GTK to test account
	fmt.Println("\n🪙 Minting 1000 GTK...")
	amount := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))
	tx, err := client.MintGTK(context.Background(), testAddress, amount)
	if err != nil {
		log.Fatal("Failed to mint:", err)
	}
	
	fmt.Printf("✅ Mint transaction: %s\n", tx.Hash().Hex())
	
	// Wait for confirmation
	fmt.Println("⏳ Waiting for confirmation...")
	receipt, err := client.WaitForTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal("Transaction failed:", err)
	}
	
	fmt.Printf("✅ Confirmed in block %d\n", receipt.BlockNumber.Uint64())
	
	// Check new balance
	newBalance, err := client.GetGTKBalance(context.Background(), testAddress)
	if err != nil {
		log.Fatal("Failed to get new balance:", err)
	}
	
	fmt.Printf("\n📊 New Balance: %s wei\n", newBalance.String())
	fmt.Printf("📈 Increase: %s GTK\n", new(big.Int).Div(new(big.Int).Sub(newBalance, balance), big.NewInt(1e18)).String())
	
	fmt.Println("\n🎉 Web3 integration test passed!")
}
EOF

go run /tmp/test_web3.go
