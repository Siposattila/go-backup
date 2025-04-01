package config

import (
	"os"
	"strings"

	"github.com/Siposattila/gobkup/io"
	"github.com/Siposattila/gobkup/log"
)

type Client struct {
	ClientId string `json:"clientId"`
	Token    string `json:"token"`
	Endpoint string `json:"endpoint"`
}

type Backup struct {
	WhenToBackup      string   `json:"whenToBackup"`
	WhatToBackup      []string `json:"whatToBackup"`
	ExcludeExtensions []string `json:"excludeExtensions"`
	ExcludeFiles      []string `json:"excludeFiles"`
}

func getBackupConfigName(clientId string) string {
	return strings.Replace(BACKUP_CONFIG_FILENAME, "?", clientId, 1)
}

func getClientName() string {
	name, nameError := os.Hostname()
	if nameError != nil {
		log.GetLogger().Fatal(nameError.Error())
	}

	return name
}

func (c *Client) Get() *Client {
	var config Client
	rawConfig, readError := io.ReadFile(CONFIG_PATH, CLIENT_CONFIG_FILENAME)
	if readError != nil {
		config = Client{
			ClientId: getClientName(),
			Token:    "",
			Endpoint: "https://localhost:2000",
		}

		generationError := generateConfig(&config, CLIENT_CONFIG_FILENAME)
		if generationError != nil {
			log.GetLogger().Fatal(generationError.Error())
		}
	} else {
		loadedConfig, loadError := loadConfig[*Client](rawConfig)
		if loadError != nil {
			log.GetLogger().Fatal(loadError.Error())
		}

		config = **loadedConfig
	}

	return &config
}

func (b *Backup) Get(clientId string) *Backup {
	var config Backup
	rawConfig, readError := io.ReadFile(CONFIG_PATH, getBackupConfigName(clientId))
	if readError != nil {
		config = Backup{
			WhenToBackup:      "0 0 * * *",
			WhatToBackup:      []string{},
			ExcludeExtensions: []string{},
			ExcludeFiles:      []string{},
		}

		generationError := generateConfig(&config, getBackupConfigName(clientId))
		if generationError != nil {
			log.GetLogger().Fatal(generationError.Error())
		}
	} else {
		loadedConfig, loadError := loadConfig[*Backup](rawConfig)
		if loadError != nil {
			log.GetLogger().Fatal(loadError.Error())
		}

		config = **loadedConfig
	}

	return &config
}
