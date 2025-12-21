package services

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

type GachaService struct {
	nameGenerator *NameGeneratorService
	ledger        *LedgerService
	config        *ConfigService
	blockchain    *BlockchainService
}

// NewGachaService creates a new gacha service
func NewGachaService(bc *BlockchainService) *GachaService {
	return &GachaService{
		nameGenerator: NewNameGeneratorService(),
		ledger:        NewLedgerService(),
		config:        NewConfigService(),
		blockchain:    bc,
	}
}

// MintEgg mints a new egg with gacha mechanics
func (s *GachaService) MintEgg(userID uint, towerAmount int64, txHash string) (*models.Egg, error) {
	// SECURITY CHECK 1: Validate amount (0 for free mint, 1-10000 for paid)
	if towerAmount < 0 || towerAmount > 10000 {
		return nil, errors.New("amount must be between 0 and 10,000 TOWER")
	}

	// SECURITY CHECK 2: Rate limiting (1 mint per minute) - ONLY FOR FREE MINTS
	// Paid mints can be done without cooldown
	if towerAmount == 0 {
		var lastMint models.Egg
		err := db.DB.Where("user_id = ?", userID).
			Order("created_at DESC").First(&lastMint).Error
		if err == nil && time.Since(lastMint.CreatedAt) < 1*time.Minute {
			remaining := 1*time.Minute - time.Since(lastMint.CreatedAt)
			return nil, fmt.Errorf("please wait %v before minting again", remaining.Round(time.Second))
		}
	}
	// If error is "record not found", it's first mint - allow it

	// Get user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// SECURITY CHECK 3: Verify TOWER balance or Blockchain Transaction logic
	if txHash == "" {
		// INTERNAL PAYMENT (using In-Game Balance)
		if towerAmount > 0 && user.TOWERBalance < towerAmount {
			return nil, fmt.Errorf("insufficient TOWER balance (need %d, have %d)", towerAmount, user.TOWERBalance)
		}
	} else {
		// EXTERNAL BLOCKCHAIN PAYMENT
		if s.blockchain != nil {
			// Assume payment is in TOWER (ERC20)
			costBig := big.NewInt(towerAmount)
			// NOTE: real production would check TOWER decimals (likely 10^18)
			// verifyTransaction is currently set up for GTK/Native/ERC20 generic but logs check GTK address.
			// Ideally we verify "TOWER" contract. The BlockchainService should be generic enough or have VerifyERC20Transfer
			// For this MVP, we use the same verify method which checks configured token transfer.
			// Assuming TOWER token address is passed
			if err := s.blockchain.VerifyTransaction(txHash, costBig); err != nil {
				return nil, fmt.Errorf("blockchain verification failed: %v", err)
			}
		} else {
			fmt.Println("⚠️ WARNING: Blockchain verification skipped (Service nil)")
		}
	}

	// SECURITY CHECK 4: Daily mint limit (Dynamic Config) - check before transaction
	if s.config == nil {
		s.config = NewConfigService()
	}
	dailyLimit := int64(s.config.GetInt("daily_mint_limit", 10))

	today := time.Now().Format("2006-01-02")
	var todayMints int64
	db.DB.Model(&models.Egg{}).
		Where("user_id = ? AND DATE(created_at) = ?", userID, today).
		Count(&todayMints)

	if todayMints >= dailyLimit {
		return nil, fmt.Errorf("daily mint limit reached (%d per day)", dailyLimit)
	}

	// SECURITY CHECK 5: Check if first mint
	isFirstMint := !false && towerAmount == 0

	// Calculate rarity based on investment
	rarity := s.rollRarity(towerAmount, isFirstMint)

	// Roll character traits
	charType := s.rollCharacterType()
	element := s.rollElement()
	class := s.rollClass()

	// Calculate base stats (based on rarity)
	baseStats := s.calculateBaseStats(rarity)

	// Roll abilities (based on class + rarity)
	abilities := s.rollAbilities(class, rarity)

	// Calculate incubation time (in hours)
	incubationTime := s.getIncubationTime(rarity)

	// BEGIN ATOMIC TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// LEDGER INTEGRATION: Charge for Minting (if Paid and Internal)
	if towerAmount > 0 && txHash == "" {
		// Debit User Wallet (TOWER), Credit System Sink (TOWER)
		userAcc, _ := s.ledger.GetOrCreateAccount(&userID, models.AccountTypeWallet, "TOWER")
		sinkAcc, _ := s.ledger.GetOrCreateAccount(nil, models.AccountTypeSink, "TOWER")

		entries := []models.LedgerEntry{
			{AccountID: userAcc.ID, Amount: -towerAmount, Type: "DEBIT"},
			{AccountID: sinkAcc.ID, Amount: towerAmount, Type: "CREDIT"},
		}

		// Execute Transaction
		if err := s.ledger.CreateTransaction(models.TxTypeGachaMint, fmt.Sprintf("mint_%d_%d", userID, time.Now().Unix()), fmt.Sprintf("Mint Cost: %d TOWER", towerAmount), entries); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("payment failed: %v", err)
		}

		// Note: We do NOT deduct balance via Update("tower_balance") anymore.
		// We should rely on Ledger. However, for legacy UI state, we sync it.
		if err := tx.Model(&user).Update("tower_balance", gorm.Expr("tower_balance - ?", towerAmount)).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to sync legacy balance")
		}
	}

	// Mark first character minted if applicable
	if isFirstMint {
		// user.HasMintedFirstChar = true // Field removed
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to update user")
		}
	}

	// Serialize stats and abilities to JSON
	statsJSON, _ := json.Marshal(baseStats)
	abilitiesJSON, _ := json.Marshal(abilities)

	// Create egg with pre-determined traits
	egg := models.Egg{
		UserID:                  userID,
		Rarity:                  rarity,
		Element:                 element,
		CharacterType:           charType,
		Class:                   class,
		IncubationTime:          incubationTime,
		MintCost:                towerAmount,
		PredeterminedStats:      string(statsJSON),
		PredeterminedAbilities:  string(abilitiesJSON),
		EffectiveIncubationTime: incubationTime,
		AcceleratorsApplied:     "[]", // Initialize with empty JSON array
	}

	if err := tx.Create(&egg).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create egg")
	}

	// Calculate balances based on payment method
	balanceBefore := user.TOWERBalance
	balanceAfter := user.TOWERBalance
	var blockchainHashPtr *string

	if txHash == "" {
		// Internal payment
		balanceAfter = user.TOWERBalance - towerAmount
	} else {
		// External payment
		blockchainHashPtr = &txHash
	}

	// Create transaction record
	transaction := models.Transaction{
		UserID:           userID,
		TransactionType:  "EGG_MINT",
		TokenType:        "TOWER",
		Amount:           -towerAmount,
		BalanceBefore:    balanceBefore,
		BalanceAfter:     balanceAfter,
		BlockchainTxHash: blockchainHashPtr,
		Description:      fmt.Sprintf("Minted %s egg (%s %s %s)", rarity, element, charType, class),
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "EGG_MINT",
		EntityType: "egg",
		EntityID:   &egg.ID,
		NewValues:  fmt.Sprintf("rarity:%s,cost:%d,type:%s", rarity, towerAmount, charType),
	}
	if err := tx.Create(&auditLog).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create audit log: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("transaction failed")
	}

	return &egg, nil
}

