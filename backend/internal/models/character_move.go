package models

import (
	"time"
)

// CharacterMove represents a move/attack that a character can use in battle (Phase 10.2)
type CharacterMove struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Ownership
	CharacterID uint `gorm:"not null;index" json:"character_id"`
	MoveSlot    int  `gorm:"not null" json:"move_slot"` // 1-4

	// Base Move Data
	Name      string `gorm:"size:50;not null" json:"name"`       // "Flamethrower", "Aqua Jet"
	Type      string `gorm:"size:20;not null;index" json:"type"` // "FIRE", "WATER", "GRASS", etc.
	Category  string `gorm:"size:20;not null" json:"category"`   // "PHYSICAL", "SPECIAL", "STATUS"
	Power     int    `gorm:"default:0" json:"power"`             // 80, 100, 120 (0 for status moves)
	Accuracy  int    `gorm:"default:100" json:"accuracy"`        // 100, 95, 90 (%)
	BasePP    int    `gorm:"default:15" json:"base_pp"`          // 15, 20, 25 (max uses)
	CurrentPP int    `gorm:"default:15" json:"current_pp"`       // Decrements each use
	Priority  int    `gorm:"default:0" json:"priority"`          // +1, 0, -1 (for move order)

	// Secondary Effects
	EffectChance    int    `gorm:"default:0" json:"effect_chance"`    // 30% chance to burn
	EffectType      string `gorm:"size:30" json:"effect_type"`        // "burn", "paralyze", "lower_def"
	EffectMagnitude int    `gorm:"default:0" json:"effect_magnitude"` // -1 stage, +2 stages

	// Metadata
	Description string `gorm:"type:text" json:"description"`
	Animation   string `gorm:"size:50" json:"animation"` // "fire_blast", "water_pulse" (for frontend)
}

// MoveTemplate for seeding default moves (not a database table)
type MoveTemplate struct {
	Name            string
	Type            string
	Category        string
	Power           int
	Accuracy        int
	PP              int
	Priority        int
	EffectChance    int
	EffectType      string
	EffectMagnitude int
	Description     string
	Animation       string
}

// Common move templates for seeding
var DefaultMoves = []MoveTemplate{
	// Fire Moves
	{Name: "Ember", Type: "FIRE", Category: "SPECIAL", Power: 40, Accuracy: 100, PP: 25, EffectChance: 10, EffectType: "burn", Description: "The target is scorched with intense fire. May burn."},
	{Name: "Flamethrower", Type: "FIRE", Category: "SPECIAL", Power: 90, Accuracy: 100, PP: 15, EffectChance: 10, EffectType: "burn", Description: "A powerful fire attack. May burn the target."},
	{Name: "Fire Blast", Type: "FIRE", Category: "SPECIAL", Power: 110, Accuracy: 85, PP: 5, EffectChance: 10, EffectType: "burn", Description: "The most powerful Fire attack. May burn."},
	{Name: "Will-O-Wisp", Type: "FIRE", Category: "STATUS", Power: 0, Accuracy: 85, PP: 15, EffectChance: 100, EffectType: "burn", Description: "Inflicts a burn on the target."},

	// Water Moves
	{Name: "Water Gun", Type: "WATER", Category: "SPECIAL", Power: 40, Accuracy: 100, PP: 25, Description: "Squirts water to attack the target."},
	{Name: "Surf", Type: "WATER", Category: "SPECIAL", Power: 90, Accuracy: 100, PP: 15, Description: "A massive wave that hits everything."},
	{Name: "Hydro Pump", Type: "WATER", Category: "SPECIAL", Power: 110, Accuracy: 80, PP: 5, Description: "The ultimate Water attack. Low accuracy."},
	{Name: "Aqua Jet", Type: "WATER", Category: "PHYSICAL", Power: 40, Accuracy: 100, PP: 20, Priority: 1, Description: "A high-priority attack that always strikes first."},

	// Grass Moves
	{Name: "Vine Whip", Type: "GRASS", Category: "PHYSICAL", Power: 45, Accuracy: 100, PP: 25, Description: "Strikes with slender vines."},
	{Name: "Razor Leaf", Type: "GRASS", Category: "PHYSICAL", Power: 55, Accuracy: 95, PP: 25, Description: "Sharp-edged leaves are launched. High crit ratio."},
	{Name: "Solar Beam", Type: "GRASS", Category: "SPECIAL", Power: 120, Accuracy: 100, PP: 10, Description: "Powerful but requires charging."},
	{Name: "Leech Seed", Type: "GRASS", Category: "STATUS", Power: 0, Accuracy: 90, PP: 10, EffectChance: 100, EffectType: "leech", Description: "Drains HP each turn."},

	// Electric Moves
	{Name: "Thunder Shock", Type: "ELECTRIC", Category: "SPECIAL", Power: 40, Accuracy: 100, PP: 30, EffectChance: 10, EffectType: "paralyze", Description: "An electric shock. May paralyze."},
	{Name: "Thunderbolt", Type: "ELECTRIC", Category: "SPECIAL", Power: 90, Accuracy: 100, PP: 15, EffectChance: 10, EffectType: "paralyze", Description: "A strong electric attack. May paralyze."},
	{Name: "Thunder", Type: "ELECTRIC", Category: "SPECIAL", Power: 110, Accuracy: 70, PP: 10, EffectChance: 30, EffectType: "paralyze", Description: "The ultimate Electric attack. May paralyze."},
	{Name: "Thunder Wave", Type: "ELECTRIC", Category: "STATUS", Power: 0, Accuracy: 90, PP: 20, EffectChance: 100, EffectType: "paralyze", Description: "Paralyzes the target."},

	// Normal Moves
	{Name: "Tackle", Type: "NORMAL", Category: "PHYSICAL", Power: 40, Accuracy: 100, PP: 35, Description: "A basic full-body charge attack."},
	{Name: "Quick Attack", Type: "NORMAL", Category: "PHYSICAL", Power: 40, Accuracy: 100, PP: 30, Priority: 1, Description: "A fast attack that strikes first."},
	{Name: "Hyper Beam", Type: "NORMAL", Category: "SPECIAL", Power: 150, Accuracy: 90, PP: 5, Description: "Devastatingly powerful but user must rest next turn."},

	// Status Moves
	{Name: "Swords Dance", Type: "NORMAL", Category: "STATUS", Power: 0, Accuracy: 100, PP: 20, EffectChance: 100, EffectType: "raise_atk_2", Description: "Sharply raises ATK."},
	{Name: "Iron Defense", Type: "STEEL", Category: "STATUS", Power: 0, Accuracy: 100, PP: 15, EffectChance: 100, EffectType: "raise_def_2", Description: "Sharply raises DEF."},
	{Name: "Agility", Type: "PSYCHIC", Category: "STATUS", Power: 0, Accuracy: 100, PP: 30, EffectChance: 100, EffectType: "raise_spd_2", Description: "Sharply raises SPD."},
	{Name: "Growl", Type: "NORMAL", Category: "STATUS", Power: 0, Accuracy: 100, PP: 40, EffectChance: 100, EffectType: "lower_atk_1", Description: "Lowers target's ATK."},
}
