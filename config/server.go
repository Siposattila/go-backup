package config

import "github.com/Siposattila/gobkup/io"

type Server struct {
	Port                          string `json:"port"`
	Domain                        string `json:"domain"`
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
		generateConfig(&config, SERVER_CONFIG_FILENAME)
	} else {

	}

	return &config
}
