package services

import (
	cryptoRand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// SkillActivationService handles skill usage in battles
type SkillActivationService struct{}

// NewSkillActivationService creates a new skill activation service
func NewSkillActivationService() *SkillActivationService {
	return &SkillActivationService{}
}

// SkillActivationRequest represents a request to use a skill
type SkillActivationRequest struct {
	CharacterID uint
	AbilityID   uint
	TargetID    uint   // Single target (0 for self/AOE)
	TargetIDs   []uint // Multiple targets for AOE
	BattleID    uint
	TurnNumber  int
}

// SkillActivationResult contains the result of skill activation
type SkillActivationResult struct {
	Success        bool
	Damage         int
	Healing        int
	ManaUsed       int
	EffectsApplied []string
	CriticalHit    bool
	Message        string
	AnimationName  string
}

// ActivateSkill validates and executes a skill
func (s *SkillActivationService) ActivateSkill(req SkillActivationRequest) (*SkillActivationResult, error) {
	// 1. Load character
	var character models.Character
	if err := db.DB.First(&character, req.CharacterID).Error; err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	// 2. Load ability
	var ability models.Ability
	if err := db.DB.First(&ability, req.AbilityID).Error; err != nil {
		return nil, fmt.Errorf("ability not found: %w", err)
	}

	// 3. Validate skill activation
	if err := s.ValidateSkillActivation(&character, &ability); err != nil {
		return nil, err
	}

	// 4. Check if skill is active
	if !s.IsSkillActive(character.ID, ability.ID) {
		return nil, errors.New("skill not equipped in active slots")
	}

	// 5. Check cooldown
	if s.IsOnCooldown(character.ID, ability.ID) {
		return nil, errors.New("skill is on cooldown")
	}

	// 6. Deduct mana
	if err := s.DeductMana(&character, ability.ManaCost); err != nil {
		return nil, err
	}

	// 7. Calculate logic
	result := s.CalculateSkillOutcome(&character, &ability)

	// 7b. Apply side effects (Buffs to DB)
	for _, effect := range result.EffectsApplied {
		if effect == ability.AppliesBuff {
			s.ApplyBuff(character.ID, effect, ability.BuffDuration)
		}
	}

	// 8. Start cooldown
	s.StartCooldown(character.ID, ability.ID, ability.Cooldown)

	// 9. Update character state
	db.DB.Save(&character)

	// 10. Log skill usage
	s.LogSkillUsage(req, result)

	return result, nil
}

// CalculateSkillOutcome computes the result of a skill (damage, healing, effects) without side effects
// This can be used by PvE, PvP, and Raids
func (s *SkillActivationService) CalculateSkillOutcome(char *models.Character, ability *models.Ability) *SkillActivationResult {
	result := &SkillActivationResult{
		Success:        true,
		ManaUsed:       ability.ManaCost,
		AnimationName:  ability.AnimationName,
		EffectsApplied: []string{},
	}

	// Calculate damage with element bonuses
	if ability.BaseDamage > 0 {
		damage := s.CalculateDamage(char, ability)
		result.Damage = damage
		result.Message = fmt.Sprintf("%s dealt %d damage!", ability.Name, damage)
	}

	// Calculate healing
	if ability.BaseHeal > 0 {
		healing := s.CalculateHealing(char, ability)
		result.Healing = healing
		result.Message = fmt.Sprintf("%s healed %d HP!", ability.Name, healing)
	}

	// Apply buffs logic (just recording them for result)
	if ability.AppliesBuff != "" {
		// In a real generic function, we might just return "Intent to buff"
		// But for now we simulate the legacy ApplySkillEffects behavior minus the DB writes?
		// Wait, ApplySkillEffects did DB writes for buffs.
		// For pure calculation, we should just return what *should* happen.
		// The caller (ActivateSkill or RaidService) handles the DB/State updates.
		result.EffectsApplied = append(result.EffectsApplied, ability.AppliesBuff)
	}

	// Apply debuffs
	if ability.AppliesDebuff != "" {
		result.EffectsApplied = append(result.EffectsApplied, ability.AppliesDebuff)
	}

	// Check for critical hit
	if ability.BaseDamage > 0 {
		result.CriticalHit = s.RollCritical(char)
		if result.CriticalHit {
			result.Damage = int(float64(result.Damage) * 1.5)
			result.Message += " CRITICAL HIT!"
		}
	}

	return result
}

// ValidateSkillActivation checks if a skill can be used
func (s *SkillActivationService) ValidateSkillActivation(char *models.Character, ability *models.Ability) error {
	// Check mana
	if char.CurrentMana < ability.ManaCost {
		return fmt.Errorf("insufficient mana: have %d, need %d", char.CurrentMana, ability.ManaCost)
	}

	// Check level requirement
	if char.Level < ability.UnlockLevel {
		return fmt.Errorf("level too low: have %d, need %d", char.Level, ability.UnlockLevel)
	}

	// Check class compatibility
	if len(ability.RequiredClass) > 0 {
		classMatch := false
		for _, reqClass := range ability.RequiredClass {
			if reqClass == char.Class {
				classMatch = true
				break
			}
		}
		if !classMatch {
			return fmt.Errorf("incompatible class: %s cannot use %s skill", char.Class, ability.Class)
		}
	}

	// Check element compatibility
	if len(ability.RequiredElement) > 0 {
		elementMatch := false
		for _, reqElement := range ability.RequiredElement {
			if reqElement == char.Element {
				elementMatch = true
				break
			}
		}
		if !elementMatch {
			return errors.New("incompatible element")
		}
	}

	// Check if character is alive
	if char.IsFainted || char.IsDead {
		return errors.New("character is unable to act")
	}

	return nil
}

// IsSkillActive checks if a skill is equipped in active slots (1-4)
func (s *SkillActivationService) IsSkillActive(characterID, abilityID uint) bool {
	var char models.Character
	if err := db.DB.Select("equipped_ability_1, equipped_ability_2, equipped_ability_3, equipped_ability_4").
		First(&char, characterID).Error; err != nil {
		return false
	}

	// Check if abilityID matches any of the 4 slots
	if char.EquippedAbility1 != nil && *char.EquippedAbility1 == abilityID {
		return true
	}
	if char.EquippedAbility2 != nil && *char.EquippedAbility2 == abilityID {
		return true
	}
	if char.EquippedAbility3 != nil && *char.EquippedAbility3 == abilityID {
		return true
	}
	if char.EquippedAbility4 != nil && *char.EquippedAbility4 == abilityID {
		return true
	}

	return false
}

// IsOnCooldown checks if a skill is on cooldown
func (s *SkillActivationService) IsOnCooldown(characterID, abilityID uint) bool {
	var cooldown models.CharacterSkillCooldown
	err := db.DB.Where("character_id = ? AND ability_id = ? AND cooldown_remaining > 0", characterID, abilityID).
		First(&cooldown).Error
	return err == nil
}

// DeductMana removes mana from character
func (s *SkillActivationService) DeductMana(char *models.Character, manaCost int) error {
	if char.CurrentMana < manaCost {
		return errors.New("insufficient mana")
	}
	char.CurrentMana -= manaCost
	return nil
}

// CalculateDamage calculates final damage with all modifiers
func (s *SkillActivationService) CalculateDamage(char *models.Character, ability *models.Ability) int {
	baseDamage := ability.BaseDamage

	// Add character attack stat for physical skills
	if ability.DamageType == "physical" {
		baseDamage += char.CurrentAttack / 2
	}

	// Element bonus (simplified - would parse JSON in production)
	elementMultiplier := 1.0
	if ability.Class == char.Class {
		elementMultiplier = 1.2 // 20% bonus for matching class
	}

	return int(float64(baseDamage) * elementMultiplier)
}

// CalculateHealing calculates final healing amount
func (s *SkillActivationService) CalculateHealing(char *models.Character, ability *models.Ability) int {
	baseHeal := ability.BaseHeal

	// Scale with character level
	levelBonus := 1.0 + (float64(char.Level-1) * 0.02)

	return int(float64(baseHeal) * levelBonus)
}

// RollCritical determines if attack is critical
func (s *SkillActivationService) RollCritical(char *models.Character) bool {
	// Base 10% crit chance, +1% per 10 levels
	critChance := 10 + (char.Level / 10)

	// Use cryptographically secure random number generator
	roll := secureRandomInt(100)
	return roll < critChance
}

// secureRandomInt generates a cryptographically secure random integer in range [0, max)
func secureRandomInt(max int) int {
	if max <= 0 {
		return 0
	}

	// Use crypto/rand for security-critical randomness
	nBig, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(max)))
	if err != nil {
		// Fallback to time-based (should never happen in practice)
		return int(time.Now().UnixNano() % int64(max))
	}
	return int(nBig.Int64())
}

