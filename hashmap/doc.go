// Package hashmap implements a generic hash map data structure.
//
// HashMap provides efficient key-value mapping operations using a hash table.
// It achieves average O(1) time complexity for insertion, deletion, and lookup.
// Unlike ordered maps (BTreeMap, SkipMap), iteration order is not guaranteed.
//
// # Time Complexity
//
//   - Insert: O(1) average, O(n) worst case (hash collisions)
//   - Delete: O(1) average, O(n) worst case
//   - Get: O(1) average, O(n) worst case
//   - Traversal: O(n)
//
// # Features
//
//   - Fast key-value lookup with average O(1) time complexity
//   - Generic type support with custom hashers
//   - Batch operations (InsertBatch, DeleteBatch)
//   - Load factor management for performance tuning
//   - Iterator support for traversal
//
// # Usage
//
// HashMap requires a hasher implementation for key types:
//
//	// For string keys, use the provided StringHasher
//	m := hashmap.New[string, int](hashutil.StringHasher)
//	m.Insert("apple", 5)
//	m.Insert("banana", 3)
//
//	// For custom types, implement the Hasher interface
//	type Point struct { X, Y int }
//	pointHasher := hashutil.NewStructHasher(Point{})
//	m := hashmap.New[Point, string](pointHasher)
//
// # Differences from Ordered Maps
//
// HashMap provides faster average-case performance but lacks ordering:
//
//	// BTreeMap/SkipMap: sorted iteration, range queries
//	// HashMap: faster operations, no ordering guarantee
//
// When to use HashMap:
//   - Fast lookups are the primary concern
//   - Iteration order doesn't matter
//   - Keys are not naturally ordered
package hashmap
