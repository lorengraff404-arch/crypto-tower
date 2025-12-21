package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

// Battle represents a game battle
type Battle struct {
	ID             int
	Mode           string
	Player1ID      int
	Player2ID      *int
	WagerAmount    float64
	EscrowTx       string
	WinnerID       *int
	ReplayData     json.RawMessage
	ReplayChecksum string
	AntiCheatFlags json.RawMessage
	Status         string
	ResultData     json.RawMessage
	CreatedAt      time.Time
	StartedAt      *time.Time
	CompletedAt    *time.Time
}

// BattleService handles battle logic
type BattleService struct {
	db             *sql.DB
	revenueService *RevenueService
}

// NewBattleService creates a new battle service
func NewBattleService(db *sql.DB, revenueService *RevenueService) *BattleService {
	return &BattleService{
		db:             db,
		revenueService: revenueService,
	}
}

// StartFreeBattle starts a Free/PvE battle
func (s *BattleService) StartFreeBattle(userID int) (*Battle, error) {
	query := `
		INSERT INTO battles (mode, player1_id, status)
		VALUES ('free', $1, 'in_progress')
		RETURNING id, mode, player1_id, status, created_at
	`

	battle := &Battle{}
	err := s.db.QueryRow(query, userID).Scan(
		&battle.ID,
		&battle.Mode,
		&battle.Player1ID,
		&battle.Status,
		&battle.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create free battle: %w", err)
	}

	log.Printf("Free battle started: ID=%d, Player=%d", battle.ID, userID)
	return battle, nil
}

// StartRankedBattle starts a Ranked battle
func (s *BattleService) StartRankedBattle(player1ID, player2ID int) (*Battle, error) {
	query := `
		INSERT INTO battles (mode, player1_id, player2_id, status)
		VALUES ('ranked', $1, $2, 'in_progress')
		RETURNING id, mode, player1_id, player2_id, status, created_at
	`

	battle := &Battle{}
	err := s.db.QueryRow(query, player1ID, player2ID).Scan(
		&battle.ID,
		&battle.Mode,
		&battle.Player1ID,
		&battle.Player2ID,
		&battle.Status,
		&battle.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create ranked battle: %w", err)
	}

	log.Printf("Ranked battle started: ID=%d, P1=%d, P2=%d", battle.ID, player1ID, player2ID)
	return battle, nil
}

// StartWagerBattle starts a Wager battle with escrow
func (s *BattleService) StartWagerBattle(player1ID, player2ID int, wagerAmount float64) (*Battle, error) {
	// Validate wager amount
	var minBet, maxBet float64
	err := s.db.QueryRow("SELECT min_bet, max_bet FROM battle_modes WHERE mode = 'wager'").Scan(&minBet, &maxBet)
	if err != nil {
		return nil, err
	}

	if wagerAmount < minBet || wagerAmount > maxBet {
		return nil, fmt.Errorf("wager amount %.2f outside allowed range [%.2f, %.2f]", wagerAmount, minBet, maxBet)
	}

	// TODO: Create blockchain escrow transaction
	// For now, simulate
	escrowTx := fmt.Sprintf("0x%x", time.Now().UnixNano())

	query := `
		INSERT INTO battles (mode, player1_id, player2_id, wager_amount, escrow_tx, status)
		VALUES ('wager', $1, $2, $3, $4, 'in_progress')
		RETURNING id, mode, player1_id, player2_id, wager_amount, escrow_tx, status, created_at
	`

	battle := &Battle{}
	err = s.db.QueryRow(query, player1ID, player2ID, wagerAmount, escrowTx).Scan(
		&battle.ID,
		&battle.Mode,
		&battle.Player1ID,
		&battle.Player2ID,
		&battle.WagerAmount,
		&battle.EscrowTx,
		&battle.Status,
		&battle.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create wager battle: %w", err)
	}

	log.Printf("Wager battle started: ID=%d, P1=%d, P2=%d, Wager=%.2f TOWER", battle.ID, player1ID, player2ID, wagerAmount)
	return battle, nil
}

// CompleteBattle completes a battle and distributes rewards
func (s *BattleService) CompleteBattle(battleID, winnerID int, replayData []byte) error {
	// Get battle
	battle, err := s.getBattle(battleID)
	if err != nil {
		return err
	}

	if battle.Status != "in_progress" {
		return errors.New("battle is not in progress")
	}

	// Validate winner
	if battle.Player2ID != nil {
		if winnerID != battle.Player1ID && winnerID != *battle.Player2ID {
			return errors.New("winner must be one of the players")
		}
	} else if winnerID != battle.Player1ID {
		return errors.New("winner must be player1 in PvE")
	}

	// Calculate replay checksum
	checksum := sha256.Sum256(replayData)
	checksumStr := hex.EncodeToString(checksum[:])

	// TODO: Anti-cheat validation
	// For now, skip

	// Update battle
	query := `
		UPDATE battles
		SET winner_id = $1, replay_data = $2, replay_checksum = $3, 
		    status = 'completed', completed_at = NOW()
		WHERE id = $4
	`

	_, err = s.db.Exec(query, winnerID, replayData, checksumStr, battleID)
	if err != nil {
		return err
	}

	// Handle rewards based on mode
	switch battle.Mode {
	case "free":
		return s.handleFreeReward(battle, winnerID)
	case "ranked":
		return s.handleRankedReward(battle, winnerID)
	case "wager":
		return s.handleWagerPayout(battle, winnerID)
	}

	return nil
}