// ApplyBuff applies a buff to a character
func (s *SkillActivationService) ApplyBuff(characterID uint, buffType string, duration int) {
	buff := models.CharacterBuff{
		CharacterID:    characterID,
		BuffType:       buffType,
		Multiplier:     1.5, // +50% (would vary by buff type)
		TurnsRemaining: duration,
	}
	db.DB.Create(&buff)
}

// StartCooldown sets a skill on cooldown
func (s *SkillActivationService) StartCooldown(characterID, abilityID uint, cooldownTurns int) {
	now := time.Now()

	// Check if cooldown record exists
	var cooldown models.CharacterSkillCooldown
	err := db.DB.Where("character_id = ? AND ability_id = ?", characterID, abilityID).
		First(&cooldown).Error

	if err != nil {
		// Create new cooldown record
		cooldown = models.CharacterSkillCooldown{
			CharacterID:       characterID,
			AbilityID:         abilityID,
			CooldownRemaining: cooldownTurns,
			LastUsedAt:        &now,
		}
		db.DB.Create(&cooldown)
	} else {
		// Update existing record
		cooldown.CooldownRemaining = cooldownTurns
		cooldown.LastUsedAt = &now
		db.DB.Save(&cooldown)
	}
}

// ReduceCooldowns reduces all cooldowns by 1 (called at start of turn)
func (s *SkillActivationService) ReduceCooldowns(characterID uint) {
	db.DB.Model(&models.CharacterSkillCooldown{}).
		Where("character_id = ? AND cooldown_remaining > 0", characterID).
		UpdateColumn("cooldown_remaining", db.DB.Raw("cooldown_remaining - 1"))
}

