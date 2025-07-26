package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type EmailConfig struct {
	SenderEmail    string
	SenderPassword string
	ReceiverEmail  string
	SMTPHost       string
	SMTPPort       string
}

func LoadEmailConfig() *EmailConfig {
	envErr := godotenv.Load(".env.dev")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	return &EmailConfig{
		SenderEmail:    os.Getenv("SENDER_EMAIL"),
		SenderPassword: os.Getenv("SENDER_EMAIL_PASSWORD"),
		ReceiverEmail:  os.Getenv("RECEIVER_EMAIL"),
		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       os.Getenv("SMTP_PORT"),
	}
}
