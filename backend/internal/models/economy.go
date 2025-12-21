package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction represents economy transactions for audit trail
type Transaction struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// User
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	// Transaction Details
	TransactionType string `gorm:"type:varchar(30);not null;index" json:"transaction_type"`
	// Types: BATTLE_REWARD, ISLAND_REWARD, BREEDING_FEE, CRAFTING_FEE, MARKETPLACE_BUY,
	//        MARKETPLACE_SELL, WITHDRAWAL, DEPOSIT, TOWER_TO_GTK_CONVERSION, DAILY_QUEST

	TokenType     string `gorm:"type:varchar(10);not null;index" json:"token_type"` // GTK or TOWER
	Amount        int64  `gorm:"not null" json:"amount"`                            // Positive for credit, negative for debit
	BalanceBefore int64  `gorm:"not null" json:"balance_before"`
	BalanceAfter  int64  `gorm:"not null" json:"balance_after"`

	// Related Entities (nullable)
	BattleID    *uint `gorm:"index" json:"battle_id,omitempty"`
	CharacterID *uint `gorm:"index" json:"character_id,omitempty"`
	ItemID      *uint `gorm:"index" json:"item_id,omitempty"`

	// Metadata
	Description string `gorm:"type:varchar(255)" json:"description"`
	Metadata    string `gorm:"type:text" json:"metadata,omitempty"` // JSON for additional data

	// Blockchain (for withdrawals/deposits)
	BlockchainTxHash *string `gorm:"type:varchar(66);uniqueIndex" json:"blockchain_tx_hash,omitempty"`
	ChainID          int     `json:"chain_id,omitempty"` // 204 for opBNB, 56 for BSC
	IsOnChain        bool    `gorm:"default:false;index" json:"is_on_chain"`
}

// DailyLimit tracks daily emission caps per user
type DailyLimit struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID uint   `gorm:"uniqueIndex:idx_user_date" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"-"`
	Date   string `gorm:"type:date;uniqueIndex:idx_user_date" json:"date"` // YYYY-MM-DD

	GTKEarned   int64 `gorm:"default:0;not null" json:"gtk_earned"`
	TOWEREarned int64 `gorm:"default:0;not null" json:"tower_earned"`

	// Caps from config
	GTKLimit   int64 `gorm:"default:1000;not null" json:"gtk_limit"`
	TOWERLimit int64 `gorm:"default:5000;not null" json:"tower_limit"`

	CanEarnMore bool `gorm:"default:true" json:"can_earn_more"`
}

// MarketplaceListing represents an active listing
type MarketplaceListing struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Seller
	SellerID uint `gorm:"not null;index" json:"seller_id"`
	Seller   User `gorm:"foreignKey:SellerID" json:"seller,omitempty"`

	// Asset
	AssetType   string     `gorm:"type:varchar(20);not null;index" json:"asset_type"` // CHARACTER, ITEM, EGG
	CharacterID *uint      `gorm:"index" json:"character_id,omitempty"`
	Character   *Character `gorm:"foreignKey:CharacterID" json:"character,omitempty"`
	ItemID      *uint      `gorm:"index" json:"item_id,omitempty"`
	Item        *Item      `gorm:"foreignKey:ItemID" json:"item,omitempty"`

	// Pricing
	Price      int64 `gorm:"not null" json:"price"`         // In TOWER tokens
	ListingFee int64 `gorm:"default:10" json:"listing_fee"` // Flat fee in GTK

	// Status
	Status    string     `gorm:"type:varchar(20);not null;index;default:'ACTIVE'" json:"status"` // ACTIVE, SOLD, CANCELLED, EXPIRED
	ExpiresAt time.Time  `gorm:"index" json:"expires_at"`                                        // Max 7 days
	SoldAt    *time.Time `json:"sold_at,omitempty"`
	BuyerID   *uint      `json:"buyer_id,omitempty"`
	Buyer     *User      `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`

	// Featured (paid promotion)
}

// LoginReward tracks daily login rewards (Phase 21)
type LoginReward struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	DayNumber    int       `gorm:"not null" json:"day_number"`
	ClaimedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"claimed_at"`
	RewardTokens int       `json:"reward_tokens"`
	RewardItems  string    `gorm:"type:text" json:"reward_items"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TradeHistory tracks completed trades (Phase 19)
type TradeHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ListingID   *uint     `json:"listing_id"`
	SellerID    uint      `gorm:"not null" json:"seller_id"`
	BuyerID     uint      `gorm:"not null" json:"buyer_id"`
	ItemType    string    `gorm:"size:30;not null" json:"item_type"`
	ItemID      uint      `gorm:"not null" json:"item_id"`
	Price       int64     `gorm:"not null" json:"price"`
	Currency    string    `gorm:"size:20;not null" json:"currency"`
	CompletedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"completed_at"`

	Seller User `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Buyer  User `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
}
