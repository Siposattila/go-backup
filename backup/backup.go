package backup

import (
	"time"

	"github.com/Siposattila/gobkup/log"
	"github.com/robfig/cron/v3"
)

type BackupInterface interface {
	Backup()
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
	var newBackup = backup{
		CronExpression:    cronExpression,
		WhatToBackup:      whatToBackup,
		ExcludeExtensions: excludeExtensions,
		ExcludeFiles:      excludeFiles,
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

func (b *backup) Backup() {
	var timeSignal = time.After(time.Until(b.Cron.Next(time.Now())))
	log.GetLogger().Success("Next backup will be at: " + b.Cron.Next(time.Now()).String())
	select {
	case <-timeSignal:
		log.GetLogger().Normal("Backing up...")
		// fmt.Sprintf("%s_backup.zip", time.Now().String()
		//for _, path := range b.WhatToBackup {
		// TODO: implement backup process
		//}
		log.GetLogger().Success("Backup finished successfully!")
		b.Backup()
	case <-b.stopChannel:
		log.GetLogger().Normal("Stopping backup process...")
	}
}

func (b *backup) Stop() {
	b.stopChannel <- true
}
