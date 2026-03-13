// Package skipset implements a sorted set based on skip lists.
//
// SkipSet is a collection type that stores unique elements in a skip list structure,
// providing ordered iteration and efficient set operations. All operations run in O(log n) average time.
// Internally, SkipSet is implemented using SkipMap with empty structs as values.
//
// # Time Complexity
//
//   - Insert: O(log n) average
//   - Delete: O(log n) average
//   - Contains: O(log n) average
//   - Union/Intersection/Difference: O(m log(m+n)) average where m, n are set sizes
//   - Traversal: O(n)
//
// # Features
//
//   - Ordered storage with sorted element iteration (IterAsc, IterDesc)
//   - Efficient range queries with configurable boundaries (RangeAsc, RangeDesc)
//   - Standard set operations: Union, Intersection, Difference, SymmetricDifference
//   - Subset/superset checks: IsSubsetOf, IsSupersetOf
//   - Generic type support for any cmp.Ordered element type
//   - Simpler implementation compared to BTreeSet
//
// # Usage
//
// SkipSet is simpler than BTreeSet and often faster for smaller datasets:
//
//	// For cmp.Ordered elements (recommended)
//	s := skipset.NewOrdered[string]()
//	s.Insert("apple")
//	s.Insert("banana")
//
//	// Check existence
//	if s.Contains("apple") {
//	    fmt.Println("Found apple")
//	}
//
// # Set Operations
//
// SkipSet provides rich set operations:
//
//	set1 := skipset.NewOrdered[int]()
//	set1.Insert(1, 2, 3)
//
//	set2 := skipset.NewOrdered[int]()
//	set2.Insert(2, 3, 4)
//
//	union := set1.Union(set2)           // {1, 2, 3, 4}
//	intersection := set1.Intersection(set2) // {2, 3}
//	difference := set1.Difference(set2)    // {1}
package skipset
