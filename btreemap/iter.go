package btreemap

import (
	"iter"

	diter "github.com/go-board/ds/internal/iter"
)

// Iter returns an iterator that yields all key/value pairs in the BTreeMap
// in ascending order by key.
//
// Return value:
//   - An `iter.Seq2[K, V]` iterator that produces all key/value pairs
func (m *BTreeMap[K, V]) Iter() iter.Seq2[K, V] {
	return diter.Split(m.btree.Iter(), (*node[K, V]).kv)
}

// IterMut returns a mutable iterator that yields all key/value pairs in the
// BTreeMap in ascending key order, providing pointers to values for in-place
// modification.
//
// Return value:
//   - An `iter.Seq2[K, *V]` iterator that produces each key and a pointer to its value
func (m *BTreeMap[K, V]) IterMut() iter.Seq2[K, *V] {
	return diter.Split(m.btree.Iter(), (*node[K, V]).kvMut)
}

// IterBack returns a reverse iterator that yields all key/value pairs in the
// BTreeMap in descending key order.
//
// Return value:
//   - An `iter.Seq2[K, V]` iterator that produces all key/value pairs
func (m *BTreeMap[K, V]) IterBack() iter.Seq2[K, V] {
	return diter.Split(m.btree.IterBack(), (*node[K, V]).kv)
}

// IterBackMut returns a mutable reverse iterator that yields all key/value
// pairs in the BTreeMap in descending key order, providing pointers to values
// for in-place modification.
//
// Return value:
//   - An `iter.Seq2[K, *V]` iterator that produces each key and a pointer to its value
func (m *BTreeMap[K, V]) IterBackMut() iter.Seq2[K, *V] {
	return diter.Split(m.btree.IterBack(), (*node[K, V]).kvMut)
}

// RangeMut returns a mutable iterator over key-value pairs where the key is in the [lowerBound, upperBound) range.
// Parameters:
//   - lowerBound: Lower bound of the range (inclusive), nil means no lower bound
//   - upperBound: Upper bound of the range (exclusive), nil means no upper bound
//
// Return value:
//   - Mutable iterator over key-value pairs within the specified range, sorted by key in ascending order
func (m *BTreeMap[K, V]) RangeMut(lowerBound, upperBound *K) iter.Seq2[K, *V] {
	return diter.Split(m.rangeNode(lowerBound, upperBound), (*node[K, V]).kvMut)
}

// Range returns an iterator over key-value pairs where the key is in the [lowerBound, upperBound) range.
//
// Parameters:
//   - lowerBound: Lower bound of the range (inclusive), nil means no lower bound
//   - upperBound: Upper bound of the range (exclusive), nil means no upper bound
//
// Return value:
//   - Iterator over key-value pairs within the specified range, sorted by key in ascending order
func (m *BTreeMap[K, V]) Range(lowerBound, upperBound *K) iter.Seq2[K, V] {
	return diter.Split(m.rangeNode(lowerBound, upperBound), (*node[K, V]).kv)
}

func (m *BTreeMap[K, V]) rangeNode(lowerBound, upperBound *K) iter.Seq[*node[K, V]] {
	var lowerNode, upperNode *node[K, V]
	if lowerBound != nil {
		lowerNode = &node[K, V]{Key: *lowerBound}
	}
	if upperBound != nil {
		upperNode = &node[K, V]{Key: *upperBound}
	}
	var lowerPtr, upperPtr **node[K, V]
	if lowerNode != nil {
		lowerPtr = &lowerNode
	}
	if upperNode != nil {
		upperPtr = &upperNode
	}
	return m.btree.Range(lowerPtr, upperPtr)
}

// Keys returns an iterator over all keys in the BTreeMap.
//
// Return value:
//   - Iterator over keys, of type iter.Seq[K]
func (m *BTreeMap[K, V]) Keys() iter.Seq[K] {
	return diter.Keys(m.Iter())
}

// Values returns an iterator over all values in the BTreeMap.
//
// Return value:
//   - Iterator over values, of type iter.Seq[V]
func (m *BTreeMap[K, V]) Values() iter.Seq[V] {
	return diter.Values(m.Iter())
}

// ValuesMut returns a mutable iterator over all values in the BTreeMap.
//
// Return value:
//   - Iterator over mutable values, of type iter.Seq[*V]
func (m *BTreeMap[K, V]) ValuesMut() iter.Seq[*V] {
	return diter.Values(m.IterMut())
}

// KeysBack returns an iterator over all keys in the BTreeMap, in reverse order.
//
// Return value:
//   - Iterator over keys, of type iter.Seq[K]s
func (m *BTreeMap[K, V]) KeysBack() iter.Seq[K] {
	return diter.Keys(m.IterBack())
}

// ValuesBack returns an iterator over all values in the BTreeMap, in reverse order.
//
// Return value:
//   - Iterator over values, of type iter.Seq[V]
func (m *BTreeMap[K, V]) ValuesBack() iter.Seq[V] {
	return diter.Values(m.IterBack())
}

// ValuesBackMut returns a mutable iterator over all values in the BTreeMap, in reverse order.
//
// Return value:
//   - Iterator over mutable values, of type iter.Seq[*V]
func (m *BTreeMap[K, V]) ValuesBackMut() iter.Seq[*V] {
	return diter.Values(m.IterBackMut())
}

// Extend adds another iterable collection of key-value pairs to the current map.
// Parameters:
//   - it: Iterator providing key-value pairs
//
// Behavior:
//   - For each key-value pair, if the key exists, its value is updated; otherwise, a new key-value pair is added
func (m *BTreeMap[K, V]) Extend(it iter.Seq2[K, V]) {
	for k, v := range it {
		m.Insert(k, v)
	}
}
