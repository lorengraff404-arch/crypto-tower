package models

import (
	"time"

	"gorm.io/gorm"
)

// Character represents a collectible unit/guardian
type Character struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Ownership
	OwnerID    uint       `gorm:"not null;index" json:"owner_id"`
	Owner      User       `gorm:"foreignKey:OwnerID" json:"-"`
	NFTTokenID *string    `gorm:"uniqueIndex" json:"nft_token_id,omitempty"` // Minted on BSC (Pointer for NULL support)
	IsMinted   bool       `gorm:"default:false;index" json:"is_minted"`
	MintedAt   *time.Time `json:"minted_at,omitempty"`

	// NFT/Blockchain fields
	OnChainTokenID *uint64 `gorm:"uniqueIndex" json:"on_chain_token_id,omitempty"`
	MetadataURI    string  `gorm:"type:varchar(200)" json:"metadata_uri,omitempty"`
	MintTxHash     string  `gorm:"type:varchar(66)" json:"mint_tx_hash,omitempty"`

	// Character Identity
	Name          string  `gorm:"type:varchar(50)" json:"name"`
	UniqueName    *string `gorm:"type:varchar(100);uniqueIndex" json:"unique_name,omitempty"` // Generated unique name
	Class         string  `gorm:"size:20;not null;index" json:"class"`                        // Warrior, Mage, Tank
	Element       string  `gorm:"size:20;not null;index" json:"element"`                      // Fire, Water, Ice, etc.
	CharacterType string  `gorm:"size:20;default:'BEAST';index" json:"character_type"`        // BEAST, DRAGON, INSECT, etc.
	Rarity        string  `gorm:"size:10;not null;index" json:"rarity"`                       // C, B, A, S, SS, SSS
	// Base Stats (at level 1, rarity C)
	BaseAttack  int `gorm:"not null" json:"base_attack"`
	BaseDefense int `gorm:"not null" json:"base_defense"`
	BaseHP      int `gorm:"not null" json:"base_hp"`
	BaseSpeed   int `gorm:"not null" json:"base_speed"`

	// Current Stats (with level, rarity, equipment bonuses)
	CurrentAttack  int `gorm:"not null" json:"current_attack"`
	CurrentDefense int `gorm:"not null" json:"current_defense"`
	CurrentHP      int `gorm:"not null" json:"current_hp"`
	CurrentSpeed   int `gorm:"not null" json:"current_speed"`

	// Progression
	Level          int `gorm:"default:1;not null" json:"level"`
	Experience     int `gorm:"default:0;not null" json:"experience"` // Current level progress
	TotalXP        int `gorm:"default:0;not null" json:"total_xp"`   // Total XP earned (for level calculation)
	EvolutionStage int `gorm:"default:0" json:"evolution_stage"`     // 0-3 (Base, Growth, Mature, Ultimate)

	// Tier System (SSS to C)
	Tier        string `gorm:"size:3;default:'C'" json:"tier"`      // SSS, SS, S, A, B, C
	CombatPower int    `gorm:"default:0;index" json:"combat_power"` // For matchmaking
	ELORating   int    `gorm:"default:1000" json:"elo_rating"`      // For ranked PvP

	// Mana System (Phase 20 - Skill System)
	CurrentMana   int `gorm:"default:100;not null" json:"current_mana"`   // Current mana available
	MaxMana       int `gorm:"default:100;not null" json:"max_mana"`       // Maximum mana capacity (scales with rarity)
	ManaRegenRate int `gorm:"default:10;not null" json:"mana_regen_rate"` // Mana regenerated per turn (higher for higher rarity)

	// Equipped Abilities (4 Slots) - Like Pokémon
	EquippedAbility1 *uint `json:"equipped_ability_1,omitempty"`
	EquippedAbility2 *uint `json:"equipped_ability_2,omitempty"`
	EquippedAbility3 *uint `json:"equipped_ability_3,omitempty"`
	EquippedAbility4 *uint `json:"equipped_ability_4,omitempty"`

	// Passive Ability (Phase 15.2)
	PassiveAbility string `gorm:"size:50" json:"passive_ability"` // Berserker, Mana Surge, etc.

	// Breeding (Phase 17)
	BreedCount int        `gorm:"default:0" json:"breed_count"`
	LastBredAt *time.Time `json:"last_bred_at"`

	// Durability & Fatigue
	Durability   int        `gorm:"default:100;not null" json:"durability"` // 0-100
	Fatigue      int        `gorm:"default:0;not null" json:"fatigue"`      // 0-100
	LastBattleAt *time.Time `json:"last_battle_at,omitempty"`
	IsDead       bool       `gorm:"default:false;index" json:"is_dead"` // Durability reached 0
	CanBeRevived bool       `gorm:"default:true" json:"can_be_revived"`

	// Fainted status (for battle)
	IsFainted bool       `gorm:"default:false;index" json:"is_fainted"`
	FaintedAt *time.Time `json:"fainted_at,omitempty"`

	// Abilities (JSON array of ability IDs)
	Abilities         string `gorm:"type:text" json:"abilities"`           // JSON: ["ability_1", "ability_2"]
	UnlockedAbilities string `gorm:"type:jsonb" json:"unlocked_abilities"` // JSON array of unlocked ability IDs

	// Visual
	ImageURL      string `gorm:"type:varchar(255)" json:"image_url"`
	VisualTraits  string `gorm:"type:text" json:"visual_traits,omitempty"`  // Generated unique traits for prompt
	AnimationData string `gorm:"type:text" json:"animation_data,omitempty"` // JSON for animations

	// Sprite Sheets (for battle animations)
	SpriteIdle      string `gorm:"type:varchar(255)" json:"sprite_idle,omitempty"`
	SpriteWalk      string `gorm:"type:varchar(255)" json:"sprite_walk,omitempty"`
	SpriteRun       string `gorm:"type:varchar(255)" json:"sprite_run,omitempty"`
	SpriteAttack    string `gorm:"type:varchar(255)" json:"sprite_attack,omitempty"`
	SpriteSkill     string `gorm:"type:varchar(255)" json:"sprite_skill,omitempty"`
	SpriteHit       string `gorm:"type:varchar(255)" json:"sprite_hit,omitempty"`
	SpriteBlock     string `gorm:"type:varchar(255)" json:"sprite_block,omitempty"`
	SpriteDodge     string `gorm:"type:varchar(255)" json:"sprite_dodge,omitempty"`
	SpriteDeath     string `gorm:"type:varchar(255)" json:"sprite_death,omitempty"`
	SpriteVictory   string `gorm:"type:varchar(255)" json:"sprite_victory,omitempty"`
	SpriteGenStatus string `gorm:"size:20;default:'pending'" json:"sprite_gen_status"` // pending, processing, completed, failed
	SpriteGenJobID  *uint  `gorm:"index" json:"sprite_gen_job_id,omitempty"`

	// Relationships
	Moves []CharacterMove `gorm:"foreignKey:CharacterID" json:"moves,omitempty"` // Phase 10.2

	// Marketplace
	IsListed  bool       `gorm:"default:false;index" json:"is_listed"`
	ListPrice int64      `json:"list_price,omitempty"`
	ListedAt  *time.Time `json:"listed_at,omitempty"`

	// Incubation (for eggs/seeds)
	IsEgg             bool       `gorm:"default:false;index" json:"is_egg"` // True if not hatched yet
	HatchTime         *time.Time `json:"hatch_time,omitempty"`
	CareSlotCalibrate bool       `gorm:"default:false" json:"care_slot_calibrate"`
	CareSlotNurture   bool       `gorm:"default:false" json:"care_slot_nurture"`
	CareSlotStabilize bool       `gorm:"default:false" json:"care_slot_stabilize"`

	// Relationships
	EquippedItems []Item `gorm:"many2many:character_equipped_items;" json:"-"`
}

// BeforeCreate hook
func (c *Character) BeforeCreate(tx *gorm.DB) error {
	// Set current stats equal to base stats at creation
	if c.CurrentAttack == 0 {
		c.CurrentAttack = c.BaseAttack
	}
	if c.CurrentDefense == 0 {
		c.CurrentDefense = c.BaseDefense
	}
	if c.CurrentHP == 0 {
		c.CurrentHP = c.BaseHP
	}
	if c.CurrentSpeed == 0 {
		c.CurrentSpeed = c.BaseSpeed
	}
	return nil
}
