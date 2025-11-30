package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/Siposattila/go-backup/alert"
	"github.com/Siposattila/go-backup/config"
	"github.com/Siposattila/go-backup/disk"
	"github.com/Siposattila/go-backup/log"
	"github.com/Siposattila/go-backup/proto"
	"github.com/Siposattila/go-backup/request"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

type Server interface {
	Start(wg *sync.WaitGroup)
	Stop()
}

type server struct {
	Transport webtransport.Server
	Config    *proto.Server
	Discord   alert.AlertInterface
	Email     alert.AlertInterface
}

func NewServer() Server {
	var s server
	s.Config = config.GetServerConfig()
	s.Transport = webtransport.Server{
		H3: http3.Server{Addr: s.Config.Port, QUICConfig: &quic.Config{
			EnableDatagrams:                true,
			Allow0RTT:                      true,
			MaxIdleTimeout:                 time.Duration(30 * time.Second),
			InitialStreamReceiveWindow:     20 << 20,  // 20 megabytes
			InitialConnectionReceiveWindow: 20 << 20,  // 20 megabytes
			MaxStreamReceiveWindow:         60 << 20,  // 60 megabytes
			MaxConnectionReceiveWindow:     150 << 20, // 150 megabytes
		}},
	}
	s.getTlsConfig()

	if s.Config.DiscordAlert {
		s.Discord = alert.NewDiscord(s.Config.Discord.DiscordWebHookId, s.Config.Discord.DiscordWebHookToken)
		s.Discord.Start()
	}

	if s.Config.EmailAlert {
		s.Email = alert.NewEmail(
			s.Config.Email.EmailReceiver,
			s.Config.Email.EmailSender,
			s.Config.Email.EmailUser,
			s.Config.Email.EmailPassword,
			s.Config.Email.EmailHost,
			int(s.Config.Email.EmailPort),
		)
		s.Email.Start()
	}

	return &s
}

func (s *server) Start(serverWg *sync.WaitGroup) {
	log.GetLogger().Normal("Starting up server...")
	serverWg.Add(1)

	s.setupEndpoint()
	go func() {
		defer serverWg.Done()
		serverError := s.Transport.ListenAndServe()
		if serverError != nil {
			log.GetLogger().Fatal("Error during server listen and serve.", serverError.Error())
		}
	}()

	log.GetLogger().Success("Server is up and running! Ready to handle connections on port " + s.Config.Port)
}

func (s *server) Stop() {
	log.GetLogger().Normal("Stopping server...")
	if err := s.Transport.Close(); err != nil {
		log.GetLogger().Error(err.Error())
	}

	if s.Discord != nil {
		s.Discord.Stop()
	}

	if s.Email != nil {
		s.Email.Stop()
	}
}

func (s *server) setupEndpoint() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			log.GetLogger().Error("Authentication error occured.", r.RemoteAddr)
			w.WriteHeader(401)

			return
		}

		if username != s.Config.Username || password != s.Config.Password {
			log.GetLogger().Error("Authentication error occured (Wrong creds).", r.RemoteAddr)
			w.WriteHeader(401)

			return
		}

		connection, error := s.Transport.Upgrade(w, r)
		if error != nil {
			log.GetLogger().Error("Upgrading failed.", error.Error(), r.RemoteAddr)
			w.WriteHeader(500)

			return
		}

		stream, streamError := connection.AcceptStream(context.Background())
		if streamError != nil {
			log.GetLogger().Error("There was an error during accepting the stream.", streamError.Error(), r.RemoteAddr)
			w.WriteHeader(500)

			return
		}

		go s.handleStream(stream)
	})
}

type clientStream struct {
	ClientId string
}

