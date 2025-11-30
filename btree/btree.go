// Package btree implements a generic B-tree data structure.
//
// A B-tree is a self-balancing tree data structure that maintains sorted data and allows efficient insertion,
// deletion, and lookup operations. It's particularly suitable for external storage systems like databases
// and file systems because it reduces disk I/O operations.
//
// Features:
//   - Ordered storage with support for sequential element traversal
//   - Efficient insertion, deletion, and lookup operations with O(log n) time complexity
//   - Support for range queries
//   - Support for retrieving first and last elements
//
// Usage Example:
//
//	// Create an ordered B-tree for integers
//	tree := btree.NewOrdered[int]()
//
//	// Insert elements
//	tree.Insert(5)
//	tree.Insert(3)
//	tree.Insert(7)
//
//	// Check if an element exists
//	val, found := tree.Search(5) // 5, true
//
//	// Traverse elements (in order)
//	for val := range tree.Iter() {
//	    fmt.Println(val) // Output: 3, 5, 7
//	}
//
//	// Get elements within a range
//	lower := 4
//	upper := 8
//	for val := range tree.Range(&lower, &upper) {
//	    fmt.Println(val) // Output: 5, 7
//	}
package btree

import (
	"cmp"
)

// defaultOrder is the default B-tree order, set to 3 according to common implementations
const defaultOrder = 3

// BTree is a generic B-tree implementation
// A B-tree is a self-balancing tree data structure that maintains sorted data and allows efficient insertion,
// deletion, and lookup operations, all with O(log n) time complexity
type BTree[E any] struct {
	root       *node[E]       // Root node
	order      int            // B-tree order, initialized with default value
	comparator func(E, E) int // Element comparison function
	size       int            // Number of elements in the tree
}

// node represents a node in the B-tree
type node[E any] struct {
	keys     []E        // Array of keys in the node
	children []*node[E] // Array of child node pointers
	isLeaf   bool       // Whether the node is a leaf node
}

// New creates a new B-tree with the given comparator function
// The B-tree order is set to a default value and cannot be configured by users
//
// Parameters:
//   - comparator: A function to compare elements, returning a negative number if a < b, zero if a == b,
//     and a positive number if a > b
//
// Returns:
//   - A newly created BTree instance
//
// Note: The comparator function cannot be nil, otherwise it will panic
//
// Time Complexity: O(1)
func New[E any](comparator func(E, E) int) *BTree[E] {
	if comparator == nil {
		panic("comparator function cannot be nil")
	}

	return &BTree[E]{
		root:       &node[E]{keys: make([]E, 0), children: make([]*node[E], 0), isLeaf: true},
		order:      defaultOrder,
		comparator: comparator,
		size:       0,
	}
}

// NewOrdered creates a new B-tree for ordered element types
// This is a convenience function that uses cmp.Compare as the comparator
//
// Type Parameters:
//   - E: The element type, must implement the cmp.Ordered interface
//
// Returns:
//   - A newly created BTree instance
//
// Time Complexity: O(1)
func NewOrdered[E cmp.Ordered]() *BTree[E] {
	return New(cmp.Compare[E])
}

// Len returns the number of elements in the B-tree
//
// Returns:
//   - The count of elements in the tree
//
// Time Complexity: O(1)
func (t *BTree[E]) Len() int {
	return t.size
}

// Insert adds an element to the B-tree
//
// Parameters:
//   - key: The element value to insert
//
// Time Complexity: O(log n)
func (t *BTree[E]) Insert(key E) {
	// If root node is full, split it
	if len(t.root.keys) == 2*t.order-1 {
		newRoot := &node[E]{keys: make([]E, 0), children: make([]*node[E], 0), isLeaf: false}
		newRoot.children = append(newRoot.children, t.root)
		t.splitChild(newRoot, 0)
		t.root = newRoot
	}

	// Insert into non-full root node
	t.insertNonFull(t.root, key)
	t.size++
}

// insertNonFull inserts an element into a non-full node
// This is an internal helper method and not exposed to users
func (t *BTree[E]) insertNonFull(n *node[E], key E) {
	i := len(n.keys) - 1

	// If it's a leaf node, insert directly at the appropriate position
	if n.isLeaf {
		// Extend key array
		n.keys = append(n.keys, key)

		// Move keys to find the correct insertion position
		for i >= 0 && t.comparator(key, n.keys[i]) < 0 {
			n.keys[i+1] = n.keys[i]
			i--
		}
		n.keys[i+1] = key
	} else {
		// Find the appropriate child node
		for i >= 0 && t.comparator(key, n.keys[i]) < 0 {
			i--
		}
		i++

		// If child node is full, split it first
		if len(n.children[i].keys) == 2*t.order-1 {
			t.splitChild(n, i)
			// After splitting, decide which child to insert into
			if t.comparator(key, n.keys[i]) > 0 {
				i++
			}
		}
		t.insertNonFull(n.children[i], key)
	}
}

