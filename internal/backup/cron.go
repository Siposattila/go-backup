package backup

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type bitset uint64

type bound struct {
	min int
	max int
}

type translater map[string]int

type Schedule struct {
	Minute   bitset
	Hour     bitset
	Dom      bitset
	Month    bitset
	Dow      bitset
	Location *time.Location
}

const bitsetStar = 1<<64 - 1

var boundMinute = bound{0, 59}
var boundHour = bound{0, 24}
var boundDOM = bound{1, 31}
var boundMonth = bound{1, 12}
var boundDOW = bound{0, 6}

var translaterMonth = translater{
	"JAN": 1,
	"FEB": 2,
	"MAR": 3,
	"APR": 4,
	"MAY": 5,
	"JUN": 6,
	"JUL": 7,
	"AUG": 8,
	"SEP": 9,
	"OCT": 10,
	"NOV": 11,
	"DEC": 12,
}

var translaterDay = translater{
	"SUN": 0,
	"MON": 1,
	"TUE": 2,
	"WED": 3,
	"THU": 4,
	"FRI": 5,
	"SAT": 6,
}

// ParseInLocation parse the expression in the location and
// returns a new schedule representing the given spec.
// It returns an error when loading the location is failed or
// the syntax of the expression is wrong.
func parseInLocation(expr string, locName string) (*Schedule, error) {
	loc, err := time.LoadLocation(locName)
	if err != nil {
		return nil, err
	}

	schdule, err := parse(expr)
	if err != nil {
		return nil, err
	}

	schdule.Location = loc
	return schdule, nil
}

// Parse parses the expression and returns a new schedule representing the given spec.
// And the default location of a schedule is "UTC".
// It returns an error when the syntax of expression is wrong.
func parse(expr string) (*Schedule, error) {
	err := verifyExpr(expr)
	if err != nil {
		return nil, err
	}

	var (
		minute, hour, dom, month, dow bitset
	)

	fields := strings.Fields(strings.TrimSpace(expr))

	if minute, err = parseField(fields[0], boundMinute, translater{}); err != nil {
		return nil, err
	}

	if hour, err = parseField(fields[1], boundHour, translater{}); err != nil {
		return nil, err
	}

	if dom, err = parseField(fields[2], boundDOM, translater{}); err != nil {
		return nil, err
	}

	if month, err = parseField(fields[3], boundMonth, translaterMonth); err != nil {
		return nil, err
	}

	if dow, err = parseField(fields[4], boundDOW, translaterDay); err != nil {
		return nil, err
	}

	return &Schedule{
		Minute: minute,
		Hour:   hour,
		Dom:    dom,
		Month:  month,
		Dow:    dow,
	}, nil
}

// parseField returns an int with the bits set representing all of the times that
// the field represents or error parsing field value.
func parseField(field string, b bound, t translater) (bitset, error) {
	var bitsets bitset = 0

	// Split with "," (OR).
	fieldexprs := strings.Split(field, ",")
	for _, fieldexpr := range fieldexprs {
		b, err := parseFieldExpr(fieldexpr, b, t)
		if err != nil {
			return 0, err
		}

		bitsets = bitsets | b
	}

	return bitsets, nil
}

// parseFieldExpr returns the bits indicated by the given expression:
//
//	number | number "-" number [ "/" number ]
func parseFieldExpr(fieldexpr string, b bound, t translater) (bitset, error) {
	// Replace "*" into "min-max".
	newexpr := strings.Replace(fieldexpr, "*", fmt.Sprintf("%d-%d", b.min, b.max), 1)

	rangeAndStep := strings.Split(newexpr, "/")
	if !(len(rangeAndStep) == 1 || len(rangeAndStep) == 2) {
		return 0, fmt.Errorf("Failed to parse the expr '%s', too many '/'", fieldexpr)
	}

	hasStep := len(rangeAndStep) == 2

	// Parse the range, first.
	var (
		begin, end int
	)
	{
		lowAndHigh := strings.Split(rangeAndStep[0], "-")
		if !(len(lowAndHigh) == 1 || len(lowAndHigh) == 2) {
			return 0, fmt.Errorf("Failed to parse the expr '%s', too many '-'", fieldexpr)
		}

		low, err := parseInt(lowAndHigh[0], t)
		if err != nil {
			return 0, fmt.Errorf("Failed to parse the expr '%s': %w", fieldexpr, err)
		}

		begin = low

		// Special handling: "N/step" means "N-max/step".
		if len(lowAndHigh) == 1 && hasStep {
			end = b.max
		} else if len(lowAndHigh) == 1 && !hasStep {
			end = low
		} else if len(lowAndHigh) == 2 {
			high, err := parseInt(lowAndHigh[1], t)
			if err != nil {
				return 0, fmt.Errorf("Failed to parse the expr '%s': %w", fieldexpr, err)
			}

			end = high
		}
	}

	// Parse the step, second.
	step := 1
	if hasStep {
		var err error
		if step, err = strconv.Atoi(rangeAndStep[1]); err != nil {
			return 0, fmt.Errorf("Failed to parse the expr '%s': %w", fieldexpr, err)
		}
	}

	return buildBitset(begin, end, step), nil
}

