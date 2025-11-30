package btreeset

import "iter"

// Iter returns an iterator that traverses all elements in the set in ascending order.
//
// Returns:
//
//	An iterator that generates all elements in the set, ordered in ascending order
//
// Note: Modifying the set (inserting or deleting) during iteration may cause undefined behavior.
//
// Time complexity: O(n) for a full traversal, with an amortized O(1) time complexity for individual Next operations
func (s *BTreeSet[T]) Iter() iter.Seq[T] {
	return s.btree.Iter()
}

// Range returns an iterator over elements in the set that fall within the range [lowerBound, upperBound).
//
// Parameters:
//
//	lowerBound: The lower bound of the range, or nil for no lower bound
//	upperBound: The upper bound of the range, or nil for no upper bound
//
// Returns:
//
//	An iterator that generates all elements within the specified range, ordered in ascending order
func (s *BTreeSet[T]) Range(lowerBound, upperBound *T) iter.Seq[T] {
	return s.btree.Range(lowerBound, upperBound)
}

// IterBack returns an iterator that traverses all elements in the set in descending order.
//
// Returns:
//
//	An iterator that generates all elements in the set, ordered in descending order
//
// Note: Modifying the set (inserting or deleting) during iteration may cause undefined behavior.
//
// Time complexity: O(n) for a full traversal, with an amortized O(1) time complexity for individual Next operations
func (s *BTreeSet[T]) IterBack() iter.Seq[T] {
	return s.btree.IterBack()
}

// Extend adds all elements from another iterable collection to the current set.
//
// Parameters:
//
//	iter: Iterator providing elements to add
//
// For each element in the iterator, it is added only if it does not already exist in the current set.
func (s *BTreeSet[T]) Extend(it iter.Seq[T]) {
	for e := range it {
		s.Insert(e)
	}
}
