package alert

type AlertInterface interface {
	Start()
	Stop()
	Send(message string)
}
