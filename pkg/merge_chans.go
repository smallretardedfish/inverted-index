package utils

import (
	"sync"
)

func MergeChannels[T any](chans ...chan T) chan T {
	res := make(chan T)

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chans))

		for _, ch := range chans {
			go func(c chan T) {
				for val := range c {
					res <- val
				}
			}(ch)
		}

		wg.Wait()
		close(res)
	}()

	return res
}
