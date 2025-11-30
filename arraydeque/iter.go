// Package arraydeque implements a generic double-ended queue data structure.
package arraydeque

import (
	"iter"
)

// Iter returns an iterator over all elements in the queue, in normal order.
//
// Returns:
//   - An iterator over the elements, of type iter.Seq[T]
//
// Time complexity: O(n)
func (d *ArrayDeque[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < len(d.data); i++ {
			if !yield(d.data[i]) {
				return
			}
		}
	}
}

// IterMut returns a mutable iterator over all elements in the queue, in normal order.
//
// Returns:
//   - A mutable iterator over the elements, of type iter.Seq[*T]
//
// Time complexity: O(n)
func (d *ArrayDeque[T]) IterMut() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for i := 0; i < len(d.data); i++ {
			if !yield(&d.data[i]) {
				return
			}
		}
	}
}

// IterBack returns an iterator over all elements in the queue, in reverse order.
//
// Returns:
//   - An iterator over the elements in reverse order, of type iter.Seq[T]
//
// Time complexity: O(n)
func (d *ArrayDeque[T]) IterBack() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(d.data) - 1; i >= 0; i-- {
			if !yield(d.data[i]) {
				return
			}
		}
	}
}

// IterBackMut returns a mutable iterator over all elements in the queue, in reverse order.
//
// Returns:
//   - A mutable iterator over the elements in reverse order, of type iter.Seq[*T]
//
// Time complexity: O(n)
func (d *ArrayDeque[T]) IterBackMut() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for i := len(d.data) - 1; i >= 0; i-- {
			if !yield(&d.data[i]) {
				return
			}
		}
	}
}

// Extend appends all elements from the given iterator to the back of the deque
//
// Parameters:
//   - i: An iterator over elements of type T
//
// Time complexity: O(n)
func (d *ArrayDeque[T]) Extend(i iter.Seq[T]) {
	for v := range i {
		d.PushBack(v)
	}
}
