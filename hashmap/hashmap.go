// Package hashmap implements a generic hash map data structure.
// HashMap provides efficient key-value mapping operations, supporting arbitrary key and value types,
// with fast lookup, insertion, and deletion through hashing algorithms, achieving average O(1) time complexity.
//
// Example:
//
//	// Create a new string-to-integer hash map
//	stringHasher := hashutil.StrHasher{}
//	m := hashmap.NewHashMap[string, int](stringHasher)
//
//	// Insert key-value pairs
//	_, updated := m.Insert("apple", 5)
//	fmt.Println(updated) // false, since it's a new insertion
//
//	// Update an existing key
//	oldValue, updated := m.Insert("apple", 10)
//	fmt.Println(oldValue, updated) // 5 true
//
//	// Get value
//	val, found := m.Get("apple")
//	if found {
//		fmt.Println(val) // 10
//	}
//
//	// Iterate through all key-value pairs
//	for k, v := range m.Iter() {
//		fmt.Printf("%s: %d\n", k, v)
//	}
package hashmap

import (
	"hash/maphash"

	"github.com/go-board/ds/hashutil"
)

// HashMap is an efficient generic hash map implementation supporting arbitrary key and value types.
// It achieves fast query, insertion, and deletion operations through hashing algorithms with average O(1) time complexity.
// Uses type parameter H (implementing Hasher[K] interface) for key hashing and equality comparison,
// supporting custom hashing strategies for improved flexibility and customizability.
// Internally uses a map for bucket storage and provides a soft deletion mechanism to optimize performance.
type HashMap[K, V any, H hashutil.Hasher[K]] struct {
	buckets      map[uint64]*bucket[K, V] // Bucket mapping table (using hash values as keys) - 8 bytes
	hasher       H                        // Specific implementation instance for hash computation and equality comparison - at least 8 bytes
	seed         maphash.Seed             // Hash seed - 16 bytes
	size         int                      // Number of valid elements - 8 bytes
	deletedCount int                      // Number of deleted elements - 8 bytes
}

// bucket represents a bucket in the hash map, storing key-value pairs with hash collisions.
// Each bucket internally maintains an array of nodes for storing key-value pair data.
type bucket[K, V any] struct {
	nodes []*node[K, V] // Array of nodes in the bucket
}

// New creates a new empty HashMap instance.
// Parameters:
//   - hasher: Specific implementation instance for hash computation and equality comparison
//
// Type Parameters:
//   - K: Key type
//   - V: Value type
//   - H: Hasher type, must implement hashutil.Hasher[K] interface
//
// Returns:
//   - Pointer to the newly created HashMap
//
// Example:
//
//	stringHasher := hashutil.StrHasher{}
//	m := hashmap.New[string, int](stringHasher)
func New[K, V any, H hashutil.Hasher[K]](hasher H) *HashMap[K, V, H] {
	// Since map is used as the underlying storage, there's no need to pre-allocate a fixed-size bucket array
	return &HashMap[K, V, H]{
		buckets: make(map[uint64]*bucket[K, V]), // Initialize as empty map
		hasher:  hasher,                         // Store the provided hasher instance
		seed:    maphash.MakeSeed(),
	}
}

// NewComparable creates a new HashMap instance with a default hasher for comparable key types.
//
// Type Parameters:
//   - K: Key type, must be comparable
//   - V: Value type
//
// Returns:
//   - Pointer to the newly created HashMap
//
// Example:
//
//	m := hashmap.NewComparable[string, int]()
func NewComparable[K comparable, V any]() *HashMap[K, V, hashutil.Default[K]] {
	return New[K, V](hashutil.Default[K]{})
}

// NewFromMap creates a new HashMap instance from an existing map.
// Parameters:
//   - m: The map to copy key-value pairs from
//
// Type Parameters:
//   - K: Key type, must be comparable
//   - V: Value type
//   - M: Map type, must be a map with comparable keys and any value type
//
// Returns:
//   - Pointer to the newly created HashMap
//
// Example:
//
//	m := map[string]int{"apple": 5, "banana": 10}
//	hm := hashmap.NewFromMap(m)
func NewFromMap[K comparable, V any, M ~map[K]V](m M) *HashMap[K, V, hashutil.Default[K]] {
	hm := NewComparable[K, V]()
	for k, v := range m {
		hm.Insert(k, v)
	}
	return hm
}

// hash computes the hash value for a key, for internal use.
// Parameters:
//   - key: The key to compute the hash value for
//
// Returns:
//   - The computed hash value (uint64)
func (hm *HashMap[K, V, H]) hash(key K) uint64 {
	var h maphash.Hash
	h.SetSeed(hm.seed)
	hm.hasher.Hash(&h, key) // Use the hasher instance from the struct
	return h.Sum64()
}

