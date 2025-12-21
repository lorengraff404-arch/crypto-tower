package services

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// TurnEntry represents one participant in thebattle turn queue (Axie-style)
type TurnEntry struct {
	Type      string `json:"type"`       // "player" or "enemy"
	CharID    uint   `json:"char_id"`    // Character ID (0 for enemy)
	Speed     int    `json:"speed"`      // Speed stat for ordering
	Name      string `json:"name"`       // Display name
	CurrentHP int64  `json:"current_hp"` // Current health
}

// CharacterState tracks individual character HP during battle
type CharacterState struct {
	CharID      uint  `json:"char_id"`
	CurrentHP   int64 `json:"current_hp"`
	MaxHP       int64 `json:"max_hp"`
	CurrentMana int   `json:"current_mana"`
	MaxMana     int   `json:"max_mana"`
	Speed       int   `json:"speed"`

	// Real-time Combat Stats (Base + Buffs)
	CurrentAttack  int    `json:"current_attack"`
	CurrentDefense int    `json:"current_defense"`
	Level          int    `json:"level"`
	Element        string `json:"element"`
	Class          string `json:"class"`

	IsDead bool `json:"is_dead"`
}

// buildTurnQueue creates the initial turn order based on speed stats (Axie-style)
func (s *RaidService) buildTurnQueue(team *models.Team, mission *models.IslandMission) ([]TurnEntry, error) {
	queue := []TurnEntry{}

	// Add player characters (active members only)
	for _, member := range team.Members {
		if !member.IsBackup {
			char := member.Character
			queue = append(queue, TurnEntry{
				Type:      "player",
				CharID:    char.ID,
				Speed:     char.CurrentSpeed,
				Name:      char.Class,
				CurrentHP: int64(char.CurrentHP),
			})
		}
	}

	if len(queue) == 0 {
		return nil, errors.New("team has no active characters")
	}

	// Add enemy
	queue = append(queue, TurnEntry{
		Type:      "enemy",
		CharID:    0, // Enemy doesn't have a char ID
		Speed:     mission.EnemySpeed,
		Name:      mission.EnemyName,
		CurrentHP: mission.EnemyHP,
	})

	// Sort by Speed (DESC) - fastest character acts first
	sort.Slice(queue, func(i, j int) bool {
		if queue[i].Speed == queue[j].Speed {
			// Tie-breaker: player goes first
			return queue[i].Type == "player"
		}
		return queue[i].Speed > queue[j].Speed
	})

	return queue, nil
}

// initializeCharacterStates creates initial HP tracking for all team members
func (s *RaidService) initializeCharacterStates(team *models.Team) ([]CharacterState, error) {
	states := []CharacterState{}

	for _, member := range team.Members {
		if !member.IsBackup {
			char := member.Character
			states = append(states, CharacterState{
				CharID:         char.ID,
				CurrentHP:      int64(char.CurrentHP),
				MaxHP:          int64(char.CurrentHP), // Assume current = max at start
				CurrentMana:    char.MaxMana,
				MaxMana:        char.MaxMana,
				Speed:          char.CurrentSpeed,
				CurrentAttack:  char.CurrentAttack,
				CurrentDefense: char.CurrentDefense,
				Level:          char.Level,
				Element:        char.Element,
				Class:          char.Class,
				IsDead:         false,
			})
		}
	}

	return states, nil
}

// getCurrentTurn gets the current turn entry from the queue
func (s *RaidService) getCurrentTurn(session *models.RaidSession) (*TurnEntry, error) {
	if session.TurnQueue == "" {
		return nil, errors.New("turn queue not initialized")
	}

	var queue []TurnEntry
	if err := json.Unmarshal([]byte(session.TurnQueue), &queue); err != nil {
		return nil, err
	}

	if len(queue) == 0 {
		return nil, errors.New("turn queue is empty")
	}

	index := session.CurrentTurnIndex % len(queue) // Wrap around
	return &queue[index], nil
}

// advanceTurn moves to the next character in the turn queue
func (s *RaidService) advanceTurn(session *models.RaidSession) error {
	var queue []TurnEntry
	if err := json.Unmarshal([]byte(session.TurnQueue), &queue); err != nil {
		return err
	}

	session.CurrentTurnIndex++
	if session.CurrentTurnIndex >= len(queue) {
		session.CurrentTurnIndex = 0 // Start new round
		session.TurnCount++          // Increment round counter
	}

	return nil
}

// updateCharacterHP updates a specific character's HP in CharacterStates
func (s *RaidService) updateCharacterHP(session *models.RaidSession, charID uint, newHP int64) error {
	var states []CharacterState
	if err := json.Unmarshal([]byte(session.CharacterStates), &states); err != nil {
		return err
	}

	for i := range states {
		if states[i].CharID == charID {
			states[i].CurrentHP = newHP
			if newHP < 0 {
				states[i].CurrentHP = 0
			}
			break
		}
	}

	// Re-serialize
	data, err := json.Marshal(states)
	if err != nil {
		return err
	}
	session.CharacterStates = string(data)

	// Also update CurrentTeamHP (legacy field)
	totalHP := int64(0)
	for _, state := range states {
		totalHP += state.CurrentHP
	}
	session.CurrentTeamHP = totalHP

	return nil
}

// GetCharacterHP retrieves a specific character's current HP (exported for future use)
func (s *RaidService) GetCharacterHP(session *models.RaidSession, charID uint) (int64, error) {
	var states []CharacterState
	if err := json.Unmarshal([]byte(session.CharacterStates), &states); err != nil {
		return 0, err
	}

	for _, state := range states {
		if state.CharID == charID {
			return state.CurrentHP, nil
		}
	}

	return 0, errors.New("character not found in battle")
}

// checkAllPlayerCharactersFainted checks if all player characters have 0 HP
func (s *RaidService) checkAllPlayerCharactersFainted(session *models.RaidSession) bool {
	var states []CharacterState
	if err := json.Unmarshal([]byte(session.CharacterStates), &states); err != nil {
		return true // Assume defeat on error
	}

	for _, state := range states {
		if state.CurrentHP > 0 {
			return false // At least one alive
		}
	}

	return true // All fainted
}
