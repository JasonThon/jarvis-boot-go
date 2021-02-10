package utils

import (
	"errors"
	"github.com/robfig/cron"
)

func CheckCron(pattern string) error {
	if len(pattern) == 0 {
		return errors.New("cron is NOT specified")
	}
	_, err := cron.Parse(pattern)

	return err
}
