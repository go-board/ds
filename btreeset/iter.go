package btreeset

import (
	"iter"

	"github.com/go-board/ds/bound"
)

// IterAsc returns an iterator over all elements in ascending order.
//
// Returns:
//   - An iter.Seq[T] that yields elements in ascending order.
//
// Time Complexity: O(n)
func (s *BTreeSet[T]) IterAsc() iter.Seq[T] {
	return s.btree.IterAsc()
}

// RangeAsc returns an iterator over elements within the given bounds in ascending order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper limits.
//
// Returns:
//   - An iter.Seq[T] that yields elements in ascending order within the bounds.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (s *BTreeSet[T]) RangeAsc(bounds bound.RangeBounds[T]) iter.Seq[T] {
	return s.btree.RangeAsc(bounds)
}

// IterDesc returns an iterator over all elements in descending order.
//
// Returns:
//   - An iter.Seq[T] that yields elements in descending order.
//
// Time Complexity: O(n)
func (s *BTreeSet[T]) IterDesc() iter.Seq[T] {
	return s.btree.IterDesc()
}

// RangeDesc returns an iterator over elements within the given bounds in descending order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper limits.
//
// Returns:
//   - An iter.Seq[T] that yields elements in descending order within the bounds.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (s *BTreeSet[T]) RangeDesc(bounds bound.RangeBounds[T]) iter.Seq[T] {
	return s.btree.RangeDesc(bounds)
}

// Extend inserts all elements from the iterator into the set.
//
// Parameters:
//   - it: An iterator yielding elements to insert.
func (s *BTreeSet[T]) Extend(it iter.Seq[T]) {
	for e := range it {
		s.Insert(e)
	}
}
