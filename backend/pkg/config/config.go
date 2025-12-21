package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost     string
	RedisPassword string

	// JWT
	JWTSecret string

	// Blockchain
	OpBNBTestnetRPC    string
	BSCTestnetRPC      string
	DeployerPrivateKey string

	// Smart Contract Addresses
	GTKTokenAddress     string
	TowerTokenAddress   string
	CharacterNFTAddress string
	ItemNFTAddress      string
	DeployerAddress     string

	// CORS
	AllowedOrigins []string

	// Rate Limiting
	RateLimitRequests int
	RateLimitDuration time.Duration

	// Game Economy
	TowerToGTKRatio      int
	MaxGTKPerUserDaily   int
	MaxTowerFromPvPDaily int
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{
		// Server
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "tower_defense_dev"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", ""),

		// Blockchain
		OpBNBTestnetRPC:    getEnv("OPBNB_TESTNET_RPC", "https://opbnb-testnet-rpc.bnbchain.org"),
		BSCTestnetRPC:      getEnv("BSC_TESTNET_RPC", "https://data-seed-prebsc-1-s1.binance.org:8545"),
		DeployerPrivateKey: getEnv("DEPLOYER_PRIVATE_KEY", ""),

		// Smart Contract Addresses
		GTKTokenAddress:     getEnv("GTK_TOKEN_ADDRESS", "0xFb0a39aE8c44a0E83a1445d4d272294345fA2207"),
		TowerTokenAddress:   getEnv("TOWER_TOKEN_ADDRESS", "0x4300536b909FbA47e042fCa31B97c09F64643110"),
		CharacterNFTAddress: getEnv("CHARACTER_NFT_ADDRESS", "0xe3765f851977Ed7B377D0234e9275845fc960775"),
		ItemNFTAddress:      getEnv("ITEM_NFT_ADDRESS", "0x8467806e70FbE05Ca5e17f5d316C09F5bD2391bC"),
		DeployerAddress:     getEnv("DEPLOYER_ADDRESS", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),

		// CORS
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},

		// Rate Limiting
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitDuration: time.Duration(getEnvAsInt("RATE_LIMIT_DURATION", 60)) * time.Second,

		// Game Economy
		TowerToGTKRatio:      getEnvAsInt("TOWER_TO_GTK_RATIO", 100),
		MaxGTKPerUserDaily:   getEnvAsInt("MAX_GTK_PER_USER_DAILY", 1000),
		MaxTowerFromPvPDaily: getEnvAsInt("MAX_TOWER_FROM_PVP_DAILY", 5000),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
