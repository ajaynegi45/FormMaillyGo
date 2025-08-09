package service

import (
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/model"
	"Form-Mailly-Go/internal/template"
	"errors"
	"net/smtp"
)

type EmailService interface {
	Send(form *model.ContactForm) error
}

type SMTPEmailService struct {
	config *config.EnvironmentVariable
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
