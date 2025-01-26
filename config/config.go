package config

import (
	"strings"

	"github.com/Siposattila/gobkup/io"
	"github.com/Siposattila/gobkup/serializer"
)

const (
	CONFIG_PATH            = "./config.d"
	SERVER_CONFIG_FILENAME = "server.json"
	CLIENT_CONFIG_FILENAME = "client_?.json"
)

type Config[T Server | Client] interface {
	Get() T
}

func getNodeConfigName(clientId string) string {
	return strings.Replace(CLIENT_CONFIG_FILENAME, "?", clientId, 1)
}

func generateConfig[T *Server | *Client](config T, configName string) {
	io.CreateDir(CONFIG_PATH)
	buffer := serializer.Json.Deserialize(config)
	io.WriteFile(CONFIG_PATH, configName, buffer)
}

//func loadConfig(rawConfig []byte, configName string) {
//	io.CreateDir(CONFIG_PATH)
//	buffer := serializer.Json.Deserialize(rawConfig)
//	io.WriteFile(CONFIG_PATH, configName, buffer)
//}
