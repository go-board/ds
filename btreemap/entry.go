package btreemap

// node represents a key-value pair entry in BTreeMap
type node[K, V any] struct {
	Key   K
	Value V
}

func (n *node[K, V]) kv() (K, V) {
	return n.Key, n.Value
}

func (n *node[K, V]) kvMut() (K, *V) {
	return n.Key, &n.Value
}

// Entry represents a possibly existing or non-existing key entry in BTreeMap, similar to Rust's Entry enum
type Entry[K, V any] struct {
	mapRef *BTreeMap[K, V]
	key    K
	node   *node[K, V] // nil indicates Vacant
}

// OrInsert inserts a value if the key doesn't exist and returns a reference to the value; if the key exists, returns a reference to the existing value
func (e Entry[K, V]) OrInsert(value V) *V {
	// Use Search to directly find if the key exists
	targetEntry := &node[K, V]{Key: e.key}
	existingEntry, found := e.mapRef.btree.Search(targetEntry)
	if found {
		// Key exists, return a reference to the existing value
		return &existingEntry.Value
	}

	// Key doesn't exist, insert the value
	e.mapRef.Insert(e.key, value)

	// Search again to get the reference
	existingEntry, _ = e.mapRef.btree.Search(targetEntry)
	return &existingEntry.Value
}

// OrInsertWith creates a value through a function and inserts it if the key doesn't exist, returns a reference to the value; if the key exists, returns a reference to the existing value
func (e Entry[K, V]) OrInsertWith(f func() V) *V {
	// Use Search to directly find if the key exists
	targetEntry := &node[K, V]{Key: e.key}
	existingEntry, found := e.mapRef.btree.Search(targetEntry)
	if found {
		// Key exists, return a reference to the existing value
		return &existingEntry.Value
	}

	// Key doesn't exist, use the function to create and insert the value
	value := f()
	e.mapRef.Insert(e.key, value)

	// Search again to get a reference
	existingEntry, _ = e.mapRef.btree.Search(targetEntry)
	return &existingEntry.Value
}

// Get retrieves the current value and an existence flag
// If the key exists, returns the value and true; if the key doesn't exist, returns zero value and false
func (e Entry[K, V]) Get() (V, bool) {
	var zero V
	// Use Search to directly find if the key exists
	targetEntry := &node[K, V]{Key: e.key}
	existingEntry, found := e.mapRef.btree.Search(targetEntry)
	if !found {
		return zero, false
	}
	return existingEntry.Value, true
}

// Insert inserts or updates the value and returns the old value if one existed.
func (e Entry[K, V]) Insert(value V) (V, bool) {
	return e.mapRef.insert(e.key, value)
}

// OrInsertWithKey creates a value through a key-related function and inserts it if the key doesn't exist, returns a reference to the value; if the key exists, returns a reference to the existing value
func (e Entry[K, V]) OrInsertWithKey(f func(K) V) *V {
	// search directly using the key to check if it exists
	targetEntry := &node[K, V]{Key: e.key}
	existingEntry, found := e.mapRef.btree.Search(targetEntry)
	if found {
		// Key exists, return a reference to the existing value
		return &existingEntry.Value
	}

	// Key doesn't exist, use the key-related function to create and insert the value
	value := f(e.key)
	e.mapRef.Insert(e.key, value)

	// Search again to get a reference
	existingEntry, _ = e.mapRef.btree.Search(targetEntry)
	return &existingEntry.Value
}

// AndModify modifies the value if the key exists, returns the Entry itself to support method chaining
func (e Entry[K, V]) AndModify(modifyFn func(*V)) Entry[K, V] {
	// Check if the key exists
	if e.node != nil {
		// Key exists, modify the value
		modifyFn(&e.node.Value)
	}
	return e
}

// Delete removes the key and reports whether it existed.
func (e Entry[K, V]) Delete() bool {
	_, ok := e.mapRef.Remove(e.key)
	return ok
}

// GetKey retrieves the key of the Entry
func (e node[K, V]) GetKey() K {
	return e.Key
}

// GetValue retrieves the value of the Entry
func (e node[K, V]) GetValue() V {
	return e.Value
}
