package config

type MasterConfig struct {
	Nodes                         map[string]string `json:"nodes"`
	Port                          string            `json:"port"`
	Domain                        string            `json:"domain"`
	BackupPath                    string            `json:"backupPath"`
	WhenToBackup                  string            `json:"whenToBackup"`
	ExcludeExtensions             []string          `json:"excludeExtensions"`
	ExcludeFiles                  []string          `json:"excludeFiles"`
	StorageAlertTresholdInPercent int               `json:"storageAlertTresholdInPercent"`
	EmailAlert                    bool              `json:"emailAlert"`
	DiscordAlert                  bool              `json:"discordAlert"`
	RegisterNodeIfKnown           bool              `json:"registerNodeIfNotKnown"`
	Debug                         bool              `json:"-"`
}

type NodeConfig struct {
	NodeId            string   `json:"nodeId"`
	WhenToBackup      string   `json:"whenToBackup"`
	WhatToBackup      []string `json:"whatToBackup"`
	ExcludeExtensions []string `json:"excludeExtensions"`
	ExcludeFiles      []string `json:"excludeFiles"`
	Token             string   `json:"-"`
	Debug             bool     `json:"-"`
}

type EmailConfig struct {
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type DiscordConfig struct {
	WebHookId    string `json:"webHookId"`
	WebHookToken string `json:"webHookToken"`
}
