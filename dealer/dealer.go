package dealer

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Siposattila/gobkup/client"
	"github.com/Siposattila/gobkup/server"
)

type dealer interface {
	Start(wg *sync.WaitGroup)
	Stop()
}

func Run(isServer, isClient bool) {
	var dealer dealer
	if isServer && !isClient {
		dealer = server.NewServer()
	} else {
		dealer = client.NewClient()
	}

	var wg sync.WaitGroup
	dealer.Start(&wg)

	// Setup logic for stopping dealer
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		dealer.Stop()
		os.Exit(0)
	}()

	wg.Wait()
}
