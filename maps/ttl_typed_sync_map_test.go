package maps

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestTtlTypedSyncMap_StoreLoad(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exp := 10 * time.Millisecond
	sanitizeInterval := time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)

	m.Store(1, "foo")
	v, ok := m.Load(1)
	if !ok || v != "foo" {
		t.Fatalf("want 'foo', got %v", v)
	}
}

func TestTtlTypedSyncMap_ExpireOnLoad(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exp := 30 * time.Millisecond
	sanitizeInterval := 10 * time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)

	m.Store(2, "bar")
	time.Sleep(40 * time.Millisecond)
	_, ok := m.Load(2)
	if ok {
		t.Fatalf("expected expired value")
	}
}

func TestTtlTypedSyncMap_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exp := time.Second
	sanitizeInterval := 10 * time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)

	m.Store(3, "baz")
	m.Delete(3)
	_, ok := m.Load(3)
	if ok {
		t.Fatalf("expected deleted value")
	}
}

func TestTtlTypedSyncMap_Len(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exp := time.Second
	sanitizeInterval := 10 * time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)

	if m.Len() != 0 {
		t.Fatalf("expected 0")
	}
	m.Store(4, "qux")
	if m.Len() != 1 {
		t.Fatalf("expected 1")
	}
	m.Delete(4)
	if m.Len() != 0 {
		t.Fatalf("expected 0 after delete")
	}
}

func TestTtlTypedSyncMap_SanitizeJanitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exp := 20 * time.Millisecond
	sanitizeInterval := 5 * time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)

	m.Store(5, "e")
	time.Sleep(30 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	if m.Len() != 0 {
		t.Fatalf("expected map to be empty after sanitize, got %d", m.Len())
	}
}

func TestTtlTypedSyncMap_Concurrency(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exp := 10 * time.Millisecond
	sanitizeInterval := 1 * time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)

	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			m.Store(k, "x")
			m.Load(k)
			m.Delete(k)
		}(i)
	}
	wg.Wait()
	if m.Len() != 0 {
		t.Fatalf("expected len 0 after concurrent store/delete, got %d", m.Len())
	}
}

func TestTtlTypedSyncMap_Range(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exp := 10 * time.Millisecond
	sanitizeInterval := 1 * time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)

	// Store two keys; key1 will expire
	m.Store(1, "one")
	m.Store(2, "two")

	// expire key 1, remain key 2
	time.Sleep(time.Millisecond * 5)
	// load will update ttl for key 2
	m.Load(2)
	time.Sleep(time.Millisecond * 5)

	collected := make(map[int]string)
	m.Range(func(k int, v string) bool {
		collected[k] = v
		return true
	})

	// key1 expired, key2 should remain
	if _, ok := collected[1]; ok {
		t.Fatalf("expected expired key 1 to be skipped in Range")
	}
	if val, ok := collected[2]; !ok || val != "two" {
		t.Fatalf("expected key2=\"\"two\"\", got %v (ok=%v)", val, ok)
	}

	// Test early stop
	times := 0
	m.Range(func(k int, v string) bool {
		times++
		return false
	})
	if times != 1 {
		t.Fatalf("expected Range to stop after 1 iteration, got %d", times)
	}
}

func TestTtlTypedSyncMap_TTLResetOnLoad(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exp := 20 * time.Millisecond
	sanitizeInterval := time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)

	m.Store(10, "ten")
	// sleep less than exp
	time.Sleep(10 * time.Millisecond)
	// Load resets TTL
	if v, ok := m.Load(10); !ok || v != "ten" {
		t.Fatalf("expected Load to reset TTL and return (\"ten\",true), got (%v,%v)", v, ok)
	}
	// sleep more than exp from original store but less than after reset
	time.Sleep(15 * time.Millisecond)
	if v, ok := m.Load(10); !ok || v != "ten" {
		t.Fatalf("expected value alive after TTL reset, got (%v,%v)", v, ok)
	}
}

// Test that sanitize stops on context cancel
func TestTtlTypedSyncMap_SanitizeStopsOnCancelDirect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	exp := 10 * time.Millisecond
	sanitizeInterval := 1 * time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)
	cancel()
	// Direct call to sanitize should return immediately without panic or blocking
	m.sanitize()
}

