package main

import (
	"os"

	"github.com/dmcclung/pixelparade/models"
	"github.com/go-mail/mail/v2"
)

func main() {
	from := "dylan@pixelparade"
	to := "guest@pixelparade"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`

	m := mail.NewMessage()
	m.SetHeader("To", to)
	m.SetHeader("From", from)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", plaintext)
	m.AddAlternative("text/html", html)
	m.WriteTo(os.Stdout)

	emailService, err := models.GetEmailService()
	if err != nil {
		panic(err)
	}

	err = emailService.SendEmail(m)
	if err != nil {
		panic(err)
	}
}
