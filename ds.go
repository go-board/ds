package ds

import (
	"cmp"

	"github.com/go-board/ds/arraydeque"
	"github.com/go-board/ds/arraystack"
	"github.com/go-board/ds/bound"
	"github.com/go-board/ds/btree"
	"github.com/go-board/ds/btreemap"
	"github.com/go-board/ds/btreeset"
	"github.com/go-board/ds/hashmap"
	"github.com/go-board/ds/hashset"
	"github.com/go-board/ds/hashutil"
	"github.com/go-board/ds/linkedlist"
	"github.com/go-board/ds/priorityqueue"
	"github.com/go-board/ds/skipmap"
	"github.com/go-board/ds/skipset"
	"github.com/go-board/ds/trie"
)

// Double-ended Queue (ArrayDeque)

// ArrayDeque is an alias for [arraydeque.ArrayDeque],
// providing slice-based double-ended queue functionality.
type ArrayDeque[T any] = arraydeque.ArrayDeque[T]

// NewArrayDeque creates a new slice-based double-ended queue instance.
func NewArrayDeque[T any]() *ArrayDeque[T] {
	return arraydeque.New[T]()
}

// Stack (ArrayStack)

// ArrayStack is an alias for [arraystack.ArrayStack],
// providing slice-based stack functionality.
type ArrayStack[T any] = arraystack.ArrayStack[T]

// NewArrayStack creates a new slice-based stack instance.
func NewArrayStack[T any]() *ArrayStack[T] {
	return arraystack.New[T]()
}

// B-Tree

// Bound Kind

// BoundKind is an alias for [bound.Kind].
type BoundKind = bound.Kind

const (
	// Unbounded means no boundary on this side.
	Unbounded BoundKind = bound.Unbounded
	// Included means boundary value is included.
	Included BoundKind = bound.Included
	// Excluded means boundary value is excluded.
	Excluded BoundKind = bound.Excluded
)

// Bound is an alias for [bound.Bound].
type Bound[T any] = bound.Bound[T]

// RangeBounds is an alias for [bound.RangeBounds].
type RangeBounds[T any] = bound.RangeBounds[T]

// NewUnbounded creates an unbounded boundary.
func NewUnbounded[T any]() Bound[T] { return bound.NewUnbounded[T]() }

// NewIncluded creates an inclusive boundary.
func NewIncluded[T any](value T) Bound[T] { return bound.NewIncluded(value) }

// NewExcluded creates an exclusive boundary.
func NewExcluded[T any](value T) Bound[T] { return bound.NewExcluded(value) }

// NewRangeBounds creates a start/end range bounds object.
func NewRangeBounds[T any](start, end Bound[T]) RangeBounds[T] {
	return bound.NewRangeBounds(start, end)
}

// B-Tree

// BTree is an alias for [btree.BTree],
// providing ordered data storage functionality.
type BTree[T any] = btree.BTree[T]

// NewBTree creates a new B-tree instance.
// comparator is used to compare element sizes and cannot be nil.
func NewBTree[T any](comparator func(T, T) int) *BTree[T] {
	return btree.New(comparator)
}

// NewOrderedBTree creates a new ordered B-tree instance.
// The element type must implement the [cmp.Ordered] interface.
func NewOrderedBTree[T cmp.Ordered]() *BTree[T] {
	return btree.NewOrdered[T]()
}

// B-Tree Map

// BTreeMap is a type alias for [btreemap.BTreeMap],
// providing ordered key-value mapping functionality.
type BTreeMap[K any, V any] = btreemap.BTreeMap[K, V]

// NewBTreeMap creates a new B-tree map instance.
// comparator is used to compare key sizes, returning negative,
// zero, or positive to indicate whether the first key is less than,
// equal to, or greater than the second key.
func NewBTreeMap[K any, V any](comparator func(K, K) int) *BTreeMap[K, V] {
	return btreemap.New[K, V](comparator)
}

// NewOrderedBTreeMap creates a new ordered B-tree map instance.
// The key type must implement the [cmp.Ordered] interface.
func NewOrderedBTreeMap[K cmp.Ordered, V any]() *BTreeMap[K, V] {
	return btreemap.NewOrdered[K, V]()
}

// B-Tree Set

// BTreeSet is a type alias for [btreeset.BTreeSet],
// providing ordered set functionality.
type BTreeSet[T any] = btreeset.BTreeSet[T]

// NewBTreeSet creates a new B-tree set instance.
// comparator is used to compare element sizes, returning negative,
// zero, or positive to indicate whether the first element is less than,
// equal to, or greater than the second element.
func NewBTreeSet[T any](comparator func(T, T) int) *BTreeSet[T] {
	return btreeset.New(comparator)
}

// NewOrderedBTreeSet creates a new ordered B-tree set instance.
// The element type must implement the [cmp.Ordered] interface.
func NewOrderedBTreeSet[T cmp.Ordered]() *BTreeSet[T] {
	return btreeset.NewOrdered[T]()
}

// Hash Map

// HashMap is an alias for [hashmap.HashMap],
// providing key-value mapping functionality.
type HashMap[K any, V any, H hashutil.Hasher[K]] = hashmap.HashMap[K, V, H]

// NewHashMap creates a new hash map instance.
// h is the hasher used to calculate key hash values.
func NewHashMap[K any, V any, H hashutil.Hasher[K]](h H) *HashMap[K, V, H] {
	return hashmap.New[K, V](h)
}

// NewComparableHashMap creates a new hash map instance with a default hasher for comparable key types.
func NewComparableHashMap[K comparable, V any]() *HashMap[K, V, hashutil.Default[K]] {
	return hashmap.NewComparable[K, V]()
}

