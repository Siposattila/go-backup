package alert

type Alert interface {
	Start()
	Close()
	Send(message string)
}
