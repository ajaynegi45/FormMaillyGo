package handler

import (
	"Form-Mailly-Go/internal/monitoring"
	"encoding/json"
	"net/http"
	"os"
	"runtime"
)

// HealthHandler provides comprehensive health and performance information
func HealthHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-cache")

	// Get version from environment or build info
	version := os.Getenv("VERSION")
	if version == "" {
		version = "development"
	}

	healthStatus := monitoring.PerformHealthCheck(request.Context(), version)

	// Set appropriate HTTP status based on health
	switch healthStatus.Status {
	case "healthy":
		response.WriteHeader(http.StatusOK)
	case "degraded":
		response.WriteHeader(http.StatusOK) // Still functional
	case "unhealthy":
		response.WriteHeader(http.StatusServiceUnavailable)
	default:
		response.WriteHeader(http.StatusInternalServerError)
	}

	// Return health status as JSON
	if err := json.NewEncoder(response).Encode(healthStatus); err != nil {
		http.Error(response, `{"error": "Failed to encode health status"}`, http.StatusInternalServerError)
	}
}

// RuntimeInfoHandler provides Go runtime information
func RuntimeInfoHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	runtimeInfo := map[string]interface{}{
		"go_version":    runtime.Version(),
		"num_cpu":       runtime.NumCPU(),
		"num_goroutine": runtime.NumGoroutine(),
		"gomaxprocs":    runtime.GOMAXPROCS(0),

		// Memory statistics
		"memory": map[string]interface{}{
			"alloc_bytes":       m.Alloc,
			"total_alloc_bytes": m.TotalAlloc,
			"sys_bytes":         m.Sys,
			"heap_alloc_bytes":  m.HeapAlloc,
			"heap_sys_bytes":    m.HeapSys,
			"heap_idle_bytes":   m.HeapIdle,
			"heap_inuse_bytes":  m.HeapInuse,
			"heap_objects":      m.HeapObjects,
			"stack_inuse_bytes": m.StackInuse,
			"stack_sys_bytes":   m.StackSys,
			"num_gc":            m.NumGC,
			"gc_cpu_fraction":   m.GCCPUFraction,
		},

		// Environment
		"environment": map[string]interface{}{
			"aws_lambda_function_name":    os.Getenv("AWS_LAMBDA_FUNCTION_NAME"),
			"aws_lambda_function_version": os.Getenv("AWS_LAMBDA_FUNCTION_VERSION"),
			"aws_region":                  os.Getenv("AWS_REGION"),
			"aws_execution_env":           os.Getenv("AWS_EXECUTION_ENV"),
		},
	}

	if err := json.NewEncoder(response).Encode(runtimeInfo); err != nil {
		http.Error(response, `{"error": "Failed to encode runtime info"}`, http.StatusInternalServerError)
	}
}

// MetricsHandler provides detailed performance metrics
func MetricsHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-cache, max-age=10")

	metrics := monitoring.GetMetrics()

	if err := json.NewEncoder(response).Encode(metrics); err != nil {
		http.Error(response, `{"error": "Failed to encode metrics"}`, http.StatusInternalServerError)
	}
}
