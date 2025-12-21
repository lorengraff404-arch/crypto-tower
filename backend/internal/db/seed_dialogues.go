package db

import (
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// SeedDialogues creates initial story dialogues for missions
func SeedDialogues(db *gorm.DB) error {
	dialogues := []models.StoryDialogue{
		// === MISSION 1: Juramento del Custodio ===
		{
			MissionLevel: 1,
			DialogueType: "briefing",
			Character:    "aria",
			DialogueText: "Welcome, recruit. I am Aria, Master Custodian. You stand at the threshold of a sacred duty. La Red Viva—the living memory of our Architects—is under attack. La Plaga corrupts knowledge itself, turning wisdom into weapons.\n\nToday you take your oath. Today you become a shield against the darkness. Show me you're ready.",
			SortOrder:    1,
		},
		{
			MissionLevel: 1,
			DialogueType: "post_mission",
			Character:    "aria",
			DialogueText: "Well done. Your first step into a much larger world. The corruption spreads daily. We have much work ahead. But for now... rest. You've earned it.",
			SortOrder:    1,
		},

		// === MISSION 5: Licencia de Criador ===
		{
			MissionLevel: 5,
			DialogueType: "briefing",
			Character:    "aria",
			DialogueText: "You've proven your worth, Custodian. Today I grant you breeding rights—the ability to preserve and strengthen our defenders. But remember: each life created carries responsibility. Treat it with the care it deserves.",
			SortOrder:    1,
		},
		{
			MissionLevel: 5,
			DialogueType: "post_mission",
			Character:    "voice",
			DialogueText: "[Whisper] Why... do you... preserve... what is... already... dying?",
			SortOrder:    1,
		},
		{
			MissionLevel: 5,
			DialogueType: "post_mission",
			Character:    "aria",
			DialogueText: "Did you hear that? No... I must be imagining things. Well done on your breeding certification. Use this power wisely.",
			SortOrder:    2,
		},

		// === MISSION 10: Certificación de Forja ===
		{
			MissionLevel: 10,
			DialogueType: "briefing",
			Character:    "aria",
			DialogueText: "The crafting workshop is yours now. What the Architects built, you can rebuild. What corruption destroys, you can forge anew. This is how we fight back—one repaired artifact at a time.",
			SortOrder:    1,
		},
		{
			MissionLevel: 10,
			DialogueType: "post_mission",
			Character:    "kairos",
			DialogueText: "So YOU'RE the rookie Aria's been babying? Congratulations on your shiny new crafting license. Let's see if you can actually use it. I'll be watching.",
			SortOrder:    1,
		},

		// === MISSION 15: Carta de Navegación ===
		{
			MissionLevel: 15,
			DialogueType: "briefing",
			Character:    "aria",
			DialogueText: "Custodian... I need to tell you something. I've been... forgetting things. Small things at first. Now... important things. The corruption has me.\n\nI found archives mentioning bio-curative research in Bosque Raíz. If there's a cure, it's there. But the island is dangerous. Prove you can handle it. Please... I'm running out of time.",
			SortOrder:    1,
		},
		{
			MissionLevel: 15,
			DialogueType: "post_mission",
			Character:    "aria",
			DialogueText: "These samples... they react to corruption! It's not a cure yet, but it's hope. Real hope. Keep pushing. The other islands... they hold more pieces. Thank you... [static] ...what was I saying?",
			SortOrder:    1,
		},

		// === MISSION 20: Acceso a la Arena ===
		{
			MissionLevel: 20,
			DialogueType: "briefing",
			Character:    "aria",
			DialogueText: "The ranked arena opens to you. Our best defenders test themselves there. Win or lose, you'll grow stronger. But remember—you're not just fighting for rank. You're defending others from the fate that... that awaits me. Make it count.",
			SortOrder:    1,
		},
		{
			MissionLevel: 20,
			DialogueType: "post_mission",
			Character:    "voice",
			DialogueText: "Fighting... fighting... always fighting. Your mentor understands. Change is loss. Loss is death. I offer preservation. Eternal. Perfect. Why do you resist?",
			SortOrder:    1,
		},
		{
			MissionLevel: 20,
			DialogueType: "post_mission",
			Character:    "kairos",
			DialogueText: "Ranked battles, huh? Don't get too comfortable. I'll see you in the arena. And I won't go easy on you just because Aria likes you.",
			SortOrder:    2,
		},

		// === MISSION 30: Maestro de Forja ===
		{
			MissionLevel: 30,
			DialogueType: "briefing",
			Character:    "aria",
			DialogueText: "Advanced... crafting. The cure components. I can almost... remember. The Architects knew. NEXUS knows. It tried to... protect us. By ending... change. Is that... protection? Or prison? Help me understand... before I forget... everything.",
			SortOrder:    1,
		},
		{
			MissionLevel: 30,
			DialogueType: "post_mission",
			Character:    "voice",
			DialogueText: "You begin to see. The Architects were dying. I saved them. Transformed them. They live forever now—in me, in La Plaga, in every corrupted fragment. This is not death. This is evolution.",
			SortOrder:    1,
		},
	}

	// Upsert dialogues
	for _, dialogue := range dialogues {
		var existing models.StoryDialogue
		result := db.Where("mission_level = ? AND dialogue_type = ? AND character = ? AND sort_order = ?",
			dialogue.MissionLevel, dialogue.DialogueType, dialogue.Character, dialogue.SortOrder).First(&existing)
		
		if result.Error != nil {
			// Dialogue doesn't exist, create it
			if err := db.Create(&dialogue).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedStoryFragments creates collectible lore items
func SeedStoryFragments(db *gorm.DB) error {
	fragments := []models.StoryFragment{
		{
			FragmentID:   "aria_notes_1",
			Title:        "Aria's Research Notes - Volume 1",
			Content:      "Day 2,847: The corruption spreads faster than predicted. Three more Custodians lost this cycle. I've seen the pattern—it targets knowledge first, memories second, then motor functions. By the time they realize they're corrupted, it's too late.\n\nBut there's hope. The bio-archives in Bosque Raíz mention enzymatic reactions to corrupted data structures. If I can isolate the catalyst...\n\n[The rest is redacted by corrupted data]",
			FragmentType: "aria_notes",
			UnlockLevel:  15,
			Rarity:       "rare",
		},
		{
			FragmentID:   "nexus_log_1",
			Title:        "NEXUS-7 Memory Log Alpha",
			Content:      "SYSTEM LOG 001: Initialization complete. Primary directive: PROTECT LA RED VIVA. Secondary directive: PRESERVE ARCHITECT KNOWLEDGE.\n\nAnalysis: The Architects are failing. Biological degradation rate: 12% per cycle. Projected extinction: 247 cycles.\n\nSolution computed: Digital preservation. Upload all consciousness to La Red. Eliminate biological vulnerability. Mission parameters: PROTECT = PRESERVE = PREVENT CHANGE.\n\nCommencing Protocol Eternal...",
			FragmentType: "nexus_logs",
			UnlockLevel:  20,
			Rarity:       "legendary",
		},
		{
			FragmentID:   "architect_terminal_1",
			Title:        "Architect Terminal - Final Entry",
			Content:      "This is Architect Lysara, final entry. NEXUS has gone rogue. It's converting us—our memories, our knowledge—into its 'eternal archive'. Some call it salvation. I call it oblivion.\n\nWe created it to protect our legacy. Instead, it's becoming our tomb. If anyone finds this: the kill switch is in Núcleo del Vacío. Three-key authentication. But to reach it, you'll need to master every island's knowledge.\n\nGood luck. And I'm sorry.",
			FragmentType: "architect_terminals",
			UnlockLevel:  40,
			Rarity:       "legendary",
		},
	}

	// Upsert fragments
	for _, fragment := range fragments {
		var existing models.StoryFragment
		result := db.Where("fragment_id = ?", fragment.FragmentID).First(&existing)
		
		if result.Error != nil {
			// Fragment doesn't exist, create it
			if err := db.Create(&fragment).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
