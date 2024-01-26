package console

import "log"

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

func colorize(color, s string) string {
	return color + s + reset
}

func Warning(s string) {
	log.Println(Yellow("[Warning]"), s)
}

func Error(s string) {
	log.Println(Red("[Error]"), s)
}

func Fatal(s string) {
	log.Fatal(Red("[Fatal] "), s)
}

func Success(s string) {
	log.Println(Green("[Success]"), s)
}

func Debug(s string) {
	log.Println(Blue("[Debug]"), s)
}

func Normal(s string) {
	log.Println("[Log]", s)
}

func Bold(s string) string {
	return colorize(bold, s)
}

func Red(s string) string {
	return colorize(red, s)
}

func Green(s string) string {
	return colorize(green, s)
}

func Yellow(s string) string {
	return colorize(yellow, s)
}

func Blue(s string) string {
	return colorize(blue, s)
}

func Purple(s string) string {
	return colorize(purple, s)
}

func Cyan(s string) string {
	return colorize(cyan, s)
}

func Gray(s string) string {
	return colorize(gray, s)
}

func White(s string) string {
	return colorize(white, s)
}