// splitChild splits the ith child of node n
// This is an internal helper method and not exposed to users
func (t *BTree[E]) splitChild(n *node[E], i int) {
	order := t.order
	child := n.children[i]

	// Create new node to store the right half after splitting
	newChild := &node[E]{
		keys:     make([]E, order-1),
		children: make([]*node[E], 0),
		isLeaf:   child.isLeaf,
	}

	// Move middle key up to parent node
	medianKey := child.keys[order-1]

	// Move right half of keys to new node
	copy(newChild.keys, child.keys[order:])
	child.keys = child.keys[:order-1]

	// If not a leaf node, also move child node pointers
	if !child.isLeaf {
		newChild.children = append(newChild.children, child.children[order:]...)
		child.children = child.children[:order]
	}

	// Insert new node into parent's children list
	n.children = append(n.children, nil)
	for j := len(n.children) - 1; j > i+1; j-- {
		n.children[j] = n.children[j-1]
	}
	n.children[i+1] = newChild

	// Insert middle key into parent node
	n.keys = append(n.keys, medianKey)
	for j := len(n.keys) - 1; j > i; j-- {
		n.keys[j] = n.keys[j-1]
	}
	n.keys[i] = medianKey
}

// Search looks up an element in the B-tree, returning the value and true if found, or zero value and false if not found
//
// Parameters:
//   - key: The element value to search for
//
// Returns:
//   - If the element exists, returns the element value and true
//   - If the element doesn't exist, returns zero value and false
//
// Time Complexity: O(log n)
func (t *BTree[E]) Search(key E) (E, bool) {
	var zero E
	result, found := t.searchNode(t.root, key)
	if !found {
		return zero, false
	}
	return result, true
}

// searchNode searches for an element in the node and its subtree
// This is an internal helper method and not exposed to users
func (t *BTree[E]) searchNode(n *node[E], key E) (E, bool) {
	i := 0
	// Find the position of the key in the current node
	for i < len(n.keys) && t.comparator(key, n.keys[i]) > 0 {
		i++
	}

	// If key is found
	if i < len(n.keys) && t.comparator(key, n.keys[i]) == 0 {
		return n.keys[i], true
	}

	// If it's a leaf node and not found, key doesn't exist
	if n.isLeaf {
		var zero E
		return zero, false
	}

	// Recursively search in the corresponding subtree
	return t.searchNode(n.children[i], key)
}

// Remove deletes an element from the B-tree
//
// Parameters:
//   - key: The element value to delete
//
// Returns:
//   - true if the element existed and was successfully deleted
//   - false if the element didn't exist
//
// Time Complexity: O(log n)
func (t *BTree[E]) Remove(key E) bool {
	if t.root == nil || len(t.root.keys) == 0 {
		return false
	}

	deleted := t.deleteNode(t.root, key)
	if deleted {
		t.size--
		// If root node has no keys but has one child, update root node
		if len(t.root.keys) == 0 && !t.root.isLeaf {
			t.root = t.root.children[0]
		}
	}
	return deleted
}

// Range returns an iterator over all elements in the B-tree within the specified range, in in-order traversal
//
// Parameters:
//   - lowerBound: The lower bound of the range, or nil for no lower bound
//   - upperBound: The upper bound of the range, or nil for no upper bound
//
// Returns:
//   - An iterator over elements within the specified range, in ascending order
//
// Range Rules:
//   - Closed interval lower bound: includes elements equal to lowerBound
//   - Open interval upper bound: excludes elements equal to upperBound
//
// Example:
//
//	// Get all elements greater than or equal to 4 and less than 8
//	lower := 4
//	upper := 8
//	for val := range tree.Range(&lower, &upper) {
//	    fmt.Println(val)
//	}
//
// iter相关方法已移至iter.go文件中

// First returns the smallest element in the B-tree
//
// Returns:
//   - If the tree is non-empty, returns the smallest element and true
//   - If the tree is empty, returns zero value and false
//
// Time Complexity: O(log n)
func (t *BTree[E]) First() (E, bool) {
	var zero E
	if t.root == nil || len(t.root.keys) == 0 {
		return zero, false
	}

	// Start from root node and go all the way left to the bottom
	current := t.root
	for !current.isLeaf {
		current = current.children[0]
	}

	// The first key of the leaf node is the smallest element
	return current.keys[0], true
}

