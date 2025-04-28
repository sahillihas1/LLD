package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter interface {
	Allow() bool
}

type SlidingWindowLimiter struct {
	mu         sync.Mutex
	windowSize time.Duration
	limit      int
	timestamps []time.Time
}

func NewSlidingWindowLimiter(limit int, windowSize time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		windowSize: windowSize,
		limit:      limit,
		timestamps: []time.Time{},
	}
}

func (r *SlidingWindowLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.windowSize)

	validTimestamps := r.timestamps[:0]
	for _, ts := range r.timestamps {
		if ts.After(cutoff) {
			validTimestamps = append(validTimestamps, ts)
		}
	}
	r.timestamps = validTimestamps

	if len(r.timestamps) < r.limit {
		r.timestamps = append(r.timestamps, now)
		return true
	}

	return false
}

type TokenBucketLimiter struct {
	mu         sync.Mutex
	tokens     int
	capacity   int
	rate       int // tokens per second
	lastRefill time.Time
}

func NewTokenBucketLimiter(rate int, capacity int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

func (t *TokenBucketLimiter) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(t.lastRefill).Seconds()
	newTokens := int(elapsed * float64(t.rate))

	if newTokens > 0 {
		t.tokens = min(t.capacity, t.tokens+newTokens)
		t.lastRefill = now
	}

	if t.tokens > 0 {
		t.tokens--
		return true
	}

	return false
}
