package config

type MasterConfig struct {
	Nodes                         []string `json:"nodes"`
	BackupPath                    string   `json:"backupPath"`
	WhenToBackup                  string   `json:"whenToBackup"`
	ExcludeExtensions             []string `json:"excludeExtensions"`
	ExcludeFiles                  []string `json:"excludeFiles"`
	StorageAlertTresholdInPercent int      `json:"storageAlertTresholdInPercent"`
	EmailAlert                    bool     `json:"emailAlert"`
	DiscordAlert                  bool     `json:"discordAlert"`
	Debug                         bool     `json:"-"`
}

type NodeConfig struct {
	NodeId            string   `json:"nodeId"`
	WhenToBackup      string   `json:"whenToBackup"`
	WhatToBackup      []string `json:"whatToBackup"`
	ExcludeExtensions []string `json:"excludeExtensions"`
	ExcludeFiles      []string `json:"excludeFiles"`
	Debug             bool     `json:"-"`
}

type NotificationEmailSenderAuth struct {
	EmailReceiver string `json:"emailReceiver"`
	EmailSender   string `json:"emailSender"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
}

type NotificationDiscordSenderAuth struct {
	WebHookId    string `json:"webHookId"`
	WebHookToken string `json:"webHookToken"`
}
