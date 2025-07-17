package model

// LeaderboardEntry represents an entry in the leaderboard
type LeaderboardEntry struct {
	Rank     int    `json:"rank"`
	PlayerID string `json:"player_id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

// Leaderboard represents the leaderboard for a game
type Leaderboard struct {
	GameID  string             `json:"game_id"`
	Entries []LeaderboardEntry `json:"entries"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}