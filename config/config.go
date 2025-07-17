// Package config provides application configuration management.
// It loads configuration from environment variables with optional .env file support.
package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds the complete application configuration.
type Config struct {
	Server Server
}

// Server contains server-related configuration options.
type Server struct {
	Port string `envconfig:"PORT" default:"8080"`
}

var cfg Config

func init() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()
	
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to process environment configuration: %v", err)
	}
}

// Get returns the current application configuration.
func Get() Config {
	return cfg
}
