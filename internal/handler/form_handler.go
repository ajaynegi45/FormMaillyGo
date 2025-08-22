package handler

import (
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
	errMsg := validateContactForm(&form)
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

	if err := service.Send(&form); err != nil {
		http.Error(response, `{"error": "Failed to send service"}`, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
	_, err := response.Write([]byte(`{"message": "Email sent successfully"}`))
	if err != nil {
		return
	}
}

func validateContactForm(form *model.ContactForm) string {
	validator := validation.NewValidator()

	fields := []validation.Field{
		{
			Name:  "name",
			Value: &form.Name,
			Rules: []validation.Rule{
				validation.RequiredRule(),
				validation.MaxLengthRule(100),
			},
		},
		{
			Name:  "email",
			Value: &form.Email,
			Rules: []validation.Rule{
				validation.RequiredRule(),
				validation.EmailRule(),
				validation.MaxLengthRule(255),
			},
		},
		{
			Name:  "subject",
			Value: &form.Subject,
			Rules: []validation.Rule{
				validation.RequiredRule(),
				validation.MaxLengthRule(300),
			},
		},
		{
			Name:  "message",
			Value: &form.Message,
			Rules: []validation.Rule{
				validation.RequiredRule(),
			},
		},
		{
			Name:  "product_name",
			Value: &form.ProductName,
			Rules: []validation.Rule{
				validation.ProductNameRule(),
			},
		},
		{
			Name:  "product_website",
			Value: &form.ProductWebsite,
			Rules: []validation.Rule{
				validation.UrlRule(),
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
