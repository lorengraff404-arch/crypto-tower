package models

import (
	"time"

	"gorm.io/gorm"
)

// Team represents a user's battle team configuration
type Team struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	Name      string `gorm:"size:50;not null;default:'My Team'" json:"name"`
	IsDefault bool   `gorm:"default:false" json:"is_default"`

	// Members
	Members []TeamMember `gorm:"foreignKey:TeamID" json:"members"`

	// Calculated Stats (cached)
	TotalPower int `json:"total_power"`
	AvgLevel   int `json:"avg_level"`

	// Synergies (calculated on load)
	Synergies []TeamSynergy `gorm:"-" json:"synergies"`
}

// TeamMember links a character to a team with a specific position
type TeamMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	TeamID uint `gorm:"not null;index" json:"team_id"`

	CharacterID uint      `gorm:"not null;index" json:"character_id"`
	Character   Character `gorm:"foreignKey:CharacterID" json:"character"`

	// Position logic
	// Slot 0-2: Active Team (Front, Middle, Back)
	// Slot 3-5: Backup Team
	Slot int `gorm:"not null" json:"slot"`

	IsBackup bool `gorm:"default:false" json:"is_backup"`
}

// TeamSynergy represents an active bonus for the team
type TeamSynergy struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Tier        int     `json:"tier"` // 1 (Bronze), 2 (Silver), 3 (Gold)
	Icon        string  `json:"icon"`
	Condition   string  `json:"condition"`
	BonusStat   string  `json:"bonus_stat"`
	BonusValue  float64 `json:"bonus_value"`
}
