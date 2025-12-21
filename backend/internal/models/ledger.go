package models

import (
	"time"

	"gorm.io/gorm"
)

// AccountType defines the type of a ledger account
type AccountType string

const (
	AccountTypeWallet   AccountType = "WALLET"   // Player's main wallet
	AccountTypeTreasury AccountType = "TREASURY" // Game treasury
	AccountTypeEscrow   AccountType = "ESCROW"   // Temporary hold (e.g., during battle)
	AccountTypeReward   AccountType = "REWARD"   // Source of rewards (inflationary)
	AccountTypeSink     AccountType = "SINK"     // Destination for burnt tokens
)

// LedgerAccount represents a logical account in the double-entry system
// A User can have multiple accounts (e.g., Wallet, Escrow)
type LedgerAccount struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    *uint          `gorm:"index" json:"user_id,omitempty"` // Null for system accounts
	Type      AccountType    `gorm:"type:varchar(20);not null" json:"type"`
	Currency  string         `gorm:"type:varchar(10);default:'GTK'" json:"currency"`
	Balance   int64          `gorm:"default:0" json:"balance"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TransactionType defines the business context of a transaction
type TransactionType string

const (
	TxTypeDeposit      TransactionType = "DEPOSIT"
	TxTypeWithdraw     TransactionType = "WITHDRAWAL" // Changed from WITHDRAW
	TxTypeWagerEnter   TransactionType = "WAGER_ENTER"
	TxTypeWagerWin     TransactionType = "WAGER_WIN"
	TxTypeWagerRefund  TransactionType = "WAGER_REFUND" // Added
	TxTypeWagerFee     TransactionType = "WAGER_FEE"
	TxTypeGachaMint    TransactionType = "GACHA_MINT"
	TxTypeBreedingFee  TransactionType = "BREEDING_FEE"    // Added
	TxTypeShopBuy      TransactionType = "SHOP_PURCHASE"   // Added
	TxTypeMarketBuy    TransactionType = "MARKET_PURCHASE" // Changed from MARKET_BUY
	TxTypeMarketSell   TransactionType = "MARKET_SALE"     // Added
	TxTypeMarketFee    TransactionType = "MARKET_FEE"
	TxTypeRevenueDist  TransactionType = "REVENUE_DIST" // Added
	TxTypeAdminAdj     TransactionType = "ADMIN_ADJUSTMENT"
	TxTypeRaidReward   TransactionType = "RAID_REWARD"   // Added
	TxTypeRankedReward TransactionType = "RANKED_REWARD" // Added
	TxTypeReward       TransactionType = "REWARD"        // Generic reward
)

// LedgerTransaction groups entries required to balance (Sum Debits = Sum Credits)
type LedgerTransaction struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Type        TransactionType `gorm:"type:varchar(30);not null" json:"type"`
	ReferenceID string          `gorm:"index" json:"reference_id"` // e.g., BattleID, TxHash
	Description string          `json:"description"`
	Timestamp   time.Time       `json:"timestamp"`
	Metadata    string          `gorm:"type:jsonb" json:"metadata"` // Flexible details

	Entries []LedgerEntry `gorm:"foreignKey:TransactionID" json:"entries"`
}

// LedgerEntry is a single debit or credit line
type LedgerEntry struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `gorm:"index;not null" json:"transaction_id"`
	AccountID     uint      `gorm:"index;not null" json:"account_id"`
	Amount        int64     `gorm:"not null" json:"amount"`       // Positive = Credit, Negative = Debit (or vice-versa, depending on convention. Here: Amount > 0 is Credit to account, Amount < 0 is Debit)
	Type          string    `gorm:"type:varchar(10)" json:"type"` // "DEBIT" or "CREDIT" explicitly
	CreatedAt     time.Time `json:"created_at"`
}
