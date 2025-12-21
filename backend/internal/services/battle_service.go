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
		s.CompleteBattle(battle.ID, winnerID)
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
		MaxHP:         c.BaseHP, // Simplified
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

// CompleteBattle handles victory
func (s *BattleService) CompleteBattle(battleID uint, winnerID uint) error {
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

		if err := tx.Save(&battle).Error; err != nil {
			return err
		}

		// Handle Rewards
		if battle.BattleType == "wager" {
			// WAGER LOGIC: Release Escrow (Winner takes pot minus fee)
			// Pot = BetAmount * 2
			pot := battle.BetAmount * 2
			fee := int64(float64(pot) * 0.05) // 5% fee
			payout := pot - fee

			// Transfer Escrow -> Winner & Treasury
			escrowAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeEscrow, "GTK")
			winnerAcc, _ := s.ledger.GetOrCreateAccount(&winnerID, models.AccountTypeWallet, "GTK")
			treasuryAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeTreasury, "GTK")

			entries := []models.LedgerEntry{
				{AccountID: escrowAcc.ID, Amount: -pot, Type: "DEBIT"}, // Drain Escrow
				{AccountID: winnerAcc.ID, Amount: payout, Type: "CREDIT"},
				{AccountID: treasuryAcc.ID, Amount: fee, Type: "CREDIT"},
			}

			if err := s.ledger.CreateTransactionWithTx(tx, models.TxTypeWagerWin, fmt.Sprintf("wager_win_%d", battleID), "Wager Payout", entries); err != nil {
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

		return nil
	})
}
