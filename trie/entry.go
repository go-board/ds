// Package trie implements a generic trie (prefix tree) data structure.
package trie

// Entry represents a potential entry in the TrieMap. It can be used to check if a key exists
// and to conditionally modify, insert, or delete values.
type Entry[K any, V any] struct {
	m   *TrieMap[K, V]
	key []K
}

// Get returns the current value for the key and a boolean indicating whether it exists.
//
// Returns:
//   - value: The value associated with the key (undefined if not found)
//   - found: A boolean indicating whether the key exists
func (e Entry[K, V]) Get() (V, bool) {
	return e.m.Get(e.key)
}

// OrInsert inserts the value into the map if the key is not already present, and returns
// a reference to the existing or inserted value.
//
// Parameters:
//   - value: The value to insert if the key is not present
//
// Returns:
//   - A pointer to the existing value if the key was present, or the new value if it was inserted
func (e Entry[K, V]) OrInsert(value V) *V {
	// 尝试获取现有值
	if _, found := e.m.Get(e.key); found {
		// 键已存在，返回值的指针
		node, _, _ := e.m.findNode(e.key)
		if node != nil {
			return &node.value
		}
		// 如果由于某种原因无法获取节点，直接插入并再次尝试
	}

	// 插入新值
	e.m.Insert(e.key, value)

	// 获取插入后的节点并返回值的指针
	node, _, _ := e.m.findNode(e.key)
	if node != nil {
		return &node.value
	}

	// 理论上不应该到达这里，但为了安全返回nil
	return nil
}

// OrInsertWith inserts a value using the provided function if the key is not already present,
// and returns a reference to the existing or inserted value.
//
// Parameters:
//   - f: A function that returns the value to insert
//
// Returns:
//   - A pointer to the existing value if the key was present, or the new value if it was inserted
func (e Entry[K, V]) OrInsertWith(f func() V) *V {
	// 尝试获取现有值
	if _, found := e.m.Get(e.key); found {
		// 键已存在，返回值的指针
		node, _, _ := e.m.findNode(e.key)
		if node != nil {
			return &node.value
		}
		// 如果由于某种原因无法获取节点，直接插入并再次尝试
	}

	// 使用函数生成值并插入
	value := f()
	e.m.Insert(e.key, value)

	// 获取插入后的节点并返回值的指针
	node, _, _ := e.m.findNode(e.key)
	if node != nil {
		return &node.value
	}

	// 理论上不应该到达这里，但为了安全返回nil
	return nil
}

// OrInsertWithKey inserts a value using the provided function with the key as an argument
// if the key is not already present, and returns a reference to the existing or inserted value.
//
// Parameters:
//   - f: A function that takes the key and returns the value to insert
//
// Returns:
//   - A pointer to the existing value if the key was present, or the new value if it was inserted
func (e Entry[K, V]) OrInsertWithKey(f func([]K) V) *V {
	// 尝试获取现有值
	if _, found := e.m.Get(e.key); found {
		// 键已存在，返回值的指针
		node, _, _ := e.m.findNode(e.key)
		if node != nil {
			return &node.value
		}
		// 如果由于某种原因无法获取节点，直接插入并再次尝试
	}

	// 使用函数生成值并插入
	value := f(e.key)
	e.m.Insert(e.key, value)

	// 获取插入后的节点并返回值的指针
	node, _, _ := e.m.findNode(e.key)
	if node != nil {
		return &node.value
	}

	// 理论上不应该到达这里，但为了安全返回nil
	return nil
}

// AndModify applies the provided function to the value if the key exists, and returns
// the Entry itself to support chaining.
//
// Parameters:
//   - modifyFn: A function that takes a pointer to the value and modifies it
//
// Returns:
//   - The same Entry, allowing for chained calls to other methods
func (e Entry[K, V]) AndModify(modifyFn func(*V)) Entry[K, V] {
	// 使用m.findNode查找节点，避免重复实现查找逻辑
	node, _, _ := e.m.findNode(e.key)
	if node != nil && node.hasValue {
		// 获取值指针并应用修改函数
		valPtr := &node.value
		modifyFn(valPtr)
	}
	return e
}

// Insert unconditionally inserts or updates the key-value pair and returns the previous value (if any) and a boolean indicating if the key existed.
//
// Parameters:
//   - value: The value to insert or update
//
// Returns:
//   - If the key already existed, returns the old value and true
//   - If the key did not exist, returns the zero value and false
func (e Entry[K, V]) Insert(value V) (V, bool) {
	// 获取旧值
	oldVal, found := e.m.Get(e.key)
	// 插入或更新值
	e.m.Insert(e.key, value)
	return oldVal, found
}

// Delete removes the key from the map and returns whether the key was present.
//
// Returns:
//   - A boolean indicating whether the key was present and removed
func (e Entry[K, V]) Delete() bool {
	// 使用Remove方法并提取布尔结果
	_, removed := e.m.Remove(e.key)
	return removed
}