// getBucket retrieves or creates a bucket for the specified hash value, for internal use.
// Parameters:
//   - hash: Hash value
//
// Returns:
//   - Pointer to the corresponding bucket, creates a new bucket if it doesn't exist
func (hm *HashMap[K, V, H]) getBucket(hash uint64) *bucket[K, V] {
	b, exists := hm.buckets[hash]
	if !exists {
		// Create new bucket
		b = &bucket[K, V]{nodes: make([]*node[K, V], 0)}
		hm.buckets[hash] = b
	}
	return b
}

// Get retrieves the value associated with the specified key.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists and is not deleted, returns the associated value and true
//   - If the key does not exist or has been deleted, returns the zero value and false
//
// Average Time Complexity: O(1)
func (hm *HashMap[K, V, H]) Get(key K) (v V, found bool) {
	hash := hm.hash(key)
	bucket, exists := hm.buckets[hash]
	if !exists {
		// Bucket doesn't exist, meaning the key doesn't exist
		return
	}

	// Iterate through nodes in the bucket to find the matching key (skipping deleted nodes)
	for _, node := range bucket.nodes {
		if !node.deleted && hm.hasher.Equal(node.key, key) {
			return node.value, true
		}
	}

	// Key does not exist, return zero value and false
	return
}

// GetMut retrieves a mutable reference to the value associated with the specified key, allowing in-place value modification.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists and is not deleted, returns a pointer to the value and true
//   - If the key does not exist or has been deleted, returns nil and false
//
// Average Time Complexity: O(1)
// Note: The returned pointer is only valid as long as the hash map is not modified by other operations.
func (hm *HashMap[K, V, H]) GetMut(key K) (*V, bool) {
	hash := hm.hash(key)
	bucket, exists := hm.buckets[hash]
	if !exists {
		// Bucket doesn't exist, indicating the key doesn't exist
		return nil, false
	}

	// Iterate through nodes in the bucket to find the matching key (skip deleted nodes)
	for _, node := range bucket.nodes {
		if !node.deleted && hm.hasher.Equal(node.key, key) {
			// Return a pointer to the value, allowing direct modification
			return &node.value, true
		}
	}

	// Key does not exist
	return nil, false
}

// GetKeyValue returns the key, value, and existence flag.
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - If the key exists and is not deleted, returns the key, associated value, and true
//   - If the key does not exist or has been deleted, returns the input key, zero value, and false
//
// Average Time Complexity: O(1)
func (hm *HashMap[K, V, H]) GetKeyValue(key K) (K, V, bool) {
	hash := hm.hash(key)
	bucket, exists := hm.buckets[hash]
	if !exists {
		var zeroV V
		return key, zeroV, false
	}

	// Iterate through nodes in the bucket to find the matching key (skip deleted nodes)
	for _, node := range bucket.nodes {
		if !node.deleted && hm.hasher.Equal(node.key, key) {
			return node.key, node.value, true
		}
	}

	var zeroV V
	return key, zeroV, false
}

// ContainsKey checks if the hash map contains the specified key.
// Parameters:
//   - key: The key to check
//
// Return value:
//   - true if the key exists and is not deleted, otherwise false
//
// Average time complexity: O(1)
func (hm *HashMap[K, V, H]) ContainsKey(key K) bool {
	hash := hm.hash(key)
	bucket, exists := hm.buckets[hash]
	if !exists {
		return false
	}

	// Iterate through nodes in the bucket to find a matching undeleted key
	for _, node := range bucket.nodes {
		if !node.deleted && hm.hasher.Equal(node.key, key) {
			return true
		}
	}

	return false
}

// Insert inserts or updates a key-value pair.
// Parameters:
//   - key: The key to insert or update
//   - value: The value to associate
//
// Return values:
//   - The old value and true if the key already exists
//   - Zero value and false if the key does not exist
//
// Average time complexity: O(1)
// Note: When there are many deleted nodes in the hash map, compression is automatically triggered.
func (hm *HashMap[K, V, H]) Insert(key K, value V) (V, bool) {
	// Automatically trigger compression when there are many deleted nodes.
	if hm.deletedCount > hm.size {
		hm.Compact()
	}
	return hm.Entry(key).Insert(value)
}

