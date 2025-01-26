package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
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

var (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	gray   = "\033[37m"
	white  = "\033[97m"
)

var instance Logger

func GetLogger() Logger {
	if instance == nil {
		instance = newLogger("gobkup.log")
	}

	return instance
}

func newLogger(logFileName string) Logger {
	logger := &logger{}
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	multi := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multi)

	return logger
}

func colorize(color string, label string) string {
	return color + label + reset
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
	case reflect.Slice:
		if value.Type() == reflect.TypeOf([]byte(nil)) {
			resultMessage = string(value.Bytes()[:])
		}
	case reflect.Interface:
		if value.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			resultMessage = value.String()
		}
	}

	log.Println(label + " " + resultMessage)
}

func (l *logger) Warning(message ...any) {
	for _, value := range message {
		writeMessageToLog(colorize(yellow, "[Warning]"), value)
	}
}

func (l *logger) Error(message ...any) {
	for _, value := range message {
		writeMessageToLog(colorize(red, "[Error]"), value)
	}
}

func (l *logger) Fatal(message ...any) {
	for _, value := range message {
		writeMessageToLog(colorize(red, "[Fatal] "), value)
	}
	log.Fatal()
}

func (l *logger) Success(message ...any) {
	for _, value := range message {
		writeMessageToLog(colorize(green, "[Success]"), value)
	}
}

func (l *logger) Debug(message ...any) {
	for _, value := range message {
		writeMessageToLog(colorize(blue, "[Debug]"), value)
	}
}

func (l *logger) Normal(message ...any) {
	for _, value := range message {
		writeMessageToLog(colorize(white, "[Log]"), value)
	}
}
