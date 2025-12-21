package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

type MissionService struct {
	db *gorm.DB
}

func NewMissionService(db *gorm.DB) *MissionService {
	return &MissionService{db: db}
}

// GetAvailableMissions returns all missions available for the user
func (s *MissionService) GetAvailableMissions(userID uint) ([]models.UserMissionProgress, error) {
	var progress []models.UserMissionProgress

	err := s.db.
		Preload("Mission").
		Where("user_id = ? AND (status = 'available' OR status = 'in_progress')", userID).
		Order("mission_id ASC").
		Find(&progress).Error

	return progress, err
}

// GetCurrentMission returns the user's current active mission
func (s *MissionService) GetCurrentMission(userID uint) (*models.UserMissionProgress, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var progress models.UserMissionProgress
	err := s.db.
		Preload("Mission").
		Where("user_id = ? AND status = 'in_progress'", userID).
		First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		// No active mission, get next available
		return s.getNextAvailableMission(userID, user.Level)
	}

	return &progress, err
}

// getNextAvailableMission finds the next mission the user should do
func (s *MissionService) getNextAvailableMission(userID uint, userLevel int) (*models.UserMissionProgress, error) {
	var mission models.Mission

	// Find first mission at or above user level that isn't completed
	err := s.db.
		Where("level = ? AND is_active = true", userLevel).
		First(&mission).Error

	if err != nil {
		return nil, err
	}

	// Create progress entry
	progress := &models.UserMissionProgress{
		UserID:    userID,
		MissionID: mission.ID,
		Status:    "available",
		Mission:   mission,
	}

	// Initialize objectives progress
	var objectives []models.Objective
	if err := json.Unmarshal([]byte(mission.Objectives), &objectives); err == nil {
		progressJSON, _ := json.Marshal(objectives)
		progress.ObjectivesProgress = string(progressJSON)
		s.db.Create(progress)
	}

	return progress, nil
}

// StartMission starts a mission for the user
func (s *MissionService) StartMission(userID, missionID uint) error {
	var progress models.UserMissionProgress

	err := s.db.
		Where("user_id = ? AND mission_id = ?", userID, missionID).
		First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		// Create new progress
		var mission models.Mission
		if err := s.db.First(&mission, missionID).Error; err != nil {
			return err
		}

		now := time.Now()
		progress = models.UserMissionProgress{
			UserID:    userID,
			MissionID: missionID,
			Status:    "in_progress",
			StartedAt: &now,
		}

		// Initialize objectives
		var objectives []models.Objective
		if err := json.Unmarshal([]byte(mission.Objectives), &objectives); err == nil {
			progressJSON, _ := json.Marshal(objectives)
			progress.ObjectivesProgress = string(progressJSON)
		}

		return s.db.Create(&progress).Error
	}

	// Update existing progress
	now := time.Now()
	return s.db.Model(&progress).Updates(map[string]interface{}{
		"status":     "in_progress",
		"started_at": &now,
	}).Error
}

// UpdateProgress updates progress for a specific objective type
func (s *MissionService) UpdateProgress(userID, missionID uint, objectiveType string, increment int) error {
	var progress models.UserMissionProgress

	err := s.db.
		Where("user_id = ? AND mission_id = ? AND status = 'in_progress'", userID, missionID).
		First(&progress).Error

	if err != nil {
		return err
	}

	// Parse current objectives
	var objectives []models.Objective
	if err := json.Unmarshal([]byte(progress.ObjectivesProgress), &objectives); err != nil {
		return err
	}

	// Update matching objective
	updated := false
	for i := range objectives {
		if objectives[i].Type == objectiveType && objectives[i].Current < objectives[i].Target {
			objectives[i].Current += increment
			if objectives[i].Current > objectives[i].Target {
				objectives[i].Current = objectives[i].Target
			}
			updated = true
		}
	}

	if !updated {
		return nil // No matching objective
	}

	// Save updated progress
	progressJSON, _ := json.Marshal(objectives)
	err = s.db.Model(&progress).Update("objectives_progress", string(progressJSON)).Error
	if err != nil {
		return err
	}

	// Check if mission is complete
	s.CheckCompletion(userID, missionID)

	return nil
}

// CheckCompletion checks if all objectives are complete
func (s *MissionService) CheckCompletion(userID, missionID uint) (bool, error) {
	var progress models.UserMissionProgress

	err := s.db.
		Preload("Mission").
		Where("user_id = ? AND mission_id = ?", userID, missionID).
		First(&progress).Error

	if err != nil {
		return false, err
	}

	// Parse objectives
	var objectives []models.Objective
	if err := json.Unmarshal([]byte(progress.ObjectivesProgress), &objectives); err != nil {
		return false, err
	}

	// Check if all complete
	allComplete := true
	for _, obj := range objectives {
		if obj.Current < obj.Target {
			allComplete = false
			break
		}
	}

	if allComplete && progress.Status != "completed" {
		// Auto-complete mission
		s.CompleteMission(userID, missionID)
		return true, nil
	}

	return allComplete, nil
}

