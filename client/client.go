package client

import (
	"github.com/Siposattila/gobkup/log"
	"github.com/quic-go/webtransport-go"
)

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

func (c *client) Start() {
	log.GetLogger().Normal("Starting up client...")

	// TODO: connect and handle
}

func (c *client) Stop() {
	// TODO: Stop
}

func (c *client) handleStream() {
	// TODO: handle
	//console.Normal("Trying to request for config from master...")
	//node.writeToStream(node.makeRequest(request.REQUEST_ID_CONFIG, "PLEASE"))
	//for {
	//	var response request.MasterResponse
	//	serializer.Json.Serialize(node.readFromStream(), &response)

	//	switch response.Id {
	//	case request.REQUEST_ID_CONFIG:
	//		if node.Backup != nil {
	//			node.Backup.Stop()
	//		}
	//		config := config.NodeConfig{}
	//		serializer.Json.Serialize([]byte(response.Data), &config)
	//		config.Debug = node.Config.Debug
	//		config.Token = node.Config.Token
	//		node.Config = config
	//		console.Success("Got config from master!")
	//		node.Backup = backup.NewBackup(node.Config.WhenToBackup, node.Config.WhatToBackup, node.Config.ExcludeExtensions, node.Config.ExcludeFiles)
	//		node.Backup.BackupProcess(node.Config.NodeId + "_backup_" + time.Now().String() + ".zip")
	//		break
	//	case request.REQUEST_ID_NODE_REGISTERED:
	//		console.Warning("This node is now registered at the master! The token was generated for this node at the master.")
	//		return
	//	case request.REQUEST_ID_AUTH_ERROR:
	//		console.Fatal("Wrong credentials!")
	//		break
	//	}
	//}
}
