package gamehdl

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

func TestGameHandler_ListGames(t *testing.T) {
	app := fiber.New()
	storage := storage.NewInMemoryStorage()
	handler := New(storage)

	app.Get("/games", handler.ListGames)

	req := httptest.NewRequest("GET", "/games", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "games")
	assert.Contains(t, response, "total")
	assert.Equal(t, float64(2), response["total"]) // 2 sample games
}

func TestGameHandler_GetGame(t *testing.T) {
	app := fiber.New()
	storage := storage.NewInMemoryStorage()
	handler := New(storage)

	app.Get("/games/:id", handler.GetGame)

	// Test existing game
	req := httptest.NewRequest("GET", "/games/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var game model.Game
	err = json.NewDecoder(resp.Body).Decode(&game)
	assert.NoError(t, err)
	assert.Equal(t, "1", game.ID)
	assert.Equal(t, "Two Sum", game.Title)

	// Test non-existing game
	req = httptest.NewRequest("GET", "/games/999", nil)
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGameHandler_CreateGame(t *testing.T) {
	app := fiber.New()
	storage := storage.NewInMemoryStorage()
	handler := New(storage)

	app.Post("/games", handler.CreateGame)

	// Test valid game creation
	gameReq := model.CreateGameRequest{
		Title:       "Test Game",
		Description: "A test game",
		Difficulty:  "easy",
		Language:    "go",
		MaxScore:    100,
	}
	body, _ := json.Marshal(gameReq)

	req := httptest.NewRequest("POST", "/games", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var game model.Game
	err = json.NewDecoder(resp.Body).Decode(&game)
	assert.NoError(t, err)
	assert.Equal(t, "Test Game", game.Title)
	assert.Equal(t, "easy", game.Difficulty)

	// Test invalid game creation
	invalidReq := model.CreateGameRequest{
		Title:      "Test Game",
		Difficulty: "invalid",
		MaxScore:   -1,
	}
	body, _ = json.Marshal(invalidReq)

	req = httptest.NewRequest("POST", "/games", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGameHandler_ListGamesWithFilters(t *testing.T) {
	app := fiber.New()
	storage := storage.NewInMemoryStorage()
	handler := New(storage)

	app.Get("/games", handler.ListGames)

	// Test difficulty filter
	req := httptest.NewRequest("GET", "/games?difficulty=easy", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	
	games := response["games"].([]interface{})
	assert.Equal(t, 1, len(games)) // Only one easy game

	// Test language filter
	req = httptest.NewRequest("GET", "/games?language=python", nil)
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	
	games = response["games"].([]interface{})
	assert.Equal(t, 1, len(games)) // Only one python game
}