![License](https://img.shields.io/badge/Licence-MIT-brightgreen)
[![Go Reference](https://pkg.go.dev/badge/github.com/NLipatov/goutils.svg)](https://pkg.go.dev/github.com/NLipatov/goutils)
![Build](https://github.com/NLipatov/goutils/actions/workflows/main.yml/badge.svg)
![codecov](https://codecov.io/gh/NLipatov/goutils/branch/main/graph/badge.svg)

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

# queues

Generic FIFO queues and LIFO stacks for Go.

## Types

### `Queue[T any]`

FIFO (first-in, first-out) queue implemented on top of a slice.

#### Features

* O(1) amortized enqueue/dequeue.
* Safe for single-goroutine use.
* Automatically grows underlying slice.
* Exposes `Size()` method.

#### Example

```go
import "github.com/NLipatov/goutils/queues"

q, _ := queues.NewQueue // initial capacity 10
q.Enqueue(1)
q.Enqueue(2)
val, ok := q.Dequeue() // val == 1, ok == true
peek, ok := q.Peek()   // peek == 2, ok == true
length := q.Size()     // length == 1
```

---

### `RingQueue[T any]`

FIFO queue implemented as a ring buffer.

#### Features

* True O(1) enqueue/dequeue without slice-shifting.
* No stale elements; helps GC.
* Automatically doubles capacity when full.
* Exposes `Size()` and `Capacity()`.

#### Example

```go
import "github.com/NLipatov/goutils/queues"

rq, _ := queues.NewRingQueue
rq.Enqueue("a")
rq.Enqueue("b")
first, _ := rq.Dequeue()   // "a"
currentCap := rq.Capacity() // 4 (or 8 after growth)
```

---

### `Stack[T any]`

LIFO (last-in, first-out) stack implemented on top of a slice.

#### Features

* O(1) amortized push/pop.
* Safe for single-goroutine use.
* Automatically grows underlying slice.
* Exposes `Size()`, `Push()`, `Pop()`, and `Peek()`.

#### Example

```go
import "github.com/NLipatov/goutils/queues"

s, _ := queues.NewStack // initial capacity 5
s.Push("first")
s.Push("second")
top, ok := s.Peek() // top == "second", ok == true
val, ok := s.Pop()  // val == "second", ok == true
length := s.Size()  // length == 1
```

---

## License

MIT
