// internal/database/connection.go
package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/igorschechtel/clearflow-backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx PostgreSQL driver
	"github.com/sirupsen/logrus"
)

// NewConnection establishes a connection to the database
func NewConnection(config config.DatabaseConfig) (*sql.DB, error) {
	connStr := config.ConnectionString()

	// Use pgx driver for compatibility with Supabase Pooler (Transaction Mode)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("Database connection established")
	return db, nil
}

// RunMigrations runs the database migrations from the db/migrations directory
func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logrus.Info("Database migrations completed successfully")
	return nil
}
