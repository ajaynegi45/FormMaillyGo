package main

import (
	"Form-Mailly-Go"
	"Form-Mailly-Go/internal/handler"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

func main() {

	// Mux Router with optimized routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", Form_Mailly_Go.HomeHandler)
	mux.HandleFunc("GET /api/health", handler.HealthHandler)
	mux.HandleFunc("POST /api/contact", handler.ContactHandler)

	// Load .env files
	envErr := godotenv.Load(".env.dev")
	if envErr != nil {
		log.Fatal("Error loading .env file!")
	}

	// Configure server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      securityHeadersMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	fmt.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %v", err)
	}
}

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}
