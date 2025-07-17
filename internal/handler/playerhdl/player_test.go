package playerhdl

import (
	"bytes"
	"coding-games/internal/model"
	"coding-games/internal/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestPlayerHandler_ListPlayers(t *testing.T) {
	app := fiber.New()
	storage := storage.NewInMemoryStorage()
	handler := New(storage)

	app.Get("/players", handler.ListPlayers)

	req := httptest.NewRequest("GET", "/players", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "players")
	assert.Contains(t, response, "total")
	assert.Equal(t, float64(2), response["total"]) // 2 sample players
}

func TestPlayerHandler_CreatePlayer(t *testing.T) {
	app := fiber.New()
	storage := storage.NewInMemoryStorage()
	handler := New(storage)

	app.Post("/players", handler.CreatePlayer)

	// Test valid player creation
	playerReq := model.CreatePlayerRequest{
		Username: "test_user",
		Email:    "test@example.com",
	}
	body, _ := json.Marshal(playerReq)

	req := httptest.NewRequest("POST", "/players", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var player model.Player
	err = json.NewDecoder(resp.Body).Decode(&player)
	assert.NoError(t, err)
	assert.Equal(t, "test_user", player.Username)
	assert.Equal(t, "test@example.com", player.Email)
	assert.Equal(t, 0, player.TotalScore)

	// Test duplicate username
	req = httptest.NewRequest("POST", "/players", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestPlayerHandler_SubmitScore(t *testing.T) {
	app := fiber.New()
	storage := storage.NewInMemoryStorage()
	handler := New(storage)

	app.Post("/players/:id/scores", handler.SubmitScore)

	// Test valid score submission
	scoreReq := model.SubmitScoreRequest{
		GameID: "1",
		Score:  150,
	}
	body, _ := json.Marshal(scoreReq)

	req := httptest.NewRequest("POST", "/players/1/scores", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response model.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Score submitted successfully", response.Message)

	// Test invalid player ID
	req = httptest.NewRequest("POST", "/players/999/scores", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestPlayerHandler_GetPlayerScores(t *testing.T) {
	app := fiber.New()
	storage := storage.NewInMemoryStorage()
	handler := New(storage)

	app.Get("/players/:id/scores", handler.GetPlayerScores)

	req := httptest.NewRequest("GET", "/players/1/scores", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "scores")
	assert.Contains(t, response, "total")
	assert.Equal(t, "1", response["player_id"])
}