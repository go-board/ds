package hashmap

import (
	"github.com/go-board/ds/hashutil"
	"github.com/go-board/ds/internal/kv"
)

// Entry represents the state of a key in the hash map, which can be Occupied or Vacant
type Entry[K, V any, H hashutil.Hasher[K]] struct {
	hashMap  *HashMap[K, V, H]
	hash     uint64
	key      K
	index    int
	occupied bool
}

func (e Entry[K, V, H]) pair() *kv.Pair[K, V] {
	if !e.occupied {
		return nil
	}

	bucket := e.hashMap.buckets[e.hash]
	if bucket == nil {
		return nil
	}

	if e.index >= 0 && e.index < len(bucket.nodes) {
		node := &bucket.nodes[e.index]
		if e.hashMap.hasher.Equal(node.Key, e.key) {
			return node
		}
	}

	for i := range bucket.nodes {
		node := &bucket.nodes[i]
		if e.hashMap.hasher.Equal(node.Key, e.key) {
			return node
		}
	}

	return nil
}

// OrInsert inserts the value if the key doesn't exist, returns a mutable reference to the value
func (e Entry[K, V, H]) OrInsert(defaultValue V) *V {
	if node := e.pair(); node != nil {
		// Key exists, return reference to existing value
		return &node.Value
	}

	// Key doesn't exist, insert new value
	node := kv.Pair[K, V]{
		Key:   e.key,
		Value: defaultValue,
	}
	bucket := e.hashMap.getBucket(e.hash)
	bucket.nodes = append(bucket.nodes, node)
	e.hashMap.size++
	return &bucket.nodes[len(bucket.nodes)-1].Value
}

// OrInsertWith creates and inserts a value using the function if the key doesn't exist, returns a mutable reference to the value
func (e Entry[K, V, H]) OrInsertWith(defaultValueFn func() V) *V {
	if node := e.pair(); node != nil {
		// Key exists, return reference to existing value
		return &node.Value
	}

	// Key doesn't exist, create and insert new value using the function
	defaultValue := defaultValueFn()
	node := kv.Pair[K, V]{
		Key:   e.key,
		Value: defaultValue,
	}
	bucket := e.hashMap.getBucket(e.hash)
	bucket.nodes = append(bucket.nodes, node)
	e.hashMap.size++
	return &bucket.nodes[len(bucket.nodes)-1].Value
}

// OrInsertWithKey creates and inserts a value using the key-related function if the key doesn't exist, returns a mutable reference to the value
func (e Entry[K, V, H]) OrInsertWithKey(defaultValueFn func(K) V) *V {
	if node := e.pair(); node != nil {
		// Key exists, return reference to existing value
		return &node.Value
	}

	// Key doesn't exist, create and insert new value using the key-related function
	defaultValue := defaultValueFn(e.key)
	node := kv.Pair[K, V]{
		Key:   e.key,
		Value: defaultValue,
	}
	bucket := e.hashMap.getBucket(e.hash)
	bucket.nodes = append(bucket.nodes, node)
	e.hashMap.size++
	return &bucket.nodes[len(bucket.nodes)-1].Value
}

// AndModify modifies the value if the key exists, returns Entry itself to support chaining
func (e Entry[K, V, H]) AndModify(modifyFn func(*V)) Entry[K, V, H] {
	if node := e.pair(); node != nil {
		// Key exists, modify the value
		modifyFn(&node.Value)
	}
	return e
}

// Get retrieves the current value and a flag indicating existence
// If the key exists, returns the value and true; if not, returns zero value and false
func (e Entry[K, V, H]) Get() (V, bool) {
	var zero V
	node := e.pair()
	if node == nil {
		return zero, false
	}
	return node.Value, true
}

// Insert inserts or updates the value and returns the old value if one existed.
func (e Entry[K, V, H]) Insert(value V) (V, bool) {
	if node := e.pair(); node != nil {
		old := node.Value
		node.Value = value
		return old, true
	}

	bucket := e.hashMap.getBucket(e.hash)
	var zero V
	bucket.nodes = append(bucket.nodes, kv.Pair[K, V]{Key: e.key, Value: value})
	e.hashMap.size++
	return zero, false
}

// Delete removes the entry value and reports whether the key existed.
func (e Entry[K, V, H]) Delete() bool {
	_, ok := e.hashMap.Remove(e.key)
	return ok
}
