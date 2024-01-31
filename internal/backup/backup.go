package backup

import (
	"time"

	"github.com/Siposattila/gobkup/internal/console"
	"github.com/robfig/cron/v3"
)

type BackupInterface interface {
	BackupProcess()
}

type backup struct {
	CronExpression    string
	WhatToBackup      []string
	ExcludeExtensions []string
	ExcludeFiles      []string
	Cron              cron.Schedule
}

func NewBackup(cronExpression string, whatToBackup []string, excludeExtensions []string, excludeFiles []string) BackupInterface {
	return backup{
		CronExpression:    cronExpression,
		WhatToBackup:      whatToBackup,
		ExcludeExtensions: excludeExtensions,
		ExcludeFiles:      excludeFiles,
	}
}

func (b backup) BackupProcess() {
	var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	var schedule, parseError = cronParser.Parse(b.CronExpression)
    if parseError != nil {
        console.Fatal("Your cron expression is invalid or an error occured: " + parseError.Error())
    }
    b.Cron = schedule
	console.Debug(b.Cron.Next(time.Now()).String())

	var timeSignal = time.After(b.Cron.Next(time.Now()).Sub(time.Now()))
	select {
	case <-timeSignal:
		// TODO: implement backup
		console.Debug("TODO backup")
		//BackupProcess()
	}

	return
}
