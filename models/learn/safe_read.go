package learn

import (
	"time"
)

const (
	defaultSleepDuration = 100 * time.Millisecond
)

type SafeReadWrite struct {
	I    int
	Intc chan int
}

func (s *SafeReadWrite) Write(i int) {
	s.Intc <- i
}

func (s *SafeReadWrite) Close() {
	close(s.Intc)
}

func (s *SafeReadWrite) Read() (val int) {
	select {
	case v, ok := <-s.Intc:
		// If Channel is Not Closed Update I
		if ok {
			// Update New Value in Cache
			s.I = v
		}
		// Serve Updated I
		val = s.I
	default:
		// Runs if no other event is ready run
		time.Sleep(defaultSleepDuration)
		// Serve Cached I
		val = s.I
	}
	return
}
