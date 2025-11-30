package hashset

import (
	"iter"
)

// Extend adds elements from another iterable collection to the current set.
//
// Parameters:
//
//	iter: Iterator providing elements to add
//
// For each element in the iterator, it is added only if it does not already exist in the current set.
func (hs *HashSet[E, H]) Extend(it iter.Seq[E]) {
	for e := range it {
		hs.Insert(e)
	}
}

// Iter returns an iterator over all elements in the set.
//
// Return value:
//
//	Iterator generating all elements in the set, with non-deterministic order
func (hs *HashSet[E, H]) Iter() iter.Seq[E] {
	return hs.table.Keys()
}
