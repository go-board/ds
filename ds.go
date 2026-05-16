package ds

import (
	"cmp"

	"github.com/go-board/ds/arraydeque"
	"github.com/go-board/ds/arraymap"
	"github.com/go-board/ds/arrayset"
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

// ===== Bound utilities =====

type BoundKind = bound.Kind

const (
	Unbounded BoundKind = bound.Unbounded
	Included  BoundKind = bound.Included
	Excluded  BoundKind = bound.Excluded
)

type Bound[T any] = bound.Bound[T]
type RangeBounds[T any] = bound.RangeBounds[T]

func NewUnbounded[T any]() Bound[T]       { return bound.NewUnbounded[T]() }
func NewIncluded[T any](value T) Bound[T] { return bound.NewIncluded(value) }
func NewExcluded[T any](value T) Bound[T] { return bound.NewExcluded(value) }
func NewRangeBounds[T any](start, end Bound[T]) RangeBounds[T] {
	return bound.NewRangeBounds(start, end)
}

// ===== Sequence types =====

type ArrayDeque[T any] = arraydeque.ArrayDeque[T]

func NewArrayDeque[T any]() *ArrayDeque[T] { return arraydeque.New[T]() }

type ArrayStack[T any] = arraystack.ArrayStack[T]

func NewArrayStack[T any]() *ArrayStack[T] { return arraystack.New[T]() }

type LinkedList[T any] = linkedlist.LinkedList[T]

func NewLinkedList[T any]() *LinkedList[T] { return linkedlist.New[T]() }

// ===== Map types =====

type ArrayMap[K any, V any] = arraymap.ArrayMap[K, V]

func NewArrayMap[K any, V any](comparator func(K, K) int) *ArrayMap[K, V] {
	return arraymap.New[K, V](comparator)
}
func NewOrderedArrayMap[K cmp.Ordered, V any]() *ArrayMap[K, V] {
	return arraymap.NewOrdered[K, V]()
}
func NewArrayMapFromMap[K cmp.Ordered, V any, M ~map[K]V](m M) *ArrayMap[K, V] {
	return arraymap.NewFromMap[K, V](m)
}

type HashMap[K any, V any, H hashutil.Hasher[K]] = hashmap.HashMap[K, V, H]

func NewHashMap[K any, V any, H hashutil.Hasher[K]](h H) *HashMap[K, V, H] {
	return hashmap.New[K, V](h)
}
func NewComparableHashMap[K comparable, V any]() *HashMap[K, V, hashutil.Default[K]] {
	return hashmap.NewComparable[K, V]()
}
func NewHashMapFromMap[K comparable, V any, M ~map[K]V](m M) *HashMap[K, V, hashutil.Default[K]] {
	return hashmap.NewFromMap(m)
}

type BTreeMap[K any, V any] = btreemap.BTreeMap[K, V]

func NewBTreeMap[K any, V any](comparator func(K, K) int) *BTreeMap[K, V] {
	return btreemap.New[K, V](comparator)
}
func NewOrderedBTreeMap[K cmp.Ordered, V any]() *BTreeMap[K, V] {
	return btreemap.NewOrdered[K, V]()
}

type SkipMap[K any, V any] = skipmap.SkipMap[K, V]

func NewSkipMap[K any, V any](comparator func(K, K) int) *SkipMap[K, V] {
	return skipmap.New[K, V](comparator)
}
func NewOrderedSkipMap[K cmp.Ordered, V any]() *SkipMap[K, V] {
	return skipmap.NewOrdered[K, V]()
}

type TrieMap[K any, V any] = trie.TrieMap[K, V]

func NewTrieMap[K any, V any](comparator func(K, K) int) *TrieMap[K, V] {
	return trie.New[K, V](comparator)
}
func NewOrderedTrieMap[K cmp.Ordered, V any]() *TrieMap[K, V] {
	return trie.NewOrdered[K, V]()
}

// ===== Set types =====

type ArraySet[E any] = arrayset.ArraySet[E]

func NewArraySet[E any](comparator func(E, E) int) *ArraySet[E] { return arrayset.New[E](comparator) }
func NewOrderedArraySet[E cmp.Ordered]() *ArraySet[E]           { return arrayset.NewOrdered[E]() }

type HashSet[T any, H hashutil.Hasher[T]] = hashset.HashSet[T, H]

func NewHashSet[T any, H hashutil.Hasher[T]](h H) *HashSet[T, H] { return hashset.New(h) }
func NewComparableHashSet[T comparable]() *HashSet[T, hashutil.Default[T]] {
	return hashset.NewComparable[T]()
}

type BTreeSet[T any] = btreeset.BTreeSet[T]

func NewBTreeSet[T any](comparator func(T, T) int) *BTreeSet[T] { return btreeset.New(comparator) }
func NewOrderedBTreeSet[T cmp.Ordered]() *BTreeSet[T]           { return btreeset.NewOrdered[T]() }

type SkipSet[E any] = skipset.SkipSet[E]

func NewSkipSet[E any](comparator func(E, E) int) *SkipSet[E] { return skipset.New(comparator) }
func NewOrderedSkipSet[E cmp.Ordered]() *SkipSet[E]           { return skipset.NewOrdered[E]() }

// ===== Tree / heap types =====

type BTree[T any] = btree.BTree[T]

func NewBTree[T any](comparator func(T, T) int) *BTree[T] { return btree.New(comparator) }
func NewOrderedBTree[T cmp.Ordered]() *BTree[T]           { return btree.NewOrdered[T]() }

type PriorityQueue[T any] = priorityqueue.PriorityQueue[T]

func NewMinPriorityQueue[T any](cmp func(T, T) int) *PriorityQueue[T] {
	return priorityqueue.NewMin(cmp)
}
func NewMaxPriorityQueue[T any](cmp func(T, T) int) *PriorityQueue[T] {
	return priorityqueue.NewMax(cmp)
}
func NewMinOrderedPriorityQueue[T cmp.Ordered]() *PriorityQueue[T] {
	return priorityqueue.NewMinOrdered[T]()
}
func NewMaxOrderedPriorityQueue[T cmp.Ordered]() *PriorityQueue[T] {
	return priorityqueue.NewMaxOrdered[T]()
}

// ===== Hash utility aliases =====

type Hasher[E any] = hashutil.Hasher[E]
type DefaultHasher[E comparable] = hashutil.Default[E]
type SliceHasher[E ~[]T, T any, H Hasher[T]] = hashutil.SliceHasher[E, T, H]
type MapHasher[E ~map[K]V, K comparable, V any, H Hasher[V]] = hashutil.MapHasher[E, K, V, H]
