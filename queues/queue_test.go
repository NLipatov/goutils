package queues

import (
	"errors"
	"testing"
)

// Test that NewQueue panics on invalid capacity (≤ 0).
func TestNewQueue_PanicOnInvalidCapacity(t *testing.T) {
	_, err := NewQueue[int](-1)
	if !errors.Is(err, ErrInvalidCapacity) {
		t.Fatalf("expected 'ErrInvalidCapacity' error, got '%v'", err)
	}
}

// Dequeue on an empty queue should return ok=false.
func TestDequeue_EmptyQueue(t *testing.T) {
	q, qErr := NewQueue[int](1)
	if qErr != nil {
		t.Fatal(qErr)
	}
	v, ok := q.Dequeue()
	if ok {
		t.Errorf("expected ok=false, got v=%v, ok=%v", v, ok)
	}
}

// Enqueue then Dequeue a single element.
func TestEnqueueDequeue_Single(t *testing.T) {
	q, qErr := NewQueue[int](1)
	if qErr != nil {
		t.Fatal(qErr)
	}
	q.Enqueue(42)
	if got := q.Size(); got != 1 {
		t.Errorf("expected size 1 after Enqueue, got %d", got)
	}
	v, ok := q.Dequeue()
	if !ok || v != 42 {
		t.Errorf("expected Dequeue to return 42,true; got %v,%v", v, ok)
	}
	if got := q.Size(); got != 0 {
		t.Errorf("expected size 0 after Dequeue, got %d", got)
	}
}

// Enqueue multiple values and ensure they come out in FIFO order.
func TestEnqueueDequeue_Multiple(t *testing.T) {
	q, qErr := NewQueue[int](1)
	if qErr != nil {
		t.Fatal(qErr)
	}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if got := q.Size(); got != 3 {
		t.Errorf("expected size 3 after three Enqueue calls, got %d", got)
	}
	v1, o1 := q.Dequeue()
	v2, o2 := q.Dequeue()
	v3, o3 := q.Dequeue()
	if !(o1 && o2 && o3) {
		t.Error("expected ok=true for all three Dequeue calls")
	}
	if v1 != 1 || v2 != 2 || v3 != 3 {
		t.Errorf("expected sequence 1,2,3; got %v,%v,%v", v1, v2, v3)
	}
}

// Peek on an empty queue returns ok=false.
func TestPeek_Empty(t *testing.T) {
	q, qErr := NewQueue[int](1)
	if qErr != nil {
		t.Fatal(qErr)
	}
	v, ok := q.Peek()
	if ok {
		t.Errorf("expected ok=false when peeking empty queue, got %v,%v", v, ok)
	}
}

// Peek on a non-empty queue returns the first element without removing it.
func TestPeek_NonEmpty(t *testing.T) {
	q, qErr := NewQueue[int](1)
	if qErr != nil {
		t.Fatal(qErr)
	}
	q.Enqueue(99)
	v, ok := q.Peek()
	if !ok || v != 99 {
		t.Errorf("expected Peek to return 99,true; got %v,%v", v, ok)
	}
	if got := q.Size(); got != 1 {
		t.Errorf("expected size 1 after Peek, got %d", got)
	}
}

// Size should reflect the number of elements currently in the queue.
func TestSize(t *testing.T) {
	q, qErr := NewQueue[int](1)
	if qErr != nil {
		t.Fatal(qErr)
	}
	if got := q.Size(); got != 0 {
		t.Errorf("expected size 0 for new queue, got %d", got)
	}
	q.Enqueue(7)
	if got := q.Size(); got != 1 {
		t.Errorf("expected size 1 after Enqueue, got %d", got)
	}
	q.Dequeue()
	if got := q.Size(); got != 0 {
		t.Errorf("expected size 0 after Dequeue, got %d", got)
	}
}

func TestNewQueue_ValidCapacity(t *testing.T) {
	q, err := NewQueue[int](1)
	if err != nil {
		t.Fatalf("did not expect error for capacity 1, got %v", err)
	}
	if q == nil {
		t.Fatal("expected non-nil queue pointer")
	}
	if got := q.Size(); got != 0 {
		t.Errorf("expected size 0 for new queue, got %d", got)
	}
}

func TestQueue_EnqueueDequeue_Empty(t *testing.T) {
	q, err := NewQueue[int](1)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	v, ok := q.Dequeue()
	if ok {
		t.Errorf("expected ok=false when dequeuing empty queue, got v=%v, ok=%v", v, ok)
	}
}

func TestQueue_EnqueueDequeue_Single(t *testing.T) {
	q, err := NewQueue[int](1)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	q.Enqueue(42)
	if got := q.Size(); got != 1 {
		t.Errorf("expected size 1 after Enqueue, got %d", got)
	}
	v, ok := q.Dequeue()
	if !ok || v != 42 {
		t.Errorf("expected Dequeue to return 42,true; got %v,%v", v, ok)
	}
	if got := q.Size(); got != 0 {
		t.Errorf("expected size 0 after Dequeue, got %d", got)
	}
}
func TestQueue_EnqueueDequeue_Multiple(t *testing.T) {
	q, err := NewQueue[int](1)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if got := q.Size(); got != 3 {
		t.Errorf("expected size 3 after three Enqueue calls, got %d", got)
	}
	v1, ok1 := q.Dequeue()
	v2, ok2 := q.Dequeue()
	v3, ok3 := q.Dequeue()
	if !(ok1 && ok2 && ok3) {
		t.Error("expected ok=true for all three Dequeue calls")
	}
	if v1 != 1 || v2 != 2 || v3 != 3 {
		t.Errorf("expected sequence 1,2,3; got %v,%v,%v", v1, v2, v3)
	}
}
func TestQueue_Peek_Empty(t *testing.T) {
	q, err := NewQueue[int](1)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	v, ok := q.Peek()
	if ok {
		t.Errorf("expected ok=false when peeking empty queue, got v=%v, ok=%v", v, ok)
	}
}
func TestQueue_Peek_NonEmpty(t *testing.T) {
	q, err := NewQueue[int](1)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	q.Enqueue(99)
	v, ok := q.Peek()
	if !ok || v != 99 {
		t.Errorf("expected Peek to return 99,true; got %v,%v", v, ok)
	}
	// ensure Peek does not remove the element
	if got := q.Size(); got != 1 {
		t.Errorf("expected size 1 after Peek, got %d", got)
	}
}

func TestQueue_Size(t *testing.T) {
	q, err := NewQueue[int](1)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if got := q.Size(); got != 0 {
		t.Errorf("expected size 0 for new queue, got %d", got)
	}
	q.Enqueue(7)
	if got := q.Size(); got != 1 {
		t.Errorf("expected size 1 after Enqueue, got %d", got)
	}
	q.Dequeue()
	if got := q.Size(); got != 0 {
		t.Errorf("expected size 0 after Dequeue, got %d", got)
	}
}
