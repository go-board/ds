package skipmap

import (
	"iter"

	diter "github.com/go-board/ds/internal/iter"
)

// Keys returns an iterator over all keys in the SkipMap in ascending order.
//
// It returns an `iter.Seq[K]` that yields each key.
func (m *SkipMap[K, V]) Keys() iter.Seq[K] {
	return diter.Keys(m.Iter())
}

// Values returns an iterator over all values in the SkipMap in ascending key order.
//
// It returns an `iter.Seq[V]` that yields each value.
func (m *SkipMap[K, V]) Values() iter.Seq[V] {
	return diter.Values(m.Iter())
}

// ValuesMut returns a mutable iterator over all values in the SkipMap in ascending key order.
//
// It returns an `iter.Seq[*V]` that yields a pointer to each value.
func (m *SkipMap[K, V]) ValuesMut() iter.Seq[*V] {
	return diter.Values(m.IterMut())
}

// Iter returns an iterator over all key/value pairs in the SkipMap in ascending key order.
//
// It returns an `iter.Seq2[K, V]` that yields each (key, value) pair.
func (m *SkipMap[K, V]) Iter() iter.Seq2[K, V] {
	return diter.Split(m.rangeNode(nil, nil), (*node[K, V]).kv)
}

// IterMut returns a mutable iterator over all key/value pairs in the SkipMap.
//
// It returns an `iter.Seq2[K, *V]` that yields (key, *value) pairs for in-place modification.
func (m *SkipMap[K, V]) IterMut() iter.Seq2[K, *V] {
	return diter.Split(m.rangeNode(nil, nil), (*node[K, V]).kvMut)
}

// Range returns an iterator over key/value pairs whose keys fall in [lower, upper).
//
// The `lower` bound is inclusive and the `upper` bound is exclusive. A nil bound
// indicates no limit on that side.
//
// It returns an `iter.Seq2[K, V]` that yields all pairs in ascending key order within the range.
func (m *SkipMap[K, V]) Range(lowerBound, upperBound *K) iter.Seq2[K, V] {
	return diter.Split(m.rangeNode(lowerBound, upperBound), (*node[K, V]).kv)
}

// RangeMut returns a mutable range iterator over key/value pairs whose keys are in [lower, upper).
//
// It returns an `iter.Seq2[K, *V]` that yields (key, *value) pairs, allowing modification of values.
func (m *SkipMap[K, V]) RangeMut(lowerBound, upperBound *K) iter.Seq2[K, *V] {
	return diter.Split(m.rangeNode(lowerBound, upperBound), (*node[K, V]).kvMut)
}

// IterBack returns a reverse iterator that yields key/value pairs from largest to smallest key.
//
// It returns an `iter.Seq2[K, V]` that produces pairs in descending key order.
func (m *SkipMap[K, V]) IterBack() iter.Seq2[K, V] {
	return diter.Split(m.iterNodeBack(), (*node[K, V]).kv)
}

func (m *SkipMap[K, V]) iterNodeBack() iter.Seq[*node[K, V]] {
	// Lazy reverse traversal: start from the last node and repeatedly find the
	// predecessor using skip-list search logic. This avoids allocating a full
	// slice of nodes up-front.
	return func(yield func(*node[K, V]) bool) {
		if m == nil || m.head == nil {
			return
		}

		// find last node at level 0
		last := m.head.next[0]
		if last == nil {
			return
		}
		for last.next[0] != nil {
			last = last.next[0]
		}

		// helper: find predecessor of a node with given key
		predecessor := func(key K) *node[K, V] {
			current := m.head
			// traverse from top level down to level 0 as in Insert/Get
			for i := m.level; i >= 0; i-- {
				for current.next[i] != nil && m.comparator(current.next[i].Key, key) < 0 {
					current = current.next[i]
				}
			}
			// current is the predecessor (largest node with key < given key),
			// may be head if there's no smaller node
			if current == m.head {
				// if head.next[0] is the node with key, predecessor is nil
				if current.next[0] != nil && m.comparator(current.next[0].Key, key) < 0 {
					return current
				}
				// explicit nil predecessor for smallest element
				if current.next[0] != nil && m.comparator(current.next[0].Key, key) >= 0 {
					return nil
				}
			}
			return current
		}

		cur := last
		for cur != nil {
			if !yield(cur) {
				return
			}
			// find predecessor of cur.Key
			prev := predecessor(cur.Key)
			if prev == nil || prev == m.head {
				// if prev is head (no real node before cur) then stop
				return
			}
			cur = prev
		}
	}
}

// IterBackMut returns a mutable reverse iterator that yields (key, *value) pairs
// from largest to smallest key.
//
// It returns an `iter.Seq2[K, *V]` for in-place modification while iterating.
func (m *SkipMap[K, V]) IterBackMut() iter.Seq2[K, *V] {
	return diter.Split(m.iterNodeBack(), (*node[K, V]).kvMut)
}

// KeysBack returns an iterator over keys in descending order.
//
// It returns an `iter.Seq[K]` that yields keys from largest to smallest.
func (m *SkipMap[K, V]) KeysBack() iter.Seq[K] {
	return diter.Keys(m.IterBack())
}

// ValuesBack returns an iterator over values in descending key order.
//
// It returns an `iter.Seq[V]` that yields values corresponding to keys from largest to smallest.
func (m *SkipMap[K, V]) ValuesBack() iter.Seq[V] {
	return diter.Values(m.IterBack())
}

// ValuesBackMut returns a mutable reverse iterator over values in descending key order.
//
// It returns an `iter.Seq[*V]` that yields pointers to values for in-place modification.
func (m *SkipMap[K, V]) ValuesBackMut() iter.Seq[*V] {
	return diter.Values(diter.Split(m.iterNodeBack(), (*node[K, V]).kvMut))
}

func (m *SkipMap[K, V]) rangeNode(lowerBound, upperBound *K) iter.Seq[*node[K, V]] {
	return func(yield func(*node[K, V]) bool) {
		current := m.head.next[0]

		// If there's a lower bound, move to the lower bound position
		if lowerBound != nil {
			for current != nil && m.comparator(current.Key, *lowerBound) < 0 {
				current = current.next[0]
			}
		}

		// Iterate through elements in the range
		for current != nil {
			// If there's an upper bound and current key >= upper bound, stop
			if upperBound != nil && m.comparator(current.Key, *upperBound) >= 0 {
				break
			}
			if !yield(current) {
				return
			}
			current = current.next[0]
		}
	}
}

// Extend adds another iterable collection of key-value pairs to the current map.
// Parameters:
//   - it: An iterator providing key-value pairs
//
// Behavior:
//   - For each key-value pair, updates the value if the key already exists, otherwise adds a new key-value pair
func (m *SkipMap[K, V]) Extend(it iter.Seq2[K, V]) {
	for k, v := range it {
		m.Insert(k, v)
	}
}
