package linkedlist

import (
	"iter"
)

// Iter returns a sequential iterator that yields the elements of the LinkedList in forward order.
//
// It returns an `iter.Seq[T]` that produces each element in the list from head to tail.
func (ll *LinkedList[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		// check if linked list is nil
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

// IterMut returns a mutable sequential iterator that yields pointers to elements
// of the LinkedList in forward order, allowing in-place modification.
//
// It returns an `iter.Seq[*T]` that produces a pointer to each element.
func (ll *LinkedList[T]) IterMut() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		// check if linked list is nil
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

// IterBack returns a reverse iterator that yields the elements of the LinkedList
// from tail to head.
//
// It returns an `iter.Seq[T]` that produces elements in reverse order.
func (ll *LinkedList[T]) IterBack() iter.Seq[T] {
	return func(yield func(T) bool) {
		// check if linked list is nil
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

// IterBackMut returns a mutable reverse iterator that yields pointers to elements
// of the LinkedList from tail to head, allowing in-place modification.
//
// It returns an `iter.Seq[*T]` that produces a pointer to each element.
func (ll *LinkedList[T]) IterBackMut() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		// check if linked list is nil
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

// Extend appends all elements from the given iterator to the end of the linked list
//
// Parameters:
//   - it: An iterator over elements of type T
//
// Time complexity: O(n)
func (l *LinkedList[T]) Extend(it iter.Seq[T]) {
	for v := range it {
		l.PushBack(v)
	}
}
