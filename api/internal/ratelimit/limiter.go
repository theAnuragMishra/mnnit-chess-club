package ratelimit

import (
	"sync"
	"time"
)

const (
	cleanupInterval = time.Minute
	idleTimeout     = 5 * time.Minute
)

// Limiter keeps a token bucket per key and shares an overflow bucket when full.
type Limiter[K comparable] struct {
	mu          sync.Mutex
	perSecond   float64
	burst       float64
	maxKeys     int
	buckets     map[K]*bucket
	overflow    bucket
	nextCleanup time.Time
}

type bucket struct {
	tokens float64
	last   time.Time
}

func New[K comparable](perSecond float64, burst, maxKeys int) *Limiter[K] {
	return &Limiter[K]{
		perSecond: perSecond,
		burst:     float64(burst),
		maxKeys:   maxKeys,
		buckets:   make(map[K]*bucket),
	}
}

// Take consumes one token or returns how long to wait for a token.
func (l *Limiter[K]) Take(key K) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.take(key, time.Now())
}

func (l *Limiter[K]) take(key K, now time.Time) time.Duration {
	if !now.Before(l.nextCleanup) {
		for key, entry := range l.buckets {
			if now.Sub(entry.last) >= idleTimeout {
				delete(l.buckets, key)
			}
		}
		l.nextCleanup = now.Add(cleanupInterval)
	}

	entry := l.buckets[key]
	if entry == nil {
		entry = &l.overflow
		if len(l.buckets) < l.maxKeys {
			entry = &bucket{}
			l.buckets[key] = entry
		}
	}
	if entry.last.IsZero() {
		entry.tokens = l.burst
	} else {
		entry.tokens = min(l.burst, entry.tokens+now.Sub(entry.last).Seconds()*l.perSecond)
	}
	entry.last = now
	if entry.tokens < 1 {
		return max(time.Nanosecond, time.Duration((1-entry.tokens)/l.perSecond*float64(time.Second)))
	}
	entry.tokens--
	return 0
}
