package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/igorschechtel/clearflow-backend/internal/api"
	"github.com/igorschechtel/clearflow-backend/internal/api/handlers"
	"github.com/igorschechtel/clearflow-backend/internal/config"
	"github.com/igorschechtel/clearflow-backend/internal/database"
	"github.com/igorschechtel/clearflow-backend/internal/repositories"
	"github.com/igorschechtel/clearflow-backend/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Initialize validator
	v := validator.New(validator.WithRequiredStructEnabled())

	// Database connection
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// If --check-db flag is passed, just verify connection and exit
	if len(os.Args) > 1 && os.Args[1] == "--check-db" {
		if err := db.Ping(); err != nil {
			logrus.WithError(err).Fatal("Database ping failed")
		}
		logrus.Info("Database connection check successful")
		return
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logrus.WithError(err).Fatal("Failed to run database migrations")
	}

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)

	// Services
	userService := services.NewUserService(userRepo)
	expenseService := services.NewExpenseService(expenseRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Handlers
	handlers := &api.Handlers{
		User:     handlers.NewUserHandler(userService, v),
		Expense:  handlers.NewExpenseHandler(expenseService, v),
		Category: handlers.NewCategoryHandler(categoryService, v),
	}
	router := api.SetupRouter(cfg, handlers, db)

	// Start server
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	logrus.Infof("Server starting on %s in %s mode", addr, cfg.Env)
	if err := http.ListenAndServe(addr, router); err != nil {
		logrus.WithError(err).Fatal("Server failed to start")
	}
}
