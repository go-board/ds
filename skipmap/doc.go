// Package skipmap implements an ordered key-value map based on skip lists.
//
// SkipMap is a sorted associative array that stores key-value pairs in a skip list structure,
// providing ordered iteration and efficient range queries. All operations run in O(log n) average time.
//
// # Time Complexity
//
//   - Insert: O(log n) average
//   - Delete: O(log n) average
//   - Get: O(log n) average
//   - Range query: O(log n + k) average where k is the number of elements in range
//   - Traversal: O(n)
//
// # Features
//
//   - Ordered key storage with sorted iteration (IterAsc, IterDesc)
//   - Efficient range queries with configurable boundaries (RangeAsc, RangeDesc)
//   - Duplicate key detection (Insert returns old value and updated flag)
//   - Batch operations (InsertBatch, DeleteBatch)
//   - Generic type support for any cmp.Ordered key type
//   - Simpler implementation compared to B-trees
//
// # Usage
//
// SkipMap is simpler than BTreeMap and often faster for smaller datasets:
//
//	// For cmp.Ordered keys (recommended)
//	m := skipmap.NewOrdered[string, int]()
//	m.Insert("apple", 5)
//	m.Insert("banana", 3)
//
//	// Get or insert atomically
//	val, existed := m.GetOrInsert("cherry", 10)
//	if existed {
//	    fmt.Println("Key already existed:", val)
//	}
//
// # Iterators
//
// SkipMap provides multiple iterator types:
//
//	// All key-value pairs in ascending key order
//	for k, v := range m.IterAsc() {
//	    fmt.Printf("%s: %d\n", k, v)
//	}
//
//	// Range query
//	lower := bound.NewIncluded("apple")
//	upper := bound.NewExcluded("cherry")
//	for k, v := range m.RangeAsc(bound.NewRangeBounds(lower, upper)) {
//	    fmt.Printf("%s: %d\n", k, v)
//	}
package skipmap
