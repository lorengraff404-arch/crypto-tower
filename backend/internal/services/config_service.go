package services

import (
	"strconv"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// GetConfig returns string value
func (s *ConfigService) GetConfig(key string, defaultValue string) string {
	var setting models.SystemSetting
	// TODO: Add caching layer (Redis or in-memory map with TTL)
	if err := db.DB.First(&setting, "key = ?", key).Error; err != nil {
		return defaultValue
	}
	return setting.Value
}

// GetInt returns int value (helper)
func (s *ConfigService) GetInt(key string, defaultValue int) int {
	valStr := s.GetConfig(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

// GetFloat returns float value (helper)
func (s *ConfigService) GetFloat(key string, defaultValue float64) float64 {
	valStr := s.GetConfig(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

// SetConfig updates or creates a setting
func (s *ConfigService) SetConfig(key, value, description string, adminID uint) error {
	setting := models.SystemSetting{
		Key:         key,
		Value:       value,
		Description: description,
		UpdatedBy:   adminID,
		UpdatedAt:   time.Now(),
	}
	return db.DB.Save(&setting).Error // Upsert (Save handles update if PK exists, but here Key isn't PK unless defined in GORM model. We might need Where().Assign())
}

// GetAllSettings returns all settings
func (s *ConfigService) GetAllSettings() ([]models.SystemSetting, error) {
	var settings []models.SystemSetting
	if err := db.DB.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}
