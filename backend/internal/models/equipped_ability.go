package models

import "time"

// EquippedAbility represents abilities that are actively equipped for battle
// Separate from CharacterAbility which tracks all learned abilities
type EquippedAbility struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CharacterID  uint      `gorm:"not null;index:idx_equipped_character" json:"character_id"`
	AbilityID    uint      `gorm:"not null" json:"ability_id"`
	SlotPosition int       `gorm:"not null;check:slot_position >= 1 AND slot_position <= 16" json:"slot_position"`
	EquippedAt   time.Time `gorm:"default:now()" json:"equipped_at"`

	// Relationships
	Character Character `gorm:"foreignKey:CharacterID" json:"-"`
	Ability   Ability   `gorm:"foreignKey:AbilityID" json:"ability"`
}

// TableName specifies the table name
func (EquippedAbility) TableName() string {
	return "equipped_abilities"
}
