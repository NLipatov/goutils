package queues

import (
	"testing"
)

func TestQueue_Empty(t *testing.T) {
	q := NewQueue[int](0)
	if size := q.Size(); size != 0 {
		t.Fatalf("expected size 0, got %d", size)
	}

	if _, ok := q.Dequeue(); ok {
		t.Fatal("Dequeue on empty queue should return false")
	}
	if _, ok := q.Peek(); ok {
		t.Fatal("Peek on empty queue should return false")
	}
}

func TestQueue_Ints(t *testing.T) {
	q := NewQueue[int](0)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	if size := q.Size(); size != 3 {
		t.Fatalf("expected size 3, got %d", size)
	}

	// FIFO check
	if val, ok := q.Peek(); !ok || val != 1 {
		t.Fatalf("Peek should return 1, got %v", val)
	}
	if val, ok := q.Dequeue(); !ok || val != 1 {
		t.Fatalf("Dequeue should return 1, got %v", val)
	}
	if val, ok := q.Dequeue(); !ok || val != 2 {
		t.Fatalf("Dequeue should return 2, got %v", val)
	}
	if val, ok := q.Dequeue(); !ok || val != 3 {
		t.Fatalf("Dequeue should return 3, got %v", val)
	}

	if size := q.Size(); size != 0 {
		t.Fatalf("expected size 0 after Dequeue all, got %d", size)
	}
	if _, ok := q.Peek(); ok {
		t.Fatal("Peek on empty queue should return false after Dequeue all")
	}
	if _, ok := q.Dequeue(); ok {
		t.Fatal("Dequeue on empty queue should return false after Dequeue all")
	}
}

func TestQueue_ReuseAfterEmpty(t *testing.T) {
	q := NewQueue[string](0)
	q.Enqueue("a")
	q.Dequeue()
	q.Enqueue("b")
	q.Enqueue("c")
	if v, ok := q.Peek(); !ok || v != "b" {
		t.Fatalf("Peek should return 'b', got '%v'", v)
	}
	if v, ok := q.Dequeue(); !ok || v != "b" {
		t.Fatalf("Dequeue should return 'b', got '%v'", v)
	}
	if v, ok := q.Dequeue(); !ok || v != "c" {
		t.Fatalf("Dequeue should return 'c', got '%v'", v)
	}
	if _, ok := q.Dequeue(); ok {
		t.Fatal("Dequeue on empty queue should return false after reuse")
	}
}

func TestQueue_SingleElement(t *testing.T) {
	q := NewQueue[int](0)
	q.Enqueue(42)
	if v, ok := q.Peek(); !ok || v != 42 {
		t.Fatalf("Peek should return 42, got %v", v)
	}
	if v, ok := q.Dequeue(); !ok || v != 42 {
		t.Fatalf("Dequeue should return 42, got %v", v)
	}
	if _, ok := q.Dequeue(); ok {
		t.Fatal("Dequeue on empty queue should return false")
	}
	if q.Size() != 0 {
		t.Fatalf("expected size 0 after Dequeue, got %d", q.Size())
	}
}
