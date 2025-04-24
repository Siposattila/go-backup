package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"gopkg.in/natefinch/lumberjack.v2"
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

func writeMessageToLog(label string, message any) {
	var resultMessage string
	value := reflect.ValueOf(message)
	switch value.Kind() {
	case reflect.String:
		resultMessage = value.String()
	case reflect.Bool:
		resultMessage = fmt.Sprintf("%t", value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		resultMessage = fmt.Sprintf("%d", value.Int())
	case reflect.Float32, reflect.Float64:
		resultMessage = fmt.Sprintf("%f", value.Float())
	case reflect.Uint32, reflect.Uint64:
		resultMessage = fmt.Sprintf("%d", value.Uint())
	case reflect.Slice:
		if value.Type() == reflect.TypeOf([]byte(nil)) {
			resultMessage = string(value.Bytes()[:])
		}
	}

	log.Println(label + " " + resultMessage)
}

func (l *logger) Warning(message ...any) {
	for _, value := range message {
		writeMessageToLog("[Warning]", value)
	}
}

func (l *logger) Error(message ...any) {
	for _, value := range message {
		writeMessageToLog("[Error]", value)
	}
}

func (l *logger) Fatal(message ...any) {
	for _, value := range message {
		writeMessageToLog("[Fatal] ", value)
	}
	log.Fatal()
}

func (l *logger) Success(message ...any) {
	for _, value := range message {
		writeMessageToLog("[Success]", value)
	}
}

func (l *logger) Debug(message ...any) {
	for _, value := range message {
		writeMessageToLog("[Debug]", value)
	}
}

func (l *logger) Normal(message ...any) {
	for _, value := range message {
		writeMessageToLog("[Log]", value)
	}
}
