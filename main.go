package main

import (
	"net/http"
	"strconv"

	"github.com/igorschechtel/finance-tracker-backend/internal/api"
	"github.com/igorschechtel/finance-tracker-backend/internal/api/handlers"
	"github.com/igorschechtel/finance-tracker-backend/internal/config"
	"github.com/igorschechtel/finance-tracker-backend/internal/database"
	"github.com/igorschechtel/finance-tracker-backend/internal/repositories"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Database connection
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)

	// Handlers
	handlers := &api.Handlers{
		User:     handlers.NewUserHandler(userRepo),
		Expense:  handlers.NewExpenseHandler(expenseRepo),
		Category: handlers.NewCategoryHandler(categoryRepo),
	}
	router := api.SetupRouter(cfg, handlers)

	// Start server
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	logrus.Infof("Server starting on %s in %s mode", addr, cfg.Env)
	if err := http.ListenAndServe(addr, router); err != nil {
		logrus.WithError(err).Fatal("Server failed to start")
	}
}
