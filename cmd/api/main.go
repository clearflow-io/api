package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/igorschechtel/finance-tracker-backend/internal/handlers"
)

func main() {
	// Create a new router
	r := chi.NewRouter()

	// Add some middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set up routes
	handlers.SetupRoutes(r)

	// Start the server
	port := "8080"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case <-shutdown:
		log.Println("Starting shutdown...")
		// Gracefully shutdown the server
		if err := server.Close(); err != nil {
			log.Fatalf("Could not stop server gracefully: %v", err)
		}
		log.Println("Server stopped")
	}
}
