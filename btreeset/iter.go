package btreeset

import "iter"

import "github.com/go-board/ds/bound"

// IterAsc returns all elements in ascending order.
func (s *BTreeSet[T]) IterAsc() iter.Seq[T] {
	return s.btree.IterAsc()
}

// RangeAsc returns elements in ascending order within [lowerBound, upperBound).
func (s *BTreeSet[T]) RangeAsc(bounds bound.RangeBounds[T]) iter.Seq[T] {
	return s.btree.RangeAsc(bounds)
}

// IterDesc returns all elements in descending order.
func (s *BTreeSet[T]) IterDesc() iter.Seq[T] {
	return s.btree.IterDesc()
}

// RangeDesc returns elements in descending order within [lowerBound, upperBound).
func (s *BTreeSet[T]) RangeDesc(bounds bound.RangeBounds[T]) iter.Seq[T] {
	return s.btree.RangeDesc(bounds)
}

// Extend inserts all elements from the iterator into the set.
func (s *BTreeSet[T]) Extend(it iter.Seq[T]) {
	for e := range it {
		s.Insert(e)
	}
}
