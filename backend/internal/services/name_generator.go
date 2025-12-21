package services

import (
	"fmt"
	"strings"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
)

// NameGeneratorService generates unique character names
type NameGeneratorService struct{}

// NewNameGeneratorService creates a new name generator
func NewNameGeneratorService() *NameGeneratorService {
	return &NameGeneratorService{}
}

// Generate creates a unique name based on character traits
func (s *NameGeneratorService) Generate(element, charType, class, rarity string) string {
	prefix := s.getElementPrefix(element)
	core := s.getTypeCoreName(charType)
	suffix := s.getClassRaritySuffix(class, rarity)

	baseName := fmt.Sprintf("%s %s %s", prefix, core, suffix)

	// Check uniqueness and add numeric suffix if needed
	return s.ensureUniqueness(baseName)
}

// getElementPrefix returns prefix based on element
func (s *NameGeneratorService) getElementPrefix(element string) string {
	prefixes := map[string][]string{
		"Fire":     {"Inferno", "Blaze", "Ember", "Flame", "Pyro", "Scorch"},
		"Water":    {"Aqua", "Hydro", "Tide", "Wave", "Frost", "Marine"},
		"Earth":    {"Terra", "Geo", "Stone", "Boulder", "Quake", "Granite"},
		"Air":      {"Aero", "Gale", "Storm", "Zephyr", "Tempest", "Cyclone"},
		"Light":    {"Lux", "Radiant", "Divine", "Holy", "Celestial", "Lumina"},
		"Dark":     {"Shadow", "Umbra", "Void", "Eclipse", "Nocturnal", "Tenebris"},
		"Electric": {"Volt", "Thunder", "Spark", "Lightning", "Electro", "Plasma"},
		"Ice":      {"Cryo", "Glacial", "Frozen", "Arctic", "Blizzard", "Rime"},
	}

	options := prefixes[element]
	if len(options) == 0 {
		return "Mystic"
	}

	// Use first option for consistency (can be randomized later)
	return options[0]
}

// getTypeCoreName returns core name based on character type
func (s *NameGeneratorService) getTypeCoreName(charType string) string {
	coreNames := map[string][]string{
		"BEAST":   {"Lion", "Tiger", "Wolf", "Bear", "Panther", "Lynx"},
		"DRAGON":  {"Drake", "Wyrm", "Wyvern", "Serpent", "Leviathan", "Basilisk"},
		"BIRD":    {"Phoenix", "Eagle", "Falcon", "Hawk", "Raven", "Condor"},
		"INSECT":  {"Mantis", "Beetle", "Wasp", "Hornet", "Scorpion", "Scarab"},
		"AQUATIC": {"Kraken", "Shark", "Whale", "Dolphin", "Orca", "Leviathan"},
		"MINERAL": {"Golem", "Titan", "Colossus", "Sentinel", "Guardian", "Monolith"},
		"SPIRIT":  {"Phantom", "Wraith", "Specter", "Ghost", "Shade", "Apparition"},
		"AVIAN":   {"Griffin", "Harpy", "Garuda", "Thunderbird", "Roc", "Simurgh"},
		"PLANT":   {"Treant", "Dryad", "Blossom", "Thorn", "Vine", "Sprout"},
		"MACHINE": {"Automaton", "Construct", "Mech", "Golem", "Engine", "Gear"},
	}

	options := coreNames[charType]
	if len(options) == 0 {
		return "Creature"
	}

	return options[0]
}

// getClassRaritySuffix returns suffix based on class and rarity
func (s *NameGeneratorService) getClassRaritySuffix(class, rarity string) string {
	suffixes := map[string]map[string]string{
		"Warrior": {
			"SSS": "the Destroyer",
			"SS":  "the Conqueror",
			"S":   "the Brave",
			"A":   "the Valiant",
			"B":   "the Fighter",
			"C":   "the Recruit",
		},
		"Mage": {
			"SSS": "the Archmage",
			"SS":  "the Sorcerer",
			"S":   "the Wizard",
			"A":   "the Enchanter",
			"B":   "the Apprentice",
			"C":   "the Novice",
		},
		"Archer": {
			"SSS": "the Deadeye",
			"SS":  "the Sharpshooter",
			"S":   "the Marksman",
			"A":   "the Hunter",
			"B":   "the Scout",
			"C":   "the Trainee",
		},
		"Tank": {
			"SSS": "the Immovable",
			"SS":  "the Fortress",
			"S":   "the Defender",
			"A":   "the Guardian",
			"B":   "the Shield",
			"C":   "the Bulwark",
		},
		"Support": {
			"SSS": "the Savior",
			"SS":  "the Healer",
			"S":   "the Cleric",
			"A":   "the Medic",
			"B":   "the Helper",
			"C":   "the Aide",
		},
		"Rogue": {
			"SSS": "the Assassin",
			"SS":  "the Shadow",
			"S":   "the Thief",
			"A":   "the Rogue",
			"B":   "the Sneaky",
			"C":   "the Pickpocket",
		},
		"Paladin": {
			"SSS": "the Crusader",
			"SS":  "the Templar",
			"S":   "the Knight",
			"A":   "the Protector",
			"B":   "the Squire",
			"C":   "the Initiate",
		},
		"Berserker": {
			"SSS": "the Unstoppable",
			"SS":  "the Rampager",
			"S":   "the Furious",
			"A":   "the Wild",
			"B":   "the Fierce",
			"C":   "the Angry",
		},
	}

	classSuffixes := suffixes[class]
	if classSuffixes == nil {
		return "the Unknown"
	}

	suffix := classSuffixes[rarity]
	if suffix == "" {
		return "the Common"
	}

	return suffix
}

