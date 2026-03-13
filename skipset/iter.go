package skipset

import (
	"iter"

	"github.com/go-board/ds/bound"
)

// Extend inserts all elements from the iterator into the set.
//
// Parameters:
//   - it: An iterator yielding elements to insert.
func (ss *SkipSet[E]) Extend(it iter.Seq[E]) {
	for key := range it {
		ss.Insert(key)
	}
}

// IterAsc returns an iterator over all elements in ascending order.
//
// Returns:
//   - An iter.Seq[E] that yields elements in ascending order.
//
// Time Complexity: O(n)
func (ss *SkipSet[E]) IterAsc() iter.Seq[E] {
	return ss.iterFromMap(ss.mapImpl.KeysAsc())
}

// IterDesc returns an iterator over all elements in descending order.
//
// Returns:
//   - An iter.Seq[E] that yields elements in descending order.
//
// Time Complexity: O(n)
func (ss *SkipSet[E]) IterDesc() iter.Seq[E] {
	return ss.iterFromMap(ss.mapImpl.KeysDesc())
}

// RangeAsc returns an iterator over elements within the given bounds in ascending order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper limits.
//
// Returns:
//   - An iter.Seq[E] that yields elements in ascending order within the bounds.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (ss *SkipSet[E]) RangeAsc(bounds bound.RangeBounds[E]) iter.Seq[E] {
	return ss.iterFromMapPairs(ss.mapImpl.RangeAsc(bounds))
}

// RangeDesc returns an iterator over elements within the given bounds in descending order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper limits.
//
// Returns:
//   - An iter.Seq[E] that yields elements in descending order within the bounds.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (ss *SkipSet[E]) RangeDesc(bounds bound.RangeBounds[E]) iter.Seq[E] {
	return ss.iterFromMapPairs(ss.mapImpl.RangeDesc(bounds))
}

func (ss *SkipSet[E]) iterFromMap(seq iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for key := range seq {
			if !yield(key) {
				return
			}
		}
	}
}

func (ss *SkipSet[E]) iterFromMapPairs(seq iter.Seq2[E, struct{}]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for key := range seq {
			if !yield(key) {
				return
			}
		}
	}
}
