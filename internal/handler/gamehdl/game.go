// Package gamehdl provides HTTP handlers for game-related operations.
package gamehdl

import "github.com/gofiber/fiber/v3"

// handler implements the Handler interface for game operations.
type handler struct {
}

// Handler defines the interface for game-related HTTP handlers.
type Handler interface {
	HandleGameTest(ctx fiber.Ctx) error
}

// New creates a new game handler instance.
func New() Handler {
	return &handler{}
}

// HandleGameTest handles the test endpoint for game functionality.
// This endpoint can be used for health checks or basic game service validation.
func (h *handler) HandleGameTest(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "Game service is operational",
		"status":  "success",
	})
}
