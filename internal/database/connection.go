// internal/database/connection.go
package database

import (
	"database/sql"
	"fmt"

	"github.com/igorschechtel/clearflow-backend/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sirupsen/logrus"
)

// NewConnection establishes a connection to the database
func NewConnection(config config.DatabaseConfig) (*sql.DB, error) {
	connStr := config.ConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("Database connection established")
	return db, nil
}
