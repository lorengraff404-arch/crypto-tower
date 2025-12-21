package models

import "time"

// AuditLog tracks all important actions for security
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     *uint     `gorm:"index" json:"user_id"`
	Action     string    `gorm:"size:50;not null" json:"action"`
	EntityType string    `gorm:"size:50;not null" json:"entity_type"`
	EntityID   *uint     `json:"entity_id"`
	OldValues  string    `gorm:"type:text" json:"old_values"`
	NewValues  string    `gorm:"type:text" json:"new_values"`
	IPAddress  string    `gorm:"size:45" json:"ip_address"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// RateLimit prevents abuse
type RateLimit struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	ActionType  string    `gorm:"size:50;not null" json:"action_type"`
	Count       int       `gorm:"default:1" json:"count"`
	WindowStart time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"window_start"`
}

// Notification for user alerts
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Type      string    `gorm:"size:30;not null" json:"type"`
	Title     string    `gorm:"size:100" json:"title"`
	Message   string    `gorm:"type:text" json:"message"`
	Data      string    `gorm:"type:text" json:"data"`
	Read      bool      `gorm:"default:false;index" json:"read"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// Leaderboard for rankings
type Leaderboard struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Category  string    `gorm:"size:30;not null;index" json:"category"`
	Score     int       `gorm:"not null" json:"score"`
	Rank      int       `gorm:"index" json:"rank"`
	Season    int       `gorm:"default:1" json:"season"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Friendship for social features
type Friendship struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	FriendID  uint      `gorm:"not null;index" json:"friend_id"`
	Status    string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`

	User   User `gorm:"foreignKey:UserID" json:"-"`
	Friend User `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

// Referral for referral system
type Referral struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ReferrerID    uint      `gorm:"not null" json:"referrer_id"`
	ReferredID    uint      `gorm:"not null" json:"referred_id"`
	RewardClaimed bool      `gorm:"default:false" json:"reward_claimed"`
	CreatedAt     time.Time `json:"created_at"`

	Referrer User `gorm:"foreignKey:ReferrerID" json:"referrer,omitempty"`
	Referred User `gorm:"foreignKey:ReferredID" json:"referred,omitempty"`
}

// BattleReplay stores battle data for replay
type BattleReplay struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RaidSessionID uint      `gorm:"not null;index" json:"raid_session_id"`
	ReplayData    string    `gorm:"type:text;not null" json:"replay_data"`
	CreatedAt     time.Time `json:"created_at"`
}
