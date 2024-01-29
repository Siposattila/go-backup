package node

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Siposattila/gobkup/internal/alert"
	"github.com/Siposattila/gobkup/internal/certification"
	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

type Node struct {
	Dialer       webtransport.Dialer
	ServerStream webtransport.Stream
	Config       config.NodeConfig
	DiscordAlert alert.AlertInterface
	EmailAlert   alert.AlertInterface
	initOnce     sync.Once
}

func (node *Node) init(token string, debug bool) {
    node.Config = config.NodeConfig{
        Token: token,
        Debug: debug,
    }
	node.Dialer.RoundTripper = &http3.RoundTripper{}
	node.getTlsConfig()
}

func (node *Node) Run(endpoint string, token string, debug bool) {
	console.Normal("Setting up and starting node server...")
	node.initOnce.Do(func() { node.init(token, debug) })

	if node.Config.Debug {
		console.Debug("Debug mode is active!")
	}

	node.listenForKill()
	var response, connection, error = node.Dialer.Dial(context.Background(), endpoint, nil)
	if error != nil {
		console.Fatal("Unable to connect to master: " + error.Error())
	}

	if response.StatusCode < 200 && response.StatusCode >= 300 {
		console.Fatal("The response status code was not 2xx the error is: " + error.Error())
	}

	var stream, streamError = connection.OpenStream()
	if streamError != nil {
		console.Fatal("There was an error during opening the stream: " + streamError.Error())
	}
	node.ServerStream = stream

	console.Success("Node is up and running! Ready to communicate with the master!")
	node.handleStream()

	return
}

func (node *Node) getClientName() string {
	var name, nameError = os.Hostname()
	if nameError != nil {
		console.Fatal("Can't get client name. This means there is no hostname??")
	}

	return name
}

func (node *Node) Close() {
    console.Normal("Shutting down node client...")
	if node.ServerStream != nil {
		node.ServerStream.Close()
	}
	node.Dialer.Close()

	return
}

func (node *Node) listenForKill() {
	var channel = make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-channel
		node.Close()
		os.Exit(1)
	}()

	return
}

func (node *Node) getTlsConfig() {
	var ca, _, caError = certification.GenerateCA()
	if caError != nil {
		console.Fatal("Unable to generate CA certificate: " + caError.Error())
	}

	var certPool = x509.NewCertPool()
	certPool.AddCert(ca)
	var tlsConfig = &tls.Config{RootCAs: certPool}
	if node.Config.Debug {
		tlsConfig.InsecureSkipVerify = true
	}
	node.Dialer.RoundTripper.TLSClientConfig = tlsConfig

	return
}
