package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// SetupRoutes configures all the routes for our application
func SetupRoutes(r chi.Router) {
	r.Get("/", HomeHandler)
	r.Get("/health", HealthCheckHandler)
}

// HomeHandler handles the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the API!"))
}

// HealthCheckHandler handles the health check endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
