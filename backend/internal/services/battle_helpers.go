package services

import "github.com/lorengraff/crypto-tower-defense/internal/models"

// Helper function to convert CharacterState to BattleParticipant
func convertToParticipant(state *CharacterState) models.BattleParticipant {
	return models.BattleParticipant{
		CharacterID:   state.CharID,
		CharacterName: "", // No Name field in CharacterState
		CurrentHP:     int(state.CurrentHP),
		MaxHP:         int(state.MaxHP),
		CurrentMana:   state.CurrentMana,
		MaxMana:       state.MaxMana,
		Attack:        state.CurrentAttack,
		Defense:       state.CurrentDefense,
		Speed:         state.Speed, // FIXED: Use Speed directly
		IsFainted:     state.CurrentHP <= 0,
		IsActive:      true,
		Position:      0,
		TeamID:        0,
		Buffs:         []models.Buff{},
		Debuffs:       []models.Buff{},
		StatusEffect:  "",
	}
}
