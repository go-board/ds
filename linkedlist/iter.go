package linkedlist

import (
	"iter"
)

// Iter returns an iterator over all elements in forward order (head to tail).
//
// Returns:
//   - An iter.Seq[T] that yields elements in forward order.
//
// Time Complexity: O(n)
func (ll *LinkedList[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if ll == nil || ll.head == nil {
			return
		}
		current := ll.head
		for current != nil {
			if !yield(current.value) {
				return
			}
			current = current.next
		}
	}
}

// IterMut returns a mutable iterator over all elements in forward order.
//
// Returns:
//   - An iter.Seq[*T] that yields pointers to elements in forward order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (ll *LinkedList[T]) IterMut() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		if ll == nil || ll.head == nil {
			return
		}
		current := ll.head
		for current != nil {
			if !yield(&current.value) {
				return
			}
			current = current.next
		}
	}
}

// IterBack returns an iterator over all elements in reverse order (tail to head).
//
// Returns:
//   - An iter.Seq[T] that yields elements in reverse order.
//
// Time Complexity: O(n)
func (ll *LinkedList[T]) IterBack() iter.Seq[T] {
	return func(yield func(T) bool) {
		if ll == nil || ll.tail == nil {
			return
		}
		current := ll.tail
		for current != nil {
			if !yield(current.value) {
				return
			}
			current = current.prev
		}
	}
}

// IterBackMut returns a mutable iterator over all elements in reverse order.
//
// Returns:
//   - An iter.Seq[*T] that yields pointers to elements in reverse order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (ll *LinkedList[T]) IterBackMut() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		if ll == nil || ll.tail == nil {
			return
		}
		current := ll.tail
		for current != nil {
			if !yield(&current.value) {
				return
			}
			current = current.prev
		}
	}
}

// Extend appends all elements from the iterator to the end of the linked list.
//
// Parameters:
//   - it: An iterator yielding elements to append.
//
// Time Complexity: O(n)
func (l *LinkedList[T]) Extend(it iter.Seq[T]) {
	for v := range it {
		l.PushBack(v)
	}
}
