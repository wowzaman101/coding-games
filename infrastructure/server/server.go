package server

import (
	"coding-games/internal/handler/gamehdl"
	"coding-games/internal/handler/leaderboardhdl"
	"coding-games/internal/handler/playerhdl"

	"github.com/gofiber/fiber/v3"
)

func New(gh gamehdl.Handler, ph playerhdl.Handler, lh leaderboardhdl.Handler) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error":   "internal_error",
				"message": err.Error(),
				"code":    code,
			})
		},
	})

	// Health check route
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Coding Games API is running",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Game routes
	games := api.Group("/games")
	games.Get("/", gh.ListGames)           // GET /api/v1/games
	games.Post("/", gh.CreateGame)         // POST /api/v1/games
	games.Get("/:id", gh.GetGame)          // GET /api/v1/games/:id
	games.Put("/:id", gh.UpdateGame)       // PUT /api/v1/games/:id
	games.Delete("/:id", gh.DeleteGame)    // DELETE /api/v1/games/:id

	// Player routes
	players := api.Group("/players")
	players.Get("/", ph.ListPlayers)                    // GET /api/v1/players
	players.Post("/", ph.CreatePlayer)                  // POST /api/v1/players
	players.Get("/:id", ph.GetPlayer)                   // GET /api/v1/players/:id
	players.Put("/:id", ph.UpdatePlayer)                // PUT /api/v1/players/:id
	players.Delete("/:id", ph.DeletePlayer)             // DELETE /api/v1/players/:id
	players.Post("/:id/scores", ph.SubmitScore)         // POST /api/v1/players/:id/scores
	players.Get("/:id/scores", ph.GetPlayerScores)      // GET /api/v1/players/:id/scores

	// Leaderboard routes
	leaderboard := api.Group("/leaderboards")
	leaderboard.Get("/global", lh.GetGlobalLeaderboard)         // GET /api/v1/leaderboards/global
	leaderboard.Get("/games/:gameId", lh.GetGameLeaderboard)    // GET /api/v1/leaderboards/games/:gameId

	// Legacy routes for backward compatibility
	app.Get("/game/test", gh.Test)

	return app
}
