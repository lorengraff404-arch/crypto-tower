package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// ==================== BATTLE MANAGEMENT ====================

// GetActiveBattles returns all currently active battles
func (s *AdminService) GetActiveBattles() ([]models.Battle, error) {
	var battles []models.Battle
	// Preload players for display
	if err := db.DB.Preload("Player1").Preload("Player2").
		Where("status = ?", "active").
		Order("created_at desc").
		Find(&battles).Error; err != nil {
		return nil, err
	}
	return battles, nil
}

// GetBattleHistory returns global battle history with pagination
func (s *AdminService) GetBattleHistory(limit int) ([]models.Battle, error) {
	var battles []models.Battle
	if err := db.DB.Preload("Player1").Preload("Player2").Preload("Winner").
		Where("status != ?", "active").
		Order("created_at desc").
		Limit(limit).
		Find(&battles).Error; err != nil {
		return nil, err
	}
	return battles, nil
}

// TerminateBattle forces a battle to end
func (s *AdminService) TerminateBattle(battleID uint, reason string, adminID uint) error {
	var battle models.Battle
	if err := db.DB.First(&battle, battleID).Error; err != nil {
		return errors.New("battle not found")
	}

	if battle.Status != "active" {
		return errors.New("battle is not active")
	}

	// Force end
	battle.Status = "TERMINATED_BY_ADMIN"
	battle.ActionLog += fmt.Sprintf("\n[ADMIN] Battle terminated by admin %d. Reason: %s", adminID, reason)
	now := time.Now()
	battle.EndedAt = &now

	if err := db.DB.Save(&battle).Error; err != nil {
		return err
	}

	// Refund bets if Wager match?
	// For safety, we should probably REFUND if it was a forced termination unless specified otherwise.
	// Implementing Safe Refund Logic if bets exist and no winner declared
	if battle.BattleType == "wager" && (battle.Player1Bet > 0 || battle.Player2Bet > 0) {
		s.ls.UnlockFunds(battle.Player1ID, battle.Player1Bet, "GTK")
		s.ls.UnlockFunds(battle.Player2ID, battle.Player2Bet, "GTK")
	}

	s.CreateAuditLog(adminID, "TERMINATE_BATTLE", strconv.Itoa(int(battleID)), "active", "TERMINATED")
	return nil
}
