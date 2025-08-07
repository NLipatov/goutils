package rate_limiters

import (
	"sync/atomic"
	"time"
)

// LockFreeFixedWindow implements a lock-free fixed window rate limiter.
// Allows a small overshoot under concurrent load due to lock-free design.
type LockFreeFixedWindow struct {
	window      time.Duration // Duration of each rate limit window
	windowStart atomic.Int64  // Start time of the current window (UnixNano)
	counter     atomic.Int64  // Requests counted in the current window
	limit       int64         // Maximum requests allowed per window
}

// NewLockFreeFixedWindow creates a new lock-free fixed window rate limiter.
func NewLockFreeFixedWindow(
	window time.Duration,
	limit int64,
) *LockFreeFixedWindow {
	rl := &LockFreeFixedWindow{
		window: window,
		limit:  limit,
	}
	rl.windowStart.Store(time.Now().UnixNano())
	return rl
}

// Allow returns true if the request is allowed in the current window, false otherwise.
// May allow a small number of requests above the limit under concurrent load.
func (l *LockFreeFixedWindow) Allow() bool {
	if time.Since(time.Unix(0, l.windowStart.Load())) >= l.window {
		l.counter.Store(0)
		l.windowStart.Store(time.Now().UnixNano())
	}

	if l.counter.Add(1) > l.limit {
		return false
	}

	return true
}
