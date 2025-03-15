// internal/api/router.go
package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/igorschechtel/finance-tracker-backend/internal/api/handlers"
)

func SetupRouter(
	userHandler *handlers.UserHandler,
	expenseHandler *handlers.ExpenseHandler,
	categoryHandler *handlers.CategoryHandler,
) *chi.Mux {
	r := chi.NewRouter()

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
