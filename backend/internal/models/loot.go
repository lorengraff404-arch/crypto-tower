package models

import "time"

// LootDrop represents a single dropped item (Phase 16)
type LootDrop struct {
	ItemID   uint   `json:"item_id"`
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
	Rarity   string `json:"rarity"`
}

// Equipment represents equippable items (Phase 18)
type Equipment struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	ItemID        uint   `gorm:"not null" json:"item_id"`
	Slot          string `gorm:"size:20;not null" json:"slot"` // weapon, armor, accessory
	RequiredLevel int    `gorm:"default:1" json:"required_level"`
	RequiredClass string `gorm:"size:20" json:"required_class"`

	// Stat Bonuses
	BonusAttack  int `gorm:"default:0" json:"bonus_attack"`
	BonusDefense int `gorm:"default:0" json:"bonus_defense"`
	BonusHP      int `gorm:"default:0" json:"bonus_hp"`
	BonusSpeed   int `gorm:"default:0" json:"bonus_speed"`

	// Special Effects
	SpecialEffect string `gorm:"size:50" json:"special_effect"`
	EffectValue   int    `gorm:"default:0" json:"effect_value"`

	// Upgrade System
	UpgradeLevel    int `gorm:"default:0" json:"upgrade_level"`
	MaxUpgradeLevel int `gorm:"default:5" json:"max_upgrade_level"`

	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Item Item `gorm:"foreignKey:ItemID" json:"item"`
}

// CharacterEquipment tracks equipped items
type CharacterEquipment struct {
	ID          uint  `gorm:"primaryKey" json:"id"`
	CharacterID uint  `gorm:"not null;uniqueIndex" json:"character_id"`
	WeaponID    *uint `json:"weapon_id"`
	ArmorID     *uint `json:"armor_id"`
	AccessoryID *uint `json:"accessory_id"`

	// Relationships
	Character Character  `gorm:"foreignKey:CharacterID" json:"-"`
	Weapon    *Equipment `gorm:"foreignKey:WeaponID" json:"weapon,omitempty"`
	Armor     *Equipment `gorm:"foreignKey:ArmorID" json:"armor,omitempty"`
	Accessory *Equipment `gorm:"foreignKey:AccessoryID" json:"accessory,omitempty"`
}
