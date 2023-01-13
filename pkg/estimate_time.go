package utils

import "time"

func EstimateExecutionTime(f func()) (t time.Duration) {
	start := time.Now()
	f()
	return time.Since(start)
}
