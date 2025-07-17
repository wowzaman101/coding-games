//go:build wireinject
// +build wireinject

package main

import (
	"coding-games/config"
	"coding-games/infrastructure/server"
	"coding-games/internal/handler/gamehdl"
	"coding-games/internal/handler/leaderboardhdl"
	"coding-games/internal/handler/playerhdl"
	"coding-games/internal/storage"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
)

type dependencies struct {
	server *fiber.App
}

func initialize() *dependencies {
	wire.Build(
		storage.NewInMemoryStorage,
		gamehdl.New,
		playerhdl.New,
		leaderboardhdl.New,
		server.New,
		wire.Struct(new(dependencies), "server"),
	)
	return nil
}

func main() {
	d := initialize()

	log.Printf("Server is running on port %s", config.Get().Server.Port)

	// Graceful shutdown and other server logic can be added here
	defer func() {
		if err := d.server.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %s", err.Error())
		}
		log.Println("Server gracefully stopped")
	}()

	// Start the server on the configured port
	if err := d.server.Listen(":" + config.Get().Server.Port); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
