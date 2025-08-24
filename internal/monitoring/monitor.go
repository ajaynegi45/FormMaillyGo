package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

// PerformanceMonitor tracks system performance metrics
type PerformanceMonitor struct {
	startTime      time.Time
	requestCount   int64
	errorCount     int64
	totalLatency   int64
	maxLatency     int64
	minLatency     int64
	emailsSent     int64
	emailsFailed   int64
	memoryPeak     int64
	goroutinesPeak int64
}

var globalMonitor = &PerformanceMonitor{
	startTime:  time.Now(),
	minLatency: int64(^uint64(0) >> 1), // Max int64 value
}

// RecordRequest records a request with its latency
func RecordRequest(duration time.Duration, success bool) {
	latency := duration.Nanoseconds()

	atomic.AddInt64(&globalMonitor.requestCount, 1)
	atomic.AddInt64(&globalMonitor.totalLatency, latency)

	// Update max latency
	for {
		current := atomic.LoadInt64(&globalMonitor.maxLatency)
		if latency <= current || atomic.CompareAndSwapInt64(&globalMonitor.maxLatency, current, latency) {
			break
		}
	}

	// Update min latency
	for {
		current := atomic.LoadInt64(&globalMonitor.minLatency)
		if latency >= current || atomic.CompareAndSwapInt64(&globalMonitor.minLatency, current, latency) {
			break
		}
	}

	if !success {
		atomic.AddInt64(&globalMonitor.errorCount, 1)
	}
}

// RecordEmail records service sending statistics
func RecordEmail(success bool) {
	if success {
		atomic.AddInt64(&globalMonitor.emailsSent, 1)
	} else {
		atomic.AddInt64(&globalMonitor.emailsFailed, 1)
	}
}

// UpdateSystemMetrics updates system-level metrics
func UpdateSystemMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Update memory peak
	memUsage := int64(m.Alloc)
	for {
		current := atomic.LoadInt64(&globalMonitor.memoryPeak)
		if memUsage <= current || atomic.CompareAndSwapInt64(&globalMonitor.memoryPeak, current, memUsage) {
			break
		}
	}

	// Update goroutines peak
	goroutines := int64(runtime.NumGoroutine())
	for {
		current := atomic.LoadInt64(&globalMonitor.goroutinesPeak)
		if goroutines <= current || atomic.CompareAndSwapInt64(&globalMonitor.goroutinesPeak, current, goroutines) {
			break
		}
	}
}

// Metrics represents current performance metrics
type Metrics struct {
	Uptime           time.Duration `json:"uptime"`
	RequestCount     int64         `json:"request_count"`
	ErrorCount       int64         `json:"error_count"`
	ErrorRate        float64       `json:"error_rate"`
	AverageLatency   time.Duration `json:"average_latency"`
	MinLatency       time.Duration `json:"min_latency"`
	MaxLatency       time.Duration `json:"max_latency"`
	RequestsPerSec   float64       `json:"requests_per_second"`
	EmailsSent       int64         `json:"emails_sent"`
	EmailsFailed     int64         `json:"emails_failed"`
	EmailSuccessRate float64       `json:"email_success_rate"`

	// System metrics
	MemoryUsage    int64  `json:"memory_usage_bytes"`
	MemoryPeak     int64  `json:"memory_peak_bytes"`
	Goroutines     int    `json:"goroutines"`
	GoroutinesPeak int64  `json:"goroutines_peak"`
	CPUs           int    `json:"cpus"`
	GoVersion      string `json:"go_version"`
}

