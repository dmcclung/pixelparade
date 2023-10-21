package models

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@pixelparade"
)

type SMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type Email struct {
	To        string
	From      string
	Subject   string
	Plaintext string
	HTML      string
}

func GetEmailConfig() (SMTPConfig, error) {
	var config SMTPConfig

	host, exists := os.LookupEnv("EMAIL_HOST")
	if !exists {
		return config, fmt.Errorf("EMAIL_HOST environment var not found")
	}
	config.Host = host

	portStr, exists := os.LookupEnv("EMAIL_PORT")
	if !exists {
		return config, fmt.Errorf("EMAIL_PORT environment var not found")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return config, fmt.Errorf("converting port: %w", err)
	}
	config.Port = port

	user, exists := os.LookupEnv("EMAIL_USER")
	if !exists {
		return config, fmt.Errorf("EMAIL_USER environment var not found")
	}
	config.User = user

	pass, exists := os.LookupEnv("EMAIL_PASSWORD")
	if !exists {
		return config, fmt.Errorf("EMAIL_PASSWORD environment var not found")
	}
	config.Pass = pass

	return config, nil
}

type EmailService struct {
	dialer        *mail.Dialer
	DefaultSender string
}

func GetEmailService(config SMTPConfig) (*EmailService, error) {
	dialer := mail.NewDialer(config.Host, config.Port, config.User, config.Pass)

	return &EmailService{
		dialer: dialer,
	}, nil
}

func (es *EmailService) SendResetEmail(to, resetLink string) error {
	email := Email{
		To:        to,
		Subject:   "Password Reset",
		Plaintext: fmt.Sprintf("To reset your password visit this link %v", resetLink),
		HTML:      fmt.Sprintf("<p>To reset your password visit this <a href=\"%v\">link</a></p>", resetLink),
	}

	err := es.SendEmail(email)
	if err != nil {
		return fmt.Errorf("reset email: %w", err)
	}
	return nil
}

func (es *EmailService) SendEmail(email Email) error {
	m := mail.NewMessage()
	m.SetHeader("To", email.To)

	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	m.SetHeader("From", from)

	m.SetHeader("Subject", email.Subject)
	switch {
	case email.HTML != "" && email.Plaintext != "":
		m.SetBody("text/html", email.HTML)
		m.AddAlternative("text/plain", email.Plaintext)
	case email.HTML == "" && email.Plaintext != "":
		m.SetBody("text/plain", email.Plaintext)
	case email.HTML != "" && email.Plaintext == "":
		m.SetBody("text/html", email.HTML)
	}

	err := es.dialer.DialAndSend(m)
	return err
}
