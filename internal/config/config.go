// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Clerk    ClerkConfig
	Env      string
}

type DatabaseConfig struct {
	URL string
}

type ServerConfig struct {
	Port           int
	AllowedOrigins []string
}

type ClerkConfig struct {
	SecretKey string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	// Database config
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}
	dbConfig := DatabaseConfig{
		URL: dbURL,
	}

	// Server config
	serverPort, _ := strconv.Atoi(getEnv("PORT", "8080"))
	allowedOrigins := strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ",")
	serverConfig := ServerConfig{
		Port:           serverPort,
		AllowedOrigins: allowedOrigins,
	}

	clerkConfig := ClerkConfig{
		SecretKey: getEnv("CLERK_SECRET_KEY", ""),
	}

	return &Config{
		Database: dbConfig,
		Server:   serverConfig,
		Clerk:    clerkConfig,
		Env:      getEnv("ENV", "development"),
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ConnectionString returns the database connection string
func (c *DatabaseConfig) ConnectionString() string {
	return c.URL
}
