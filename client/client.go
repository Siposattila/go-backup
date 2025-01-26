package client

import "github.com/quic-go/webtransport-go"

type Client interface {
	Start()
	Stop()
}

type client struct {
	Dialer webtransport.Dialer
	Stream webtransport.Stream
}

func NewClient() Client {
	return &client{}
}

func (c *client) Start() {}

func (c *client) Stop() {}
