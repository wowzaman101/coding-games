package gamehdl

import (
	"coding-games/internal/model"
	"coding-games/internal/storage"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type handler struct {
	storage storage.Storage
}

type Handler interface {
	Test(ctx fiber.Ctx) error
	ListGames(ctx fiber.Ctx) error
	GetGame(ctx fiber.Ctx) error
	CreateGame(ctx fiber.Ctx) error
	UpdateGame(ctx fiber.Ctx) error
	DeleteGame(ctx fiber.Ctx) error
}

func New(storage storage.Storage) Handler {
	return &handler{
		storage: storage,
	}
}

func (h *handler) Test(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "Test endpoint hit",
		"time":    time.Now(),
	})
}

func (h *handler) ListGames(ctx fiber.Ctx) error {
	// Get query parameters for filtering
	difficulty := ctx.Query("difficulty")
	language := ctx.Query("language")
	limitStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var games []model.Game
	for _, game := range h.storage.GetGames() {
		// Apply filters
		if difficulty != "" && game.Difficulty != difficulty {
			continue
		}
		if language != "" && game.Language != language {
			continue
		}
		games = append(games, *game)
	}

	// Apply pagination
	total := len(games)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedGames := games[start:end]

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"games":  paginatedGames,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *handler) GetGame(ctx fiber.Ctx) error {
	gameID := ctx.Params("id")
	if gameID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Game ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	game, exists := h.storage.GetGames()[gameID]
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: "Game not found",
			Code:    fiber.StatusNotFound,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(game)
}

func (h *handler) CreateGame(ctx fiber.Ctx) error {
	var req model.CreateGameRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Basic validation
	if req.Title == "" || req.Description == "" || req.Difficulty == "" || req.Language == "" || req.MaxScore <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "validation_error",
			Message: "All fields are required and max_score must be positive",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Validate difficulty
	if req.Difficulty != "easy" && req.Difficulty != "medium" && req.Difficulty != "hard" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "validation_error",
			Message: "Difficulty must be one of: easy, medium, hard",
			Code:    fiber.StatusBadRequest,
		})
	}

	gameID := uuid.New().String()
	now := time.Now()
	
	game := &model.Game{
		ID:          gameID,
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Language:    req.Language,
		MaxScore:    req.MaxScore,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	h.storage.GetGames()[gameID] = game

	return ctx.Status(fiber.StatusCreated).JSON(game)
}

func (h *handler) UpdateGame(ctx fiber.Ctx) error {
	gameID := ctx.Params("id")
	if gameID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Game ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	game, exists := h.storage.GetGames()[gameID]
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: "Game not found",
			Code:    fiber.StatusNotFound,
		})
	}

	var req model.UpdateGameRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Update fields if provided
	if req.Title != nil {
		game.Title = *req.Title
	}
	if req.Description != nil {
		game.Description = *req.Description
	}
	if req.Difficulty != nil {
		if *req.Difficulty != "easy" && *req.Difficulty != "medium" && *req.Difficulty != "hard" {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "validation_error",
				Message: "Difficulty must be one of: easy, medium, hard",
				Code:    fiber.StatusBadRequest,
			})
		}
		game.Difficulty = *req.Difficulty
	}
	if req.Language != nil {
		game.Language = *req.Language
	}
	if req.MaxScore != nil {
		if *req.MaxScore <= 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "validation_error",
				Message: "max_score must be positive",
				Code:    fiber.StatusBadRequest,
			})
		}
		game.MaxScore = *req.MaxScore
	}

	game.UpdatedAt = time.Now()
	h.storage.GetGames()[gameID] = game

	return ctx.Status(fiber.StatusOK).JSON(game)
}

func (h *handler) DeleteGame(ctx fiber.Ctx) error {
	gameID := ctx.Params("id")
	if gameID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Game ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	_, exists := h.storage.GetGames()[gameID]
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: "Game not found",
			Code:    fiber.StatusNotFound,
		})
	}

	delete(h.storage.GetGames(), gameID)

	return ctx.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Message: fmt.Sprintf("Game %s deleted successfully", gameID),
	})
}
