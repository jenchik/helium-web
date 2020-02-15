package core

import "time"

var startTime = time.Now().UTC()

func Uptime() time.Duration {
	return time.Since(startTime)
}
