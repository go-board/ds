// Package trie implements a generic trie (prefix tree) data structure.
package trie

import (
	"cmp"
)

// node represents a node in the trie.
// It contains a key (which is a single element of the path), a value (if the node represents the end of a key),
// and a list of child nodes.
//   - key: The key element stored in this node
//   - children: List of child nodes
//   - value: The value stored at this node (if it's a leaf with a value)
//   - hasValue: Indicates if this node has a value
//   - isLeaf: Indicates if this node is a leaf (no children)
type node[K, V any] struct {
	key      K
	children []*node[K, V]
	value    V
	hasValue bool
	isLeaf   bool
}

// TrieMap represents a trie (prefix tree) data structure that maps sequences of keys to values.
// It supports efficient prefix-based operations.
//   - root: The root node of the trie
//   - size: The number of key-value pairs in the trie
//   - comparator: Function used to compare keys
//
// Type parameters:
//   - K: The type of individual key elements
//   - V: The type of values stored in the trie
type TrieMap[K, V any] struct {
	root       *node[K, V]
	size       int
	comparator func(K, K) int
}

// New creates a new TrieMap with the specified key comparator.
// Parameters:
//   - comparator: Function used to compare individual key elements
//
// Return value:
//   - A new empty TrieMap instance
//
// Time complexity: O(1)
func New[K, V any](comparator func(K, K) int) *TrieMap[K, V] {
	return &TrieMap[K, V]{
		root:       newNode[K, V](),
		size:       0,
		comparator: comparator,
	}
}

// NewOrdered creates a new TrieMap with a default comparator for ordered types.
// Return value:
//   - A new empty TrieMap instance with a default comparator
func NewOrdered[K cmp.Ordered, V any]() *TrieMap[K, V] {
	return New[K, V](cmp.Compare[K])
}

// newNode creates a new empty node with default values.
func newNode[K, V any]() *node[K, V] {
	return &node[K, V]{
		children: make([]*node[K, V], 0),
		hasValue: false,
		isLeaf:   true,
	}
}

// findNode finds the node corresponding to the given key path.
// Parameters:
//   - key: The key path to search for
//
// Return values:
//   - The node corresponding to the key path, or nil if not found
//   - The parent node of the found node, or nil if the node is the root or not found
//   - The index of the found node in its parent's children list, or -1 if not found
func (m *TrieMap[K, V]) findNode(key []K) (*node[K, V], *node[K, V], int) {
	if len(key) == 0 {
		return m.root, nil, -1
	}

	current := m.root
	parent := m.root
	childIndex := -1

	for i, k := range key {
		found := false
		for j, child := range current.children {
			if m.comparator(child.key, k) == 0 {
				parent = current
				current = child
				childIndex = j
				found = true
				break
			}
		}

		if !found {
			return nil, parent, -1
		}

		// Check if we've reached the end of the key
		if i == len(key)-1 {
			return current, parent, childIndex
		}
	}

	return current, parent, childIndex
}

// Get retrieves the value associated with the given key.
// Parameters:
//   - key: The key to search for
//
// Return values:
//   - The value associated with the key, or the zero value of V if not found
//   - A boolean indicating whether the key was found
//
// Time complexity: O(L), where L is the length of the key
func (m *TrieMap[K, V]) Get(key []K) (V, bool) {
	n, _, _ := m.findNode(key)
	if n == nil || !n.hasValue {
		var zero V
		return zero, false
	}
	return n.value, true
}

// GetMut retrieves a pointer to the value associated with the given key.
// Parameters:
//   - key: The key to search for
//
// Return values:
//   - A pointer to the value associated with the key, or nil if not found
//
// Time complexity: O(L), where L is the length of the key
func (m *TrieMap[K, V]) GetMut(key []K) *V {
	n, _, _ := m.findNode(key)
	if n == nil || !n.hasValue {
		return nil
	}
	return &n.value
}

// Insert adds a key-value pair to the trie. If the key already exists, its value is updated.
// Parameters:
//   - key: The key to insert
//   - value: The value to associate with the key
//
// Return value:
//   - A boolean indicating whether the key was newly added (true) or updated (false)
//
// Time complexity: O(L), where L is the length of the key
func (m *TrieMap[K, V]) Insert(key []K, value V) bool {
	// Special case for empty key
	if len(key) == 0 {
		if !m.root.hasValue {
			m.root.hasValue = true
			m.size++
		}
		m.root.value = value
		return !m.root.hasValue // Return true if it was a new insertion
	}

	current := m.root

	for i, k := range key {
		// Search for the child with the current key element
		found := false
		for _, child := range current.children {
			if m.comparator(child.key, k) == 0 {
				current = child
				found = true
				break
			}
		}

		if !found {
			// Create a new node for this key element
			newChild := newNode[K, V]()
			newChild.key = k
			current.children = append(current.children, newChild)
			current.isLeaf = false
			current = newChild
		}

		// If we've reached the end of the key, set the value
		if i == len(key)-1 {
			if !current.hasValue {
				m.size++
			}
			current.value = value
			current.hasValue = true
			return !current.hasValue // Return true if it was a new insertion
		}
	}

	// This should not be reachable due to the check at the end of the key
	if !current.hasValue {
		m.size++
	}
	current.value = value
	current.hasValue = true
	return !current.hasValue
}

