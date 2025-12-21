package services

import (
	"encoding/json"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// AuditService handles audit logging for security and debugging
type AuditService struct{}

// LogAction logs a user action to audit trail
func (s *AuditService) LogAction(userID uint, action, entityType string, entityID uint, oldValues, newValues interface{}, ipAddress, userAgent string) error {
	oldJSON, _ := json.Marshal(oldValues)
	newJSON, _ := json.Marshal(newValues)

	log := models.AuditLog{
		UserID:     &userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   &entityID,
		OldValues:  string(oldJSON),
		NewValues:  string(newJSON),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	return db.DB.Create(&log).Error
}

// RateLimitService prevents abuse
type RateLimitService struct{}

// CheckRateLimit checks if action is allowed
func (s *RateLimitService) CheckRateLimit(userID uint, actionType string, maxCount int, windowMinutes int) (bool, error) {
	var limit models.RateLimit
	err := db.DB.Where("user_id = ? AND action_type = ?", userID, actionType).First(&limit).Error

	if err != nil {
		// Create new limit
		limit = models.RateLimit{
			UserID:     userID,
			ActionType: actionType,
			Count:      1,
		}
		db.DB.Create(&limit)
		return true, nil
	}

	// Check if window expired
	windowDuration := time.Duration(windowMinutes) * time.Minute
	if time.Since(limit.WindowStart) > windowDuration {
		// Reset window
		limit.Count = 1
		limit.WindowStart = time.Now()
		db.DB.Save(&limit)
		return true, nil
	}

	// Check count
	if limit.Count >= maxCount {
		return false, nil
	}

	// Increment count
	limit.Count++
	db.DB.Save(&limit)
	return true, nil
}

// NotificationService handles user notifications
type NotificationService struct{}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(userID uint, notifType, title, message string, data interface{}) error {
	dataJSON, _ := json.Marshal(data)

	notif := models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		Data:    string(dataJSON),
		Read:    false,
	}

	return db.DB.Create(&notif).Error
}

// GetUnreadNotifications returns unread notifications for user
func (s *NotificationService) GetUnreadNotifications(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := db.DB.Where("user_id = ? AND read = ?", userID, false).
		Order("created_at DESC").
		Limit(50).
		Find(&notifications).Error

	return notifications, err
}

// MarkAsRead marks notification as read
func (s *NotificationService) MarkAsRead(userID, notificationID uint) error {
	return db.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("read", true).Error
}

// LeaderboardService handles rankings
type LeaderboardService struct{}

// UpdateLeaderboard updates user's leaderboard entry
func (s *LeaderboardService) UpdateLeaderboard(userID uint, category string, score int) error {
	season := s.getCurrentSeason()

	var entry models.Leaderboard
	err := db.DB.Where("user_id = ? AND category = ? AND season = ?", userID, category, season).First(&entry).Error

	if err != nil {
		// Create new entry
		entry = models.Leaderboard{
			UserID:   userID,
			Category: category,
			Score:    score,
			Season:   season,
		}
		db.DB.Create(&entry)
	} else {
		// Update if score is better
		if score > entry.Score {
			entry.Score = score
			entry.UpdatedAt = time.Now()
			db.DB.Save(&entry)
		}
	}

	// Recalculate ranks
	s.recalculateRanks(category, season)

	return nil
}

// GetTopPlayers returns top players for category
func (s *LeaderboardService) GetTopPlayers(category string, limit int) ([]models.Leaderboard, error) {
	season := s.getCurrentSeason()
	var entries []models.Leaderboard

	err := db.DB.Preload("User").
		Where("category = ? AND season = ?", category, season).
		Order("rank ASC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

// getCurrentSeason returns current season number
func (s *LeaderboardService) getCurrentSeason() int {
	// Season changes every 3 months
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	return (year-2024)*4 + (month-1)/3 + 1
}

// recalculateRanks recalculates all ranks for category
func (s *LeaderboardService) recalculateRanks(category string, season int) {
	db.DB.Exec(`
		UPDATE leaderboards l1
		SET rank = (
			SELECT COUNT(*) + 1
			FROM leaderboards l2
			WHERE l2.category = l1.category
			AND l2.season = l1.season
			AND l2.score > l1.score
		)
		WHERE category = ? AND season = ?
	`, category, season)
}