// GetMetrics returns current performance metrics
func GetMetrics() *Metrics {
	UpdateSystemMetrics()

	uptime := time.Since(globalMonitor.startTime)
	requestCount := atomic.LoadInt64(&globalMonitor.requestCount)
	errorCount := atomic.LoadInt64(&globalMonitor.errorCount)
	totalLatency := atomic.LoadInt64(&globalMonitor.totalLatency)
	maxLatency := atomic.LoadInt64(&globalMonitor.maxLatency)
	minLatency := atomic.LoadInt64(&globalMonitor.minLatency)
	emailsSent := atomic.LoadInt64(&globalMonitor.emailsSent)
	emailsFailed := atomic.LoadInt64(&globalMonitor.emailsFailed)
	memoryPeak := atomic.LoadInt64(&globalMonitor.memoryPeak)
	goroutinesPeak := atomic.LoadInt64(&globalMonitor.goroutinesPeak)

	// Calculate derived metrics
	var errorRate float64
	if requestCount > 0 {
		errorRate = float64(errorCount) / float64(requestCount) * 100
	}

	var averageLatency time.Duration
	if requestCount > 0 {
		averageLatency = time.Duration(totalLatency / requestCount)
	}

	requestsPerSec := float64(requestCount) / uptime.Seconds()

	var emailSuccessRate float64
	totalEmails := emailsSent + emailsFailed
	if totalEmails > 0 {
		emailSuccessRate = float64(emailsSent) / float64(totalEmails) * 100
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &Metrics{
		Uptime:           uptime,
		RequestCount:     requestCount,
		ErrorCount:       errorCount,
		ErrorRate:        errorRate,
		AverageLatency:   averageLatency,
		MinLatency:       time.Duration(minLatency),
		MaxLatency:       time.Duration(maxLatency),
		RequestsPerSec:   requestsPerSec,
		EmailsSent:       emailsSent,
		EmailsFailed:     emailsFailed,
		EmailSuccessRate: emailSuccessRate,
		MemoryUsage:      int64(m.Alloc),
		MemoryPeak:       memoryPeak,
		Goroutines:       runtime.NumGoroutine(),
		GoroutinesPeak:   goroutinesPeak,
		CPUs:             runtime.NumCPU(),
		GoVersion:        runtime.Version(),
	}
}

// HealthCheck performs a comprehensive health check
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
	Checks    []Check   `json:"checks"`
	Metrics   *Metrics  `json:"metrics"`
}

type Check struct {
	Name    string        `json:"name"`
	Status  string        `json:"status"`
	Latency time.Duration `json:"latency"`
	Error   string        `json:"error,omitempty"`
}

// PerformHealthCheck performs comprehensive health checks
func PerformHealthCheck(ctx context.Context, version string) *HealthStatus {
	checks := []Check{}
	overallStatus := "healthy"

	// Memory check
	memCheck := checkMemoryUsage()
	checks = append(checks, memCheck)
	if memCheck.Status != "healthy" {
		overallStatus = "degraded"
	}

	// Goroutine check
	goroutineCheck := checkGoroutines()
	checks = append(checks, goroutineCheck)
	if goroutineCheck.Status != "healthy" {
		overallStatus = "degraded"
	}

	// Error rate check
	errorRateCheck := checkErrorRate()
	checks = append(checks, errorRateCheck)
	if errorRateCheck.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	return &HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   version,
		Checks:    checks,
		Metrics:   GetMetrics(),
	}
}

func checkMemoryUsage() Check {
	start := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Convert to MB for readability
	allocMB := float64(m.Alloc) / 1024 / 1024
	sysMB := float64(m.Sys) / 1024 / 1024

	status := "healthy"
	var errorMsg string

	// Alert if using more than 400MB (Lambda has 512MB limit)
	if allocMB > 400 {
		status = "unhealthy"
		errorMsg = fmt.Sprintf("High memory usage detected: your application is currently using %.2fMB of memory (out of %.2fMB reserved by the runtime).", allocMB, sysMB)
	} else if allocMB > 300 {
		status = "degraded"
		errorMsg = fmt.Sprintf("Elevated memory usage: your application is currently using %.2fMB of memory (out of %.2fMB reserved by the runtime).", allocMB, sysMB)
	}

	return Check{
		Name:    "memory",
		Status:  status,
		Latency: time.Since(start),
		Error:   errorMsg,
	}
}

