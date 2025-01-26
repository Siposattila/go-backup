package config

type Client struct {
	ClientId          string   `json:"clientId"`
	WhenToBackup      string   `json:"whenToBackup"`
	WhatToBackup      []string `json:"whatToBackup"`
	ExcludeExtensions []string `json:"excludeExtensions"`
	ExcludeFiles      []string `json:"excludeFiles"`
}
