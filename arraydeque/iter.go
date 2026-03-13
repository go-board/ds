package arraydeque

import (
	"iter"
)

// Iter returns an iterator over all elements in the deque in normal order.
//
// Returns:
//   - An iter.Seq[T] that yields all elements.
//
// Time Complexity: O(n)
func (d *ArrayDeque[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < len(d.data); i++ {
			if !yield(d.data[i]) {
				return
			}
		}
	}
}

// IterMut returns a mutable iterator over all elements in the deque in normal order.
//
// Returns:
//   - An iter.Seq[*T] that yields pointers to all elements.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (d *ArrayDeque[T]) IterMut() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for i := 0; i < len(d.data); i++ {
			if !yield(&d.data[i]) {
				return
			}
		}
	}
}

// IterBack returns an iterator over all elements in the deque in reverse order.
//
// Returns:
//   - An iter.Seq[T] that yields all elements in reverse order.
//
// Time Complexity: O(n)
func (d *ArrayDeque[T]) IterBack() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(d.data) - 1; i >= 0; i-- {
			if !yield(d.data[i]) {
				return
			}
		}
	}
}

// IterBackMut returns a mutable iterator over all elements in the deque in reverse order.
//
// Returns:
//   - An iter.Seq[*T] that yields pointers to all elements in reverse order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (d *ArrayDeque[T]) IterBackMut() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for i := len(d.data) - 1; i >= 0; i-- {
			if !yield(&d.data[i]) {
				return
			}
		}
	}
}

// Extend appends all elements from the iterator to the back of the deque.
//
// Parameters:
//   - it: An iterator yielding elements to append.
//
// Time Complexity: O(n)
func (d *ArrayDeque[T]) Extend(it iter.Seq[T]) {
	for v := range it {
		d.PushBack(v)
	}
}