// ensureUniqueness checks if name exists and adds numeric suffix if needed
func (s *NameGeneratorService) ensureUniqueness(baseName string) string {
	// Check if base name exists
	var existing models.Character
	err := db.DB.Where("unique_name = ?", baseName).First(&existing).Error

	if err != nil {
		// Name is unique
		return baseName
	}

	// Name exists, find next available number
	counter := 2
	for {
		numberedName := fmt.Sprintf("%s #%d", baseName, counter)
		err := db.DB.Where("unique_name = ?", numberedName).First(&existing).Error

		if err != nil {
			// This numbered version is unique
			return numberedName
		}

		counter++

		// Safety limit
		if counter > 1000 {
			// Fallback to timestamp-based uniqueness
			return fmt.Sprintf("%s #%d", baseName, db.DB.NowFunc().Unix())
		}
	}
}

// GenerateVariation creates a name variation with different prefix/core
func (s *NameGeneratorService) GenerateVariation(element, charType, class, rarity string, variation int) string {
	prefixes := s.getAllPrefixes(element)
	cores := s.getAllCores(charType)

	prefixIndex := variation % len(prefixes)
	coreIndex := (variation / len(prefixes)) % len(cores)

	prefix := prefixes[prefixIndex]
	core := cores[coreIndex]
	suffix := s.getClassRaritySuffix(class, rarity)

	baseName := fmt.Sprintf("%s %s %s", prefix, core, suffix)
	return s.ensureUniqueness(baseName)
}

// getAllPrefixes returns all prefixes for an element
func (s *NameGeneratorService) getAllPrefixes(element string) []string {
	allPrefixes := map[string][]string{
		"Fire":     {"Inferno", "Blaze", "Ember", "Flame", "Pyro", "Scorch"},
		"Water":    {"Aqua", "Hydro", "Tide", "Wave", "Frost", "Marine"},
		"Earth":    {"Terra", "Geo", "Stone", "Boulder", "Quake", "Granite"},
		"Air":      {"Aero", "Gale", "Storm", "Zephyr", "Tempest", "Cyclone"},
		"Light":    {"Lux", "Radiant", "Divine", "Holy", "Celestial", "Lumina"},
		"Dark":     {"Shadow", "Umbra", "Void", "Eclipse", "Nocturnal", "Tenebris"},
		"Electric": {"Volt", "Thunder", "Spark", "Lightning", "Electro", "Plasma"},
		"Ice":      {"Cryo", "Glacial", "Frozen", "Arctic", "Blizzard", "Rime"},
	}

	options := allPrefixes[element]
	if len(options) == 0 {
		return []string{"Mystic"}
	}
	return options
}

// getAllCores returns all core names for a type
func (s *NameGeneratorService) getAllCores(charType string) []string {
	allCores := map[string][]string{
		"BEAST":   {"Lion", "Tiger", "Wolf", "Bear", "Panther", "Lynx"},
		"DRAGON":  {"Drake", "Wyrm", "Wyvern", "Serpent", "Leviathan", "Basilisk"},
		"BIRD":    {"Phoenix", "Eagle", "Falcon", "Hawk", "Raven", "Condor"},
		"INSECT":  {"Mantis", "Beetle", "Wasp", "Hornet", "Scorpion", "Scarab"},
		"AQUATIC": {"Kraken", "Shark", "Whale", "Dolphin", "Orca", "Leviathan"},
		"MINERAL": {"Golem", "Titan", "Colossus", "Sentinel", "Guardian", "Monolith"},
		"SPIRIT":  {"Phantom", "Wraith", "Specter", "Ghost", "Shade", "Apparition"},
		"AVIAN":   {"Griffin", "Harpy", "Garuda", "Thunderbird", "Roc", "Simurgh"},
		"PLANT":   {"Treant", "Dryad", "Blossom", "Thorn", "Vine", "Sprout"},
		"MACHINE": {"Automaton", "Construct", "Mech", "Golem", "Engine", "Gear"},
	}

	options := allCores[charType]
	if len(options) == 0 {
		return []string{"Creature"}
	}
	return options
}

// ParseName extracts components from a generated name
func (s *NameGeneratorService) ParseName(name string) (prefix, core, suffix string) {
	parts := strings.Split(name, " ")
	if len(parts) >= 3 {
		prefix = parts[0]
		core = parts[1]
		suffix = strings.Join(parts[2:], " ")
	}
	return
}
