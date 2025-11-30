package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/Siposattila/go-backup/backup"
	"github.com/Siposattila/go-backup/config"
	"github.com/Siposattila/go-backup/log"
	"github.com/Siposattila/go-backup/proto"
	"github.com/Siposattila/go-backup/request"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/webtransport-go"
)

type Client interface {
	Start(wg *sync.WaitGroup)
	Stop()
}

type client struct {
	Dialer               webtransport.Dialer
	ServerStream         webtransport.Stream
	Config               *proto.Client
	BackupConfig         *proto.BackupConfig
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
	c.ServerStream = stream

	log.GetLogger().Success("Client is up and running! Ready to communicate with the server!")
	go c.handleStream()
}

func (c *client) Stop() {
	log.GetLogger().Normal("Stopping client...")
	if err := c.ServerStream.Close(); err != nil {
		log.GetLogger().Error(err.Error())
	}
	if err := c.Dialer.Close(); err != nil {
		log.GetLogger().Error(err.Error())
	}
	c.Backup.Stop()
}

func (c *client) handleStream() {
	log.GetLogger().Normal("Trying to request backup config from server...")
	if _, err := request.Write(c.ServerStream, &proto.Envelope{
		ClientId: c.Config.ClientId,
		Message: &proto.Envelope_BackupConfigRequest{
			BackupConfigRequest: &proto.BackupConfigRequest{}}}); err != nil {
		log.GetLogger().Fatal("Failed writing config getter request to stream: ", err.Error())
	}

	for {
		envelope := proto.Envelope{}
		_, readError := request.Read(c.ServerStream, &envelope)
		if readError != nil {
			log.GetLogger().Fatal("Read error occured during stream handling. Server error occured!", readError.Error())
		}

		switch message := envelope.Message.(type) {
		case *proto.Envelope_BackupConfigResponse:
			c.BackupConfig = message.BackupConfigResponse.BackupConfig

			log.GetLogger().Success("Got backup config from server!")
			c.startBackup()
		case *proto.Envelope_BackupChunkResponse:
			if message.BackupChunkResponse.IsOk {
				log.GetLogger().Success(fmt.Sprintf("Server processed %s chunk!", message.BackupChunkResponse.ChunkName))
			} else {
				log.GetLogger().Error(fmt.Sprintf("Server was not able to process %s chunk!", message.BackupChunkResponse.ChunkName))
			}

			if err := os.Remove(path.Join(CHUNK_TEMP_DIR, message.BackupChunkResponse.ChunkName)); err != nil {
				log.GetLogger().Error("Failed to remove chunk temp file: ", err.Error())
			}
		case *proto.Envelope_BackupEndResponse:
			if message.BackupEndResponse.IsOk {
				log.GetLogger().Success("Server processed the backup!")
			} else {
				log.GetLogger().Error("Server was not able to process the backup!")
			}
		}
	}
}

func (c *client) startBackup() {
	c.Backup = backup.NewBackup(c.BackupConfig)

	c.newBackupPathChannel = make(chan string)
	go c.Backup.Backup(c.newBackupPathChannel)
	go c.handNewBackup()
}
