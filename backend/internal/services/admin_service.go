package services

import (
	"fmt"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// ConfigService moved to config_service.go

type AdminService struct {
	ls *LedgerService
}

func NewAdminService() *AdminService {
	return &AdminService{ls: NewLedgerService()}
}

func (s *AdminService) BanUser(targetID uint, reason string, durationHours int, adminID uint) error {
	var user models.User
	if err := db.DB.First(&user, targetID).Error; err != nil {
		return err
	}

	tx := db.DB.Begin()

	expires := time.Now().Add(time.Duration(durationHours) * time.Hour)
	user.Status = "BANNED"
	user.BanReason = reason
	user.BanExpiresAt = &expires
	user.Nonce = "" // Invalidate auth immediately

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Audit Log
	log := models.AdminAuditLog{
		AdminID:   adminID,
		Action:    "BAN_USER",
		TargetID:  fmt.Sprintf("%d", targetID),
		NewValue:  fmt.Sprintf("Reason: %s, Duration: %dh", reason, durationHours),
		CreatedAt: time.Now(),
	}
	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	tx.Commit()
	return nil
}

// GetRevenueStats returns total platform revenue from the Treasury account
func (s *AdminService) GetRevenueStats() (map[string]interface{}, error) {
	// 1. Get Treasury Account
	treasury, err := s.ls.GetOrCreateAccount(nil, models.AccountTypeTreasury, "GTK")
	if err != nil {
		return nil, err
	}
	towerTreasury, err := s.ls.GetOrCreateAccount(nil, models.AccountTypeTreasury, "TOWER")
	if err != nil {
		return nil, err
	}

	// 2. Aggregate Credits (Revenue)
	var totalGTK int64
	var totalTOWER int64

	// LedgerEntry where AccountID = Treasury AND Type = CREDIT
	db.DB.Model(&models.LedgerEntry{}).
		Where("account_id = ? AND type = 'CREDIT'", treasury.ID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalGTK)

	db.DB.Model(&models.LedgerEntry{}).
		Where("account_id = ? AND type = 'CREDIT'", towerTreasury.ID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalTOWER)

	// 3. User & Game Stats
	var totalUsers int64
	db.DB.Model(&models.User{}).Count(&totalUsers)

	var activeBattles int64
	// Check if Battle table exists (it should)
	db.DB.Model(&models.Battle{}).Where("status = ?", "ACTIVE").Count(&activeBattles)

	var activeRaids int64
	db.DB.Model(&models.RaidSession{}).Where("status = ?", "ACTIVE").Count(&activeRaids)

	return map[string]interface{}{
		"total_revenue_gtk":      totalGTK,
		"total_revenue_tower":    totalTOWER,
		"treasury_balance_gtk":   treasury.Balance,
		"treasury_balance_tower": towerTreasury.Balance,
		"total_users":            totalUsers,
		"active_battles":         activeBattles + activeRaids, // Combined metric
		"active_raids":           activeRaids,
	}, nil
}
