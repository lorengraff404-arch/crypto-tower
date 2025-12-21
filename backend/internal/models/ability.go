package models

import (
	"time"

	"github.com/lib/pq"
)

// Ability represents a learnable skill for characters
type Ability struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:50;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Class       string `gorm:"size:20;not null;index" json:"class"` // Warrior, Mage, Tank
	UnlockLevel int    `gorm:"not null" json:"unlock_level"`
	AbilityType string `gorm:"size:20;not null" json:"ability_type"` // ACTIVE, PASSIVE, ULTIMATE
	Damage      int    `gorm:"default:0" json:"damage"`              // NEW: Base damage of ability
	Element     string `gorm:"size:20" json:"element"`               // NEW: Fire, Water, Grass, etc.
	Category    string `gorm:"size:20" json:"category"`              // physical, special, status
	Accuracy    int    `gorm:"default:100" json:"accuracy"`          // Hit chance %
	Priority    int    `gorm:"default:0" json:"priority"`            // Higher goes first
	TargetType  string `gorm:"size:20" json:"target_type"`           // single, all_opponents, self, allySINGLE_ENEMY, AOE, ALL_ALLIES
	Cooldown    int    `json:"cooldown"`                             // Seconds
	ManaCost    int    `gorm:"not null" json:"mana_cost"`
	MaxPP       int    `gorm:"default:10" json:"max_pp"` // Maximum PP for this ability

	// Effect Values (base, modified by element)
	BaseDamage   int `json:"base_damage"`
	BaseHeal     int `json:"base_heal"`
	DurationSecs int `json:"duration_secs"`
	EffectPower  int `json:"effect_power"` // Generic power value

	// Status Effects
	AppliesBuff   string `gorm:"size:50" json:"applies_buff,omitempty"`   // Buff name
	AppliesDebuff string `gorm:"size:50" json:"applies_debuff,omitempty"` // Debuff name

	// Visual & Audio
	IconURL       string `gorm:"size:255" json:"icon_url"`
	AnimationName string `gorm:"size:50" json:"animation_name"` // VFX identifier
	SoundEffect   string `gorm:"size:100" json:"sound_effect"`

	// Element Bonuses (JSON)
	ElementBonuses string `gorm:"type:jsonb" json:"element_bonuses"` // {"Fire": 1.2, "Water": 0.8}

	// Rarity & Requirements (Phase 20)
	Rarity          string         `gorm:"size:10;default:'C';index" json:"rarity"` // C, B, A, S, SS, SSS
	IsUltimate      bool           `gorm:"default:false" json:"is_ultimate"`
	SynergyTags     pq.StringArray `gorm:"type:text[]" json:"synergy_tags"`     // For combo detection
	RequiredElement pq.StringArray `gorm:"type:text[]" json:"required_element"` // NULL = any element
	RequiredClass   pq.StringArray `gorm:"type:text[]" json:"required_class"`   // NULL = any class

	// Battle Mechanics (Phase 20)
	DamageType         string `gorm:"size:20;default:'physical'" json:"damage_type"` // physical, magical, true
	AOERadius          int    `gorm:"default:0" json:"aoe_radius"`                   // 0 = single target
	StatusEffectChance int    `gorm:"default:0" json:"status_effect_chance"`         // percentage
	BuffDuration       int    `gorm:"default:0" json:"buff_duration"`                // turns
	DebuffDuration     int    `gorm:"default:0" json:"debuff_duration"`              // turns

	CreatedAt time.Time `json:"created_at"`
}

// CharacterAbility tracks which abilities a character has learned
type CharacterAbility struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CharacterID uint      `gorm:"not null;index:idx_char_ability" json:"character_id"`
	AbilityID   uint      `gorm:"not null;index:idx_char_ability" json:"ability_id"`
	LearnedAt   time.Time `json:"learned_at"`
	TimesUsed   int       `gorm:"default:0" json:"times_used"`

	// Relationships
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
	Ability   Ability   `gorm:"foreignKey:AbilityID" json:"ability"`
}

