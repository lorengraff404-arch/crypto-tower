package models

import (
	"time"
)

// SystemSetting represents a dynamic configuration value (Live Ops)
type SystemSetting struct {
	Key         string    `gorm:"primaryKey;type:varchar(50)" json:"key"`
	Value       string    `gorm:"type:text;not null" json:"value"` // Stored as string, casted in app
	Type        string    `gorm:"type:varchar(20)" json:"type"`    // string, int, float, bool, json
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   uint      `json:"updated_by"` // Admin User ID
}

// AuditLog records all administrative actions for security
type AdminAuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AdminID   uint      `gorm:"index;not null" json:"admin_id"`
	Action    string    `gorm:"type:varchar(50);not null" json:"action"` // BAN_USER, UPDATE_CONFIG, REFUND
	TargetID  string    `gorm:"index" json:"target_id"`                  // UserID or ConfigKey
	OldValue  string    `gorm:"type:text" json:"old_value"`
	NewValue  string    `gorm:"type:text" json:"new_value"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// UserReport represents a report filed by a player against another
type UserReport struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	ReporterID uint       `gorm:"index;not null" json:"reporter_id"`
	TargetID   uint       `gorm:"index;not null" json:"target_id"`
	Reason     string     `gorm:"type:varchar(50);not null" json:"reason"` // CHEATING, TOXICITY, AFK
	Details    string     `gorm:"type:text" json:"details"`
	Evidence   string     `gorm:"type:text" json:"evidence"`                        // URL or Hash
	Status     string     `gorm:"type:varchar(20);default:'PENDING'" json:"status"` // PENDING, RESOLVED, DISMISSED
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at"`
	ResolvedBy *uint      `json:"resolved_by"`
}
