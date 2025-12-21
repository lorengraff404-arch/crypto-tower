package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// RevenueDistribution represents a GTK revenue split
type RevenueDistribution struct {
	TxHash         string
	Source         string
	TotalAmount    float64
	GrowthFund     float64
	SecurityFund   float64
	Operations     float64
	RewardsPool    float64
	DevTeam        float64
	TowerLiquidity float64
}

// RevenueService handles GTK revenue distribution
type RevenueService struct {
	db *sql.DB
}

// NewRevenueService creates a new revenue service
func NewRevenueService(db *sql.DB) *RevenueService {
	return &RevenueService{db: db}
}

// DistributeGTKRevenue distributes GTK according to the defined percentages
func (s *RevenueService) DistributeGTKRevenue(source string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %f", amount)
	}

	// Calculate distribution (percentages from tokenomics)
	dist := RevenueDistribution{
		Source:         source,
		TotalAmount:    amount,
		GrowthFund:     amount * 0.10, // 10%
		SecurityFund:   amount * 0.10, // 10%
		Operations:     amount * 0.05, // 5%
		RewardsPool:    amount * 0.30, // 30% (Top 10 players)
		DevTeam:        amount * 0.20, // 20%
		TowerLiquidity: amount * 0.25, // 25%
	}

	// TODO: Execute blockchain transfers
	// For now, just simulate
	txHash := fmt.Sprintf("0x%x", time.Now().UnixNano())
	dist.TxHash = txHash

	// Record distribution in database
	err := s.recordDistribution(dist)
	if err != nil {
		return fmt.Errorf("failed to record distribution: %w", err)
	}

	log.Printf("GTK Revenue distributed: %.2f from %s (tx: %s)", amount, source, txHash)
	return nil
}

// recordDistribution saves the distribution to the database
func (s *RevenueService) recordDistribution(dist RevenueDistribution) error {
	query := `
		INSERT INTO revenue_distributions 
		(tx_hash, source, total_amount, growth_fund, security_fund, operations, rewards_pool, dev_team, tower_liquidity)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := s.db.Exec(query,
		dist.TxHash,
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

// GetRevenueStats returns revenue statistics
func (s *RevenueService) GetRevenueStats() (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_distributions,
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COALESCE(SUM(growth_fund), 0) as total_growth,
			COALESCE(SUM(security_fund), 0) as total_security,
			COALESCE(SUM(operations), 0) as total_operations,
			COALESCE(SUM(rewards_pool), 0) as total_rewards,
			COALESCE(SUM(dev_team), 0) as total_dev,
			COALESCE(SUM(tower_liquidity), 0) as total_liquidity
		FROM revenue_distributions
	`

	var stats struct {
		TotalDistributions int
		TotalRevenue       float64
		TotalGrowth        float64
		TotalSecurity      float64
		TotalOperations    float64
		TotalRewards       float64
		TotalDev           float64
		TotalLiquidity     float64
	}

	err := s.db.QueryRow(query).Scan(
		&stats.TotalDistributions,
		&stats.TotalRevenue,
		&stats.TotalGrowth,
		&stats.TotalSecurity,
		&stats.TotalOperations,
		&stats.TotalRewards,
		&stats.TotalDev,
		&stats.TotalLiquidity,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_distributions": stats.TotalDistributions,
		"total_revenue":       stats.TotalRevenue,
		"growth_fund":         stats.TotalGrowth,
		"security_fund":       stats.TotalSecurity,
		"operations":          stats.TotalOperations,
		"rewards_pool":        stats.TotalRewards,
		"dev_team":            stats.TotalDev,
		"tower_liquidity":     stats.TotalLiquidity,
	}, nil
}

// GetDistributionsBySource returns distributions grouped by source
func (s *RevenueService) GetDistributionsBySource() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			source,
			COUNT(*) as count,
			COALESCE(SUM(total_amount), 0) as total
		FROM revenue_distributions
		GROUP BY source
		ORDER BY total DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var source string
		var count int
		var total float64

		err := rows.Scan(&source, &count, &total)
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"source": source,
			"count":  count,
			"total":  total,
		})
	}

	return results, nil
}
