package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

type BattleService struct {
	engine        *BattleEngine
	ledger        *LedgerService
	skillService  *SkillActivationService
	statusService *StatusEffectService
}

func NewBattleService() *BattleService {
	return &BattleService{
		engine:        NewBattleEngine(),
		ledger:        NewLedgerService(),
		skillService:  NewSkillActivationService(),
		statusService: NewStatusEffectService(),
	}
}

// CreatePvEBattle initializes a new PvE battle session
func (s *BattleService) CreatePvEBattle(player1ID uint, mode string) (*models.Battle, error) {
	// Create Battle Record
	battle := models.Battle{
		BattleType:          mode,
		Status:              "active",
		Player1ID:           player1ID,
		Player2ID:           player1ID, // PvE usually has same ID or System ID
		CurrentTurnPlayerID: player1ID, // P1 starts
		TurnNumber:          1,
		CreatedAt:           time.Now(),
	}

	// Initialize State
	if err := s.InitializeBattleState(&battle); err != nil {
		return nil, err
	}

	if err := db.DB.Create(&battle).Error; err != nil {
		return nil, err
	}

	return &battle, nil
}

// InitializeBattleState generates team snapshots for the battle participants
func (s *BattleService) InitializeBattleState(battle *models.Battle) error {
	// 1. Snapshot Player 1
	p1Team, err := s.GetTeamSnapshot(battle.Player1ID)
	if err != nil {
		return fmt.Errorf("failed to get player 1 team: %w", err)
	}
	p1Json, _ := json.Marshal(p1Team)
	battle.PlayerStateP1 = string(p1Json)

	// Calculate Dynamic Stakes if Wager
	if battle.BattleType == "wager" && battle.Status == "pending" { // Only calc if new
		p1Stake, p2Stake, err := s.CalculateDynamicStakes(battle.Player1ID, battle.Player2ID, p1Team)
		if err == nil {
			battle.Player1Bet = p1Stake
			battle.Player2Bet = p2Stake
		}
	}

	// 2. Snapshot Player 2 (or AI)
	var p2Team []models.BattleParticipant
	if strings.Contains(battle.BattleType, "PVE") {
		p2Team = s.generateAITeam(battle.Player1ID, battle.BattleType)
	} else {
		// PvP / Wager
		p2Team, err = s.GetTeamSnapshot(battle.Player2ID)
		if err != nil {
			return fmt.Errorf("failed to get player 2 team: %w", err)
		}
	}
	p2Json, _ := json.Marshal(p2Team)
	battle.PlayerStateP2 = string(p2Json)

	return nil
}

// Helper: Get Team Snapshot
func (s *BattleService) GetTeamSnapshot(userID uint) ([]models.BattleParticipant, error) {
	// Fetch Active Team
	var team models.Team
	if err := db.DB.Preload("Members.Character").
		Where("user_id = ? AND is_active = true", userID).
		First(&team).Error; err != nil {
		return nil, errors.New("no active team found")
	}

	var participants []models.BattleParticipant
	for _, member := range team.Members {
		if member.Character.ID != 0 {
			participants = append(participants, *s.toParticipant(&member.Character))
		}
	}
	return participants, nil
}

// Helper: Generate AI Team
func (s *BattleService) generateAITeam(_ uint, _ string) []models.BattleParticipant {
	// Simplified: Create 1 Dummy Enemy
	// Real Logic: Scale based on Player Level or Raid Difficulty
	return []models.BattleParticipant{
		{
			CharacterName: "Goblin Scout",
			MaxHP:         100,
			CurrentHP:     100,
			Attack:        15,
			Defense:       5,
			Speed:         10,
			IsActive:      true,
		},
	}
}

