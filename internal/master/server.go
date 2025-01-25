package master

import (
	"crypto/tls"
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

type Master struct {
	Server       webtransport.Server
	Config       config.MasterConfig
	DiscordAlert alert.AlertInterface
	EmailAlert   alert.AlertInterface
	initOnce     sync.Once
}

func (master *Master) init(debug bool) {
	master.Config = config.LoadMasterConfig()
	master.Config.Debug = debug
	master.Server = webtransport.Server{
		H3: http3.Server{Addr: master.Config.Port},
	}
	master.getTlsConfig()

	if master.Config.DiscordAlert {
		master.DiscordAlert = alert.NewDiscord()
		master.DiscordAlert.Start()
	}

	if master.Config.EmailAlert {
		master.EmailAlert = alert.NewEmail()
		master.EmailAlert.Start()
	}
}

func (master *Master) Run(debug bool) {
	console.Normal("Setting up and starting master server...")
	master.initOnce.Do(func() { master.init(debug) })

	if master.Config.Debug {
		console.Debug("Debug mode is active!")
	}

	master.setupEndpoint()
	master.listenForKill()
	console.Success("Master server is up and running! Ready to handle connections on port " + master.Config.Port)
	serverError := master.Server.ListenAndServe()
	if serverError != nil {
		console.Fatal("Error during server listen and serve: " + serverError.Error())
	}
}

func (master *Master) Close() {
	console.Normal("Shutting down master server...")
	master.Server.Close()
	if master.Config.DiscordAlert {
		master.DiscordAlert.Close()
	}

	if master.Config.EmailAlert {
		master.EmailAlert.Close()
	}
}

func (master *Master) listenForKill() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-channel
		master.Close()
		os.Exit(1)
	}()
}

func (master *Master) getTlsConfig() {
	ca, caPrivateKey, caError := certification.GenerateCA()
	if caError != nil {
		console.Fatal("Unable to generate CA certificate: " + caError.Error())
	}

	leafCert, leafPrivateKey, leafError := certification.GenerateLeafCert(master.Config.Domain, ca, caPrivateKey)
	if leafError != nil {
		console.Fatal("Unable to generate leaf certificate: " + leafError.Error())
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{leafCert.Raw},
			PrivateKey:  leafPrivateKey,
		}},
		NextProtos: []string{"webtransport-go / quic-go"},
	}

	if master.Config.Debug {
		tlsConfig.InsecureSkipVerify = true
	}
	master.Server.H3.TLSConfig = tlsConfig
	console.Success("Tls config was obtained successfully!")
}
