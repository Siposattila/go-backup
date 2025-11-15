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
	"github.com/Siposattila/go-backup/generatedproto"
	"github.com/Siposattila/go-backup/log"
	"github.com/Siposattila/go-backup/request"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/webtransport-go"
	"google.golang.org/protobuf/proto"
)

type Client interface {
	Start(wg *sync.WaitGroup)
	Stop()
}

type client struct {
	Dialer               webtransport.Dialer
	Stream               webtransport.Stream
	Config               *generatedproto.Client
	BackupConfig         *generatedproto.Backup
	Backup               backup.BackupInterface
	newBackupPathChannel chan string
}

func NewClient() Client {
	var c client
	c.Config = config.GetClientConfig()

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
	if err := c.Stream.Close(); err != nil {
		log.GetLogger().Error(err.Error())
	}
	if err := c.Dialer.Close(); err != nil {
		log.GetLogger().Error(err.Error())
	}
	c.Backup.Stop()
}

func (c *client) handleStream() {
	log.GetLogger().Normal("Trying to request backup config from server...")
	if _, err := request.Write(c.Stream, request.NewRequest(c.Config.ClientId, generatedproto.RequestType_ID_CONFIG, &generatedproto.Backup{})); err != nil {
		log.GetLogger().Fatal("Failed writing config getter request to stream: ", err.Error())
	}

	for {
		r := generatedproto.Response{}
		_, readError := request.Read(c.Stream, &r)
		if readError != nil {
			log.GetLogger().Fatal("Read error occured during stream handling. Server error occured!", readError.Error())
		}

		switch r.Id {
		case generatedproto.RequestType_ID_CONFIG:
			protoHelper := &generatedproto.Backup{}
			serializerError := proto.Unmarshal([]byte(r.Data), protoHelper)
			if serializerError != nil {
				log.GetLogger().Fatal("Error occured during getting backup config.", serializerError.Error())
			}
			c.BackupConfig = &generatedproto.Backup{
				WhenToBackup: protoHelper.WhenToBackup,
				WhatToBackup: protoHelper.WhatToBackup,
				Exclude:      protoHelper.Exclude,
			}

			log.GetLogger().Success("Got backup config from server!")
			c.startBackup()
		case generatedproto.RequestType_ID_BACKUP_CHUNK_PROCESSED:
			chunk := &generatedproto.Chunk{}
			serializerError := proto.Unmarshal([]byte(r.Data), chunk)
			if serializerError != nil {
				log.GetLogger().Fatal("Error occured during getting chunk info from server.", serializerError.Error())
			} else {
				if err := os.Remove(path.Join(CHUNK_TEMP_DIR, chunk.ChunkName)); err != nil {
					log.GetLogger().Error("Failed to remove temp dir for chunks: ", err.Error())
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