func (s *server) handleStream(stream webtransport.Stream) {
	client := clientStream{
		ClientId: "unknown",
	}

	for {
		envelope := proto.Envelope{}
		_, readError := request.Read(stream, &envelope)
		if readError != nil {
			log.GetLogger().Error("Read error occured during stream handling.", fmt.Sprintf("Client: %s", client.ClientId), readError.Error())
			s.alertSystem(fmt.Sprintf("Error connection suddenly closed for client!\nClient: %s\n%s", client.ClientId, readError.Error()))
			break
		}

		if envelope.ClientId == "" {
			log.GetLogger().Error("Should provide a valid clientId. (Can't use empty clientId)")
			break
		}

		if client.ClientId == "unknown" {
			client.ClientId = envelope.ClientId
		}

		switch message := envelope.Message.(type) {
		case *proto.Envelope_BackupConfigRequest:
			log.GetLogger().Normal(fmt.Sprintf("%s sent a request for it's backup config.", client.ClientId))

			response := &proto.Envelope{
				Message: &proto.Envelope_BackupConfigResponse{
					BackupConfigResponse: &proto.BackupConfigResponse{
						BackupConfig: config.GetBackupConfig(client.ClientId),
					},
				},
			}
			if _, err := request.Write(stream, response); err != nil {
				log.GetLogger().Error(fmt.Sprintf("An error occured during sending the backup config back to %s", client.ClientId), err.Error())
			}

			log.GetLogger().Success(fmt.Sprintf("Backup config sent to %s", client.ClientId))
		case *proto.Envelope_BackupStartRequest:
			log.GetLogger().Normal(fmt.Sprintf("%s started sending backup...", client.ClientId))

			diskUsage := disk.NewDiskUsage("/")
			usageAfterBackupTransfer := int32(diskUsage.Used()+uint64(message.BackupStartRequest.Size)) * 100 / int32(diskUsage.Size())
			if diskUsage.Usage() >= s.Config.StorageAlertTresholdInPercent || usageAfterBackupTransfer >= s.Config.StorageAlertTresholdInPercent {
				// TODO: if the threshold was hit then should do something about it
				log.GetLogger().Warning(fmt.Sprintf("A backup from this client %s will put the storage above the set threshold.", client.ClientId))
				s.alertSystem(fmt.Sprintf("Warning the storage alert threshold was met! The current usage is: %d%s", usageAfterBackupTransfer, "%"))
			}
		case *proto.Envelope_BackupChunkRequest:
			writeChunkError := s.writeChunk(message.BackupChunkRequest.Chunk)
			if writeChunkError != nil {
				log.GetLogger().Error(
					fmt.Sprintf("Failed to write chunk %s from %s", message.BackupChunkRequest.Chunk.ChunkName, client.ClientId),
					writeChunkError.Error(),
				)
			} else {
				log.GetLogger().Success(fmt.Sprintf("Processed chunk %s from %s", message.BackupChunkRequest.Chunk.ChunkName, client.ClientId))
			}

			response := &proto.Envelope{
				Message: &proto.Envelope_BackupChunkResponse{
					BackupChunkResponse: &proto.BackupChunkResponse{
						ChunkName: message.BackupChunkRequest.Chunk.ChunkName,
						IsOk:      writeChunkError == nil,
					},
				},
			}
			if _, err := request.Write(stream, response); err != nil {
				log.GetLogger().Error(fmt.Sprintf("Failed to write backup processed response to %s", client.ClientId), err.Error())
			}
		case *proto.Envelope_BackupEndRequest:
			log.GetLogger().Success(fmt.Sprintf("Received backup from %s...", client.ClientId))

			renameError := os.Rename(
				path.Join(
					s.Config.BackupPath,
					fmt.Sprintf(TEMP_FILE, message.BackupEndRequest.Name),
				),
				path.Join(s.Config.BackupPath, message.BackupEndRequest.Name),
			)
			if renameError != nil {
				log.GetLogger().Error(fmt.Sprintf("Failed to finish saving the backup that came from client: %s", client.ClientId), renameError.Error())
			}

			response := &proto.Envelope{
				Message: &proto.Envelope_BackupEndResponse{
					BackupEndResponse: &proto.BackupEndResponse{
						IsOk: renameError == nil,
					},
				},
			}
			if _, err := request.Write(stream, response); err != nil {
				log.GetLogger().Error(fmt.Sprintf("Failed to write backup end response to %s", client.ClientId), err.Error())
			}
		}
	}

	if err := stream.Close(); err != nil {
		log.GetLogger().Fatal(fmt.Sprintf("Failed to close stream for client (%s): ", client.ClientId), err.Error())
	}
}
