package rate_limiters

import (
	"sync"
	"time"
)

// FixedWindow implements a thread-safe fixed window rate limiter using a mutex.
// This implementation strictly enforces the rate limit with zero overshoot.
type FixedWindow struct {
	window      time.Duration // Duration of each rate limit window
	windowStart time.Time     // Start time of the current window
	counter     int64         // Requests counted in the current window
	limit       int64         // Maximum requests allowed per window
	mutex       sync.Mutex
}

// NewFixedWindow creates a new mutex-based fixed window rate limiter.
func NewFixedWindow(
	window time.Duration,
	limit int64,
) *FixedWindow {
	return &FixedWindow{
		window:      window,
		windowStart: time.Now(),
		limit:       limit,
	}
}

// Allow returns true if the request is allowed in the current window, false otherwise.
// Strictly enforces the rate limit. This method is thread-safe.
func (l *FixedWindow) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if time.Now().Sub(l.windowStart) > l.window {
		l.counter = 0
		l.windowStart = time.Now()
	}

	l.counter++
	if l.counter > l.limit {
		return false
	}

	return true
}
