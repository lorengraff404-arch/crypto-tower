package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/lorengraff/crypto-tower-defense/internal/blockchain/contracts"
)

// Client wraps Ethereum client with contract instances
type Client struct {
	ethClient    *ethclient.Client
	chainID      *big.Int
	privateKey   *ecdsa.PrivateKey
	fromAddress  common.Address
	
	// Contract instances
	gameToken    *contracts.GameToken
	towerToken   *contracts.TowerToken
	characterNFT *contracts.CharacterNFT
	itemNFT      *contracts.ItemNFT
	
	// Gas settings
	gasLimit     uint64
	maxGasPrice  *big.Int
}

// NewClient creates a new blockchain client
func NewClient(rpcURL string) (*Client, error) {
	// Connect to RPC
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}
	
	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}
	
	// Load private key from env
	privateKeyHex := os.Getenv("MINTER_PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, fmt.Errorf("MINTER_PRIVATE_KEY not set")
	}
	
	// Remove 0x prefix if present
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	
	// Load contract addresses
	gtkAddress := common.HexToAddress(os.Getenv("GAME_TOKEN_ADDRESS"))
	towerAddress := common.HexToAddress(os.Getenv("TOWER_TOKEN_ADDRESS"))
	charNFTAddress := common.HexToAddress(os.Getenv("CHARACTER_NFT_ADDRESS"))
	itemNFTAddress := common.HexToAddress(os.Getenv("ITEM_NFT_ADDRESS"))
	
	// Create contract instances
	gameToken, err := contracts.NewGameToken(gtkAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load GameToken contract: %w", err)
	}
	
	towerToken, err := contracts.NewTowerToken(towerAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load TowerToken contract: %w", err)
	}
	
	characterNFT, err := contracts.NewCharacterNFT(charNFTAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load CharacterNFT contract: %w", err)
	}
	
	itemNFT, err := contracts.NewItemNFT(itemNFTAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load ItemNFT contract: %w", err)
	}
	
	// Gas settings from env
	maxGasPriceGwei := int64(20) // Default 20 gwei
	if envGas := os.Getenv("MAX_GAS_PRICE_GWEI"); envGas != "" {
		fmt.Sscanf(envGas, "%d", &maxGasPriceGwei)
	}
	
	gasLimit := uint64(500000) // Default 500k
	if envLimit := os.Getenv("GAS_LIMIT"); envLimit != "" {
		fmt.Sscanf(envLimit, "%d", &gasLimit)
	}
	
	log.Printf("Blockchain client initialized:")
	log.Printf("  Chain ID: %s", chainID.String())
	log.Printf("  From Address: %s", fromAddress.Hex())
	log.Printf("  GameToken: %s", gtkAddress.Hex())
	log.Printf("  TowerToken: %s", towerAddress.Hex())
	log.Printf("  CharacterNFT: %s", charNFTAddress.Hex())
	log.Printf("  ItemNFT: %s", itemNFTAddress.Hex())
	
	return &Client{
		ethClient:    client,
		chainID:      chainID,
		privateKey:   privateKey,
		fromAddress:  fromAddress,
		gameToken:    gameToken,
		towerToken:   towerToken,
		characterNFT: characterNFT,
		itemNFT:      itemNFT,
		gasLimit:     gasLimit,
		maxGasPrice:  big.NewInt(maxGasPriceGwei * 1e9), // Convert to wei
	}, nil
}

// getTransactor creates a new transactor with current gas settings
func (c *Client) getTransactor(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := c.ethClient.PendingNonceAt(ctx, c.fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}
	
	gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	
	// Cap gas price
	if gasPrice.Cmp(c.maxGasPrice) > 0 {
		log.Printf("Warning: capping gas price from %s to %s wei", gasPrice.String(), c.maxGasPrice.String())
		gasPrice = c.maxGasPrice
	}
	
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = c.gasLimit
	auth.GasPrice = gasPrice
	
	return auth, nil
}

// MintGTK mints GTK tokens to specified address
func (c *Client) MintGTK(ctx context.Context, to common.Address, amount *big.Int) (*types.Transaction, error) {
	auth, err := c.getTransactor(ctx)
	if err != nil {
		return nil, err
	}
	
	tx, err := c.gameToken.Mint(auth, to, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to mint GTK: %w", err)
	}
	
	log.Printf("GTK mint transaction sent: %s", tx.Hash().Hex())
	return tx, nil
}

// MintTower mints TOWER tokens to specified address
func (c *Client) MintTower(ctx context.Context, to common.Address, amount *big.Int) (*types.Transaction, error) {
	auth, err := c.getTransactor(ctx)
	if err != nil {
		return nil, err
	}
	
	tx, err := c.towerToken.Mint(auth, to, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to mint TOWER: %w", err)
	}
	
	log.Printf("TOWER mint transaction sent: %s", tx.Hash().Hex())
	return tx, nil
}

// MintCharacterNFT mints a character NFT
func (c *Client) MintCharacterNFT(
	ctx context.Context,
	to common.Address,
	gameCharacterID *big.Int,
	characterType string,
	element string,
	rarity string,
	level *big.Int,
	tokenURI string,
) (*types.Transaction, error) {
	auth, err := c.getTransactor(ctx)
	if err != nil {
		return nil, err
	}
	
	tx, err := c.characterNFT.MintCharacter(auth, to, gameCharacterID, characterType, element, rarity, level, tokenURI)
	if err != nil {
		return nil, fmt.Errorf("failed to mint character NFT: %w", err)
	}
	
	log.Printf("CharacterNFT mint transaction sent: %s", tx.Hash().Hex())
	return tx, nil
}

// MintItemNFT mints item NFT(s)
func (c *Client) MintItemNFT(ctx context.Context, to common.Address, itemID *big.Int, amount *big.Int) (*types.Transaction, error) {
	auth, err := c.getTransactor(ctx)
	if err != nil {
		return nil, err
	}
	
	tx, err := c.itemNFT.Mint(auth, to, itemID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to mint item NFT: %w", err)
	}
	
	log.Printf("ItemNFT mint transaction sent: %s", tx.Hash().Hex())
	return tx, nil
}

// GetGTKBalance returns GTK balance of address
func (c *Client) GetGTKBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	balance, err := c.gameToken.BalanceOf(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get GTK balance: %w", err)
	}
	return balance, nil
}

// GetCharacterNFTCount returns number of character NFTs owned
func (c *Client) GetCharacterNFTCount(ctx context.Context, address common.Address) (*big.Int, error) {
	count, err := c.characterNFT.BalanceOf(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get character NFT count: %w", err)
	}
	return count, nil
}

// WaitForTransaction waits for transaction to be mined
func (c *Client) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	
	if receipt.Status == 0 {
		return receipt, fmt.Errorf("transaction failed")
	}
	
	return receipt, nil
}

// Close closes the blockchain client
func (c *Client) Close() {
	c.ethClient.Close()
}