// NewHashMapFromMap creates a new hash map instance from an existing map.
func NewHashMapFromMap[K comparable, V any, M ~map[K]V](m M) *HashMap[K, V, hashutil.Default[K]] {
	return hashmap.NewFromMap(m)
}

// Hash Set

// HashSet is an alias for [hashset.HashSet], providing set functionality.
type HashSet[T any, H hashutil.Hasher[T]] = hashset.HashSet[T, H]

// NewHashSet creates a new hash set instance.
// h is the hasher used to calculate element hash values.
func NewHashSet[T any, H hashutil.Hasher[T]](h H) *HashSet[T, H] {
	return hashset.New(h)
}

// NewComparableHashSet creates a new hash set instance with a default hasher for comparable key types.
func NewComparableHashSet[T comparable]() *HashSet[T, hashutil.Default[T]] {
	return hashset.NewComparable[T]()
}

// Hash Utilities

// Hasher is an alias for [hashutil.Hasher],
// providing a generic interface for comparing and hashing values of any type.
type Hasher[E any] = hashutil.Hasher[E]

// DefaultHasher is a default hasher implementation for comparable types.
// It uses Go's standard library maphash package for hash calculations.
type DefaultHasher[E comparable] = hashutil.Default[E]

// SliceHasher is a hasher implementation for slice types.
// It uses the element type's Hasher to calculate the hash for the entire slice.
type SliceHasher[E ~[]T, T any, H Hasher[T]] = hashutil.SliceHasher[E, T, H]

// MapHasher is a hasher implementation for map types.
// It uses the key type's Hasher and value type's Hasher to calculate the hash for the entire map.
type MapHasher[E ~map[K]V, K comparable, V any, H Hasher[V]] = hashutil.MapHasher[E, K, V, H]

// Linked List

// LinkedList is an alias for [linkedlist.LinkedList],
// providing linked list-based functionality.
type LinkedList[T any] = linkedlist.LinkedList[T]

// NewLinkedList creates a new linked list instance.
func NewLinkedList[T any]() *LinkedList[T] {
	return linkedlist.New[T]()
}

// Priority Queue

// PriorityQueue is an alias for [priorityqueue.PriorityQueue].
type PriorityQueue[T any] = priorityqueue.PriorityQueue[T]

// NewMinPriorityQueue creates a new min-heap priority queue.
// cmp is a comparison function: returns negative when a < b,
// 0 when a == b, positive when a > b.
// The largest element has the highest priority in a max-heap.
func NewMinPriorityQueue[T any](cmp func(T, T) int) *PriorityQueue[T] {
	return priorityqueue.NewMin(cmp)
}

// NewMaxPriorityQueue creates a new max-heap priority queue.
// cmp is a comparison function: returns negative when a < b,
// 0 when a == b, positive when a > b.
// The largest element has the highest priority in a max-heap.
func NewMaxPriorityQueue[T any](cmp func(T, T) int) *PriorityQueue[T] {
	return priorityqueue.NewMax(cmp)
}

// NewMinOrderedPriorityQueue creates a new min-heap priority queue,
// where the element type must implement the [cmp.Ordered] interface.
func NewMinOrderedPriorityQueue[T cmp.Ordered]() *PriorityQueue[T] {
	return priorityqueue.NewMinOrdered[T]()
}

// NewMaxOrderedPriorityQueue creates a new max-heap priority queue,
// where the element type must implement the [cmp.Ordered] interface.
func NewMaxOrderedPriorityQueue[T cmp.Ordered]() *PriorityQueue[T] {
	return priorityqueue.NewMaxOrdered[T]()
}

// Skip List Map

// SkipMap is a type alias for [skipmap.SkipMap],
// providing skip list-based key-value mapping functionality.
type SkipMap[K any, V any] = skipmap.SkipMap[K, V]

// NewSkipMap creates a new skip list map instance.
// comparator is used to compare key sizes, returning negative,
// zero, or positive to indicate whether the first key is less than,
// equal to, or greater than the second key.
func NewSkipMap[K any, V any](comparator func(K, K) int) *SkipMap[K, V] {
	return skipmap.New[K, V](comparator)
}

// NewOrderedSkipMap creates a new ordered skip list map instance.
// The key type must implement the [cmp.Ordered] interface.
func NewOrderedSkipMap[K cmp.Ordered, V any]() *SkipMap[K, V] {
	return skipmap.NewOrdered[K, V]()
}

// Skip List Set

// SkipSet is a type alias for [skipset.SkipSet],
// providing skip list-based ordered set functionality.
type SkipSet[E any] = skipset.SkipSet[E]

// NewSkipSet creates a new skip list set instance.
// comparator is used to compare element sizes, returning negative,
// zero, or positive to indicate whether the first element is less than,
// equal to, or greater than the second element.
func NewSkipSet[E any](comparator func(E, E) int) *SkipSet[E] {
	return skipset.New(comparator)
}

// NewOrderedSkipSet creates a new ordered skip list set instance.
// The element type must implement the [cmp.Ordered] interface.
func NewOrderedSkipSet[E cmp.Ordered]() *SkipSet[E] {
	return skipset.NewOrdered[E]()
}

// Trie Map

// TrieMap is an alias for [trie.TrieMap],
// providing trie-based key-value mapping functionality.
type TrieMap[K any, V any] = trie.TrieMap[K, V]

// NewTrieMap creates a new trie map instance.
// comparator is used to compare individual key elements.
func NewTrieMap[K any, V any](comparator func(K, K) int) *TrieMap[K, V] {
	return trie.New[K, V](comparator)
}

// NewOrderedTrieMap creates a new trie map instance for ordered key types.
// The key type must implement the [cmp.Ordered] interface.
func NewOrderedTrieMap[K cmp.Ordered, V any]() *TrieMap[K, V] {
	return trie.NewOrdered[K, V]()
}
