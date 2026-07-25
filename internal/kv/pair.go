package kv

// Pair stores one key-value item for internal map implementations.
type Pair[K, V any] struct {
	Key   K
	Value V
}

// NewPair creates a key-value pair.
func NewPair[K, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{Key: key, Value: value}
}

// KV returns the key and value.
func (p *Pair[K, V]) KV() (K, V) {
	return p.Key, p.Value
}

// KVMut returns the key and a mutable value pointer.
func (p *Pair[K, V]) KVMut() (K, *V) {
	return p.Key, &p.Value
}
