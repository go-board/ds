// Package btreeset implements a B-tree based ordered set data structure.
//
// BTreeSet is a collection type that stores unique elements in a B-tree structure,
// providing ordered iteration and efficient set operations. All operations run in O(log n) time.
//
// # Time Complexity
//
//   - Insert: O(log n)
//   - Delete: O(log n)
//   - Contains: O(log n)
//   - Union/Intersection/Difference: O(m + n) where m, n are set sizes
//   - Traversal: O(n)
//
// # Features
//
//   - Ordered storage with sorted element iteration (IterAsc, IterDesc)
//   - Efficient range queries with configurable boundaries (RangeAsc, RangeDesc)
//   - Standard set operations: Union, Intersection, Difference, SymmetricDifference
//   - Subset/superset checks: IsSubsetOf, IsSupersetOf
//   - Generic type support for any cmp.Ordered element type
//
// # Usage
//
//	BTreeSet supports both custom comparators and ordered types:
//
//		// For cmp.Ordered elements (recommended)
//		set := btreeset.NewOrdered[int]()
//		set.Insert(3)
//		set.Insert(1)
//		set.Insert(4)
//
//		// For custom element types, provide a comparator
//		type Person struct { Name string; Age int }
//		set := btreeset.New[Person](func(a, b Person) int {
//		    return strings.Compare(a.Name, b.Name)
//		})
//
// # Set Operations
//
//	BTreeSet provides rich set operations:
//
//		set1 := btreeset.NewOrdered[int]()
//		set1.Insert(1, 2, 3)
//
//		set2 := btreeset.NewOrdered[int]()
//		set2.Insert(2, 3, 4)
//
//		union := set1.Union(set2)           // {1, 2, 3, 4}
//		intersection := set1.Intersection(set2) // {2, 3}
//		difference := set1.Difference(set2)    // {1}
//		symDiff := set1.SymmetricDifference(set2) // {1, 4}
package btreeset
