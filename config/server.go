package config

import (
	"github.com/Siposattila/go-backup/generatedproto"
	"github.com/Siposattila/go-backup/io"
	"github.com/Siposattila/go-backup/log"
)

func GetServerConfig() *generatedproto.Server {
	var config *generatedproto.Server
	rawConfig, readError := io.ReadFile(CONFIG_PATH, SERVER_CONFIG_FILENAME)
	if readError != nil {
		config = &generatedproto.Server{
			Port:                          ":2000",
			Domain:                        "localhost",
			Username:                      "backup",
			Password:                      "123456",
			BackupPath:                    ".",
			StorageAlertTresholdInPercent: 95,
			EmailAlert:                    false,
			Email: &generatedproto.Email{
				EmailReceiver: "example@example.com",
				EmailSender:   "example@example.com",
				EmailUser:     "",
				EmailPassword: "",
				EmailPort:     25,
				EmailHost:     "",
			},
			DiscordAlert: false,
			Discord: &generatedproto.Discord{
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
		loadedConfig, loadError := loadConfig[*generatedproto.Server](rawConfig)
		if loadError != nil {
			log.GetLogger().Fatal("Failed to load server config: ", loadError.Error())
		}

		config = *loadedConfig
	}

	return config
}