// ProcessTurn executes a turn in a PvP battle
// Uses BattleEngine to calculate damage and update state
func (s *BattleService) ProcessTurn(battleID uint, userID uint, actionData map[string]interface{}) (*models.Battle, error) {
	var battle models.Battle
	if err := db.DB.Preload("Player1").Preload("Player2").First(&battle, battleID).Error; err != nil {
		return nil, errors.New("battle not found")
	}

	if battle.Status != "active" {
		return nil, errors.New("battle is not active")
	}

	// 1. Verify Turn
	if battle.CurrentTurnPlayerID != userID {
		return nil, errors.New("not your turn")
	}

	// 2. Parse Action
	actionType, ok := actionData["action"].(string)
	if !ok {
		return nil, errors.New("invalid action format")
	}

	// 3. Load Participants
	charIDVal, ok := actionData["character_id"].(float64)
	if !ok {
		return nil, errors.New("missing character_id")
	}
	charID := uint(charIDVal)

	targetIDVal, ok := actionData["target_id"].(float64)
	if !ok {
		return nil, errors.New("missing target_id")
	}
	targetID := uint(targetIDVal)

	// Fetch characters
	var attacker, defender models.Character
	if err := db.DB.First(&attacker, charID).Error; err != nil {
		return nil, errors.New("attacker not found")
	}
	if err := db.DB.First(&defender, targetID).Error; err != nil {
		return nil, errors.New("defender not found")
	}

	// SECURITY: Verify ownership and state
	if attacker.OwnerID != userID {
		return nil, errors.New("you do not own the attacking character")
	}
	if attacker.IsFainted || attacker.IsDead {
		return nil, errors.New("character cannot act")
	}

	// --- TURN START: Process Status Effects (DoT) on Attacker ---
	// Only process effects if it's the start of their turn logic
	dotDamage, expiredEffects, err := s.statusService.ProcessTurnEffects(attacker.ID)
	if err != nil {
		fmt.Printf("Error processing effects: %v\n", err)
	}

	// Refresh attacker from DB after status effects (HP might have changed)
	if dotDamage > 0 {
		db.DB.First(&attacker, charID) // Reload
		if attacker.IsFainted {
			// If died from Poison, turn ends immediately? Or prevent action?
			return nil, errors.New("character fainted from status effects")
		}
	}
	// Reduce Cooldowns via SkillService
	s.skillService.ReduceCooldowns(attacker.ID)
	// Regenerate Mana
	s.skillService.RegenerateMana(attacker.ID)
	// Reload attacker again to get fresh Mana/CDs
	db.DB.First(&attacker, charID)

	var logMsg string

	switch actionType {
	case "skill":
		skillIDVal, ok := actionData["skill_id"].(float64)
		if !ok {
			return nil, errors.New("missing skill_id")
		}
		skillID := uint(skillIDVal)

		req := SkillActivationRequest{
			CharacterID: attacker.ID,
			AbilityID:   skillID,
			TargetID:    defender.ID,
			BattleID:    battle.ID,
			TurnNumber:  battle.TurnNumber,
		}

		// Activate Skill (Handles Mana, CD, Buffs, DB Save for Attacker)
		result, err := s.skillService.ActivateSkill(req)
		if err != nil {
			return nil, err
		}

		// Apply Damage to Defender
		if result.Damage > 0 {
			defender.CurrentHP -= result.Damage
			if defender.CurrentHP < 0 {
				defender.CurrentHP = 0
			}
			if defender.CurrentHP == 0 {
				defender.IsFainted = true
			}
		}

		// Apply Healing (if any target - usually self for verify simplicity)
		if result.Healing > 0 {
			// If target was self, reload attacker to see healing?
			// ActivateSkill only calculates, doesn't apply healing to DB?
			// Checking ActivateSkill... it only returns Result.
			// So WE must apply healing.
			if targetID == attacker.ID {
				// Reload attacker (ActivateSkill saved mana deduction)
				db.DB.First(&attacker, charID)
				attacker.CurrentHP += result.Healing
				if attacker.CurrentHP > attacker.BaseHP {
					attacker.CurrentHP = attacker.BaseHP
				}
				db.DB.Save(&attacker)
			} else {
				defender.CurrentHP += result.Healing
				if defender.CurrentHP > defender.BaseHP {
					defender.CurrentHP = defender.BaseHP
				}
			}
		}

		db.DB.Save(&defender)
		logMsg = result.Message

	case "attack":
		// Basic Attack (Physical, No Mana, No CD)
		// Use Engine or create dummy ability
		// Adapting to simple physical hit
		pAttacker := s.toParticipant(&attacker)
		pDefender := s.toParticipant(&defender)

		// Basic Attack Ability
		ability := models.Ability{Name: "Attack", Damage: 10, DamageType: "physical", Element: "Normal"}

		// Use Engine for logic
		res, err := s.engine.ExecuteAbility(pAttacker, pDefender, ability)
		if err != nil {
			return nil, err
		}

		defender.CurrentHP = pDefender.CurrentHP
		if pDefender.IsFainted {
			defender.IsFainted = true
		}
		db.DB.Save(&defender) // Only save defender. Attacker not changed in basic attack (no mana)
		logMsg = res.Message

	case "item":
		// Item Usage Logic
		itemIDVal, ok := actionData["item_id"].(float64)
		if !ok {
			return nil, errors.New("missing item_id")
		}
		itemID := uint(itemIDVal)

		// 1. Verify Inventory
		var inventory models.UserInventory
		// Start Transaction for Item consumption
		err := db.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("user_id = ? AND item_id = ?", userID, itemID).First(&inventory).Error; err != nil {
				return errors.New("item not owned or empty")
			}
			if inventory.Quantity <= 0 {
				return errors.New("item quantity is 0")
			}

			// 2. Fetch Item Details
			var shopItem models.ShopItem
			if err := tx.First(&shopItem, itemID).Error; err != nil {
				return errors.New("item details not found")
			}

			// 3. Apply Effect
			if shopItem.EffectType == "heal_hp" {
				heal := shopItem.EffectValue
				attacker.CurrentHP += heal
				if attacker.CurrentHP > attacker.BaseHP {
					attacker.CurrentHP = attacker.BaseHP
				}
				logMsg = fmt.Sprintf("Used %s. Healed %d HP.", shopItem.Name, heal)
			} else {
				logMsg = fmt.Sprintf("Used %s.", shopItem.Name)
			}

			// 4. Consume Item
			inventory.Quantity--
			if inventory.Quantity == 0 {
				if err := tx.Delete(&inventory).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Save(&inventory).Error; err != nil {
					return err
				}
			}

			// Save Character
			if err := tx.Save(&attacker).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

	case "switch":
		logMsg = "Switched character!"
		// Placeholder for active slot logic

	default:
		return nil, errors.New("unknown action type")
	}

	// 4. Update Battle State

	winnerID := uint(0)
	gameEnded := false

	if defender.IsFainted {
		// Check if team is wiped
		var count int64
		db.DB.Model(&models.Character{}).
			Where("owner_id = ? AND is_fainted = false", defender.OwnerID).
			Count(&count)

		if count == 0 {
			winnerID = attacker.OwnerID
			gameEnded = true
		}
	}

	if gameEnded {
		s.CompleteBattle(battle.ID, winnerID, "")
		db.DB.First(&battle, battleID) // Reload
	} else {
		// Toggle Turn
		if battle.CurrentTurnPlayerID == battle.Player1ID {
			battle.CurrentTurnPlayerID = battle.Player2ID
		} else {
			battle.CurrentTurnPlayerID = battle.Player1ID
		}
		battle.TurnNumber++
	}

	// Include Status Info in Log
	if len(expiredEffects) > 0 {
		logMsg += fmt.Sprintf(" (Expired: %v)", expiredEffects)
	}
	if dotDamage > 0 {
		logMsg += fmt.Sprintf(" (Took %d DoT damage)", dotDamage)
	}
	if gameEnded {
		logMsg += " BATTLE ENDED!"
	}

	newState := map[string]interface{}{
		"last_action": actionType,
		"attacker":    charID,
		"target":      targetID,
		"log":         logMsg,
		"defender_hp": defender.CurrentHP,
		"game_ended":  gameEnded,
		"winner_id":   winnerID,
	}
	stateBytes, _ := json.Marshal(newState)
	battle.LastTurnData = string(stateBytes)

	// Save
	db.DB.Save(&battle)

	// --- AI TURN TRIGGER ---
	if !gameEnded && battle.WinnerID == nil && strings.Contains(battle.BattleType, "PVE") && battle.CurrentTurnPlayerID == battle.Player2ID {
		if err := s.executeAITurn(&battle); err != nil {
			fmt.Printf("AI Execution Failed: %v\n", err)
		}
	}

	return &battle, nil
}

