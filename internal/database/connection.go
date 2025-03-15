// internal/database/connection.go
package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/igorschechtel/finance-tracker-backend/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
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

	log.Println("Database connection established")
	return db, nil
}
