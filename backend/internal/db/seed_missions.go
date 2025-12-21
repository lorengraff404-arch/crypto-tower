package db

import (
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
)

// SeedMissions creates the initial tutorial and progression missions (Levels 1-30)
func SeedMissions(db *gorm.DB) error {
	missions := []models.Mission{
		// LEVEL 1 - Juramento del Custodio
		{
			Level:         1,
			Name:          "Juramento del Custodio",
			Description:   "Activate Custodian role and enable Simulator",
			Story:         "Activas tu rol de Custodio y se te habilita el 'Simulador'.",
			MissionType:   "tutorial",
			UnlockFeature: "",
			Objectives: `[
				{"type": "battle_waves", "target": 3, "current": 0, "description": "Complete 3 waves in Simulator mode"},
				{"type": "deploy_units", "target": 3, "current": 0, "description": "Deploy 3 units"},
				{"type": "cast_spells", "target": 1, "current": 0, "description": "Cast 1 spell"}
			]`,
			Rewards: `{
				"gtk": 120,
				"items": [{"type": "material", "name": "Basic Kit", "quantity": 15, "rarity": "C"}]
			}`,
			RequiredLevel: 0,
			IsActive:      true,
		},

		// LEVEL 2 - Primer Mazo
		{
			Level:         2,
			Name:          "Primer Mazo",
			Description:   "Build your first operational deck and complete 5 waves",
			Story:         "El Mentor te enseña a construir tu primer mazo operativo.",
			MissionType:   "tutorial",
			UnlockFeature: "",
			Objectives: `[
				{"type": "build_deck", "target": 1, "current": 0, "description": "Build deck of 20 cards (16 units, 4 spells)"},
				{"type": "battle_waves", "target": 5, "current": 0, "description": "Complete 5 waves including mini-boss"}
			]`,
			Rewards: `{
				"gtk": 180,
				"items": [
					{"type": "material", "name": "Common Material", "quantity": 25, "rarity": "C"},
					{"type": "consumable", "name": "Healing Potion", "quantity": 1, "rarity": "C"}
				]
			}`,
			RequiredLevel: 1,
			IsActive:      true,
		},

		// LEVEL 3 - Rastros en el Patio
		{
			Level:         3,
			Name:          "Rastros en el Patio",
			Description:   "First embryo signs appear - complete tracking events",
			Story:         "Aparece la primera señal de embriones en el 'Patio Silvestre'.",
			MissionType:   "tutorial",
			UnlockFeature: "",
			Objectives: `[
				{"type": "tracking_events", "target": 2, "current": 0, "description": "Complete 2 tracking events"},
				{"type": "battle_waves", "target": 5, "current": 0, "description": "Win with at least 2 lanes"}
			]`,
			Rewards: `{
				"gtk": 220,
				"items": [
					{"type": "seed", "name": "Seed", "quantity": 1, "rarity": "C"},
					{"type": "material", "name": "Nature Essence", "quantity": 10, "rarity": "C"}
				]
			}`,
			RequiredLevel: 2,
			IsActive:      true,
		},

		// LEVEL 4 - Incubación Controlada
		{
			Level:         4,
			Name:          "Incubación Controlada",
			Description:   "Unlock Nursery/Incubator and learn hatching with care",
			Story:         "Desbloqueas el Vivero/Incubadora y aprendes el hatch con 'cuidado'.",
			MissionType:   "tutorial",
			UnlockFeature: "",
			Objectives: `[
				{"type": "start_incubation", "target": 1, "current": 0, "description": "Start 1 incubation (seed/egg)"},
				{"type": "complete_care", "target": 1, "current": 0, "description": "Complete 1 care task to improve stats"}
			]`,
			Rewards: `{
				"gtk": 250,
				"items": [
					{"type": "consumable", "name": "Accelerator", "quantity": 1, "rarity": "C"},
					{"type": "material", "name": "Common Material", "quantity": 20, "rarity": "C"}
				]
			}`,
			RequiredLevel: 3,
			IsActive:      true,
		},

		// LEVEL 5 - Licencia de Criador (UNLOCK: Breeding)
		{
			Level:         5,
			Name:          "Licencia de Criador",
			Description:   "Receive permission to breed units - UNLOCK: Breeding",
			Story:         "Recibes permiso oficial para combinar unidades y producir embriones.",
			MissionType:   "tutorial",
			UnlockFeature: "breeding",
			Objectives: `[
				{"type": "battle_waves", "target": 10, "current": 0, "description": "Complete 10 waves (bosses at 5 and 10)"},
				{"type": "breeding", "target": 1, "current": 0, "description": "Perform 1 breeding (2 parents → 1 egg/seed)"}
			]`,
			Rewards: `{
				"gtk": 300,
				"items": [
					{"type": "egg", "name": "Egg", "quantity": 1, "rarity": "C"},
					{"type": "material", "name": "Mutation Material", "quantity": 5, "rarity": "C"}
				]
			}`,
			RequiredLevel: 4,
			IsActive:      true,
		},

		// LEVEL 6 - Fatiga y Rotación
		{
			Level:         6,
			Name:          "Fatiga y Rotación",
			Description:   "Learn team stamina management to avoid infinite grinding",
			Story:         "Se introduce la gestión de resistencia del equipo.",
			MissionType:   "tutorial",
			UnlockFeature: "",
			Objectives: `[
				{"type": "use_recovery", "target": 1, "current": 0, "description": "Use 1 Energy Drink or send unit to recovery"},
				{"type": "battle_waves", "target": 10, "current": 0, "description": "Complete 2 matches of 5 waves with different teams"}
			]`,
			Rewards: `{
				"gtk": 280,
				"items": [
					{"type": "consumable", "name": "Energy Drink", "quantity": 2, "rarity": "C"},
					{"type": "material", "name": "Common Material", "quantity": 25, "rarity": "C"}
				]
			}`,
			RequiredLevel: 5,
			IsActive:      true,
		},

		// LEVEL 7 - Farm de Soporte
		{
			Level:         7,
			Name:          "Farm de Soporte",
			Description:   "Enable Farming Mode for passive income",
			Story:         "Habilitas el 'Farming Mode' como soporte económico y de descanso.",
			MissionType:   "tutorial",
			UnlockFeature: "",
			Objectives: `[
				{"type": "assign_farming", "target": 2, "current": 0, "description": "Assign 2 units to farming (1-4 hours)"},
				{"type": "claim_farm_rewards", "target": 1, "current": 0, "description": "Claim passive rewards"}
			]`,
			Rewards: `{
				"gtk": 200,
				"items": [
					{"type": "material", "name": "Resource Pack", "quantity": 1, "rarity": "C"},
					{"type": "material", "name": "Common Material", "quantity": 20, "rarity": "C"}
				]
			}`,
			RequiredLevel: 6,
			IsActive:      true,
		},

		// LEVEL 8 - Sinergias de Escuadrón
		{
			Level:         8,
			Name:          "Sinergias de Escuadrón",
			Description:   "Learn to win by composition, not spam",
			Story:         "Aprendes a ganar por composición, no por spam.",
			MissionType:   "tutorial",
			UnlockFeature: "",
			Objectives: `[
				{"type": "activate_synergy", "target": 1, "current": 0, "description": "Activate 1 type synergy (2 of same type)"},
				{"type": "no_damage_waves", "target": 10, "current": 0, "description": "Complete 10 waves without losing more than 2 HP"}
			]`,
			Rewards: `{
				"gtk": 320,
				"items": [
					{"type": "rune", "name": "Rune", "quantity": 1, "rarity": "C"},
					{"type": "material", "name": "Elemental Material", "quantity": 10, "rarity": "C"}
				]
			}`,
			RequiredLevel: 7,
			IsActive:      true,
		},

		// LEVEL 9 - Prueba de Resistencia
		{
			Level:         9,
			Name:          "Prueba de Resistencia",
			Description:   "Extended simulation with enemy buffs",
			Story:         "Simulación extendida con buffs enemigos (preparación para islas).",
			MissionType:   "tutorial",
			UnlockFeature: "",
			Objectives: `[
				{"type": "battle_waves", "target": 15, "current": 0, "description": "Complete 15 waves (bosses at 5/10/15)"},
				{"type": "defeat_buffed_enemies", "target": 1, "current": 0, "description": "Defeat enemies with Amped/Bulked/Warded buffs"}
			]`,
			Rewards: `{
				"gtk": 400,
				"items": [
					{"type": "egg", "name": "Egg/Seed", "quantity": 1, "rarity": "B"},
					{"type": "material", "name": "Rare Material", "quantity": 3, "rarity": "B"}
				]
			}`,
			RequiredLevel: 8,
			IsActive:      true,
		},

		// LEVEL 10 - Certificación de Forja (UNLOCK: Crafting)
		{
			Level:         10,
			Name:          "Certificación de Forja",
			Description:   "Access workshop to convert loot - UNLOCK: Crafting",
			Story:         "Accedes al taller para convertir loot en progreso tangible.",
			MissionType:   "tutorial",
			UnlockFeature: "crafting",
			Objectives: `[
				{"type": "obtain_items", "target": 3, "current": 0, "description": "Obtain 3 items of same type in combat"},
				{"type": "craft", "target": 1, "current": 0, "description": "Perform 1 basic craft (3 items + fee → higher rarity)"}
			]`,
			Rewards: `{
				"gtk": 450,
				"items": [
					{"type": "weapon", "name": "Weapon", "quantity": 1, "rarity": "B"},
					{"type": "consumable", "name": "Experience Booster 24h", "quantity": 1, "rarity": "B"}
				]
			}`,
			RequiredLevel: 9,
			IsActive:      true,
		},

		// === PROGRESSION MISSIONS (11-30) ===

		// LEVEL 11 - Rutina de Taller
		{
			Level:         11,
			Name:          "Rutina de Taller",
			Description:   "Master crafting basics with multiple attempts",
			Story:         "Dominas lo básico del crafteo realizando varias creaciones.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "craft", "target": 2, "current": 0, "description": "Craft 2 times (any type)"},
				{"type": "battle_waves", "target": 10, "current": 0, "description": "Complete 1 PvE of 10 waves"}
			]`,
			Rewards: `{
				"gtk": 500,
				"items": [
					{"type": "material", "name": "Crafting Materials", "quantity": 30, "rarity": "B"},
					{"type": "consumable", "name": "Experience Booster 24h", "quantity": 1, "rarity": "B"}
				]
			}`,
			RequiredLevel: 10,
			IsActive:      true,
		},

		// LEVEL 12 - Disciplina del Mazo
		{
			Level:         12,
			Name:          "Disciplina del Mazo",
			Description:   "Win with diverse team compositions",
			Story:         "Aprendes versatilidad ganando con diferentes arquetipos.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "win_archetypes", "target": 3, "current": 0, "description": "Win 3 matches with 3 different archetypes (Fortress/Glass Cannon/Balanced)"}
			]`,
			Rewards: `{
				"gtk": 550,
				"items": [
					{"type": "rune", "name": "Rune", "quantity": 1, "rarity": "B"},
					{"type": "item", "name": "Random Item", "quantity": 1, "rarity": "B"}
				]
			}`,
			RequiredLevel: 11,
			IsActive:      true,
		},

		// LEVEL 13 - Control de Desgaste
		{
			Level:         13,
			Name:          "Control de Desgaste",
			Description:   "Manage character durability effectively",
			Story:         "Aprendes a gestionar el desgaste de tus personajes.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "preserve_durability", "target": 5, "current": 0, "description": "Finish 5 matches without 2 units dropping below critical durability"},
				{"type": "use_item", "target": 1, "current": 0, "description": "Use 1 Healing Potion"}
			]`,
			Rewards: `{
				"gtk": 600,
				"items": [
					{"type": "consumable", "name": "Healing Potion", "quantity": 3, "rarity": "B"},
					{"type": "consumable", "name": "Energy Drink", "quantity": 2, "rarity": "B"}
				]
			}`,
			RequiredLevel: 12,
			IsActive:      true,
		},

		// LEVEL 14 - Preparación de Expedición
		{
			Level:         14,
			Name:          "Preparación de Expedición",
			Description:   "Use farming mode for resource gathering",
			Story:         "Preparas tu expedición mediante farming estratégico.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "assign_farming", "target": 3, "current": 0, "description": "Send 3 units to farming (minimum 4h)"},
				{"type": "claim_farm_rewards", "target": 1, "current": 0, "description": "Claim resources by type"}
			]`,
			Rewards: `{
				"gtk": 650,
				"items": [
					{"type": "material", "name": "Mixed Materials", "quantity": 50, "rarity": "B"}
				]
			}`,
			RequiredLevel: 13,
			IsActive:      true,
		},

		// LEVEL 15 - Carta de Navegación (UNLOCK: Island Raids)
		{
			Level:         15,
			Name:          "Carta de Navegación",
			Description:   "Unlock Island Raids - first expedition to Bosque Raíz",
			Story:         "Recibes acceso a las Incursiones de Islas.",
			MissionType:   "progression",
			UnlockFeature: "island_raids",
			Objectives: `[
				{"type": "island_raid", "target": 1, "current": 0, "description": "Complete 1 island raid (10-20 waves)"}
			]`,
			Rewards: `{
				"gtk": 700,
				"items": [
					{"type": "seed", "name": "Island Egg/Seed", "quantity": 1, "rarity": "B"},
					{"type": "material", "name": "Rare Materials", "quantity": 5, "rarity": "A"}
				]
			}`,
			RequiredLevel: 14,
			IsActive:      true,
		},

		// LEVEL 16 - Incubación Estratégica
		{
			Level:         16,
			Name:          "Incubación Estratégica",
			Description:   "Strategic incubation and marketplace preparation",
			Story:         "Aprendes incubación estratégica y preparación de mercado.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "start_incubation", "target": 2, "current": 0, "description": "Start 2 incubations with care tasks"},
				{"type": "keep_unopened", "target": 1, "current": 0, "description": "Keep 1 egg/seed unopened for market"}
			]`,
			Rewards: `{
				"gtk": 750,
				"items": [
					{"type": "consumable", "name": "Accelerator", "quantity": 2, "rarity": "B"}
				]
			}`,
			RequiredLevel: 15,
			IsActive:      true,
		},

		// LEVEL 17 - Límite y Eficiencia
		{
			Level:         17,
			Name:          "Límite y Eficiencia",
			Description:   "Optimize raid efficiency within daily limits",
			Story:         "Optimizas raids sin exceder el límite diario.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "daily_raids", "target": 3, "current": 0, "description": "Complete 3 raids in one day"}
			]`,
			Rewards: `{
				"gtk": 800,
				"items": [
					{"type": "material", "name": "Premium Materials", "quantity": 10, "rarity": "A"}
				]
			}`,
			RequiredLevel: 16,
			IsActive:      true,
		},

		// LEVEL 18 - Criador Responsable
		{
			Level:         18,
			Name:          "Criador Responsable",
			Description:   "Responsible breeding with cooldowns and evolution",
			Story:         "Practicas crianza responsable con límites y evolución.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "breeding", "target": 2, "current": 0, "description": "2 breedings with cooldown limits"},
				{"type": "evolution", "target": 1, "current": 0, "description": "1 character evolution"}
			]`,
			Rewards: `{
				"gtk": 850,
				"items": [
					{"type": "egg", "name": "Premium Egg/Seed", "quantity": 1, "rarity": "A"},
					{"type": "material", "name": "Mutation Material", "quantity": 10, "rarity": "A"}
				]
			}`,
			RequiredLevel: 17,
			IsActive:      true,
		},

		// LEVEL 19 - Operación Anti-Plaga
		{
			Level:         19,
			Name:          "Operación Anti-Plaga",
			Description:   "Master island raids with perfect execution",
			Story:         "Dominas las raids con ejecución perfecta.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "perfect_raid", "target": 1, "current": 0, "description": "Clear 1 Beginner island without losing HP"},
				{"type": "advanced_raid_trial", "target": 1, "current": 0, "description": "Complete 1 Advanced island trial"}
			]`,
			Rewards: `{
				"gtk": 900,
				"items": [
					{"type": "item", "name": "Guaranteed Item", "quantity": 1, "rarity": "A"}
				]
			}`,
			RequiredLevel: 18,
			IsActive:      true,
		},

		// LEVEL 20 - Acceso a la Arena (UNLOCK: Ranked Battles)
		{
			Level:         20,
			Name:          "Acceso a la Arena",
			Description:   "Unlock Ranked PvP battles - enter competitive scene",
			Story:         "Accedes a batallas clasificatorias competitivas.",
			MissionType:   "progression",
			UnlockFeature: "ranked_pvp",
			Objectives: `[
				{"type": "ranked_matches", "target": 5, "current": 0, "description": "Play 5 ranked matches"},
				{"type": "register_rank", "target": 1, "current": 0, "description": "Register first rank"}
			]`,
			Rewards: `{
				"gtk": 1000,
				"items": [
					{"type": "cosmetic", "name": "Rank Badge", "quantity": 1, "rarity": "A"},
					{"type": "cosmetic", "name": "Basic Cosmetic", "quantity": 1, "rarity": "B"}
				]
			}`,
			RequiredLevel: 19,
			IsActive:      true,
		},

		// LEVEL 21 - Racha Controlada
		{
			Level:         21,
			Name:          "Racha Controlada",
			Description:   "Learn controlled win streaks without excessive risk",
			Story:         "Aprendes a jugar seguro, no siempre all-in.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "win_streak", "target": 3, "current": 0, "description": "Achieve 3-win streak"}
			]`,
			Rewards: `{
				"gtk": 1050,
				"items": [
					{"type": "consumable", "name": "Recovery Consumable", "quantity": 1, "rarity": "A"}
				]
			}`,
			RequiredLevel: 20,
			IsActive:      true,
		},

		// LEVEL 22 - Apuesta Opcional
		{
			Level:         22,
			Name:          "Apuesta Opcional",
			Description:   "Introduction to PvP betting system",
			Story:         "Aprendes el sistema de apuestas PvP (pool, burn, payout).",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "pvp_bet", "target": 1, "current": 0, "description": "Place 1 small PvP bet (TOWER)"}
			]`,
			Rewards: `{
				"gtk": 1100,
				"tower": 10,
				"items": [
					{"type": "cosmetic", "name": "Betting Cosmetic", "quantity": 1, "rarity": "B"}
				]
			}`,
			RequiredLevel: 21,
			IsActive:      true,
		},

		// LEVEL 23 - Gestión de Caps
		{
			Level:         23,
			Name:          "Gestión de Caps",
			Description:   "Learn daily reward caps and optimization",
			Story:         "Dominas la gestión de límites y recompensas diarias.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "first_win_day", "target": 1, "current": 0, "description": "Complete First Win of the Day"},
				{"type": "daily_quests", "target": 3, "current": 0, "description": "Complete 3 daily quests"}
			]`,
			Rewards: `{
				"gtk": 1150,
				"items": [
					{"type": "consumable", "name": "Experience Booster 24h", "quantity": 1, "rarity": "A"}
				]
			}`,
			RequiredLevel: 22,
			IsActive:      true,
		},

		// LEVEL 24 - Arsenal Especializado
		{
			Level:         24,
			Name:          "Arsenal Especializado",
			Description:   "Complete equipment set optimization",
			Story:         "Optimizas tu arsenal con set completo.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "equip_full_set", "target": 1, "current": 0, "description": "Equip full set (weapon/armor/accessory/rune)"},
				{"type": "win_ranked", "target": 2, "current": 0, "description": "Win 2 ranked matches"}
			]`,
			Rewards: `{
				"gtk": 1200,
				"items": [
					{"type": "item", "name": "Guaranteed B Item", "quantity": 1, "rarity": "B"}
				]
			}`,
			RequiredLevel: 23,
			IsActive:      true,
		},

		// LEVEL 25 - Fusión Semanal
		{
			Level:         25,
			Name:          "Fusión Semanal",
			Description:   "Master character fusion mechanic",
			Story:         "Dominas la fusión de personajes (límite semanal).",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "character_fusion", "target": 1, "current": 0, "description": "Execute 1 Character Fusion (burns parents)"}
			]`,
			Rewards: `{
				"gtk": 1250,
				"items": [
					{"type": "character", "name": "Fusion Result", "quantity": 1, "rarity": "A"},
					{"type": "material", "name": "Rare Materials", "quantity": 15, "rarity": "A"}
				]
			}`,
			RequiredLevel: 24,
			IsActive:      true,
		},

		// LEVEL 26 - Mercado Responsable
		{
			Level:         26,
			Name:          "Mercado Responsable",
			Description:   "Learn marketplace mechanics and fees",
			Story:         "Comprendes el mercado (5% fee, sin tienda oficial).",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "list_marketplace", "target": 1, "current": 0, "description": "List 1 egg/seed or item"}
			]`,
			Rewards: `{
				"gtk": 1300,
				"items": [
					{"type": "bonus", "name": "Fee Discount Voucher", "quantity": 1, "rarity": "B"}
				]
			}`,
			RequiredLevel: 25,
			IsActive:      true,
		},

		// LEVEL 27 - Asalto Coordinado
		{
			Level:         27,
			Name:          "Asalto Coordinado",
			Description:   "Weekly raid coordination within daily limits",
			Story:         "Coordinas raids semanales sin exceder 5/día.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "weekly_raids", "target": 5, "current": 0, "description": "Complete 5 raids in one week"}
			]`,
			Rewards: `{
				"gtk": 1350,
				"items": [
					{"type": "material", "name": "Material Pack", "quantity": 1, "rarity": "A"},
					{"type": "item", "name": "Rare Item", "quantity": 1, "rarity": "A"}
				]
			}`,
			RequiredLevel: 26,
			IsActive:      true,
		},

		// LEVEL 28 - Escalón de Rangos
		{
			Level:         28,
			Name:          "Escalón de Rangos",
			Description:   "Climb the ranked ladder",
			Story:         "Asciendes en la escalera clasificatoria.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "reach_rank", "target": 1, "current": 0, "description": "Reach Tier B/A in ladder"},
				{"type": "target_winrate", "target": 1, "current": 0, "description": "Maintain target win rate for week"}
			]`,
			Rewards: `{
				"gtk": 1400,
				"tower": 50,
				"items": [
					{"type": "cosmetic", "name": "Rank Cosmetic", "quantity": 1, "rarity": "A"}
				]
			}`,
			RequiredLevel: 27,
			IsActive:      true,
		},

		// LEVEL 29 - Preparación de Maestría
		{
			Level:         29,
			Name:          "Preparación de Maestría",
			Description:   "High-risk crafting with decreasing success rates",
			Story:         "Practicas crafteo de alto riesgo con gestión.",
			MissionType:   "progression",
			UnlockFeature: "",
			Objectives: `[
				{"type": "consecutive_crafts", "target": 3, "current": 0, "description": "3 consecutive crafts with risk management"}
			]`,
			Rewards: `{
				"gtk": 1450,
				"items": [
					{"type": "material", "name": "Premium Materials", "quantity": 20, "rarity": "A"},
					{"type": "item", "name": "A Item", "quantity": 1, "rarity": "A"}
				]
			}`,
			RequiredLevel: 28,
			IsActive:      true,
		},

		// LEVEL 30 - Maestro de Forja (UNLOCK: Advanced Crafting)
		{
			Level:         30,
			Name:          "Maestro de Forja",
			Description:   "Master craftsman - UNLOCK: Advanced Crafting",
			Story:         "Te conviertes en maestro de forja, desbloqueando crafteo avanzado.",
			MissionType:   "progression",
			UnlockFeature: "advanced_crafting",
			Objectives: `[
				{"type": "advanced_craft", "target": 1, "current": 0, "description": "Complete 1 advanced craft (S/SS path)"},
				{"type": "win_ranked_week", "target": 3, "current": 0, "description": "Win 3 ranked same week"}
			]`,
			Rewards: `{
				"gtk": 1500,
				"tower": 100,
				"items": [
					{"type": "recipe", "name": "Exclusive Recipe", "quantity": 1, "rarity": "S"},
					{"type": "cosmetic", "name": "Master Cosmetic", "quantity": 1, "rarity": "S"},
					{"type": "material", "name": "Legendary Material", "quantity": 5, "rarity": "S"}
				]
			}`,
			RequiredLevel: 29,
			IsActive:      true,
		},
	}

	// Seed missions (upsert behavior - insert if not exists)
	for _, mission := range missions {
		var existing models.Mission
		result := db.Where("level = ?", mission.Level).First(&existing)
		if result.Error != nil {
			// Mission doesn't exist, create it
			if err := db.Create(&mission).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}
