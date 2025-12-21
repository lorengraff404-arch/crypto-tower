package services

// Maná Base y Crecimiento por Rank
const (
	// Maná Base por Rank (según plan maestro)
	ManaBaseC   = 50
	ManaBaseB   = 70
	ManaBaseA   = 100
	ManaBaseS   = 140
	ManaBaseSS  = 190
	ManaBaseSSS = 250

	// Crecimiento de Maná por Nivel
	ManaGrowthC   = 2
	ManaGrowthB   = 4
	ManaGrowthA   = 6
	ManaGrowthS   = 9
	ManaGrowthSS  = 13
	ManaGrowthSSS = 18
)

// ManaService handles mana calculation logic
type ManaService struct{}

// NewManaService creates a new mana service
func NewManaService() *ManaService {
	return &ManaService{}
}

// CalculateMaxMana calcula el maná máximo basado en rank y nivel
func (s *ManaService) CalculateMaxMana(rank string, level int) int {
	baseMap := map[string]int{
		"C": ManaBaseC, "B": ManaBaseB, "A": ManaBaseA,
		"S": ManaBaseS, "SS": ManaBaseSS, "SSS": ManaBaseSSS,
	}
	growthMap := map[string]int{
		"C": ManaGrowthC, "B": ManaGrowthB, "A": ManaGrowthA,
		"S": ManaGrowthS, "SS": ManaGrowthSS, "SSS": ManaGrowthSSS,
	}

	base := baseMap[rank]
	if base == 0 {
		base = ManaBaseC // Default C
	}

	growth := growthMap[rank]
	if growth == 0 {
		growth = ManaGrowthC // Default C
	}

	return base + (growth * (level - 1))
}

// CalculateMaxHP calcula el HP máximo con multiplicadores de rank y nivel
func (s *ManaService) CalculateMaxHP(baseHP int, rank string, level int) int {
	rankMultiplier := map[string]float64{
		"C": 1.0, "B": 1.2, "A": 1.5,
		"S": 2.0, "SS": 2.8, "SSS": 4.0,
	}

	multiplier := rankMultiplier[rank]
	if multiplier == 0 {
		multiplier = 1.0 // Default C
	}

	levelBonus := float64(level-1) * 0.1 // 10% por nivel

	return int(float64(baseHP) * multiplier * (1 + levelBonus))
}

// GetManaRegenRate retorna regeneración de maná por turno basado en rank
func (s *ManaService) GetManaRegenRate(rank string, level int) int {
	baseRegen := map[string]int{
		"C": 5, "B": 7, "A": 10,
		"S": 15, "SS": 20, "SSS": 30,
	}

	base := baseRegen[rank]
	if base == 0 {
		base = 5 // Default C
	}

	// Pequeño bonus por nivel (0.5 por nivel)
	levelBonus := (level - 1) / 2

	return base + levelBonus
}
