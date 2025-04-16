package client

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Siposattila/gobkup/backup"
	"github.com/Siposattila/gobkup/config"
	"github.com/Siposattila/gobkup/log"
	"github.com/Siposattila/gobkup/request"
	"github.com/Siposattila/gobkup/serializer"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/webtransport-go"
)

type Client interface {
	Start(wg *sync.WaitGroup)
	Stop()
}

type client struct {
	Dialer       webtransport.Dialer
	Stream       webtransport.Stream
	Config       *config.Client
	BackupConfig *config.Backup
	Backup       backup.BackupInterface
}

func NewClient() Client {
	var c client
	c.Config = c.Config.Get()

	c.Dialer = webtransport.Dialer{QUICConfig: &quic.Config{
		KeepAlivePeriod: time.Duration(25 * time.Second),
		EnableDatagrams: true,
	}}
	c.getTlsConfig()

	return &c
}

func (c *client) Start(clientWg *sync.WaitGroup) {
	log.GetLogger().Normal("Starting up client...")

	h := http.Header{}
	h.Add("Authorization", "Basic "+c.Config.Token)
	res, conn, err := c.Dialer.Dial(context.Background(), c.Config.Endpoint, h)
	if err != nil {
		log.GetLogger().Fatal("Unable to connect to server.", err.Error())
	}

	if res.StatusCode < 200 && res.StatusCode >= 300 {
		log.GetLogger().Fatal("The response status code was not 2xx.", err.Error())
	}

	stream, streamError := conn.OpenStream()
	if streamError != nil {
		log.GetLogger().Fatal("There was an error during opening the stream.", streamError.Error())
	}
	c.Stream = stream

	log.GetLogger().Success("Client is up and running! Ready to communicate with the server!")
	c.handleStream()
}

func (c *client) Stop() {
	log.GetLogger().Normal("Stopping client...")
	c.Stream.Close()
	c.Dialer.Close()
	c.Backup.Stop()
}

func (c *client) handleStream() {
	log.GetLogger().Normal("Trying to request backup config from server...")
	request.Write(c.Stream, request.NewRequest(c.Config.ClientId, request.REQUEST_ID_CONFIG, ""))

	for {
		resp := request.Response{}
		n, readError := request.Read(c.Stream, &resp)
		if readError != nil {
			log.GetLogger().Error("Read error occured during stream handling.", readError.Error())
			break
		}

		log.GetLogger().Debug("Read length: ", n)

		switch resp.Id {
		case request.REQUEST_ID_CONFIG:
			config := config.Backup{}
			serializerError := serializer.Json.Serialize([]byte(resp.Data), &config)
			if serializerError != nil {
				log.GetLogger().Fatal("Error occured during getting backup config.", serializerError.Error())
			}
			c.BackupConfig = &config
			log.GetLogger().Success("Got backup config from server!")
			c.startBackup()
		}
	}
}

func (c *client) startBackup() {
	c.Backup = backup.NewBackup(
		c.BackupConfig.WhenToBackup,
		&c.BackupConfig.WhatToBackup,
		&c.BackupConfig.Exclude,
	)

	c.Backup.Backup()
}
