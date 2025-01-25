package node

import (
	"time"

	"github.com/Siposattila/gobkup/internal/backup"
	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/Siposattila/gobkup/internal/request"
	"github.com/Siposattila/gobkup/internal/serializer"
)

func (node *Node) handleStream() {
	console.Normal("Trying to request for config from master...")
	node.writeToStream(node.makeRequest(request.REQUEST_ID_CONFIG, "PLEASE"))
	for {
		var response request.MasterResponse
		serializer.Json.Serialize(node.readFromStream(), &response)

		switch response.Id {
		case request.REQUEST_ID_CONFIG:
			if node.Backup != nil {
				node.Backup.Stop()
			}
			config := config.NodeConfig{}
			serializer.Json.Serialize([]byte(response.Data), &config)
			config.Debug = node.Config.Debug
			config.Token = node.Config.Token
			node.Config = config
			console.Success("Got config from master!")
			node.Backup = backup.NewBackup(node.Config.WhenToBackup, node.Config.WhatToBackup, node.Config.ExcludeExtensions, node.Config.ExcludeFiles)
			node.Backup.BackupProcess(node.Config.NodeId + "_backup_" + time.Now().String() + ".zip")
			break
		case request.REQUEST_ID_NODE_REGISTERED:
			console.Warning("This node is now registered at the master! The token was generated for this node at the master.")
			return
		case request.REQUEST_ID_AUTH_ERROR:
			console.Fatal("Wrong credentials!")
			break
		}
	}
}

func (node *Node) makeRequest(id int, data string) request.NodeRequest {
	return request.NodeRequest{
		Id:     id,
		Data:   data,
		NodeId: node.getClientName(),
		Token:  node.Config.Token,
	}
}

func (node *Node) writeToStream(data any) {
	_, writeError := node.ServerStream.Write(serializer.Json.Deserialize(data))
	if writeError != nil {
		console.Error("Error during write to master: " + writeError.Error())
	}
}

func (node *Node) readFromStream() []byte {
	buffer := make([]byte, 1024)
	n, readError := node.ServerStream.Read(buffer)
	if readError != nil {
		console.Error("Error during reading from master: " + readError.Error())
	}

	return buffer[:n]
}