// RegenerateMana restores mana at start of turn
func (s *SkillActivationService) RegenerateMana(characterID uint) error {
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return err
	}

	// Restore mana up to max
	character.CurrentMana += character.ManaRegenRate
	if character.CurrentMana > character.MaxMana {
		character.CurrentMana = character.MaxMana
	}

	return db.DB.Save(&character).Error
}

// LogSkillUsage logs skill usage for analytics
func (s *SkillActivationService) LogSkillUsage(req SkillActivationRequest, result *SkillActivationResult) {
	// Update usage count in character_abilities
	db.DB.Model(&models.CharacterAbility{}).
		Where("character_id = ? AND ability_id = ?", req.CharacterID, req.AbilityID).
		UpdateColumn("times_used", db.DB.Raw("times_used + 1"))
}

// GetActiveSkills returns all active skills for a character
func (s *SkillActivationService) GetActiveSkills(characterID uint) ([]models.Ability, error) {
	var char models.Character
	if err := db.DB.First(&char, characterID).Error; err != nil {
		return nil, err
	}

	var abilityIDs []uint
	if char.EquippedAbility1 != nil {
		abilityIDs = append(abilityIDs, *char.EquippedAbility1)
	}
	if char.EquippedAbility2 != nil {
		abilityIDs = append(abilityIDs, *char.EquippedAbility2)
	}
	if char.EquippedAbility3 != nil {
		abilityIDs = append(abilityIDs, *char.EquippedAbility3)
	}
	if char.EquippedAbility4 != nil {
		abilityIDs = append(abilityIDs, *char.EquippedAbility4)
	}

	if len(abilityIDs) == 0 {
		return []models.Ability{}, nil
	}

	var abilities []models.Ability
	if err := db.DB.Where("id IN ?", abilityIDs).Find(&abilities).Error; err != nil {
		return nil, err
	}

	return abilities, nil
}

// SwapActiveSkill swaps a skill in/out of active slots
// Note: This duplicates functionality from AbilityEquipService but is kept for Handler compatibility
// Ideally, the Handler should call AbilityEquipService directly
func (s *SkillActivationService) SwapActiveSkill(userID, characterID, oldAbilityID, newAbilityID uint) error {
	equipService := NewAbilityEquipService()

	// Find which slot holds the old ability
	var char models.Character
	if err := db.DB.First(&char, characterID).Error; err != nil {
		return err
	}

	slotToSwap := 0
	if char.EquippedAbility1 != nil && *char.EquippedAbility1 == oldAbilityID {
		slotToSwap = 1
	} else if char.EquippedAbility2 != nil && *char.EquippedAbility2 == oldAbilityID {
		slotToSwap = 2
	} else if char.EquippedAbility3 != nil && *char.EquippedAbility3 == oldAbilityID {
		slotToSwap = 3
	} else if char.EquippedAbility4 != nil && *char.EquippedAbility4 == oldAbilityID {
		slotToSwap = 4
	}

	if slotToSwap == 0 {
		return errors.New("old ability not found in any active slot")
	}

	return equipService.EquipAbility(userID, characterID, newAbilityID, slotToSwap)
}

// AutoEquipNewSkill equips a newly learned skill if an empty slot is available
func (s *SkillActivationService) AutoEquipNewSkill(tx *gorm.DB, characterID, abilityID uint) error {
	var char models.Character
	if err := tx.First(&char, characterID).Error; err != nil {
		return err
	}

	// Check for empty slots and equip immediately
	if char.EquippedAbility1 == nil {
		char.EquippedAbility1 = &abilityID
	} else if char.EquippedAbility2 == nil {
		char.EquippedAbility2 = &abilityID
	} else if char.EquippedAbility3 == nil {
		char.EquippedAbility3 = &abilityID
	} else if char.EquippedAbility4 == nil {
		char.EquippedAbility4 = &abilityID
	} else {
		// No empty slots
		return nil
	}

	return tx.Save(&char).Error
}

// GetUsableSkills returns skills that are active, off cooldown, and affordable
func (s *SkillActivationService) GetUsableSkills(characterID uint, currentMana int) ([]models.Ability, error) {
	// 1. Get Active Skills
	activeSkills, err := s.GetActiveSkills(characterID)
	if err != nil {
		return nil, err
	}

	var usable []models.Ability
	for _, ability := range activeSkills {
		// Check Mana
		if ability.ManaCost > currentMana {
			continue
		}

		// Check Cooldown
		if s.IsOnCooldown(characterID, ability.ID) {
			continue
		}

		usable = append(usable, ability)
	}

	return usable, nil
}
