package btree

import (
	"iter"
)

// Range returns a range iterator that traverses all elements in the BTree with values in [lowerBound, upperBound).
// Parameters:
//   - lowerBound: The lower bound of the range (inclusive); if nil, no lower bound
//   - upperBound: The upper bound of the range (exclusive); if nil, no upper bound
//
// Returns:
//   - An iter.Seq[E] iterator that produces all elements in the specified range in order
func (bt *BTree[E]) Range(lowerBound, upperBound *E) iter.Seq[E] {
	return func(yield func(E) bool) {
		if bt.root == nil {
			return
		}
		bt.rangeNode(bt.root, lowerBound, upperBound, yield)
	}
}

// Iter returns a sequential iterator that traverses all elements in the BTree.
// Returns:
//   - An iter.Seq[E] iterator that produces all elements in the BTree in order
func (bt *BTree[E]) Iter() iter.Seq[E] {
	return bt.Range(nil, nil)
}

// RangeBack returns a reverse iterator that traverses all elements in the BTree with values in [lowerBound, upperBound)
// 参数:
//   - lowerBound: 范围的下界（包含），如果为nil则没有下界
//   - upperBound: 范围的上界（不包含），如果为nil则没有上界
//
// 返回:
//   - 一个iter.Seq[E]类型的迭代器，按逆序产生指定范围内的所有元素
func (bt *BTree[E]) RangeBack(lowerBound, upperBound *E) iter.Seq[E] {
	return func(yield func(E) bool) {
		if bt.root == nil {
			return
		}
		bt.rangeNodeBack(bt.root, lowerBound, upperBound, yield)
	}
}

// IterBack 返回一个逆序迭代器，用于从后向前遍历BTree中的所有元素
// 返回:
//   - 一个iter.Seq[E]类型的迭代器，按逆序产生BTree中的所有元素
func (bt *BTree[E]) IterBack() iter.Seq[E] {
	return func(yield func(E) bool) {
		if bt.root == nil {
			return
		}
		bt.rangeNodeBack(bt.root, nil, nil, yield)
	}
}
