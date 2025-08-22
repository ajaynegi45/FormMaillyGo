package handler

import (
	"Form-Mailly-Go/internal/model"
	"Form-Mailly-Go/internal/service"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	// Validate the Data
	// ---------------------

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
	wg.Add(1)
	go (func() error {
		defer wg.Done()
		DataFeeding := time.Now()
		for _, email := range emailList {
			emailChan <- email
		}
		fmt.Println("DataFeeding time taken is", time.Since(DataFeeding))

		defer close(emailChan)
		return nil
	})()

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
