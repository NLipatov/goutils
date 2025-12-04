package maps

import (
	"context"
	"sync"
	"time"
)

// TtlTypedSyncMap is a sliding-expiration TTL map:
// every successful Load prolongs item's expiration.
type TtlTypedSyncMap[K comparable, V any] struct {
	ctx              context.Context
	sanitizeInterval time.Duration
	expDuration      time.Duration
	mu               sync.Mutex
	items            map[K]ttlEntry[V]
}

type ttlEntry[V any] struct {
	value     V
	expiresAt time.Time
}

func NewTtlTypedSyncMap[K comparable, V any](
	ctx context.Context,
	expDuration time.Duration,
	sanitizeInterval time.Duration,
) *TtlTypedSyncMap[K, V] {
	if expDuration <= 0 {
		expDuration = time.Second
	}
	if sanitizeInterval <= 0 {
		sanitizeInterval = expDuration / 2
		if sanitizeInterval <= 0 {
			sanitizeInterval = time.Second / 2
		}
	}

	res := &TtlTypedSyncMap[K, V]{
		ctx:              ctx,
		expDuration:      expDuration,
		sanitizeInterval: sanitizeInterval,
		items:            make(map[K]ttlEntry[V]),
	}
	go res.sanitize()
	return res
}

func (t *TtlTypedSyncMap[K, V]) Store(key K, value V) {
	now := time.Now()
	t.mu.Lock()
	t.items[key] = ttlEntry[V]{
		value:     value,
		expiresAt: now.Add(t.expDuration),
	}
	t.mu.Unlock()
}

func (t *TtlTypedSyncMap[K, V]) Load(key K) (V, bool) {
	var zero V
	now := time.Now()

	t.mu.Lock()
	entry, ok := t.items[key]
	if !ok {
		t.mu.Unlock()
		return zero, false
	}

	if now.After(entry.expiresAt) {
		delete(t.items, key)
		t.mu.Unlock()
		return zero, false
	}

	// sliding TTL
	entry.expiresAt = now.Add(t.expDuration)
	t.items[key] = entry
	v := entry.value

	t.mu.Unlock()
	return v, true
}

func (t *TtlTypedSyncMap[K, V]) Delete(key K) {
	t.mu.Lock()
	delete(t.items, key)
	t.mu.Unlock()
}

func (t *TtlTypedSyncMap[K, V]) Len() int64 {
	t.mu.Lock()
	n := len(t.items)
	t.mu.Unlock()
	return int64(n)
}

func (t *TtlTypedSyncMap[K, V]) Range(f func(key K, value V) bool) {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	for k, entry := range t.items {
		if now.After(entry.expiresAt) {
			delete(t.items, k)
			continue
		}
		entry.expiresAt = now.Add(t.expDuration)
		t.items[k] = entry

		if !f(k, entry.value) {
			break
		}
	}
}

func (t *TtlTypedSyncMap[K, V]) sanitize() {
	ticker := time.NewTicker(t.sanitizeInterval)
	defer ticker.Stop()
	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			t.mu.Lock()
			for k, entry := range t.items {
				if now.After(entry.expiresAt) {
					delete(t.items, k)
				}
			}
			t.mu.Unlock()
		}
	}
}
