package storage

import (
	"coding-games/internal/model"
	"time"
)

// InMemoryStorage represents an in-memory storage for the example API
type InMemoryStorage struct {
	players map[string]*model.Player
	scores  map[string][]model.PlayerScore
	games   map[string]*model.Game
}

type Storage interface {
	// Game methods
	GetGames() map[string]*model.Game
	
	// Player methods
	GetPlayers() map[string]*model.Player
	GetPlayerScores() map[string][]model.PlayerScore
}

func NewInMemoryStorage() Storage {
	// Initialize with sample data
	games := map[string]*model.Game{
		"1": {
			ID:          "1",
			Title:       "Two Sum",
			Description: "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
			Difficulty:  "easy",
			Language:    "go",
			MaxScore:    100,
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
		},
		"2": {
			ID:          "2",
			Title:       "Binary Tree Traversal",
			Description: "Implement inorder, preorder, and postorder traversal of a binary tree.",
			Difficulty:  "medium",
			Language:    "python",
			MaxScore:    200,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
			UpdatedAt:   time.Now().Add(-12 * time.Hour),
		},
	}

	players := map[string]*model.Player{
		"1": {
			ID:         "1",
			Username:   "alice_coder",
			Email:      "alice@example.com",
			TotalScore: 250,
			Rank:       1,
			CreatedAt:  time.Now().Add(-72 * time.Hour),
			UpdatedAt:  time.Now().Add(-1 * time.Hour),
		},
		"2": {
			ID:         "2",
			Username:   "bob_dev",
			Email:      "bob@example.com",
			TotalScore: 180,
			Rank:       2,
			CreatedAt:  time.Now().Add(-48 * time.Hour),
			UpdatedAt:  time.Now().Add(-2 * time.Hour),
		},
	}

	scores := map[string][]model.PlayerScore{
		"1": {
			{
				PlayerID:  "1",
				GameID:    "1",
				Score:     100,
				Completed: true,
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			{
				PlayerID:  "1",
				GameID:    "2",
				Score:     150,
				Completed: true,
				CreatedAt: time.Now().Add(-12 * time.Hour),
				UpdatedAt: time.Now().Add(-12 * time.Hour),
			},
		},
		"2": {
			{
				PlayerID:  "2",
				GameID:    "1",
				Score:     80,
				Completed: true,
				CreatedAt: time.Now().Add(-20 * time.Hour),
				UpdatedAt: time.Now().Add(-20 * time.Hour),
			},
			{
				PlayerID:  "2",
				GameID:    "2",
				Score:     100,
				Completed: false,
				CreatedAt: time.Now().Add(-6 * time.Hour),
				UpdatedAt: time.Now().Add(-6 * time.Hour),
			},
		},
	}

	return &InMemoryStorage{
		games:   games,
		players: players,
		scores:  scores,
	}
}

func (s *InMemoryStorage) GetGames() map[string]*model.Game {
	return s.games
}

func (s *InMemoryStorage) GetPlayers() map[string]*model.Player {
	return s.players
}

func (s *InMemoryStorage) GetPlayerScores() map[string][]model.PlayerScore {
	return s.scores
}