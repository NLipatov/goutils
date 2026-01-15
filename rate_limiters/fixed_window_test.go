package rate_limiters

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestFixedWindow_Allow_WithinLimit(t *testing.T) {
	limiter := NewFixedWindow(100*time.Millisecond, 3)

	for i := 0; i < 3; i++ {
		if !limiter.Allow() {
			t.Fatalf("expected Allow() to return true for request #%d", i+1)
		}
	}
}

func TestFixedWindow_Allow_OverLimit(t *testing.T) {
	limiter := NewFixedWindow(100*time.Millisecond, 2)

	_ = limiter.Allow()
	_ = limiter.Allow()
	if limiter.Allow() {
		t.Fatal("expected Allow() to return false when over the limit")
	}
}

func TestFixedWindow_Allow_WindowReset(t *testing.T) {
	limiter := NewFixedWindow(30*time.Millisecond, 1)

	if !limiter.Allow() {
		t.Fatal("expected Allow() to return true for first request")
	}
	if limiter.Allow() {
		t.Fatal("expected Allow() to return false for second request in same window")
	}

	time.Sleep(40 * time.Millisecond) // wait for window to reset

	if !limiter.Allow() {
		t.Fatal("expected Allow() to return true after window reset")
	}
}

func TestFixedWindow_Allow_StrictNoOvershoot(t *testing.T) {
	limiter := NewFixedWindow(100*time.Millisecond, 1000)

	var allowed, denied int64
	var wg sync.WaitGroup

	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				atomic.AddInt64(&allowed, 1)
			} else {
				atomic.AddInt64(&denied, 1)
			}
		}()
	}
	wg.Wait()

	if allowed != 1000 {
		t.Errorf("allowed = %d, want exactly 1000 (no overshoot expected)", allowed)
	}
	if denied != 1000 {
		t.Errorf("denied = %d, want exactly 1000", denied)
	}
}

func TestFixedWindow_WindowStartsNow(t *testing.T) {
	limiter := NewFixedWindow(100*time.Millisecond, 1)
	time.Sleep(10 * time.Millisecond)
	if !limiter.Allow() {
		t.Fatal("expected Allow() to return true on first request")
	}
}
