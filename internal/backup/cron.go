package backup

import (
	"strconv"
	"strings"

	"github.com/Siposattila/gobkup/internal/console"
)

const (
	EVERY     = "*/"
	UNDEFINED = "*"
	ZERO      = "0"
)

type expression struct {
	Value   int8
	Special string
}

type cron struct {
	Minute     expression // minute (0-59)
	Hour       expression // hour (0 - 23)
	DayOfMonth expression // day of the month (1 - 31)
	Month      expression // month (1 - 12)
	DayOfWeek  expression // day of the week (0 - 6)
}

func newCron() cron {
	return cron{
		Minute:     expression{},
		Hour:       expression{},
		DayOfMonth: expression{},
		Month:      expression{},
		DayOfWeek:  expression{},
	}
}

// This is a very simple cron parser. It only handles very simple number based cron expressions.
//
// * * * * * (every minute)
//
// */3 * * * * (every 3 minutes)
//
// 3 * * * * (at 3 minutes past the hour)
func (b backup) parseCron() {
	var cron = cron{}
	var helper = strings.Split(b.CronExpression, " ")
	console.Debug(strings.Join(helper, " "))
	console.Debug(strconv.FormatInt(int64(cron.Minute.Value), 10) + " " + cron.Minute.Special)

    b.ParsedCronExpression = cron

    return
}
