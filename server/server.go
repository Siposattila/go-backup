package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/Siposattila/gobkup/alert"
	"github.com/Siposattila/gobkup/client"
	"github.com/Siposattila/gobkup/config"
	"github.com/Siposattila/gobkup/disk"
	"github.com/Siposattila/gobkup/log"
	"github.com/Siposattila/gobkup/request"
	"github.com/Siposattila/gobkup/serializer"
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
	Config    *config.Server
	Discord   alert.AlertInterface
	Email     alert.AlertInterface
}

func NewServer() Server {
	var s server
	s.Config = s.Config.Get()
	s.Transport = webtransport.Server{
		H3: http3.Server{Addr: s.Config.Port, QUICConfig: &quic.Config{
			EnableDatagrams: true,
			Allow0RTT:       true,
			MaxIdleTimeout:  time.Duration(30 * time.Second),
			// TODO: need more research on quic flow control and BDP
			InitialStreamReceiveWindow:     20 << 20,  // 20 MB
			InitialConnectionReceiveWindow: 20 << 20,  // 20 MB
			MaxStreamReceiveWindow:         60 << 20,  // 60 MB
			MaxConnectionReceiveWindow:     150 << 20, // 150 MB
		}},
	}
	s.getTlsConfig()

	if s.Config.DiscordAlert {
		s.Discord = alert.NewDiscord(s.Config.DiscordWebHookId, s.Config.DiscordWebHookToken)
		s.Discord.Start()
	}

	if s.Config.EmailAlert {
		s.Email = alert.NewEmail(
			s.Config.EmailReceiver,
			s.Config.EmailSender,
			s.Config.EmailUser,
			s.Config.EmailPassword,
			s.Config.EmailHost,
			s.Config.EmailPort,
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
	s.Transport.Close()

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

func (s *server) handleStream(stream webtransport.Stream) {
	defer stream.Close()
	var clientId string

	for {
		r := request.Request{}
		_, readError := request.Read(stream, &r)
		if readError != nil {
			source := "unknown"
			if clientId != "" {
				source = clientId
			}

			log.GetLogger().Error("Read error occured during stream handling.", fmt.Sprintf("Client: %s", source), readError.Error())
			s.alertSystem(fmt.Sprintf("Error connection suddenly closed for client!\nClient: %s\n%s", source, readError.Error()))
			break
		}

		if clientId == "" {
			clientId = r.ClientId
		}

		switch r.Id {
		case request.ID_CONFIG:
			log.GetLogger().Normal(fmt.Sprintf("%s sent a request for it's backup config.", clientId))

			var backupConfig config.Backup
			backupConfig = *backupConfig.Get(clientId)
			request.Write(stream, request.NewResponse(request.ID_CONFIG, backupConfig))

			log.GetLogger().Success("Backup config sent to " + clientId)
		case request.ID_BACKUP_START:
			log.GetLogger().Normal(fmt.Sprintf("%s started sending backup...", clientId))

			info := client.Info{}
			serializerError := serializer.Json.Serialize([]byte(r.Data), &info)
			if serializerError != nil {
				log.GetLogger().Error("Error occured during getting file info.", serializerError.Error())
			} else {
				diskUsage := disk.NewDiskUsage("/")
				usageAfterBackupTransfer := int(diskUsage.Used()+uint64(info.Size)) * 100 / int(diskUsage.Size())
				if diskUsage.Usage() >= s.Config.StorageAlertTresholdInPercent || usageAfterBackupTransfer >= s.Config.StorageAlertTresholdInPercent {
					s.alertSystem(fmt.Sprintf("Warning the storage alert treshold was met! The current usage is: %d%s", usageAfterBackupTransfer, "%"))
				}
			}
		case request.ID_BACKUP_CHUNK:
			chunk := client.Chunk{}
			serializerError := serializer.Json.Serialize([]byte(r.Data), &chunk)
			if serializerError != nil {
				log.GetLogger().Error("Error occured during getting chunk.", serializerError.Error())
			} else {
				s.writeChunk(&chunk)
				chunk.Data = nil // do not need to send back the chunk data

				request.Write(stream, request.NewResponse(request.ID_BACKUP_CHUNK_PROCESSED, chunk))
			}
		case request.ID_BACKUP_END:
			log.GetLogger().Success(fmt.Sprintf("Received backup from %s...", clientId))

			info := client.Info{}
			serializerError := serializer.Json.Serialize([]byte(r.Data), &info)
			if serializerError != nil {
				log.GetLogger().Error("Error occured during getting file info.", serializerError.Error())
			} else {
				if err := os.Rename(path.Join(s.Config.BackupPath, fmt.Sprintf(TEMP_FILE, info.Name)), path.Join(s.Config.BackupPath, info.Name)); err != nil {
					log.GetLogger().Error(err.Error())
				}
			}
		}
	}
}
