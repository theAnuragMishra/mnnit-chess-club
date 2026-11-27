package ratelimit

import (
	"testing"
	"time"
)

func TestLimiterRefillsAndBoundsKeys(t *testing.T) {
	limiter := New[string](2, 2, 1)
	now := time.Now()
	for _, check := range []struct {
		key   string
		at    time.Time
		retry time.Duration
	}{
		{"first", now, 0},
		{"first", now, 0},
		{"first", now, 500 * time.Millisecond},
		{"second", now, 0},
		{"third", now, 0},
		{"fourth", now, 500 * time.Millisecond},
		{"first", now.Add(250 * time.Millisecond), 250 * time.Millisecond},
		{"first", now.Add(500 * time.Millisecond), 0},
		{"second", now.Add(6 * time.Minute), 0},
	} {
		if retry := limiter.take(check.key, check.at); retry != check.retry {
			t.Fatalf("Take(%q) returned %s, want %s", check.key, retry, check.retry)
		}
		if len(limiter.buckets) > 1 {
			t.Fatal("limiter exceeded its key capacity")
		}
	}
	if _, ok := limiter.buckets["second"]; !ok {
		t.Fatal("a new client did not get its own bucket after cleanup")
	}
}
