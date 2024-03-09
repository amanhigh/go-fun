package util

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ScheduleJob(secondsWait int, callback func(exit bool)) {
	//Ticker Based on Wait Time
	ticker := time.NewTicker(time.Duration(secondsWait) * time.Second)
	defer ticker.Stop()

	// Create a channel to listen for Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			callback(false)
		case <-sigChan:
			callback(true)
			return
		}
	}
}
