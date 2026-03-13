// Package hashset implements an unordered set data structure based on hash tables.
//
// HashSet is a collection type that stores unique elements in a hash table.
// It achieves average O(1) time complexity for insertion, deletion, and lookup.
// Unlike ordered sets (BTreeSet, SkipSet), iteration order is not guaranteed.
// Internally, HashSet is implemented using HashMap with empty structs as values.
//
// # Time Complexity
//
//   - Insert: O(1) average, O(n) worst case (hash collisions)
//   - Delete: O(1) average, O(n) worst case
//   - Contains: O(1) average, O(n) worst case
//   - Union/Intersection/Difference: O(m + n) where m, n are set sizes
//   - Traversal: O(n)
//
// # Features
//
//   - Fast membership testing with average O(1) time complexity
//   - Generic type support with custom hashers
//   - Standard set operations: Union, Intersection, Difference, SymmetricDifference
//   - Subset/superset checks: IsSubsetOf, IsSupersetOf
//   - Batch operations
//
// # Usage
//
// HashSet requires a hasher implementation for element types:
//
//	// For string elements, use the provided StringHasher
//	set := hashset.New[string](hashutil.StringHasher)
//	set.Insert("apple")
//	set.Insert("banana")
//
//	// Check existence
//	if set.Contains("apple") {
//	    fmt.Println("Found apple")
//	}
//
// # Differences from Ordered Sets
//
// HashSet provides faster average-case performance but lacks ordering:
//
//	// BTreeSet/SkipSet: sorted iteration, range queries
//	// HashSet: faster operations, no ordering guarantee
//
// When to use HashSet:
//   - Fast membership testing is the primary concern
//   - Iteration order doesn't matter
//   - Elements are not naturally ordered
package hashset