func TestTtlTypedSyncMap_RangeDeletesMissing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exp := 50 * time.Millisecond
	sanitizeInterval := 10 * time.Millisecond
	m := NewTtlTypedSyncMap[int, string](ctx, exp, sanitizeInterval)

	// Store and then remove underlying before Range
	m.Store(42, "value")
	// Simulate missing underlying entry: delete directly on inner map
	m.Delete(42)

	collected := make(map[int]string)
	m.Range(func(k int, v string) bool {
		collected[k] = v
		return true
	})
	// collected must be empty, and exp map entry removed
	if len(collected) != 0 {
		t.Fatalf("expected no entries collected, got %v", collected)
	}
	// underlying Len should be zero
	if m.Len() != 0 {
		t.Fatalf("expected Len()=0 after missing cleanup, got %d", m.Len())
	}
}

func TestNewTtlTypedSyncMap_Defaults(t *testing.T) {
	// Passing non-positive durations should fall back to defaults
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// expDuration and sanitizeInterval are <= 0
	m := NewTtlTypedSyncMap[string, int](ctx, 0, -10*time.Second)

	// Immediately after Store the entry should be available
	m.Store("foo", 42)
	if v, ok := m.Load("foo"); !ok || v != 42 {
		t.Fatalf("expected to load immediately, got ok=%v, v=%d", ok, v)
	}

	// After ~1 second (default expDuration) the entry should expire
	time.Sleep(1100 * time.Millisecond)
	if _, ok := m.Load("foo"); ok {
		t.Fatal("expected entry to be expired after default 1s TTL")
	}
}

func TestNewTtlTypedSyncMap_CustomIntervals(t *testing.T) {
	// Passing custom intervals should use them
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exp := 200 * time.Millisecond
	sanitizeInterval := 50 * time.Millisecond
	m := NewTtlTypedSyncMap[string, int](ctx, exp, sanitizeInterval)

	// Store an entry
	m.Store("bar", 100)
	if v, ok := m.Load("bar"); !ok || v != 100 {
		t.Fatalf("expected to load immediately, got ok=%v, v=%d", ok, v)
	}

	// Wait just beyond expDuration but before sanitizeInterval triggers
	time.Sleep(exp + 20*time.Millisecond)
	// Although the sanitize goroutine may not have run yet,
	// Load itself should remove the expired key.
	if _, ok := m.Load("bar"); ok {
		t.Fatal("expected entry to be expired after custom TTL")
	}
}

func TestRangeProlongatesTTL(t *testing.T) {
	// Verify that Range prolongs the TTL on each access
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exp := 10 * time.Millisecond
	sanitizeInterval := 5 * time.Millisecond
	m := NewTtlTypedSyncMap[string, int](ctx, exp, sanitizeInterval)

	m.Store("baz", 7)
	// Call Range a few times before the original expDuration elapses
	for i := 0; i < 3; i++ {
		time.Sleep(exp / 2)
		called := false
		m.Range(func(k string, v int) bool {
			called = true
			if k != "baz" || v != 7 {
				t.Fatalf("unexpected key/value %s/%d", k, v)
			}
			return true
		})
		if !called {
			t.Fatal("expected Range to see the entry before expiration")
		}
	}

	// Now without further Range calls, wait for TTL to elapse
	time.Sleep(exp + 10*time.Millisecond)
	var seen bool
	m.Range(func(k string, v int) bool {
		seen = true
		return true
	})
	if seen {
		t.Fatal("expected no entries after TTL expired")
	}
}

func TestTtlTypedSyncMap_Range_RemovesExpiredEntries(t *testing.T) {
	// expDuration is very short, sanitizeInterval is long so sanitize()
	// won't remove entries before we call Range.
	ttl := NewTtlTypedSyncMap[string, int](
		context.Background(),
		10*time.Millisecond,
		1*time.Second,
	)
	ttl.Store("expiredKey", 100)

	// wait until after expiration
	time.Sleep(20 * time.Millisecond)

	// Range should see the entry as expired, delete it, and continue
	called := false
	ttl.Range(func(k string, v int) bool {
		called = true
		return true
	})

	if called {
		t.Errorf("expected Range not to call function on expired entries")
	}
	if got := ttl.Len(); got != 0 {
		t.Errorf("expected Len() == 0 after Range removed expired entry, got %d", got)
	}
}

func TestTtlTypedSyncMap_Sanitize_RemovesExpiredEntries(t *testing.T) {
	// expDuration and sanitizeInterval both short
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ttl := NewTtlTypedSyncMap[string, int](
		ctx,
		10*time.Millisecond,
		10*time.Millisecond,
	)
	ttl.Store("expiredKey", 200)

	// wait enough time for the key to expire and for at least one sanitize tick
	time.Sleep(35 * time.Millisecond)

	// sanitize goroutine should have deleted the expired key
	if got := ttl.Len(); got != 0 {
		t.Errorf("expected Len() == 0 after sanitize removed expired entry, got %d", got)
	}
}
