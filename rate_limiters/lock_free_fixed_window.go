package rate_limiters

import (
	"sync/atomic"
	"time"
)

type LockFreeFixedWindow struct {
	duration time.Duration
	start    atomic.Int64
	rate     atomic.Int64
	maxRate  int64
}

func NewLockFreeFixedWindow(
	periodDuration time.Duration,
	maxRate int64,
) *LockFreeFixedWindow {
	rl := &LockFreeFixedWindow{
		duration: periodDuration,
		maxRate:  maxRate,
	}
	rl.start.Store(time.Now().UnixNano())
	return rl
}

func (l *LockFreeFixedWindow) Allow() bool {
	if time.Since(time.Unix(0, l.start.Load())) >= l.duration {
		l.rate.Store(0)
		l.start.Store(time.Now().UnixNano())
	}

	if l.rate.Add(1) > l.maxRate {
		return false
	}

	return true
}
