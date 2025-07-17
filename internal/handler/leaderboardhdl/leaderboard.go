package leaderboardhdl

import (
	"coding-games/internal/model"
	"coding-games/internal/storage"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type handler struct {
	storage storage.Storage
}

type Handler interface {
	GetGlobalLeaderboard(ctx fiber.Ctx) error
	GetGameLeaderboard(ctx fiber.Ctx) error
}

func New(storage storage.Storage) Handler {
	return &handler{
		storage: storage,
	}
}

func (h *handler) GetGlobalLeaderboard(ctx fiber.Ctx) error {
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

	var entries []model.LeaderboardEntry
	for _, player := range h.storage.GetPlayers() {
		entries = append(entries, model.LeaderboardEntry{
			PlayerID: player.ID,
			Username: player.Username,
			Score:    player.TotalScore,
		})
	}

	// Sort by total score (descending)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	// Apply pagination
	total := len(entries)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedEntries := entries[start:end]

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"leaderboard": model.Leaderboard{
			GameID:  "global",
			Entries: paginatedEntries,
		},
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *handler) GetGameLeaderboard(ctx fiber.Ctx) error {
	gameID := ctx.Params("gameId")
	if gameID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Game ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

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

	// Collect scores for the specific game
	gameScores := make(map[string]int) // playerID -> best score
	for playerID, scores := range h.storage.GetPlayerScores() {
		bestScore := 0
		for _, score := range scores {
			if score.GameID == gameID && score.Score > bestScore {
				bestScore = score.Score
			}
		}
		if bestScore > 0 {
			gameScores[playerID] = bestScore
		}
	}

	var entries []model.LeaderboardEntry
	for playerID, score := range gameScores {
		if player, exists := h.storage.GetPlayers()[playerID]; exists {
			entries = append(entries, model.LeaderboardEntry{
				PlayerID: playerID,
				Username: player.Username,
				Score:    score,
			})
		}
	}

	// Sort by score (descending)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	// Apply pagination
	total := len(entries)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedEntries := entries[start:end]

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"leaderboard": model.Leaderboard{
			GameID:  gameID,
			Entries: paginatedEntries,
		},
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}