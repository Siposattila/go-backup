package client

import (
	"context"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/Siposattila/go-backup/backup"
	"github.com/Siposattila/go-backup/config"
	"github.com/Siposattila/go-backup/log"
	"github.com/Siposattila/go-backup/request"
	"github.com/Siposattila/go-backup/serializer"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/webtransport-go"
)

type Client interface {
	Start(wg *sync.WaitGroup)
	Stop()
}

type client struct {
	Dialer               webtransport.Dialer
	Stream               webtransport.Stream
	Config               *config.Client
	BackupConfig         *config.Backup
	Backup               backup.BackupInterface
	newBackupPathChannel chan string
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
	clientWg.Add(1)

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
	go c.handleStream()
}

func (c *client) Stop() {
	log.GetLogger().Normal("Stopping client...")
	c.Stream.Close()
	c.Dialer.Close()
	c.Backup.Stop()
}

func (c *client) handleStream() {
	log.GetLogger().Normal("Trying to request backup config from server...")
	request.Write(c.Stream, request.NewRequest(c.Config.ClientId, request.ID_CONFIG, ""))

	for {
		r := request.Response{}
		_, readError := request.Read(c.Stream, &r)
		if readError != nil {
			log.GetLogger().Fatal("Read error occured during stream handling. Server error occured!", readError.Error())
		}

		switch r.Id {
		case request.ID_CONFIG:
			config := config.Backup{}
			serializerError := serializer.Json.Serialize([]byte(r.Data), &config)
			if serializerError != nil {
				log.GetLogger().Fatal("Error occured during getting backup config.", serializerError.Error())
			}
			c.BackupConfig = &config

			log.GetLogger().Success("Got backup config from server!")
			c.startBackup()
		case request.ID_BACKUP_CHUNK_PROCESSED:
			chunk := Chunk{}
			serializerError := serializer.Json.Serialize([]byte(r.Data), &chunk)
			if serializerError != nil {
				log.GetLogger().Fatal("Error occured during getting chunk info from server.", serializerError.Error())
			} else {
				if err := os.Remove(path.Join(CHUNK_TEMP_DIR, chunk.ChunkName)); err != nil {
					log.GetLogger().Error(err.Error())
				}
			}
		}
	}
}

func (c *client) startBackup() {
	c.Backup = backup.NewBackup(
		c.BackupConfig.WhenToBackup,
		&c.BackupConfig.WhatToBackup,
		&c.BackupConfig.Exclude,
	)

	c.newBackupPathChannel = make(chan string)
	go c.Backup.Backup(c.newBackupPathChannel)
	go c.handNewBackup()
}
