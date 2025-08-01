package handler

import (
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/model"
	"Form-Mailly-Go/internal/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	var form model.ContactForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validator
	validationErrors := service.ValidateFormData(&form)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": validationErrors,
		})
		if err != nil {
			return
		}
		return
	}

	// Load configuration
	emailConfig := config.LoadEmailConfig()

	// Initialize services
	emailService := service.NewSMTPEmailService(emailConfig)

	if err := emailService.Send(&form); err != nil {
		http.Error(w, `{"error": "Failed to send email"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte(`{"message": "Email sent successfully"}`))
	if err != nil {
		return
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte("FormMaillyGo service is running fine. ❤️"))
	if err != nil {
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Open file on demand
	f, err := os.Open("./public/index.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	// Get file size for Content-Length
	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	size := fi.Size()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))

	// Zero-copy transfer via ReaderFrom → sendfile under the hood
	if rf, ok := w.(io.ReaderFrom); ok {
		if _, err := rf.ReadFrom(f); err != nil {
			log.Printf("sendfile error: %v", err)
		}
	} else {
		// Fallback if ReaderFrom not supported (unlikely)
		if _, err := io.Copy(w, f); err != nil {
			log.Printf("io.Copy error: %v", err)
		}
	}
}
