package auth

import (
	"sync"
	"time"
)

var (
	stateMap = make(map[string]struct{})
	mu       sync.Mutex
)

func StartStateMapCleaner(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			mu.Lock()
			stateMap = make(map[string]struct{})
			mu.Unlock()
		}
	}()
}
