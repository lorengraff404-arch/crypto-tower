package models

import (
	"time"

	"gorm.io/gorm"
)

// StatusEffect represents a buff or debuff applied to a character
type StatusEffect struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Effect details
	EffectType string `gorm:"size:20;not null;index" json:"effect_type"` // BUFF or DEBUFF
	EffectName string `gorm:"size:30;not null;index" json:"effect_name"` // AMPED, BURN, STUN, etc.

	// Target
	CharacterID uint      `gorm:"not null;index" json:"character_id"`
	Character   Character `gorm:"foreignKey:CharacterID" json:"-"`

	// Stacking
	Stacks int `gorm:"default:1;not null" json:"stacks"` // 1-3 stacks allowed

	// Duration
	Duration       int       `gorm:"not null" json:"duration"`         // Turns/seconds remaining
	MaxDuration    int       `gorm:"not null" json:"max_duration"`     // Original duration
	TurnsRemaining int       `gorm:"default:0" json:"turns_remaining"` // For turn-based effects
	ExpiresAt      time.Time `json:"expires_at"`

	// Source
	SourceAbilityID *uint `json:"source_ability_id,omitempty"` // Which ability caused it
	CasterID        *uint `json:"caster_id,omitempty"`         // Who cast it
	BattleID        *uint `json:"battle_id,omitempty"`         // Which battle

	// Effect values
	StatModifier  float64 `json:"stat_modifier"`   // e.g., 0.3 for +30% or -0.3 for -30%
	DamagePerTurn int     `json:"damage_per_turn"` // For DoT effects
}

// StatusEffectDefinition holds the definitions for all effects
type StatusEffectDefinition struct {
	Name            string
	Type            string // BUFF or DEBUFF
	Icon            string // Emoji
	Description     string
	Modifier        float64 // Stat modifier
	DamagePerTurn   int     // DoT damage
	DefaultDuration int     // Default turns
}

var BuffDefinitions = map[string]StatusEffectDefinition{
	"AMPED": {
		Name:            "Amped",
		Type:            "BUFF",
		Icon:            "⚡",
		Description:     "+30% Attack",
		Modifier:        0.3,
		DefaultDuration: 3,
	},
	"BULKED": {
		Name:            "Bulked",
		Type:            "BUFF",
		Icon:            "🛡️",
		Description:     "+30% Defense",
		Modifier:        0.3,
		DefaultDuration: 3,
	},
	"HASTE": {
		Name:            "Haste",
		Type:            "BUFF",
		Icon:            "💨",
		Description:     "+30% Speed",
		Modifier:        0.3,
		DefaultDuration: 3,
	},
	"FLEET": {
		Name:            "Fleet",
		Type:            "BUFF",
		Icon:            "🦅",
		Description:     "+30% Evasion",
		Modifier:        0.3,
		DefaultDuration: 3,
	},
	"WARDED": {
		Name:            "Warded",
		Type:            "BUFF",
		Icon:            "✨",
		Description:     "+30% Resistance",
		Modifier:        0.3,
		DefaultDuration: 3,
	},
	"REGEN": {
		Name:            "Regen",
		Type:            "BUFF",
		Icon:            "💚",
		Description:     "+5% HP per turn",
		Modifier:        0.05,
		DefaultDuration: 5,
	},
}

var DebuffDefinitions = map[string]StatusEffectDefinition{
	"BURN": {
		Name:            "Burn",
		Type:            "DEBUFF",
		Icon:            "🔥",
		Description:     "5% max HP damage per turn",
		DamagePerTurn:   5, // Percentage
		DefaultDuration: 3,
	},
	"POISON": {
		Name:            "Poison",
		Type:            "DEBUFF",
		Icon:            "☠️",
		Description:     "8% max HP damage per turn",
		DamagePerTurn:   8,
		DefaultDuration: 3,
	},
	"BLEED": {
		Name:            "Bleed",
		Type:            "DEBUFF",
		Icon:            "🩸",
		Description:     "3% max HP damage per turn",
		DamagePerTurn:   3,
		DefaultDuration: 4,
	},
	"SLOW": {
		Name:            "Slow",
		Type:            "DEBUFF",
		Icon:            "🐌",
		Description:     "-30% Speed",
		Modifier:        -0.3,
		DefaultDuration: 3,
	},
	"FEEBLE": {
		Name:            "Feeble",
		Type:            "DEBUFF",
		Icon:            "😰",
		Description:     "-30% Attack",
		Modifier:        -0.3,
		DefaultDuration: 3,
	},
	"FRAGILE": {
		Name:            "Fragile",
		Type:            "DEBUFF",
		Icon:            "💔",
		Description:     "-30% Defense",
		Modifier:        -0.3,
		DefaultDuration: 3,
	},
	"STUN": {
		Name:            "Stun",
		Type:            "DEBUFF",
		Icon:            "💫",
		Description:     "Cannot act",
		DefaultDuration: 1,
	},
	"SLEEP": {
		Name:            "Sleep",
		Type:            "DEBUFF",
		Icon:            "😴",
		Description:     "Cannot act, wakes on damage",
		DefaultDuration: 2,
	},
	"FREEZE": {
		Name:            "Freeze",
		Type:            "DEBUFF",
		Icon:            "🧊",
		Description:     "Cannot act, 2x fire damage",
		DefaultDuration: 1,
	},
	"PARALYZE": {
		Name:            "Paralyze",
		Type:            "DEBUFF",
		Icon:            "⚡",
		Description:     "25% chance to skip turn",
		DefaultDuration: 3,
	},
	"SILENCE": {
		Name:            "Silence",
		Type:            "DEBUFF",
		Icon:            "🤐",
		Description:     "Cannot use abilities",
		DefaultDuration: 2,
	},
	"BLIND": {
		Name:            "Blind",
		Type:            "DEBUFF",
		Icon:            "🌫️",
		Description:     "-50% Hit chance",
		Modifier:        -0.5,
		DefaultDuration: 2,
	},
}

// GetEffectDefinition returns the definition for an effect
func GetEffectDefinition(effectName string) (StatusEffectDefinition, bool) {
	if def, ok := BuffDefinitions[effectName]; ok {
		return def, true
	}
	if def, ok := DebuffDefinitions[effectName]; ok {
		return def, true
	}
	return StatusEffectDefinition{}, false
}
