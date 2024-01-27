package alert

import (
	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"gopkg.in/gomail.v2"
)

var emailClient *gomail.Dialer

func InitEmailClient() {
	console.Normal("Starting email client...")
	config.LoadConfig("email")
	console.Normal("Email client version: gomail")
	emailClient = gomail.NewDialer(config.Email.Host, config.Email.Port, config.Email.User, config.Email.Password)
	console.Success("Email client started! Ready to send alerts!")

	return
}

func CloseEmailClient() {
	if emailClient != nil {
		emailClient = nil
	}
	console.Normal("Getting rid of the email client...")

	return
}

func SendEmailAlert(subject string, message string) {
	var email = gomail.NewMessage()
	email.SetHeader("From", config.Email.EmailSender)
	email.SetHeader("To", config.Email.EmailReceiver)
	email.SetHeader("Subject", subject)
	email.SetBody("text/html", message)

	var emailError = emailClient.DialAndSend(email)
	if emailError != nil {
		console.Error("There was an error during the dialing or sending of the email alert: " + emailError.Error())
	}

	return
}
