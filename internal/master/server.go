package master

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Siposattila/gobkup/internal/alert"
	"github.com/Siposattila/gobkup/internal/certification"
	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

var server webtransport.Server

func SetupAndRunServer() {
	console.Normal("Setting up and starting master server...")
	listenForKill()
	server = webtransport.Server{
		H3: http3.Server{Addr: config.Master.Port},
	}

	certification.GetServerTlsConfig()
	if config.Master.Debug {
		console.Debug("Debug mode is active!")
		certification.TlsConfig.InsecureSkipVerify = true
	}
	server.H3.TLSConfig = certification.TlsConfig

	setupAlertSystem()

	http.HandleFunc("/master", func(writer http.ResponseWriter, request *http.Request) {
		var connection, error = server.Upgrade(writer, request)
		if error != nil {
			console.Error("Upgrading failed: " + error.Error())
			writer.WriteHeader(500)

			return
		}

		var stream, streamError = connection.AcceptStream(context.Background())
		if streamError != nil {
			console.Error("There was an error during accepting the stream: " + streamError.Error())
		}

		go handleStream(stream)
	})

	console.Success("Master server is up and running! Ready to handle connections on port " + config.Master.Port)
	var serverError = server.ListenAndServe()
	if serverError != nil {
		console.Fatal("Error during server listen and serve: " + serverError.Error())
	}

	return
}

func setupAlertSystem() {
	if config.Master.DiscordAlert {
		alert.InitDiscordClient()
	}

	if config.Master.EmailAlert {
		alert.InitEmailClient()
	}
}

func CloseServer() {
	server.Close()
	console.Normal("Shutting down master server...")

	return
}

func listenForKill() {
	var channel = make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-channel
		CloseServer()
		if config.Master.DiscordAlert {
			alert.CloseDiscordClient()
		}

		if config.Master.EmailAlert {
			alert.CloseEmailClient()
		}
		os.Exit(1)
	}()

	return
}
