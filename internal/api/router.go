// internal/api/router.go
package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/igorschechtel/finance-tracker-backend/internal/api/handlers"
)

func SetupRouter(
	userHandler *handlers.UserHandler,
	expenseHandler *handlers.ExpenseHandler,
	categoryHandler *handlers.CategoryHandler,
) *chi.Mux {
	r := chi.NewRouter()
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

	// Add global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.List)
			r.Post("/", userHandler.Create)

			// User expense routes
			r.Route("/{userId}/expenses", func(r chi.Router) {
				r.Get("/", expenseHandler.ListByUser)
				r.Post("/", expenseHandler.Create)
			})

			// User category routes
			r.Route("/{userId}/categories", func(r chi.Router) {
				r.Get("/", categoryHandler.ListByUser)
				r.Post("/", categoryHandler.Create)
			})
		})

		// Other entity routes...
	})

	return r
}
