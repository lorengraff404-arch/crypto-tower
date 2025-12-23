package services

import (
	"errors"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// TeamService handles team management logic
type TeamService struct{}

// NewTeamService creates a new team service
func NewTeamService() *TeamService {
	return &TeamService{}
}

// CreateTeam creates a new team for a user
func (s *TeamService) CreateTeam(userID uint, name string) (*models.Team, error) {
	// Check team limit (max 5 teams per user)
	var count int64
	db.DB.Model(&models.Team{}).Where("user_id = ?", userID).Count(&count)
	if count >= 5 {
		return nil, errors.New("max team limit reached (5)")
	}

	team := models.Team{
		UserID:    userID,
		Name:      name,
		IsDefault: count == 0, // First team is default
	}

	if err := db.DB.Create(&team).Error; err != nil {
		return nil, err
	}

	return &team, nil
}

// AddMember adds a character to a team slot
func (s *TeamService) AddMember(teamID, characterID uint, slot int, isBackup bool) error {
	var team models.Team
	if err := db.DB.First(&team, teamID).Error; err != nil {
		return errors.New("team not found")
	}

	// Verify character ownership
	var character models.Character
	if err := db.DB.First(&character, characterID).Error; err != nil {
		return errors.New("character not found")
	}
	if character.OwnerID != team.UserID {
		return errors.New("character does not belong to user")
	}

	// Check if character is listed on marketplace
	if character.IsListed {
		return errors.New("cannot add listed character to team")
	}

	// Validate slot
	if isBackup {
		if slot < 0 || slot > 2 {
			return errors.New("invalid backup slot (0-2)")
		}
	} else {
		if slot < 0 || slot > 2 {
			return errors.New("invalid active slot (0-2)")
		}
	}

	// Check if character is already in team
	var existingMember models.TeamMember
	err := db.DB.Where("team_id = ? AND character_id = ?", teamID, characterID).First(&existingMember).Error
	if err == nil {
		return errors.New("character already in team")
	}

	// Check if slot is occupied
	db.DB.Where("team_id = ? AND slot = ? AND is_backup = ?", teamID, slot, isBackup).Delete(&models.TeamMember{})

	member := models.TeamMember{
		TeamID:      teamID,
		CharacterID: characterID,
		Slot:        slot,
		IsBackup:    isBackup,
	}

	return db.DB.Create(&member).Error
}

// RemoveMember removes a character from a team
func (s *TeamService) RemoveMember(teamID, characterID uint) error {
	return db.DB.Where("team_id = ? AND character_id = ?", teamID, characterID).
		Delete(&models.TeamMember{}).Error
}

// GetUserTeams returns all teams for a user with members loaded
func (s *TeamService) GetUserTeams(userID uint) ([]models.Team, error) {
	var teams []models.Team
	err := db.DB.Preload("Members.Character").
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at ASC").
		Find(&teams).Error

	if err != nil {
		return nil, err
	}

	// Calculate stats and synergies for each team
	for i := range teams {
		s.calculateTeamStats(&teams[i])
	}

	return teams, nil
}

func (s *TeamService) calculateTeamStats(team *models.Team) {
	totalPower := 0
	totalLevel := 0
	activeCount := 0

	for _, member := range team.Members {
		// Only active members contribute to main stats
		if !member.IsBackup {
			char := member.Character
			cp := char.CurrentAttack + char.CurrentDefense + int(char.CurrentHP) + char.CurrentSpeed
			totalPower += cp
			totalLevel += char.Level
			activeCount++
		}
	}

	team.TotalPower = totalPower
	if activeCount > 0 {
		team.AvgLevel = totalLevel / activeCount
	}

	// Calculate Synergies
	team.Synergies = s.calculateSynergies(team)
}

func (s *TeamService) calculateSynergies(team *models.Team) []models.TeamSynergy {
	synergies := []models.TeamSynergy{}

	// Count types and elements
	typeCounts := make(map[string]int)
	elementCounts := make(map[string]int)
	classCounts := make(map[string]int)

	for _, member := range team.Members {
		if !member.IsBackup {
			char := member.Character
			typeCounts[char.CharacterType]++
			elementCounts[char.Element]++
			classCounts[char.Class]++
		}
	}

	// Mono-Type Synergy (3 same type)
	for t, count := range typeCounts {
		if count >= 3 {
			synergies = append(synergies, models.TeamSynergy{
				ID:          "MONO_TYPE_" + t,
				Name:        "Tribal Unity (" + t + ")",
				Description: "All team members share the same lineage.",
				Tier:        3,
				Icon:        "👑",
				Condition:   "3 " + t + " Types",
				BonusStat:   "ALL_STATS",
				BonusValue:  0.15, // +15% All Stats
			})
		}
	}

	// Elemental Harmony (3 same element)
	for e, count := range elementCounts {
		if count >= 3 {
			synergies = append(synergies, models.TeamSynergy{
				ID:          "ELEMENTAL_HARMONY_" + e,
				Name:        "Elemental Harmony (" + e + ")",
				Description: "The team resonates with a single element.",
				Tier:        3,
				Icon:        "✨",
				Condition:   "3 " + e + " Elements",
				BonusStat:   "RESISTANCE",
				BonusValue:  0.20, // +20% Resistance
			})
		}
	}

	// Trinity (Warrior + Mage + Tank)
	if classCounts["Warrior"] >= 1 && classCounts["Mage"] >= 1 && classCounts["Tank"] >= 1 {
		synergies = append(synergies, models.TeamSynergy{
			ID:          "TRINITY_BALANCE",
			Name:        "Perfect Trinity",
			Description: "A perfectly balanced team composition.",
			Tier:        3,
			Icon:        "⚖️",
			Condition:   "Warrior + Mage + Tank",
			BonusStat:   "DAMAGE",
			BonusValue:  0.10, // +10% Damage
		})
	}

	return synergies
}
