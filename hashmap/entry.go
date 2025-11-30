package hashmap

import "github.com/go-board/ds/hashutil"

// node represents a key-value pair node in the hash map
type node[K, V any] struct {
	key     K
	value   V
	deleted bool // Mark if the node is deleted (soft delete flag)
}

// Entry represents the state of a key in the hash map, which can be Occupied or Vacant
type Entry[K, V any, H hashutil.Hasher[K]] struct {
	hashMap *HashMap[K, V, H]
	hash    uint64
	key     K
	node    *node[K, V] // nil indicates Vacant
}

// OrInsert inserts the value if the key doesn't exist, returns a mutable reference to the value
func (e Entry[K, V, H]) OrInsert(defaultValue V) *V {
	if e.node != nil {
		// Key exists, return reference to existing value
		return &e.node.value
	}

	// Key doesn't exist, insert new value
	node := &node[K, V]{
		key:     e.key,
		value:   defaultValue,
		deleted: false,
	}
	bucket := e.hashMap.getBucket(e.hash)
	bucket.nodes = append(bucket.nodes, node)
	e.hashMap.size++
	return &node.value
}

// OrInsertWith creates and inserts a value using the function if the key doesn't exist, returns a mutable reference to the value
func (e Entry[K, V, H]) OrInsertWith(defaultValueFn func() V) *V {
	if e.node != nil {
		// Key exists, return reference to existing value
		return &e.node.value
	}

	// Key doesn't exist, create and insert new value using the function
	defaultValue := defaultValueFn()
	node := &node[K, V]{
		key:     e.key,
		value:   defaultValue,
		deleted: false,
	}
	bucket := e.hashMap.getBucket(e.hash)
	bucket.nodes = append(bucket.nodes, node)
	e.hashMap.size++
	return &node.value
}

// OrInsertWithKey creates and inserts a value using the key-related function if the key doesn't exist, returns a mutable reference to the value
func (e Entry[K, V, H]) OrInsertWithKey(defaultValueFn func(K) V) *V {
	if e.node != nil {
		// Key exists, return reference to existing value
		return &e.node.value
	}

	// Key doesn't exist, create and insert new value using the key-related function
	defaultValue := defaultValueFn(e.key)
	node := &node[K, V]{
		key:     e.key,
		value:   defaultValue,
		deleted: false,
	}
	bucket := e.hashMap.getBucket(e.hash)
	bucket.nodes = append(bucket.nodes, node)
	e.hashMap.size++
	return &node.value
}

// AndModify modifies the value if the key exists, returns Entry itself to support chaining
func (e Entry[K, V, H]) AndModify(modifyFn func(*V)) Entry[K, V, H] {
	if e.node != nil {
		// Key exists, modify the value
		modifyFn(&e.node.value)
	}
	return e
}

// Get retrieves the current value and a flag indicating existence
// If the key exists, returns the value and true; if not, returns zero value and false
func (e Entry[K, V, H]) Get() (V, bool) {
	var zero V
	if e.node == nil {
		return zero, false
	}
	return e.node.value, true
}

// Insert inserts the value and returns a reference to it, regardless of whether the key exists
func (e Entry[K, V, H]) Insert(value V) *V {
	if e.node != nil {
		// Key exists, update the value
		e.node.value = value
		return &e.node.value
	}

	// Key doesn't exist, insert new value
	node := &node[K, V]{
		key:     e.key,
		value:   value,
		deleted: false,
	}
	bucket := e.hashMap.getBucket(e.hash)
	bucket.nodes = append(bucket.nodes, node)
	e.hashMap.size++
	return &node.value
}
