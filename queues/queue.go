package queues

// Queue is a generic FIFO queue.
// Not safe for concurrent use.
type Queue[T any] struct {
	arr []T
}

// NewQueue returns a new empty Queue with the specified initial capacity.
func NewQueue[T any](capacity int) *Queue[T] {
	if capacity <= 0 {
		panic("capacity must be > 0")
	}
	return &Queue[T]{
		arr: make([]T, 0, capacity),
	}
}

// Enqueue adds value to the end of the queue.
func (q *Queue[T]) Enqueue(value T) {
	q.arr = append(q.arr, value)
}

// Dequeue removes and returns the first element from the queue.
// The boolean result is false if the queue is empty.
func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if len(q.arr) == 0 {
		return zero, false
	}

	head := q.arr[0]
	q.arr[0] = zero
	q.arr = q.arr[1:]
	return head, true
}

// Peek returns the first element of the queue without removing it.
// The boolean result is false if the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.arr) == 0 {
		var zero T
		return zero, false
	}

	return q.arr[0], true
}

// Size returns the number of elements in the queue.
func (q *Queue[T]) Size() int {
	return len(q.arr)
}
