package config

import (
	"github.com/Siposattila/go-backup/io"
	"github.com/Siposattila/go-backup/log"
	"github.com/Siposattila/go-backup/proto"
)

func GetServerConfig() *proto.Server {
	var config *proto.Server
	rawConfig, readError := io.ReadFile(CONFIG_PATH, SERVER_CONFIG_FILENAME)
	if readError != nil {
		config = &proto.Server{
			Port:                          ":2000",
			Domain:                        "localhost",
			Username:                      "backup",
			Password:                      "123456",
			BackupPath:                    ".",
			StorageAlertTresholdInPercent: 95,
			EmailAlert:                    false,
			Email: &proto.Email{
				EmailReceiver: "example@example.com",
				EmailSender:   "example@example.com",
				EmailUser:     "",
				EmailPassword: "",
				EmailPort:     25,
				EmailHost:     "",
			},
			DiscordAlert: false,
			Discord: &proto.Discord{
				DiscordWebHookId:    "",
				DiscordWebHookToken: "",
			},
			RegisterNodeIfNotKnown: true,
		}

		generationError := generateConfig(config, SERVER_CONFIG_FILENAME)
		if generationError != nil {
			log.GetLogger().Fatal("Failed to generate server config: ", generationError.Error())
		}
	} else {
		loadedConfig, loadError := loadConfig[*proto.Server](rawConfig)
		if loadError != nil {
			log.GetLogger().Fatal("Failed to load server config: ", loadError.Error())
		}

		config = *loadedConfig
	}

	return config
}
