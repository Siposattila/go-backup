package node

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
)

func AddNode(nodeId string) {
	config.Master.Nodes[nodeId] = getMD5Hash(nodeId + time.Now().String())
	config.UpdateConfig("master")
	config.GenerateNodeConfig(nodeId)
	console.Success("Node " + nodeId + " has been successfully added!")

	return
}

func getMD5Hash(plain string) string {
	var hash = md5.Sum([]byte(plain))

	return hex.EncodeToString(hash[:])
}
