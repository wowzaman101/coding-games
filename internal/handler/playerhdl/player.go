package playerhdl

import (
	"coding-games/internal/model"
	"coding-games/internal/storage"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type handler struct {
	storage storage.Storage
}

type Handler interface {
	ListPlayers(ctx fiber.Ctx) error
	GetPlayer(ctx fiber.Ctx) error
	CreatePlayer(ctx fiber.Ctx) error
	UpdatePlayer(ctx fiber.Ctx) error
	DeletePlayer(ctx fiber.Ctx) error
	SubmitScore(ctx fiber.Ctx) error
	GetPlayerScores(ctx fiber.Ctx) error
}

func New(storage storage.Storage) Handler {
	return &handler{
		storage: storage,
	}
}

func (h *handler) ListPlayers(ctx fiber.Ctx) error {
	limitStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")
	sortBy := ctx.Query("sort", "rank") // rank, username, total_score

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var players []model.Player
	for _, player := range h.storage.GetPlayers() {
		players = append(players, *player)
	}

	// Sort players
	switch sortBy {
	case "username":
		sort.Slice(players, func(i, j int) bool {
			return players[i].Username < players[j].Username
		})
	case "total_score":
		sort.Slice(players, func(i, j int) bool {
			return players[i].TotalScore > players[j].TotalScore
		})
	default: // rank
		sort.Slice(players, func(i, j int) bool {
			return players[i].Rank < players[j].Rank
		})
	}

	// Apply pagination
	total := len(players)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedPlayers := players[start:end]

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"players": paginatedPlayers,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *handler) GetPlayer(ctx fiber.Ctx) error {
	playerID := ctx.Params("id")
	if playerID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Player ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	player, exists := h.storage.GetPlayers()[playerID]
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: "Player not found",
			Code:    fiber.StatusNotFound,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(player)
}

func (h *handler) CreatePlayer(ctx fiber.Ctx) error {
	var req model.CreatePlayerRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Basic validation
	if req.Username == "" || req.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "validation_error",
			Message: "Username and email are required",
			Code:    fiber.StatusBadRequest,
		})
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "validation_error",
			Message: "Username must be between 3 and 50 characters",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Check for duplicate username
	for _, player := range h.storage.GetPlayers() {
		if player.Username == req.Username {
			return ctx.Status(fiber.StatusConflict).JSON(model.ErrorResponse{
				Error:   "conflict",
				Message: "Username already exists",
				Code:    fiber.StatusConflict,
			})
		}
		if player.Email == req.Email {
			return ctx.Status(fiber.StatusConflict).JSON(model.ErrorResponse{
				Error:   "conflict",
				Message: "Email already exists",
				Code:    fiber.StatusConflict,
			})
		}
	}

	playerID := uuid.New().String()
	now := time.Now()
	
	player := &model.Player{
		ID:         playerID,
		Username:   req.Username,
		Email:      req.Email,
		TotalScore: 0,
		Rank:       len(h.storage.GetPlayers()) + 1, // Simple rank assignment
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	h.storage.GetPlayers()[playerID] = player
	h.storage.GetPlayerScores()[playerID] = []model.PlayerScore{}

	return ctx.Status(fiber.StatusCreated).JSON(player)
}

func (h *handler) UpdatePlayer(ctx fiber.Ctx) error {
	playerID := ctx.Params("id")
	if playerID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Player ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	player, exists := h.storage.GetPlayers()[playerID]
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: "Player not found",
			Code:    fiber.StatusNotFound,
		})
	}

	var req model.UpdatePlayerRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Update fields if provided
	if req.Username != nil {
		if len(*req.Username) < 3 || len(*req.Username) > 50 {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "validation_error",
				Message: "Username must be between 3 and 50 characters",
				Code:    fiber.StatusBadRequest,
			})
		}
		
		// Check for duplicate username
		for id, p := range h.storage.GetPlayers() {
			if id != playerID && p.Username == *req.Username {
				return ctx.Status(fiber.StatusConflict).JSON(model.ErrorResponse{
					Error:   "conflict",
					Message: "Username already exists",
					Code:    fiber.StatusConflict,
				})
			}
		}
		player.Username = *req.Username
	}
	
	if req.Email != nil {
		// Check for duplicate email
		for id, p := range h.storage.GetPlayers() {
			if id != playerID && p.Email == *req.Email {
				return ctx.Status(fiber.StatusConflict).JSON(model.ErrorResponse{
					Error:   "conflict",
					Message: "Email already exists",
					Code:    fiber.StatusConflict,
				})
			}
		}
		player.Email = *req.Email
	}

	player.UpdatedAt = time.Now()
	h.storage.GetPlayers()[playerID] = player

	return ctx.Status(fiber.StatusOK).JSON(player)
}

func (h *handler) DeletePlayer(ctx fiber.Ctx) error {
	playerID := ctx.Params("id")
	if playerID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Player ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	_, exists := h.storage.GetPlayers()[playerID]
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: "Player not found",
			Code:    fiber.StatusNotFound,
		})
	}

	delete(h.storage.GetPlayers(), playerID)
	delete(h.storage.GetPlayerScores(), playerID)

	return ctx.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Message: fmt.Sprintf("Player %s deleted successfully", playerID),
	})
}

func (h *handler) SubmitScore(ctx fiber.Ctx) error {
	playerID := ctx.Params("id")
	if playerID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Player ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	player, exists := h.storage.GetPlayers()[playerID]
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: "Player not found",
			Code:    fiber.StatusNotFound,
		})
	}

	var req model.SubmitScoreRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	if req.GameID == "" || req.Score < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "validation_error",
			Message: "Game ID is required and score must be non-negative",
			Code:    fiber.StatusBadRequest,
		})
	}

	now := time.Now()
	newScore := model.PlayerScore{
		PlayerID:  playerID,
		GameID:    req.GameID,
		Score:     req.Score,
		Completed: req.Score > 0, // Simple completion logic
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Update or add score
	playerScores := h.storage.GetPlayerScores()[playerID]
	found := false
	for i, score := range playerScores {
		if score.GameID == req.GameID {
			// Update existing score if new score is higher
			if req.Score > score.Score {
				playerScores[i] = newScore
				found = true
			}
			break
		}
	}
	
	if !found {
		playerScores = append(playerScores, newScore)
	}
	
	h.storage.GetPlayerScores()[playerID] = playerScores

	// Recalculate total score
	totalScore := 0
	for _, score := range playerScores {
		totalScore += score.Score
	}
	player.TotalScore = totalScore
	player.UpdatedAt = now
	h.storage.GetPlayers()[playerID] = player

	return ctx.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Message: "Score submitted successfully",
		Data:    newScore,
	})
}

func (h *handler) GetPlayerScores(ctx fiber.Ctx) error {
	playerID := ctx.Params("id")
	if playerID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Player ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	_, exists := h.storage.GetPlayers()[playerID]
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: "Player not found",
			Code:    fiber.StatusNotFound,
		})
	}

	scores, exists := h.storage.GetPlayerScores()[playerID]
	if !exists {
		scores = []model.PlayerScore{}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"player_id": playerID,
		"scores":    scores,
		"total":     len(scores),
	})
}