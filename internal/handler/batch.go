package handler

import (
	"Form-Mailly-Go/internal/model"
	"Form-Mailly-Go/internal/service"
	"Form-Mailly-Go/internal/validation"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	numWorkers = 2
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func BatchEmailProcessor(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	totalStart := time.Now()

	// Json to object Processing
	var emailList []model.Email
	if err := json.NewDecoder(request.Body).Decode(&emailList); err != nil {
		http.Error(response, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// ---------------------
	// Validate the Email Data
	for _, email := range emailList {

		errMsg := validateBatchEmailData(email)
		if errMsg != "" {
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(response).Encode(struct {
				Error string `json:"error"`
			}{Error: errMsg})
			if err != nil {
				return
			}
			return
		}
	}

	// Setting headers
	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set("Connection", "keep-alive")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := response.(http.Flusher)
	if !ok {
		http.Error(response, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	emailChan := make(chan model.Email, len(emailList))
	resultChan := make(chan *model.EmailResult, len(emailList))

	// Detect client disconnect
	notify := request.Context().Done()

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			EachSMTPWorkerStart := time.Now()
			conn, err := service.SetupNewSMTPConnection()
			if err != nil {
				fmt.Printf("Worker %d: failed to setup SMTP connection: %v\n", workerID, err)
				return // Bail out this worker to avoid spinning with a nil client
			}
			fmt.Printf("Worker %d setup time: %v\n", workerID, time.Since(EachSMTPWorkerStart))

			defer service.CloseSMTPConnection(conn)

			for {
				select {

				case <-notify:
					fmt.Println("Client disconnected - Email processing stopped")
					return

				case email, ok := <-emailChan:
					EmailSentEach := time.Now()

					if !ok {
						return // channel closed, no more jobs
					}

					err := service.SendEmailUsingWorker(conn, &email)

					res := &model.EmailResult{Email: email.SentTo}
					if err != nil {
						res.Status = "failed"
						res.Error = err.Error()
					} else {
						res.Status = "success"
					}
					resultChan <- res
					fmt.Println("Each time taken is", time.Since(EmailSentEach))
				}
			}

		}(i)
	}

	// feeder goroutine
	wg.Go(func() {
		for _, email := range emailList {
			emailChan <- email
		}
		close(emailChan)
	})

	// result sender: drains resultChan and streams to client
	var resultWG sync.WaitGroup
	resultWG.Add(1)
	go (func() {
		defer resultWG.Done()
		ResultSendingTime := time.Now()

		// Detect client disconnect
		notify := request.Context().Done()

		// Stream updates
		for result := range resultChan {
			select {
			case <-notify:
				fmt.Println("Client disconnected - Result sending stopped")
				return
			default:
				// Get and reset buffer
				buf := bufPool.Get().(*bytes.Buffer)
				buf.Reset() // Clears old data before reuse — avoids data corruption or leaks.

				// Encode JSON into buffer
				err := json.NewEncoder(buf).Encode(result)
				if err == nil {
					// Stream it using SSE format: "data: <json>\n\n"
					_, err = fmt.Fprintf(response, "data: %s\n\n", buf.Bytes())
					if err != nil {
						// If the write fails (e.g. client gone), exit
						bufPool.Put(buf)
						return
					}

					flusher.Flush()
				}
				// Return buffer to pool
				bufPool.Put(buf)
			}
		}
		fmt.Println("Result Sending time taken is", time.Since(ResultSendingTime))
	})()

	// Once all producers are done, close resultChan so result-sender can finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// wait for result-sender to finish (which happens once resultChan is closed/drained)
	resultWG.Wait()
	totalDuration := time.Since(totalStart)
	fmt.Println("Total time taken:", totalDuration)
}

func validateBatchEmailData(email model.Email) string {

	validator := validation.NewValidator()
	fields := []validation.Field{
		{
			Name:  "email",
			Value: &email.SentTo,
			Rules: []validation.Rule{
				validation.RequiredRule(),
				validation.EmailRule(),
				validation.MaxLengthRule(255),
			},
		},
		{
			Name:  "subject",
			Value: &email.Subject,
			Rules: []validation.Rule{
				validation.RequiredRule(),
				validation.MaxLengthRule(300),
			},
		},
		{
			Name:  "message",
			Value: &email.Message,
			Rules: []validation.Rule{
				validation.RequiredRule(),
			},
		},
		{
			Name:  "product_name",
			Value: &email.ProductName,
			Rules: []validation.Rule{
				validation.ProductNameRule(),
			},
		},
	}

	for _, field := range fields {
		validator.ValidateField(field)
		if !validator.IsValid() {
			// Return immediately once an error occurs
			return validator.Error
		}
	}

	return "" // no error found, valid form
}

// getNumberOfWorkers calculates the optimal number of SMTP workers for email batch processing.
// This function is designed to handle batches ranging from 1 email to 3000+ emails efficiently.
// It considers SMTP connection setup overhead, CPU cores, and diminishing returns from too many workers.
func getNumberOfWorkers(batchSize int) int {
	numCPU := runtime.NumCPU()

	// Core constants based on SMTP processing characteristics
	const (
		// Each SMTP connection has setup/teardown overhead, so we need minimum emails per worker
		minEmailsPerWorker = 8

		// Hard ceiling to prevent SMTP server overload and resource exhaustion
		maxWorkers = 25

		// Always have at least one worker
		minWorkers = 1

		// For I/O bound SMTP operations, we can use more workers than CPU cores
		// But not too many since each worker holds a persistent SMTP connection
		ioMultiplier = 2.0
	)

	var workers int

	// The core algorithm: balance connection overhead against parallelism benefits
	if batchSize <= 5 {
		// Micro batches: Single worker is most efficient
		// Connection setup time would dominate with multiple workers
		workers = 1
	} else if batchSize <= 25 {
		// Small batches: Limited workers to avoid setup overhead waste
		// Use at most half the batch size, but respect CPU limits
		workers = min(batchSize/2, numCPU)
		if workers == 0 {
			workers = 1
		}
	} else {
		// Medium to large batches: Calculate based on efficiency and CPU capacity

		// Calculate workers needed for efficient batch processing
		efficiencyBasedWorkers := batchSize / minEmailsPerWorker

		// Calculate workers based on CPU capacity for I/O bound work
		cpuBasedWorkers := int(float64(numCPU) * ioMultiplier)

		// Use the smaller of the two calculations to avoid over-provisioning
		workers = min(efficiencyBasedWorkers, cpuBasedWorkers)

		// For very large batches (1000+), we can be more aggressive with worker count
		// since the connection setup cost becomes negligible per email
		if batchSize >= 1000 {
			workers = min(maxWorkers, batchSize/minEmailsPerWorker)
		}
	}

	// Apply absolute bounds to prevent edge cases
	if workers > maxWorkers {
		workers = maxWorkers
	} else if workers < minWorkers {
		workers = minWorkers
	}

	return workers
}
