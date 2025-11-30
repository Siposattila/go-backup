package config

import (
	"github.com/Siposattila/go-backup/io"
	"github.com/Siposattila/go-backup/proto"
	"github.com/Siposattila/go-backup/serializer"
)

const (
	CONFIG_PATH            = "./config.d"
	SERVER_CONFIG_FILENAME = "server.json"
	CLIENT_CONFIG_FILENAME = "client.json"
	BACKUP_CONFIG_FILENAME = "backup_?.json"
)

func generateConfig[T *proto.Server | *proto.Client | *proto.BackupConfig](config T, configName string) error {
	io.CreateDir(CONFIG_PATH)
	buffer, serializerError := serializer.Json.Deserialize(config)
	if serializerError != nil {
		return serializerError
	}

	io.WriteFile(CONFIG_PATH, configName, buffer)

	return nil
}

func loadConfig[T *proto.Server | *proto.Client | *proto.BackupConfig](rawConfig []byte) (*T, error) {
	config := new(T)
	serializerError := serializer.Json.Serialize(rawConfig, config)
	if serializerError != nil {
		return nil, serializerError
	}

	return config, nil
}
