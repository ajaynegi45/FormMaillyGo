package main

import (
	"Form-Mailly-Go/internal/handler"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "FORM-MAILER: ", log.LstdFlags|log.Lshortfile)

	// Create handlers
	contactHandler := handler.NewContactHandler(logger)
	healthHandler := &handler.HealthHandler{}

	// Setup router
	mux := http.NewServeMux()

	mux.Handle("GET /api/health", healthHandler)
	mux.Handle("POST /api/contact", contactHandler)

	// Configure server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      securityHeadersMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	logger.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}
