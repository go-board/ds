// Package arraystack implements a generic stack data structure using ArrayDeque as the underlying storage.
package arraystack

import (
	"iter"
)

// Iter returns a sequential [iter.Seq] over the elements in LIFO (Last-In-First-Out) order
//
// Returns:
//
//	An iterator that yields elements in LIFO order
//
// Time Complexity: O(n) for full iteration
func (s *ArrayStack[T]) Iter() iter.Seq[T] {
	return s.deque.IterBack()
}

// IterMut returns a mutable [iter.Seq] over the elements in LIFO (Last-In-First-Out) order
//
// Returns:
//
//	An iterator that yields pointers to elements in LIFO order, allowing modification
//
// Time Complexity: O(n) for full iteration
func (s *ArrayStack[T]) IterMut() iter.Seq[*T] {
	return s.deque.IterBackMut()
}

// Extend adds multiple elements to the top of the stack in the order they are provided
//
// Parameters:
//
//	values: An [iter.Seq] of elements to push onto the stack
//
// Time Complexity: O(n) where n is the number of elements to add, amortized
func (s *ArrayStack[T]) Extend(it iter.Seq[T]) {
	s.deque.Extend(it)
}
