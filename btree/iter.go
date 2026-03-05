package btree

import (
	"iter"

	"github.com/go-board/ds/bound"
)

// RangeAsc returns an ascending iterator that traverses elements in [lowerBound, upperBound).
// Parameters:
//   - lowerBound: The lower bound of the range (inclusive); if nil, no lower bound.
//   - upperBound: The upper bound of the range (exclusive); if nil, no upper bound.
//
// Returns:
//   - An iter.Seq[E] iterator that yields all elements in ascending order within the range.
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

// IterAsc returns an ascending iterator that traverses all elements in the BTree.
// Returns:
//   - An iter.Seq[E] iterator that yields all elements in ascending order.
func (bt *BTree[E]) IterAsc() iter.Seq[E] {
	return bt.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[E](), bound.NewUnbounded[E]()))
}

// RangeDesc returns a descending iterator that traverses elements in [lowerBound, upperBound).
// Parameters:
//   - lowerBound: inclusive lower bound; nil means unbounded.
//   - upperBound: exclusive upper bound; nil means unbounded.
//
// Returns:
//   - An iter.Seq[E] iterator that yields elements in descending order within the range.
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

// IterDesc returns a descending iterator that walks all elements from back to front.
// Returns:
//   - An iter.Seq[E] iterator that yields all elements in descending order.
func (bt *BTree[E]) IterDesc() iter.Seq[E] {
	return func(yield func(E) bool) {
		if bt.root == nil {
			return
		}
		bt.rangeNodeBack(bt.root, nil, nil, yield)
	}
}
