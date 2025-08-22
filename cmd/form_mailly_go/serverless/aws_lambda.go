package main

import (
	"Form-Mailly-Go"
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
		headers.Set("X-XSS-Protection", "1; mode=block")                  // Enable XSS protection in browsers
		headers.Set("Referrer-Policy", "strict-origin-when-cross-origin") // Control referrer information

		// Enable HSTS for HTTPS connections
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			headers.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}

		// CORS headers for cross-origin requests
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		headers.Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
