package hashmap

import (
	"iter"

	diter "github.com/go-board/ds/internal/iter"
)

// Keys returns an iterator over all keys in the hash map.
// Return value:
//   - Iterator over keys, of type iter.Seq[K]
//
// Time complexity: O(n) for iterating all elements
func (hm *HashMap[K, V, H]) Keys() iter.Seq[K] {
	return diter.Keys(hm.Iter())
}

// Values returns an iterator over all values in the hash map.
// Return value:
//   - Iterator over values, of type iter.Seq[V]
//
// Time complexity: O(n) for iterating all elements
func (hm *HashMap[K, V, H]) Values() iter.Seq[V] {
	return diter.Values(hm.Iter())
}

// ValuesMut returns a mutable iterator over all values in the hash map.
// Return value:
//   - Iterator over mutable values, of type iter.Seq[*V]
//
// Time complexity: O(n) for iterating all elements
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
// Return value:
//   - Iterator over key-value pairs, of type iter.Seq2[K, V]
//
// Time complexity: O(n) for iterating all elements
func (hm *HashMap[K, V, H]) Iter() iter.Seq2[K, V] {
	return diter.Split(hm.iterNode(), func(n *node[K, V]) (K, V) {
		return n.key, n.value
	})
}

// IterMut returns a mutable iterator over all key-value pairs in the hash map.
// Return value:
//   - Mutable iterator over key-value pairs, of type iter.Seq2[K, *V]
//
// Behavior:
//   - Allows modifying values during iteration
//
// Time complexity: O(n) for iterating all elements
func (hm *HashMap[K, V, H]) IterMut() iter.Seq2[K, *V] {
	return diter.Split(hm.iterNode(), func(n *node[K, V]) (K, *V) {
		return n.key, &n.value
	})
}

// Extend adds another iterable key-value pair collection to the current hash map.
// Parameters:
//   - it: Iterator providing key-value pairs
//
// Behavior:
//   - For each key-value pair, if the key exists, update its value; otherwise add a new key-value pair
func (hm *HashMap[K, V, H]) Extend(it iter.Seq2[K, V]) {
	if hm.deletedCount > hm.size {
		hm.Compact()
	}
	for k, v := range it {
		hm.Insert(k, v)
	}
}