// executeAITurn handles AI logic for PvE
func (s *BattleService) executeAITurn(battle *models.Battle) error {
	// 1. Identify AI Character (Player 2's active char)
	// Simplified: Fetch Player 2's FIRST active character
	// In real logic, we'd check PlayerStateP2 or a dedicated ActiveCharacter table
	// For MVP, we assume Player 2 has one character active for now or fetch from DB.
	var aiChar models.Character
	// Assuming Player 2 has characters.
	// Find FIRST non-fainted character owned by Player 2
	err := db.DB.Where("owner_id = ? AND is_fainted = false", battle.Player2ID).First(&aiChar).Error
	if err != nil {
		// AI has no chars? AI signs of surrender/loss?
		// CheckBattleEnd should handle it.
		return nil
	}

	// 2. Identify Target (Player 1's active char)
	var playerChar models.Character
	// Pick one random? or First?
	err = db.DB.Where("owner_id = ? AND is_fainted = false", battle.Player1ID).First(&playerChar).Error
	if err != nil {
		return nil // Player 1 dead?
	}

	// 3. Decide Action
	// Get Usable Skills
	usableSkills, err := s.skillService.GetUsableSkills(aiChar.ID, aiChar.CurrentMana)

	actionType := "attack"
	var skillID uint = 0

	if err == nil && len(usableSkills) > 0 {
		// Simple AI: Pick random usable skill
		// Or pick highest damage?
		// Random for unpredictability
		idx := rand.Intn(len(usableSkills))
		skill := usableSkills[idx]

		// 30% chance to just basic attack to save mana?
		if rand.Float32() > 0.3 {
			actionType = "skill"
			skillID = skill.ID
		}
	}

	// 4. Construct Action Data
	actionData := map[string]interface{}{
		"action":       actionType,
		"character_id": float64(aiChar.ID),
		"target_id":    float64(playerChar.ID),
	}
	if actionType == "skill" {
		actionData["skill_id"] = float64(skillID)
	}

	// 5. Calculate AI Move (Reuse ProcessTurn logic parts?)
	// Calling ProcessTurn recursively for the AI user (battle.Player2ID)
	// Ensure we don't infinitely recurse. ProcessTurn checks IsTurn logic.
	// Since we updated CurrentTurnPlayerID to P2 before calling this, it should pass verification.
	// And at end of AI turn, it sets P1. P1 is NOT "PVE" trigger (only P2 is AI).
	// So recursion depth = 1. Safe.

	_, err = s.ProcessTurn(battle.ID, battle.Player2ID, actionData)
	return err
}

