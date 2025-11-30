// Package skipmap implements ordered key-value mapping based on skip lists.
// SkipMap provides ordered operations and efficient query performance with time complexity approaching O(log n).
// Thread safety for all operations depends on usage scenarios; no concurrency safety is guaranteed by default.
//
// Example:
//
//	// Create a new string-to-integer map
//	m := skipmap.NewOrdered[string, int]()
//
//	// Insert a key-value pair
//	_, updated := m.Insert("apple", 5)
//	fmt.Println(updated) // false, because it's a new insertion
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
//	// Iterate over all key-value pairs (sorted by key)
//	for k, v := range m.Iter() {
//		fmt.Printf("%s: %d\n", k, v)
//	}
package skipmap

import (
	"cmp"
	"math/rand"
	"time"
)

const (
	// Maximum level
	maxLevel = 32
	// Level increase probability
	p = 0.5
)

var (
	// Random number generator
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// node represents a node in the skip list
// Each node contains a key, a value, and an array of pointers to subsequent nodes
// Array indices correspond to levels, e.g., next[0] is the successor node at level 0 (lowest level)
// Higher levels have fewer nodes and are used to accelerate lookups

type node[K, V any] struct {
	Key   K             // Key
	Value V             // Value
	next  []*node[K, V] // Array of pointers to subsequent nodes at each level
}

// newNode creates a new skip list node
// Parameters:
//   - key: The node's key
//   - value: The node's value
//   - level: The node's level
//
// Returns:
//   - Pointer to the newly created node
func newNode[K, V any](key K, value V, level int) *node[K, V] {
	return &node[K, V]{
		Key:   key,
		Value: value,
		next:  make([]*node[K, V], level+1),
	}
}

func (n *node[K, V]) kv() (K, V) {
	return n.Key, n.Value
}

func (n *node[K, V]) kvMut() (K, *V) {
	return n.Key, &n.Value
}

// SkipMap implements an ordered map based on skip lists.
// A skip list is a data structure that allows for fast lookups, functioning as a multi-level linked list,
// where each level is ordered and upper levels are subsets of lower levels, used to accelerate the lookup process.
// It provides key-ordered iteration and efficient range queries.
// Key ordering is achieved through a custom comparison function provided during construction.
type SkipMap[K, V any] struct {
	head       *node[K, V]    // Head node, doesn't store actual data
	level      int            // Current maximum level of the skip list
	comparator func(K, K) int // Key comparison function
	length     int            // Element count
}

// New creates a new empty SkipMap using the specified key comparison function.
// Parameters:
//   - comparator: Function for comparing keys, cannot be nil. Should return a negative number if a < b, zero if a == b, and a positive number if a > b.
//
// Returns:
//   - Pointer to the newly created [SkipMap]
func New[K, V any](comparator func(K, K) int) *SkipMap[K, V] {
	if comparator == nil {
		panic("comparator function cannot be nil")
	}

	// Create head node with initial level as maxLevel
	head := &node[K, V]{
		next: make([]*node[K, V], maxLevel+1),
	}

	return &SkipMap[K, V]{
		head:       head,
		level:      0,
		comparator: comparator,
		length:     0,
	}
}

// NewOrdered creates a new empty SkipMap for key types that support ordered comparison.
// This is a convenience function that uses [cmp.Compare] as the comparator.
// Type parameters:
//   - K: Key type, must implement [cmp.Ordered] interface
//   - V: Value type, can be any type
//
// Returns:
//   - Pointer to the newly created [SkipMap]
func NewOrdered[K cmp.Ordered, V any]() *SkipMap[K, V] {
	return New[K, V](cmp.Compare[K])
}

// NewFromMap creates a new [SkipMap] instance from an existing map.
// Parameters:
//   - m: The map to copy key-value pairs from
//
// Type Parameters:
//   - K: Key type, must be comparable
//   - V: Value type
//   - M: Map type, must be a map with comparable keys and any value type
//
// Returns:
//   - Pointer to the newly created [SkipMap]
func NewFromMap[K cmp.Ordered, V any, M ~map[K]V](m M) *SkipMap[K, V] {
	hm := NewOrdered[K, V]()
	for k, v := range m {
		hm.Insert(k, v)
	}
	return hm
}

// randomLevel randomly generates the level for a node
// Returns:
//   - Randomly generated level, ranging from 0 to maxLevel
func randomLevel() int {
	level := 0
	// Increase level with probability p until reaching max level or probability not met
	for random.Float64() < p && level < maxLevel {
		level++
	}
	return level
}

// Insert inserts or updates a key-value pair.
// Parameters:
//   - key: The key to insert or update
//   - value: The value to associate with the key
//
// Returns:
//   - If the key existed, returns the old value and true
//   - If the key did not exist, returns the zero value and false
func (sm *SkipMap[K, V]) Insert(key K, value V) (V, bool) {
	// Used to track nodes that need updates at each level
	update := make([]*node[K, V], maxLevel+1)
	current := sm.head

	// Start searching from the highest level
	for i := sm.level; i >= 0; i-- {
		// Move forward along the current level until finding a node >= key or reaching end
		for current.next[i] != nil && sm.comparator(current.next[i].Key, key) < 0 {
			current = current.next[i]
		}
		update[i] = current
	}

	// Reached level 0, current.next[0] is the first node >= key
	current = current.next[0]

	// If the same key is found, update the value and return the old value
	var oldValue V
	var updated bool
	if current != nil && sm.comparator(current.Key, key) == 0 {
		oldValue = current.Value
		current.Value = value
		updated = true
		return oldValue, updated
	}

	// Generate level for new node
	newLevel := randomLevel()

	// If new node's level is higher than current max level, update max level and update array
	if newLevel > sm.level {
		for i := sm.level + 1; i <= newLevel; i++ {
			update[i] = sm.head
		}
		sm.level = newLevel
	}

	// Create new node
	newNode := newNode(key, value, newLevel)

	// Insert new node at each level
	for i := 0; i <= newLevel; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}

	// Increase element count
	sm.length++

	return oldValue, updated
}

