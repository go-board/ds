package arraystack

import (
	"github.com/go-board/ds/arraydeque"
)

// ArrayStack implements a stack data structure using ArrayDeque as the underlying storage
// T is the type of elements stored in the stack
//
// Note: ArrayStack is not thread-safe. Additional synchronization mechanisms are required for concurrent access.
type ArrayStack[T any] struct {
	// deque is the underlying storage using ArrayDeque
	deque *arraydeque.ArrayDeque[T]
}

// New creates a new empty stack with the default initial capacity
// T is the type of elements to be stored in the stack
//
// Returns:
//
//	A newly created ArrayStack instance
//
// Time Complexity: O(1)
func New[T any]() *ArrayStack[T] {
	return &ArrayStack[T]{
		deque: arraydeque.New[T](),
	}
}

// Push adds an element to the top of the stack
//
// Parameters:
//
//	value: The element value to push onto the stack
//
// Time Complexity: O(1) amortized
func (s *ArrayStack[T]) Push(value T) {
	s.deque.PushBack(value)
}

// Pop removes and returns the element at the top of the stack
//
// Returns:
//
//	The removed element and true if the stack is not empty
//	The zero value of type T and false if the stack is empty
//
// Time Complexity: O(1)
func (s *ArrayStack[T]) Pop() (T, bool) {
	return s.deque.PopBack()
}

// Peek returns the element at the top of the stack without removing it
//
// Returns:
//
//	The top element and true if the stack is not empty
//	The zero value of type T and false if the stack is empty
//
// Time Complexity: O(1)
func (s *ArrayStack[T]) Peek() (T, bool) {
	return s.deque.Back()
}

// Clear removes all elements from the stack, leaving it empty
//
// Time Complexity: O(n)
func (s *ArrayStack[T]) Clear() {
	s.deque.Clear()
}

// Len returns the number of elements in the stack
//
// Returns:
//
//	The number of elements in the stack
//
// Time Complexity: O(1)
func (s *ArrayStack[T]) Len() int {
	return s.deque.Len()
}

// IsEmpty checks if the stack contains no elements
//
// Returns:
//
//	true if the stack is empty, false otherwise
//
// Time Complexity: O(1)
func (s *ArrayStack[T]) IsEmpty() bool {
	return s.deque.IsEmpty()
}

// Clone creates a new stack that is a shallow copy of the original stack
//
// Returns:
//
//	A new ArrayStack instance with the same elements as the original
//
// Time Complexity: O(n)
func (s *ArrayStack[T]) Clone() *ArrayStack[T] {
	clone := &ArrayStack[T]{
		deque: s.deque.Clone(),
	}
	return clone
}
