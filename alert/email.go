package alert

import (
	"github.com/Siposattila/go-backup/log"
	"gopkg.in/gomail.v2"
)

type email struct {
	Client   *gomail.Dialer
	Receiver string
	Sender   string
	User     string
	Password string
	Host     string
	Port     int
}

func NewEmail(
	receiver string,
	sender string,
	user string,
	password string,
	host string,
	port int) AlertInterface {
	return &email{
		Receiver: receiver,
		Sender:   sender,
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (e *email) Start() {
	log.GetLogger().Normal("Starting email client...")
	log.GetLogger().Normal("Email client version: (gomail)v2")
	e.Client = gomail.NewDialer(e.Host, e.Port, e.User, e.Password)
	log.GetLogger().Success("Email client started! Ready to send alerts!")
}

func (e *email) Stop() {
	if e.Client != nil {
		e.Client = nil
	}
	log.GetLogger().Normal("Stopping email client...")
}

func (e *email) Send(message string) {
	var mail = gomail.NewMessage()
	mail.SetHeader("From", e.Sender)
	mail.SetHeader("To", e.Receiver)
	mail.SetHeader("Subject", "Gobkup email alert")
	mail.SetBody("text/html", message)

	if emailError := e.Client.DialAndSend(mail); emailError != nil {
		log.GetLogger().Error("There was an error during the dialing or sending of the email alert: " + emailError.Error())
	}
}
