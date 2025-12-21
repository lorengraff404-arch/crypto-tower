package models

import (
	"time"
)

// Mission represents a game mission/quest
type Mission struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Level          int       `gorm:"unique;not null" json:"level"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	Description    string    `json:"description"`
	Story          string    `json:"story"`
	MissionType    string    `gorm:"size:50;not null" json:"mission_type"` // tutorial, progression, daily, weekly
	UnlockFeature  string    `gorm:"size:50" json:"unlock_feature,omitempty"` // breeding, crafting, raids, ranked, advanced_crafting
	Objectives     string    `gorm:"type:jsonb" json:"objectives"` // JSONB array of objectives
	Rewards        string    `gorm:"type:jsonb" json:"rewards"` // JSONB object for rewards
	RequiredLevel  int       `gorm:"default:0" json:"required_level"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

// UserMissionProgress tracks user progress on missions
type UserMissionProgress struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	UserID             uint       `gorm:"not null;index:idx_user_mission" json:"user_id"`
	MissionID          uint       `gorm:"not null;index:idx_user_mission" json:"mission_id"`
	Status             string     `gorm:"size:20;default:'locked';index" json:"status"` // locked, available, in_progress, completed
	ObjectivesProgress string     `gorm:"type:jsonb" json:"objectives_progress"` // Current progress per objective
	StartedAt          *time.Time `json:"started_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Mission Mission `gorm:"foreignKey:MissionID" json:"mission,omitempty"`
}

// TableName specifies the table name for Mission
func (Mission) TableName() string {
	return "missions"
}

// TableName specifies the table name for UserMissionProgress
func (UserMissionProgress) TableName() string {
	return "user_mission_progress"
}

// Objective represents a mission objective
type Objective struct {
	Type        string `json:"type"`        // battle_waves, deploy_units, cast_spells, hatch_egg, breeding, craft, etc
	Target      int    `json:"target"`      // Target value
	Current     int    `json:"current"`     // Current progress
	Description string `json:"description"` // Human-readable description
}

// MissionRewards represents rewards for completing a mission
type MissionRewards struct {
	GTK   int64        `json:"gtk"`
	Items []RewardItem `json:"items,omitempty"`
}

// RewardItem represents a single reward item
type RewardItem struct {
	Type     string `json:"type"`     // material, consumable, weapon, armor, etc
	Name     string `json:"name"`     // Item name
	Quantity int    `json:"quantity"` // Amount
	Rarity   string `json:"rarity,omitempty"`
}
