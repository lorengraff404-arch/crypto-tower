package services

import (
	"errors"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// --- Missing Methods for BattleHandler Compatibility ---

// FindMatch attempts to find an opponent for the given user
// FindMatch attempts to find an opponent for the given user based on ELO
func (s *BattleService) FindMatch(userID uint, betAmount int64) (*models.User, error) {
	// 1. Get User's ELO
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Search for opponent within ELO range (+/- 200)
	var opponent models.User
	err := db.DB.Where("id != ? AND elo_rating BETWEEN ? AND ?", userID, user.ELO-200, user.ELO+200).
		Order("RANDOM()").
		First(&opponent).Error

	// 3. Fallback: If no close match, widen search (+/- 500)
	if err != nil {
		err = db.DB.Where("id != ? AND elo_rating BETWEEN ? AND ?", userID, user.ELO-500, user.ELO+500).
			Order("RANDOM()").
			First(&opponent).Error
	}

	// 4. Final Fallback: Any opponent
	if err != nil {
		err = db.DB.Where("id != ?", userID).
			Order("RANDOM()").
			First(&opponent).Error
	}

	if err != nil {
		return nil, errors.New("no opponent found")
	}

	return &opponent, nil
}

// CreatePvPBattle creates a battle between two players
func (s *BattleService) CreatePvPBattle(p1ID, p2ID uint, betAmount int64) (*models.Battle, error) {
	battle := &models.Battle{
		Player1ID:  p1ID,
		Player2ID:  p2ID,
		BattleType: "pvp", // or "wager" if bet > 0
		Status:     "active",
		Player1Bet: betAmount,
		Player2Bet: betAmount, // Simplification for legacy handler
		CreatedAt:  time.Now(),
	}
	if betAmount > 0 {
		battle.BattleType = "wager"
	}

	if err := s.InitializeBattleState(battle); err != nil {
		return nil, err
	}

	if err := db.DB.Create(battle).Error; err != nil {
		return nil, err
	}
	return battle, nil
}

// StartBattle is a placeholder to satisfy the handler (logic might be in Initialize)
func (s *BattleService) StartBattle(battleID uint) (*models.Battle, error) {
	return s.GetBattleByID(battleID)
}

// SurrenderBattle handles a player surrendering
func (s *BattleService) SurrenderBattle(battleID, userID uint) error {
	battle, err := s.GetBattleByID(battleID)
	if err != nil {
		return err
	}
	if battle.Status != "active" {
		return errors.New("battle not active")
	}

	// Set winner to the other player
	var winnerID uint
	if battle.Player1ID == userID {
		winnerID = battle.Player2ID
	} else {
		winnerID = battle.Player1ID
	}

	battle.WinnerID = &winnerID
	battle.Status = "completed"
	return db.DB.Save(battle).Error
}

// GetBattleByID retrieves a battle by ID with preloads
func (s *BattleService) GetBattleByID(id uint) (*models.Battle, error) {
	var battle models.Battle
	err := db.DB.Preload("Player1").Preload("Player2").First(&battle, id).Error
	if err != nil {
		return nil, err
	}
	return &battle, nil
}

// GetBattleHistory retrieves completed battles for a user
func (s *BattleService) GetBattleHistory(userID uint, limit int) ([]models.Battle, error) {
	var battles []models.Battle
	err := db.DB.Where("(player1_id = ? OR player2_id = ?) AND status = ?", userID, userID, "completed").
		Order("created_at desc").
		Limit(limit).
		Find(&battles).Error
	return battles, err
}

// RequestRematch creates a new battle with the same participants
func (s *BattleService) RequestRematch(battleID, userID uint) (*models.Battle, error) {
	// Re-create battle with same opponents
	oldBattle, err := s.GetBattleByID(battleID)
	if err != nil {
		return nil, err
	}
	// Swap logic or same?
	return s.CreatePvPBattle(oldBattle.Player1ID, oldBattle.Player2ID, oldBattle.Player1Bet)
}
