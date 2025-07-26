package handler

import (
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/model"
	"Form-Mailly-Go/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type ContactHandler struct {
	logger *log.Logger
}

func NewContactHandler(l *log.Logger) *ContactHandler {
	return &ContactHandler{
		logger: l,
	}
}

func (h *ContactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var form model.ContactForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		h.logger.Printf("JSON decode error: %v", err)
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validator
	validationErrors := service.ValidateFormData(&form)
	if len(validationErrors) > 0 {
		h.logger.Printf("Validation errors: %v", validationErrors)
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
		h.logger.Printf("Email send error: %v", err)
		http.Error(w, `{"error": "Failed to send email"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte(`{"message": "Email sent successfully"}`))
	if err != nil {
		return
	}
}

type HealthHandler struct{}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Service is healthy"))
	if err != nil {
		return
	}
}