func parseInt(s string, t translater) (int, error) {
	if i, err := strconv.Atoi(s); err == nil {
		return i, nil
	}

	i, ok := t[strings.ToUpper(s)]
	if !ok {
		return 0, fmt.Errorf("'%s' is out of reserved words", s)
	}

	return i, nil
}

func buildBitset(min, max, step int) bitset {
	var b bitset

	for i := min; i <= max; i += step {
		b = b | (1 << i)
	}

	return b
}

func verifyExpr(expr string) error {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return fmt.Errorf("The length of fields must be five.")
	}

	return nil
}

// Next returns the next time matched with the expression.
func (s *Schedule) next(t time.Time) time.Time {
	loc := time.UTC
	if s.Location != nil {
		loc = s.Location
	}

	origLoc := t.Location()
	t = t.In(loc)

	added := false

	// Start at the earliest possible time (the upcoming second).
	t = t.Add(1*time.Minute - time.Duration(t.Nanosecond())*time.Nanosecond)

	yearLimit := t.Year() + 5

L:
	if t.Year() > yearLimit {
		return time.Time{}
	}

	// Find the first applicable month.
	// If it's this month, then do nothing.
	year := t.Year()
	for 1<<uint(t.Month())&s.Month == 0 {
		// If we have to add a month, reset the other parts to 0.
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
		}

		t = t.AddDate(0, 1, 0)

		if t.Year() != year {
			goto L
		}
	}

	// Now get a day in that month.
	month := t.Month()
	for !dayMatches(s, t) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		}

		t = t.AddDate(0, 0, 1)

		if t.Month() != month {
			goto L
		}
	}

	day := t.Day()
	for 1<<uint(t.Hour())&s.Hour == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
		}

		t = t.Add(1 * time.Hour)

		if t.Day() != day {
			goto L
		}
	}

	hour := t.Hour()
	for 1<<uint(t.Minute())&s.Minute == 0 {
		if !added {
			added = true
			t = t.Truncate(time.Minute)
		}

		t = t.Add(1 * time.Minute)

		if t.Hour() != hour {
			goto L
		}
	}

	return t.
		Truncate(time.Minute).
		In(origLoc)
}

// Next returns the previous time matched with the expression.
func (s *Schedule) prev(t time.Time) time.Time {
	loc := time.UTC
	if s.Location != nil {
		loc = s.Location
	}

	origLoc := t.Location()
	t = t.In(loc)

	subtracted := false

	// Start at the earliest possible time (the upcoming second).
	t = t.Add(-1*time.Minute + time.Duration(t.Nanosecond())*time.Nanosecond)

	yearLimit := t.Year() - 5

L:
	if t.Year() < yearLimit {
		return time.Time{}
	}

	year := t.Year()
	for 1<<uint(t.Month())&s.Month == 0 {
		// If we have to add a month, reset with the next month before.
		if !subtracted {
			subtracted = true
			t = time.Date(t.Year(), t.Month()+1, 0, 23, 59, 0, 0, loc)
		}

		// Change the time into the last day of the previous month.
		// Note that AddDate(0, -1, 0) has a bug by the normalization.
		// E.g) time.Date(2021, 6, 0, 23, 59, 59, 0, time.UTC).AddDate(0, -1, 0)
		t = time.Date(t.Year(), t.Month(), 0, 23, 59, 0, 0, loc)

		if t.Year() != year {
			goto L
		}
	}

	// Now get a day in that month.
	month := t.Month()
	for !dayMatches(s, t) {
		if !subtracted {
			subtracted = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 0, 0, loc)
		}

		t = t.AddDate(0, 0, -1)

		if t.Month() != month {
			goto L
		}
	}

	day := t.Day()
	for 1<<uint(t.Hour())&s.Hour == 0 {
		if !subtracted {
			subtracted = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 0, 0, loc)
		}

		t = t.Add(-1 * time.Hour)

		if t.Day() != day {
			goto L
		}
	}

	hour := t.Hour()
	for 1<<uint(t.Minute())&s.Minute == 0 {
		if !subtracted {
			subtracted = true
			t = t.Truncate(-time.Minute)
		}

		t = t.Add(-1 * time.Minute)

		if t.Hour() != hour {
			goto L
		}
	}

	return t.
		Truncate(time.Minute).
		In(origLoc)
}

// dayMatches returns true if the schedule's day-of-week and day-of-month
// restrictions are satisfied by the given time.
func dayMatches(s *Schedule, t time.Time) bool {
	var (
		domMatch bool = 1<<uint(t.Day())&s.Dom > 0
		dowMatch bool = 1<<uint(t.Weekday())&s.Dow > 0
	)
	if s.Dom&bitsetStar > 0 || s.Dow&bitsetStar > 0 {
		return domMatch && dowMatch
	}
	return domMatch || dowMatch
}
