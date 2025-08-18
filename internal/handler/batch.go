package handler

import (
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/model"
	"Form-Mailly-Go/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	numWorkers     = 2
	jobQueueLength = 100 // buffer for jobs
)

type StatusUpdate struct {
	Email  string `json:"email"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func BatchEmailHandler(response http.ResponseWriter, request *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(request.Body)

	totalStart := time.Now()

	var forms []model.ContactForm
	if err := json.NewDecoder(request.Body).Decode(&forms); err != nil {
		http.Error(response, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// SSE headers
	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set("Connection", "keep-alive")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := response.(http.Flusher)
	if !ok {
		http.Error(response, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	emailConfig := config.LoadEnvironmentVariable()

	jobChan := make(chan model.ContactForm, jobQueueLength)
	statusChan := make(chan StatusUpdate, jobQueueLength)
	var wg sync.WaitGroup

	// Create service connections
	services := make([]*service.SMTPEmailService, numWorkers)
	for i := 0; i < numWorkers; i++ {
		services[i] = service.NewSMTPEmailService(emailConfig)
		if err := services[i].Connect(); err != nil {
			log.Fatalf("Worker %d failed to connect: %v", i, err)
		}
	}

	// Start workers (fixed loop)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(services[i], jobChan, statusChan, &wg)
	}

	// Feed jobs
	go func() {
		for _, form := range forms {
			jobChan <- form
		}
		close(jobChan)
	}()

	// Close statusChan when all workers done
	go func() {
		wg.Wait()
		close(statusChan)
	}()

	// Detect client disconnect
	notify := request.Context().Done()

	// Stream updates
	for {
		select {
		case <-notify:
			fmt.Println("Client disconnected")
			return
		case update, ok := <-statusChan:
			if !ok {
				time.Sleep(1 * time.Second) // simulate network delay
				totalDuration := time.Since(totalStart)
				fmt.Println("Total time taken:", totalDuration)

				// Close all after workers are done
				for _, s := range services {
					s.Close()
				}
				return
			}
			if jsonData, err := json.Marshal(update); err == nil {
				_, err := fmt.Fprintf(response, "data: %s\n\n", jsonData)
				if err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}
}

func worker(emailService service.EmailService, jobs <-chan model.ContactForm, statusChan chan<- StatusUpdate, wg *sync.WaitGroup) {
	defer wg.Done()
	for contact := range jobs {
		start := time.Now()
		status := StatusUpdate{Email: contact.Email}

		if err := emailService.SendBatchEmail(&contact); err != nil {
			status.Status = "failed"
			status.Error = err.Error()
		} else {
			status.Status = "success"
		}

		fmt.Println("Each time taken is", time.Since(start))
		statusChan <- status
	}
}