// Get retrieves the value associated with the specified key.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists, returns the associated value and true
//   - If the key does not exist, returns the zero value and false
func (sm *SkipMap[K, V]) Get(key K) (V, bool) {
	current := sm.head

	// Start searching from the highest level
	for i := sm.level; i >= 0; i-- {
		// Move forward along current level until finding a node >= key or reaching end
		for current.next[i] != nil && sm.comparator(current.next[i].Key, key) < 0 {
			current = current.next[i]
		}
	}

	// Reached level 0, check if next node is the key we're looking for
	current = current.next[0]

	// If key is found, return value and true
	if current != nil && sm.comparator(current.Key, key) == 0 {
		return current.Value, true
	}

	// Key does not exist, return zero value and false
	var zero V
	return zero, false
}

// GetMut retrieves a mutable reference to the value associated with the specified key, allowing in-place modification.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists, returns a pointer to the value and true
//   - If the key does not exist, returns nil and false
func (sm *SkipMap[K, V]) GetMut(key K) (*V, bool) {
	current := sm.head

	for i := sm.level; i >= 0; i-- {
		for current.next[i] != nil && sm.comparator(current.next[i].Key, key) < 0 {
			current = current.next[i]
		}
	}

	// Reached level 0, check if next node is the key we're looking for
	current = current.next[0]

	// If key is found, return pointer to value and true
	if current != nil && sm.comparator(current.Key, key) == 0 {
		return &current.Value, true
	}

	// Key does not exist, return nil and false
	return nil, false
}

// GetKeyValue returns the key, value, and existence flag.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists, returns the key, associated value, and true
//   - If the key does not exist, returns zero value key, zero value value, and false
func (sm *SkipMap[K, V]) GetKeyValue(key K) (k K, v V, found bool) {
	current := sm.head

	// Search from the highest level
	for i := sm.level; i >= 0; i-- {
		// Move forward in the current level until we find a node >= key or reach the end of the level
		for current.next[i] != nil && sm.comparator(current.next[i].Key, key) < 0 {
			current = current.next[i]
		}
	}

	// Reached level 0, check if the next node is the key we're looking for
	current = current.next[0]

	// If key is found, return key, value and true
	if current != nil && sm.comparator(current.Key, key) == 0 {
		return current.Key, current.Value, true
	}

	// Key does not exist, return zero key, zero value and false
	return
}

