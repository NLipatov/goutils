package maps

import (
	"sync"
	"sync/atomic"
)

type TypedSyncMap[K comparable, V any] struct {
	m     *sync.Map
	count atomic.Int64
}

func NewTypedSyncMap[K comparable, V any]() *TypedSyncMap[K, V] {
	return &TypedSyncMap[K, V]{
		m:     new(sync.Map),
		count: atomic.Int64{},
	}
}

func NewFromSyncMap[K comparable, V any](m *sync.Map) *TypedSyncMap[K, V] {
	t := &TypedSyncMap[K, V]{m: m}
	var cnt int64
	m.Range(func(_, _ any) bool {
		cnt++
		return true
	})
	t.count.Store(cnt)
	return t
}

func (t *TypedSyncMap[K, V]) Store(key K, value V) {
	_, loaded := t.m.LoadOrStore(key, value)
	if loaded {
		t.m.Store(key, value)
	} else {
		t.count.Add(1)
	}
}

func (t *TypedSyncMap[K, V]) Load(key K) (V, bool) {
	v, ok := t.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true
}

func (t *TypedSyncMap[K, V]) Delete(key K) {
	if _, existed := t.m.Load(key); existed {
		t.m.Delete(key)
		t.count.Add(-1)
	}
}

func (t *TypedSyncMap[K, V]) Len() int64 {
	return t.count.Load()
}

func (t *TypedSyncMap[K, V]) Range(f func(key K, value V) bool) {
	if t.m == nil || f == nil {
		return
	}

	t.m.Range(func(k, v any) bool {
		key := k.(K)
		val := v.(V)
		return f(key, val)
	})
}
