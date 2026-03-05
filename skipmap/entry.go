// Package skipmap implements ordered key-value mapping based on skip lists.
// This file implements the Entry API, providing Rust-like entry operations for conditional key-value manipulation.
package skipmap

// Entry represents an entry for a specific key in the map, allowing conditional manipulation of the key's value.
// The Entry API provides a higher-level abstraction for common conditional operations,
// such as "update the value if the key exists, otherwise insert a new value".

type Entry[K, V any] struct {
	mapRef *SkipMap[K, V] // Reference to the owning map
	key    K              // Key
	node   *node[K, V]    // Points to the node if key exists; nil otherwise
}

// OrInsert inserts the specified value if the key does not exist and returns a reference to the value;
// if the key already exists, it returns a reference to the existing value without performing insertion.
// Parameters:
//   - value: The value to insert if the key does not exist
//
// Returns:
//   - A pointer to the existing value or the newly inserted value
func (e Entry[K, V]) OrInsert(value V) *V {
	// If the node exists, directly return a reference to the value
	if e.node != nil {
		return &e.node.Value
	}

	// Key does not exist, perform insertion
	e.mapRef.Insert(e.key, value)

	// Re-fetch the node and return a reference to the value
	n, ok := e.mapRef.getNode(e.key)
	if ok {
		return &n.Value
	}

	// Theoretically shouldn't reach here, but return reference to value for safety
	return &value
}

// OrInsertWith uses the provided function to generate a value and inserts it if the key does not exist,
// then returns a reference to the value; if the key already exists, it returns a reference to the existing value
// without performing insertion.
// This is useful when generating the value might be expensive, as it avoids unnecessary computation.
// Parameters:
//   - f: Function used to generate the value if the key does not exist
//
// Returns:
//   - A pointer to the existing value or the newly generated value
func (e Entry[K, V]) OrInsertWith(f func() V) *V {
	// If node exists, directly return reference to the value
	if e.node != nil {
		return &e.node.Value
	}

	// Key doesn't exist, generate value using function and insert
	value := f()
	e.mapRef.Insert(e.key, value)

	// Retrieve node again and return reference to value
	n, ok := e.mapRef.getNode(e.key)
	if ok {
		return &n.Value
	}

	// Theoretically shouldn't reach here, but return reference to value for safety
	return &value
}

// OrInsertWithKey uses the provided function and the key to generate a value and inserts it if the key does not exist,
// then returns a reference to the value; if the key already exists, it returns a reference to the existing value
// without performing insertion.
// This is useful when the value needs to be generated based on the key.
// Parameters:
//   - f: Function used to generate the value if the key does not exist, taking the key as an argument
//
// Returns:
//   - A pointer to the existing value or the newly generated value
func (e Entry[K, V]) OrInsertWithKey(f func(K) V) *V {
	// If node exists, directly return reference to the value
	if e.node != nil {
		return &e.node.Value
	}

	// Key doesn't exist, generate value using function and key and insert
	value := f(e.key)
	e.mapRef.Insert(e.key, value)

	// Retrieve node again and return reference to value
	n, ok := e.mapRef.getNode(e.key)
	if ok {
		return &n.Value
	}

	// Theoretically shouldn't reach here, but return reference to value for safety
	return &value
}

// AndModify applies the specified modification function to the value if the key exists;
// if the key does not exist, no operation is performed.
// Parameters:
//   - modifyFn: Function used to modify the existing value
//
// Returns:
//   - The same Entry, allowing for chained calls to other methods
func (e Entry[K, V]) AndModify(modifyFn func(*V)) Entry[K, V] {
	// Only perform modification if the node exists
	if e.node != nil {
		modifyFn(&e.node.Value)
	}
	return e
}

// Get retrieves the value associated with the key (if it exists).
func (e Entry[K, V]) Get() (V, bool) {
	var zero V
	if e.node != nil {
		return e.node.Value, true
	}
	return zero, false
}

// Insert unconditionally inserts or updates the key-value pair.
// Parameters:
//   - value: The value to insert or update
//
// Returns:
//   - If the key already existed, returns the old value and true
//   - If the key did not exist, returns the zero value and false
func (e Entry[K, V]) Insert(value V) (V, bool) {
	return e.mapRef.insert(e.key, value)
}

// Delete removes the key and reports whether it existed.
func (e Entry[K, V]) Delete() bool {
	_, ok := e.mapRef.Remove(e.key)
	return ok
}

// getNode is an internal method of SkipMap for retrieving the node associated with the specified key.
// This is a helper method added for the needs of the Entry API.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists, returns the node pointer and true
//   - If the key does not exist, returns nil and false

func (sm *SkipMap[K, V]) getNode(key K) (*node[K, V], bool) {
	current := sm.head

	// Start searching from the highest level
	for i := sm.level; i >= 0; i-- {
		// Move forward along the current level until finding a node >= key or reaching end
		for current.next[i] != nil && sm.comparator(current.next[i].Key, key) < 0 {
			current = current.next[i]
		}
	}

	// Reached level 0, current.next[0] is the first node >= key
	current = current.next[0]

	// If the same key is found, return node and true
	if current != nil && sm.comparator(current.Key, key) == 0 {
		return current, true
	}

	// Key does not exist, return nil and false
	return nil, false
}
