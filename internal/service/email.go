package service

import (
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/model"
	"Form-Mailly-Go/internal/template"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
)

type EmailService interface {
	Send(form *model.ContactForm) error
	SendBatchEmail(form *model.ContactForm) error
}

type SMTPEmailService struct {
	config *config.EnvironmentVariable
	client *smtp.Client
	auth   smtp.Auth
}

func NewSMTPEmailService(cfg *config.EnvironmentVariable) *SMTPEmailService {
	return &SMTPEmailService{config: cfg}
}

func (s *SMTPEmailService) Send(form *model.ContactForm) error {
	if s.config.SenderEmail == "" || s.config.SenderPassword == "" {
		return errors.New("email configuration not initialized")
	}

	auth := smtp.PlainAuth("", s.config.SenderEmail, s.config.SenderPassword, s.config.SMTPHost)
	to := []string{s.config.ReceiverEmail}
	msg := []byte(
		"From: " + form.ProductName + " <" + s.config.SenderEmail + ">\r\n" +
			"To: " + to[0] + "\r\n" +
			"Subject: " + form.Subject + "\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			template.BuildContactFormMessage2(form),
	)

	return smtp.SendMail(
		s.config.SMTPHost+":"+s.config.SMTPPort,
		auth,
		s.config.SenderEmail,
		to,
		msg,
	)
}

// Connect Initialize and reuse SMTP connection
func (s *SMTPEmailService) Connect() error {
	if s.config.SenderEmail == "" || s.config.SenderPassword == "" {
		return errors.New("email configuration not initialized")
	}

	addr := s.config.SMTPHost + ":" + s.config.SMTPPort

	// 1️⃣ TCP connect
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP: %v", err)
	}

	// 2️⃣ Create SMTP client
	c, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}

	// 3️⃣ STARTTLS upgrade
	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: s.config.SMTPHost}
		if err = c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %v", err)
		}
	}

	// 4️⃣ Authenticate
	auth := smtp.PlainAuth("", s.config.SenderEmail, s.config.SenderPassword, s.config.SMTPHost)
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	s.client = c
	s.auth = auth
	return nil
}

// SendBatchEmail email using the existing connection
func (s *SMTPEmailService) SendBatchEmail(form *model.ContactForm) error {
	if s.client == nil {
		if err := s.Connect(); err != nil {
			return err
		}
	}

	to := []string{s.config.ReceiverEmail}

	if err := s.client.Mail(s.config.SenderEmail); err != nil {
		return err
	}
	if err := s.client.Rcpt(to[0]); err != nil {
		return err
	}

	w, err := s.client.Data()
	if err != nil {
		return err
	}

	msg := []byte(
		"From: " + form.ProductName + " <" + s.config.SenderEmail + ">\r\n" +
			"To: " + to[0] + "\r\n" +
			"Subject: " + form.Subject + "\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			form.Message,
	)

	if _, err = w.Write(msg); err != nil {
		return err
	}

	return w.Close()
}

// Close the connection when done
func (s *SMTPEmailService) Close() {
	if s.client != nil {
		_ = s.client.Quit()
	}
}

//
//func (s *SMTPEmailService) SendBatchEmail(form *model.ContactForm) error {
//	if s.config.SenderEmail == "" || s.config.SenderPassword == "" {
//		return errors.New("email configuration not initialized")
//	}
//
//	auth := smtp.PlainAuth("", s.config.SenderEmail, s.config.SenderPassword, s.config.SMTPHost)
//	to := []string{s.config.ReceiverEmail}
//	msg := []byte(
//		"From: " + form.ProductName + " <" + s.config.SenderEmail + ">\r\n" +
//			"To: " + to[0] + "\r\n" +
//			"Subject: " + form.Subject + "\r\n" +
//			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
//			"\r\n" +
//			form.Message,
//	)
//
//	return smtp.SendMail(
//		s.config.SMTPHost+":"+s.config.SMTPPort,
//		auth,
//		s.config.SenderEmail,
//		to,
//		msg,
//	)
//}
