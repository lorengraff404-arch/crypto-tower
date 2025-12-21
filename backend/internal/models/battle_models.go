package models

import (
	"time"
)

// BattleParticipant represents a character in battle with their current state
type BattleParticipant struct {
	CharacterID   uint   `json:"character_id"`
	CharacterName string `json:"character_name"`
	TeamID        uint   `json:"team_id"`
	Position      int    `json:"position"`  // 0-2 for 3v3
	IsActive      bool   `json:"is_active"` // Currently in battle (not switched out)

	// Current Battle Stats
	CurrentHP   int `json:"current_hp"`
	MaxHP       int `json:"max_hp"`
	CurrentMana int `json:"current_mana"`
	MaxMana     int `json:"max_mana"`
	Attack      int `json:"attack"`
	Defense     int `json:"defense"`
	Speed       int `json:"speed"`

	// Battle State
	IsFainted    bool   `json:"is_fainted"`
	Buffs        []Buff `json:"buffs"`
	Debuffs      []Buff `json:"debuffs"`
	StatusEffect string `json:"status_effect"` // burn, poison, paralysis, sleep, freeze

	// Equipped Abilities
	Ability1ID *uint `json:"ability1_id"`
	Ability2ID *uint `json:"ability2_id"`
	Ability3ID *uint `json:"ability3_id"`
	Ability4ID *uint `json:"ability4_id"`
}

// Buff represents a temporary stat modifier
type Buff struct {
	Name           string    `json:"name"`
	Stat           string    `json:"stat"`     // attack, defense, speed
	Modifier       float64   `json:"modifier"` // 1.5 = +50%, 0.5 = -50%
	TurnsRemaining int       `json:"turns_remaining"`
	AppliedAt      time.Time `json:"applied_at"`
}

// TurnAction represents an action taken during a turn
type TurnAction struct {
	ParticipantID uint   `json:"participant_id"`
	ActionType    string `json:"action_type"` // ability, item, switch, flee
	TargetID      *uint  `json:"target_id"`
	AbilityID     *uint  `json:"ability_id"`
	ItemID        *uint  `json:"item_id"`
	SwitchToPos   *int   `json:"switch_to_pos"`
}

// TurnResult represents the outcome of a turn
type TurnResult struct {
	AttackerID      uint   `json:"attacker_id"`
	DefenderID      uint   `json:"defender_id"`
	AbilityUsed     string `json:"ability_used"`
	Damage          int    `json:"damage"`
	Critical        bool   `json:"critical"`
	Effectiveness   string `json:"effectiveness"` // super_effective, not_very_effective, normal, immune
	DefenderHP      int    `json:"defender_hp"`
	DefenderFainted bool   `json:"defender_fainted"`
	Message         string `json:"message"`
}

// BattleTurnQueue manages turn order
type BattleTurnQueue struct {
	Turns []TurnEntry `json:"turns"`
}

// TurnEntry represents a participant's turn in queue
type TurnEntry struct {
	ParticipantID uint `json:"participant_id"`
	Speed         int  `json:"speed"`
	Priority      int  `json:"priority"` // Higher priority moves go first
}
