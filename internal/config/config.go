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
	masterConfig := &MasterConfig{
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
}

func GenerateNodeConfig(nodeId string) {
	nodeConfig := &NodeConfig{
		NodeId:            nodeId,
		WhenToBackup:      "",
		WhatToBackup:      []string{},
		ExcludeExtensions: []string{},
		ExcludeFiles:      []string{},
	}

	generateConfig(nodeConfig, getNodeConfigName(nodeId))
	console.Success("Node config generated!")
}

func generateDiscordAlertConfig() {
	discordConfig := &DiscordConfig{
		WebHookId:    "",
		WebHookToken: "",
	}

	generateConfig(discordConfig, discordAlertConfigName)
	console.Success("Discord config generated!")
}

func generateEmailAlertConfig() {
	emailConfig := &EmailConfig{
		Receiver: "example@example.com",
		Sender:   "example@example.com",
		User:     "",
		Password: "",
		Port:     25,
		Host:     "",
	}

	generateConfig(emailConfig, emailAlertConfigName)
	console.Success("Email config generated!")
}

func getNodeConfigName(nodeId string) string {
	return strings.Replace(nodeConfigName, "?", nodeId, 1)
}

func generateConfig(config any, configName string) {
	io.CreateDir(configPath)
	buffer := serializer.Json.Deserialize(config)
	io.WriteFile(configPath, configName, buffer)
}

func UpdateMasterConfig(config MasterConfig) {
	io.WriteFile(configPath, masterConfigName, serializer.Json.Deserialize(config))
	console.Success("Master config updated!")
}

func LoadMasterConfig() MasterConfig {
	config := MasterConfig{}
	serializer.Json.Serialize(io.ReadFile(configPath, masterConfigName), &config)
	console.Success("Master config loaded!")

	return config
}

// FIXME: do not load the master config again its not needed but for now its okay
func LoadNodeConfig(nodeId string) NodeConfig {
	config := NodeConfig{}
	masterConfig := MasterConfig{}
	serializer.Json.Serialize(io.ReadFile(configPath, getNodeConfigName(nodeId)), &config)
	serializer.Json.Serialize(io.ReadFile(configPath, masterConfigName), &masterConfig)

	config.ExcludeFiles = append(config.ExcludeFiles, masterConfig.ExcludeFiles...)
	config.ExcludeExtensions = append(config.ExcludeExtensions, masterConfig.ExcludeExtensions...)
	if config.WhenToBackup == "" {
		config.WhenToBackup = masterConfig.WhenToBackup
	}
	console.Success("Node config loaded!")

	return config
}

func LoadDiscordAlertConfig() DiscordConfig {
	config := DiscordConfig{}
	serializer.Json.Serialize(io.ReadFile(configPath, discordAlertConfigName), &config)
	console.Success("Discord config loaded!")

	return config
}

func LoadEmailAlertConfig() EmailConfig {
	config := EmailConfig{}
	serializer.Json.Serialize(io.ReadFile(configPath, emailAlertConfigName), &config)
	console.Success("Email config loaded!")

	return config
}
