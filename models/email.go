package models

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-mail/mail/v2"
	"github.com/joho/godotenv"
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
	To string
	From string
	Subject string
	Plaintext string
	HTML string
}

func GetEmailConfig() (*SMTPConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("loading env: %w", err)
	}

	envVars := []string{"EMAIL_HOST", "EMAIL_PORT", "EMAIL_USER", "EMAIL_PASS"}
	args := []string{}
	for _, envVar := range envVars {
		val, exists := os.LookupEnv(envVar)
		if !exists {
			return nil, fmt.Errorf("%v not found", envVar)
		}
		args = append(args, val)
	}

	p, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, fmt.Errorf("converting port: %w", err)
	}

	config := SMTPConfig{args[0], p, args[2], args[3]}

	return &config, nil
}

type EmailService struct {
	dialer *mail.Dialer
}

func GetEmailService() (*EmailService, error) {
	config, err := GetEmailConfig()
	if err != nil {
		return nil, fmt.Errorf("email config: %w", err)
	}

	dialer := mail.NewDialer(config.Host, config.Port, config.User, config.Pass)

	return &EmailService{
		dialer: dialer,
	}, nil
}

func (es *EmailService) SendEmail(email Email) error {
	m := mail.NewMessage()
	m.SetHeader("To", email.To)
	m.SetHeader("From", email.From)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", email.Plaintext)
	m.AddAlternative("text/html", email.HTML)

	err := es.dialer.DialAndSend(m)
	return err
}
