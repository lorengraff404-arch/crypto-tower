package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// TokenService handles TOWER↔GTK conversion and withdrawals with security
type TokenService struct {
	conversionRate int64 // TOWER:GTK ratio (from config, default 100:1)
}

// NewTokenService creates a new token service
func NewTokenService() *TokenService {
	return &TokenService{
		conversionRate: 100, // 1 TOWER = 100 GTK (configurable)
	}
}

// ConvertTowerToGTK converts TOWER tokens to GTK for in-game purchases
func (s *TokenService) ConvertTowerToGTK(userID uint, towerAmount int64) error {
	// SECURITY CHECK 1: Validate amount
	if towerAmount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if towerAmount > 1000000 {
		return errors.New("amount exceeds maximum conversion limit (1,000,000 TOWER)")
	}

	// SECURITY CHECK 2: Rate limiting - max 10 conversions per hour
	var recentConversions int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	db.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND transaction_type = 'TOWER_TO_GTK_CONVERSION' AND created_at > ?", userID, oneHourAgo).
		Count(&recentConversions)

	if recentConversions >= 10 {
		return errors.New("conversion rate limit exceeded (max 10 per hour)")
	}

	// SECURITY CHECK 3: Daily conversion limit
	today := time.Now().Format("2006-01-02")
	var todayTotal int64
	db.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND transaction_type = 'TOWER_TO_GTK_CONVERSION' AND DATE(created_at) = ?", userID, today).
		Select("COALESCE(SUM(ABS(amount)), 0)").
		Scan(&todayTotal)

	dailyLimit := int64(50000) // 50,000 TOWER per day
	if todayTotal+towerAmount > dailyLimit {
		return fmt.Errorf("daily conversion limit exceeded (limit: %d, used: %d)", dailyLimit, todayTotal)
	}

	// Get user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// SECURITY CHECK 4: Verify sufficient TOWER balance
	if user.TOWERBalance < towerAmount {
		return fmt.Errorf("insufficient TOWER balance (have: %d, need: %d)", user.TOWERBalance, towerAmount)
	}

	// Calculate GTK amount
	gtkAmount := towerAmount * s.conversionRate

	// BEGIN ATOMIC TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Deduct TOWER
	if err := tx.Model(&user).Update("tower_balance", gorm.Expr("tower_balance - ?", towerAmount)).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to deduct TOWER")
	}

	// Add GTK
	if err := tx.Model(&user).Update("gtk_balance", gorm.Expr("gtk_balance + ?", gtkAmount)).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to add GTK")
	}

	// Create transaction record
	transaction := models.Transaction{
		UserID:          userID,
		TransactionType: "TOWER_TO_GTK_CONVERSION",
		TokenType:       "TOWER",
		Amount:          -towerAmount,
		BalanceBefore:   user.TOWERBalance,
		BalanceAfter:    user.TOWERBalance - towerAmount,
		Description:     fmt.Sprintf("Converted %d TOWER to %d GTK", towerAmount, gtkAmount),
		Metadata:        fmt.Sprintf("{\"gtk_amount\":%d,\"rate\":%d}", gtkAmount, s.conversionRate),
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to create transaction record")
	}

	// Audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "TOWER_TO_GTK_CONVERSION",
		EntityType: "transaction",
		EntityID:   &transaction.ID,
		NewValues:  fmt.Sprintf("tower:%d,gtk:%d", towerAmount, gtkAmount),
	}
	tx.Create(&auditLog)

	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction failed")
	}

	return nil
}

