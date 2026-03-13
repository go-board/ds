package arraystack

import (
	"iter"
)

// Iter returns an iterator over all elements in LIFO (last-in-first-out) order.
//
// Returns:
//   - An iter.Seq[T] that yields elements in LIFO order.
//
// Time Complexity: O(n)
func (s *ArrayStack[T]) Iter() iter.Seq[T] {
	return s.deque.IterBack()
}

// IterMut returns a mutable iterator over all elements in LIFO order.
//
// Returns:
//   - An iter.Seq[*T] that yields pointers to elements in LIFO order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (s *ArrayStack[T]) IterMut() iter.Seq[*T] {
	return s.deque.IterBackMut()
}

// Extend pushes all elements from the iterator onto the stack.
//
// Parameters:
//   - it: An iterator yielding elements to push onto the stack.
//
// Time Complexity: O(n) amortized
func (s *ArrayStack[T]) Extend(it iter.Seq[T]) {
	s.deque.Extend(it)
}
