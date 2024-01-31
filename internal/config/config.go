package config

import (
	"strings"

	"github.com/Siposattila/gobkup/internal/console"
	"github.com/Siposattila/gobkup/internal/io"
	"github.com/Siposattila/gobkup/internal/serializer"
)

const configPath = "./config"
const masterConfigName = "master.json"
const nodeConfigName = "node_?.json"
const discordAlertConfigName = "discord.json"
const emailAlertConfigName = "email.json"

func GenerateMasterConfig() {
	var masterConfig = &MasterConfig{
		Nodes:                         map[string]string{},
		Port:                          ":2000",
		Domain:                        "localhost",
		BackupPath:                    ".",
		WhenToBackup:                  "*/2 * * * *",
		ExcludeExtensions:             []string{},
		ExcludeFiles:                  []string{},
		StorageAlertTresholdInPercent: 95,
		EmailAlert:                    false,
		DiscordAlert:                  false,
		RegisterNodeIfKnown:           true,
	}

	generateConfig(masterConfig, masterConfigName)
	console.Success("Master config generated!")
	generateDiscordAlertConfig()
	generateEmailAlertConfig()

	return
}

func GenerateNodeConfig(nodeId string) {
	var nodeConfig = &NodeConfig{
		NodeId:            nodeId,
		WhenToBackup:      "",
		WhatToBackup:      []string{},
		ExcludeExtensions: []string{},
		ExcludeFiles:      []string{},
	}

	generateConfig(nodeConfig, getNodeConfigName(nodeId))
	console.Success("Node config generated!")

	return
}

func generateDiscordAlertConfig() {
	var discordConfig = &DiscordConfig{
		WebHookId:    "",
		WebHookToken: "",
	}

	generateConfig(discordConfig, discordAlertConfigName)
	console.Success("Discord config generated!")

	return
}

func generateEmailAlertConfig() {
	var emailConfig = &EmailConfig{
		Receiver: "example@example.com",
		Sender:   "example@example.com",
		User:     "",
		Password: "",
		Port:     25,
		Host:     "",
	}

	generateConfig(emailConfig, emailAlertConfigName)
	console.Success("Email config generated!")

	return
}

func getNodeConfigName(nodeId string) string {
	return strings.Replace(nodeConfigName, "?", nodeId, 1)
}

func generateConfig(config any, configName string) {
	io.CreateDir(configPath)
	var buffer = serializer.Deserialize(config)
	io.WriteFile(configPath, configName, buffer)

	return
}

func UpdateMasterConfig(config MasterConfig) {
	io.WriteFile(configPath, masterConfigName, serializer.Deserialize(config))
	console.Success("Master config updated!")

	return
}

func LoadMasterConfig() MasterConfig {
	var config = MasterConfig{}
	serializer.Serialize(io.ReadFile(configPath, masterConfigName), &config)
	console.Success("Master config loaded!")

	return config
}

// FIXME: do not load the master config again its not needed but for now its okay
func LoadNodeConfig(nodeId string) NodeConfig {
	var config = NodeConfig{}
	var masterConfig = MasterConfig{}
	serializer.Serialize(io.ReadFile(configPath, getNodeConfigName(nodeId)), &config)
	serializer.Serialize(io.ReadFile(configPath, masterConfigName), &masterConfig)

	config.ExcludeFiles = append(config.ExcludeFiles, masterConfig.ExcludeFiles...)
	config.ExcludeExtensions = append(config.ExcludeExtensions, masterConfig.ExcludeExtensions...)
	if config.WhenToBackup == "" {
		config.WhenToBackup = masterConfig.WhenToBackup
	}
	console.Success("Node config loaded!")

	return config
}

func LoadDiscordAlertConfig() DiscordConfig {
	var config = DiscordConfig{}
	serializer.Serialize(io.ReadFile(configPath, discordAlertConfigName), &config)
	console.Success("Discord config loaded!")

	return config
}

func LoadEmailAlertConfig() EmailConfig {
	var config = EmailConfig{}
	serializer.Serialize(io.ReadFile(configPath, emailAlertConfigName), &config)
	console.Success("Email config loaded!")

	return config
}
