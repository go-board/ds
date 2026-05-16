package arrayset

import (
	"iter"

	"github.com/go-board/ds/bound"
)

// Extend inserts all elements from iterator.
func (s *ArraySet[E]) Extend(it iter.Seq[E]) {
	for e := range it {
		s.Insert(e)
	}
}

// IterAsc returns an iterator over all elements in ascending order.
func (s *ArraySet[E]) IterAsc() iter.Seq[E] {
	return s.mapImpl.KeysAsc()
}

// IterDesc returns an iterator over all elements in descending order.
func (s *ArraySet[E]) IterDesc() iter.Seq[E] {
	return s.mapImpl.KeysDesc()
}

// RangeAsc returns elements within bounds in ascending order.
func (s *ArraySet[E]) RangeAsc(bounds bound.RangeBounds[E]) iter.Seq[E] {
	return s.iterFromPairs(s.mapImpl.RangeAsc(bounds))
}

// RangeDesc returns elements within bounds in descending order.
func (s *ArraySet[E]) RangeDesc(bounds bound.RangeBounds[E]) iter.Seq[E] {
	return s.iterFromPairs(s.mapImpl.RangeDesc(bounds))
}

func (s *ArraySet[E]) iterFromPairs(seq iter.Seq2[E, struct{}]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}