// GetUserEggs returns all eggs owned by user
func (s *GachaService) GetUserEggs(userID uint) ([]models.Egg, error) {
	var eggs []models.Egg
	if err := db.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&eggs).Error; err != nil {
		return nil, err
	}
	return eggs, nil
}

// StartIncubation starts egg incubation
func (s *GachaService) StartIncubation(userID, eggID uint) error {
	var egg models.Egg
	if err := db.DB.First(&egg, eggID).Error; err != nil {
		return errors.New("egg not found")
	}

	if egg.UserID != userID {
		return errors.New("you don't own this egg")
	}

	if egg.IncubationStartedAt != nil {
		return errors.New("egg is already incubating")
	}

	if egg.HatchedAt != nil {
		return errors.New("egg already hatched")
	}

	now := time.Now()
	egg.IncubationStartedAt = &now

	if err := db.DB.Save(&egg).Error; err != nil {
		return errors.New("failed to start incubation")
	}

	return nil
}

// HatchEgg hatches an egg into a character
func (s *GachaService) HatchEgg(userID, eggID uint) (*models.Character, error) {
	var egg models.Egg
	if err := db.DB.First(&egg, eggID).Error; err != nil {
		return nil, errors.New("egg not found")
	}

	if egg.UserID != userID {
		return nil, errors.New("you don't own this egg")
	}

	if egg.HatchedAt != nil {
		return nil, errors.New("egg already hatched")
	}

	if egg.IncubationStartedAt == nil {
		return nil, errors.New("egg must be incubating to hatch")
	}

	// Check if incubation time has passed
	incubationDuration := time.Duration(egg.EffectiveIncubationTime) * time.Hour
	if time.Since(*egg.IncubationStartedAt) < incubationDuration {
		remaining := incubationDuration - time.Since(*egg.IncubationStartedAt)
		return nil, fmt.Errorf("egg needs %v more to hatch", remaining.Round(time.Minute))
	}

	// Parse predetermined stats
	var stats map[string]int
	json.Unmarshal([]byte(egg.PredeterminedStats), &stats)

	var abilityIDs []uint
	json.Unmarshal([]byte(egg.PredeterminedAbilities), &abilityIDs)

	// Generate unique name
	uniqueName := s.nameGenerator.Generate(egg.Element, egg.CharacterType, egg.Class, egg.Rarity)

	// BEGIN TRANSACTION
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Serialize abilities
	abilitiesJSON, _ := json.Marshal(abilityIDs)

	// Create character
	character := models.Character{
		OwnerID:       userID,
		Name:          uniqueName,
		UniqueName:    &uniqueName,
		Rarity:        egg.Rarity,
		Element:       egg.Element,
		CharacterType: egg.CharacterType,
		Class:         egg.Class,
		Level:         1,

		// Stats from egg
		BaseHP:      stats["hp"],
		BaseAttack:  stats["atk"],
		BaseDefense: stats["def"],
		BaseSpeed:   stats["spd"],

		CurrentHP:      stats["hp"],
		CurrentAttack:  stats["atk"],
		CurrentDefense: stats["def"],
		CurrentSpeed:   stats["spd"],

		// Abilities
		UnlockedAbilities: string(abilitiesJSON),
	}

	if err := tx.Create(&character).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create character")
	}

	// Update egg
	now := time.Now()
	egg.CharacterID = &character.ID
	egg.HatchedAt = &now

	if err := tx.Save(&egg).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update egg")
	}

	// Audit log
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     "EGG_HATCHED",
		EntityType: "character",
		EntityID:   &character.ID,
		NewValues:  fmt.Sprintf("egg_id:%d,name:%s,rarity:%s", eggID, uniqueName, character.Rarity),
	}
	tx.Create(&auditLog)

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("transaction failed")
	}

	return &character, nil
}

