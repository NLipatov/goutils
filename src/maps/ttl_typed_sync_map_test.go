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
