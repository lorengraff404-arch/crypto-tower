package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// AntiCheatFlag represents a detected cheat flag
type AntiCheatFlag struct {
	ID          int
	BattleID    int
	UserID      int
	FlagType    string
	Severity    string
	Details     json.RawMessage
	AutoFlagged bool
	Reviewed    bool
	CreatedAt   time.Time
}

// GameAction represents a player action in a battle
type GameAction struct {
	Timestamp  time.Time
	ActionType string
	Data       map[string]interface{}
}

// AntiCheatService handles cheat detection
type AntiCheatService struct {
	db *sql.DB
}

// NewAntiCheatService creates a new anti-cheat service
func NewAntiCheatService(db *sql.DB) *AntiCheatService {
	return &AntiCheatService{db: db}
}

// ValidateBattle performs comprehensive validation on a battle
func (s *AntiCheatService) ValidateBattle(battleID int, player1ID, player2ID int, replayData []byte) ([]AntiCheatFlag, error) {
	var flags []AntiCheatFlag

	// Parse replay data
	var actions []GameAction
	err := json.Unmarshal(replayData, &actions)
	if err != nil {
		return nil, fmt.Errorf("invalid replay data: %w", err)
	}

	// Run all detection systems
	botFlags1, err := s.DetectBot(player1ID, actions)
	if err != nil {
		return nil, err
	}
	flags = append(flags, botFlags1...)

	if player2ID > 0 {
		botFlags2, err := s.DetectBot(player2ID, actions)
		if err != nil {
			return nil, err
		}
		flags = append(flags, botFlags2...)

		// Check collusion between players
		collusionFlags, err := s.DetectCollusion(battleID, player1ID, player2ID)
		if err != nil {
			return nil, err
		}
		flags = append(flags, collusionFlags...)
	}

	// Save all flags to database
	for i := range flags {
		flags[i].BattleID = battleID
		err = s.saveFlag(&flags[i])
		if err != nil {
			log.Printf("Failed to save anti-cheat flag: %v", err)
		}
	}

	return flags, nil
}

// DetectBot detects bot-like behavior
func (s *AntiCheatService) DetectBot(userID int, actions []GameAction) ([]AntiCheatFlag, error) {
	var flags []AntiCheatFlag

	if len(actions) == 0 {
		return flags, nil
	}

	// Calculate APM (Actions Per Minute)
	duration := actions[len(actions)-1].Timestamp.Sub(actions[0].Timestamp).Minutes()
	if duration <= 0 {
		duration = 1
	}
	apm := float64(len(actions)) / duration

	// Superhuman APM threshold (500 APM is very high for tower defense)
	if apm > 500 {
		details, _ := json.Marshal(map[string]interface{}{
			"apm":              apm,
			"actions":          len(actions),
			"duration_minutes": duration,
		})

		flags = append(flags, AntiCheatFlag{
			UserID:      userID,
			FlagType:    "bot_detection",
			Severity:    "high",
			Details:     details,
			AutoFlagged: true,
		})
	}

	// Check timing consistency (bots have very consistent timing)
	if s.isTimingTooConsistent(actions) {
		details, _ := json.Marshal(map[string]interface{}{
			"reason": "inhuman_timing_consistency",
			"apm":    apm,
		})

		flags = append(flags, AntiCheatFlag{
			UserID:      userID,
			FlagType:    "bot_detection",
			Severity:    "medium",
			Details:     details,
			AutoFlagged: true,
		})
	}

	// Check for impossible reaction times
	if s.hasImpossibleReactions(actions) {
		details, _ := json.Marshal(map[string]interface{}{
			"reason": "impossible_reaction_time",
		})

		flags = append(flags, AntiCheatFlag{
			UserID:      userID,
			FlagType:    "bot_detection",
			Severity:    "high",
			Details:     details,
			AutoFlagged: true,
		})
	}

	return flags, nil
}

// DetectCollusion detects collusion between players
func (s *AntiCheatService) DetectCollusion(battleID, player1ID, player2ID int) ([]AntiCheatFlag, error) {
	var flags []AntiCheatFlag

	// Check if same IP address
	if s.sameIP(player1ID, player2ID) {
		details, _ := json.Marshal(map[string]interface{}{
			"reason":  "same_ip_address",
			"player1": player1ID,
			"player2": player2ID,
		})

		flags = append(flags, AntiCheatFlag{
			UserID:      player1ID,
			FlagType:    "collusion",
			Severity:    "high",
			Details:     details,
			AutoFlagged: true,
		})
	}

	// Check wallet connection patterns
	if s.connectedWallets(player1ID, player2ID) {
		details, _ := json.Marshal(map[string]interface{}{
			"reason":  "connected_wallets",
			"player1": player1ID,
			"player2": player2ID,
		})

		flags = append(flags, AntiCheatFlag{
			UserID:      player1ID,
			FlagType:    "collusion",
			Severity:    "medium",
			Details:     details,
			AutoFlagged: true,
		})
	}

	// Check suspicious win rate patterns
	if s.suspiciousWinRate(player1ID, player2ID) {
		details, _ := json.Marshal(map[string]interface{}{
			"reason":  "unnatural_win_rate_pattern",
			"player1": player1ID,
			"player2": player2ID,
		})

		flags = append(flags, AntiCheatFlag{
			UserID:      player1ID,
			FlagType:    "collusion",
			Severity:    "high",
			Details:     details,
			AutoFlagged: true,
		})
	}

	return flags, nil
}