// ConvertGTKToTower converts GTK back to TOWER for withdrawal
func (s *TokenService) ConvertGTKToTower(userID uint, gtkAmount int64) error {
	// SECURITY CHECK 1: Validate amount
	if gtkAmount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	// SECURITY CHECK 2: GTK must be divisible by conversion rate
	if gtkAmount%s.conversionRate != 0 {
		return fmt.Errorf("GTK amount must be divisible by %d", s.conversionRate)
	}

	// SECURITY CHECK 3: Rate limiting
	var recentConversions int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	db.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND transaction_type = 'GTK_TO_TOWER_CONVERSION' AND created_at > ?", userID, oneHourAgo).
		Count(&recentConversions)

	if recentConversions >= 10 {
		return errors.New("conversion rate limit exceeded (max 10 per hour)")
	}

	// Get user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// SECURITY CHECK 4: Verify sufficient GTK balance
	if user.GTKBalance < gtkAmount {
		return fmt.Errorf("insufficient GTK balance (have: %d, need: %d)", user.GTKBalance, gtkAmount)
	}

	// Calculate TOWER amount (with 1% conversion fee)
	towerAmount := gtkAmount / s.conversionRate
	conversionFee := towerAmount / 100 // 1% fee
	finalTowerAmount := towerAmount - conversionFee

	// BEGIN ATOMIC TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Deduct GTK
	if err := tx.Model(&user).Update("gtk_balance", gorm.Expr("gtk_balance - ?", gtkAmount)).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to deduct GTK")
	}

	// Add TOWER (minus fee)
	if err := tx.Model(&user).Update("tower_balance", gorm.Expr("tower_balance + ?", finalTowerAmount)).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to add TOWER")
	}

	// Create transaction record
	transaction := models.Transaction{
		UserID:          userID,
		TransactionType: "GTK_TO_TOWER_CONVERSION",
		TokenType:       "GTK",
		Amount:          -gtkAmount,
		BalanceBefore:   user.GTKBalance,
		BalanceAfter:    user.GTKBalance - gtkAmount,
		Description:     fmt.Sprintf("Converted %d GTK to %d TOWER (fee: %d)", gtkAmount, finalTowerAmount, conversionFee),
		Metadata:        fmt.Sprintf("{\"tower_amount\":%d,\"fee\":%d}", finalTowerAmount, conversionFee),
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to create transaction record")
	}

	// Audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "GTK_TO_TOWER_CONVERSION",
		EntityType: "transaction",
		EntityID:   &transaction.ID,
		NewValues:  fmt.Sprintf("gtk:%d,tower:%d,fee:%d", gtkAmount, finalTowerAmount, conversionFee),
	}
	tx.Create(&auditLog)

	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction failed")
	}

	return nil
}

