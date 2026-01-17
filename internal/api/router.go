// internal/api/router.go
package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/igorschechtel/clearflow-backend/internal/api/handlers"
	"github.com/igorschechtel/clearflow-backend/internal/auth"
	"github.com/igorschechtel/clearflow-backend/internal/config"
	"github.com/igorschechtel/clearflow-backend/internal/services"
	u "github.com/igorschechtel/clearflow-backend/internal/utils"
	"database/sql"
	"net/http"
)

type Handlers struct {
	User     *handlers.UserHandler
	Expense  *handlers.ExpenseHandler
	Category *handlers.CategoryHandler
}

func SetupRouter(
	cfg *config.Config,
	handlers *Handlers,
	db *sql.DB,
) *chi.Mux {
	r := chi.NewRouter()

	auth.InitializeClerk()

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			u.WriteJSONError(w, http.StatusServiceUnavailable, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.Server.AllowedOrigins,
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
	r.Use(u.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", handlers.User.List)
			r.Post("/", handlers.User.Create)
		})

		// Protected routes
		protected := r.With(auth.ClerkAuthMiddleware())

		// User expense routes
		protected.Route("/expenses", func(r chi.Router) {
			r.Get("/", handlers.Expense.ListByUser)
			r.Post("/", handlers.Expense.Create)
		})

		// User category routes
		protected.Route("/categories", func(r chi.Router) {
			r.Get("/", handlers.Category.ListByUser)
			r.Post("/", handlers.Category.Create)
		})

	})

	return r
}
