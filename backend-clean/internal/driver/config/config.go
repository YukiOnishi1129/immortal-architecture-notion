// Package config handles application configuration from environment variables.
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration.
type Config struct {
	// Database configuration
	DatabaseURL string

	// API server configuration
	ServerPort int

	// CORS configuration
	AllowedOrigins []string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	port := 8080
	if portStr := os.Getenv("API_PORT"); portStr != "" {
		parsedPort, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, errors.New("API_PORT must be a valid integer")
		}
		port = parsedPort
	}

	origins := parseAllowedOrigins(os.Getenv("CLIENT_ORIGIN"))

	return &Config{
		DatabaseURL:    dbURL,
		ServerPort:     port,
		AllowedOrigins: origins,
	}, nil
}

// parseAllowedOrigins parses comma-separated allowed origins from environment variable.
// Returns default localhost origins if empty.
func parseAllowedOrigins(fromEnv string) []string {
	if strings.TrimSpace(fromEnv) == "" {
		return []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	parts := strings.Split(fromEnv, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	if len(origins) == 0 {
		return []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	return origins
}
