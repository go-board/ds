package hashutil

import (
	"hash/maphash"
	"maps"
	"slices"
)

// Hasher defines a generic interface for comparing and hashing values of any type
type Hasher[E any] interface {
	// Equal compares two values for equality
	//
	// Parameters:
	//   - x, y: The two values to compare
	//
	// Returns:
	//   - true if x and y are equal, false otherwise
	Equal(x, y E) bool
	// Hash computes a hash value for the given value using maphash.Hash
	//
	// Parameters:
	//   - h: maphash.Hash instance for computing hash values
	//   - v: The value to hash
	Hash(h *maphash.Hash, v E)
}

// DefaultHasher is a default hasher implementation for comparable types
// It uses Go's standard maphash package for hash computation
type Default[E comparable] struct {
	_ [0]func(*E)
}

// Equal compares two comparable values for equality
//
// Parameters:
//   - x, y: The two values to compare
//
// Returns:
//   - true if x and y are equal, false otherwise
//
// Time Complexity: O(1)
func (Default[E]) Equal(x, y E) bool {
	return x == y
}

// Hash computes hash for comparable values using maphash.Hash
//
// Parameters:
//   - h: maphash.Hash instance for computing hash values
//   - v: The value to hash
//
// Time Complexity: O(1)
func (Default[E]) Hash(h *maphash.Hash, v E) {
	maphash.WriteComparable(h, v)
}

// SliceHasher is a hasher implementation for slice types
// It uses the element type's Hasher to compute the hash for the entire slice
type SliceHasher[E ~[]T, T any, H Hasher[T]] struct {
	hasher H
}

// NewSliceHasher creates a new slice hasher instance
//
// Parameters:
//   - hasher: Hasher for computing hash values of slice elements
//
// Returns:
//   - Newly created slice hasher instance
//
// Time Complexity: O(1)
func NewSliceHasher[E ~[]T, T any, H Hasher[T]](hasher H) SliceHasher[E, T, H] {
	return SliceHasher[E, T, H]{hasher: hasher}
}

// Equal compares two slices for equality
//
// Parameters:
//   - x, y: The two slices to compare
//
// Returns:
//   - true if x and y are equal, false otherwise
//
// Time Complexity: O(n), where n is the length of the slices
func (sh SliceHasher[E, T, H]) Equal(x, y E) bool {
	return slices.EqualFunc(x, y, sh.hasher.Equal)
}

// Hash computes hash for a slice using maphash.Hash
//
// Parameters:
//   - h: maphash.Hash instance for computing hash values
//   - v: The slice to hash
//
// Time Complexity: O(n), where n is the length of the slice
func (sh SliceHasher[E, T, H]) Hash(h *maphash.Hash, v E) {
	for i := range v {
		sh.hasher.Hash(h, v[i])
	}
}

// MapHasher is a hasher implementation for map types
// It uses the key and value types' Hasher to compute the hash for the entire map
type MapHasher[E ~map[K]V, K comparable, V any, H Hasher[V]] struct {
	hasher H
}

// NewMapHasher creates a new map hasher instance
//
// Parameters:
//   - hasher: Hasher for computing hash values of map values
//
// Returns:
//   - Newly created map hasher instance
//
// Time Complexity: O(1)
func NewMapHasher[E ~map[K]V, K comparable, V any, H Hasher[V]](hasher H) MapHasher[E, K, V, H] {
	return MapHasher[E, K, V, H]{hasher: hasher}
}

// Equal compares two maps for equality
//
// Parameters:
//   - x, y: The two maps to compare
//
// Returns:
//   - true if x and y are equal, false otherwise
//
// Time Complexity: O(n), where n is the number of keys in the maps
func (mh MapHasher[E, K, V, H]) Equal(x, y E) bool {
	return maps.EqualFunc(x, y, mh.hasher.Equal)
}

// Hash computes hash for a map using maphash.Hash
//
// Parameters:
//   - h: maphash.Hash instance for computing hash values
//   - v: The map to hash
//
// Time Complexity: O(n), where n is the number of keys in the map
func (mh MapHasher[E, K, V, H]) Hash(h *maphash.Hash, v E) {
	for k := range v {
		mh.hasher.Hash(h, v[k])
	}
}
