package node

import (
	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
)

func AddNode(nodeId string) {
	config.Master.Nodes = append(config.Master.Nodes, nodeId)
	config.UpdateConfig("master")
	config.GenerateNodeConfig(nodeId)
	console.Success("Node " + nodeId + " has been successfully added!")

	return
}
