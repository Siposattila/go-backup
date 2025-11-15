package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	WARNING = "[Warning]"
	ERROR   = "[Error]"
	FATAL   = "[Fatal]"
	SUCCESS = "[Success]"
	DEBUG   = "[Debug]"
	NORMAL  = "[Log]"
)

type Logger interface {
	Warning(message ...any)
	Error(message ...any)
	Fatal(message ...any)
	Success(message ...any)
	Debug(message ...any)
	Normal(message ...any)
}

type logger struct{}

var instance Logger

func GetLogger() Logger {
	if instance == nil {
		instance = newLogger("gobkup.log")
	}

	return instance
}

func newLogger(logFileName string) Logger {
	logger := &logger{}
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    5, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
	}

	multi := io.MultiWriter(lumberjackLogger, os.Stdout)
	log.SetOutput(multi)

	return logger
}

func writeMessageToLog(label string, message string) {
	log.Printf("%s %s", label, message)
}

func processMessages(label string, message ...any) {
	var logMessage string
	for i, v := range message {
		if i > 0 {
			logMessage += "\n"
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.String:
			logMessage += value.String()
		case reflect.Bool:
			logMessage += fmt.Sprintf("%t", value.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			logMessage += fmt.Sprintf("%d", value.Int())
		case reflect.Float32, reflect.Float64:
			logMessage += fmt.Sprintf("%f", value.Float())
		case reflect.Uint32, reflect.Uint64:
			logMessage += fmt.Sprintf("%d", value.Uint())
		case reflect.Slice:
			if value.Type() == reflect.TypeOf([]byte(nil)) {
				logMessage += string(value.Bytes()[:])
			}
		}
	}
	writeMessageToLog(label, logMessage)
}

func (l *logger) Warning(message ...any) {
	processMessages(WARNING, message...)
}

func (l *logger) Error(message ...any) {
	processMessages(ERROR, message...)
}

func (l *logger) Fatal(message ...any) {
	processMessages(FATAL, message...)
	log.Fatal()
}

func (l *logger) Success(message ...any) {
	processMessages(SUCCESS, message...)
}

func (l *logger) Debug(message ...any) {
	processMessages(DEBUG, message...)
}

func (l *logger) Normal(message ...any) {
	processMessages(NORMAL, message...)
}
