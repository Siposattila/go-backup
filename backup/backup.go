package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/Siposattila/gobkup/log"
	"github.com/robfig/cron/v3"
)

type BackupInterface interface {
	Backup(newBackupPath chan<- string)
	Stop()
}

type backup struct {
	CronExpression string
	WhatToBackup   *[]string
	Exclude        *[]string
	Cron           cron.Schedule
	stopChannel    chan bool
}

func NewBackup(cronExpression string, whatToBackup *[]string, exclude *[]string) BackupInterface {
	var newBackup = backup{
		CronExpression: cronExpression,
		WhatToBackup:   whatToBackup,
		Exclude:        exclude,
	}

	var schedule, parseError = cron.
		NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).
		Parse(newBackup.CronExpression)
	if parseError != nil {
		log.GetLogger().Fatal("Your cron expression is invalid or an error occured!", parseError.Error())
	}
	newBackup.Cron = schedule

	return &newBackup
}

func (b *backup) Backup(newBackupPathChannel chan<- string) {
	var timeSignal = time.After(time.Until(b.Cron.Next(time.Now())))
	log.GetLogger().Success("Next backup will be at: " + b.Cron.Next(time.Now()).Format("15:04:05"))
	select {
	case <-timeSignal:
		log.GetLogger().Normal("Backing up...")

		c := compression{
			BackupPath: os.TempDir(),
			Paths:      b.WhatToBackup,
			Exclude:    b.Exclude,
			Store:      getStore(),
		}
		zipPath := c.zipCompress(fmt.Sprintf("%s_backup.zip", time.Now().Format("20060102150405")))

		log.GetLogger().Success("Backup finished successfully!")
		if zipPath != "" {
			newBackupPathChannel <- zipPath
		}

		b.Backup(newBackupPathChannel)
	case <-b.stopChannel:
		log.GetLogger().Normal("Stopping backup process...")
	}
}

func (b *backup) Stop() {
	b.stopChannel <- true
}
