package services

import (
	"strconv"
	"sync"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// ConfigService manages dynamic system settings with caching
type ConfigService struct {
	cache      map[string]cachedSetting
	cacheMutex sync.RWMutex
}

type cachedSetting struct {
	Value     string
	ExpiresAt time.Time
}

var (
	CacheTTL = 1 * time.Minute
	instance *ConfigService
	once     sync.Once
)

// GetConfigService returns the singleton instance
func GetConfigService() *ConfigService {
	once.Do(func() {
		instance = &ConfigService{
			cache: make(map[string]cachedSetting),
		}
	})
	return instance
}

// GetValue retrieves a setting string, checking cache first
func (s *ConfigService) GetValue(key, defaultValue string) string {
	// 1. Check Cache
	s.cacheMutex.RLock()
	cached, found := s.cache[key]
	s.cacheMutex.RUnlock()

	if found && time.Now().Before(cached.ExpiresAt) {
		return cached.Value
	}

	// 2. Fetch from DB
	var setting models.SystemSetting
	if err := db.DB.First(&setting, "key = ?", key).Error; err != nil {
		return defaultValue
	}

	// 3. Update Cache
	s.cacheMutex.Lock()
	s.cache[key] = cachedSetting{
		Value:     setting.Value,
		ExpiresAt: time.Now().Add(CacheTTL),
	}
	s.cacheMutex.Unlock()

	return setting.Value
}

// SetValue updates a setting in DB and invalidates cache
func (s *ConfigService) SetValue(key, value, valueType string, adminID uint) error {
	setting := models.SystemSetting{
		Key:       key,
		Value:     value,
		Type:      valueType,
		UpdatedBy: adminID,
		UpdatedAt: time.Now(),
	}

	// Upsert (Postgres compatible)
	if err := db.DB.Save(&setting).Error; err != nil {
		return err
	}

	// Invalidate Cache
	s.cacheMutex.Lock()
	delete(s.cache, key)
	s.cacheMutex.Unlock()

	return nil
}

// Typed Helpers

func (s *ConfigService) GetInt(key string, defaultVal int) int {
	valStr := s.GetValue(key, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func (s *ConfigService) GetFloat(key string, defaultVal float64) float64 {
	valStr := s.GetValue(key, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

func (s *ConfigService) GetBool(key string, defaultVal bool) bool {
	valStr := s.GetValue(key, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetAllSettings returns all settings (Admin Console)
func (s *ConfigService) GetAllSettings() ([]models.SystemSetting, error) {
	var settings []models.SystemSetting
	if err := db.DB.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}