// Remove softly deletes the key-value pair with the specified key.
// Parameters:
//   - key: The key to delete
//
// Return values:
//   - The deleted value and true if the key exists and is not deleted
//   - Zero value and false if the key does not exist or is already deleted
//
// Average time complexity: O(1)
// Note: Soft deletion does not immediately free memory but marks the node as deleted.
// Memory can be freed by calling the Compact method.
func (hm *HashMap[K, V, H]) Remove(key K) (V, bool) {
	hash := hm.hash(key)
	bucket, exists := hm.buckets[hash]
	if !exists {
		var zeroV V
		return zeroV, false
	}

	// Find and mark the node as deleted
	for _, node := range bucket.nodes {
		if !node.deleted && hm.hasher.Equal(node.key, key) {
			// Soft deletion: mark the node as deleted
			node.deleted = true
			oldValue := node.value
			hm.size--
			hm.deletedCount++
			return oldValue, true
		}
	}

	// Key doesn't exist
	var zeroV V
	return zeroV, false
}

// Len returns the number of valid elements in the hash map.
// Return value:
//   - The current number of undeleted key-value pairs in the hash map
//
// Time complexity: O(1)
func (hm *HashMap[K, V, H]) Len() int {
	return hm.size
}

// IsEmpty checks if the hash map is empty (contains no valid key-value pairs).
// Return value:
//   - true if the hash map is empty, otherwise false
//
// Time complexity: O(1)
func (hm *HashMap[K, V, H]) IsEmpty() bool {
	return hm.size == 0
}

// Clear empties the hash map, removing all key-value pairs.
// Time complexity: O(1)
// Note: This operation recreates the internal map, completely freeing all memory.
func (hm *HashMap[K, V, H]) Clear() {
	hm.buckets = make(map[uint64]*bucket[K, V])
	hm.size = 0
	hm.deletedCount = 0
}

// Clone creates a deep copy of the hash map.
// Return value:
//   - A new HashMap containing all the same key-value pairs
//
// Time complexity: O(n), where n is the number of elements
// Note: The clone operation copies all key-value pairs but does not deep copy the keys and values themselves.
func (hm *HashMap[K, V, H]) Clone() *HashMap[K, V, H] {
	// Create new HashMap instance
	clone := &HashMap[K, V, H]{
		buckets:      make(map[uint64]*bucket[K, V], len(hm.buckets)),
		hasher:       hm.hasher, // Copy hasher instance
		seed:         hm.seed,
		size:         hm.size,
		deletedCount: hm.deletedCount,
	}

	// Copy all buckets and nodes
	for hash, b := range hm.buckets {
		// Create new bucket
		newBucket := &bucket[K, V]{
			nodes: make([]*node[K, V], len(b.nodes)),
		}

		// Copy all nodes
		for idx, n := range b.nodes {
			// Create new node and copy content
			newNode := &node[K, V]{
				key:     n.key,
				value:   n.value,
				deleted: n.deleted,
			}
			newBucket.nodes[idx] = newNode
		}

		clone.buckets[hash] = newBucket
	}

	return clone
}

// Compact compresses the hash map, removing all deleted nodes.
// This operation cleans up soft-deleted nodes, frees memory, and improves traversal and lookup efficiency.
// Time complexity: O(n), where n is the number of elements
func (hm *HashMap[K, V, H]) Compact() {
	// Iterate through all buckets
	for hash, bucket := range hm.buckets {
		// Filter out deleted nodes
		activeNodes := make([]*node[K, V], 0, len(bucket.nodes))
		for _, node := range bucket.nodes {
			if !node.deleted {
				activeNodes = append(activeNodes, node)
			}
		}

		if len(activeNodes) > 0 {
			// Update node list in bucket
			bucket.nodes = activeNodes
		} else {
			// If bucket is empty, remove it from the map
			delete(hm.buckets, hash)
		}
	}

	// Reset deleted count
	hm.deletedCount = 0
}

// Extend adds another iterable key-value pair collection to the current hash map.
// Parameters:
//   - iter: Iterator providing key-value pairs
//
// Behavior:
//   - For each key-value pair, if the key exists, update its value; otherwise add a new key-value pair
//
// Average time complexity: O(n), where n is the number of elements in the iterator
// Note: Before operation, it checks if compression is needed to optimize performance.

// Entry gets the Entry state for a key, used for flexible handling of insertion/update operations.
// Provides a Rust-like Entry API, supporting more complex conditional operations.
// Parameters:
//   - key: The key to operate on
//
// Return value:
//   - Entry object for the corresponding key, which can be used to perform various conditional operations
func (hm *HashMap[K, V, H]) Entry(key K) Entry[K, V, H] {
	hash := hm.hash(key)
	entry := Entry[K, V, H]{
		hashMap: hm,
		hash:    hash,
		key:     key,
		node:    nil,
	}

	// Find if key exists
	bucket, exists := hm.buckets[hash]
	if exists {
		for _, node := range bucket.nodes {
			if !node.deleted && hm.hasher.Equal(node.key, key) {
				entry.node = node
				break
			}
		}
	}

	return entry
}
