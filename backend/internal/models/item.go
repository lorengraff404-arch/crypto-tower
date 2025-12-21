package models

import (
	"time"

	"gorm.io/gorm"
)

// Item represents equipment, consumables, or materials
type Item struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Ownership
	OwnerID    uint   `gorm:"not null;index" json:"owner_id"`
	Owner      User   `gorm:"foreignKey:OwnerID" json:"-"`
	NFTTokenID string `gorm:"uniqueIndex" json:"nft_token_id,omitempty"` // For rare items
	IsMinted   bool   `gorm:"default:false" json:"is_minted"`

	// Item Identity
	ItemType string `gorm:"type:varchar(20);not null;index" json:"item_type"` // WEAPON, ARMOR, ACCESSORY, RUNE, CONSUMABLE, MATERIAL
	Name     string `gorm:"type:varchar(50);not null" json:"name"`
	Rarity   string `gorm:"type:varchar(5);not null;index" json:"rarity"` // SSS, SS, S, A, B, C
	IconURL  string `gorm:"type:varchar(255)" json:"icon_url,omitempty"`

	// Stat Modifiers (for equipment)
	AttackBonus  int `gorm:"default:0" json:"attack_bonus"`
	DefenseBonus int `gorm:"default:0" json:"defense_bonus"`
	HPBonus      int `gorm:"default:0" json:"hp_bonus"`
	SpeedBonus   int `gorm:"default:0" json:"speed_bonus"`

	// Special Effects (JSON array)
	SpecialEffects string `gorm:"type:text" json:"special_effects,omitempty"` // e.g., "+20% Fire Damage"

	// Durability (for equipment)
	Durability int  `gorm:"default:100" json:"durability"` // 0-100
	IsBroken   bool `gorm:"default:false" json:"is_broken"`

	// Consumable Properties
	IsConsumable  bool   `gorm:"default:false" json:"is_consumable"`
	ConsumeEffect string `gorm:"type:varchar(50)" json:"consume_effect,omitempty"` // REVIVE, REDUCE_FATIGUE, XP_BOOST

	// Crafting Material
	IsCraftingMaterial bool `gorm:"default:false" json:"is_crafting_material"`

	// Stackable (for materials/consumables)
	IsStackable bool `gorm:"default:false" json:"is_stackable"`
	Quantity    int  `gorm:"default:1;not null" json:"quantity"`

	// Marketplace
	IsListed  bool       `gorm:"default:false;index" json:"is_listed"`
	ListPrice int64      `json:"list_price,omitempty"`
	ListedAt  *time.Time `json:"listed_at,omitempty"`

	// Equipment Status
	IsEquipped   bool  `gorm:"default:false;index" json:"is_equipped"`
	EquippedByID *uint `gorm:"index" json:"equipped_by_id,omitempty"` // Character ID
}

// UserInventory represents items owned by a user
type UserInventory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	ItemID     uint      `gorm:"not null;index" json:"item_id"`
	Quantity   int       `gorm:"default:1;check:quantity >= 0" json:"quantity"`
	AcquiredAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"acquired_at"`

	// Relationships
	User User     `gorm:"foreignKey:UserID" json:"-"`
	Item ShopItem `gorm:"foreignKey:ItemID" json:"item"` // Changed from Item to ShopItem
}

// ShopItem represents an item available in the shop
type ShopItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:50;not null;index" json:"name"` // Changed from uniqueIndex to index
	Description  string    `gorm:"type:text" json:"description"`
	Category     string    `gorm:"size:20;not null;index" json:"category"` // healing, status, pp, battle, egg, evolution
	EffectType   string    `gorm:"size:20;not null" json:"effect_type"`    // heal_hp, cure_status, restore_pp, buff, accelerate, etc
	EffectValue  int       `json:"effect_value"`                           // Amount to heal/restore/reduce
	GTKCost      int64     `gorm:"not null" json:"gtk_cost"`
	IsConsumable bool      `gorm:"default:true" json:"is_consumable"`
	MaxStack     int       `gorm:"default:99" json:"max_stack"`
	IsAvailable  bool      `gorm:"default:true;index" json:"is_available"`
	IconURL      string    `gorm:"size:200" json:"icon_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LootTable defines what can drop from a mission (Phase 16)
type LootTable struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	MissionID *uint     `json:"mission_id"`
	DropType  string    `gorm:"size:20;default:'random'" json:"drop_type"`
	CreatedAt time.Time `json:"created_at"`

	Entries []LootEntry `gorm:"foreignKey:LootTableID" json:"entries"`
}

// LootEntry represents a single item in a loot table (Phase 16)
type LootEntry struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	LootTableID  uint    `gorm:"not null" json:"loot_table_id"`
	ItemID       *uint   `json:"item_id"`
	DropChance   float64 `gorm:"type:decimal(5,2);not null" json:"drop_chance"`
	MinQuantity  int     `gorm:"default:1" json:"min_quantity"`
	MaxQuantity  int     `gorm:"default:1" json:"max_quantity"`
	RarityWeight int     `gorm:"default:1" json:"rarity_weight"`

	LootTable LootTable `gorm:"foreignKey:LootTableID" json:"-"`
	Item      *Item     `gorm:"foreignKey:ItemID" json:"item"`
}

// BattleReward tracks rewards earned from battles (Phase 16)
type BattleReward struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	RaidSessionID    uint      `gorm:"not null;index" json:"raid_session_id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	TokensEarned     int       `gorm:"default:0" json:"tokens_earned"`
	XPEarned         int       `gorm:"default:0" json:"xp_earned"`
	ItemsJSON        string    `gorm:"type:text" json:"items_json"`
	PerformanceGrade string    `gorm:"size:2" json:"performance_grade"`
	CreatedAt        time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
