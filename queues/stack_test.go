package queues

import (
	"errors"
	"testing"
)

func TestNewStack_ErrorOnInvalidCapacity(t *testing.T) {
	_, err := NewStack[int](-1)
	if !errors.Is(err, ErrInvalidCapacity) {
		t.Fatalf("expected 'ErrInvalidCapacity' error, got '%v'", err)
	}
}

func TestStack_PushPopPeek_Size(t *testing.T) {
	stack, stackErr := NewStack[int](1)
	if stackErr != nil {
		t.Fatal(stackErr)
	}

	if _, ok := stack.Pop(); ok {
		t.Fatal("Pop on empty stack should return false")
	}
	if _, ok := stack.Peek(); ok {
		t.Fatal("Peek on empty stack should return false")
	}
	if stack.Size() != 0 {
		t.Fatalf("Expected size 0, got %d", stack.Size())
	}

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if stack.Size() != 3 {
		t.Fatalf("Expected size 3, got %d", stack.Size())
	}

	if v, ok := stack.Peek(); !ok || v != 3 {
		t.Fatalf("Peek should return 3, got %v", v)
	}
	if stack.Size() != 3 {
		t.Fatalf("Peek should not change size, got %d", stack.Size())
	}

	if v, ok := stack.Pop(); !ok || v != 3 {
		t.Fatalf("Pop should return 3, got %v", v)
	}
	if v, ok := stack.Pop(); !ok || v != 2 {
		t.Fatalf("Pop should return 2, got %v", v)
	}
	if v, ok := stack.Pop(); !ok || v != 1 {
		t.Fatalf("Pop should return 1, got %v", v)
	}

	if stack.Size() != 0 {
		t.Fatalf("Expected size 0 after pops, got %d", stack.Size())
	}
	if _, ok := stack.Pop(); ok {
		t.Fatal("Pop on empty stack should return false")
	}
}

func TestStack_WithStrings(t *testing.T) {
	stack, stackErr := NewStack[string](1)
	if stackErr != nil {
		t.Fatal(stackErr)
	}
	stack.Push("a")
	stack.Push("b")
	if v, ok := stack.Pop(); !ok || v != "b" {
		t.Fatalf("Pop should return b, got %v", v)
	}
	if v, ok := stack.Pop(); !ok || v != "a" {
		t.Fatalf("Pop should return a, got %v", v)
	}
}
