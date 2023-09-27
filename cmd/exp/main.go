package main

import (
	"github.com/dmcclung/pixelparade/models"
)

func main() {
	from := "dylan@pixelparade"
	to := "guest@pixelparade"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`

	email := models.Email{
		To:        to,
		From:      from,
		Subject:   subject,
		Plaintext: plaintext,
		HTML:      html,
	}

	config, err := models.GetEmailConfig()
	if err != nil {
		panic(err)
	}

	emailService, err := models.GetEmailService(config)
	if err != nil {
		panic(err)
	}

	err = emailService.SendEmail(email)
	if err != nil {
		panic(err)
	}

	err = emailService.SendResetEmail("test@pixelparade", "https://google.com")
	if err != nil {
		panic(err)
	}
}