// Last returns the largest element in the B-tree
//
// Returns:
//   - If the tree is non-empty, returns the largest element and true
//   - If the tree is empty, returns zero value and false
//
// Time Complexity: O(log n)
func (t *BTree[E]) Last() (E, bool) {
	var zero E
	if t.root == nil || len(t.root.keys) == 0 {
		return zero, false
	}

	// Start from root node and go all the way right to the bottom
	current := t.root
	for !current.isLeaf {
		current = current.children[len(current.keys)]
	}

	// The last key of the leaf node is the largest element
	return current.keys[len(current.keys)-1], true
}

// PopFirst retrieves and removes the smallest element from the B-tree
//
// Returns:
//   - If the tree is non-empty, returns the removed smallest element and true
//   - If the tree is empty, returns zero value and false
//
// Time Complexity: O(log n)
func (t *BTree[E]) PopFirst() (E, bool) {
	// First get the smallest element
	elem, found := t.First()
	if !found {
		return elem, false
	}

	// Remove that element
	t.Remove(elem)
	return elem, true
}

// PopLast retrieves and removes the largest element from the B-tree
//
// Returns:
//   - If the tree is non-empty, returns the removed largest element and true
//   - If the tree is empty, returns zero value and false
//
// Time Complexity: O(log n)
func (t *BTree[E]) PopLast() (E, bool) {
	// First get the largest element
	elem, found := t.Last()
	if !found {
		return elem, false
	}

	// Remove that element
	t.Remove(elem)
	return elem, true
}

// rangeNode recursively traverses the node and its subtree, returning only elements within the specified range
// This is an internal helper method and not exposed to users
func (t *BTree[E]) rangeNode(n *node[E], lowerBound, upperBound *E, yield func(E) bool) bool {
	for i := 0; i < len(n.keys); i++ {
		// If upper bound exists and current key is >= upper bound, stop traversal
		if upperBound != nil && t.comparator(n.keys[i], *upperBound) >= 0 {
			// For internal nodes, also consider left subtree may have elements < upper bound
			if !n.isLeaf {
				return t.rangeNode(n.children[i], lowerBound, upperBound, yield)
			}
			return true
		}

		// Process left subtree
		if !n.isLeaf {
			if !t.rangeNode(n.children[i], lowerBound, upperBound, yield) {
				return false
			}
		}

		// If lower bound exists and current key < lower bound, skip
		if lowerBound != nil && t.comparator(n.keys[i], *lowerBound) < 0 {
			continue
		}

		// Yield current key
		if !yield(n.keys[i]) {
			return false
		}
	}

	// Process last subtree
	if !n.isLeaf {
		return t.rangeNode(n.children[len(n.keys)], lowerBound, upperBound, yield)
	}
	return true
}

// rangeNodeBack recursively traverses the node and its subtree in reverse order, returning only elements within the specified range
// This is an internal helper method and not exposed to users
func (t *BTree[E]) rangeNodeBack(n *node[E], lowerBound, upperBound *E, yield func(E) bool) bool {
	// Start from the last child and key
	for i := len(n.keys) - 1; i >= 0; i-- {
		// If lower bound exists and current key is < lower bound, stop traversal
		if lowerBound != nil && t.comparator(n.keys[i], *lowerBound) < 0 {
			// For internal nodes, also consider right subtree may have elements > lower bound
			if !n.isLeaf {
				return t.rangeNodeBack(n.children[i+1], lowerBound, upperBound, yield)
			}
			return true
		}

		// Process right subtree
		if !n.isLeaf {
			if !t.rangeNodeBack(n.children[i+1], lowerBound, upperBound, yield) {
				return false
			}
		}

		// If upper bound exists and current key >= upper bound, skip
		if upperBound != nil && t.comparator(n.keys[i], *upperBound) >= 0 {
			continue
		}

		// Yield current key
		if !yield(n.keys[i]) {
			return false
		}
	}

	// Process first subtree
	if !n.isLeaf {
		return t.rangeNodeBack(n.children[0], lowerBound, upperBound, yield)
	}
	return true
}

