package util

import "time"

func RunWithTicker(secondsWait int, callback func()) {
	ticker := time.NewTicker(time.Duration(secondsWait) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		callback()
	}
}
