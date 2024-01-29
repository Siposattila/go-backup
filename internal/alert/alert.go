package alert

type AlertInterface interface {
	Start()
	Close()
	Send(message string)
}
