package queues

// Stack is a generic LIFO stack.
type Stack[T any] struct {
	arr []T
}

// NewStack creates a new empty Stack with the given capacity.
func NewStack[T any](capacity int) (*Stack[T], error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	return &Stack[T]{
		arr: make([]T, 0, capacity),
	}, nil
}

// Push adds a value to the top of the stack.
func (stack *Stack[T]) Push(value T) {
	stack.arr = append(stack.arr, value)
}

// Pop removes and returns the value from the top of the stack.
// The boolean result is false if the stack is empty.
func (stack *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(stack.arr) == 0 {
		return zero, false
	}

	headIdx := stack.headIndex()
	head := stack.arr[headIdx]
	stack.arr[headIdx] = zero
	stack.arr = stack.arr[:headIdx]
	return head, true
}

// Size returns the number of elements in the stack.
func (stack *Stack[T]) Size() int {
	return len(stack.arr)
}

// Peek returns the value at the top of the stack without removing it.
// The boolean result is false if the stack is empty.
func (stack *Stack[T]) Peek() (T, bool) {
	if len(stack.arr) == 0 {
		var zero T
		return zero, false
	}

	headIdx := stack.headIndex()
	return stack.arr[headIdx], true
}

func (stack *Stack[T]) headIndex() int {
	return len(stack.arr) - 1
}
