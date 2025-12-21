package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// BattleState stores the complete state of an ongoing battle
type BattleState struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Battle reference
	BattleID uint   `gorm:"uniqueIndex;not null" json:"battle_id"`
	Battle   Battle `gorm:"foreignKey:BattleID" json:"-"`

	// Turn tracking
	CurrentTurn int    `gorm:"default:1" json:"current_turn"`
	Phase       string `gorm:"size:20;default:'SELECT_ACTION'" json:"phase"` // SELECT_ACTION, EXECUTING, FINISHED

	// State JSON (full battle state)
	StateJSON string `gorm:"type:text;not null" json:"-"`

	// Quick access fields
	Player1ActiveIndex int  `gorm:"default:0" json:"player1_active_index"`
	Player2ActiveIndex int  `gorm:"default:0" json:"player2_active_index"`
	IsFinished         bool `gorm:"default:false" json:"is_finished"`
}

// BattleStateData is the full state structure (stored as JSON)
type BattleStateData struct {
	BattleID    uint   `json:"battle_id"`
	Mode        string `json:"mode"` // PVE, RANKED, WAGER
	CurrentTurn int    `json:"current_turn"`
	Phase       string `json:"phase"`

	// Players
	Player1ID   uint                   `json:"player1_id"`
	Player2ID   uint                   `json:"player2_id"`
	Player1Team []CharacterBattleState `json:"player1_team"`
	Player2Team []CharacterBattleState `json:"player2_team"`

	Player1ActiveIndex int `json:"player1_active_index"`
	Player2ActiveIndex int `json:"player2_active_index"`

	// Battle log
	ActionLog []BattleAction `json:"action_log"`

	// Turn queue (sorted by speed)
	TurnQueue []TurnQueueEntry `json:"turn_queue"`

	// Wager specific
	Player1Stake int `json:"player1_stake,omitempty"`
	Player2Stake int `json:"player2_stake,omitempty"`

	// Result
	Finished bool  `json:"finished"`
	Winner   *uint `json:"winner,omitempty"`
}

// CharacterBattleState represents a character's state during battle
type CharacterBattleState struct {
	CharacterID uint   `json:"character_id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`

	// Current stats (modified by buffs/debuffs)
	CurrentHP   int `json:"current_hp"`
	MaxHP       int `json:"max_hp"`
	CurrentMana int `json:"current_mana"`
	MaxMana     int `json:"max_mana"`

	// Base stats
	BaseAttack    int `json:"base_attack"`
	BaseDefense   int `json:"base_defense"`
	BaseSpAttack  int `json:"base_sp_attack"`
	BaseSpDefense int `json:"base_sp_defense"`
	BaseSpeed     int `json:"base_speed"`

	// Type
	ElementType string `json:"element_type"` // FIRE, WATER, etc.
	Class       string `json:"class"`        // WARRIOR, MAGE, etc.

	// Abilities
	Abilities []AbilityData `json:"abilities"`

	// Active effects
	Buffs        []ActiveBuff   `json:"buffs"`
	Debuffs      []ActiveDebuff `json:"debuffs"`
	StatusEffect string         `json:"status_effect,omitempty"` // POISON, BURN, etc.

	// State
	IsFainted bool `json:"is_fainted"`
}

// AbilityData represents an ability in battle
type AbilityData struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ElementType string `json:"element_type"`
	DamageType  string `json:"damage_type"` // PHYSICAL, SPECIAL
	BaseDamage  int    `json:"base_damage"`
	Accuracy    int    `json:"accuracy"`
	ManaCost    int    `json:"mana_cost"`
	EffectType  string `json:"effect_type"` // DAMAGE, BUFF, DEBUFF, HEAL, STATUS
	EffectValue int    `json:"effect_value"`
	EffectStat  string `json:"effect_stat,omitempty"` // ATK, DEF, etc.
	EffectTurns int    `json:"effect_turns,omitempty"`
	TargetType  string `json:"target_type"` // SINGLE, ALL_ENEMIES, SELF, ALL_ALLIES
}

// ActiveBuff represents a buff on a character
type ActiveBuff struct {
	Stat       string  `json:"stat"`       // ATK, DEF, SP_ATK, SP_DEF, SPEED
	Multiplier float64 `json:"multiplier"` // 1.5 = +50%
	Duration   int     `json:"duration"`   // Turns remaining
	Source     string  `json:"source"`     // Ability name
}

// ActiveDebuff represents a debuff on a character
type ActiveDebuff struct {
	Stat       string  `json:"stat"`
	Multiplier float64 `json:"multiplier"` // 0.5 = -50%
	Duration   int     `json:"duration"`
	Source     string  `json:"source"`
}

// BattleAction represents a single action in battle
type BattleAction struct {
	Turn      int    `json:"turn"`
	ActorID   uint   `json:"actor_id"`
	ActorName string `json:"actor_name"`
	Action    string `json:"action"` // "used Fireball", "fainted", etc.
	TargetID  uint   `json:"target_id,omitempty"`
	Damage    int    `json:"damage,omitempty"`
	Healing   int    `json:"healing,omitempty"`
	Effects   string `json:"effects,omitempty"` // "Critical hit!", "Super effective!"
}

// TurnQueueEntry represents who acts next
type TurnQueueEntry struct {
	CharacterID uint `json:"character_id"`
	IsPlayer1   bool `json:"is_player1"`
	Speed       int  `json:"speed"`
}

// GetStateData deserializes the JSON state
func (bs *BattleState) GetStateData() (*BattleStateData, error) {
	var state BattleStateData
	if err := json.Unmarshal([]byte(bs.StateJSON), &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// SetStateData serializes the state to JSON
func (bs *BattleState) SetStateData(state *BattleStateData) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	bs.StateJSON = string(data)
	bs.CurrentTurn = state.CurrentTurn
	bs.Phase = state.Phase
	bs.Player1ActiveIndex = state.Player1ActiveIndex
	bs.Player2ActiveIndex = state.Player2ActiveIndex
	bs.IsFinished = state.Finished
	return nil
}
