package arraymap

import (
	"cmp"
	"iter"

	"github.com/go-board/ds/internal/kv"
)

// ArrayMap is an ordered map implementation backed by a sorted slice.
//
// Keys are ordered by the provided comparator.
type ArrayMap[K any, V any] struct {
	items      []kv.Pair[K, V]
	comparator func(K, K) int
}

// New creates an empty ArrayMap using comparator.
func New[K any, V any](comparator func(K, K) int) *ArrayMap[K, V] {
	if comparator == nil {
		panic("comparator function cannot be nil")
	}
	return &ArrayMap[K, V]{
		items:      make([]kv.Pair[K, V], 0),
		comparator: comparator,
	}
}

// NewOrdered creates an empty ArrayMap for ordered key types.
func NewOrdered[K cmp.Ordered, V any]() *ArrayMap[K, V] {
	return New[K, V](cmp.Compare[K])
}

// NewFromMap creates an ArrayMap from an existing map.
func NewFromMap[K cmp.Ordered, V any, M ~map[K]V](m M) *ArrayMap[K, V] {
	am := NewOrdered[K, V]()
	for k, v := range m {
		am.Insert(k, v)
	}
	return am
}

func (m *ArrayMap[K, V]) search(key K) (idx int, found bool) {
	left, right := 0, len(m.items)
	for left < right {
		mid := left + (right-left)/2
		cmp := m.comparator(m.items[mid].Key, key)
		if cmp < 0 {
			left = mid + 1
		} else {
			right = mid
		}
	}
	if left < len(m.items) && m.comparator(m.items[left].Key, key) == 0 {
		return left, true
	}
	return left, false
}

func (m *ArrayMap[K, V]) insertAt(idx int, key K, value V) {
	m.items = append(m.items, kv.Pair[K, V]{})
	copy(m.items[idx+1:], m.items[idx:])
	m.items[idx] = kv.NewPair(key, value)
}

// Insert inserts or updates a key-value pair.
func (m *ArrayMap[K, V]) Insert(key K, value V) (V, bool) {
	return m.Entry(key).Insert(value)
}

// Get retrieves the value associated with key.
func (m *ArrayMap[K, V]) Get(key K) (V, bool) {
	idx, found := m.search(key)
	if !found {
		var zero V
		return zero, false
	}
	return m.items[idx].Value, true
}

// GetMut retrieves a writable pointer to the value associated with key.
func (m *ArrayMap[K, V]) GetMut(key K) (*V, bool) {
	idx, found := m.search(key)
	if !found {
		return nil, false
	}
	return &m.items[idx].Value, true
}

// GetKeyValue returns key, value, and whether the key exists.
func (m *ArrayMap[K, V]) GetKeyValue(key K) (K, V, bool) {
	idx, found := m.search(key)
	if !found {
		var zeroV V
		return key, zeroV, false
	}
	p := m.items[idx]
	return p.Key, p.Value, true
}

// Remove removes key and returns the removed value if key exists.
func (m *ArrayMap[K, V]) Remove(key K) (V, bool) {
	idx, found := m.search(key)
	if !found {
		var zero V
		return zero, false
	}
	old := m.items[idx].Value
	m.items = append(m.items[:idx], m.items[idx+1:]...)
	return old, true
}

// ContainsKey reports whether key exists.
func (m *ArrayMap[K, V]) ContainsKey(key K) bool {
	_, found := m.search(key)
	return found
}

// Len returns number of elements.
func (m *ArrayMap[K, V]) Len() int {
	return len(m.items)
}

// IsEmpty reports whether the map has no elements.
func (m *ArrayMap[K, V]) IsEmpty() bool {
	return len(m.items) == 0
}

// Clear removes all elements.
func (m *ArrayMap[K, V]) Clear() {
	m.items = m.items[:0]
}

// Clone creates a shallow copy of the map.
func (m *ArrayMap[K, V]) Clone() *ArrayMap[K, V] {
	clone := &ArrayMap[K, V]{
		items:      make([]kv.Pair[K, V], len(m.items)),
		comparator: m.comparator,
	}
	copy(clone.items, m.items)
	return clone
}

// Entry returns an Entry for key.
func (m *ArrayMap[K, V]) Entry(key K) Entry[K, V] {
	idx, found := m.search(key)
	return Entry[K, V]{mapRef: m, key: key, index: idx, found: found}
}

// Extend inserts all key-value pairs from iterator.
func (m *ArrayMap[K, V]) Extend(it iter.Seq2[K, V]) {
	for k, v := range it {
		m.Insert(k, v)
	}
}

// First returns the first key-value pair in ascending key order.
func (m *ArrayMap[K, V]) First() (k K, v V, found bool) {
	if len(m.items) == 0 {
		return
	}
	p := m.items[0]
	return p.Key, p.Value, true
}

// Last returns the last key-value pair in ascending key order.
func (m *ArrayMap[K, V]) Last() (k K, v V, found bool) {
	if len(m.items) == 0 {
		return
	}
	p := m.items[len(m.items)-1]
	return p.Key, p.Value, true
}

// PopFirst removes and returns the first key-value pair.
func (m *ArrayMap[K, V]) PopFirst() (k K, v V, found bool) {
	if len(m.items) == 0 {
		return
	}
	p := m.items[0]
	m.items = m.items[1:]
	return p.Key, p.Value, true
}

// PopLast removes and returns the last key-value pair.
func (m *ArrayMap[K, V]) PopLast() (k K, v V, found bool) {
	if len(m.items) == 0 {
		return
	}
	idx := len(m.items) - 1
	p := m.items[idx]
	m.items = m.items[:idx]
	return p.Key, p.Value, true
}

// GetComparator returns the key comparison function used by the map.
func (m *ArrayMap[K, V]) GetComparator() func(K, K) int {
	return m.comparator
}
