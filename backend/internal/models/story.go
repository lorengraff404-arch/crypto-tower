package models

import (
	"time"
)

// StoryDialogue represents narrative content for missions
type StoryDialogue struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	MissionLevel int       `gorm:"not null;index" json:"mission_level"`
	DialogueType string    `gorm:"size:20;not null" json:"dialogue_type"` // briefing, post_mission, cutscene
	Character    string    `gorm:"size:50;not null" json:"character"` // aria, kairos, voice, player, narrator
	DialogueText string    `gorm:"type:text;not null" json:"dialogue_text"`
	AudioFile    string    `gorm:"size:255" json:"audio_file"` // Future voice acting
	SortOrder    int       `gorm:"default:0" json:"sort_order"` // For multiple dialogues in sequence
	CreatedAt    time.Time `json:"created_at"`
}

// PlayerChoice tracks narrative choices made by players
type PlayerChoice struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	ChoiceID    string    `gorm:"size:50;not null" json:"choice_id"` // kairos_friend, kairos_enemy, etc.
	ChoiceValue string    `gorm:"size:100" json:"choice_value"` // Selected option
	ChoiceMade  time.Time `gorm:"not null" json:"choice_made"`
	CreatedAt   time.Time `json:"created_at"`
	
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// StoryFragment represents collectible lore items
type StoryFragment struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	FragmentID     string    `gorm:"size:50;not null;unique" json:"fragment_id"` // aria_notes_1, nexus_log_3
	Title          string    `gorm:"size:100;not null" json:"title"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	FragmentType   string    `gorm:"size:30" json:"fragment_type"` // aria_notes, nexus_logs, architect_terminals
	UnlockLevel    int       `gorm:"not null" json:"unlock_level"` // Min level to receive
	UnlockCondition string   `gorm:"type:jsonb" json:"unlock_condition"` // Additional conditions
	Rarity         string    `gorm:"size:10" json:"rarity"` // Common, Rare, Legendary
	CreatedAt      time.Time `json:"created_at"`
}

// UserStoryProgress tracks user's story progression
type UserStoryProgress struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	UserID             uint       `gorm:"not null;uniqueIndex" json:"user_id"`
	CurrentAct         int        `gorm:"default:1" json:"current_act"` // 1-4
	AriaCorruption     int        `gorm:"default:0" json:"aria_corruption"` // 0-100%
	KairosRelationship string     `gorm:"size:20;default:'neutral'" json:"kairos_relationship"` // neutral, friend, enemy
	VoiceEncounters    int        `gorm:"default:0" json:"voice_encounters"` // Times met The Voice
	CollectedFragments string     `gorm:"type:jsonb" json:"collected_fragments"` // Array of fragment IDs
	ViewedCutscenes    string     `gorm:"type:jsonb" json:"viewed_cutscenes"` // Array of cutscene IDs
	LastStoryUpdate    *time.Time `json:"last_story_update"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
