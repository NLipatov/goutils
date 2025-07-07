package queues

import (
	"errors"
	"testing"
)

// ErrInvalidCapacity must come from the same package.
func TestNewRingQueue_InvalidCapacity(t *testing.T) {
	q, err := NewRingQueue[int](-1)
	if q != nil {
		t.Errorf("expected nil queue, got %v", q)
	}
	if !errors.Is(err, ErrInvalidCapacity) {
		t.Errorf("expected ErrInvalidCapacity, got %v", err)
	}
}

func TestNewRingQueue_ValidCapacity(t *testing.T) {
	q, err := NewRingQueue[int](3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q == nil {
		t.Fatal("expected non‐nil queue")
	}
	if got := q.Size(); got != 0 {
		t.Errorf("expected size 0, got %d", got)
	}
	if got := q.Capacity(); got != 3 {
		t.Errorf("expected capacity 3, got %d", got)
	}
}

func TestRingQueue_BasicOps(t *testing.T) {
	q, _ := NewRingQueue[string](3)

	// empty Dequeue
	v, ok := q.Dequeue()
	if ok || v != "" {
		t.Errorf("empty Dequeue = (%q, %v), want (\"\", false)", v, ok)
	}
	// empty Peek
	p, ok := q.Peek()
	if ok || p != "" {
		t.Errorf("empty Peek = (%q, %v), want (\"\", false)", p, ok)
	}

	// single Enqueue/Peek/Dequeue
	q.Enqueue("a")
	if got := q.Size(); got != 1 {
		t.Errorf("size after 1 Enqueue = %d, want 1", got)
	}
	p, ok = q.Peek()
	if !ok || p != "a" {
		t.Errorf("Peek = (%q, %v), want (\"a\", true)", p, ok)
	}
	d, ok := q.Dequeue()
	if !ok || d != "a" {
		t.Errorf("Dequeue = (%q, %v), want (\"a\", true)", d, ok)
	}

	// two more Enqueue & Dequeue
	q.Enqueue("b")
	q.Enqueue("c")
	d1, ok1 := q.Dequeue()
	d2, ok2 := q.Dequeue()
	if !ok1 || !ok2 || d1 != "b" || d2 != "c" {
		t.Errorf("Dequeue sequence = (%q,%q), (%v,%v), want (\"b\",\"c\"),(true,true)", d1, d2, ok1, ok2)
	}
}

func TestRingQueue_Grow_ElseBranch(t *testing.T) {
	q, _ := NewRingQueue[int](2)
	q.Enqueue(1)
	q.Enqueue(2)
	// third Enqueue forces grow → else‐branch (wrap‐around logic)
	q.Enqueue(3)

	if got := q.Capacity(); got != 4 {
		t.Errorf("capacity after grow = %d, want 4", got)
	}
	if got := q.Size(); got != 3 {
		t.Errorf("size after grow = %d, want 3", got)
	}

	for want := 1; want <= 3; want++ {
		v, ok := q.Dequeue()
		if !ok || v != want {
			t.Errorf("after grow Dequeue = (%d, %v), want (%d, true)", v, ok, want)
		}
	}
}

func TestGrow_TailGreaterThanHead_Branch(t *testing.T) {
	// simulate contiguous block: head < tail
	q := &RingQueue[int]{
		data:     []int{1, 2, 3, 4, 5},
		head:     1,
		tail:     4,
		size:     3,
		capacity: 5,
	}
	q.grow() // hits if tail>head branch

	if q.capacity != 10 {
		t.Errorf("capacity = %d, want 10", q.capacity)
	}
	if q.head != 0 {
		t.Errorf("head = %d, want 0", q.head)
	}
	if q.tail != 3 {
		t.Errorf("tail = %d, want 3", q.tail)
	}

	want := []int{2, 3, 4}
	for i := 0; i < 3; i++ {
		if q.data[i] != want[i] {
			t.Errorf("data[%d] = %d, want %d", i, q.data[i], want[i])
		}
	}
}

func TestGrow_HeadGreaterThanTail_Branch(t *testing.T) {
	// simulate wrap‐around: head > tail
	q := &RingQueue[int]{
		data:     []int{1, 2, 3, 4, 5},
		head:     3,
		tail:     2,
		size:     3,
		capacity: 5,
	}
	q.grow() // hits else branch for wrap‐around

	if q.capacity != 10 {
		t.Errorf("capacity = %d, want 10", q.capacity)
	}
	if q.head != 0 {
		t.Errorf("head = %d, want 0", q.head)
	}
	if q.tail != 3 {
		t.Errorf("tail = %d, want 3", q.tail)
	}

	want := []int{4, 5, 1}
	for i := 0; i < 3; i++ {
		if q.data[i] != want[i] {
			t.Errorf("data[%d] = %d, want %d", i, q.data[i], want[i])
		}
	}
}

func TestMustDequeue_ReturnAndPanic(t *testing.T) {
	q, _ := NewRingQueue[string](3)
	q.Enqueue("x")

	// non‐empty MustDequeue
	if got := q.MustDequeue(); got != "x" {
		t.Errorf("MustDequeue = %q, want \"x\"", got)
	}

	// empty MustDequeue panics ErrEmpty
	defer func() {
		if r := recover(); r != ErrEmptyRingQueue {
			t.Errorf("panic = %v, want ErrEmpty", r)
		}
	}()
	q.MustDequeue()
}
