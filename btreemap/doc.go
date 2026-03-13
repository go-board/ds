// Package btreemap implements an ordered key-value map based on B-trees.
//
// BTreeMap is a sorted associative array that stores key-value pairs in a B-tree structure,
// providing ordered iteration and efficient range queries. All operations run in O(log n) time.
//
// # Time Complexity
//
//   - Insert: O(log n)
//   - Delete: O(log n)
//   - Get: O(log n)
//   - Range query: O(log n + k) where k is the number of elements in range
//   - Traversal: O(n)
//
// # Features
//
//   - Ordered key storage with sorted iteration (IterAsc, IterDesc)
//   - Efficient range queries with configurable boundaries (RangeAsc, RangeDesc)
//   - Duplicate key detection (Insert returns old value and updated flag)
//   - Batch operations (InsertBatch, DeleteBatch)
//   - Generic type support for any cmp.Ordered key type
//
// # Usage
//
// BTreeMap supports both custom comparators and ordered types:
//
//	// For cmp.Ordered keys (recommended)
//	m := btreemap.NewOrdered[string, int]()
//	m.Insert("apple", 5)
//	m.Insert("banana", 3)
//
//	// For custom key types, provide a comparator
//	type Point struct { x, y int }
//	m := btreemap.New[Point, string](func(a, b Point) int {
//	    if a.x != b.x {
//	        return a.x - b.x
//	    }
//	    return a.y - b.y
//	})
//
// # Iterators
//
// BTreeMap provides multiple iterator types:
//
//	// All key-value pairs in ascending key order
//	for k, v := range m.IterAsc() {
//	    fmt.Printf("%s: %d\n", k, v)
//	}
//
//	// Range query
//	lower := bound.NewIncluded("apple")
//	upper := bound.NewExcluded("banana")
//	for k, v := range m.RangeAsc(bound.NewRangeBounds(lower, upper)) {
//	    fmt.Printf("%s: %d\n", k, v)
//	}
package btreemap
