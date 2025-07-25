package internal

import (
	"Form-Mailly-Go/internal/template"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

func SendEmail(contactForm *ContactForm, response http.ResponseWriter) {

	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	// Sender data
	from := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("SENDER_EMAIL_PASSWORD")

	// Receiver email address
	to := []string{os.Getenv("RECEIVER_EMAIL")}

	// SMTP server configuration
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Create full message with headers
	message := []byte(
		"From: " + contactForm.ProductName + " <" + from + ">\r\n" +
			"To: " + to[0] + "\r\n" +
			"Subject: " + contactForm.Subject + "\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			template.ContactFormSubmissionTemplateBuilder(contactForm.Name, contactForm.Email, contactForm.Subject, contactForm.Message, contactForm.ProductName, contactForm.ProductWebsite),
	)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}

	response.WriteHeader(http.StatusCreated)
	_, err2 := response.Write([]byte("Email Sent Successfully!"))
	if err2 != nil {
		return
	}
}
