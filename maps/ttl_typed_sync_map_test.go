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
	m := NewTtlTypedSyncMap[int, string](ctx, 100*time.Millisecond)
	m.sanitizeInterval = 10 * time.Millisecond

	m.Store(1, "foo")
	v, ok := m.Load(1)
	if !ok || v != "foo" {
		t.Fatalf("want 'foo', got %v", v)
	}
}

func TestTtlTypedSyncMap_ExpireOnLoad(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := NewTtlTypedSyncMap[int, string](ctx, 30*time.Millisecond)
	m.sanitizeInterval = 10 * time.Millisecond

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
	m := NewTtlTypedSyncMap[int, string](ctx, time.Second)
	m.sanitizeInterval = 10 * time.Millisecond

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
	m := NewTtlTypedSyncMap[int, string](ctx, time.Second)
	m.sanitizeInterval = 10 * time.Millisecond

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
	m := NewTtlTypedSyncMap[int, string](ctx, 20*time.Millisecond)
	m.sanitizeInterval = 5 * time.Millisecond

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
	m := NewTtlTypedSyncMap[int, string](ctx, 100*time.Millisecond)
	m.sanitizeInterval = 10 * time.Millisecond

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
	m := NewTtlTypedSyncMap[int, string](ctx, exp)
	m.sanitizeInterval = 1 * time.Millisecond

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
	m := NewTtlTypedSyncMap[int, string](ctx, exp)
	m.sanitizeInterval = 1 * time.Millisecond

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
	m := NewTtlTypedSyncMap[int, string](ctx, 10*time.Millisecond)
	m.sanitizeInterval = 1 * time.Millisecond
	cancel()
	// Direct call to sanitize should return immediately without panic or blocking
	m.sanitize()
}
