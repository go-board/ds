package hashmap

import (
	"iter"

	diter "github.com/go-board/ds/internal/iter"
)

// Keys returns an iterator over all keys in the hash map.
//
// Returns:
//   - An iter.Seq[K] that yields all keys.
//
// Time Complexity: O(n)
func (hm *HashMap[K, V, H]) Keys() iter.Seq[K] {
	return diter.Keys(hm.Iter())
}

// Values returns an iterator over all values in the hash map.
//
// Returns:
//   - An iter.Seq[V] that yields all values.
//
// Time Complexity: O(n)
func (hm *HashMap[K, V, H]) Values() iter.Seq[V] {
	return diter.Values(hm.Iter())
}

// ValuesMut returns a mutable iterator over all values in the hash map.
//
// Returns:
//   - An iter.Seq[*V] that yields pointers to all values.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (hm *HashMap[K, V, H]) ValuesMut() iter.Seq[*V] {
	return diter.Values(hm.IterMut())
}

func (hm *HashMap[K, V, H]) iterNode() iter.Seq[*node[K, V]] {
	return func(yield func(*node[K, V]) bool) {
		for _, bucket := range hm.buckets {
			for _, node := range bucket.nodes {
				if !node.deleted {
					if !yield(node) {
						return
					}
				}
			}
		}
	}
}

// Iter returns an iterator over all key-value pairs in the hash map.
//
// Returns:
//   - An iter.Seq2[K, V] that yields (key, value) pairs.
//
// Time Complexity: O(n)
func (hm *HashMap[K, V, H]) Iter() iter.Seq2[K, V] {
	return diter.Split(hm.iterNode(), func(n *node[K, V]) (K, V) {
		return n.key, n.value
	})
}

// IterMut returns a mutable iterator over all key-value pairs in the hash map.
//
// Returns:
//   - An iter.Seq2[K, *V] that yields (key, pointer to value) pairs.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n)
func (hm *HashMap[K, V, H]) IterMut() iter.Seq2[K, *V] {
	return diter.Split(hm.iterNode(), func(n *node[K, V]) (K, *V) {
		return n.key, &n.value
	})
}

// Extend inserts all key/value pairs from the iterator into the map.
//
// Parameters:
//   - it: An iterator yielding key/value pairs to insert.
//
// Behavior:
//   - If a key already exists, its value will be updated.
func (hm *HashMap[K, V, H]) Extend(it iter.Seq2[K, V]) {
	if hm.deletedCount > hm.size {
		hm.Compact()
	}
	for k, v := range it {
		hm.Insert(k, v)
	}
}
