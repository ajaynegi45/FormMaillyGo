package service

import (
	"Form-Mailly-Go/internal/config"
	"Form-Mailly-Go/internal/model"
	"Form-Mailly-Go/internal/template"
	"net/smtp"
)

func Send(form *model.ContactForm) error {

	auth := smtp.PlainAuth("", config.EnvVar.SenderEmail, config.EnvVar.SenderPassword, config.EnvVar.SMTPHost)
	to := []string{config.EnvVar.ReceiverEmail}
	msg := []byte(
		"From: " + form.ProductName + " <" + config.EnvVar.SenderEmail + ">\r\n" +
			"To: " + to[0] + "\r\n" +
			"Subject: " + form.Subject + "\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"MIME-Version: 1.0\r\n" +
			"\r\n" +
			template.BuildContactFormMessage2(form),
	)

	return smtp.SendMail(
		config.EnvVar.SMTPHost+":"+config.EnvVar.SMTPPort,
		auth,
		config.EnvVar.SenderEmail,
		to,
		msg,
	)
}