// Remove deletes the key-value pair for the specified key.
// Parameters:
//   - key: The key to delete
//
// Returns:
//   - If the key existed, returns the deleted value and true
//   - If the key did not exist, returns the zero value and false
func (sm *SkipMap[K, V]) Remove(key K) (V, bool) {
	// Used to track nodes that need updates at each level
	update := make([]*node[K, V], maxLevel+1)
	current := sm.head

	// Start searching from the highest level
	for i := sm.level; i >= 0; i-- {
		// Move forward in the current level until we find a node >= key or reach the end of the level
		for current.next[i] != nil && sm.comparator(current.next[i].Key, key) < 0 {
			current = current.next[i]
		}
		update[i] = current
	}

	// Reached level 0, current.next[0] is the first node >= key
	current = current.next[0]

	var oldValue V
	var found bool
	if current != nil && sm.comparator(current.Key, key) == 0 {
		// Save the value to be deleted
		oldValue = current.Value
		found = true

		// Delete node at each level
		for i := 0; i <= sm.level; i++ {
			// If predecessor's next node at current level is not the node to delete, stop
			if update[i].next[i] != current {
				break
			}
			// Update pointer to skip the node to delete
			update[i].next[i] = current.next[i]
		}

		// Update max level of the skip list (if the highest level node was deleted)
		for sm.level > 0 && sm.head.next[sm.level] == nil {
			sm.level--
		}

		// Decrease element count
		sm.length--
	}

	return oldValue, found
}

// ContainsKey checks if the map contains the specified key.
// Parameters:
//   - key: The key to check for
//
// Returns:
//   - true if the key exists, false otherwise
func (sm *SkipMap[K, V]) ContainsKey(key K) bool {
	current := sm.head

	// Start searching from the highest level
	for i := sm.level; i >= 0; i-- {
		// Move forward in the current level until we find a node >= key or reach the end of the level
		for current.next[i] != nil && sm.comparator(current.next[i].Key, key) < 0 {
			current = current.next[i]
		}
	}

	// Reached level 0, check if the next node is the key we're looking for
	current = current.next[0]

	// If key is found, return true
	return current != nil && sm.comparator(current.Key, key) == 0
}

// Entry gets an Entry object for the specified key, used for conditionally inserting or updating values.
// Similar to Rust's entry API, it provides more flexible key-value operations.
// Parameters:
//   - key: The key to operate on
//
// Returns:
//   - An Entry object for the key, which can be used to perform various conditional operations
func (sm *SkipMap[K, V]) Entry(key K) Entry[K, V] {
	// Check if key exists
	current := sm.head
	var foundNode *node[K, V]

	// Start searching from the highest level
	for i := sm.level; i >= 0; i-- {
		// Move forward in the current level until we find a node >= key or reach the end of the level
		for current.next[i] != nil && sm.comparator(current.next[i].Key, key) < 0 {
			current = current.next[i]
		}
	}

	// Reached level 0, check if the next node is the key we're looking for
	current = current.next[0]

	// If key is found, save node reference
	if current != nil && sm.comparator(current.Key, key) == 0 {
		foundNode = current
	}

	return Entry[K, V]{
		mapRef: sm,
		key:    key,
		node:   foundNode,
	}
}

// Len returns the number of key-value pairs in the map.
// Returns:
//   - The number of elements in the map
func (sm *SkipMap[K, V]) Len() int {
	return sm.length
}

// IsEmpty checks if the map is empty (contains no key-value pairs).
// Returns:
//   - true if the map is empty, false otherwise
func (sm *SkipMap[K, V]) IsEmpty() bool {
	return sm.length == 0
}

// Clear removes all key-value pairs from the map, making it empty.
func (sm *SkipMap[K, V]) Clear() {
	// Recreate head node
	sm.head = &node[K, V]{
		next: make([]*node[K, V], maxLevel+1),
	}
	sm.level = 0
	sm.length = 0
}