// CompleteMission marks mission as complete and distributes rewards
func (s *MissionService) CompleteMission(userID, missionID uint) (*models.MissionRewards, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var progress models.UserMissionProgress
	err := tx.
		Preload("Mission").
		Where("user_id = ? AND mission_id = ?", userID, missionID).
		First(&progress).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if progress.Status == "completed" {
		tx.Rollback()
		return nil, errors.New("mission already completed")
	}

	// Mark as completed
	now := time.Now()
	tx.Model(&progress).Updates(map[string]interface{}{
		"status":       "completed",
		"completed_at": &now,
	})

	// Parse and distribute rewards
	var rewards models.MissionRewards
	if err := json.Unmarshal([]byte(progress.Mission.Rewards), &rewards); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Add GTK
	if rewards.GTK > 0 {
		tx.Model(&models.User{}).
			Where("id = ?", userID).
			Update("gtk_balance", gorm.Expr("gtk_balance + ?", rewards.GTK))
	}

	// Update user level and missions completed
	tx.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"current_level":            progress.Mission.Level + 1,
			"total_missions_completed": gorm.Expr("total_missions_completed + 1"),
		})

	// Unlock feature if applicable
	if progress.Mission.UnlockFeature != "" {
		var user models.User
		tx.First(&user, userID)

		var unlocked []string
		json.Unmarshal([]byte(""), &unlocked)
		unlocked = append(unlocked, progress.Mission.UnlockFeature)
		unlockedJSON, _ := json.Marshal(unlocked)

		tx.Model(&models.User{}).
			Where("id = ?", userID).
			Update("unlocked_features", string(unlockedJSON))
	}

	tx.Commit()
	return &rewards, nil
}

// CheckFeatureUnlocked checks if a feature is unlocked for user
func (s *MissionService) CheckFeatureUnlocked(userID uint, feature string) (bool, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return false, err
	}

	var unlocked []string
	if err := json.Unmarshal([]byte(""), &unlocked); err != nil {
		return false, nil // Default to locked if can't parse
	}

	for _, f := range unlocked {
		if f == feature {
			return true, nil
		}
	}

	return false, nil
}

// Event hooks for automatic progress tracking

func (s *MissionService) OnBattleComplete(userID uint, waves int, hasJefe bool) error {
	// Find active missions
	var progressList []models.UserMissionProgress
	s.db.Where("user_id = ? AND status = 'in_progress'", userID).Find(&progressList)

	for _, progress := range progressList {
		// Update battle_waves objective
		s.UpdateProgress(userID, progress.MissionID, "battle_waves", waves)

		if hasJefe {
			s.UpdateProgress(userID, progress.MissionID, "defeat_boss", 1)
		}
	}

	return nil
}

func (s *MissionService) OnEggHatched(userID uint) error {
	var progressList []models.UserMissionProgress
	s.db.Where("user_id = ? AND status = 'in_progress'", userID).Find(&progressList)

	for _, progress := range progressList {
		s.UpdateProgress(userID, progress.MissionID, "hatch_egg", 1)
	}

	return nil
}

func (s *MissionService) OnBreeding(userID uint) error {
	// Check if breeding is unlocked
	unlocked, _ := s.CheckFeatureUnlocked(userID, "breeding")
	if !unlocked {
		return errors.New("breeding unlocks at level 5")
	}

	var progressList []models.UserMissionProgress
	s.db.Where("user_id = ? AND status = 'in_progress'", userID).Find(&progressList)

	for _, progress := range progressList {
		s.UpdateProgress(userID, progress.MissionID, "breeding", 1)
	}

	return nil
}

func (s *MissionService) OnCraft(userID uint) error {
	// Check if crafting is unlocked
	unlocked, _ := s.CheckFeatureUnlocked(userID, "crafting")
	if !unlocked {
		return errors.New("crafting unlocks at level 10")
	}

	var progressList []models.UserMissionProgress
	s.db.Where("user_id = ? AND status = 'in_progress'", userID).Find(&progressList)

	for _, progress := range progressList {
		s.UpdateProgress(userID, progress.MissionID, "craft", 1)
	}

	return nil
}

func (s *MissionService) OnUnitDeployed(userID uint) error {
	var progressList []models.UserMissionProgress
	s.db.Where("user_id = ? AND status = 'in_progress'", userID).Find(&progressList)

	for _, progress := range progressList {
		s.UpdateProgress(userID, progress.MissionID, "deploy_units", 1)
	}

	return nil
}

func (s *MissionService) OnSpellCast(userID uint) error {
	var progressList []models.UserMissionProgress
	s.db.Where("user_id = ? AND status = 'in_progress'", userID).Find(&progressList)

	for _, progress := range progressList {
		s.UpdateProgress(userID, progress.MissionID, "cast_spells", 1)
	}

	return nil
}
