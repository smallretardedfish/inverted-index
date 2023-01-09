package utils

import "time"

func EstimateExecutionTime(f func()) (t time.Duration) {
	start := time.Now()
	defer func() {
		t = time.Since(start)
	}()

	f()
	return t
}