// Clone creates a deep copy of the map.
// Returns:
//   - A new SkipMap containing all the same key-value pairs
//
// Note: The clone operation copies all key-value pairs, but does not deep copy the keys and values themselves.
func (sm *SkipMap[K, V]) Clone() *SkipMap[K, V] {
	// Create new SkipMap instance
	clone := New[K, V](sm.comparator)

	// Copy all key-value pairs (through iteration)
	for k, v := range sm.Iter() {
		clone.Insert(k, v)
	}

	return clone
}

// Extend adds another iterable collection of key-value pairs to the current map.
// Parameters:
//   - iter: An iterator providing key-value pairs
//
// Behavior:
//   - For each key-value pair, updates the value if the key already exists, otherwise adds a new key-value pair

// iterator-related methods have been moved to `iter.go`

// iterator-related methods have been moved to `iter.go`

// First returns the first (smallest) key-value pair in the map.
// Returns:
//   - If the map is empty, returns zero value key, zero value value, and false
//   - Otherwise returns the smallest key, corresponding value, and true
func (sm *SkipMap[K, V]) First() (k K, v V, found bool) {
	current := sm.head.next[0]
	if current == nil {
		return
	}
	return current.Key, current.Value, true
}

// Last returns the last (largest) key-value pair in the map.
// Returns:
//   - If the map is empty, returns zero value key, zero value value, and false
//   - Otherwise returns the largest key, corresponding value, and true
func (sm *SkipMap[K, V]) Last() (k K, v V, found bool) {
	// For skip list, we need to iterate from level 0 to the end
	current := sm.head.next[0]
	if current == nil {
		return
	}

	// Traverse to the last node
	for current.next[0] != nil {
		current = current.next[0]
	}

	return current.Key, current.Value, true
}

// PopFirst removes and returns the first (smallest) key-value pair.
// Returns:
//   - If the map is empty, returns zero value key, zero value value, and false
//   - Otherwise returns the removed key, corresponding value, and true
func (sm *SkipMap[K, V]) PopFirst() (k K, v V, found bool) {
	current := sm.head.next[0]
	if current == nil {
		return
	}

	// Save the key-value pair to return
	key := current.Key
	value := current.Value

	// Delete node at each level
	for i := 0; i <= sm.level; i++ {
		if sm.head.next[i] != current {
			break
		}
		sm.head.next[i] = current.next[i]
	}

	// Update max level of the skip list
	for sm.level > 0 && sm.head.next[sm.level] == nil {
		sm.level--
	}

	// Decrease element count
	sm.length--

	return key, value, true
}

// GetComparator returns the key comparison function used by the map.
// Returns:
//   - Function for comparing keys
func (sm *SkipMap[K, V]) GetComparator() func(K, K) int {
	return sm.comparator
}

// PopLast removes and returns the last (largest) key-value pair.
// Returns:
//   - If the map is empty, returns zero value key, zero value value, and false
//   - Otherwise returns the removed key, corresponding value, and true
func (sm *SkipMap[K, V]) PopLast() (k K, v V, found bool) {
	// Check if the map is empty
	if sm.IsEmpty() {
		return
	}

	// First find the last node from level 0
	var last *node[K, V] = sm.head.next[0]

	// Traverse to the last node
	for last != nil && last.next[0] != nil {
		last = last.next[0]
	}

	if last == nil {
		return
	}

	// Save the key-value pair to return
	key := last.Key
	value := last.Value

	// Delete the node at each level
	// We need to find the predecessor node of 'last' at each level
	for i := 0; i <= sm.level; i++ {
		current := sm.head
		for current.next[i] != nil && current.next[i] != last {
			current = current.next[i]
		}
		// If predecessor node is found, update the pointer
		if current.next[i] == last {
			current.next[i] = last.next[i]
		}
	}

	// Update the maximum level of the skip list
	for sm.level > 0 && sm.head.next[sm.level] == nil {
		sm.level--
	}

	// Decrease the element count
	sm.length--

	return key, value, true
}
