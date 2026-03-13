package trie

import (
	"iter"
)

// Iter returns an iterator over all key-value pairs in the trie.
//
// Returns:
//   - An iter.Seq2[[]K, V] that yields (key, value) pairs.
//
// Time Complexity: O(n * L) where n is the number of entries and L is average key length.
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
//
// Returns:
//   - An iter.Seq2[[]K, *V] that yields (key, pointer to value) pairs.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n * L) where n is the number of entries and L is average key length.
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
//
// Returns:
//   - An iter.Seq[[]K] that yields all keys.
//
// Time Complexity: O(n * L) where n is the number of entries and L is average key length.
func (m *TrieMap[K, V]) Keys() iter.Seq[[]K] {
	return func(yield func([]K) bool) {
		for k := range m.makeIter(make([]K, 0)) {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator over all values in the trie.
//
// Returns:
//   - An iter.Seq[V] that yields all values.
//
// Time Complexity: O(n * L) where n is the number of entries and L is average key length.
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
//
// Returns:
//   - An iter.Seq[*V] that yields pointers to all values.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(n * L) where n is the number of entries and L is average key length.
func (m *TrieMap[K, V]) ValuesMut() iter.Seq[*V] {
	return func(yield func(*V) bool) {
		for _, n := range m.makeIter(make([]K, 0)) {
			if !yield(&n.value) {
				return
			}
		}
	}
}

// KeysByPrefix returns an iterator over all keys with the given prefix.
//
// Parameters:
//   - prefix: The prefix to filter keys by.
//
// Returns:
//   - An iter.Seq[[]K] that yields keys starting with the prefix.
//
// Time Complexity: O(p + m*L) where p is prefix length and m is number of matches.
func (m *TrieMap[K, V]) KeysByPrefix(prefix []K) iter.Seq[[]K] {
	return func(yield func([]K) bool) {
		for k := range m.makeIter(prefix) {
			if !yield(k) {
				return
			}
		}
	}
}

// ValuesByPrefix returns an iterator over all values with the given prefix.
//
// Parameters:
//   - prefix: The prefix to filter values by.
//
// Returns:
//   - An iter.Seq[V] that yields values with keys starting with the prefix.
//
// Time Complexity: O(p + m*L) where p is prefix length and m is number of matches.
func (m *TrieMap[K, V]) ValuesByPrefix(prefix []K) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, n := range m.makeIter(prefix) {
			if !yield(n.value) {
				return
			}
		}
	}
}

// IterByPrefix returns an iterator over all key-value pairs with the given prefix.
//
// Parameters:
//   - prefix: The prefix to filter key-value pairs by.
//
// Returns:
//   - An iter.Seq2[[]K, V] that yields (key, value) pairs with keys starting with the prefix.
//
// Time Complexity: O(p + m*L) where p is prefix length and m is number of matches.
func (m *TrieMap[K, V]) IterByPrefix(prefix []K) iter.Seq2[[]K, V] {
	return func(yield func([]K, V) bool) {
		for k, n := range m.makeIter(prefix) {
			if !yield(k, n.value) {
				return
			}
		}
	}
}

// IterMutByPrefix returns a mutable iterator over all key-value pairs with the given prefix.
//
// Parameters:
//   - prefix: The prefix to filter key-value pairs by.
//
// Returns:
//   - An iter.Seq2[[]K, *V] that yields (key, pointer to value) pairs with keys starting with the prefix.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(p + m*L) where p is prefix length and m is number of matches.
func (m *TrieMap[K, V]) IterMutByPrefix(prefix []K) iter.Seq2[[]K, *V] {
	return func(yield func([]K, *V) bool) {
		for k, n := range m.makeIter(prefix) {
			if !yield(k, &n.value) {
				return
			}
		}
	}
}

// ValuesMutByPrefix returns a mutable iterator over all values with the given prefix.
//
// Parameters:
//   - prefix: The prefix to filter values by.
//
// Returns:
//   - An iter.Seq[*V] that yields pointers to values with keys starting with the prefix.
//   - The yielded values can be modified in place.
//
// Time Complexity: O(p + m*L) where p is prefix length and m is number of matches.
func (m *TrieMap[K, V]) ValuesMutByPrefix(prefix []K) iter.Seq[*V] {
	return func(yield func(*V) bool) {
		for _, n := range m.makeIter(prefix) {
			if !yield(&n.value) {
				return
			}
		}
	}
}

// Extend inserts all key-value pairs from the iterator into the trie.
//
// Parameters:
//   - it: An iterator yielding key/value pairs to insert.
//
// Behavior:
//   - If a key already exists, its value will be updated.
//
// Time Complexity: O(N*L) where N is the number of pairs and L is average key length.
func (m *TrieMap[K, V]) Extend(it iter.Seq2[[]K, V]) {
	for k, v := range it {
		m.Insert(k, v)
	}
}

// makeIter returns an iterator over nodes with the given prefix.
// It performs a depth-first search (DFS) starting from the node corresponding to the prefix.
func (m *TrieMap[K, V]) makeIter(prefix []K) iter.Seq2[[]K, *node[K, V]] {
	return func(yield func([]K, *node[K, V]) bool) {
		if m.root == nil {
			return
		}

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
				return
			}
		}

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
				keyCopy := make([]K, len(frame.prefix))
				copy(keyCopy, frame.prefix)
				if !yield(keyCopy, frame.n) {
					return
				}
			}

			for i := len(frame.n.children) - 1; i >= 0; i-- {
				child := frame.n.children[i]
				newPrefix := append(append([]K{}, frame.prefix...), child.key)
				stack = append(stack, stackFrame{n: child, prefix: newPrefix})
			}
		}
	}
}
