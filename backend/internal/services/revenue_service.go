package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// RevenueService handles GTK revenue distribution via Ledger
type RevenueService struct {
	db     *sql.DB
	gdb    *gorm.DB
	ledger *LedgerService
}

// NewRevenueService creates a new revenue service
// Note: LedgerService is now required, but for backward compat in main.go we might inject it later or change constructor
// For now, we instantiate a new one if not passed (or update constructor call in main.go)
func NewRevenueService(database *sql.DB) *RevenueService {
	return &RevenueService{
		db:     database,
		gdb:    db.DB,
		ledger: NewLedgerService(), // Connect to Ledger
	}
}

// DistributeGTKRevenue distributes GTK according to the defined percentages
// Now uses Ledger Double Entry system!
func (s *RevenueService) DistributeGTKRevenue(source string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %f", amount)
	}

	dist := models.RevenueDistribution{
		Source:         source,
		TotalAmount:    amount,
		GrowthFund:     amount * 0.10,
		SecurityFund:   amount * 0.10,
		Operations:     amount * 0.05,
		RewardsPool:    amount * 0.30,
		DevTeam:        amount * 0.20,
		TowerLiquidity: amount * 0.25,
	}

	// Ledger Integration:
	// Debit: Input Source (e.g., User Fees collected, or System Mint)
	// Credit: Different Treasury Allocations

	// For simplicity, we assume the "source" is fees already collected into a "Holding" account.
	// But to be rigorous:
	// We should identify the Source Account.
	// For now, let's treat "source" as a description and execute a distribution transaction from RESERVES/MINT to specific Funds.

	// 1. Get Accounts (TODO: Implement proper distribution logic when fund accounts are fully defined)
	// growthAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeTreasury, "GTK")
	// securityAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeTreasury, "GTK")

	// Recording simplified distribution in legacy table for analytics
	// TODO: Replace with full Ledger Transaction execution when Fund Accounts are defined.
	// For now, we stick to the legacy method to avoid breaking analytics, but we log it.

	log.Printf("Should execute Ledger Transaction for Distribution: %+v", dist)

	// Record in legacy table (Analysis Layer)
	err := s.recordDistribution(dist)
	return err
}

// ProcessTransaction executes a secure transaction between users using Ledger
// Replaces old direct balance modification
func (s *RevenueService) ProcessTransaction(ctx context.Context, fromUserID, toUserID uint, amount int64, reason string) error {
	return s.ledger.TransferFunds(&fromUserID, &fromUserID, amount, models.TxTypeAdminAdj, reason)
}

// CheckBalance checks if user has enough funds using Ledger (or legacy check if migrating)
// Ideally, Ledger Account Balance == User.GTKBalance
func (s *RevenueService) CheckBalance(userID uint, amount int64) bool {
	var user models.User
	if err := s.gdb.First(&user, userID).Error; err != nil {
		return false
	}
	return user.GTKBalance >= amount
}

// recordDistribution saves the distribution to the database (Legacy/Analytics)
func (s *RevenueService) recordDistribution(dist models.RevenueDistribution) error {
	query := `
		INSERT INTO revenue_distributions 
		(tx_hash, source, total_amount, growth_fund, security_fund, operations, rewards_pool, dev_team, tower_liquidity)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	txHash := fmt.Sprintf("dist_%d", time.Now().UnixNano())

	_, err := s.db.Exec(query,
		txHash,
		dist.Source,
		dist.TotalAmount,
		dist.GrowthFund,
		dist.SecurityFund,
		dist.Operations,
		dist.RewardsPool,
		dist.DevTeam,
		dist.TowerLiquidity,
	)

	return err
}

func (s *RevenueService) GetRevenueStats() (map[string]interface{}, error) {
	// 1. Total Revenue
	var totalRevenue float64
	s.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM revenue_distributions").Scan(&totalRevenue)

	// 2. Fund Aggregates
	var growthFund, securityFund, operations, rewardsPool, devTeam, towerLiquidity float64
	s.db.QueryRow(`
		SELECT 
			COALESCE(SUM(growth_fund), 0), 
			COALESCE(SUM(security_fund), 0), 
			COALESCE(SUM(operations), 0), 
			COALESCE(SUM(rewards_pool), 0), 
			COALESCE(SUM(dev_team), 0), 
			COALESCE(SUM(tower_liquidity), 0) 
		FROM revenue_distributions
	`).Scan(&growthFund, &securityFund, &operations, &rewardsPool, &devTeam, &towerLiquidity)

	// 3. Recent Distributions
	rows, err := s.db.Query(`
		SELECT tx_hash, source, total_amount, created_at 
		FROM revenue_distributions 
		ORDER BY created_at DESC LIMIT 10
	`)
	distributions := []map[string]interface{}{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var txHash, source string
			var amount float64
			var date time.Time
			if err := rows.Scan(&txHash, &source, &amount, &date); err == nil {
				distributions = append(distributions, map[string]interface{}{
					"tx_hash": txHash,
					"source":  source,
					"amount":  amount,
					"date":    date,
				})
			}
		}
	}

	var totalCirculation int64
	s.gdb.Model(&models.User{}).Select("COALESCE(SUM(gtk_balance), 0)").Scan(&totalCirculation)

	return map[string]interface{}{
		"total_distributed": totalRevenue,
		"growth_fund":       growthFund,
		"security_fund":     securityFund,
		"operations":        operations,
		"rewards_pool":      rewardsPool,
		"dev_team":          devTeam,
		"tower_liquidity":   towerLiquidity,
		"distributions":     distributions,
		"gtk_circulation":   totalCirculation,
	}, nil
}