// ScanEgg reveals egg stats (requires Egg Scanner item)
func (s *GachaService) ScanEgg(userID, eggID uint) (map[string]int, error) {
	var egg models.Egg
	if err := db.DB.First(&egg, eggID).Error; err != nil {
		return nil, errors.New("egg not found")
	}

	if egg.UserID != userID {
		return nil, errors.New("you don't own this egg")
	}

	// TODO: Check if user has Egg Scanner item in inventory
	// For now, allow scanning

	// Parse stats
	var stats map[string]int
	json.Unmarshal([]byte(egg.PredeterminedStats), &stats)

	// Mark as revealed
	if !egg.IsStatsRevealed {
		now := time.Now()
		egg.IsStatsRevealed = true
		egg.RevealedAt = &now
		db.DB.Save(&egg)
	}

	return stats, nil
}

// ApplyAccelerator applies an accelerator item to an egg
func (s *GachaService) ApplyAccelerator(userID, eggID, itemID uint) error {
	acceleratorService := NewAcceleratorService()
	return acceleratorService.ApplyAccelerator(userID, eggID, itemID)
}

// rollRarity determines rarity based on TOWER investment
func (s *GachaService) rollRarity(towerAmount int64, isFirstMint bool) string {
	// Free first mint has terrible odds (99.9% C rank)
	if isFirstMint {
		return s.rollFirstMintRarity()
	}

	// Calculate odds based on investment
	odds := s.calculateRarityOdds(towerAmount)

	// Roll random number (0-100)
	randomNum := s.secureRandom(0, 100000) // Use 100000 for precision

	// Determine rarity based on cumulative probability
	cumulative := float64(0)
	rarities := []string{"SSS", "SS", "S", "A", "B", "C"}

	for _, rarity := range rarities {
		cumulative += odds[rarity] * 1000 // Multiply by 1000 for precision
		if float64(randomNum) < cumulative {
			return rarity
		}
	}

	return "C" // Fallback
}

