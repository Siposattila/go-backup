package server

import (
	"context"
	"net/http"

	"github.com/Siposattila/gobkup/config"
	"github.com/Siposattila/gobkup/log"
	"github.com/Siposattila/gobkup/request"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

type Server interface {
	Start()
	Stop()
}

type server struct {
	Transport webtransport.Server
	Config    *config.Server
}

func NewServer() Server {
	var server server
	server.Config = server.Config.Get()
	server.Transport = webtransport.Server{
		H3: http3.Server{Addr: server.Config.Port},
	}
	server.getTlsConfig()

	if server.Config.DiscordAlert {
		// TODO: discord alert
	}

	if server.Config.EmailAlert {
		// TODO: email alert
	}

	return &server
}

func (s *server) Start() {
	log.GetLogger().Normal("Starting up server...")

	s.setupEndpoint()
	serverError := s.Transport.ListenAndServe()
	if serverError != nil {
		log.GetLogger().Fatal("Error during server listen and serve: ", serverError)
	}

	log.GetLogger().Success("Server is up and running! Ready to handle connections on port :" + s.Config.Port)
}

func (s *server) Stop() {
	log.GetLogger().Normal("Stopping server...")
	s.Transport.Close()
}

func (s *server) setupEndpoint() {
	http.HandleFunc("/server", func(w http.ResponseWriter, r *http.Request) {
		connection, error := s.Transport.Upgrade(w, r)
		if error != nil {
			w.WriteHeader(500)
			log.GetLogger().Fatal("Upgrading failed: ", error)
		}

		stream, streamError := connection.AcceptStream(context.Background())
		if streamError != nil {
			log.GetLogger().Fatal("There was an error during accepting the stream: ", streamError)
		}

		go s.handleStream(stream)
	})
}

func (s *server) handleStream(stream webtransport.Stream) {
	defer stream.Close()

	for {
		var r request.Request
		n, readError := request.Read(stream, r)
		if readError != nil {
			log.GetLogger().Error("Read error occured during stream handling. Stream ID: ", stream.StreamID())
			break
		}

		log.GetLogger().Debug("Read length: ", n)

		//authError := s.authClient(clientStream, incomingRequest.NodeId, incomingRequest.Token)
		//if authError != nil {
		//	log.GetLogger().Error("Authentication error occured during stream handling. Stream ID: ", stream.StreamID())
		//	break
		//}

		switch r.Id {
		case request.REQUEST_ID_CONFIG:
			log.GetLogger().Normal(r.DealerId + " sent a request for the config.")
			// nodeConfig := config.LoadNodeConfig(incomingRequest.NodeId)
			// s.writeToStream(stream, s.makeResponse(request.REQUEST_ID_CONFIG, string(serializer.Json.Deserialize(nodeConfig))))
			log.GetLogger().Success("Config sent to " + r.DealerId)
			break
		}
	}
}

//func (s *server) authClient(stream webtransport.Stream, nodeId string, requestToken string) error {
//	token, ok := s.Config.Nodes[nodeId]
//	if s.Config.RegisterNodeIfKnown && !ok {
//		s.AddNode(nodeId)
//		s.writeToStream(stream, s.makeResponse(request.REQUEST_ID_NODE_REGISTERED, "TOKEN"))
//
//		return errors.New("NEW")
//	}
//
//	if ok && requestToken != token {
//		s.writeToStream(stream, s.makeResponse(request.REQUEST_ID_AUTH_ERROR, "TOKEN"))
//
//		return errors.New("AUTH")
//	}
//
//	return nil
//}
