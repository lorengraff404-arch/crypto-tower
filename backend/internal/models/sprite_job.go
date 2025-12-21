package models

import (
	"time"
)

// SpriteGenerationJob represents an async job for generating character sprites
type SpriteGenerationJob struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CharacterID uint       `gorm:"not null;index" json:"character_id"`
	Status      string     `gorm:"size:20;not null;index" json:"status"` // pending, processing, completed, failed
	Progress    int        `gorm:"default:0" json:"progress"`            // 0-100
	ErrorMsg    string     `gorm:"type:text" json:"error_msg,omitempty"`
	Provider    string     `gorm:"size:50" json:"provider"` // openai, mock, etc.
	RetryCount  int        `gorm:"default:0" json:"retry_count"`
	MaxRetries  int        `gorm:"default:3" json:"max_retries"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Relationships
	Character Character `gorm:"foreignKey:CharacterID" json:"character,omitempty"`
}

// SpriteAnimationType represents the type of animation sprite
type SpriteAnimationType string

const (
	SpriteIdle    SpriteAnimationType = "idle"
	SpriteWalk    SpriteAnimationType = "walk"
	SpriteRun     SpriteAnimationType = "run"
	SpriteAttack  SpriteAnimationType = "attack"
	SpriteSkill   SpriteAnimationType = "skill"
	SpriteHit     SpriteAnimationType = "hit"
	SpriteBlock   SpriteAnimationType = "block"
	SpriteDodge   SpriteAnimationType = "dodge"
	SpriteDeath   SpriteAnimationType = "death"
	SpriteVictory SpriteAnimationType = "victory"
)

// AnimationSpec defines the specifications for each animation type
type AnimationSpec struct {
	Type        SpriteAnimationType
	FrameCount  int
	IsLoop      bool
	Description string
}

// GetAnimationSpecs returns specifications for all animation types
func GetAnimationSpecs() []AnimationSpec {
	return []AnimationSpec{
		{Type: SpriteIdle, FrameCount: 8, IsLoop: true, Description: "Breathing, micro movement, weapon sway"},
		{Type: SpriteWalk, FrameCount: 8, IsLoop: true, Description: "Walking cycle, arms and legs alternating"},
		{Type: SpriteRun, FrameCount: 8, IsLoop: true, Description: "Running with longer stride, leaning forward"},
		{Type: SpriteAttack, FrameCount: 10, IsLoop: false, Description: "2f anticipation, 5f strike, 3f recovery"},
		{Type: SpriteSkill, FrameCount: 10, IsLoop: false, Description: "Channeling energy, culminates in release"},
		{Type: SpriteHit, FrameCount: 6, IsLoop: false, Description: "Recoil, damage expression, return to idle"},
		{Type: SpriteBlock, FrameCount: 6, IsLoop: true, Description: "Defensive stance, slight tremor on impact"},
		{Type: SpriteDodge, FrameCount: 8, IsLoop: false, Description: "Quick lateral dash with motion blur"},
		{Type: SpriteDeath, FrameCount: 10, IsLoop: false, Description: "Losing strength, falling, final frame still"},
		{Type: SpriteVictory, FrameCount: 10, IsLoop: true, Description: "Triumphant pose, weapon gesture"},
	}
}

// TableName specifies the table name for SpriteGenerationJob
func (SpriteGenerationJob) TableName() string {
	return "sprite_generation_jobs"
}