// deleteNode removes an element from the node and its subtree
// This is an internal helper method and not exposed to users
func (t *BTree[E]) deleteNode(n *node[E], key E) bool {
	order := t.order
	i := 0

	// Find the position of the key in the current node
	for i < len(n.keys) && t.comparator(key, n.keys[i]) > 0 {
		i++
	}

	// Case 1: Key found in current node
	if i < len(n.keys) && t.comparator(key, n.keys[i]) == 0 {
		if n.isLeaf {
			// Case 1a: Leaf node, delete directly
			n.keys = append(n.keys[:i], n.keys[i+1:]...)
			return true
		} else {
			// Case 1b: Internal node
			// Case 1b1: Predecessor has at least order keys, replace current key with predecessor
			if len(n.children[i].keys) >= order {
				predecessor := t.getPredecessor(n, i)
				n.keys[i] = predecessor
				return t.deleteNode(n.children[i], predecessor)
			} else if len(n.children[i+1].keys) >= order {
				// Case 1b2: Successor has at least order keys, replace current key with successor
				successor := t.getSuccessor(n, i)
				n.keys[i] = successor
				return t.deleteNode(n.children[i+1], successor)
			} else {
				// Case 1b3: Both predecessor and successor have only order-1 keys, merge
				t.mergeChildren(n, i)
				return t.deleteNode(n.children[i], key)
			}
		}
	} else {
		// Case 2: Key not found in current node
		if n.isLeaf {
			// Leaf node and not found, key doesn't exist
			return false
		}

		// Determine if borrowing or merging is needed
		needFix := i == len(n.keys) && len(n.children[i].keys) == order-1 ||
			i < len(n.keys) && len(n.children[i].keys) == order-1

		if needFix {
			if i > 0 && len(n.children[i-1].keys) >= order {
				// Case 2a1: Borrow from left sibling
				t.borrowFromLeft(n, i)
			} else if i < len(n.keys) && len(n.children[i+1].keys) >= order {
				// Case 2a2: Borrow from right sibling
				t.borrowFromRight(n, i)
			} else {
				// Case 2a3: Merge children
				if i > 0 {
					t.mergeChildren(n, i-1)
					i--
				} else {
					t.mergeChildren(n, i)
				}
			}
		}

		// Recursive deletion
		return t.deleteNode(n.children[i], key)
	}
}

// getPredecessor gets the predecessor of a key
// This is an internal helper method and not exposed to users
func (t *BTree[E]) getPredecessor(n *node[E], i int) E {
	node := n.children[i]
	// Go all the way to the right
	for !node.isLeaf {
		node = node.children[len(node.keys)]
	}
	return node.keys[len(node.keys)-1]
}

// getSuccessor gets the successor of a key
// This is an internal helper method and not exposed to users
func (t *BTree[E]) getSuccessor(n *node[E], i int) E {
	node := n.children[i+1]
	// Go all the way to the left
	for !node.isLeaf {
		node = node.children[0]
	}
	return node.keys[0]
}

// borrowFromLeft borrows a key from the left sibling
// This is an internal helper method and not exposed to users
func (t *BTree[E]) borrowFromLeft(n *node[E], i int) {
	child := n.children[i]
	sibling := n.children[i-1]

	// Parent node's key moves down to child node
	child.keys = append([]E{n.keys[i-1]}, child.keys...)

	// Sibling node's maximum key moves up to parent node
	n.keys[i-1] = sibling.keys[len(sibling.keys)-1]
	sibling.keys = sibling.keys[:len(sibling.keys)-1]

	// If not a leaf node, also move child node pointers
	if !sibling.isLeaf {
		child.children = append([]*node[E]{sibling.children[len(sibling.children)-1]}, child.children...)
		sibling.children = sibling.children[:len(sibling.children)-1]
	}
}

// borrowFromRight borrows a key from the right sibling
// This is an internal helper method and not exposed to users
func (t *BTree[E]) borrowFromRight(n *node[E], i int) {
	child := n.children[i]
	sibling := n.children[i+1]

	// Parent node's key moves down to child node
	child.keys = append(child.keys, n.keys[i])

	// Sibling node's minimum key moves up to parent node
	n.keys[i] = sibling.keys[0]
	sibling.keys = sibling.keys[1:]

	// If not a leaf node, also move child node pointers
	if !sibling.isLeaf {
		child.children = append(child.children, sibling.children[0])
		sibling.children = sibling.children[1:]
	}
}

// mergeChildren merges two child nodes
// This is an internal helper method and not exposed to users
func (t *BTree[E]) mergeChildren(n *node[E], i int) {
	child := n.children[i]
	sibling := n.children[i+1]

	// Parent node's key moves down to child node
	child.keys = append(child.keys, n.keys[i])

	// Merge sibling node's keys
	child.keys = append(child.keys, sibling.keys...)

	// If not a leaf node, merge child node pointers
	if !sibling.isLeaf {
		child.children = append(child.children, sibling.children...)
	}

	// Update parent node's keys and children list
	n.keys = append(n.keys[:i], n.keys[i+1:]...)
	n.children = append(n.children[:i+1], n.children[i+2:]...)
}