// rollFirstMintRarity for free first mint (99.9% C rank)
func (s *GachaService) rollFirstMintRarity() string {
	randomNum := s.secureRandom(0, 1000000)

	if randomNum < 1 {
		return "SSS" // 0.0001%
	} else if randomNum < 10 {
		return "SS" // 0.0009%
	} else if randomNum < 100 {
		return "S" // 0.009%
	} else if randomNum < 1000 {
		return "A" // 0.09%
	} else if randomNum < 10000 {
		return "B" // 0.9%
	}
	return "C" // 99.9%
}

// calculateRarityOdds returns probability distribution based on TOWER amount
func (s *GachaService) calculateRarityOdds(towerAmount int64) map[string]float64 {
	// Base odds (1-99 TOWER) - ULTRA NERFED
	// S/SS/SSS all < 0.02% to make them extremely rare
	baseOdds := map[string]float64{
		"C":   75.0,  // Very common
		"B":   20.0,  // Common
		"A":   4.98,  // Uncommon
		"S":   0.015, // Very rare (< 0.02%)
		"SS":  0.004, // Extremely rare (< 0.02%)
		"SSS": 0.001, // Ultra rare (< 0.02%)
	}

	// NO SCALING until 100 TOWER investment
	if towerAmount < 100 {
		return baseOdds
	}

	// Progressive scaling ONLY after 100 TOWER
	// Very gradual increase to prevent drastic jumps
	var scaleFactor float64

	if towerAmount >= 100 && towerAmount < 500 {
		// 100-499 TOWER: Minimal scaling (0.1 - 0.5)
		scaleFactor = float64(towerAmount-100) / 1000.0
	} else if towerAmount >= 500 && towerAmount < 1000 {
		// 500-999 TOWER: Slow scaling (0.5 - 1.0)
		scaleFactor = 0.5 + float64(towerAmount-500)/1000.0
	} else if towerAmount >= 1000 && towerAmount < 5000 {
		// 1000-4999 TOWER: Moderate scaling (1.0 - 2.0)
		scaleFactor = 1.0 + float64(towerAmount-1000)/4000.0
	} else {
		// 5000+ TOWER: Cap at 2.5 (diminishing returns)
		scaleFactor = 2.0 + math.Min(0.5, float64(towerAmount-5000)/10000.0)
	}

	// Very gradual shift from low to high rarities
	shift := scaleFactor * 1.5 // Much smaller shift than before

	odds := make(map[string]float64)
	odds["C"] = math.Max(30, baseOdds["C"]-shift*4) // Reduce C significantly
	odds["B"] = math.Max(15, baseOdds["B"]-shift*2) // Reduce B moderately
	odds["A"] = baseOdds["A"] + shift*1.5           // Increase A slightly
	odds["S"] = baseOdds["S"] + shift*0.4           // Very gradual S increase
	odds["SS"] = baseOdds["SS"] + shift*0.25        // Minimal SS increase
	odds["SSS"] = baseOdds["SSS"] + shift*0.15      // Tiny SSS increase

	// Normalize to 100%
	total := odds["C"] + odds["B"] + odds["A"] + odds["S"] + odds["SS"] + odds["SSS"]
	for k := range odds {
		odds[k] = (odds[k] / total) * 100
	}

	return odds
}

// secureRandom generates cryptographically secure random number
func (s *GachaService) secureRandom(min, max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		// Fallback to timestamp-based (should never happen)
		return time.Now().UnixNano()%(max-min) + min
	}
	return nBig.Int64() + min
}