// Helper to adapt DB Character to BattleParticipant
func (s *BattleService) toParticipant(c *models.Character) *models.BattleParticipant {
	return &models.BattleParticipant{
		CharacterID:   c.ID,
		CharacterName: c.Name,
		Element:       c.Element, // NEW
		MaxHP:         c.BaseHP,  // Simplified
		CurrentHP:     c.CurrentHP,
		MaxMana:       100, // Fixed for now
		CurrentMana:   c.CurrentMana,
		Attack:        c.CurrentAttack,
		Defense:       c.CurrentDefense,
		Speed:         c.CurrentSpeed,
		IsFainted:     c.IsFainted,
		IsActive:      true,
	}
}

// CompleteBattle handles victory with Anti-Cheat validation
func (s *BattleService) CompleteBattle(battleID uint, winnerID uint, replayData string) error {
	// 1. Anti-Cheat: Validate Replay Data (Basic sanity checks for now)
	if replayData != "" {
		if err := s.ValidateReplay(battleID, winnerID, replayData); err != nil {
			// Log security event
			fmt.Printf("SECURTY ALERT: Battle %d failed validation: %v\n", battleID, err)
			return fmt.Errorf("security check failed: %v", err)
		}
	} else {
		// Log missing replay (Soft warning for legacy clients, Hard error for new GDevelop clients)
		// For consistency, we require it for ranked/wager
		var checkBattle models.Battle
		if err := db.DB.First(&checkBattle, battleID).Error; err == nil {
			if checkBattle.BattleType == "wager" || checkBattle.BattleType == "ranked" {
				// Strict mode for sensitive battles
				// return errors.New("missing replay data") // Uncomment when client is ready
				fmt.Printf("WARNING: Missing replay data for sensitive battle %d\n", battleID)
			}
		}
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var battle models.Battle
		if err := tx.First(&battle, battleID).Error; err != nil {
			return err
		}

		if battle.Status != "active" {
			// Already completed (idempotency check)
			return nil
		}

		battle.Status = "completed"
		battle.WinnerID = &winnerID
		now := time.Now()
		battle.EndedAt = &now

		// Save Replay Data for audit
		if replayData != "" {
			// In a real app, save to S3 or separate table. Here we log or simplify.
			// battle.ReplayLog = replayData // Assuming field exists or we add it
		}

		if err := tx.Save(&battle).Error; err != nil {
			return err
		}

		// Handle Rewards
		if battle.BattleType == "wager" {
			// WAGER LOGIC: Release Escrow (Winner takes pot minus fee)
			// Pot = P1Bet + P2Bet
			pot := battle.Player1Bet + battle.Player2Bet

			// If old data (using BetAmount), fallback
			if pot == 0 && battle.BetAmount > 0 {
				pot = battle.BetAmount * 2
				battle.Player1Bet = battle.BetAmount
				battle.Player2Bet = battle.BetAmount
			}

			// Dynamic Payout based on Difficulty/Risk
			// Logic: Winner gets their own bet back + Opponent's bet (minus fee)
			// Fee is 5% logic on winnings

			feeRate := 0.05
			var winnerBet, loserBet int64

			if winnerID == battle.Player1ID {
				winnerBet = battle.Player1Bet
				loserBet = battle.Player2Bet
			} else {
				winnerBet = battle.Player2Bet
				loserBet = battle.Player1Bet
			}

			fee := int64(float64(loserBet) * feeRate) // Fee taken from WINNINGS (loser's money)
			winnings := loserBet - fee
			winnerPayout := winnerBet + winnings
			treasuryAmount := fee

			// Transfer Escrow -> Winner & Treasury
			escrowAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeEscrow, "GTK")
			winnerAcc, _ := s.ledger.GetOrCreateAccount(&winnerID, models.AccountTypeWallet, "GTK")
			treasuryAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeTreasury, "GTK")

			entries := []models.LedgerEntry{
				{AccountID: escrowAcc.ID, Amount: -pot, Type: "DEBIT"}, // Drain total pot
				{AccountID: winnerAcc.ID, Amount: winnerPayout, Type: "CREDIT"},
				{AccountID: treasuryAcc.ID, Amount: treasuryAmount, Type: "CREDIT"},
			}

			desc := "Wager Win Payout (Risk Reward)"
			if err := s.ledger.CreateTransactionWithTx(tx, models.TxTypeWagerWin, fmt.Sprintf("wager_win_%d", battleID), desc, entries); err != nil {
				return err
			}
		} else if battle.BattleType == "ranked" {
			// Rank Reward: 25 GTK
			userAcc, _ := s.ledger.GetOrCreateAccount(&winnerID, models.AccountTypeWallet, "GTK")
			rewardAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeReward, "GTK")

			entries := []models.LedgerEntry{
				{AccountID: rewardAcc.ID, Amount: -25, Type: "DEBIT"},
				{AccountID: userAcc.ID, Amount: 25, Type: "CREDIT"},
			}
			s.ledger.CreateTransactionWithTx(tx, models.TxTypeRankedReward, fmt.Sprintf("battle_%d", battleID), "Ranked Win", entries)
		} else if strings.Contains(battle.BattleType, "PVE") {
			// PvE Reward: Small Token + XP?
			// For MVP: 10 GTK
			userAcc, _ := s.ledger.GetOrCreateAccount(&winnerID, models.AccountTypeWallet, "GTK")
			rewardAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeReward, "GTK")

			entries := []models.LedgerEntry{
				{AccountID: rewardAcc.ID, Amount: -10, Type: "DEBIT"},
				{AccountID: userAcc.ID, Amount: 10, Type: "CREDIT"},
			}
			s.ledger.CreateTransactionWithTx(tx, models.TxTypeReward, fmt.Sprintf("pve_win_%d", battleID), "PvE Victory Reward", entries)
		}

		// --- POST-BATTLE HOOKS: Stats, Elo, XP ---
		// 1. Fetch Users
		var winner, loser models.User
		if err := tx.First(&winner, winnerID).Error; err != nil {
			return err
		}
		// Loser might be AI/System (0) or actual player?
		// For PvP/Wager/Ranked, loser is real.
		// For PvE, loser is System (ID 0? No, usually not stored as User).
		// check if PvP
		isPvP := battle.BattleType == "wager" || battle.BattleType == "ranked" || battle.BattleType == "pvp"

		loserID := battle.Player1ID
		if winnerID == battle.Player1ID {
			loserID = battle.Player2ID
		}

		if isPvP {
			if err := tx.First(&loser, loserID).Error; err != nil {
				return err
			}

			// 2. Update Elo (Ranked/Wager only)
			if battle.BattleType == "ranked" || battle.BattleType == "wager" {
				kFactor := 32
				expectedWin := 1.0 / (1.0 + float64(loser.ELO-winner.ELO)/400.0)
				// expectedLoss := 1.0 - expectedWin // unused

				eloChange := int(float64(kFactor) * (1.0 - expectedWin))
				winner.ELO += eloChange
				loser.ELO -= eloChange
				if loser.ELO < 0 {
					loser.ELO = 0
				} // No negative Elo
			}

			// 3. Update Stats
			winner.PvPWins++
			winner.CurrentWinStreak++
			loser.PvPLosses++
			loser.CurrentWinStreak = 0
		}

		// 4. Grant XP (Winner)
		// Simple flat XP for now. Complex formula can go in ProgressionService
		winner.Experience += 100
		// Check Level Up? (Simplified)
		needed := winner.Level * 1000
		if winner.Experience >= needed {
			winner.Level++
			winner.Experience -= needed
		}

		if err := tx.Save(&winner).Error; err != nil {
			return err
		}
		if isPvP {
			if err := tx.Save(&loser).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ValidateReplay performs basic anti-cheat checks
func (s *BattleService) ValidateReplay(battleID, winnerID uint, replayData string) error {
	// 1. Check Data Size
	if len(replayData) < 10 {
		return errors.New("invalid replay data size")
	}

	// 2. STUB: Parse JSON (Assume simple list of actions)
	// In future: Unmarshal to []Action and simulate
	// For now, checks are rudimentary

	// 3. Verify Winner matches logic? (Needs simulation)

	// 4. Time Check (Optional) - verify battle duration matches log
	// This requires created_at vs now check in handler before calling validation?

	return nil
}

// CheckTimeouts checks for abandoned battles (> 30s inactive)
func (s *BattleService) CheckTimeouts() error {
	threshold := time.Now().Add(-30 * time.Second)

	var staleBattles []models.Battle
	// Find active battles updated before threshold
	// Note: UpdatedAt is updated on every Save() which happens on every Turn.
	if err := db.DB.Where("status = ? AND updated_at < ?", "active", threshold).Find(&staleBattles).Error; err != nil {
		return err
	}

	for _, battle := range staleBattles {
		fmt.Printf("Battle %d timed out. Last update: %v\n", battle.ID, battle.UpdatedAt)

		// Auto-Surrender Current Turn Player
		var winnerID uint
		if battle.CurrentTurnPlayerID == battle.Player1ID {
			winnerID = battle.Player2ID
		} else {
			winnerID = battle.Player1ID
		}

		if strings.Contains(battle.BattleType, "PVE") {
			// For PVE, if timeout, user loses. Winner = 0.
			winnerID = 0
		}

		if err := s.CompleteBattle(battle.ID, winnerID, ""); err != nil {
			fmt.Printf("Failed to complete timed out battle %d: %v\n", battle.ID, err)
		}
	}

	return nil
}

// CalculateDynamicStakes determines how much each player MUST risk
func (s *BattleService) CalculateDynamicStakes(p1ID, p2ID uint, p1Team []models.BattleParticipant) (int64, int64, error) {
	// 1. Get P2 Team
	p2Team, err := s.GetTeamSnapshot(p2ID)
	if err != nil {
		return 0, 0, err
	}

	// 2. Calculate Combat Power (CP)
	cp1 := calculateTeamCP(p1Team)
	cp2 := calculateTeamCP(p2Team)

	// 3. Determine Risk/Odds
	// Base Stake = 100 GTK (Standard Unit)
	baseStake := int64(100)

	// Avoid div by zero
	if cp1 < 1 {
		cp1 = 1
	}
	if cp2 < 1 {
		cp2 = 1
	}

	ratio := float64(cp1) / float64(cp2)

	var p1Stake, p2Stake int64

	if ratio > 1.0 {
		// P1 is Stronger
		if ratio > 5.0 {
			ratio = 5.0
		}
		p1Stake = int64(float64(baseStake) * ratio)
		p2Stake = int64(float64(baseStake) / ratio)
	} else {
		// P2 is Stronger (or equal)
		ratio = float64(cp2) / float64(cp1)
		if ratio > 5.0 {
			ratio = 5.0
		}
		p2Stake = int64(float64(baseStake) * ratio)
		p1Stake = int64(float64(baseStake) / ratio)
	}

	// Ensure Minimums
	if p1Stake < 10 {
		p1Stake = 10
	}
	if p2Stake < 10 {
		p2Stake = 10
	}

	return p1Stake, p2Stake, nil
}

func calculateTeamCP(team []models.BattleParticipant) int {
	total := 0
	for _, m := range team {
		cp := m.MaxHP + (m.Attack * 2) + (m.Defense * 2) + m.Speed
		total += cp
	}
	return total
}
