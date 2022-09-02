package utils

import "time"

func CombineDateAndTime(date time.Time, t time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		time.UTC,
	)
}
