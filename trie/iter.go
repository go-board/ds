// Package trie implements a generic trie (prefix tree) data structure.
package trie

import (
	"iter"
)

// Iter returns an iterator over all key-value pairs in the trie.
// Return value:
//   - Iterator over key-value pairs, of type iter.Seq2[[]K, V]
func (m *TrieMap[K, V]) Iter() iter.Seq2[[]K, V] {
	return func(yield func([]K, V) bool) {
		for k, n := range m.makeIter(make([]K, 0)) {
			if !yield(k, n.value) {
				return
			}
		}
	}
}

// IterMut returns a mutable iterator over all key-value pairs in the trie.
// Return value:
//   - Mutable iterator over key-value pairs, of type iter.Seq2[[]K, *V]
func (m *TrieMap[K, V]) IterMut() iter.Seq2[[]K, *V] {
	return func(yield func([]K, *V) bool) {
		for k, n := range m.makeIter(make([]K, 0)) {
			if !yield(k, &n.value) {
				return
			}
		}
	}
}

// Keys returns an iterator over all keys in the trie.
// Return value:
//   - Iterator over keys, of type iter.Seq[[]K]
func (m *TrieMap[K, V]) Keys() iter.Seq[[]K] {
	return func(yield func([]K) bool) {
		for k, _ := range m.makeIter(make([]K, 0)) {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator over all values in the trie.
// Return value:
//   - Iterator over values, of type iter.Seq[V]
func (m *TrieMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, n := range m.makeIter(make([]K, 0)) {
			if !yield(n.value) {
				return
			}
		}
	}
}

// ValuesMut returns a mutable iterator over all values in the trie.
// Return value:
//   - Mutable iterator over values, of type iter.Seq[*V]
//
// ValuesMut returns an iterator over all values in the trie, allowing mutation.
// Return value:
//   - Iterator over value pointers, of type iter.Seq[*V]
func (m *TrieMap[K, V]) ValuesMut() iter.Seq[*V] {
	return func(yield func(*V) bool) {
		for _, n := range m.makeIter(make([]K, 0)) {
			if !yield(&n.value) {
				return
			}
		}
	}
}

// KeysByPrefix returns all keys in the trie with the given prefix as an iterator.
// Parameters:
//   - prefix: The prefix to filter keys by
//
// Return value:
//   - Iterator over keys with the given prefix, of type iter.Seq[[]K]
func (m *TrieMap[K, V]) KeysByPrefix(prefix []K) iter.Seq[[]K] {
	return func(yield func([]K) bool) {
		for k, _ := range m.makeIter(prefix) {
			if !yield(k) {
				return
			}
		}
	}
}

// ValuesByPrefix returns all values in the trie with the given prefix as an iterator.
// Parameters:
//   - prefix: The prefix to filter values by
//
// Return value:
//   - Iterator over values with the given prefix, of type iter.Seq[V]
func (m *TrieMap[K, V]) ValuesByPrefix(prefix []K) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, n := range m.makeIter(prefix) {
			if !yield(n.value) {
				return
			}
		}
	}
}

// IterByPrefix returns an iterator over all key-value pairs in the trie with the given prefix.
// Parameters:
//   - prefix: The prefix to filter key-value pairs by
//
// Return value:
//   - Iterator over key-value pairs with the given prefix, of type iter.Seq2[[]K, V]
func (m *TrieMap[K, V]) IterByPrefix(prefix []K) iter.Seq2[[]K, V] {
	return func(yield func([]K, V) bool) {
		for k, n := range m.makeIter(prefix) {
			if !yield(k, n.value) {
				return
			}
		}
	}
}

// IterMutByPrefix returns a mutable iterator over all key-value pairs in the trie with the given prefix.
// Parameters:
//   - prefix: The prefix to filter key-value pairs by
//
// Return value:
//   - Mutable iterator over key-value pairs with the given prefix, of type iter.Seq2[[]K, *V]
func (m *TrieMap[K, V]) IterMutByPrefix(prefix []K) iter.Seq2[[]K, *V] {
	return func(yield func([]K, *V) bool) {
		for k, n := range m.makeIter(prefix) {
			if !yield(k, &n.value) {
				return
			}
		}
	}
}

// ValuesMutByPrefix returns a mutable iterator over all values in the trie with the given prefix.
// Parameters:
//   - prefix: The prefix to filter values by
//
// Return value:
//   - Mutable iterator over values with the given prefix, of type iter.Seq[*V]
func (m *TrieMap[K, V]) ValuesMutByPrefix(prefix []K) iter.Seq[*V] {
	return func(yield func(*V) bool) {
		for _, n := range m.makeIter(prefix) {
			if !yield(&n.value) {
				return
			}
		}
	}
}

// Extend adds all key-value pairs from the provided iterator to the trie.
// For each key-value pair in the iterator:
//   - If the key already exists, its value is updated
//   - Otherwise, a new key-value pair is added
//
// Parameters:
//   - it: An iterator over key-value pairs to add to the trie
//
// Time complexity: O(N*L), where N is the number of key-value pairs in the iterator and L is the average length of the keys
func (m *TrieMap[K, V]) Extend(it iter.Seq2[[]K, V]) {
	for k, v := range it {
		m.Insert(k, v)
	}
}

// makeIter returns an iterator over nodes in the trie with the given prefix.
// It performs a depth-first search (DFS) starting from the node corresponding to the prefix.
// Return value:
//   - Iterator over key-node pairs, of type iter.Seq2[[]K, *node[K, V]]
func (m *TrieMap[K, V]) makeIter(prefix []K) iter.Seq2[[]K, *node[K, V]] {
	return func(yield func([]K, *node[K, V]) bool) {
		if m.root == nil {
			return
		}

		// Find the node corresponding to the prefix
		n := m.root
		currentPrefix := make([]K, 0, len(prefix))

		for _, k := range prefix {
			found := false
			for _, child := range n.children {
				if m.comparator(child.key, k) == 0 {
					n = child
					currentPrefix = append(currentPrefix, k)
					found = true
					break
				}
			}

			if !found {
				// Prefix not found
				return
			}
		}

		// Now perform DFS starting from the node with the prefix
		type stackFrame struct {
			n      *node[K, V]
			prefix []K
		}

		stack := []stackFrame{
			{n: n, prefix: append([]K{}, currentPrefix...)},
		}

		for len(stack) > 0 {
			frame := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if frame.n.hasValue {
				// Make a copy of the prefix to avoid modifying the original
				keyCopy := make([]K, len(frame.prefix))
				copy(keyCopy, frame.prefix)
				if !yield(keyCopy, frame.n) {
					return
				}
			}

			// Push children in reverse order to ensure correct traversal order
			for i := len(frame.n.children) - 1; i >= 0; i-- {
				child := frame.n.children[i]
				// Add the child's key to the prefix
				newPrefix := append(append([]K{}, frame.prefix...), child.key)
				stack = append(stack, stackFrame{n: child, prefix: newPrefix})
			}
		}
	}
}
