package config

import (
	"Form-Mailly-Go/internal/validation"
	"log"
	"os"
)

// EnvironmentVariable holds all configuration needed for service sending
type EnvironmentVariable struct {
	SenderEmail    string // Gmail address used to send emails
	SenderPassword string // Gmail app password
	ReceiverEmail  string // Default recipient service (can be overridden in batch)
	SMTPHost       string // SMTP server hostname
	SMTPPort       string // SMTP server port
}

var EnvVar *EnvironmentVariable

// LoadEnvironmentVariable reads configuration from environment variables
// This function should be called once during application initialization
func LoadEnvironmentVariable() {

	EnvVar = &EnvironmentVariable{
		// Required: Gmail credentials
		SenderEmail:    os.Getenv("SENDER_EMAIL"),
		SenderPassword: os.Getenv("SENDER_EMAIL_PASSWORD"),

		// Optional: Default receiver (can be overridden in batch)
		ReceiverEmail: os.Getenv("RECEIVER_EMAIL"),

		// SMTP configuration with sensible defaults for Gmail
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),
	}

	if !EnvVar.IsValid() {
		log.Println("❌ Invalid environment configuration")
		os.Exit(1)
	}
	log.Println("✅ Environment configuration loaded successfully.")
}

// IsValid checks if the configuration has all required fields
// This is useful for validating configuration during startup
func (env *EnvironmentVariable) IsValid() bool {
	validator := validation.NewValidator()

	fields := []validation.Field{
		{
			Name:  "SENDER_EMAIL",
			Value: &env.SenderEmail,
			Rules: []validation.Rule{
				validation.RequiredRule(),
				validation.EmailRule(),
			},
		}, {
			Name:  "SENDER_EMAIL_PASSWORD",
			Value: &env.SenderPassword,
			Rules: []validation.Rule{
				validation.RequiredRule(),
			},
		}, {
			Name:  "RECEIVER_EMAIL",
			Value: &env.ReceiverEmail,
			Rules: []validation.Rule{
				validation.EmailRule(),
			},
		},
		{
			Name:  "SMTP_HOST",
			Value: &env.SMTPHost,
			Rules: []validation.Rule{
				validation.RequiredRule(),
			},
		},
		{
			Name:  "SMTP_PORT",
			Value: &env.SMTPPort,
			Rules: []validation.Rule{
				validation.RequiredRule(),
				validation.NumericRule(),
			},
		},
	}

	for _, field := range fields {
		validator.ValidateField(field)
		if !validator.IsValid() {
			// Return immediately once an error occurs
			return false
		}
	}

	return true
}
