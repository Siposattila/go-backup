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
		n, readError := request.Read(stream, &r)
		if readError != nil {
			source := "unknown"
			if clientId != "" {
				source = clientId
			}

			log.GetLogger().Error("Read error occured during stream handling.", "Client: "+source, readError.Error())
			s.alertSystem("Error connection suddenly closed for client!\nClient: " + source + "\n" + readError.Error())
			break
		}

		log.GetLogger().Debug("Read length: ", n)
		if clientId == "" {
			clientId = r.ClientId
		}

		switch r.Id {
		case request.ID_CONFIG:
			log.GetLogger().Normal(clientId + " sent a request for it's backup config.")

			var backupConfig config.Backup
			backupConfig = *backupConfig.Get(clientId)
			request.Write(stream, request.NewResponse(request.ID_CONFIG, backupConfig))

			log.GetLogger().Success("Backup config sent to " + clientId)
		case request.ID_BACKUP_START:
			log.GetLogger().Debug("Backup receiving start!")

			info := client.Info{}
			serializerError := serializer.Json.Serialize([]byte(r.Data), &info)
			if serializerError != nil {
				log.GetLogger().Error("Error occured during getting file info.", serializerError.Error())
			} else {
				log.GetLogger().Debug(info.Name, info.Size)
				// TODO: check if enough space is available
				// if not or after this no more space is available then send alert
			}
		case request.ID_BACKUP_CHUNK:
			log.GetLogger().Debug("Received new chunk!")

			chunk := client.Chunk{}
			serializerError := serializer.Json.Serialize([]byte(r.Data), &chunk)
			if serializerError != nil {
				log.GetLogger().Error("Error occured during getting chunk.", serializerError.Error())
			} else {
				log.GetLogger().Debug(chunk.Name, chunk.ChunkName, chunk.Size)
				s.writeChunk(&chunk)
				request.Write(stream, request.NewResponse(request.ID_BACKUP_CHUNK_PROCESSED, chunk))
			}
		case request.ID_BACKUP_END:
			log.GetLogger().Debug("Backup receiving ended!")

			info := client.Info{}
			serializerError := serializer.Json.Serialize([]byte(r.Data), &info)
			if serializerError != nil {
				log.GetLogger().Error("Error occured during getting file info.", serializerError.Error())
			} else {
				log.GetLogger().Debug(info.Name, info.Size)
				if err := os.Rename(path.Join(s.Config.BackupPath, fmt.Sprintf(TEMP_FILE, info.Name)), path.Join(s.Config.BackupPath, info.Name)); err != nil {
					log.GetLogger().Error("Failed to rename temp to normal.")
				}
			}
		}
	}
}
