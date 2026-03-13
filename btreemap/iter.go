package btreemap

import (
	"iter"

	"github.com/go-board/ds/bound"
	diter "github.com/go-board/ds/internal/iter"
)

// IterAsc returns an iterator over all key/value pairs in ascending key order.
//
// Returns:
//   - An iter.Seq2[K, V] that yields (key, value) pairs in ascending key order.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) IterAsc() iter.Seq2[K, V] {
	return diter.Split(m.btree.IterAsc(), (*node[K, V]).kv)
}

// IterMutAsc returns a mutable iterator over all key/value pairs in ascending key order.
//
// Returns:
//   - An iter.Seq2[K, *V] that yields (key, pointer to value) pairs in ascending key order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) IterMutAsc() iter.Seq2[K, *V] {
	return diter.Split(m.btree.IterAsc(), (*node[K, V]).kvMut)
}

// IterDesc returns an iterator over all key/value pairs in descending key order.
//
// Returns:
//   - An iter.Seq2[K, V] that yields (key, value) pairs in descending key order.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) IterDesc() iter.Seq2[K, V] {
	return diter.Split(m.btree.IterDesc(), (*node[K, V]).kv)
}

// IterMutDesc returns a mutable iterator over all key/value pairs in descending key order.
//
// Returns:
//   - An iter.Seq2[K, *V] that yields (key, pointer to value) pairs in descending key order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) IterMutDesc() iter.Seq2[K, *V] {
	return diter.Split(m.btree.IterDesc(), (*node[K, V]).kvMut)
}

// RangeAsc returns an iterator over key/value pairs within the given bounds in ascending key order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper key limits.
//
// Returns:
//   - An iter.Seq2[K, V] that yields (key, value) pairs within the bounds in ascending order.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (m *BTreeMap[K, V]) RangeAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, V] {
	return diter.Split(m.rangeNode(bounds, false), (*node[K, V]).kv)
}

// RangeMutAsc returns a mutable iterator over key/value pairs within the given bounds in ascending key order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper key limits.
//
// Returns:
//   - An iter.Seq2[K, *V] that yields (key, pointer to value) pairs within the bounds in ascending order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (m *BTreeMap[K, V]) RangeMutAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, *V] {
	return diter.Split(m.rangeNode(bounds, false), (*node[K, V]).kvMut)
}

// RangeDesc returns an iterator over key/value pairs within the given bounds in descending key order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper key limits.
//
// Returns:
//   - An iter.Seq2[K, V] that yields (key, value) pairs within the bounds in descending order.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (m *BTreeMap[K, V]) RangeDesc(bounds bound.RangeBounds[K]) iter.Seq2[K, V] {
	return diter.Split(m.rangeNode(bounds, true), (*node[K, V]).kv)
}

// RangeMutDesc returns a mutable iterator over key/value pairs within the given bounds in descending key order.
//
// Parameters:
//   - bounds: The range bounds specifying the lower and upper key limits.
//
// Returns:
//   - An iter.Seq2[K, *V] that yields (key, pointer to value) pairs within the bounds in descending order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(log n + k) where k is the number of elements in the range.
func (m *BTreeMap[K, V]) RangeMutDesc(bounds bound.RangeBounds[K]) iter.Seq2[K, *V] {
	return diter.Split(m.rangeNode(bounds, true), (*node[K, V]).kvMut)
}

func (m *BTreeMap[K, V]) rangeNode(bounds bound.RangeBounds[K], desc bool) iter.Seq[*node[K, V]] {
	nb := mapBounds[K, V](bounds)
	if desc {
		return m.btree.RangeDesc(nb)
	}
	return m.btree.RangeAsc(nb)
}

func mapBounds[K, V any](b bound.RangeBounds[K]) bound.RangeBounds[*node[K, V]] {
	mapOne := func(src bound.Bound[K]) bound.Bound[*node[K, V]] {
		switch src.Kind() {
		case bound.Unbounded:
			return bound.NewUnbounded[*node[K, V]]()
		case bound.Included:
			v, _ := src.Value()
			return bound.NewIncluded(&node[K, V]{Key: v})
		default:
			v, _ := src.Value()
			return bound.NewExcluded(&node[K, V]{Key: v})
		}
	}
	return bound.NewRangeBounds(mapOne(b.Start), mapOne(b.End))
}

// KeysAsc returns an iterator over all keys in ascending order.
//
// Returns:
//   - An iter.Seq[K] that yields keys in ascending order.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) KeysAsc() iter.Seq[K] { return diter.Keys(m.IterAsc()) }

// ValuesAsc returns an iterator over all values in ascending key order.
//
// Returns:
//   - An iter.Seq[V] that yields values in ascending key order.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) ValuesAsc() iter.Seq[V] { return diter.Values(m.IterAsc()) }

// ValuesMutAsc returns a mutable iterator over all values in ascending key order.
//
// Returns:
//   - An iter.Seq[*V] that yields pointers to values in ascending key order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) ValuesMutAsc() iter.Seq[*V] { return diter.Values(m.IterMutAsc()) }

// KeysDesc returns an iterator over all keys in descending order.
//
// Returns:
//   - An iter.Seq[K] that yields keys in descending order.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) KeysDesc() iter.Seq[K] { return diter.Keys(m.IterDesc()) }

// ValuesDesc returns an iterator over all values in descending key order.
//
// Returns:
//   - An iter.Seq[V] that yields values in descending key order.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) ValuesDesc() iter.Seq[V] { return diter.Values(m.IterDesc()) }

// ValuesMutDesc returns a mutable iterator over all values in descending key order.
//
// Returns:
//   - An iter.Seq[*V] that yields pointers to values in descending key order.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) ValuesMutDesc() iter.Seq[*V] { return diter.Values(m.IterMutDesc()) }

// Extend inserts all key/value pairs from the iterator into the map.
//
// Parameters:
//   - it: An iterator yielding key/value pairs to insert.
func (m *BTreeMap[K, V]) Extend(it iter.Seq2[K, V]) {
	for k, v := range it {
		m.Insert(k, v)
	}
}
