package service

import (
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/model"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

func SetupNewSMTPConnection() (*smtp.Client, error) {
	addr := config.EnvVar.SMTPHost + ":" + config.EnvVar.SMTPPort

	// 1️⃣ TCP connect
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMTP: %v", err)
	}

	// 2️⃣ Create SMTP client
	client, err := smtp.NewClient(conn, config.EnvVar.SMTPHost)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP client: %v", err)
	}

	// 3️⃣ STARTTLS upgrade
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: config.EnvVar.SMTPHost}
		if err = client.StartTLS(tlsConfig); err != nil {
			return nil, fmt.Errorf("failed to start TLS: %v", err)
		}
	}

	// 4️⃣ Authenticate
	auth := smtp.PlainAuth("", config.EnvVar.SenderEmail, config.EnvVar.SenderPassword, config.EnvVar.SMTPHost)
	if err = client.Auth(auth); err != nil {
		return nil, fmt.Errorf("failed to authenticate: %v", err)
	}
	return client, nil
}

func SendEmailUsingWorker(client *smtp.Client, email *model.Email) error {
	if client == nil {
		return fmt.Errorf("client is nil")
	}

	// Sets up the recipient list.
	to := email.SentTo

	// Sets the sender service address in the SMTP protocol using MAIL FROM:<sender>.
	if err := client.Mail(config.EnvVar.SenderEmail); err != nil { // Starts new service (MAIL FROM)
		return err
	}

	// Adds the recipient to the envelope using RCPT TO:<recipient>.
	if err := client.Rcpt(to); err != nil { // Adds recipient (RCPT TO)
		return err
	}

	// Opens the data stream to start sending the service body.
	writer, err := client.Data() // Prepares to send service content
	if err != nil {
		return err
	}

	// Composes the service message with headers and the body.
	msg := []byte(
		"From: " + email.ProductName + " <" + config.EnvVar.SenderEmail + ">\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + email.Subject + "\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			email.Message,
	)

	// Writes the message content to the SMTP data stream.
	if _, err = writer.Write(msg); err != nil { // Sends the body
		return err
	}

	// It tells the SMTP server that the message is complete. Without this, the service won't be sent. If you don’t close the writer, the SMTP server won’t process or deliver the message.
	return writer.Close() // Ends the service
}

// CloseSMTPConnection Close the connection when done
func CloseSMTPConnection(client *smtp.Client) {
	if client != nil {
		client.Quit()
	}
}
