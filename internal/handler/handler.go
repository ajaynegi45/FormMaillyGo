package handler

import (
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/model"
	"Form-Mailly-Go/internal/service"
	"Form-Mailly-Go/internal/validation"
	"encoding/json"
	"net/http"
)

func ContactHandler(response http.ResponseWriter, request *http.Request) {
	var form model.ContactForm
	if err := json.NewDecoder(request.Body).Decode(&form); err != nil {
		http.Error(response, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validator
	errMsg := validation.ValidateContactForm(&form)
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

	// Load configuration
	emailConfig := config.LoadEnvironmentVariable()

	// Initialize services
	emailService := service.NewSMTPEmailService(emailConfig)

	if err := emailService.Send(&form); err != nil {
		http.Error(response, `{"error": "Failed to send email"}`, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
	_, err := response.Write([]byte(`{"message": "Email sent successfully"}`))
	if err != nil {
		return
	}
}

func HealthHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := response.Write([]byte("FormMaillyGo service is running fine. ❤️"))
	if err != nil {
		return
	}
}