// TableName specifies the table name for Ability
func (Ability) TableName() string {
	return "abilities"
}

// TableName specifies the table name for CharacterAbility
func (CharacterAbility) TableName() string {
	return "character_abilities"
}

// ElementBonus represents element-specific ability modifiers
type ElementBonus struct {
	Fire    float64 `json:"Fire"`
	Water   float64 `json:"Water"`
	Ice     float64 `json:"Ice"`
	Thunder float64 `json:"Thunder"`
	Dark    float64 `json:"Dark"`
	Plant   float64 `json:"Plant"`
	Earth   float64 `json:"Earth"`
	Wind    float64 `json:"Wind"`
}

// AbilityUsage tracks PP usage for character abilities
type AbilityUsage struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	CharacterID     uint       `gorm:"not null;index" json:"character_id"`
	AbilityID       uint       `gorm:"not null;index" json:"ability_id"`
	CurrentPP       int        `gorm:"not null" json:"current_pp"`
	MaxPP           int        `gorm:"not null" json:"max_pp"`
	LastUsedAt      *time.Time `json:"last_used_at"`
	LastRecoveredAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"last_recovered_at"`

	// Relationships
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
	Ability   Ability   `gorm:"foreignKey:AbilityID" json:"ability"`
}

// CharacterStatusEffect represents active status effects on characters
type CharacterStatusEffect struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	CharacterID   uint       `gorm:"not null;index" json:"character_id"`
	StatusType    string     `gorm:"size:20;not null" json:"status_type"` // poison, burn, freeze, paralysis, sleep
	AppliedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"applied_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	DamagePerTurn int        `gorm:"default:0" json:"damage_per_turn"`

	// Relationships
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
}

// CharacterBuff represents active buffs on characters
type CharacterBuff struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CharacterID    uint      `gorm:"not null;index" json:"character_id"`
	BuffType       string    `gorm:"size:20;not null" json:"buff_type"`            // x_attack, x_defense, x_speed, etc
	Multiplier     float64   `gorm:"type:decimal(3,2);not null" json:"multiplier"` // 1.5 for +50%
	AppliedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"applied_at"`
	TurnsRemaining int       `gorm:"not null" json:"turns_remaining"`

	// Relationships
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
}

// CharacterActiveSkill tracks which skills are equipped in battle slots
type CharacterActiveSkill struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CharacterID  uint      `gorm:"not null;index" json:"character_id"`
	AbilityID    uint      `gorm:"not null;index" json:"ability_id"`
	SlotPosition int       `gorm:"not null" json:"slot_position"` // 1-7 based on rarity
	IsLocked     bool      `gorm:"default:true" json:"is_locked"`
	UnlockLevel  int       `gorm:"default:1" json:"unlock_level"`
	CreatedAt    time.Time `json:"created_at"`

	// Relationships
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
	Ability   Ability   `gorm:"foreignKey:AbilityID" json:"ability"`
}

// CharacterSkillCooldown tracks cooldown state for character abilities
type CharacterSkillCooldown struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	CharacterID       uint       `gorm:"not null;index" json:"character_id"`
	AbilityID         uint       `gorm:"not null;index" json:"ability_id"`
	CooldownRemaining int        `gorm:"default:0" json:"cooldown_remaining"` // Turns remaining
	LastUsedAt        *time.Time `json:"last_used_at"`

	// Relationships
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
	Ability   Ability   `gorm:"foreignKey:AbilityID" json:"ability"`
}

// TableName specifies the table name for CharacterActiveSkill
func (CharacterActiveSkill) TableName() string {
	return "character_active_skills"
}

// TableName specifies the table name for CharacterSkillCooldown
func (CharacterSkillCooldown) TableName() string {
	return "character_skill_cooldowns"
}
