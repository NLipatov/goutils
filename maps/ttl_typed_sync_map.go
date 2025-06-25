package maps

import (
	"context"
	"sync"
	"time"
)

type TtlTypedSyncMap[K comparable, V any] struct {
	ctx              context.Context
	sanitizeInterval time.Duration
	expDuration      time.Duration
	expMu            *sync.Mutex
	exp              map[K]time.Time
	m                *TypedSyncMap[K, V]
}

func NewTtlTypedSyncMap[K comparable, V any](
	ctx context.Context,
	expDuration time.Duration,
) *TtlTypedSyncMap[K, V] {
	res := &TtlTypedSyncMap[K, V]{
		ctx:         ctx,
		expDuration: expDuration,
		exp:         make(map[K]time.Time),
		expMu:       &sync.Mutex{},
		m:           NewTypedSyncMap[K, V](),
	}
	go res.sanitize()
	return res
}

func (t *TtlTypedSyncMap[K, V]) Store(key K, value V) {
	t.expMu.Lock()
	defer t.expMu.Unlock()

	t.m.Store(key, value)
	t.exp[key] = time.Now().Add(t.expDuration)
}

func (t *TtlTypedSyncMap[K, V]) Load(key K) (V, bool) {
	t.expMu.Lock()
	defer t.expMu.Unlock()

	exp, ok := t.exp[key]
	if !ok || time.Now().After(exp) {
		// not found or expired
		t.m.Delete(key)
		delete(t.exp, key)
		var zero V
		return zero, false
	}
	// prolongate expiration time on every hit
	t.exp[key] = time.Now().Add(t.expDuration)
	return t.m.Load(key)
}

func (t *TtlTypedSyncMap[K, V]) Delete(key K) {
	t.expMu.Lock()
	defer t.expMu.Unlock()

	t.m.Delete(key)
	delete(t.exp, key)
}

func (t *TtlTypedSyncMap[K, V]) Len() int64 {
	return t.m.Len()
}

func (t *TtlTypedSyncMap[K, V]) Range(f func(key K, value V) bool) {
	now := time.Now()
	t.expMu.Lock()
	defer t.expMu.Unlock()

	for k, expireAt := range t.exp {
		if now.After(expireAt) {
			t.m.Delete(k)
			delete(t.exp, k)
			continue
		}
		v, ok := t.m.Load(k)
		if !ok {
			delete(t.exp, k)
			continue
		}
		t.exp[k] = now.Add(t.expDuration)
		if !f(k, v) {
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
			t.expMu.Lock()
			for k, v := range t.exp {
				if now.After(v) {
					t.m.Delete(k)
					delete(t.exp, k)
				}
			}
			t.expMu.Unlock()
		}
	}
}
