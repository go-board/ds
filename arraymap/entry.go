package arraymap

// Entry represents a potentially occupied key in ArrayMap.
type Entry[K any, V any] struct {
	mapRef *ArrayMap[K, V]
	key    K
	index  int
	found  bool
}

// OrInsert inserts value when key is absent and returns writable pointer.
func (e Entry[K, V]) OrInsert(value V) *V {
	if e.found {
		return &e.mapRef.items[e.index].Value
	}
	e.mapRef.insertAt(e.index, e.key, value)
	return &e.mapRef.items[e.index].Value
}

// OrInsertWith lazily inserts value when key is absent and returns writable pointer.
func (e Entry[K, V]) OrInsertWith(f func() V) *V {
	if e.found {
		return &e.mapRef.items[e.index].Value
	}
	value := f()
	e.mapRef.insertAt(e.index, e.key, value)
	return &e.mapRef.items[e.index].Value
}

// OrInsertWithKey lazily inserts value derived from key when key is absent.
func (e Entry[K, V]) OrInsertWithKey(f func(K) V) *V {
	if e.found {
		return &e.mapRef.items[e.index].Value
	}
	value := f(e.key)
	e.mapRef.insertAt(e.index, e.key, value)
	return &e.mapRef.items[e.index].Value
}

// AndModify modifies current value when key exists.
func (e Entry[K, V]) AndModify(modifyFn func(*V)) Entry[K, V] {
	if e.found {
		modifyFn(&e.mapRef.items[e.index].Value)
	}
	return e
}

// Get returns current value and existence flag.
func (e Entry[K, V]) Get() (V, bool) {
	if e.found {
		return e.mapRef.items[e.index].Value, true
	}
	var zero V
	return zero, false
}

// Insert inserts or updates key and returns old value if existed.
func (e Entry[K, V]) Insert(value V) (V, bool) {
	if e.found {
		old := e.mapRef.items[e.index].Value
		e.mapRef.items[e.index].Value = value
		return old, true
	}
	var zero V
	e.mapRef.insertAt(e.index, e.key, value)
	return zero, false
}

// Delete removes key and reports whether it existed.
func (e Entry[K, V]) Delete() bool {
	_, ok := e.mapRef.Remove(e.key)
	return ok
}
