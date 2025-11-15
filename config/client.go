package config

import (
	"os"
	"strings"

	"github.com/Siposattila/go-backup/generatedproto"
	"github.com/Siposattila/go-backup/io"
	"github.com/Siposattila/go-backup/log"
)

func getBackupConfigName(clientId string) string {
	return strings.Replace(BACKUP_CONFIG_FILENAME, "?", clientId, 1)
}

func getClientName() string {
	name, nameError := os.Hostname()
	if nameError != nil {
		log.GetLogger().Fatal("Failed to get hostname: ", nameError.Error())
	}

	return name
}

func GetClientConfig() *generatedproto.Client {
	var config *generatedproto.Client
	rawConfig, readError := io.ReadFile(CONFIG_PATH, CLIENT_CONFIG_FILENAME)
	if readError != nil {
		config = &generatedproto.Client{
			ClientId: getClientName(),
			Token:    "",
			Endpoint: "https://localhost:2000",
		}

		generationError := generateConfig(config, CLIENT_CONFIG_FILENAME)
		if generationError != nil {
			log.GetLogger().Fatal("Failed to generate client config: ", generationError.Error())
		}
	} else {
		loadedConfig, loadError := loadConfig[*generatedproto.Client](rawConfig)
		if loadError != nil {
			log.GetLogger().Fatal("Failed to load client config: ", loadError.Error())
		}

		config = *loadedConfig
	}

	return config
}

func GetBackupConfig(clientId string) *generatedproto.Backup {
	var config *generatedproto.Backup
	rawConfig, readError := io.ReadFile(CONFIG_PATH, getBackupConfigName(clientId))
	if readError != nil {
		config = &generatedproto.Backup{
			WhenToBackup: "0 0 * * *",
			WhatToBackup: []string{},
			Exclude:      []string{},
		}

		generationError := generateConfig(config, getBackupConfigName(clientId))
		if generationError != nil {
			log.GetLogger().Fatal("Failed to generate backup config: ", generationError.Error())
		}
	} else {
		loadedConfig, loadError := loadConfig[*generatedproto.Backup](rawConfig)
		if loadError != nil {
			log.GetLogger().Fatal("Failed to get backup config: ", loadError.Error())
		}

		config = *loadedConfig
	}

	return config
}
