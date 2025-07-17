package model

import "time"

// Player represents a player in the coding games platform
type Player struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	TotalScore int      `json:"total_score"`
	Rank      int       `json:"rank"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePlayerRequest represents the request payload for creating a player
type CreatePlayerRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

// UpdatePlayerRequest represents the request payload for updating a player
type UpdatePlayerRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
}

// PlayerScore represents a player's score for a specific game
type PlayerScore struct {
	PlayerID  string    `json:"player_id"`
	GameID    string    `json:"game_id"`
	Score     int       `json:"score"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SubmitScoreRequest represents the request payload for submitting a score
type SubmitScoreRequest struct {
	GameID string `json:"game_id" validate:"required"`
	Score  int    `json:"score" validate:"required,min=0"`
}