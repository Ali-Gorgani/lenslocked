package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	service := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
	return &service
}

type EmailService struct {
	// DefaultSender is the default sender of emails when no sender is provided.
	DefaultSender string

	// unexported fields
	dialer *mail.Dialer
}

func (service *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	service.setFrom(msg, email)
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}

	err := service.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (service *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		To:      to,
		Subject: "Lenslocked Password Reset",
		HTML:    fmt.Sprintf("Click <a href=\"%s\">here</a> to reset your password.", resetURL),
	}
	err := service.Send(email)
	if err != nil {
		return fmt.Errorf("forgotPassword: %w", err)
	}
	return nil
}

func (service *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case service.DefaultSender != "":
		from = service.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}
