package models

import "time"

// Egg represents an unhatched character
type Egg struct {
	ID        uint  `gorm:"primaryKey" json:"id"`
	UserID    uint  `gorm:"not null;index" json:"user_id"`
	Parent1ID *uint `json:"parent1_id"`
	Parent2ID *uint `json:"parent2_id"`

	// Mint information (for gacha mints)
	MintCost   int64  `gorm:"default:0" json:"mint_cost"`  // TOWER spent on mint
	MintTxHash string `gorm:"size:66" json:"mint_tx_hash"` // Blockchain tx hash

	// Egg properties (visible on mint)
	Rarity        string `gorm:"size:20;not null" json:"rarity"`
	Element       string `gorm:"size:20" json:"element"`
	CharacterType string `gorm:"size:20" json:"character_type"`
	Class         string `gorm:"size:20" json:"class"` // Added for gacha

	// Predetermined traits (hidden until hatch or scanned)
	PredeterminedStats     string `gorm:"type:jsonb" json:"predetermined_stats,omitempty"`     // JSON: {hp, atk, def, spd}
	PredeterminedAbilities string `gorm:"type:jsonb" json:"predetermined_abilities,omitempty"` // JSON: [ability_ids]

	// Stats reveal
	IsStatsRevealed bool       `gorm:"default:false" json:"is_stats_revealed"`
	RevealedAt      *time.Time `json:"revealed_at,omitempty"`

	// Incubation
	IncubationTime          int        `gorm:"not null" json:"incubation_time"` // Base time in hours
	EffectiveIncubationTime int        `json:"effective_incubation_time"`       // After accelerators
	IncubationStartedAt     *time.Time `json:"incubation_started_at"`
	AcceleratorsApplied     string     `json:"accelerators_applied" gorm:"type:text"`
	HatchedAt               *time.Time `json:"hatched_at"`
	CharacterID             *uint      `json:"character_id"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Parent1   *Character `gorm:"foreignKey:Parent1ID" json:"parent1,omitempty"`
	Parent2   *Character `gorm:"foreignKey:Parent2ID" json:"parent2,omitempty"`
	Character *Character `gorm:"foreignKey:CharacterID" json:"character,omitempty"`
}

// BreedingSession tracks breeding attempts
type BreedingSession struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	UserID    uint `gorm:"not null;index" json:"user_id"`
	Parent1ID uint `gorm:"not null" json:"parent1_id"`
	Parent2ID uint `gorm:"not null" json:"parent2_id"`

	// Status
	Status      string     `gorm:"size:20;default:'in_progress'" json:"status"`
	StartedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`

	// Result
	EggID *uint `json:"egg_id"`

	// Cost
	TokensSpent int `gorm:"default:0" json:"tokens_spent"`

	// Relationships
	User    User      `gorm:"foreignKey:UserID" json:"-"`
	Parent1 Character `gorm:"foreignKey:Parent1ID" json:"parent1"`
	Parent2 Character `gorm:"foreignKey:Parent2ID" json:"parent2"`
	Egg     *Egg      `gorm:"foreignKey:EggID" json:"egg,omitempty"`
}
