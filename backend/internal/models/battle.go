package models

import (
	"time"

	"gorm.io/gorm"
)

// Battle represents a PvP or PvE battle session
type Battle struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Battle Type
	BattleType string `gorm:"type:varchar(20);not null;index" json:"battle_type"` // PVP, PVE_ISLAND, PVE_TUTORIAL, RANKED, WAGER
	Status     string `gorm:"type:varchar(20);not null;index" json:"status"`      // PENDING, IN_PROGRESS, COMPLETED, SURRENDERED
	Seed       string `gorm:"type:varchar(64)" json:"seed"`                       // Deterministic battle seed for replay

	// Players (for PvP)
	Player1ID uint `gorm:"index" json:"player1_id"`
	Player1   User `gorm:"foreignKey:Player1ID" json:"-"`
	Player2ID uint `gorm:"index" json:"player2_id"`
	Player2   User `gorm:"foreignKey:Player2ID" json:"-"`

	// Winner
	WinnerID *uint `gorm:"index" json:"winner_id,omitempty"`
	Winner   *User `gorm:"foreignKey:WinnerID" json:"-"`

	// Betting (for PvP)
	BetAmount       int64 `gorm:"default:0" json:"bet_amount"` // Deprecated or used as BaseStake
	Player1Bet      int64 `gorm:"default:0" json:"player1_bet"`
	Player2Bet      int64 `gorm:"default:0" json:"player2_bet"`
	WinnerPayout    int64 `gorm:"default:0" json:"winner_payout"`
	PlatformRevenue int64 `gorm:"default:0" json:"platform_revenue"`

	// Battle State & Logs
	CurrentTurnPlayerID uint       `json:"current_turn_player_id"`
	TurnNumber          int        `gorm:"default:0" json:"turn_number"`
	ActionLog           string     `gorm:"type:text" json:"-"` // JSON log of all actions
	LastTurnData        string     `gorm:"type:text" json:"last_turn_data"`
	PlayerStateP1       string     `gorm:"type:text" json:"-"` // Serialized P1 team state
	PlayerStateP2       string     `gorm:"type:text" json:"-"` // Serialized P2 team state
	EndedAt             *time.Time `json:"ended_at"`

	// Performance Metrics
	DurationSeconds int `gorm:"default:0" json:"duration_seconds"`
}
