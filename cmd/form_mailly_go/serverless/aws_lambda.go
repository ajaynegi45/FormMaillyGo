package main

import (
	Form_Mailly_Go "Form-Mailly-Go"
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/handler"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

// GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap aws_lambda.go
// zip lambda-handler.zip bootstrap
// zip -r lambda-handler.zip bootstrap public/index.html

func main() {
	// Setup HTTP routes
	router := setupRoutes()

	// Safely load environment variables
	config.LoadEnvironmentVariable()

	// Start Lambda handler with the configured router
	lambda.Start(httpadapter.NewV2(router).ProxyWithContext)
}

// setupRoutes configures all HTTP endpoints
func setupRoutes() http.Handler {

	mux := http.NewServeMux()

	// Health check endpoint - represent home page
	mux.HandleFunc("GET /{$}", Form_Mailly_Go.HomeHandler)

	// Health check endpoint - useful for monitoring Lambda function
	mux.HandleFunc("GET /api/health", handler.HealthHandler)
	mux.HandleFunc("GET /api/runtime-info", handler.RuntimeInfoHandler)
	mux.HandleFunc("GET /api/metrics", handler.MetricsHandler)

	// For sending individual emails
	mux.HandleFunc("POST /api/contact", handler.ContactHandler)
	// For sending multiple emails efficiently
	mux.HandleFunc("POST /api/batch/contact", handler.BatchEmailProcessor)

	// Apply security middleware to all routes
	return applySecurityHeaders(mux)
}

// applySecurityHeaders adds essential security headers to all responses
func applySecurityHeaders(next http.Handler) http.Handler {
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
