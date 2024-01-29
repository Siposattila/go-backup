package backup

type BackupInterface interface {
    BackupProcess(cronExpression string)
}

type backup struct {
}

func NewBackup() BackupInterface {
    return backup{}
}

func (b backup) BackupProcess(cronExpression string) {
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
