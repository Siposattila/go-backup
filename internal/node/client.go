package node

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Siposattila/gobkup/internal/certification"
	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

var dialer webtransport.Dialer
var serverStream webtransport.Stream

func SetupAndRunClient(endpoint string) {
	console.Normal("Setting up and starting node server...")
	listenForKill()
	certification.GetClientTlsConfig()
	if config.Node.Debug {
		console.Debug("Debug mode is active!")
		certification.TlsConfig.InsecureSkipVerify = true
	}

	dialer.RoundTripper = &http3.RoundTripper{
		TLSClientConfig: certification.TlsConfig,
	}

	var ctx, _ = context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	var response, connection, error = dialer.Dial(ctx, endpoint, nil)

	if error != nil {
		console.Fatal("Unable to connect to master: " + error.Error())
	}

	if response.StatusCode < 200 && response.StatusCode >= 300 {
		console.Fatal("The response status code was not 2xx the error is: " + error.Error())
	}

	var stream, streamError = connection.OpenStream()
	if streamError != nil {
		console.Error("There was an error during opening the stream: " + streamError.Error())
	}
	serverStream = stream

	console.Success("Node is up and running! Ready to communicate with the master!")
	handleStream()

	return
}

func getClientName() string {
	var name, nameError = os.Hostname()
	if nameError != nil {
		console.Fatal("Can't get client name. This means there is no hostname??")
	}

	return name
}

func CloseClient() {
	dialer.Close()
	console.Normal("Shutting down node client...")

	return
}

func listenForKill() {
	var channel = make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-channel
		if serverStream != nil {
			serverStream.Close()
		}
		CloseClient()
		os.Exit(1)
	}()

	return
}
