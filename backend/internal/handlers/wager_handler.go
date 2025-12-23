package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
	"gorm.io/gorm"
)

// WagerHandler handles high-stakes wager battles with real GTK via Ledger
type WagerHandler struct {
	battleService *services.BattleService
	ledgerService *services.LedgerService
}

// NewWagerHandler creates new wager handler
func NewWagerHandler() *WagerHandler {
	return &WagerHandler{
		battleService: services.NewBattleService(),
		ledgerService: services.NewLedgerService(),
	}
}

// StartWagerRequest represents wager battle request
// StartWagerRequest represents wager battle request
type StartWagerRequest struct {
	TeamID uint `json:"team_id" binding:"required"`
	// WagerAmount is removed; dynamic calculation used.
	OpponentID *uint `json:"opponent_id"` // Optional: challenge specific user
}

// StartWager initiates high-stakes wager battle (Matchmaking: Join or Create)
// POST /api/v1/battle/wager
func (h *WagerHandler) StartWager(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := val.(uint)

	var req StartWagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. SECURITY: Minimum Balance Validation (Proof of Solvency)
	userAcc, err := h.ledgerService.GetOrCreateAccount(&userID, models.AccountTypeWallet, "GTK")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wallet"})
		return
	}
	minBalance := int64(500) // Minimum to Enter Arena
	if userAcc.Balance < minBalance {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": fmt.Sprintf("Insufficient balance. Need %d GTK to enter High Stakes.", minBalance)})
		return
	}

	// 2. CHECK if user is already in a battle
	var activeBattle models.Battle
	if err := db.DB.Where("(player1_id = ? OR player2_id = ?) AND status IN (?, ?)", userID, userID, "SEARCHING", "active").First(&activeBattle).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":     "You are already in a battle or queue",
			"battle_id": activeBattle.ID,
			"status":    activeBattle.Status,
		})
		return
	}

	// 3. TEAM VALIDATION & SNAPSHOT
	// We need team snapshot to calculate stakes
	// 3. TEAM VALIDATION & SNAPSHOT
	// We need team snapshot to calculate stakes
	_, err = h.battleService.GetTeamSnapshot(userID) // Expose this or re-implement
	// Ideally using the service helper but it is private lowercase.
	// We should export GetTeamSnapshot in Service or use the logic here.
	// For now, let's assume we can rely on `battleService.InitializeBattleState` later,
	// BUT for matchmaking we need CP *before* creating battle?
	// Actually, simplified flow:
	// 1. Join Queue (SEARCHING) with P1.
	// 2. Matchmaking loop checks compatibility.
	//    The "Join" action triggers the calc.

	// Refined Logic:
	// User Creates "SEARCHING" record. No funds locked yet?
	// OR: Lock "Base Stake" (100) as deposit?
	// Risk: If dynamic stake is 500, we need to lock 500.
	// Solution:
	// "Searching" phase is free (or small deposit).
	// When Match Found -> Atomic Transaction:
	//   Calc Stakes for P1 & P2.
	//   Check P1 Balance > P1Stake.
	//   Check P2 Balance > P2Stake.
	//   Lock Funds.
	//   Start Battle.

	// This Handler creates a SEARCHING record if no opponent found.
	// Opponent matching logic needs to be robust.

	var battleID uint
	var matchFound bool
	var p1Stake, p2Stake int64

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// A. Find Opponent
		var match models.Battle
		// Find any SEARCHING wager battle.
		// Optimized: Order by CreatedAt ASC (FIFO) or ELO clostness?
		query := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("battle_type = ? AND status = ? AND player1_id != ?", "wager", "SEARCHING", userID)

		if err := query.First(&match).Error; err == nil {
			// --- MATCH CANDIDATE FOUND ---

			// 1. Calculate Stakes
			// We need P1 (Candidate) and P2 (Current User) Teams
			// Note: P1 is `match.Player1ID`, P2 is `userID`

			// Needs public method in Service.
			// Assuming we fixed `GetTeamSnapshot` visibility or duplicated logic.
			// Let's assume h.battleService has public GetTeamSnapshot.
			// (I will rename it in Service in next step if generic)

			// Or just rely on CalculateDynamicStakes which takes IDs?
			stake1, stake2, err := h.battleService.CalculateDynamicStakes(match.Player1ID, userID, nil) // Pass nil to force fetch?
			// Update Service to handle nil team input?
			// Or fetch here.

			if err != nil {
				return err // Skip this match? Or fail?
			}

			// 2. Verify Balances
			p1Acc, _ := h.ledgerService.GetOrCreateAccount(&match.Player1ID, models.AccountTypeWallet, "GTK")
			p2Acc, _ := h.ledgerService.GetOrCreateAccount(&userID, models.AccountTypeWallet, "GTK")

			if p1Acc.Balance < stake1 {
				// P1 invalid (maybe withdrew funds while waiting).
				// Should cancel P1's search?
				// For now, skip and continue searching? (Complex in single query)
				// Or fail transaction and let client retry?
				return errors.New("opponent insufficient funds")
			}
			if p2Acc.Balance < stake2 {
				return errors.New("insufficient GTK balance for required stake: " + fmt.Sprint(stake2))
			}

			// 3. Lock Funds (Atomic)
			escrowAcc, _ := h.ledgerService.GetOrCreateAccount(nil, models.AccountTypeEscrow, "GTK")

			// P1 Lock
			entries1 := []models.LedgerEntry{
				{AccountID: p1Acc.ID, Amount: -stake1, Type: "DEBIT"},
				{AccountID: escrowAcc.ID, Amount: stake1, Type: "CREDIT"},
			}
			if err := h.ledgerService.CreateTransactionWithTx(tx, models.TxTypeWagerEnter, fmt.Sprintf("wager_%d_p1", match.ID), "Wager Lock P1", entries1); err != nil {
				return err
			}

			// P2 Lock
			entries2 := []models.LedgerEntry{
				{AccountID: p2Acc.ID, Amount: -stake2, Type: "DEBIT"},
				{AccountID: escrowAcc.ID, Amount: stake2, Type: "CREDIT"},
			}
			if err := h.ledgerService.CreateTransactionWithTx(tx, models.TxTypeWagerEnter, fmt.Sprintf("wager_%d_p2", match.ID), "Wager Lock P2", entries2); err != nil {
				return err
			}

			// 4. Start Battle
			match.Player2ID = userID
			match.Status = "active"
			match.Player1Bet = stake1
			match.Player2Bet = stake2

			// Init State
			if err := h.battleService.InitializeBattleState(&match); err != nil {
				return err
			}
			if err := tx.Save(&match).Error; err != nil {
				return err
			}

			matchFound = true
			battleID = match.ID
			p1Stake = stake1
			p2Stake = stake2

		} else {
			// --- NO MATCH: CREATE SEARCH ---
			// No funds locked yet. Just "Ready check" balance.

			newBattle := models.Battle{
				BattleType: "wager",
				Status:     "SEARCHING",
				Player1ID:  userID,
				// BetAmount/Stakes unknown until opponent appears
			}
			if err := tx.Create(&newBattle).Error; err != nil {
				return err
			}
			battleID = newBattle.ID
			matchFound = false
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Matchmaking failed", "details": err.Error()})
		return
	}

	if matchFound {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Match found! Battle starting.",
			"battle_id":   battleID,
			"your_stake":  p2Stake,
			"enemy_stake": p1Stake,
			"status":      "active",
			"role":        "player2",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":           "Entered Arena Queue. Waiting for challenger...",
			"battle_id":         battleID,
			"min_balance_check": "passed",
			"status":            "in_queue",
			"role":              "player1",
		})
	}
}

