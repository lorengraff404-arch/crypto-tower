package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a player in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// RBAC & Security (Added for Enterprise Governance)
	Role          string     `gorm:"type:varchar(20);default:'PLAYER';not null" json:"role"`   // PLAYER, MODERATOR, ADMIN, SUPER_ADMIN
	Status        string     `gorm:"type:varchar(20);default:'ACTIVE';not null" json:"status"` // ACTIVE, SUSPENDED, BANNED
	BanReason     string     `json:"ban_reason,omitempty"`
	BanExpiresAt  *time.Time `json:"ban_expires_at,omitempty"`
	SecurityFlags int        `gorm:"default:0" json:"-"`

	// Wallet & Authentication
	WalletAddress string     `gorm:"uniqueIndex;not null" json:"wallet_address"`
	Nonce         string     `gorm:"not null" json:"-"` // For signature verification
	LastLoginAt   *time.Time `json:"last_login_at"`

	// Player Stats
	Level      int    `gorm:"default:1;not null" json:"level"`
	Experience int    `gorm:"default:0;not null" json:"experience"`
	Rank       string `gorm:"type:varchar(20);default:'Cadete'" json:"rank"` // Cadete, Biólogo, Armero, etc.
	RankTier   int    `gorm:"type:integer;default:1" json:"rank_tier"`       // 1=C, 2=B, 3=A, 4=S, 5=SS, 6=SSS
	ELO        int    `gorm:"column:elo_rating;default:1000" json:"elo"`     // ELO rating for ranked matchmaking

	// Economy
	GTKBalance        int64      `gorm:"default:0;not null" json:"gtk_balance"`   // In-game currency
	TOWERBalance      int64      `gorm:"default:0;not null" json:"tower_balance"` // Governance token (off-chain tracking)
	OnChainGTKBalance int64      `gorm:"default:0" json:"on_chain_gtk_balance"`   // Last synced blockchain balance
	LastSyncedAt      *time.Time `json:"last_synced_at"`

	// Inventory
	Characters []Character     `gorm:"foreignKey:OwnerID" json:"-"`
	Teams      []Team          `gorm:"foreignKey:UserID" json:"-"`
	Inventory  []UserInventory `gorm:"foreignKey:UserID" json:"-"`

	// Tutorial Progress
	TutorialCompleted bool `gorm:"default:false" json:"tutorial_completed"`
	StoryProgress     int  `gorm:"default:0" json:"story_progress"`

	// Metadata
	LastKnownIP string `gorm:"type:varchar(45)" json:"-"` // For anti-cheat
	IsBanned    bool   `gorm:"default:false" json:"-"`
}
