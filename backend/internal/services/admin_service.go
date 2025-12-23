package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

type AdminService struct {
	configService *ConfigService
	ls            *LedgerService
}

func NewAdminService() *AdminService {
	return &AdminService{
		configService: GetConfigService(),
		ls:            NewLedgerService(),
	}
}

// User Action: BAN
func (s *AdminService) BanUser(targetUserID uint, reason string, durationHours int, adminID uint) error {
	var user models.User
	if err := db.DB.First(&user, targetUserID).Error; err != nil {
		return errors.New("user not found")
	}

	expiry := time.Now().Add(time.Duration(durationHours) * time.Hour)

	// Update User
	oldStatus := user.Status
	user.Status = "BANNED"
	user.IsBanned = true
	user.BanReason = reason
	user.BanExpiresAt = &expiry
	user.Nonce = "" // Invalidate auth immediately

	if err := db.DB.Save(&user).Error; err != nil {
		return err
	}

	// Audit
	s.CreateAuditLog(adminID, "BAN_USER", strconv.Itoa(int(targetUserID)), oldStatus, "BANNED: "+reason)
	return nil
}

// User Action: UNBAN
func (s *AdminService) UnbanUser(targetUserID uint, adminID uint) error {
	var user models.User
	if err := db.DB.First(&user, targetUserID).Error; err != nil {
		return errors.New("user not found")
	}

	oldStatus := user.Status
	user.Status = "ACTIVE"
	user.IsBanned = false
	user.BanReason = ""
	user.BanExpiresAt = nil

	if err := db.DB.Save(&user).Error; err != nil {
		return err
	}

	s.CreateAuditLog(adminID, "UNBAN_USER", strconv.Itoa(int(targetUserID)), oldStatus, "ACTIVE")
	return nil
}

// User Action: FREEZE FUNDS (Security Flag 1 = Withdrawal Lock)
func (s *AdminService) FreezeFunds(targetUserID uint, freeze bool, adminID uint) error {
	var user models.User
	if err := db.DB.First(&user, targetUserID).Error; err != nil {
		return errors.New("user not found")
	}

	// Toggle Bit 1 (Value 2) for Withdrawal Lock
	// Simplified: We'll use a specific Int flag for now.
	// Assume SecurityFlags: 0=None, 1=WithdrawLock, 2=LoginLock
	oldFlags := strconv.Itoa(user.SecurityFlags)

	if freeze {
		user.SecurityFlags = user.SecurityFlags | 1 // Set bit 1
	} else {
		user.SecurityFlags = user.SecurityFlags &^ 1 // Clear bit 1
	}

	if err := db.DB.Save(&user).Error; err != nil {
		return err
	}

	action := "FREEZE_FUNDS"
	if !freeze {
		action = "UNFREEZE_FUNDS"
	}
	s.CreateAuditLog(adminID, action, strconv.Itoa(int(targetUserID)), oldFlags, strconv.Itoa(user.SecurityFlags))
	return nil
}

// Config Action: UPDATE SETTING
func (s *AdminService) UpdateSystemSetting(key, value, valueType string, adminID uint) error {
	// Get old value for audit
	oldVal := s.configService.GetValue(key, "N/A")

	// Update via Service (Handlers DB + Cache)
	if err := s.configService.SetValue(key, value, valueType, adminID); err != nil {
		return err
	}

	s.CreateAuditLog(adminID, "UPDATE_CONFIG", key, oldVal, value)
	return nil
}

// Report Action: RESOLVE
func (s *AdminService) ResolveReport(reportID uint, action, notes string, adminID uint) error {
	var report models.UserReport
	if err := db.DB.First(&report, reportID).Error; err != nil {
		return errors.New("report not found")
	}

	report.Status = action // RESOLVED, DISMISSED
	report.ResolvedBy = &adminID
	now := time.Now()
	report.ResolvedAt = &now
	report.Details += "\nResolution Note: " + notes

	if err := db.DB.Save(&report).Error; err != nil {
		return err
	}

	s.CreateAuditLog(adminID, "RESOLVE_REPORT", strconv.Itoa(int(reportID)), "PENDING", action)
	return nil
}

// Helper: Audit Log
func (s *AdminService) CreateAuditLog(adminID uint, action, targetID, oldVal, newVal string) {
	log := models.AdminAuditLog{
		AdminID:   adminID,
		Action:    action,
		TargetID:  targetID,
		OldValue:  oldVal,
		NewValue:  newVal,
		CreatedAt: time.Now(),
		// IP/Agent usually come from Context, adding placeholders or extending args if needed
	}
	db.DB.Create(&log)
}