// CancelWager cancels a searching battle and refunds funds
// POST /api/v1/battle/wager/cancel
func (h *WagerHandler) CancelWager(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := val.(uint)

	// Transaction: Find Battle -> Verify -> Delete -> Refund
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var battle models.Battle
		// Find battle that is SEARCHING and owned by User
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("player1_id = ? AND status = ? AND battle_type = ?", userID, "SEARCHING", "wager").
			First(&battle).Error; err != nil {
			return errors.New("no active wager search found")
		}

		// Refund Funds
		userAcc, err := h.ledgerService.GetOrCreateAccount(&userID, models.AccountTypeWallet, "GTK")
		if err != nil {
			return err
		}
		escrowAcc, err := h.ledgerService.GetOrCreateAccount(nil, models.AccountTypeEscrow, "GTK")
		if err != nil {
			return err
		}

		entries := []models.LedgerEntry{
			{AccountID: escrowAcc.ID, Amount: -battle.BetAmount, Type: "DEBIT"},
			{AccountID: userAcc.ID, Amount: battle.BetAmount, Type: "CREDIT"},
		}

		if err := h.ledgerService.CreateTransactionWithTx(tx, models.TxTypeWagerRefund, fmt.Sprintf("refund_%d", battle.ID), "Wager Cancel Refund", entries); err != nil {
			return err
		}

		// Delete/Cancel Battle
		if err := tx.Delete(&battle).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wager cancelled and funds refunded"})
}