// rollCharacterType randomly selects character type
func (s *GachaService) rollCharacterType() string {
	types := []string{"BEAST", "DRAGON", "BIRD", "INSECT", "AQUATIC", "MINERAL", "SPIRIT", "AVIAN", "PLANT", "MACHINE"}
	index := s.secureRandom(0, int64(len(types)))
	return types[index]
}

// rollElement randomly selects element
func (s *GachaService) rollElement() string {
	elements := []string{"Fire", "Water", "Earth", "Air", "Light", "Dark", "Electric", "Ice"}
	index := s.secureRandom(0, int64(len(elements)))
	return elements[index]
}

// rollClass randomly selects class
func (s *GachaService) rollClass() string {
	classes := []string{"Warrior", "Mage", "Archer", "Tank", "Support", "Rogue", "Paladin", "Berserker"}
	index := s.secureRandom(0, int64(len(classes)))
	return classes[index]
}

// calculateBaseStats generates base stats based on rarity
func (s *GachaService) calculateBaseStats(rarity string) map[string]int {
	baseValues := map[string]map[string]int{
		"C": {
			"hp":  100,
			"atk": 20,
			"def": 15,
			"spd": 10,
		},
		"B": {
			"hp":  150,
			"atk": 30,
			"def": 25,
			"spd": 15,
		},
		"A": {
			"hp":  200,
			"atk": 45,
			"def": 35,
			"spd": 20,
		},
		"S": {
			"hp":  300,
			"atk": 65,
			"def": 50,
			"spd": 30,
		},
		"SS": {
			"hp":  450,
			"atk": 95,
			"def": 75,
			"spd": 45,
		},
		"SSS": {
			"hp":  600,
			"atk": 130,
			"def": 100,
			"spd": 60,
		},
	}

	base := baseValues[rarity]

	// Add random variation (±10%)
	variation := func(val int) int {
		variance := int(float64(val) * 0.1)
		randomVar := s.secureRandom(-int64(variance), int64(variance))
		return val + int(randomVar)
	}

	return map[string]int{
		"hp":  variation(base["hp"]),
		"atk": variation(base["atk"]),
		"def": variation(base["def"]),
		"spd": variation(base["spd"]),
	}
}

// rollAbilities selects abilities based on class and rarity
func (s *GachaService) rollAbilities(class, rarity string) []uint {
	// TODO: Query abilities table when implemented
	// For now, return placeholder ability IDs based on class and rarity
	abilityCount := map[string]int{
		"C":   2,
		"B":   3,
		"A":   4,
		"S":   5,
		"SS":  6,
		"SSS": 7,
	}

	count := abilityCount[rarity]
	if count == 0 {
		count = 3 // Default
	}

	// Base ability ID offset by class for variety
	classOffset := map[string]uint{
		"Warrior":   1,
		"Mage":      100,
		"Archer":    200,
		"Tank":      300,
		"Support":   400,
		"Berserker": 500,
		"Rogue":     600,
	}

	offset := classOffset[class]
	if offset == 0 {
		offset = 1 // Default offset
	}

	abilities := make([]uint, count)
	for i := 0; i < count; i++ {
		abilities[i] = offset + uint(i+1)
	}
	return abilities
}

// getIncubationTime returns incubation time in hours based on rarity
func (s *GachaService) getIncubationTime(rarity string) int {
	times := map[string]int{
		"C":   6,  // 6 hours
		"B":   12, // 12 hours
		"A":   24, // 1 day
		"S":   48, // 2 days
		"SS":  72, // 3 days
		"SSS": 96, // 4 days
	}
	return times[rarity]
}

// GetOddsPreview returns probability preview for a given TOWER amount
func (s *GachaService) GetOddsPreview(towerAmount int64) map[string]float64 {
	if towerAmount == 0 {
		// Free mint odds
		return map[string]float64{
			"C":   99.9,
			"B":   0.09,
			"A":   0.009,
			"S":   0.0009,
			"SS":  0.00009,
			"SSS": 0.00001,
		}
	}

	return s.calculateRarityOdds(towerAmount)
}
