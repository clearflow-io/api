// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database  DatabaseConfig
	Server    ServerConfig
	StackAuth StackAuthConfig
	Env       string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port int
}

type StackAuthConfig struct {
	ProjectID       string
	SecretServerKey string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	// Database config
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	dbConfig := DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", ""),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", ""),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Server config
	serverPort, _ := strconv.Atoi(getEnv("PORT", "8080"))
	serverConfig := ServerConfig{
		Port: serverPort,
	}

	stackAuthConfig := StackAuthConfig{
		ProjectID:       getEnv("STACK_PROJECT_ID", ""),
		SecretServerKey: getEnv("STACK_SECRET_SERVER_KEY", ""),
	}

	return &Config{
		Database:  dbConfig,
		Server:    serverConfig,
		StackAuth: stackAuthConfig,
		Env:       getEnv("ENV", "development"),
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

// ConnectionString returns a formatted connection string for the database
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
