// Package btree implements a generic B-tree data structure.
//
// A B-tree is a self-balancing tree data structure that maintains sorted data and allows
// efficient insertion, deletion, and lookup operations in O(log n) time complexity.
// It's particularly suitable for external storage systems like databases and file systems
// because it optimizes for block-oriented storage and minimizes disk I/O operations.
//
// # Time Complexity
//
//   - Insert: O(log n)
//   - Delete: O(log n)
//   - Search: O(log n)
//   - Range query: O(log n + k) where k is the number of elements in range
//   - Traversal: O(n)
//
// # Features
//
//   - Ordered storage with sorted element traversal (IterAsc, IterDesc)
//   - Efficient range queries with configurable boundaries (RangeAsc, RangeDesc)
//   - Support for finding first and last elements (First, Last)
//   - Configurable branching factor for tuning performance
//   - Generic type support for any comparable element type
//
// # Usage
//
// BTree supports both custom comparators and ordered types:
//
//	// For cmp.Ordered types, use NewOrdered (recommended)
//	tree := btree.NewOrdered[int]()
//	tree.Insert(3)
//	tree.Insert(1)
//	tree.Insert(4)
//
//	// For custom types, provide a comparator
//	tree := btree.New[Person](func(a, b Person) int {
//	    return cmp.Compare(a.Age, b.Age)
//	})
//
// # Iterators
//
// BTree provides multiple iterator types for different traversal needs:
//
//	// Ascending order
//	for val := range tree.IterAsc() {
//	    fmt.Println(val)
//	}
//
//	// Descending order
//	for val := range tree.IterDesc() {
//	    fmt.Println(val)
//	}
//
//	// Range query (ascending)
//	lower := bound.NewIncluded(3)
//	upper := bound.NewExcluded(10)
//	for val := range tree.RangeAsc(bound.NewRangeBounds(lower, upper)) {
//	    fmt.Println(val)
//	}
package btree
