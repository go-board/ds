// Package ds provides common interfaces and implementations for various data structures.
// This file defines the Map interface, which is the basic abstraction for all key-value mapping implementations.
package ds

import (
	"iter"

	"github.com/go-board/ds/bound"
)

// Map defines the common interface for key-value mappings.
//
// Type parameters:
//   - K: Key type
//   - V: Value type
type Map[K any, V any] interface {
	// Get retrieves the value associated with the specified key.
	// Parameters:
	//   - key: The key to look up
	//
	// Return values:
	//   - If the key exists, returns the associated value and true
	//   - If the key does not exist, returns the zero value and false
	Get(key K) (V, bool)

	// GetMut retrieves a mutable reference to the value associated with the specified key, allowing in-place modification.
	// Parameters:
	//   - key: The key to look up
	//
	// Return values:
	//   - If the key exists, returns a pointer to the value and true
	//   - If the key does not exist, returns nil and false
	GetMut(key K) (*V, bool)

	// GetKeyValue returns the key, value, and an existence flag.
	// Parameters:
	//   - key: The key to look up
	//
	// Return values:
	//   - If the key exists, returns the key, associated value, and true
	//   - If the key does not exist, returns the zero key, zero value, and false
	GetKeyValue(key K) (K, V, bool)

	// Insert inserts or updates a key-value pair.
	// Parameters:
	//   - key: The key to insert or update
	//   - value: The value to associate
	//
	// Return values:
	//   - If the key already exists, returns the old value and true
	//   - If the key does not exist, returns the zero value and false
	Insert(key K, value V) (V, bool)

	// Remove deletes the key-value pair for the specified key.
	// Parameters:
	//   - key: The key to delete
	//
	// Return values:
	//   - If the key exists, returns the deleted value and true
	//   - If the key does not exist, returns the zero value and false
	Remove(key K) (V, bool)

	// ContainsKey checks if the map contains the specified key.
	// Parameters:
	//   - key: The key to check
	//
	// Return value:
	//   - Returns true if the key exists, otherwise false
	ContainsKey(key K) bool

	// Len returns the number of key-value pairs in the map.
	// Return value:
	//   - The number of elements in the map
	Len() int

	// IsEmpty checks if the map is empty.
	// Return value:
	//   - Returns true if the map is empty, otherwise false
	IsEmpty() bool

	// Clear removes all key-value pairs from the map, making it empty.
	Clear()

	// Extend adds another iterable collection of key-value pairs to the current map.
	// Parameters:
	//   - iter: Iterator providing key-value pairs
	//
	// Behavior:
	//   - For each key-value pair, updates the value if the key already exists, otherwise adds a new key-value pair
	Extend(iter iter.Seq2[K, V])

	// Keys returns an iterator over all keys in the map.
	// Return value:
	//   - Iterator over keys, of type iter.Seq[K]
	//
	// Note: For ordered maps, keys are returned in order; for unordered maps, the order is uncertain.
	Keys() iter.Seq[K]

	// Values returns an iterator over all values in the map.
	// Return value:
	//   - Iterator over values, of type iter.Seq[V]
	//
	// Note: For ordered maps, values are returned in key order; for unordered maps, the order is uncertain.
	Values() iter.Seq[V]

	// ValuesMut returns an iterator over mutable references to all values in the map.
	// Return value:
	//   - Iterator over mutable references to values, of type iter.Seq[*V]
	//
	// Note: For ordered maps, values are returned in key order; for unordered maps, the order is uncertain.
	ValuesMut() iter.Seq[*V]

	// Iter returns an iterator over all key-value pairs in the map.
	// Return value:
	//   - Iterator over all key-value pairs, of type iter.Seq2[K, V]
	//
	// Note: For ordered maps, key-value pairs are returned in key order; for unordered maps, the order is uncertain.
	Iter() iter.Seq2[K, V]

	// IterMut returns a mutable iterator over all key-value pairs in the map.
	// Return value:
	//   - Mutable iterator over all key-value pairs, of type iter.Seq2[K, *V]
	//
	// Note: For ordered maps, key-value pairs are returned in key order; for unordered maps, the order is uncertain.
	IterMut() iter.Seq2[K, *V]
}

// OrderedMap defines an interface for ordered key-value mappings, extending the basic Map interface.
// Ordered maps guarantee the iteration order of key-value pairs, typically by key sorting order.
//
// Type parameters:
//   - K: Key type
//   - V: Value type
//
// Note: Not all Map implementations support ordered operations.
type OrderedMap[K any, V any] interface {
	Map[K, V]

	// First returns the first (smallest key) key-value pair in the map.
	// Return values:
	//   - If the map is non-empty, returns the key, value, and true
	//   - If the map is empty, returns zero key, zero value, and false
	First() (K, V, bool)

	// Last returns the last (largest key) key-value pair in the map.
	// Return values:
	//   - If the map is non-empty, returns the key, value, and true
	//   - If the map is empty, returns zero key, zero value, and false
	Last() (K, V, bool)

	// PopFirst removes and returns the first (smallest key) key-value pair from the map.
	// Return values:
	//   - If the map is non-empty, returns the removed key, value, and true
	//   - If the map is empty, returns zero key, zero value, and false
	PopFirst() (K, V, bool)

	// PopLast removes and returns the last (largest key) key-value pair from the map.
	// Return values:
	//   - If the map is non-empty, returns the removed key, value, and true
	//   - If the map is empty, returns zero key, zero value, and false
	PopLast() (K, V, bool)

	// RangeAsc returns key-value pairs in ascending key order within bounds.
	// Parameters:
	//   - bounds: Range boundary object (inclusive/exclusive/unbounded)
	//
	// Return value:
	//   - Iterator over key-value pairs within the specified range
	RangeAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, V]

	// RangeMutAsc returns mutable key-value pairs in ascending key order within bounds.
	// Parameters:
	//   - bounds: Range boundary object (inclusive/exclusive/unbounded)
	//
	// Return value:
	//   - Mutable iterator over key-value pairs within the specified range
	RangeMutAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, *V]

	// IterDesc returns all key-value pairs in descending key order.
	// Return value:
	//   - Iterator over all key-value pairs in reverse order, of type iter.Seq2[K, V]
	//
	// Note: For ordered maps, key-value pairs are returned in reverse key order; for unordered maps, the order is uncertain.
	IterDesc() iter.Seq2[K, V]

	// IterMutDesc returns mutable key-value pairs in descending key order.
	// Return value:
	//   - Mutable iterator over all key-value pairs in reverse order, of type iter.Seq2[K, *V]
	//
	// Note: For ordered maps, key-value pairs are returned in reverse key order; for unordered maps, the order is uncertain.
	IterMutDesc() iter.Seq2[K, *V]

	// KeysDesc returns all keys in descending order.
	// Return value:
	//   - Iterator over keys in reverse order, of type iter.Seq[K]
	//
	// Note: For ordered maps, keys are returned in reverse order; for unordered maps, the order is uncertain.
	KeysDesc() iter.Seq[K]

	// ValuesDesc returns all values in descending key order.
	// Return value:
	//   - Iterator over values in reverse order, of type iter.Seq[V]
	//
	// Note: For ordered maps, values are returned in reverse key order; for unordered maps, the order is uncertain.
	ValuesDesc() iter.Seq[V]

	// ValuesMutDesc returns mutable values in descending key order.
	// Return value:
	//   - Mutable iterator over values in reverse order, of type iter.Seq[*V]
	//
	// Note: For ordered maps, values are returned in reverse key order; for unordered maps, the order is uncertain.
	ValuesMutDesc() iter.Seq[*V]
}
