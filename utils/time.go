package utils

import "time"

const (
	DefaultTimeFormat = "2006-01-02T15:04:05+08:00"
)

func NowString() string {
	return time.Now().Format(DefaultTimeFormat)
}

func Parse(timeString string) (time.Time, error) {
	return time.ParseInLocation(DefaultTimeFormat, timeString, time.Local)
}
