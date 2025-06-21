package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/GitNinja36/wello-backend/config"
	"github.com/GitNinja36/wello-backend/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(" Error loading .env file: %v", err)
	}

	config.ConnectDB()

	router := chi.NewRouter()

	// middlewares
	router.Use(middleware.Logger)    // Logs HTTP requests
	router.Use(middleware.Recoverer) // panic handling
	router.Use(middleware.RequestID) // Adds request ID
	router.Use(middleware.RealIP)    // Gets client IP

	// all routes
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf(" Wello server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
