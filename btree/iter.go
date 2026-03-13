package btree

import (
	"iter"

	"github.com/go-board/ds/bound"
)

// RangeAsc returns an iterator that traverses elements within the given bounds in ascending order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper limits.
//
// Returns:
//   - An iter.Seq[E] that yields elements in ascending order within the specified bounds.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (bt *BTree[E]) RangeAsc(bounds bound.RangeBounds[E]) iter.Seq[E] {
	lowerBound, upperBound := coarseBounds(bounds)
	return func(yield func(E) bool) {
		if bt.root == nil {
			return
		}
		bt.rangeNode(bt.root, lowerBound, upperBound, func(e E) bool {
			if !bounds.Contains(bt.comparator, e) {
				return true
			}
			return yield(e)
		})
	}
}

// IterAsc returns an iterator that traverses all elements in ascending order.
//
// Returns:
//   - An iter.Seq[E] that yields all elements in ascending order.
//
// Time Complexity: O(n)
func (bt *BTree[E]) IterAsc() iter.Seq[E] {
	return bt.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[E](), bound.NewUnbounded[E]()))
}

// RangeDesc returns an iterator that traverses elements within the given bounds in descending order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper limits.
//
// Returns:
//   - An iter.Seq[E] that yields elements in descending order within the specified bounds.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (bt *BTree[E]) RangeDesc(bounds bound.RangeBounds[E]) iter.Seq[E] {
	lowerBound, upperBound := coarseBounds(bounds)
	return func(yield func(E) bool) {
		if bt.root == nil {
			return
		}
		bt.rangeNodeBack(bt.root, lowerBound, upperBound, func(e E) bool {
			if !bounds.Contains(bt.comparator, e) {
				return true
			}
			return yield(e)
		})
	}
}

func coarseBounds[E any](bounds bound.RangeBounds[E]) (lower, upper *E) {
	if v, ok := bounds.Start.Value(); ok {
		lower = &v
	}
	if bounds.End.IsExcluded() {
		if v, ok := bounds.End.Value(); ok {
			upper = &v
		}
	}
	return
}

// IterDesc returns an iterator that traverses all elements in descending order.
//
// Returns:
//   - An iter.Seq[E] that yields all elements in descending order.
//
// Time Complexity: O(n)
func (bt *BTree[E]) IterDesc() iter.Seq[E] {
	return func(yield func(E) bool) {
		if bt.root == nil {
			return
		}
		bt.rangeNodeBack(bt.root, nil, nil, yield)
	}
}