// isTimingTooConsistent checks if action timing is unnaturally consistent
func (s *AntiCheatService) isTimingTooConsistent(actions []GameAction) bool {
	if len(actions) < 10 {
		return false
	}

	// Calculate variance in action intervals
	var intervals []float64
	for i := 1; i < len(actions); i++ {
		interval := actions[i].Timestamp.Sub(actions[i-1].Timestamp).Seconds()
		intervals = append(intervals, interval)
	}

	// Calculate standard deviation
	var sum, mean, variance float64
	for _, interval := range intervals {
		sum += interval
	}
	mean = sum / float64(len(intervals))

	for _, interval := range intervals {
		variance += (interval - mean) * (interval - mean)
	}
	variance /= float64(len(intervals))
	stdDev := variance

	// If std deviation is very low, timing is too consistent (bot-like)
	// Humans have natural variation in timing
	return stdDev < 0.05 && mean > 0
}

// hasImpossibleReactions checks for superhuman reaction times
func (s *AntiCheatService) hasImpossibleReactions(actions []GameAction) bool {
	// Human reaction time is typically 150-300ms
	minHumanReaction := 100 * time.Millisecond

	for i := 1; i < len(actions); i++ {
		interval := actions[i].Timestamp.Sub(actions[i-1].Timestamp)
		// If consecutive actions are extremely fast, it's suspicious
		if interval < minHumanReaction && interval > 0 {
			return true
		}
	}

	return false
}

// sameIP checks if two players are using the same IP address
func (s *AntiCheatService) sameIP(player1ID, player2ID int) bool {
	// TODO: Implement actual IP checking from session logs
	// Query session_logs or user_sessions table for recent IP addresses
	// Example query:
	// SELECT ip_address FROM session_logs WHERE user_id IN ($1, $2)
	// ORDER BY created_at DESC LIMIT 2

	// Log for future implementation tracking
	if player1ID > 0 && player2ID > 0 {
		// Placeholder: When implemented, this will check IP match
		// log.Printf("Checking IP match for players %d and %d", player1ID, player2ID)
	}

	return false
}

// connectedWallets checks if wallets have transaction history
func (s *AntiCheatService) connectedWallets(player1ID, player2ID int) bool {
	// TODO: Implement blockchain wallet analysis
	// Check if wallets have sent tokens to each other
	// This requires:
	// 1. Get wallet addresses for both players
	// 2. Query blockchain for transactions between these addresses
	// 3. Check transaction history for suspicious patterns

	// Log for future implementation tracking
	if player1ID > 0 && player2ID > 0 {
		// Placeholder: When implemented, this will analyze wallet connections
		// log.Printf("Checking wallet connections for players %d and %d", player1ID, player2ID)
	}

	return false
}

// suspiciousWinRate checks for unnatural win/loss patterns
func (s *AntiCheatService) suspiciousWinRate(player1ID, player2ID int) bool {
	// Query battle history between these two players
	query := `
		SELECT 
			COUNT(*) as total_battles,
			SUM(CASE WHEN winner_id = $1 THEN 1 ELSE 0 END) as player1_wins
		FROM battles
		WHERE (player1_id = $1 AND player2_id = $2) 
		   OR (player1_id = $2 AND player2_id = $1)
	`

	var totalBattles, player1Wins int
	err := s.db.QueryRow(query, player1ID, player2ID).Scan(&totalBattles, &player1Wins)
	if err != nil {
		return false
	}

	// If they've played multiple times and one always wins, suspicious
	if totalBattles >= 5 {
		winRate := float64(player1Wins) / float64(totalBattles)
		// 100% or 0% win rate over 5+ games is suspicious
		return winRate == 1.0 || winRate == 0.0
	}

	return false
}

// saveFlag saves an anti-cheat flag to the database
func (s *AntiCheatService) saveFlag(flag *AntiCheatFlag) error {
	query := `
		INSERT INTO anti_cheat_flags 
		(battle_id, user_id, flag_type, severity, details, auto_flagged)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return s.db.QueryRow(query,
		flag.BattleID,
		flag.UserID,
		flag.FlagType,
		flag.Severity,
		flag.Details,
		flag.AutoFlagged,
	).Scan(&flag.ID, &flag.CreatedAt)
}

// GetFlagsByBattle retrieves all flags for a battle
func (s *AntiCheatService) GetFlagsByBattle(battleID int) ([]AntiCheatFlag, error) {
	query := `
		SELECT id, battle_id, user_id, flag_type, severity, details, 
		       auto_flagged, reviewed, created_at
		FROM anti_cheat_flags
		WHERE battle_id = $1
		ORDER BY severity DESC, created_at DESC
	`

	rows, err := s.db.Query(query, battleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flags []AntiCheatFlag
	for rows.Next() {
		var flag AntiCheatFlag
		err := rows.Scan(
			&flag.ID,
			&flag.BattleID,
			&flag.UserID,
			&flag.FlagType,
			&flag.Severity,
			&flag.Details,
			&flag.AutoFlagged,
			&flag.Reviewed,
			&flag.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

// GetFlagsByUser retrieves all flags for a user
func (s *AntiCheatService) GetFlagsByUser(userID int) ([]AntiCheatFlag, error) {
	query := `
		SELECT id, battle_id, user_id, flag_type, severity, details, 
		       auto_flagged, reviewed, created_at
		FROM anti_cheat_flags
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flags []AntiCheatFlag
	for rows.Next() {
		var flag AntiCheatFlag
		err := rows.Scan(
			&flag.ID,
			&flag.BattleID,
			&flag.UserID,
			&flag.FlagType,
			&flag.Severity,
			&flag.Details,
			&flag.AutoFlagged,
			&flag.Reviewed,
			&flag.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}