// WithdrawTower initiates TOWER withdrawal to wallet
func (s *TokenService) WithdrawTower(userID uint, amount int64, walletAddress string) error {
	// SECURITY CHECK 1: Validate amount
	minWithdrawal := int64(100) // Minimum 100 TOWER
	if amount < minWithdrawal {
		return fmt.Errorf("minimum withdrawal is %d TOWER", minWithdrawal)
	}

	if amount > 1000000 {
		return errors.New("amount exceeds maximum withdrawal limit (1,000,000 TOWER)")
	}

	// SECURITY CHECK 2: Validate wallet address
	if walletAddress == "" || len(walletAddress) != 42 {
		return errors.New("invalid wallet address")
	}

	// SECURITY CHECK 3: Daily withdrawal limit
	today := time.Now().Format("2006-01-02")
	var todayWithdrawals int64
	db.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND transaction_type = 'WITHDRAWAL' AND DATE(created_at) = ?", userID, today).
		Select("COALESCE(SUM(ABS(amount)), 0)").
		Scan(&todayWithdrawals)

	dailyLimit := int64(10000) // 10,000 TOWER per day
	if todayWithdrawals+amount > dailyLimit {
		return fmt.Errorf("daily withdrawal limit exceeded (limit: %d, used: %d)", dailyLimit, todayWithdrawals)
	}

	// SECURITY CHECK 4: Check for pending withdrawals
	var pendingCount int64
	db.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND transaction_type = 'WITHDRAWAL' AND blockchain_tx_hash IS NULL", userID).
		Count(&pendingCount)

	if pendingCount > 0 {
		return errors.New("you have a pending withdrawal, please wait for it to complete")
	}

	// Get user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// SECURITY CHECK 5: Verify wallet address matches user's linked wallet
	if user.WalletAddress != "" && user.WalletAddress != walletAddress {
		return errors.New("wallet address does not match linked wallet")
	}

	// SECURITY CHECK 6: Verify sufficient balance
	if user.TOWERBalance < amount {
		return fmt.Errorf("insufficient TOWER balance (have: %d, need: %d)", user.TOWERBalance, amount)
	}

	// BEGIN ATOMIC TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Deduct TOWER (mark as pending)
	if err := tx.Model(&user).Update("tower_balance", gorm.Expr("tower_balance - ?", amount)).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to deduct TOWER")
	}

	// Create withdrawal transaction (pending blockchain confirmation)
	transaction := models.Transaction{
		UserID:          userID,
		TransactionType: "WITHDRAWAL",
		TokenType:       "TOWER",
		Amount:          -amount,
		BalanceBefore:   user.TOWERBalance,
		BalanceAfter:    user.TOWERBalance - amount,
		Description:     fmt.Sprintf("Withdrawal of %d TOWER to %s", amount, walletAddress),
		Metadata:        fmt.Sprintf("{\"wallet\":\"%s\"}", walletAddress),
		IsOnChain:       false, // Will be set to true when blockchain tx confirms
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to create withdrawal request")
	}

	// Audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "WITHDRAWAL_REQUESTED",
		EntityType: "transaction",
		EntityID:   &transaction.ID,
		NewValues:  fmt.Sprintf("amount:%d,wallet:%s", amount, walletAddress),
	}
	tx.Create(&auditLog)

	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction failed")
	}

	// TODO: Queue blockchain transaction for processing
	// This would be handled by a separate worker service

	return nil
}

// DepositTower processes TOWER deposit from wallet
func (s *TokenService) DepositTower(userID uint, txHash string, amount int64) error {
	// SECURITY CHECK 1: Validate transaction hash
	if txHash == "" || len(txHash) != 66 {
		return errors.New("invalid transaction hash")
	}

	// SECURITY CHECK 2: Check if transaction already processed
	var existing models.Transaction
	err := db.DB.Where("blockchain_tx_hash = ?", txHash).First(&existing).Error
	if err == nil {
		return errors.New("transaction already processed")
	}

	// SECURITY CHECK 3: Verify transaction on blockchain
	// TODO: Implement blockchain verification
	// - Check transaction exists on opBNB
	// - Verify recipient is game contract
	// - Verify amount matches
	// - Verify confirmations >= 12 blocks

	// Get user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// BEGIN ATOMIC TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Add TOWER
	if err := tx.Model(&user).Update("tower_balance", gorm.Expr("tower_balance + ?", amount)).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to add TOWER")
	}

	// Create deposit transaction
	transaction := models.Transaction{
		UserID:           userID,
		TransactionType:  "DEPOSIT",
		TokenType:        "TOWER",
		Amount:           amount,
		BalanceBefore:    user.TOWERBalance,
		BalanceAfter:     user.TOWERBalance + amount,
		Description:      fmt.Sprintf("Deposit of %d TOWER", amount),
		BlockchainTxHash: &txHash,
		ChainID:          204, // opBNB testnet
		IsOnChain:        true,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to create deposit record")
	}

	// Audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "DEPOSIT_CONFIRMED",
		EntityType: "transaction",
		EntityID:   &transaction.ID,
		NewValues:  fmt.Sprintf("amount:%d,tx:%s", amount, txHash),
	}
	tx.Create(&auditLog)

	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction failed")
	}

	return nil
}

// GetBalance returns user's token balances
func (s *TokenService) GetBalance(userID uint) (map[string]int64, error) {
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return map[string]int64{
		"tower": user.TOWERBalance,
		"gtk":   user.GTKBalance,
	}, nil
}

// GetTransactionHistory returns user's transaction history
func (s *TokenService) GetTransactionHistory(userID uint, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := db.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
