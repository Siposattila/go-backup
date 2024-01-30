package backup

type BackupInterface interface {
	BackupProcess()
}

type backup struct {
	CronExpression       string
	WhatToBackup         []string
	ExcludeExtensions    []string
	ExcludeFiles         []string
	ParsedCronExpression cron
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
    b.parseCron()
	//var timer, parseError = parse(config.Node.WhenToBackup)
	//if parseError != nil {
	//	console.Fatal("The given cron expression is invalid!")
	//}
	//var timerNext = timer.next(time.Now())
	//var duration = timerNext.Unix() - time.Now().Unix()
	//console.Debug(time.Now().String() + " " + timerNext.String())

	//var timeSignal = time.After(time.Duration(time.Duration(duration).Seconds()) * time.Second)
	//select {
	//case <-timeSignal:
	//	// TODO: implement backup
	//	console.Debug("TODO backup")
	//	//BackupProcess()
	//}

	return
}
