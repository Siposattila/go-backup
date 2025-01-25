package alert

import (
	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"gopkg.in/gomail.v2"
)

type email struct {
	Client *gomail.Dialer
	Config config.EmailConfig
}

func NewEmail() Alert {
	return email{}
}

func (email email) Start() {
	console.Normal("Starting email client...")
	email.Config = config.LoadEmailAlertConfig()
	console.Normal("Email client version: gomail")
	email.Client = gomail.NewDialer(email.Config.Host, email.Config.Port, email.Config.User, email.Config.Password)
	console.Success("Email client started! Ready to send alerts!")
}

func (email email) Close() {
	if email.Client != nil {
		email.Client = nil
	}
	console.Normal("Getting rid of the email client...")
}

func (email email) Send(message string) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", email.Config.Sender)
	mail.SetHeader("To", email.Config.Receiver)
	mail.SetHeader("Subject", "Gobkup email alert")
	mail.SetBody("text/html", message)

	emailError := email.Client.DialAndSend(mail)
	if emailError != nil {
		console.Error("There was an error during the dialing or sending of the email alert: " + emailError.Error())
	}
}
