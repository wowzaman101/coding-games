package model

import "time"

// Game represents a coding game
type Game struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"` // "easy", "medium", "hard"
	Language    string    `json:"language"`   // "go", "python", "java", etc.
	MaxScore    int       `json:"max_score"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateGameRequest represents the request payload for creating a game
type CreateGameRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Difficulty  string `json:"difficulty" validate:"required,oneof=easy medium hard"`
	Language    string `json:"language" validate:"required"`
	MaxScore    int    `json:"max_score" validate:"required,min=1"`
}

// UpdateGameRequest represents the request payload for updating a game
type UpdateGameRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Difficulty  *string `json:"difficulty,omitempty" validate:"omitempty,oneof=easy medium hard"`
	Language    *string `json:"language,omitempty"`
	MaxScore    *int    `json:"max_score,omitempty" validate:"omitempty,min=1"`
}