// handleFreeReward gives small GTK reward for free mode
func (s *BattleService) handleFreeReward(battle *Battle, winnerID int) error {
	var gtkReward float64
	err := s.db.QueryRow("SELECT gtk_reward FROM battle_modes WHERE mode = 'free'").Scan(&gtkReward)
	if err != nil {
		return err
	}

	// TODO: Transfer GTK to winner
	log.Printf("Free battle reward: %.2f GTK to user %d", gtkReward, winnerID)

	// Record in battle history
	return s.recordBattleHistory(battle.ID, winnerID, 0, 0, 0, gtkReward, 0, true)
}

// handleRankedReward updates ELO and gives GTK reward
func (s *BattleService) handleRankedReward(battle *Battle, winnerID int) error {
	if battle.Player2ID == nil {
		return errors.New("ranked battle requires 2 players")
	}

	loserID := battle.Player1ID
	if winnerID == battle.Player1ID {
		loserID = *battle.Player2ID
	}

	// Get current ELO ratings
	winnerELO, err := s.getOrCreateELO(winnerID)
	if err != nil {
		return err
	}

	loserELO, err := s.getOrCreateELO(loserID)
	if err != nil {
		return err
	}

	// Calculate ELO changes (K-factor = 32)
	expectedWinner := 1.0 / (1.0 + math.Pow(10, float64(loserELO-winnerELO)/400))
	expectedLoser := 1.0 / (1.0 + math.Pow(10, float64(winnerELO-loserELO)/400))

	winnerChange := int(32 * (1.0 - expectedWinner))
	loserChange := int(32 * (0.0 - expectedLoser))

	newWinnerELO := winnerELO + winnerChange
	newLoserELO := loserELO + loserChange

	// Update ELO ratings
	err = s.updateELO(winnerID, newWinnerELO, true)
	if err != nil {
		return err
	}

	err = s.updateELO(loserID, newLoserELO, false)
	if err != nil {
		return err
	}

	// Get GTK reward
	var gtkReward float64
	err = s.db.QueryRow("SELECT gtk_reward FROM battle_modes WHERE mode = 'ranked'").Scan(&gtkReward)
	if err != nil {
		return err
	}

	// TODO: Transfer GTK to winner
	log.Printf("Ranked battle reward: %.2f GTK to user %d, ELO: %d -> %d (+%d)",
		gtkReward, winnerID, winnerELO, newWinnerELO, winnerChange)

	// Record in battle history
	err = s.recordBattleHistory(battle.ID, winnerID, winnerELO, newWinnerELO, winnerChange, gtkReward, 0, true)
	if err != nil {
		return err
	}

	return s.recordBattleHistory(battle.ID, loserID, loserELO, newLoserELO, loserChange, 0, 0, false)
}

// handleWagerPayout releases escrow to winner
func (s *BattleService) handleWagerPayout(battle *Battle, winnerID int) error {
	// TODO: Release escrow from blockchain
	// For now, simulate
	log.Printf("Wager payout: %.2f TOWER to user %d (escrow: %s)",
		battle.WagerAmount*2, winnerID, battle.EscrowTx)

	// Record in battle history
	return s.recordBattleHistory(battle.ID, winnerID, 0, 0, 0, 0, battle.WagerAmount*2, true)
}

// getBattle retrieves a battle by ID
func (s *BattleService) getBattle(battleID int) (*Battle, error) {
	query := `
		SELECT id, mode, player1_id, player2_id, wager_amount, escrow_tx, winner_id,
		       status, created_at
		FROM battles
		WHERE id = $1
	`

	battle := &Battle{}
	err := s.db.QueryRow(query, battleID).Scan(
		&battle.ID,
		&battle.Mode,
		&battle.Player1ID,
		&battle.Player2ID,
		&battle.WagerAmount,
		&battle.EscrowTx,
		&battle.WinnerID,
		&battle.Status,
		&battle.CreatedAt,
	)

	return battle, err
}

// getOrCreateELO gets or creates ELO rating for a user
func (s *BattleService) getOrCreateELO(userID int) (int, error) {
	var elo int
	err := s.db.QueryRow("SELECT elo_rating FROM player_elo WHERE user_id = $1", userID).Scan(&elo)

	if err == sql.ErrNoRows {
		// Create new ELO record with default 1200
		_, err = s.db.Exec("INSERT INTO player_elo (user_id, elo_rating) VALUES ($1, 1200)", userID)
		return 1200, err
	}

	return elo, err
}

// updateELO updates a player's ELO rating
func (s *BattleService) updateELO(userID, newELO int, won bool) error {
	query := `
		UPDATE player_elo
		SET elo_rating = $1,
		    games_played = games_played + 1,
		    wins = wins + CASE WHEN $2 THEN 1 ELSE 0 END,
		    losses = losses + CASE WHEN $2 THEN 0 ELSE 1 END,
		    win_streak = CASE WHEN $2 THEN win_streak + 1 ELSE 0 END,
		    highest_elo = GREATEST(highest_elo, $1),
		    updated_at = NOW()
		WHERE user_id = $3
	`

	_, err := s.db.Exec(query, newELO, won, userID)
	return err
}

// recordBattleHistory records battle result in history
func (s *BattleService) recordBattleHistory(battleID, playerID, eloBefore, eloAfter, eloChange int, gtkReward, towerReward float64, isWinner bool) error {
	query := `
		INSERT INTO battle_history 
		(battle_id, player_id, elo_before, elo_after, elo_change, gtk_reward, tower_reward, is_winner)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.db.Exec(query, battleID, playerID, eloBefore, eloAfter, eloChange, gtkReward, towerReward, isWinner)
	return err
}
