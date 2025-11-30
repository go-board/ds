package skipset

import (
	"iter"
)

// Extend adds another iterable collection of elements to the current set.
// Parameters:
//   - it: Iterator providing elements
//
// Behavior:
//   - For each element, it is only added if it is not already in the current set
func (ss *SkipSet[E]) Extend(it iter.Seq[E]) {
	for key := range it {
		ss.Insert(key)
	}
}

// Iter returns an iterator for all elements in the set, sorted in ascending order of elements.
// Returns:
//   - Element iterator of type [iter.Seq]
func (ss *SkipSet[E]) Iter() iter.Seq[E] {
	return func(yield func(E) bool) {
		for key := range ss.mapImpl.Keys() {
			if !yield(key) {
				return
			}
		}
	}
}

// Range returns an iterator for elements in the set that fall within the [lowerBound, upperBound) range.
// Parameters:
//   - lowerBound: Lower bound of the range (inclusive), nil for no lower bound
//   - upperBound: Upper bound of the range (exclusive), nil for no upper bound
//
// Returns:
//   - Iterator for elements in the specified range, sorted in ascending order
func (ss *SkipSet[E]) Range(lowerBound, upperBound *E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for key := range ss.mapImpl.Range(lowerBound, upperBound) {
			if !yield(key) {
				return
			}
		}
	}
}
