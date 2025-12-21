package main

import (
	"log"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"github.com/lorengraff/crypto-tower-defense/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed config: %v", err)
	}

	if err := db.Connect(cfg); err != nil {
		log.Fatalf("Failed db: %v", err)
	}

	log.Println("🔄 Seeding Moves for existing Characters...")

	// Get all characters
	var characters []models.Character
	if err := db.DB.Find(&characters).Error; err != nil {
		log.Fatalf("Failed to fetch characters: %v", err)
	}

	log.Printf("Found %d characters", len(characters))

	for _, char := range characters {
		// Check if character already has moves
		var existingMoves []models.CharacterMove
		db.DB.Where("character_id = ?", char.ID).Find(&existingMoves)

		if len(existingMoves) > 0 {
			log.Printf("Character %s already has %d moves, skipping", char.Name, len(existingMoves))
			continue
		}

		// Assign 4 moves based on character's element
		moves := getMovesForElement(char.Element, char.ID)

		for i, moveTemplate := range moves {
			move := models.CharacterMove{
				CharacterID:     char.ID,
				MoveSlot:        i + 1,
				Name:            moveTemplate.Name,
				Type:            moveTemplate.Type,
				Category:        moveTemplate.Category,
				Power:           moveTemplate.Power,
				Accuracy:        moveTemplate.Accuracy,
				BasePP:          moveTemplate.PP,
				CurrentPP:       moveTemplate.PP,
				Priority:        moveTemplate.Priority,
				EffectChance:    moveTemplate.EffectChance,
				EffectType:      moveTemplate.EffectType,
				EffectMagnitude: moveTemplate.EffectMagnitude,
				Description:     moveTemplate.Description,
				Animation:       moveTemplate.Animation,
			}

			if err := db.DB.Create(&move).Error; err != nil {
				log.Printf("Failed to create move for %s: %v", char.Name, err)
			} else {
				log.Printf("✅ %s learned %s", char.Name, move.Name)
			}
		}
	}

	log.Println("✅ Move seeding complete!")
}

// getMovesForElement returns 4 moves for a character based on their element
func getMovesForElement(element string, charID uint) []models.MoveTemplate {
	switch element {
	case "FIRE":
		return []models.MoveTemplate{
			models.DefaultMoves[0],  // Ember
			models.DefaultMoves[1],  // Flamethrower
			models.DefaultMoves[3],  // Will-O-Wisp (status)
			models.DefaultMoves[17], // Swords Dance (buff)
		}
	case "WATER":
		return []models.MoveTemplate{
			models.DefaultMoves[4],  // Water Gun
			models.DefaultMoves[5],  // Surf
			models.DefaultMoves[7],  // Aqua Jet (priority)
			models.DefaultMoves[18], // Iron Defense (buff)
		}
	case "GRASS":
		return []models.MoveTemplate{
			models.DefaultMoves[8],  // Vine Whip
			models.DefaultMoves[9],  // Razor Leaf
			models.DefaultMoves[10], // Solar Beam
			models.DefaultMoves[11], // Leech Seed (status)
		}
	case "ELECTRIC":
		return []models.MoveTemplate{
			models.DefaultMoves[12], // Thunder Shock
			models.DefaultMoves[13], // Thunderbolt
			models.DefaultMoves[15], // Thunder Wave (status)
			models.DefaultMoves[19], // Agility (buff)
		}
	case "ICE":
		return []models.MoveTemplate{
			models.DefaultMoves[4],  // Water Gun (ice types can use water)
			models.DefaultMoves[14], // Tackle
			models.DefaultMoves[15], // Quick Attack
			models.DefaultMoves[18], // Iron Defense
		}
	case "DRAGON":
		return []models.MoveTemplate{
			models.DefaultMoves[2],  // Fire Blast
			models.DefaultMoves[6],  // Hydro Pump
			models.DefaultMoves[16], // Hyper Beam
			models.DefaultMoves[17], // Swords Dance
		}
	default: // NORMAL or others
		return []models.MoveTemplate{
			models.DefaultMoves[14], // Tackle
			models.DefaultMoves[15], // Quick Attack
			models.DefaultMoves[17], // Swords Dance
			models.DefaultMoves[20], // Growl
		}
	}
}