// Remove removes a key-value pair from the trie.
// Parameters:
//   - key: The key to remove
//
// Return values:
//   - If the key exists, returns the deleted value and true
//   - If the key does not exist, returns the zero value and false
//
// Time complexity: O(L), where L is the length of the key
func (m *TrieMap[K, V]) Remove(key []K) (V, bool) {
	// Special case for empty key
	if len(key) == 0 {
		if !m.root.hasValue {
			var zero V
			return zero, false
		}
		value := m.root.value
		m.root.hasValue = false
		m.size--
		return value, true
	}

	// Find the node and its parent
	n, parent, childIndex := m.findNode(key)
	if n == nil || !n.hasValue {
		var zero V
		return zero, false
	}

	// Save the value to return
	value := n.value
	// Remove the value from the node
	n.hasValue = false
	m.size--

	// If the node has children, we can't remove it
	if !n.isLeaf {
		return value, true
	}

	// If the node is a leaf, we need to remove it and any ancestors that are no longer needed
	current := n
	currentParent := parent
	currentIndex := childIndex

	for currentParent != nil && !current.hasValue && current.isLeaf {
		// Remove current from parent's children
		currentParent.children = append(currentParent.children[:currentIndex], currentParent.children[currentIndex+1:]...)
		if len(currentParent.children) == 0 {
			currentParent.isLeaf = true
		}

		// Move up to the next parent
		nextParent, nextIndex := m.findParentOfNode(currentParent)
		current = currentParent
		currentParent = nextParent
		currentIndex = nextIndex
	}

	return value, true
}

// Helper function to find the parent of a node
func (m *TrieMap[K, V]) findParentOfNode(n *node[K, V]) (*node[K, V], int) {
	// This is a helper function to find the parent of a node
	// It's used in the Remove method to traverse up the tree
	// For simplicity, we'll just do a BFS here
	if n == m.root {
		return nil, -1
	}

	type nodeInfo struct {
		n     *node[K, V]
		index int
	}

	queue := []nodeInfo{
		{m.root, -1},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for i, child := range current.n.children {
			if child == n {
				return current.n, i
			}
			queue = append(queue, nodeInfo{n: child, index: i})
		}
	}

	return nil, -1
}

// ContainsKey checks if the trie contains the given key.
// Parameters:
//   - key: The key to check for
//
// Return value:
//   - A boolean indicating whether the key exists in the trie
//
// Time complexity: O(L), where L is the length of the key
func (m *TrieMap[K, V]) ContainsKey(key []K) bool {
	n, _, _ := m.findNode(key)
	return n != nil && n.hasValue
}

// IsEmpty checks if the trie is empty.
// Return value:
//   - A boolean indicating whether the trie is empty
//
// Time complexity: O(1)
func (m *TrieMap[K, V]) IsEmpty() bool {
	return m.size == 0
}

// Len returns the number of key-value pairs in the trie.
// Return value:
//   - The number of key-value pairs
//
// Time complexity: O(1)
func (m *TrieMap[K, V]) Len() int {
	return m.size
}

// Clone creates a deep copy of the trie.
// Return value:
//   - A new TrieMap that is a deep copy of the original
//
// Time complexity: O(N), where N is the number of nodes in the trie
func (m *TrieMap[K, V]) Clone() *TrieMap[K, V] {
	clone := &TrieMap[K, V]{
		root:       cloneNode(m.root, m.comparator),
		size:       m.size,
		comparator: m.comparator,
	}
	return clone
}

// cloneNode creates a deep copy of a node and all its descendants.
func cloneNode[K, V any](n *node[K, V], comparator func(K, K) int) *node[K, V] {
	clone := &node[K, V]{
		key:      n.key,
		value:    n.value,
		hasValue: n.hasValue,
		isLeaf:   n.isLeaf,
	}

	// Clone children if they exist
	if len(n.children) > 0 {
		clone.children = make([]*node[K, V], len(n.children))
		for i, child := range n.children {
			clone.children[i] = cloneNode(child, comparator)
		}
	} else {
		clone.children = make([]*node[K, V], 0)
	}

	return clone
}

// Entry creates a new Entry for the given key.
// This allows for convenient operations on a specific key.
//
// Parameters:
//   - key: The key to create an Entry for
//
// Returns:
//   - A new Entry instance for the given key
func (m *TrieMap[K, V]) Entry(key []K) Entry[K, V] {
	return Entry[K, V]{
		m:   m,
		key: key,
	}
}
