package skipset

import "iter"

import "github.com/go-board/ds/bound"

// Extend inserts all elements from the iterator into the set.
func (ss *SkipSet[E]) Extend(it iter.Seq[E]) {
	for key := range it {
		ss.Insert(key)
	}
}

// IterAsc returns all elements in ascending order.
func (ss *SkipSet[E]) IterAsc() iter.Seq[E] {
	return ss.iterFromMap(ss.mapImpl.KeysAsc())
}

// IterDesc returns all elements in descending order.
func (ss *SkipSet[E]) IterDesc() iter.Seq[E] {
	return ss.iterFromMap(ss.mapImpl.KeysDesc())
}

// RangeAsc returns elements in ascending order within [lowerBound, upperBound).
func (ss *SkipSet[E]) RangeAsc(bounds bound.RangeBounds[E]) iter.Seq[E] {
	return ss.iterFromMapPairs(ss.mapImpl.RangeAsc(bounds))
}

// RangeDesc returns elements in descending order within [lowerBound, upperBound).
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
