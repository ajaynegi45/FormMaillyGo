package main

import (
	"Form-Mailly-Go/internal/handler"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"net/http"
)

// GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap aws_lambda.go
// zip lambda-handler.zip bootstrap
// zip -r lambda-handler.zip bootstrap public/index.html

func main() {

	// Mux Router with optimized routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handler.HomeHandler)
	mux.HandleFunc("GET /api/health", handler.HealthHandler)
	mux.HandleFunc("POST /api/contact", handler.ContactHandler)

	// Wrap mux with middleware
	handlerWithMiddleware := securityHeadersMiddlewareLambda(mux)

	// Instead of http.Server.ListenAndServe(), start Lambda handler:
	lambda.Start(httpadapter.NewV2(handlerWithMiddleware).ProxyWithContext)

}

// Security headers middleware, unchanged
func securityHeadersMiddlewareLambda(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}
