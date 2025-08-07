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

func (f *LockFreeFixedWindow) Allow() bool {
	if time.Since(time.Unix(0, f.start.Load())) >= f.duration {
		f.rate.Store(0)
		f.start.Store(time.Now().UnixNano())
	}

	if f.rate.Add(1) > f.maxRate {
		return false
	}

	return true
}
