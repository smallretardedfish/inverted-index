package pkg

import "time"

func EstimateExecutionTime(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}
