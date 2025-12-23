package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorengraff/crypto-tower-defense/internal/services"
)

// BattleHandler handles battle HTTP requests
type BattleHandler struct {
	battleService *services.BattleService
}

// NewBattleHandler creates a new battle handler
func NewBattleHandler() *BattleHandler {
	return &BattleHandler{
		battleService: services.NewBattleService(),
	}
}

// FindMatch finds an opponent for PvP
// POST /api/v1/battles/matchmaking
func (h *BattleHandler) FindMatch(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		BetAmount int64 `json:"bet_amount" binding:"required,min=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	opponent, err := h.battleService.FindMatch(userID.(uint), req.BetAmount)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"opponent": gin.H{
			"id":         opponent.ID,
			"level":      opponent.Level,
			"rank":       opponent.Rank,
			"rank_tier":  opponent.RankTier,
			"pvp_wins":   opponent.PvPWins,
			"pvp_losses": opponent.PvPLosses,
			"win_streak": opponent.CurrentWinStreak,
		},
		"message": "Opponent found! Ready to battle.",
	})
}

// CreateBattle creates a new PvP battle
// POST /api/v1/battles
func (h *BattleHandler) CreateBattle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		OpponentID uint  `json:"opponent_id" binding:"required"`
		BetAmount  int64 `json:"bet_amount" binding:"required,min=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	battle, err := h.battleService.CreatePvPBattle(userID.(uint), req.OpponentID, req.BetAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"battle":  battle,
		"message": "Battle created! Both players' bets have been escrowed.",
	})
}

// StartBattle starts a pending battle
// POST /api/v1/battles/:id/start
func (h *BattleHandler) StartBattle(c *gin.Context) {
	battleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid battle ID"})
		return
	}

	if _, err := h.battleService.StartBattle(uint(battleID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Battle started!",
	})
}

// CompleteBattle completes a battle and declares winner
// POST /api/v1/battles/:id/complete
func (h *BattleHandler) CompleteBattle(c *gin.Context) {
	battleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid battle ID"})
		return
	}

	var req struct {
		WinnerID   uint   `json:"winner_id" binding:"required"`
		ReplayData string `json:"replay_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.battleService.CompleteBattle(uint(battleID), req.WinnerID, req.ReplayData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Battle completed! Rewards distributed.",
	})
}

// SurrenderBattle handles player surrender
// POST /api/v1/battles/:id/surrender
func (h *BattleHandler) SurrenderBattle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	battleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid battle ID"})
		return
	}

	if err := h.battleService.SurrenderBattle(uint(battleID), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "You have surrendered. Opponent wins.",
	})
}

// GetBattle retrieves battle details
// GET /api/v1/battles/:id
func (h *BattleHandler) GetBattle(c *gin.Context) {
	battleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid battle ID"})
		return
	}

	battle, err := h.battleService.GetBattleByID(uint(battleID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Battle not found"})
		return
	}

	c.JSON(http.StatusOK, battle)
}

// GetBattleHistory retrieves user's battle history
// GET /api/v1/battles/history
func (h *BattleHandler) GetBattleHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	battles, err := h.battleService.GetBattleHistory(userID.(uint), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"battles": battles,
		"count":   len(battles),
	})
}

// RequestRematch creates rematch request
// POST /api/v1/battles/:id/rematch
func (h *BattleHandler) RequestRematch(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	battleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid battle ID"})
		return
	}

	if _, err := h.battleService.RequestRematch(uint(battleID), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rematch requested! Waiting for opponent's response.",
	})
}

// ProcessTurn processes a single turn in battle (SKILL SYSTEM INTEGRATION)
// POST /api/v1/battles/:id/turn
func (h *BattleHandler) ProcessTurn(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	battleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid battle ID"})
		return
	}

	var req struct {
		CharacterID uint   `json:"character_id" binding:"required"`
		Action      string `json:"action" binding:"required"` // "skill" or "attack"
		SkillID     *uint  `json:"skill_id"`                  // Required if action == "skill"
		TargetID    uint   `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Verify ownership is handled by Service
	_ = userID

	// Construct Map for BattleService (which expects map[string]interface{})
	actionData := map[string]interface{}{
		"action":       req.Action,
		"character_id": float64(req.CharacterID), // Service casts to float64 then uint
		"target_id":    float64(req.TargetID),
	}
	if req.SkillID != nil {
		actionData["skill_id"] = float64(*req.SkillID)
	}

	result, err := h.battleService.ProcessTurn(uint(battleID), userID.(uint), actionData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}
