package queues

// RingQueue is a generic FIFO queue implemented as a ring buffer.
//
// Compared to a slice-based queue, RingQueue avoids memory leaks and unnecessary data copying:
// all operations are O(1), and no memory is wasted by "stale" elements in the backing array.
// The underlying buffer is automatically resized when full.
// Not safe for concurrent use.
type RingQueue[T any] struct {
	data       []T
	head, tail int
	size       int
	capacity   int
}

// NewRingQueue returns a new empty RingQueue with the specified initial capacity (>0).
func NewRingQueue[T any](capacity int) (*RingQueue[T], error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	return &RingQueue[T]{
		data:     make([]T, capacity),
		capacity: capacity,
	}, nil
}

// Enqueue adds value to the end of the queue.
// Automatically grows the buffer if full.
func (q *RingQueue[T]) Enqueue(value T) {
	if q.size == q.capacity {
		q.grow()
	}
	q.data[q.tail] = value
	q.tail = (q.tail + 1) % q.capacity
	q.size++
}

// Dequeue removes and returns the first element from the queue.
// The boolean result is false if the queue is empty.
func (q *RingQueue[T]) Dequeue() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	val := q.data[q.head]
	q.data[q.head] = zero // help GC
	q.head = (q.head + 1) % q.capacity
	q.size--
	return val, true
}

// Peek returns the first element of the queue without removing it.
// The boolean result is false if the queue is empty.
func (q *RingQueue[T]) Peek() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	return q.data[q.head], true
}

// Size returns the number of elements in the queue.
func (q *RingQueue[T]) Size() int {
	return q.size
}

// Capacity returns the current capacity of the queue.
func (q *RingQueue[T]) Capacity() int {
	return q.capacity
}

// grow doubles the capacity of the buffer.
func (q *RingQueue[T]) grow() {
	newCap := q.capacity * 2
	newData := make([]T, newCap)
	if q.tail > q.head {
		copy(newData, q.data[q.head:q.tail])
	} else {
		n := copy(newData, q.data[q.head:])
		copy(newData[n:], q.data[:q.tail])
	}
	q.head = 0
	q.tail = q.size
	q.data = newData
	q.capacity = newCap
}

// MustDequeue removes and returns the first element or panics if the queue is empty.
// Use for internal invariants or testing.
func (q *RingQueue[T]) MustDequeue() T {
	val, ok := q.Dequeue()
	if !ok {
		panic(ErrEmptyRingQueue)
	}
	return val
}
