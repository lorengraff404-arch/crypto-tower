package models

import (
	"time"

	"gorm.io/gorm"
)

// RevenueDistribution represents a GTK revenue split
type RevenueDistribution struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	TxHash         string         `gorm:"uniqueIndex" json:"tx_hash"`
	Source         string         `gorm:"index" json:"source"`
	TotalAmount    float64        `json:"total_amount"`
	GrowthFund     float64        `json:"growth_fund"`
	SecurityFund   float64        `json:"security_fund"`
	Operations     float64        `json:"operations"`
	RewardsPool    float64        `json:"rewards_pool"`
	DevTeam        float64        `json:"dev_team"`
	TowerLiquidity float64        `json:"tower_liquidity"`
}
