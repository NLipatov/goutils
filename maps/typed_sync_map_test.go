package maps

import (
	"sync"
	"testing"
)

func TestTypedSyncMap_StoreAndLoad(t *testing.T) {
	m := NewTypedSyncMap[int, string]()
	m.Store(1, "a")
	v, ok := m.Load(1)
	if !ok || v != "a" {
		t.Fatalf("expected %v, got %v (ok=%v)", "a", v, ok)
	}
}

func TestTypedSyncMap_StoreOverwrite(t *testing.T) {
	m := NewTypedSyncMap[int, string]()
	m.Store(1, "a")
	m.Store(1, "b")
	v, ok := m.Load(1)
	if !ok || v != "b" {
		t.Fatalf("expected %v, got %v (ok=%v)", "b", v, ok)
	}
	if m.Len() != 1 {
		t.Fatalf("expected len 1, got %d", m.Len())
	}
}

func TestTypedSyncMap_Len(t *testing.T) {
	m := NewTypedSyncMap[int, string]()
	if m.Len() != 0 {
		t.Fatalf("expected len 0, got %d", m.Len())
	}
	m.Store(1, "a")
	m.Store(2, "b")
	if m.Len() != 2 {
		t.Fatalf("expected len 2, got %d", m.Len())
	}
}

func TestTypedSyncMap_Delete(t *testing.T) {
	m := NewTypedSyncMap[int, string]()
	m.Store(1, "a")
	m.Store(2, "b")
	m.Delete(1)
	_, ok := m.Load(1)
	if ok {
		t.Fatalf("expected not found")
	}
	if m.Len() != 1 {
		t.Fatalf("expected len 1, got %d", m.Len())
	}
	m.Delete(1)
	if m.Len() != 1 {
		t.Fatalf("double delete should not affect len, got %d", m.Len())
	}
}

func TestTypedSyncMap_ConcurrentAccess(t *testing.T) {
	m := NewTypedSyncMap[int, string]()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			m.Store(k, "v")
			m.Load(k)
			m.Delete(k)
		}(i)
	}
	wg.Wait()
	if m.Len() != 0 {
		t.Fatalf("expected len 0 after concurrent store/delete, got %d", m.Len())
	}
}

func TestNewFromSyncMap(t *testing.T) {
	sm := &sync.Map{}
	sm.Store(1, "a")
	sm.Store(2, "b")
	m := NewFromSyncMap[int, string](sm)
	if m.Len() != 2 {
		t.Fatalf("expected len 2, got %d", m.Len())
	}
	if v, ok := m.Load(1); !ok || v != "a" {
		t.Fatalf("expected Load(1) to return (\"a\", true), got (%v, %v)", v, ok)
	}
	if v, ok := m.Load(2); !ok || v != "b" {
		t.Fatalf("expected Load(2) to return (\"b\", true), got (%v, %v)", v, ok)
	}
	m.Store(3, "c")
	if m.Len() != 3 {
		t.Fatalf("expected len 3 after adding, got %d", m.Len())
	}
	m.Delete(1)
	if m.Len() != 2 {
		t.Fatalf("expected len 2 after deleting, got %d", m.Len())
	}
}

func TestTypedSyncMap_Range(t *testing.T) {
	m := NewTypedSyncMap[int, string]()
	m.Store(1, "a")
	m.Store(2, "b")

	collected := make(map[int]string)
	m.Range(func(k int, v string) bool {
		collected[k] = v
		return true
	})
	if len(collected) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(collected))
	}
	if collected[1] != "a" || collected[2] != "b" {
		t.Fatalf("expected collected values {1:\"a\",2:\"b\"}, got %v", collected)
	}

	times := 0
	m.Range(func(k int, v string) bool {
		times++
		return false
	})
	if times != 1 {
		t.Fatalf("expected Range to stop after first iteration, got %d iterations", times)
	}
}
