package models

import (
	"time"

	"gorm.io/gorm"
)

// Island represents a location where raids take place
type Island struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"size:50;not null;unique" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Difficulty  int    `gorm:"not null;check:difficulty > 0" json:"difficulty"` // 1-10
	MinLevelReq int    `gorm:"default:1" json:"min_level_req"`
	ImageURL    string `gorm:"size:255" json:"image_url"`

	// Relationships
	Bosses   []RaidBoss      `gorm:"foreignKey:IslandID" json:"bosses,omitempty"`
	Missions []IslandMission `gorm:"foreignKey:IslandID" json:"missions,omitempty"`
}

type IslandMission struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	IslandID    uint   `gorm:"not null;index" json:"island_id"`
	Sequence    int    `gorm:"not null" json:"sequence"` // 1, 2, 3...
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`

	// Enemy Stats (for this specific mission)
	EnemyName  string `gorm:"size:100;not null" json:"enemy_name"`
	EnemyType  string `gorm:"size:20;default:'NORMAL'" json:"enemy_type"` // FIRE, WATER, GRASS, etc. (Phase 10.3)
	EnemyHP    int64  `gorm:"not null" json:"enemy_hp"`
	EnemyAtk   int    `gorm:"not null" json:"enemy_atk"`
	EnemyDef   int    `gorm:"not null" json:"enemy_def"`
	EnemySpeed int    `gorm:"not null" json:"enemy_speed"`
	EnemyImage string `gorm:"type:varchar(255)" json:"enemy_image,omitempty"`

	RewardsPool string `gorm:"type:text" json:"rewards_pool"` // JSON
}

type UserCampaignProgress struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	UserID   uint `gorm:"not null;index" json:"user_id"`
	IslandID uint `gorm:"not null;index" json:"island_id"`

	// The highest sequence number completed. 0 = none.
	HighestSequence int `gorm:"default:0" json:"highest_sequence"`

	UpdatedAt time.Time `json:"updated_at"`
}

// RaidBoss represents a powerful enemy unit (NOW DEPRECATED in favor of Last Mission, but kept for legacy compat if needed)
type RaidBoss struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	IslandID uint   `gorm:"not null;index" json:"island_id"`
	Island   Island `gorm:"foreignKey:IslandID" json:"-"`

	Name          string `gorm:"size:50;not null" json:"name"`
	Element       string `gorm:"size:20;not null" json:"element"`        // FIRE, WATER, etc.
	CharacterType string `gorm:"size:20;not null" json:"character_type"` // DRAGON, BEAST, etc.

	// Stats
	TotalHP     int64 `gorm:"not null" json:"total_hp"`
	BaseAttack  int   `gorm:"not null" json:"base_attack"`
	BaseDefense int   `gorm:"not null" json:"base_defense"`
	Speed       int   `gorm:"not null" json:"speed"`

	// Rewards (JSON stored as string for simplicity in PoC, ideally separate table)
	RewardsPool string `gorm:"type:text" json:"rewards_pool"`

	// Visuals
	ImageURL string `gorm:"size:255" json:"image_url"`
}

// RaidSession represents an active battle attempt by a user
type RaidSession struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	MissionID uint          `gorm:"not null;index" json:"mission_id"`
	Mission   IslandMission `gorm:"foreignKey:MissionID" json:"mission"`

	// Optional Override for event bosses, otherwise use Mission.Enemy*
	BossID *uint     `gorm:"index" json:"boss_id,omitempty"`
	Boss   *RaidBoss `gorm:"foreignKey:BossID" json:"boss,omitempty"`

	TeamID uint  `json:"team_id"`
	Team   *Team `json:"team,omitempty" gorm:"foreignKey:TeamID;references:ID"`

	// Battle State (Individual Character System - Phase 13 - Axie Style)
	Status            string `gorm:"size:20;default:'IN_PROGRESS';index" json:"status"` // IN_PROGRESS, COMPLETED, FAILED, EXPIRED, ABANDONED
	ActiveCharacterID *uint  `json:"active_character_id,omitempty"`                     // Current fighter
	CharacterStates   string `gorm:"type:text" json:"character_states,omitempty"`       // JSON: [{"id":1,"hp":80}]
	TurnQueue         string `gorm:"type:text" json:"turn_queue,omitempty"`             // JSON: [{"type":"player","char_id":1,"speed":30},...]
	CurrentTurnIndex  int    `gorm:"default:0" json:"current_turn_index"`               // Position in turn queue

	// Team HP (kept for backward compatibility, computed from CharacterStates)
	CurrentTeamHP int64 `gorm:"not null" json:"current_team_hp"`
	InitialTeamHP int64 `json:"initial_team_hp,omitempty"`

	// Boss/Enemy State
	CurrentBossHP int64 `gorm:"not null" json:"current_boss_hp"`
	TotalHP       int64 `json:"total_hp,omitempty"` // Boss max HP (from Mission)

	// Progress Tracking
	CurrentStage int `gorm:"default:1" json:"current_stage"`
	TotalStages  int `gorm:"default:1" json:"total_stages"`
	TurnCount    int `gorm:"default:0" json:"turn_count"`

	DamageDealt      int64 `json:"damage_dealt"`
	TotalDamageTaken int64 `json:"total_damage_taken"`

	// Rewards
	TokensEarned     int    `gorm:"default:0" json:"tokens_earned"`
	XPEarned         int    `gorm:"default:0" json:"xp_earned"`
	PerformanceGrade string `json:"performance_grade,omitempty"` // S, A, B, C, D
	RewardsClaimed   bool   `gorm:"default:false" json:"rewards_claimed"`

	// Status Effects
	ActiveStatusEffects string `gorm:"type:text" json:"active_status_effects,omitempty"`

	// Timestamps
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// Security / Validation
	BattleSeed string `gorm:"size:64" json:"-"` // Deterministic seed for verifying logic
	IsVerified bool   `gorm:"default:false" json:"is_verified"`
}
