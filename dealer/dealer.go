package dealer

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Siposattila/gobkup/client"
	"github.com/Siposattila/gobkup/server"
)

type dealer interface {
	Start()
	Stop()
}

func Run(isServer, isClient bool) {
	var dealer dealer
	if isServer && !isClient {
		dealer = server.NewServer()
	} else {
		dealer = client.NewClient()
	}

	dealer.Start()

	// Setup logic for stopping dealer
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		dealer.Stop()
		os.Exit(1)
	}()
}
