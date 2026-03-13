// Package bound provides types for representing range boundaries.
//
// This package is useful for implementing range queries and slices with
// configurable boundary inclusion/exclusion. It's commonly used with
// ordered data structures like B-trees, skip lists, and tries.
//
// # Overview
//
// The package provides two main types:
//
//   - [Bound]: Represents one side of a range (start or end)
//   - [RangeBounds]: Combines two bounds to describe a full range
//
// # Bound Types
//
// A [Bound] represents one side of a range with three possible kinds:
//
//   - [Unbounded]: No boundary on this side (represented as +∞ or -∞)
//   - [Included]: The boundary value is included in the range [a, b]
//   - [Excluded]: The boundary value is excluded from the range (a, b)
//
// # RangeBounds
//
// A [RangeBounds] combines a start and end boundary to describe a full range.
// Use [NewRangeBounds] to create a range, and [RangeBounds.Contains] to check
// if a value falls within the range.
//
// # Usage with Ordered Collections
//
// Bound types are used with ordered collections for range queries:
//
//	import (
//		"cmp"
//
//		"github.com/go-board/ds/bound"
//		"github.com/go-board/ds/btreemap"
//	)
//
//	m := btreemap.NewOrdered[string, int]()
//	m.Insert("apple", 1)
//	m.Insert("banana", 2)
//	m.Insert("cherry", 3)
//
//	// Range query: [banana, cherry)
//	lower := bound.NewIncluded("banana")
//	upper := bound.NewExcluded("cherry")
//	for k, v := range m.RangeAsc(bound.NewRangeBounds(lower, upper)) {
//	    fmt.Println(k, v) // banana 2
//	}
//
// # Examples
//
//	// Closed range: [1, 10] includes 1 and 10
//	bounds := bound.NewRangeBounds(
//	    bound.NewIncluded(1),
//	    bound.NewIncluded(10),
//	)
//
//	// Open range: (1, 10) excludes 1 and 10
//	bounds := bound.NewRangeBounds(
//	    bound.NewExcluded(1),
//	    bound.NewExcluded(10),
//	)
//
//	// Half-open range: [1, 10) includes 1, excludes 10
//	bounds := bound.NewRangeBounds(
//	    bound.NewIncluded(1),
//	    bound.NewExcluded(10),
//	)
//
//	// Unbounded: [1, +∞)
//	bounds := bound.NewRangeBounds(
//	    bound.NewIncluded(1),
//	    bound.NewUnbounded[int](),
//	)
//
//	// Check containment
//	bounds.Contains(cmp.Compare, 5)  // true
//	bounds.Contains(cmp.Compare, 1)  // true (included)
//	bounds.Contains(cmp.Compare, 10) // false (excluded)
package bound
