package hashset

import (
	"iter"
)

// Extend inserts all elements from the iterator into the set.
//
// Parameters:
//   - it: An iterator yielding elements to insert.
//
// Behavior:
//   - Elements are only added if they don't already exist in the set.
func (hs *HashSet[E, H]) Extend(it iter.Seq[E]) {
	for e := range it {
		hs.Insert(e)
	}
}

// Iter returns an iterator over all elements in the set.
//
// Returns:
//   - An iter.Seq[E] that yields all elements.
//   - Order is non-deterministic (hash order).
//
// Time Complexity: O(n)
func (hs *HashSet[E, H]) Iter() iter.Seq[E] {
	return hs.table.Keys()
}
