package models

import (
	"time"

	"gorm.io/gorm"
)

// AbilityLearning defines when a character can learn an ability
// Based on their Rank (rarity) and Level
type AbilityLearning struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Ability Reference
	AbilityID uint    `gorm:"not null;index" json:"ability_id"`
	Ability   Ability `gorm:"foreignKey:AbilityID" json:"ability,omitempty"`

	// Learning Requirements
	MinRank    string `gorm:"size:3;not null;index" json:"min_rank"` // C, B, A, S, SS, SSS
	LearnLevel int    `gorm:"not null;index" json:"learn_level"`     // Level required to learn
	IsStarting bool   `gorm:"default:false" json:"is_starting"`      // Given at character creation?
	IsUltimate bool   `gorm:"default:false" json:"is_ultimate"`      // Ultimate ability (special unlock)

	// Optional: Prerequisite ability
	PrerequisiteAbilityID *uint `json:"prerequisite_ability_id,omitempty"` // Must know this ability first
}

// CharacterLearnedAbility tracks which abilities a character has learned
// Separate from equipped abilities
type CharacterLearnedAbility struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CharacterID uint `gorm:"not null;index:idx_char_ability,unique" json:"character_id"`
	AbilityID   uint `gorm:"not null;index:idx_char_ability,unique" json:"ability_id"`

	LearnedAt time.Time `json:"learned_at"`
	TimesUsed int       `gorm:"default:0" json:"times_used"`

	// Relationships
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
	Ability   Ability   `gorm:"foreignKey:AbilityID" json:"ability,omitempty"`
}
