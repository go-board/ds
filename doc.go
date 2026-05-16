// Package ds provides a single import path that re-exports all data structures
// in this module.
//
// # Overview
//
// Sequences:
//   - [ArrayDeque]
//   - [ArrayStack]
//   - [LinkedList]
//
// Ordered maps:
//   - [ArrayMap] (sorted slice)
//   - [BTreeMap]
//   - [SkipMap]
//
// Unordered maps:
//   - [HashMap]
//
// Ordered sets:
//   - [ArraySet] (sorted slice)
//   - [BTreeSet]
//   - [SkipSet]
//
// Unordered sets:
//   - [HashSet]
//
// Other structures:
//   - [BTree]
//   - [PriorityQueue]
//   - [TrieMap]
//
// Utilities:
//   - [Bound], [RangeBounds]
//   - [Hasher], [DefaultHasher], [SliceHasher], [MapHasher]
//
// # Quick Start
//
//	import "github.com/go-board/ds"
//
//	func demo() {
//		m := ds.NewOrderedArrayMap[string, int]()
//		m.Insert("a", 1)
//
//		s := ds.NewOrderedArraySet[int]()
//		s.Insert(10)
//
//		hm := ds.NewComparableHashMap[string, int]()
//		hm.Insert("k", 1)
//
//		pq := ds.NewMinOrderedPriorityQueue[int]()
//		pq.Push(3)
//	}
package ds
