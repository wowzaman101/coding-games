//go:build wireinject
// +build wireinject

// Package main is the entry point for the coding-games application.
// It uses Google Wire for dependency injection and sets up the HTTP server.
package main

import (
	"coding-games/config"
	"coding-games/infrastructure/server"
	"coding-games/internal/handler/gamehdl"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
)

// dependencies holds all the application dependencies.
type dependencies struct {
	server *fiber.App
}

// initialize sets up dependency injection using Wire.
func initialize() *dependencies {
	wire.Build(
		server.New,
		gamehdl.New,
		wire.Struct(new(dependencies), "server"),
	)
	return nil
}

func main() {
	// Initialize dependencies
	d := initialize()
	cfg := config.Get()

	log.Printf("Starting coding-games server on port %s", cfg.Server.Port)

	// Set up graceful shutdown
	defer func() {
		log.Println("Shutting down server...")
		if err := d.server.Shutdown(); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		} else {
			log.Println("Server gracefully stopped")
		}
	}()

	// Start the server
	serverAddr := ":" + cfg.Server.Port
	if err := d.server.Listen(serverAddr); err != nil {
		log.Fatalf("Failed to start server on %s: %v", serverAddr, err)
	}
}