// ListUsers returns all users in the system
func (s *AdminService) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := db.DB.Order("id asc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetAuditLogs returns the most recent administrative actions
func (s *AdminService) GetAuditLogs(limit int) ([]models.AdminAuditLog, error) {
	var logs []models.AdminAuditLog
	if err := db.DB.Order("created_at desc").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetRevenueStats returns total platform revenue from the Treasury account (Legacy Support + Ledger)
func (s *AdminService) GetRevenueStats() (map[string]interface{}, error) {
	// 1. Get Treasury Account
	treasury, err := s.ls.GetOrCreateAccount(nil, models.AccountTypeTreasury, "GTK")
	if err != nil {
		return nil, err
	}
	towerTreasury, err := s.ls.GetOrCreateAccount(nil, models.AccountTypeTreasury, "TOWER")
	if err != nil {
		return nil, err
	}

	// 2. Aggregate Credits (Revenue)
	var totalGTK int64
	var totalTOWER int64

	// LedgerEntry where AccountID = Treasury AND Type = CREDIT
	db.DB.Model(&models.LedgerEntry{}).
		Where("account_id = ? AND type = 'CREDIT'", treasury.ID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalGTK)

	db.DB.Model(&models.LedgerEntry{}).
		Where("account_id = ? AND type = 'CREDIT'", towerTreasury.ID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalTOWER)

	// 3. User & Game Stats
	var totalUsers int64
	db.DB.Model(&models.User{}).Count(&totalUsers)

	var activeBattles int64
	// Check if Battle table exists (it should)
	db.DB.Model(&models.Battle{}).Where("status = ?", "ACTIVE").Count(&activeBattles)

	// Check raid session count if model exists
	var activeRaids int64 = 0
	// Using generic check to avoid error if RaidSession not migrated yet in this context
	// db.DB.Model(&models.RaidSession{}).Where("status = ?", "ACTIVE").Count(&activeRaids)

	return map[string]interface{}{
		"total_revenue_gtk":      totalGTK,
		"total_revenue_tower":    totalTOWER,
		"treasury_balance_gtk":   treasury.Balance,
		"treasury_balance_tower": towerTreasury.Balance,
		"total_users":            totalUsers,
		"active_battles":         activeBattles + activeRaids, // Combined metric
		"active_raids":           activeRaids,
	}, nil
}

// Shop Management
func (s *AdminService) CreateShopItem(item models.ShopItem, adminID uint) error {
	if err := db.DB.Create(&item).Error; err != nil {
		return err
	}
	s.CreateAuditLog(adminID, "CREATE_SHOP_ITEM", strconv.Itoa(int(item.ID)), "", item.Name)
	return nil
}

func (s *AdminService) UpdateShopItem(item models.ShopItem, adminID uint) error {
	var existing models.ShopItem
	if err := db.DB.First(&existing, item.ID).Error; err != nil {
		return err
	}
	if err := db.DB.Save(&item).Error; err != nil {
		return err
	}
	s.CreateAuditLog(adminID, "UPDATE_SHOP_ITEM", strconv.Itoa(int(item.ID)), existing.Name, item.Name)
	return nil
}

func (s *AdminService) DeleteShopItem(itemID uint, adminID uint) error {
	var existing models.ShopItem
	if err := db.DB.First(&existing, itemID).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&existing).Error; err != nil {
		return err
	}
	s.CreateAuditLog(adminID, "DELETE_SHOP_ITEM", strconv.Itoa(int(itemID)), existing.Name, "DELETED")
	return nil
}

func (s *AdminService) ListShopItems() ([]models.ShopItem, error) {
	var items []models.ShopItem
	if err := db.DB.Order("id asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Quest Template Management
func (s *AdminService) CreateQuestTemplate(template models.QuestTemplate, adminID uint) error {
	if err := db.DB.Create(&template).Error; err != nil {
		return err
	}
	s.CreateAuditLog(adminID, "CREATE_QUEST_TEMPLATE", strconv.Itoa(int(template.ID)), "", template.Name)
	return nil
}

func (s *AdminService) UpdateQuestTemplate(template models.QuestTemplate, adminID uint) error {
	var existing models.QuestTemplate
	if err := db.DB.First(&existing, template.ID).Error; err != nil {
		return err
	}
	if err := db.DB.Save(&template).Error; err != nil {
		return err
	}
	s.CreateAuditLog(adminID, "UPDATE_QUEST_TEMPLATE", strconv.Itoa(int(template.ID)), existing.Name, template.Name)
	return nil
}

func (s *AdminService) DeleteQuestTemplate(templateID uint, adminID uint) error {
	var existing models.QuestTemplate
	if err := db.DB.First(&existing, templateID).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&existing).Error; err != nil {
		return err
	}
	s.CreateAuditLog(adminID, "DELETE_QUEST_TEMPLATE", strconv.Itoa(int(templateID)), existing.Name, "DELETED")
	return nil
}

func (s *AdminService) ListQuestTemplates() ([]models.QuestTemplate, error) {
	var templates []models.QuestTemplate
	if err := db.DB.Order("id asc").Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// ==================== ABILITIES MANAGEMENT ====================

// GetAllAbilities returns all abilities for admin management
func (s *AdminService) GetAllAbilities() ([]models.Ability, error) {
	var abilities []models.Ability
	if err := db.DB.Order("class asc, unlock_level asc").Find(&abilities).Error; err != nil {
		return nil, err
	}
	return abilities, nil
}

// CreateAbility creates a new ability
func (s *AdminService) CreateAbility(ability *models.Ability) error {
	if err := db.DB.Create(ability).Error; err != nil {
		return err
	}
	return nil
}

// UpdateAbility updates an existing ability
func (s *AdminService) UpdateAbility(ability *models.Ability) error {
	var existing models.Ability
	if err := db.DB.First(&existing, ability.ID).Error; err != nil {
		return errors.New("ability not found")
	}
	if err := db.DB.Save(ability).Error; err != nil {
		return err
	}
	return nil
}

// DeleteAbility deletes an ability
func (s *AdminService) DeleteAbility(abilityID uint) error {
	var existing models.Ability
	if err := db.DB.First(&existing, abilityID).Error; err != nil {
		return errors.New("ability not found")
	}
	if err := db.DB.Delete(&existing).Error; err != nil {
		return err
	}
	return nil
}
