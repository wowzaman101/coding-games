// Package server provides HTTP server setup and route configuration.
package server

import (
	"coding-games/internal/handler/gamehdl"

	"github.com/gofiber/fiber/v3"
)

// New creates and configures a new Fiber application with all routes.
// It takes a game handler and sets up the necessary HTTP endpoints.
func New(gh gamehdl.Handler) *fiber.App {
	app := fiber.New()
	
	// Health check endpoint for monitoring and load balancer health checks
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"service": "coding-games",
		})
	})

	// Game-related endpoints
	app.Get("/game/test", gh.HandleGameTest)

	return app
}
