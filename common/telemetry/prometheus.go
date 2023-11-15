package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func WriteToFile(filePath string, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			prometheus.WriteToTextfile(filePath, prometheus.DefaultGatherer)
		}
	}
}
