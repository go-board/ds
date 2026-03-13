// Package ds provides a unified import path for all data structures in this module.
//
// This package re-exports all data structures behind a single, consistent import path,
// making it easy to use the library without memorizing multiple import paths.
//
// # Data Structures Overview
//
// The module provides the following categories of data structures:
//
// # Sequences
//
//   - [ArrayDeque]: Double-ended queue (deque) with O(1) amortized operations
//   - [ArrayStack]: Stack (LIFO) with O(1) push/pop
//   - [LinkedList]: Doubly linked list with O(1) insertion/deletion at ends
//
// # Sets
//
// Ordered sets (sorted iteration, range queries):
//   - [BTreeSet]: B-tree based, O(log n) operations
//   - [SkipSet]: Skip list based, O(log n) average operations
//
// Unordered sets (faster average operations):
//   - [HashSet]: Hash table based, O(1) average operations
//
// # Maps
//
// Ordered maps (sorted iteration, range queries):
//   - [BTreeMap]: B-tree based, O(log n) operations
//   - [SkipMap]: Skip list based, O(log n) average operations
//
// Unordered maps (faster average operations):
//   - [HashMap]: Hash table based, O(1) average operations
//
// # Specialized
//
//   - [PriorityQueue]: Heap-based min/max priority queue
//   - [TrieMap]: Prefix tree for string/sequence operations
//
// # Quick Start
//
//	import "github.com/go-board/ds"
//
//	// Sequences
//	queue := ds.NewArrayDeque[int]()
//	stack := ds.NewArrayStack[int]()
//	list := ds.NewLinkedList[int]()
//
//	// Ordered collections (sorted iteration, range queries)
//	btree := ds.NewOrderedBTree[int]()
//	btreemap := ds.NewOrderedBTreeMap[string, int]()
//	btreeset := ds.NewOrderedBTreeSet[int]()
//	skipmap := ds.NewOrderedSkipMap[string, int]()
//	skipset := ds.NewOrderedSkipSet[string]()
//
//	// Unordered collections (faster average ops)
//	hashmap := ds.NewHashMap[string, int](hashutil.Default[string]{})
//	hashset := ds.NewHashSet[string](hashutil.Default[string]{})
//
//	// Specialized
//	pq := ds.NewOrderedMinPriorityQueue[int]()
//	trie := ds.NewOrderedTrieMap[[]byte, int]()
//
// # Choosing a Data Structure
//
//   - Need fast lookups? → HashMap, HashSet
//   - Need sorted iteration? → BTreeMap, BTreeSet, SkipMap, SkipSet
//   - Need range queries? → BTreeMap, BTreeSet, SkipMap, SkipSet
//   - Need priority order? → PriorityQueue
//   - Need prefix matching? → TrieMap
//   - Need FIFO/LIFO? → ArrayDeque, ArrayStack, LinkedList
package ds
