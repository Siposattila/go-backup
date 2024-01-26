package config

type MasterConfig struct {
	Nodes                         map[string]NodeConfig
	BackupPath                    string
	WhenToBackup                  string
	ExcludeExtensions             []string
	ExcludeFiles                  []string
	StorageAlertTresholdInPercent int
	EmailAlert                    bool
	NotificationEmailReceiver     string
	NotificationEmailSender       string
	NotificationEmailAuth         NotificationEmailSenderAuth
	DiscordAlert                  bool
	NotificationDiscordAuth       NotificationDiscordSenderAuth
}

type NodeConfig struct {
	NodeId string
}

type NotificationEmailSenderAuth struct {
	User     string
	Password string
	Host     string
	Port     int
}

type NotificationDiscordSenderAuth struct {
	WebHookId    string
	WebHookToken string
}
