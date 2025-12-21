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
	battleEngine  *services.BattleEngine
	ledgerService *services.LedgerService
}

// NewWagerHandler creates new wager handler
func NewWagerHandler() *WagerHandler {
	return &WagerHandler{
		battleEngine:  services.NewBattleEngine(),
		ledgerService: services.NewLedgerService(),
	}
}

// StartWagerRequest represents wager battle request
type StartWagerRequest struct {
	TeamID      uint  `json:"team_id" binding:"required"`
	WagerAmount int64 `json:"wager_amount" binding:"required,min=100"` // Minimum 100 GTK
	OpponentID  *uint `json:"opponent_id"`                             // Optional: challenge specific user
}

// StartWager initiates high-stakes wager battle
// POST /api/v1/battle/wager
// Security: Ledger Reserve (Funds Hold), Durable Transaction
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

	// 1. SECURITY: Minimum wager validation
	if req.WagerAmount < 100 || req.WagerAmount > 100000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wager amount"})
		return
	}

	// 2. TEAM VALIDATION
	var team models.Team
	if err := db.DB.Where("id = ? AND owner_id = ?", req.TeamID, userID).Preload("Members.Character").First(&team).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Team not found"})
		return
	}

	// Check durability
	for _, member := range team.Members {
		if !member.IsBackup {
			if member.Character.Durability < 15 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Low durability", "char_id": member.Character.ID})
				return
			}
			if member.Character.IsDead {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Character is dead", "char_id": member.Character.ID})
				return
			}
		}
	}

	// 3. ATOMIC TRANSACTION: Create Battle + Reserve Funds
	// This ensures funds aren't lost if Battle record fails, and vice versa
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// A. Create Battle Record (Queue Entry)
		newBattle := models.Battle{
			BattleType: "wager",
			Status:     "SEARCHING", // Waiting for opponent
			Player1ID:  userID,
			BetAmount:  req.WagerAmount,
		}
		if req.OpponentID != nil {
			newBattle.Player2ID = *req.OpponentID // Targeted challenge
		}

		if err := tx.Create(&newBattle).Error; err != nil {
			return err
		}

		// B. Check & Reserve Funds (Ledger)
		userAcc, err := h.ledgerService.GetOrCreateAccount(&userID, models.AccountTypeWallet, "GTK")
		if err != nil {
			return err
		}

		escrowAcc, err := h.ledgerService.GetOrCreateAccount(nil, models.AccountTypeEscrow, "GTK")
		if err != nil {
			return err
		}

		// Manual Balance Check (Ledger Check)
		if userAcc.Balance < req.WagerAmount {
			return errors.New("insufficient GTK balance")
		}

		entries := []models.LedgerEntry{
			{AccountID: userAcc.ID, Amount: -req.WagerAmount, Type: "DEBIT"},
			{AccountID: escrowAcc.ID, Amount: req.WagerAmount, Type: "CREDIT"},
		}

		// Use Transactional Ledger Call
		txType := models.TxTypeWagerEnter
		refID := fmt.Sprintf("wager_%d_p1", newBattle.ID)
		desc := fmt.Sprintf("Wager Entry for Battle #%d", newBattle.ID)

		if err := h.ledgerService.CreateTransactionWithTx(tx, txType, refID, desc, entries); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wager failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Funds reserved. Searching for match...",
		"wager":   req.WagerAmount,
		"status":  "in_queue",
	})
}
