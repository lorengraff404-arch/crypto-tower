package models

import "time"

// DailyQuest represents a daily quest for a user
type DailyQuest struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	UserID          uint       `json:"user_id" gorm:"not null;index"`
	QuestType       string     `json:"quest_type" gorm:"not null"` // combat, collection, progression, special
	QuestName       string     `json:"quest_name" gorm:"not null"`
	Description     string     `json:"description" gorm:"not null"`
	TargetValue     int        `json:"target_value" gorm:"not null"`
	CurrentProgress int        `json:"current_progress" gorm:"default:0"`
	Difficulty      string     `json:"difficulty" gorm:"not null"` // common, uncommon, rare, epic
	RewardItemID    *int       `json:"reward_item_id"`
	RewardGTK       int        `json:"reward_gtk" gorm:"default:0"`
	RewardTOWER     float64    `json:"reward_tower" gorm:"default:0"`
	IsCompleted     bool       `json:"is_completed" gorm:"default:false"`
	IsClaimed       bool       `json:"is_claimed" gorm:"default:false"`
	CompletedAt     *time.Time `json:"completed_at"`
	ClaimedAt       *time.Time `json:"claimed_at"`
	ExpiresAt       time.Time  `json:"expires_at" gorm:"not null;index"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// QuestProgressTracking tracks individual progress events
type QuestProgressTracking struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            uint      `json:"user_id" gorm:"not null;index"`
	QuestID           uint      `json:"quest_id" gorm:"not null;index"`
	ActionType        string    `json:"action_type" gorm:"not null"`
	ProgressIncrement int       `json:"progress_increment" gorm:"default:1"`
	Metadata          string    `json:"metadata" gorm:"type:jsonb"` // JSONB for additional context
	TrackedAt         time.Time `json:"tracked_at"`
}

// TableName specifies the table name for DailyQuest
func (DailyQuest) TableName() string {
	return "daily_quests"
}

// TableName specifies the table name for QuestProgressTracking
func (QuestProgressTracking) TableName() string {
	return "quest_progress_tracking"
}

// QuestTemplate defines a quest type that can be generated
type QuestTemplate struct {
	Type        string
	Name        string
	Description string
	BaseTarget  int
	ScaleFactor float64
	Difficulty  string
	ActionType  string // What action triggers progress
}
