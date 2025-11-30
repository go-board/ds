// Package btreemap implements an ordered key-value map based on B-trees.
// BTreeMap provides ordered operations and efficient range queries with logarithmic time complexity.
// Thread safety for all operations depends on the implementation of the underlying B-tree.
//
// Example:
//
//	// Create a new string-to-integer map
//	m := btreemap.NewOrdered[string, int]()
//
//	// Insert a key-value pair
//	_, updated := m.Insert("apple", 5)
//	fmt.Println(updated) // false, because it's newly inserted
//
//	// Update an existing key
//	oldValue, updated := m.Insert("apple", 10)
//	fmt.Println(oldValue, updated) // 5 true
//
//	// Get a value
//	val, found := m.Get("apple")
//	if found {
//		fmt.Println(val) // 10
//	}
//
//	// Iterate through all key-value pairs (sorted by key)
//	for k, v := range m.Iter() {
//		fmt.Printf("%s: %d\n", k, v)
//	}
package btreemap

import (
	"cmp"

	"github.com/go-board/ds/btree"
)

// BTreeMap implements an ordered map based on B-trees. It provides key-sorted iteration and efficient range queries.
// Key ordering is achieved through a custom comparison function provided at construction time.
// Implemented using generics to support any key and value types.
type BTreeMap[K, V any] struct {
	btree      *btree.BTree[*node[K, V]]
	comparator func(K, K) int
}

// New creates a new empty BTreeMap using the specified key comparison function.
// Parameters:
//   - comparator: Function for comparing keys, must not be nil. Should return a negative number when a < b, 0 when a == b, and a positive number when a > b.
//
// Returns:
//   - Pointer to the newly created BTreeMap
//
// Note: The underlying B-tree uses the default order, determined by the btree package.
func New[K, V any](comparator func(K, K) int) *BTreeMap[K, V] {
	if comparator == nil {
		panic("comparator function cannot be nil")
	}

	// create a comparator function for node entries
	entryComparator := func(a, b *node[K, V]) int {
		return comparator(a.Key, b.Key)
	}

	return &BTreeMap[K, V]{
		btree:      btree.New(entryComparator),
		comparator: comparator,
	}
}

// NewOrdered creates a new empty BTreeMap for key types that support ordered comparison.
// This is a convenience function that uses cmp.Compare as the comparator.
// Type Parameters:
//   - K: Key type, must implement the cmp.Ordered interface
//   - V: Value type, can be any type
//
// Returns:
//   - Pointer to the newly created BTreeMap
func NewOrdered[K cmp.Ordered, V any]() *BTreeMap[K, V] {
	return New[K, V](cmp.Compare[K])
}

// NewFromMap creates a new BTreeMap instance from an existing map.
// Parameters:
//   - m: The map to copy key-value pairs from
//
// Type Parameters:
//   - K: Key type, must be comparable
//   - V: Value type
//   - M: Map type, must be a map with comparable keys and any value type
//
// Returns:
//   - Pointer to the newly created BTreeMap
//
// Example:
//
//	m := map[string]int{"apple": 5, "banana": 10}
//	hm := btreemap.NewFromMap(m)
func NewFromMap[K cmp.Ordered, V any, M ~map[K]V](m M) *BTreeMap[K, V] {
	hm := NewOrdered[K, V]()
	for k, v := range m {
		hm.Insert(k, v)
	}
	return hm
}

// Insert inserts or updates a key-value pair.
// Parameters:
//   - key: The key to insert or update
//   - value: The value to associate with the key
//
// Returns:
//   - If the key exists, returns the old value and true
//   - If the key does not exist, returns the zero value and false
//
// Time Complexity: O(log n)
func (m *BTreeMap[K, V]) Insert(key K, value V) (V, bool) {
	// search for existing entry
	targetEntry := &node[K, V]{Key: key}
	existingEntry, found := m.btree.Search(targetEntry)
	var oldValue V

	// if found, save the old value and update in-place
	if found {
		oldValue = existingEntry.Value
		// update the existing value directly to avoid remove+insert
		existingEntry.Value = value
		return oldValue, true
	}

	// insert new entry
	newEntry := &node[K, V]{Key: key, Value: value}
	m.btree.Insert(newEntry)

	// return zero value and false to indicate a new insertion
	return oldValue, false
}

// Get retrieves the value associated with the specified key.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists, returns the associated value and true
//   - If the key does not exist, returns the zero value and false
//
// Time Complexity: O(log n)
func (m *BTreeMap[K, V]) Get(key K) (V, bool) {
	var zero V
	targetEntry := &node[K, V]{Key: key}
	existingEntry, found := m.btree.Search(targetEntry)
	if !found {
		return zero, false
	}
	return existingEntry.Value, true
}

// GetMut retrieves a mutable reference to the value associated with the specified key, allowing in-place modification of the value.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists, returns a pointer to the value and true
//   - If the key does not exist, returns nil and false
//
// Time Complexity: O(log n)
// Note: The returned pointer is valid only while the map is not modified by other operations.
func (m *BTreeMap[K, V]) GetMut(key K) (*V, bool) {
	targetEntry := &node[K, V]{Key: key}
	existingEntry, found := m.btree.Search(targetEntry)
	if !found {
		return nil, false
	}
	return &existingEntry.Value, true
}

