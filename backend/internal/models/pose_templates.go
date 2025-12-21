package models

import (
	"encoding/base64"
	"fmt"
)

// PoseTemplate represents a predefined pose skeleton for ControlNet
type PoseTemplate struct {
	AnimationType SpriteAnimationType
	FrameIndex    int
	PoseData      string // Base64 encoded OpenPose skeleton JSON or image
	Description   string
}

// GetPoseTemplates returns predefined pose skeletons for all animations
// These are used with ControlNet to ensure consistent character appearance
func GetPoseTemplates() map[SpriteAnimationType][]PoseTemplate {
	return map[SpriteAnimationType][]PoseTemplate{
		SpriteIdle: {
			{AnimationType: SpriteIdle, FrameIndex: 0, Description: "Standing neutral, weight centered"},
			{AnimationType: SpriteIdle, FrameIndex: 1, Description: "Slight lean right, breathing in"},
			{AnimationType: SpriteIdle, FrameIndex: 2, Description: "Weight shifted right"},
			{AnimationType: SpriteIdle, FrameIndex: 3, Description: "Return to center"},
			{AnimationType: SpriteIdle, FrameIndex: 4, Description: "Slight lean left, breathing out"},
			{AnimationType: SpriteIdle, FrameIndex: 5, Description: "Weight shifted left"},
			{AnimationType: SpriteIdle, FrameIndex: 6, Description: "Return to center"},
			{AnimationType: SpriteIdle, FrameIndex: 7, Description: "Back to neutral (loop)"},
		},
		SpriteWalk: {
			{AnimationType: SpriteWalk, FrameIndex: 0, Description: "Contact: right foot forward"},
			{AnimationType: SpriteWalk, FrameIndex: 1, Description: "Recoil: weight shifts forward"},
			{AnimationType: SpriteWalk, FrameIndex: 2, Description: "Passing: left leg passes right"},
			{AnimationType: SpriteWalk, FrameIndex: 3, Description: "High point: left foot lifting"},
			{AnimationType: SpriteWalk, FrameIndex: 4, Description: "Contact: left foot forward"},
			{AnimationType: SpriteWalk, FrameIndex: 5, Description: "Recoil: weight shifts forward"},
			{AnimationType: SpriteWalk, FrameIndex: 6, Description: "Passing: right leg passes left"},
			{AnimationType: SpriteWalk, FrameIndex: 7, Description: "High point: right foot lifting (loop)"},
		},
		SpriteRun: {
			{AnimationType: SpriteRun, FrameIndex: 0, Description: "Contact: right foot down, leaning forward"},
			{AnimationType: SpriteRun, FrameIndex: 1, Description: "Push off: right leg extends"},
			{AnimationType: SpriteRun, FrameIndex: 2, Description: "Flight: both feet off ground"},
			{AnimationType: SpriteRun, FrameIndex: 3, Description: "Reach: left leg extending forward"},
			{AnimationType: SpriteRun, FrameIndex: 4, Description: "Contact: left foot down"},
			{AnimationType: SpriteRun, FrameIndex: 5, Description: "Push off: left leg extends"},
			{AnimationType: SpriteRun, FrameIndex: 6, Description: "Flight: both feet off ground"},
			{AnimationType: SpriteRun, FrameIndex: 7, Description: "Reach: right leg extending (loop)"},
		},
		SpriteAttack: {
			{AnimationType: SpriteAttack, FrameIndex: 0, Description: "Anticipation: weapon pulled back, torso twisted"},
			{AnimationType: SpriteAttack, FrameIndex: 1, Description: "Anticipation: coiled, ready to strike"},
			{AnimationType: SpriteAttack, FrameIndex: 2, Description: "Start swing: hips rotate forward"},
			{AnimationType: SpriteAttack, FrameIndex: 3, Description: "Mid swing: torso follows"},
			{AnimationType: SpriteAttack, FrameIndex: 4, Description: "Impact: weapon at full extension"},
			{AnimationType: SpriteAttack, FrameIndex: 5, Description: "Follow through: weapon continues arc"},
			{AnimationType: SpriteAttack, FrameIndex: 6, Description: "Follow through: body rotation completes"},
			{AnimationType: SpriteAttack, FrameIndex: 7, Description: "Recovery: weapon begins return"},
			{AnimationType: SpriteAttack, FrameIndex: 8, Description: "Recovery: returning to guard"},
			{AnimationType: SpriteAttack, FrameIndex: 9, Description: "Return to idle stance"},
		},
		SpriteSkill: {
			{AnimationType: SpriteSkill, FrameIndex: 0, Description: "Preparation: hands/weapon gathering energy"},
			{AnimationType: SpriteSkill, FrameIndex: 1, Description: "Channeling: arms extending, energy building"},
			{AnimationType: SpriteSkill, FrameIndex: 2, Description: "Channeling: stance widening, power accumulating"},
			{AnimationType: SpriteSkill, FrameIndex: 3, Description: "Channeling: peak energy, body tense"},
			{AnimationType: SpriteSkill, FrameIndex: 4, Description: "Channeling: holding power"},
			{AnimationType: SpriteSkill, FrameIndex: 5, Description: "Pre-release: body coiling"},
			{AnimationType: SpriteSkill, FrameIndex: 6, Description: "Release: arms thrust forward"},
			{AnimationType: SpriteSkill, FrameIndex: 7, Description: "Release: full extension, energy unleashed"},
			{AnimationType: SpriteSkill, FrameIndex: 8, Description: "Follow through: arms still extended"},
			{AnimationType: SpriteSkill, FrameIndex: 9, Description: "Recovery: returning to stance"},
		},
		SpriteHit: {
			{AnimationType: SpriteHit, FrameIndex: 0, Description: "Impact: body recoiling, head snapping back"},
			{AnimationType: SpriteHit, FrameIndex: 1, Description: "Recoil: stumbling backward"},
			{AnimationType: SpriteHit, FrameIndex: 2, Description: "Recoil: off balance, arms flailing"},
			{AnimationType: SpriteHit, FrameIndex: 3, Description: "Recovery: regaining balance"},
			{AnimationType: SpriteHit, FrameIndex: 4, Description: "Recovery: stabilizing"},
			{AnimationType: SpriteHit, FrameIndex: 5, Description: "Return to idle stance"},
		},
		SpriteBlock: {
			{AnimationType: SpriteBlock, FrameIndex: 0, Description: "Guard up: shield/weapon raised"},
			{AnimationType: SpriteBlock, FrameIndex: 1, Description: "Braced: knees bent, ready for impact"},
			{AnimationType: SpriteBlock, FrameIndex: 2, Description: "Impact: slight push back"},
			{AnimationType: SpriteBlock, FrameIndex: 3, Description: "Holding: maintaining guard"},
			{AnimationType: SpriteBlock, FrameIndex: 4, Description: "Tremor: slight shake from impact"},
			{AnimationType: SpriteBlock, FrameIndex: 5, Description: "Holding guard (loop to frame 3)"},
		},
		SpriteDodge: {
			{AnimationType: SpriteDodge, FrameIndex: 0, Description: "Anticipation: body tensing, preparing to move"},
			{AnimationType: SpriteDodge, FrameIndex: 1, Description: "Push off: legs extending, starting dash"},
			{AnimationType: SpriteDodge, FrameIndex: 2, Description: "Dash: body horizontal, moving fast"},
			{AnimationType: SpriteDodge, FrameIndex: 3, Description: "Dash: mid-movement, blurred"},
			{AnimationType: SpriteDodge, FrameIndex: 4, Description: "Dash: continuing lateral movement"},
			{AnimationType: SpriteDodge, FrameIndex: 5, Description: "Deceleration: slowing down"},
			{AnimationType: SpriteDodge, FrameIndex: 6, Description: "Landing: feet planting"},
			{AnimationType: SpriteDodge, FrameIndex: 7, Description: "Stable: back to ready stance"},
		},
		SpriteDeath: {
			{AnimationType: SpriteDeath, FrameIndex: 0, Description: "Hit: initial impact, body jerking"},
			{AnimationType: SpriteDeath, FrameIndex: 1, Description: "Weakening: knees buckling"},
			{AnimationType: SpriteDeath, FrameIndex: 2, Description: "Falling: body tilting forward"},
			{AnimationType: SpriteDeath, FrameIndex: 3, Description: "Falling: arms dropping"},
			{AnimationType: SpriteDeath, FrameIndex: 4, Description: "Falling: knees hitting ground"},
			{AnimationType: SpriteDeath, FrameIndex: 5, Description: "Collapse: torso falling"},
			{AnimationType: SpriteDeath, FrameIndex: 6, Description: "Collapse: body horizontal"},
			{AnimationType: SpriteDeath, FrameIndex: 7, Description: "Down: lying on ground"},
			{AnimationType: SpriteDeath, FrameIndex: 8, Description: "Final: slight settling"},
			{AnimationType: SpriteDeath, FrameIndex: 9, Description: "Still: completely motionless"},
		},
		SpriteVictory: {
			{AnimationType: SpriteVictory, FrameIndex: 0, Description: "Triumphant: weapon raised high"},
			{AnimationType: SpriteVictory, FrameIndex: 1, Description: "Celebration: arms extending"},
			{AnimationType: SpriteVictory, FrameIndex: 2, Description: "Celebration: chest out, proud stance"},
			{AnimationType: SpriteVictory, FrameIndex: 3, Description: "Gesture: weapon flourish"},
			{AnimationType: SpriteVictory, FrameIndex: 4, Description: "Gesture: weapon spinning/twirling"},
			{AnimationType: SpriteVictory, FrameIndex: 5, Description: "Gesture: weapon catch"},
			{AnimationType: SpriteVictory, FrameIndex: 6, Description: "Pose: striking victory pose"},
			{AnimationType: SpriteVictory, FrameIndex: 7, Description: "Pose: holding pose"},
			{AnimationType: SpriteVictory, FrameIndex: 8, Description: "Slight movement: breathing"},
			{AnimationType: SpriteVictory, FrameIndex: 9, Description: "Loop: back to triumphant pose"},
		},
	}
}

// GenerateOpenPoseSkeleton generates an OpenPose skeleton JSON for a given pose
// This is a simplified version - in production, you'd have pre-rendered pose images
func GenerateOpenPoseSkeleton(template PoseTemplate) (string, error) {
	// This would contain actual OpenPose keypoint data
	// For now, return a placeholder that indicates the pose
	skeleton := map[string]interface{}{
		"animation": template.AnimationType,
		"frame":     template.FrameIndex,
		"pose":      template.Description,
		// In production, add actual keypoints:
		// "keypoints": [...] // 18 or 25 body keypoints with x, y, confidence
	}

	// Convert to JSON and base64 encode
	// Simplified for now
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", skeleton))), nil
}
