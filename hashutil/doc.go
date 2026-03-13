// Package hashutil provides generic hashing utilities for hash-based collections.
//
// This package defines the Hasher interface for custom hash and equality comparison,
// and provides implementations for common types. It's required by HashMap and HashSet.
//
// # Features
//
//   - Hasher interface for custom hash implementations
//   - Default hasher for comparable types
//   - String hasher for string keys
//   - Slice hasher for slice keys
//   - Map hasher for map keys
//   - Struct hasher for struct keys
//   - Compatible with standard library maphash
//
// # Hasher Interface
//
// Any type used as a key in HashMap or HashSet must implement Hasher:
//
//	type Hasher[T any] interface {
//	    Hash(e T) uint64
//	    Equal(a, b T) bool
//	}
//
// # Built-in Hashers
//
// The package provides hashers for common types:
//
//	// Default hasher for built-in comparable types
//	hasher := hashutil.Default[int]{}
//	hasher := hashutil.Default[string]{}
//
//	// String-specific hasher
//	hasher := hashutil.StringHasher{}
//
//	// For slices (requires element hasher)
//	intHasher := hashutil.Default[int]{}
//	sliceHasher := hashutil.NewSliceHasher[int](intHasher)
//
//	// For maps (requires value hasher)
//	strHasher := hashutil.Default[string]{}
//	mapHasher := hashutil.NewMapHasher[map[string]int](strHasher)
//
//	// For structs (uses reflection)
//	structHasher := hashutil.NewStructHasher[MyStruct]{}
//
// # Custom Hashers
//
// For custom types, implement the Hasher interface:
//
//	type Point struct { X, Y int }
//	pointHasher := hashutil.NewStructHasher[Point]()
//
//	// Or implement manually
//	type PointHasher struct{}
//	func (h PointHasher) Hash(p Point) uint64 {
//	    return hash([]byte(fmt.Sprintf("%d,%d", p.X, p.Y)))
//	}
//	func (h PointHasher) Equal(a, b Point) bool {
//	    return a.X == b.X && a.Y == b.Y
//	}
package hashutil