// GetKeyValue returns the key, value, and existence flag.
// This method maintains consistency with HashMap's GetKeyValue behavior.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists, returns the key, associated value, and true
//   - If the key does not exist, returns the zero value key, zero value value, and false
//
// Time Complexity: O(log n)
func (m *BTreeMap[K, V]) GetKeyValue(key K) (k K, v V, found bool) {
	targetEntry := &node[K, V]{Key: key}
	existingEntry, found := m.btree.Search(targetEntry)
	if !found {
		return
	}
	return existingEntry.Key, existingEntry.Value, true
}

// Remove deletes the key-value pair with the specified key.
// Parameters:
//   - key: The key to delete
//
// Returns:
//   - If the key exists, returns the deleted value and true
//   - If the key does not exist, returns the zero value and false
//
// Time Complexity: O(log n)
func (m *BTreeMap[K, V]) Remove(key K) (v V, found bool) {
	targetEntry := &node[K, V]{Key: key}
	// search for the entry to delete to obtain its current value
	existingEntry, found := m.btree.Search(targetEntry)
	if !found {
		return
	}

	// save the value to return
	oldValue := existingEntry.Value

	// perform deletion
	m.btree.Remove(targetEntry)

	return oldValue, true
}

// ContainsKey checks if the map contains the specified key.
// Parameters:
//   - key: The key to check
//
// Returns:
//   - true if the key exists, false otherwise
//
// Time Complexity: O(log n)
func (m *BTreeMap[K, V]) ContainsKey(key K) bool {
	targetEntry := &node[K, V]{Key: key}
	_, found := m.btree.Search(targetEntry)
	return found
}

// Entry retrieves an Entry object for the specified key, used for conditionally inserting or updating values.
// Similar to Rust's entry API, it provides more flexible key-value operations.
// Parameters:
//   - key: The key to operate on
//
// Returns:
//   - Entry object for the corresponding key, which can be used to perform various conditional operations
func (m *BTreeMap[K, V]) Entry(key K) Entry[K, V] {
	entry := Entry[K, V]{mapRef: m, key: key}

	// search for the key and attach node if present
	targetEntry := &node[K, V]{Key: key}
	existingEntry, found := m.btree.Search(targetEntry)
	if found {
		entry.node = existingEntry
	}

	return entry
}

// Len returns the number of key-value pairs in the map.
// Returns:
//   - The number of elements in the map
//
// Time Complexity: O(1)
func (m *BTreeMap[K, V]) Len() int {
	return m.btree.Len()
}

// IsEmpty checks if the map is empty (contains no key-value pairs).
// Returns:
//   - true if the map is empty, false otherwise
//
// Time Complexity: O(1)
func (m *BTreeMap[K, V]) IsEmpty() bool {
	return m.btree.Len() == 0
}

// Clear removes all key-value pairs from the map, making it empty.
// Time Complexity: O(1)
func (m *BTreeMap[K, V]) Clear() {
	m.btree = btree.New(func(a, b *node[K, V]) int {
		return m.comparator(a.Key, b.Key)
	})
}

// Clone creates a deep copy of the map.
// Returns:
//   - New BTreeMap containing all the same key-value pairs
//
// Note: The clone operation copies all key-value pairs, but does not deep copy the keys and values themselves.
// Time Complexity: O(n)
func (m *BTreeMap[K, V]) Clone() *BTreeMap[K, V] {
	// create a new BTreeMap instance
	clone := New[K, V](m.comparator)

	// copy all key-value pairs
	for k, v := range m.Iter() {
		clone.Insert(k, v)
	}

	return clone
}

// First returns the first (smallest key) key-value pair in the map.
// Returns:
//   - If the map is not empty, returns the key, value, and true
//   - If the map is empty, returns the zero value key, zero value value, and false
//
// Time Complexity: O(log n)
func (m *BTreeMap[K, V]) First() (k K, v V, found bool) {
	n, found := m.btree.First()
	if !found {
		return
	}
	return n.Key, n.Value, true
}

// Last returns the last (largest key) key-value pair in the map.
// Returns:
//   - If the map is not empty, returns the key, value, and true
//   - If the map is empty, returns the zero value key, zero value value, and false
//
// Time Complexity: O(log n)
func (m *BTreeMap[K, V]) Last() (k K, v V, found bool) {
	n, found := m.btree.Last()
	if !found {
		return
	}
	return n.Key, n.Value, true
}

// PopFirst removes and returns the first (smallest key) key-value pair from the map.
// Returns:
//   - If the map is not empty, returns the removed key, value, and true
//   - If the map is empty, returns the zero value key, zero value value, and false
//
// Time Complexity: O(log n)
func (m *BTreeMap[K, V]) PopFirst() (k K, v V, found bool) {
	// use btree's PopFirst method directly
	n, found := m.btree.PopFirst()
	if !found {
		return
	}
	return n.Key, n.Value, true
}

// PopLast removes and returns the last (largest key) key-value pair from the map.
// Returns:
//   - If the map is not empty, returns the removed key, value, and true
//   - If the map is empty, returns the zero value key, zero value value, and false
//
// Time Complexity: O(log n)
func (m *BTreeMap[K, V]) PopLast() (k K, v V, found bool) {
	// use btree's PopLast method directly
	n, found := m.btree.PopLast()
	if !found {
		return
	}
	return n.Key, n.Value, true
}
