package server

import "github.com/quic-go/webtransport-go"

type Server interface {
	Start()
	Stop()
}

type server struct {
	Transport webtransport.Server
}

func NewServer() Server {
	return &server{}
}

func (s *server) Start() {}

func (s *server) Stop() {}
