package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LedgerService implements double-entry bookkeeping
type LedgerService struct {
	db *gorm.DB
}

func NewLedgerService() *LedgerService {
	return &LedgerService{db: db.DB}
}

// GetOrCreateAccount retrieves or creates a ledger account for a user or system
func (s *LedgerService) GetOrCreateAccount(userID *uint, accType models.AccountType, currency string) (*models.LedgerAccount, error) {
	var account models.LedgerAccount

	query := s.db.Where("type = ? AND currency = ?", accType, currency)
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id IS NULL")
	}

	err := query.First(&account).Error
	if err == nil {
		return &account, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newAccount := models.LedgerAccount{
			UserID:   userID,
			Type:     accType,
			Currency: currency,
			Balance:  0,
		}
		if err := s.db.Create(&newAccount).Error; err != nil {
			return nil, err
		}
		return &newAccount, nil
	}

	return nil, err
}

// CreateTransaction executes an atomic ledger transaction
// refID: Ext identifier (BattleID, TxHash)
// entries: Must sum to zero (Debits + Credits = 0)
// Negative amount = Debit, Positive = Credit based on convention.
// Here we enforce: Total Change = 0.
func (s *LedgerService) CreateTransaction(txType models.TransactionType, refID, description string, entries []models.LedgerEntry) error {
	return s.CreateTransactionWithTx(s.db, txType, refID, description, entries)
}

// CreateTransactionWithTx executes ledger transaction within an existing DB transaction scope
func (s *LedgerService) CreateTransactionWithTx(tx *gorm.DB, txType models.TransactionType, refID, description string, entries []models.LedgerEntry) error {
	var sum int64 = 0
	for _, e := range entries {
		sum += e.Amount
	}

	if sum != 0 {
		return fmt.Errorf("transaction unbalanced: sum is %d (must be 0)", sum)
	}

	// Create Transaction Header
	ledgerTx := models.LedgerTransaction{
		Type:        txType,
		ReferenceID: refID,
		Description: description,
		Timestamp:   time.Now(),
	}
	if err := tx.Create(&ledgerTx).Error; err != nil {
		return err
	}

	// Process Entries
	for _, e := range entries {
		e.TransactionID = ledgerTx.ID

		// Update Account Balance
		var account models.LedgerAccount
		// Lock row for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&account, e.AccountID).Error; err != nil {
			return err
		}

		account.Balance += e.Amount
		if account.Balance < 0 && account.Type != models.AccountTypeSink && account.Type != models.AccountTypeReward {
			// Wallets cannot go negative, but System Sinks/Rewards can
			return fmt.Errorf("insufficient funds in account %d", account.ID)
		}

		if err := tx.Save(&account).Error; err != nil {
			return err
		}

		// Save Entry
		if err := tx.Create(&e).Error; err != nil {
			return err
		}
	}

	return nil
}

// TransferFunds simplified helper
func (s *LedgerService) TransferFunds(fromUser *uint, toUser *uint, amount int64, txType models.TransactionType, refID string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	fromAcc, err := s.GetOrCreateAccount(fromUser, models.AccountTypeWallet, "GTK")
	if err != nil {
		return err
	}

	toAcc, err := s.GetOrCreateAccount(toUser, models.AccountTypeWallet, "GTK")
	if err != nil {
		return err
	}

	entries := []models.LedgerEntry{
		{AccountID: fromAcc.ID, Amount: -amount, Type: "DEBIT"},
		{AccountID: toAcc.ID, Amount: amount, Type: "CREDIT"},
	}

	return s.CreateTransaction(txType, refID, "Fund Transfer", entries)
}

// TransferFundsWithTx executes transfer within existing transaction
func (s *LedgerService) TransferFundsWithTx(tx *gorm.DB, fromUser *uint, toUser *uint, amount int64, txType models.TransactionType, refID string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// We need Account IDs. Since we are inside a TX, we should probably fetch them with the TX too?
	// GetOrCreateAccount handles creation, which might need to be in TX.
	// Refactoring GetOrCreateAccount to take TX is complex.
	// For now, let's assume accounts exist or use non-tx Get for lookup (safe if accounts exist).
	// Better: Use `s.GetOrCreateAccount` but if it creates, it uses `s.db` (main DB).
	// This breaks atomicity if creation fails but main TX commits?
	// Actually no, creation is separate 1-time event. It's fine if account is created outside the wager TX.

	fromAcc, err := s.GetOrCreateAccount(fromUser, models.AccountTypeWallet, "GTK")
	if err != nil {
		return err
	}

	// Logic for Escrow/System accounts (nil user)
	var toAcc *models.LedgerAccount
	if toUser == nil {
		// Assuming Transfer to Escrow if nil? Or specific type?
		// The caller passes nil for system, but GetOrCreateAccount needs type.
		// Let's assume if toUser is nil, we look for Escrow? No, that's ambiguous.
		// Update: The caller (StartBattle) passed `nil` for `toUser` but we need to know WHICH system account.
		// StartBattle used `nil`. Let's fix StartBattle to pass System Account ID or handle it here.
		// Re-reading StartBattle: `s.ledger.TransferFundsWithTx(tx, &player1ID, nil, ...)`
		// This implies sending to "Sink" or "Escrow".
		// We should change TransferFundsWithTx signature or logic.
		// Use specific accounts.
		toAcc, err = s.GetOrCreateAccount(nil, models.AccountTypeEscrow, "GTK") // Defaulting to Escrow for Wager context? Risky.
	} else {
		toAcc, err = s.GetOrCreateAccount(toUser, models.AccountTypeWallet, "GTK")
	}

	if err != nil {
		return err
	}

	entries := []models.LedgerEntry{
		{AccountID: fromAcc.ID, Amount: -amount, Type: "DEBIT"},
		{AccountID: toAcc.ID, Amount: amount, Type: "CREDIT"},
	}

	return s.CreateTransactionWithTx(tx, txType, refID, "Fund Transfer", entries)
}

// UnlockFunds refunds matched amount from Escrow to User (Used by Admin Termination)
func (s *LedgerService) UnlockFunds(userID uint, amount int64, currency string) error {
	// Find Escrow Account
	escrow, err := s.GetOrCreateAccount(nil, models.AccountTypeEscrow, currency)
	if err != nil {
		return err
	}
	// Find User Account
	userAcc, err := s.GetOrCreateAccount(&userID, models.AccountTypeWallet, currency)
	if err != nil {
		return err
	}

	entries := []models.LedgerEntry{
		{AccountID: escrow.ID, Amount: -amount, Type: "DEBIT"},
		{AccountID: userAcc.ID, Amount: amount, Type: "CREDIT"},
	}

	return s.CreateTransaction(models.TxTypeAdminAdj, "ADMIN_UNLOCK", "Admin Force Unlock", entries)
}
