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

func generateConfig[T *Server | *Client](config T, configName string) error {
	io.CreateDir(CONFIG_PATH)
	buffer, serializerError := serializer.Json.Deserialize(config)
	if serializerError != nil {
		return serializerError
	}

	io.WriteFile(CONFIG_PATH, configName, buffer)

	return nil
}

func loadConfig[T *Server | *Client](rawConfig []byte) (*T, error) {
	config := new(T)
	serializerError := serializer.Json.Serialize(rawConfig, config)
	if serializerError != nil {
		return nil, serializerError
	}

	return config, nil
}
