package apps

import (
	"time"
)

// Bod returns the time at the Begninning of Day
func Bod() time.Time {
	t := time.Now()
	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// FormatDate takes a time and returns a nicely formatted date,
// complete with day suffix (1st, 2nd, 3rd etc.)
//
// hattip: https://stackoverflow.com/a/28890625
func FormatDate(t time.Time) string {
	suffix := "th"
	switch t.Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return t.Format("Monday 2" + suffix + " January, 2006")
}
