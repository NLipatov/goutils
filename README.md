[![Go Reference](https://pkg.go.dev/badge/github.com/NLipatov/goutils.svg)](https://pkg.go.dev/github.com/NLipatov/goutils)

---
# maps

Concurrent, generic, thread-safe maps for Go, with optional TTL (time-to-live) support.

## Types

### `TypedSyncMap[K comparable, V any]`

Thread-safe, generic map with atomic length and basic CRUD operations.

#### Features

- Safe for concurrent use.
- Keeps an atomic counter of elements (O(1) Len).
- Replaces standard Go `sync.Map`, but with type safety and length support.

#### Example

```go
import "github.com/NLipatov/goutils/maps"

m := maps.NewTypedSyncMap[int, string]()
m.Store(1, "foo")
val, ok := m.Load(1)        // val == "foo", ok == true
m.Delete(1)
length := m.Len()            // length == 0
````

---

### `TtlTypedSyncMap[K comparable, V any]`

Thread-safe, generic map with TTL (time-to-live) for each entry.

#### Features

* Each entry expires automatically after a specified duration.
* TTL is extended on each access (`Load`).
* Background janitor removes expired entries periodically.
* Safe for concurrent use.

#### Example

```go
import (
    "context"
    "time"
    "github.com/NLipatov/goutils/maps"
)

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

ttlMap := maps.NewTtlTypedSyncMap[int, string](ctx, 2*time.Second)
ttlMap.Store(1, "bar")

val, ok := ttlMap.Load(1) // "bar", true

time.Sleep(3 * time.Second)
val, ok = ttlMap.Load(1)  // "", false (expired)
```

---

## Usage

* Use `TypedSyncMap` for a fast, thread-safe, type-safe map with O(1) length.
* Use `TtlTypedSyncMap` when you need automatic expiration of keys (e.g., session storage, cache).

## License

MIT