func checkGoroutines() Check {
	start := time.Now()
	numGoroutines := runtime.NumGoroutine()

	status := "healthy"
	var errorMsg string

	// Alert on goroutine leaks
	if numGoroutines > 100 {
		status = "unhealthy"
		errorMsg = "Potential goroutine leak detected"
	} else if numGoroutines > 50 {
		status = "degraded"
		errorMsg = "High number of goroutines"
	}

	return Check{
		Name:    "goroutines",
		Status:  status,
		Latency: time.Since(start),
		Error:   errorMsg,
	}
}

func checkErrorRate() Check {
	start := time.Now()
	requestCount := atomic.LoadInt64(&globalMonitor.requestCount)
	errorCount := atomic.LoadInt64(&globalMonitor.errorCount)

	status := "healthy"
	var errorMsg string

	if requestCount > 0 {
		errorRate := float64(errorCount) / float64(requestCount) * 100

		if errorRate > 10 {
			status = "unhealthy"
			errorMsg = "High error rate detected"
		} else if errorRate > 5 {
			status = "degraded"
			errorMsg = "Elevated error rate"
		}
	}

	return Check{
		Name:    "error_rate",
		Status:  status,
		Latency: time.Since(start),
		Error:   errorMsg,
	}
}

// LogMetrics logs metrics to stdout (for CloudWatch)
func LogMetrics() {
	metrics := GetMetrics()
	if data, err := json.Marshal(metrics); err == nil {
		log.Printf("METRICS: %s", string(data))
	}
}

// MetricsMiddleware wraps handlers to record performance metrics
func MetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		wrapped := &responseWriter{ResponseWriter: response, statusCode: 200}

		next(wrapped, request)

		duration := time.Since(start)
		success := wrapped.statusCode < 400

		RecordRequest(duration, success)

		// Log slow requests
		if duration > 5*time.Second {
			log.Printf("SLOW_REQUEST: %s %s took %v", request.Method, request.URL.Path, duration)
		}
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// StartPeriodicMetricsLogging starts a goroutine that logs metrics periodically
func StartPeriodicMetricsLogging(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			LogMetrics()
		}
	}()
}

func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Println()
	fmt.Println("=== Memory Usage Start ===")
	fmt.Printf("   Alloc:        %.2f MB\n", float64(m.Alloc)/(1024*1024))
	fmt.Printf("   TotalAlloc:   %.2f MB\n", float64(m.TotalAlloc)/(1024*1024))
	fmt.Printf("   Sys:          %.2f MB\n", float64(m.Sys)/(1024*1024))
	fmt.Printf("   HeapAlloc:    %.2f MB\n", float64(m.HeapAlloc)/(1024*1024))
	fmt.Printf("   HeapSys:      %.2f MB\n", float64(m.HeapSys)/(1024*1024))
	fmt.Printf("   HeapIdle:     %.2f MB\n", float64(m.HeapIdle)/(1024*1024))
	fmt.Printf("   HeapInuse:    %.2f MB\n", float64(m.HeapInuse)/(1024*1024))
	fmt.Printf("   HeapReleased: %.2f MB\n", float64(m.HeapReleased)/(1024*1024))
	fmt.Printf("   StackInuse:   %.2f MB\n", float64(m.StackInuse)/(1024*1024))
	fmt.Printf("   StackSys:     %.2f MB\n", float64(m.StackSys)/(1024*1024))
	fmt.Printf("   GC Cycles:    %v\n", m.NumGC)
	fmt.Printf("   Total Pause:  %.2f ms\n", float64(m.PauseTotalNs)/1e6)
	fmt.Printf("   NextGC:       %.2f MB\n", float64(m.NextGC)/(1024*1024))
	fmt.Printf("   LastGC:       %s\n", time.Unix(0, int64(m.LastGC)).Format(time.RFC3339))
	fmt.Println("=== Memory Usage End ===")
	fmt.Println()
}
