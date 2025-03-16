// internal/api/router.go
package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/igorschechtel/finance-tracker-backend/internal/api/handlers"
	"github.com/igorschechtel/finance-tracker-backend/internal/auth"
	"github.com/igorschechtel/finance-tracker-backend/internal/config"
	"github.com/sirupsen/logrus"
)

func SetupRouter(
	cfg *config.Config,
	userHandler *handlers.UserHandler,
	expenseHandler *handlers.ExpenseHandler,
	categoryHandler *handlers.CategoryHandler,
) *chi.Mux {
	r := chi.NewRouter()

	authClient, err := auth.InitializeHTTPClient()
	if err != nil {
		logrus.Fatal(err)
	}

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.List)
			r.Post("/", userHandler.Create)
		})

		// Protected routes
		protected := r.With(auth.AuthMiddleware(authClient, cfg))

		// User expense routes
		protected.Route("/expenses", func(r chi.Router) {
			r.Get("/", expenseHandler.ListByUser)
			r.Post("/", expenseHandler.Create)
		})

		// User category routes
		protected.Route("/categories", func(r chi.Router) {
			r.Get("/", categoryHandler.ListByUser)
			r.Post("/", categoryHandler.Create)
		})

	})

	return r
}
