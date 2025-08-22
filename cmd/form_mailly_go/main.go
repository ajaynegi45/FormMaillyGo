package main

import (
	Form_Mailly_Go "Form-Mailly-Go"
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/handler"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/joho/godotenv"
)

// init runs before main and checks environment configuration
func init() {
	// Load .env files (for local development only)
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Fatal("Error loading .env file!")
	}
	// Safely load environment variables
	config.LoadEnvironmentVariable()
}

func main() {

	// Mux Router with optimized routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", Form_Mailly_Go.HomeHandler)
	mux.HandleFunc("GET /api/health", handler.HealthHandler)
	mux.HandleFunc("GET /api/runtime-info", handler.RuntimeInfoHandler)
	mux.HandleFunc("GET /api/metrics", handler.MetricsHandler)

	//mux.HandleFunc("POST /api/contact", handler.ContactHandler)
	mux.HandleFunc("POST /api/batch/contact", handler.BatchEmailProcessor)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     securityHeadersMiddleware(mux),
		ReadTimeout: 10 * time.Second,
		//WriteTimeout: 30 * time.Second,
		IdleTimeout: 60 * time.Second,
		//MaxHeaderBytes: 1 << 20, // 1MB
	}

	// simulateLambdaLimits()

	// Start server
	fmt.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %v", err)
	}
}

// securityHeadersMiddleware adds essential security headers to all responses
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("X-Content-Type-Options", "nosniff")                  // Prevent MIME type sniffing attacks
		headers.Set("X-Frame-Options", "DENY")                            // Prevent clickjacking attacks
		headers.Set("Referrer-Policy", "strict-origin-when-cross-origin") // Control referrer information
		headers.Set("Referrer-Policy", "strict-origin-when-cross-origin") // Control referrer information
		// Baseline CSP — adjust 'script-src' as needed if you serve inline scripts
		headers.Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self'; connect-src 'self'")

		// Enable HSTS for HTTPS connections
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			headers.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}

		// CORS headers for cross-origin requests
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		headers.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Accept-Language, Last-Event-ID")
		headers.Set("Content-Type", "application/json")

		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			// Optional: cache preflight for 10 minutes
			headers.Set("Access-Control-Max-Age", "600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// simulateLambdaLimits applies Go runtime limits to mimic AWS Lambda constraints.
func simulateLambdaLimits() {
	const memoryLimit = 180 * (1 << 20) // ≈ 180 MiB
	prev := debug.SetMemoryLimit(memoryLimit)
	fmt.Printf("Memory limit set to %d bytes (previous: %d)\n", memoryLimit, prev)

	// Make GC more aggressive to respect memory cap
	//prevGOGC := debug.SetGCPercent(50)
	//fmt.Printf("GOGC set to %d (previous: %d)\n", 50, prevGOGC)

	prevProcs := runtime.GOMAXPROCS(2) // match Lambda’s 2 logical cores
	fmt.Printf("GOMAXPROCS set to %d (previous: %d)\n", 2, prevProcs)

	fmt.Println("Lambda-like constraints applied")
}
