package config

import (
	"github.com/Siposattila/go-backup/io"
	"github.com/Siposattila/go-backup/log"
)

type Server struct {
	Port                          string `json:"port"`
	Domain                        string `json:"domain"`
	Username                      string `json:"username"`
	Password                      string `json:"password"`
	BackupPath                    string `json:"backupPath"`
	StorageAlertTresholdInPercent int    `json:"storageAlertTresholdInPercent"`
	RegisterNodeIfNotKnown        bool   `json:"registerNodeIfNotKnown"`
	EmailAlert                    bool   `json:"emailAlert"`
	email
	DiscordAlert bool `json:"discordAlert"`
	discord
}

type email struct {
	EmailReceiver string `json:"emailReceiver"`
	EmailSender   string `json:"emailSender"`
	EmailUser     string `json:"emailUser"`
	EmailPassword string `json:"emailPassword"`
	EmailHost     string `json:"emailHost"`
	EmailPort     int    `json:"emailPort"`
}

type discord struct {
	DiscordWebHookId    string `json:"discordWebHookId"`
	DiscordWebHookToken string `json:"discordWebHookToken"`
}

func (s *Server) Get() *Server {
	var config Server
	rawConfig, readError := io.ReadFile(CONFIG_PATH, SERVER_CONFIG_FILENAME)
	if readError != nil {
		config = Server{
			Port:                          ":2000",
			Domain:                        "localhost",
			Username:                      "backup",
			Password:                      "123456",
			BackupPath:                    ".",
			StorageAlertTresholdInPercent: 95,
			EmailAlert:                    false,
			email: email{
				EmailReceiver: "example@example.com",
				EmailSender:   "example@example.com",
				EmailUser:     "",
				EmailPassword: "",
				EmailPort:     25,
				EmailHost:     "",
			},
			DiscordAlert: false,
			discord: discord{
				DiscordWebHookId:    "",
				DiscordWebHookToken: "",
			},
			RegisterNodeIfNotKnown: true,
		}

		generationError := generateConfig(&config, SERVER_CONFIG_FILENAME)
		if generationError != nil {
			log.GetLogger().Fatal(generationError.Error())
		}
	} else {
		loadedConfig, loadError := loadConfig[*Server](rawConfig)
		if loadError != nil {
			log.GetLogger().Fatal(loadError.Error())
		}

		config = **loadedConfig
	}

	return &config
}
