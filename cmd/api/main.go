// cmd/api/main.go
package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/igorschechtel/finance-tracker-backend/internal/api"
	"github.com/igorschechtel/finance-tracker-backend/internal/api/handlers"
	"github.com/igorschechtel/finance-tracker-backend/internal/config"
	"github.com/igorschechtel/finance-tracker-backend/internal/database"
	"github.com/igorschechtel/finance-tracker-backend/internal/repositories"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Database connection
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := repositories.NewUserRepository(db)

	// Handlers
	userHandler := handlers.NewUserHandler(userRepo)

	// Setup router
	router := api.SetupRouter(userHandler)

	// Start server
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	log.Printf("Server starting on %s in %s mode", addr, cfg.Env)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
