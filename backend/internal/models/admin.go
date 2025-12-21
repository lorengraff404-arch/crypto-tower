package models

import (
	"time"
)

// RevenueTransaction tracks all GTK revenue from various sources
type RevenueTransaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Source    string    `gorm:"size:50;not null;index" json:"source"` // shop, gacha, marketplace, breeding, battle_wager
	AmountGTK float64   `gorm:"type:decimal(20,2);not null" json:"amount_gtk"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"`
	Metadata  string    `gorm:"type:jsonb" json:"metadata,omitempty"` // Additional data as JSON
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AdminAction logs all admin actions for audit trail
type AdminAction struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AdminWallet  string    `gorm:"size:42;not null;index" json:"admin_wallet"`
	ActionType   string    `gorm:"size:50;not null;index" json:"action_type"` // ban_user, adjust_balance, resolve_flag, update_config
	TargetUserID *uint     `gorm:"index" json:"target_user_id,omitempty"`
	Details      string    `gorm:"type:jsonb" json:"details,omitempty"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`

	// Relationships
	TargetUser *User `gorm:"foreignKey:TargetUserID" json:"target_user,omitempty"`
}

// AntiCheatFlag represents a detected suspicious activity
type AntiCheatFlag struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	FlagType        string     `gorm:"size:50;not null;index" json:"flag_type"` // impossible_stats, rapid_completion, suspicious_winrate, invalid_transaction
	Severity        string     `gorm:"size:20;not null;index" json:"severity"`  // low, medium, high, critical
	BattleID        *uint      `gorm:"index" json:"battle_id,omitempty"`
	Details         string     `gorm:"type:jsonb;not null" json:"details"`
	Status          string     `gorm:"size:20;default:'pending';index" json:"status"` // pending, reviewing, resolved, false_positive
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy      string     `gorm:"size:42" json:"resolved_by,omitempty"`
	ResolutionNotes string     `gorm:"type:text" json:"resolution_notes,omitempty"`
	CreatedAt       time.Time  `gorm:"index" json:"created_at"`

	// Relationships
	User   User    `gorm:"foreignKey:UserID" json:"user"`
	Battle *Battle `gorm:"foreignKey:BattleID" json:"battle,omitempty"`
}

// GameConfig stores centralized game configuration
type GameConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ConfigKey   string    `gorm:"size:100;unique;not null;index" json:"config_key"`
	ConfigValue string    `gorm:"type:jsonb;not null" json:"config_value"`
	Category    string    `gorm:"size:50;not null;index" json:"category"` // gacha, shop, battle, breeding, sprite
	Description string    `gorm:"type:text" json:"description,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `gorm:"size:42" json:"updated_by,omitempty"`
}
