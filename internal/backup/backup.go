package backup

import (
	"time"

	"github.com/Siposattila/gobkup/internal/compression"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/robfig/cron/v3"
)

type BackupInterface interface {
	BackupProcess(zipName string)
	Stop()
}

type backup struct {
	CronExpression    string
	WhatToBackup      []string
	ExcludeExtensions []string
	ExcludeFiles      []string
	Cron              cron.Schedule
	stopChannel       chan bool
}

func NewBackup(cronExpression string, whatToBackup []string, excludeExtensions []string, excludeFiles []string) BackupInterface {
	newBackup := backup{
		CronExpression:    cronExpression,
		WhatToBackup:      whatToBackup,
		ExcludeExtensions: excludeExtensions,
		ExcludeFiles:      excludeFiles,
	}
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, parseError := cronParser.Parse(newBackup.CronExpression)
	if parseError != nil {
		console.Fatal("Your cron expression is invalid or an error occured: " + parseError.Error())
	}
	newBackup.Cron = schedule

	return newBackup
}

func (b backup) BackupProcess(zipName string) {
	timeSignal := time.After(b.Cron.Next(time.Now()).Sub(time.Now()))
	console.Success("Next backup will be at: " + b.Cron.Next(time.Now()).String())
	select {
	case <-timeSignal:
		console.Normal("Backing the stuff up. This may take a long time!!!")
		for _, path := range b.WhatToBackup {
			compress := compression.Compression{Path: path}
			compress.ZipCompress("name.zip")
		}
		console.Success("Backup finished successfully!")
		b.BackupProcess(zipName)
	case <-b.stopChannel:
		console.Normal("Stopping backup process...")
	}
}

func (b backup) Stop() {
	b.stopChannel <- true
}
