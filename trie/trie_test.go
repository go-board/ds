package trie

import (
	"testing"
)

func TestTrieMapBasicOperations(t *testing.T) {
	// 创建一个新的TrieMap，使用字符串作为键的元素类型，字符串作为值类型
	// 使用正确的comparator函数而不是nil
	tm := New[string, string](func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// 测试插入
	keys := [][]string{
		{"a"},           // 键是单元素切片 ["a"]
		{"a", "b"},      // 键是双元素切片 ["a", "b"]
		{"a", "b", "c"}, // 键是三元素切片 ["a", "b", "c"]
		{"a", "d"},      // 键是双元素切片 ["a", "d"]
		{"e"},           // 键是单元素切片 ["e"]
		{"e", "f"},      // 键是双元素切片 ["e", "f"]
	}
	values := []string{"a", "ab", "abc", "ad", "e", "ef"}

	// 插入键值对
	for i, key := range keys {
		tm.Insert(key, values[i])
	}

	// 验证插入后的大小
	if tm.Len() != len(keys) {
		t.Errorf("Len: expected %d, got %d", len(keys), tm.Len())
	}

	// 测试获取
	for i, key := range keys {
		value, exists := tm.Get(key)
		if !exists {
			t.Errorf("Get: expected key %v to exist", key)
			continue
		}
		if value != values[i] {
			t.Errorf("Get: expected value %s for key %v, got %s", values[i], key, value)
		}
	}

	// 测试不存在的键
	_, exists := tm.Get([]string{"nonexistent"})
	if exists {
		t.Errorf("Get: expected key [nonexistent] to not exist")
	}

	// 测试删除
	for i, key := range keys {
		// 删除键并验证
		value, deleted := tm.Remove(key)
		if !deleted {
			t.Errorf("Remove: expected key %v to be deleted", key)
		} else if value != values[i] {
			t.Errorf("Remove: expected value %s for key %v, got %s", values[i], key, value)
		}

		// 验证键已删除
		_, exists := tm.Get(key)
		if exists {
			t.Errorf("Get: expected key %v to be deleted", key)
		}

		// 验证大小减少
		if tm.Len() != len(keys)-i-1 {
			t.Errorf("Len after delete: expected %d, got %d", len(keys)-i-1, tm.Len())
		}
	}
}

func TestTrieMapPrefixMatching(t *testing.T) {
	// 创建一个新的TrieMap，使用字符串作为键的元素类型，字符串作为值类型
	// 使用正确的comparator函数而不是nil
	tm := New[string, string](func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// 插入一些键值对
	tm.Insert([]string{"a"}, "a")
	tm.Insert([]string{"a", "b"}, "ab")
	tm.Insert([]string{"a", "b", "c"}, "abc")
	tm.Insert([]string{"a", "d"}, "ad")
	tm.Insert([]string{"e"}, "e")
	tm.Insert([]string{"e", "f"}, "ef")

	// 测试ValuesByPrefix

	// 测试ValuesByPrefix
	valuesA := tm.ValuesByPrefix([]string{"a"})
	valuesASlice := collectSeq(valuesA)
	expectedValuesA := []string{"a", "ab", "abc", "ad"}
	if len(valuesASlice) != len(expectedValuesA) {
		t.Errorf("ValuesByPrefix [a]: expected %d values, got %d", len(expectedValuesA), len(valuesASlice))
	} else {
		// 验证所有期望的值都在结果中（不考虑顺序，因为我们在测试中没有实现排序的验证）
		valueMap := make(map[string]bool)
		for _, v := range valuesASlice {
			valueMap[v] = true
		}
		for _, v := range expectedValuesA {
			if !valueMap[v] {
				t.Errorf("ValuesByPrefix [a]: expected value %s in result, not found", v)
			}
		}
	}

	// 测试KeysByPrefix
	keysAB := tm.KeysByPrefix([]string{"a", "b"})
	keysABSlice := collectSeq2D(keysAB)
	expectedKeysAB := [][]string{
		{"a", "b"},
		{"a", "b", "c"},
	}
	if len(keysABSlice) != len(expectedKeysAB) {
		t.Errorf("KeysByPrefix [a b]: expected %d keys, got %d", len(expectedKeysAB), len(keysABSlice))
	} else {
		// 简单验证结果
		keyMap := make(map[string]bool)
		for _, key := range keysABSlice {
			// 将键转换为字符串以便比较
			keyStr := ""
			for _, k := range key {
				keyStr += k
			}
			keyMap[keyStr] = true
		}

		for _, expectedKey := range expectedKeysAB {
			expectedKeyStr := ""
			for _, k := range expectedKey {
				expectedKeyStr += k
			}
			if !keyMap[expectedKeyStr] {
				t.Errorf("KeysByPrefix [a b]: expected key %v in result, not found", expectedKey)
			}
		}
	}

	// 测试不存在的前缀
	valuesNonexistent := tm.ValuesByPrefix([]string{"nonexistent"})
	valuesNonexistentSlice := collectSeq(valuesNonexistent)
	if len(valuesNonexistentSlice) > 0 {
		t.Errorf("ValuesByPrefix [nonexistent]: expected 0 values, got %d", len(valuesNonexistentSlice))
	}

	keysNonexistent := tm.KeysByPrefix([]string{"nonexistent"})
	keysNonexistentSlice := collectSeq2D(keysNonexistent)
	if len(keysNonexistentSlice) > 0 {
		t.Errorf("KeysByPrefix [nonexistent]: expected 0 keys, got %d", len(keysNonexistentSlice))
	}
}

func TestTrieMapIterators(t *testing.T) {
	// 创建一个新的TrieMap，使用字符串作为键的元素类型，字符串作为值类型
	// 使用正确的comparator函数而不是nil
	tm := New[string, string](func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// 插入一些键值对
	tm.Insert([]string{"a"}, "a")
	tm.Insert([]string{"a", "b"}, "ab")
	tm.Insert([]string{"a", "b", "c"}, "abc")
	tm.Insert([]string{"a", "d"}, "ad")
	tm.Insert([]string{"e"}, "e")
	tm.Insert([]string{"e", "f"}, "ef")

	// 测试Keys
	keys := tm.Keys()
	keysSlice := collectSeq2D(keys)
	if len(keysSlice) != tm.Len() {
		t.Errorf("Keys: expected %d keys, got %d", tm.Len(), len(keysSlice))
	}

	// 测试Values
	values := tm.Values()
	valuesSlice := collectSeq(values)
	if len(valuesSlice) != tm.Len() {
		t.Errorf("Values: expected %d values, got %d", tm.Len(), len(valuesSlice))
	}

	// 测试Iter
	keyValues := make(map[string]string)
	for key, value := range tm.Iter() {
		keyStr := ""
		for _, k := range key {
			keyStr += k
		}
		keyValues[keyStr] = value
	}

	// 验证所有键值对都被迭代到
	expectedKeyValues := map[string]string{
		"a":   "a",
		"ab":  "ab",
		"abc": "abc",
		"ad":  "ad",
		"e":   "e",
		"ef":  "ef",
	}

	if len(keyValues) != len(expectedKeyValues) {
		t.Errorf("Iter: expected %d key-value pairs, got %d", len(expectedKeyValues), len(keyValues))
	} else {
		for k, v := range expectedKeyValues {
			if gotV, exists := keyValues[k]; !exists || gotV != v {
				t.Errorf("Iter: expected key %s to map to %s, but got %s", k, v, gotV)
			}
		}
	}

	// 测试IterByPrefix
	prefix := []string{"a"}
	prefixKeyValues := make(map[string]string)
	for key, value := range tm.IterByPrefix(prefix) {
		keyStr := ""
		for _, k := range key {
			keyStr += k
		}
		prefixKeyValues[keyStr] = value
	}

	// 验证所有具有指定前缀的键值对都被迭代到
	expectedPrefixKeyValues := map[string]string{
		"a":   "a",
		"ab":  "ab",
		"abc": "abc",
		"ad":  "ad",
	}

	if len(prefixKeyValues) != len(expectedPrefixKeyValues) {
		t.Errorf("IterByPrefix [%v]: expected %d key-value pairs, got %d", prefix, len(expectedPrefixKeyValues), len(prefixKeyValues))
	} else {
		for k, v := range expectedPrefixKeyValues {
			if gotV, exists := prefixKeyValues[k]; !exists || gotV != v {
				t.Errorf("IterByPrefix [%v]: expected key %s to map to %s, but got %s", prefix, k, v, gotV)
			}
		}
	}

	// 测试IterMut
	// 为所有值添加后缀
	for _, valuePtr := range tm.IterMut() {
		*valuePtr += "_mutated"
	}

	// 验证所有值都被修改
	for keyStr, expectedValue := range expectedKeyValues {
		// 将字符串键转回切片形式
		keySlice := make([]string, 0, len(keyStr))
		for i := 0; i < len(keyStr); i++ {
			keySlice = append(keySlice, string(keyStr[i]))
		}

		value, exists := tm.Get(keySlice)
		if !exists {
			t.Errorf("IterMut: key %s no longer exists after mutation", keyStr)
			continue
		}

		expectedMutatedValue := expectedValue + "_mutated"
		if value != expectedMutatedValue {
			t.Errorf("IterMut: expected key %s to map to %s, but got %s", keyStr, expectedMutatedValue, value)
		}
	}
}

func TestTrieMapDelete(t *testing.T) {
	// 创建一个新的TrieMap，使用字符串作为键的元素类型，字符串作为值类型
	// 使用正确的comparator函数而不是nil
	tm := New[string, string](func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// 插入一些键值对
	tm.Insert([]string{"a"}, "a")
	tm.Insert([]string{"a", "b"}, "ab")
	tm.Insert([]string{"a", "b", "c"}, "abc")
	tm.Insert([]string{"a", "d"}, "ad")

	// 测试删除叶子节点
	_, deleted := tm.Remove([]string{"a", "b", "c"})
	if !deleted {
		t.Errorf("Remove: expected to delete leaf node [a b c]")
	}

	// 验证节点已被删除
	_, exists := tm.Get([]string{"a", "b", "c"})
	if exists {
		t.Errorf("Delete: leaf node [a b c] still exists after deletion")
	}

	// 验证父节点仍然存在
	_, exists = tm.Get([]string{"a", "b"})
	if !exists {
		t.Errorf("Delete: parent node [a b] no longer exists after deleting leaf")
	}

	// 测试删除中间节点（有子节点）
	_, deleted = tm.Remove([]string{"a", "b"})
	if !deleted {
		t.Errorf("Remove: expected to delete intermediate node [a b]")
	}

	// 验证节点已被删除
	_, exists = tm.Get([]string{"a", "b"})
	if exists {
		t.Errorf("Delete: intermediate node [a b] still exists after deletion")
	}

	// 测试删除根节点的直接子节点
	_, deleted = tm.Remove([]string{"a"})
	if !deleted {
		t.Errorf("Remove: expected to delete root child node [a]")
	}

	// 验证节点已被删除
	_, exists = tm.Get([]string{"a"})
	if exists {
		t.Errorf("Delete: root child node [a] still exists after deletion")
	}

	// 验证其他键仍然存在
	_, exists = tm.Get([]string{"a", "d"})
	if !exists {
		t.Errorf("Delete: node [a d] no longer exists after deleting [a]")
	}

	// 测试删除不存在的键
	_, deleted = tm.Remove([]string{"nonexistent"})
	if deleted {
		t.Errorf("Remove: expected to return false for nonexistent key")
	}
}

func TestTrieMapClone(t *testing.T) {
	// 创建一个新的TrieMap，使用字符串作为键的元素类型，字符串作为值类型
	// 使用正确的comparator函数而不是nil
	tm := New[string, string](func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// 插入一些键值对
	tm.Insert([]string{"a"}, "a")
	tm.Insert([]string{"a", "b"}, "ab")
	tm.Insert([]string{"a", "b", "c"}, "abc")
	tm.Insert([]string{"a", "d"}, "ad")

	// 克隆TrieMap
	clone := tm.Clone()

	// 验证克隆的大小与原TrieMap相同
	if clone.Len() != tm.Len() {
		t.Errorf("Clone: expected len %d, got %d", tm.Len(), clone.Len())
	}

	// 验证克隆包含所有键值对
	for key, expectedValue := range tm.Iter() {
		value, exists := clone.Get(key)
		if !exists {
			t.Errorf("Clone: key %v not found in clone", key)
			continue
		}
		if value != expectedValue {
			t.Errorf("Clone: value mismatch for key %v, expected %s, got %s", key, expectedValue, value)
		}
	}

	// 验证克隆是深拷贝
	// 修改克隆中的值
	clone.Insert([]string{"a"}, "modified_a")
	clone.Insert([]string{"new"}, "new_value")

	// 验证原TrieMap未被修改
	value, exists := tm.Get([]string{"a"})
	if !exists || value != "a" {
		t.Errorf("Clone: original TrieMap was modified, expected value 'a', got '%v'", value)
	}

	// 验证原TrieMap不包含新插入的键
	_, exists = tm.Get([]string{"new"})
	if exists {
		t.Errorf("Clone: original TrieMap contains new key inserted in clone")
	}
}

// 辅助函数：将iter.Seq[[]K]转换为[][]K切片
func collectSeq2D[K any](seq func(func([]K) bool)) [][]K {
	result := make([][]K, 0)
	seq(func(k []K) bool {
		// 创建键的副本
		keyCopy := make([]K, len(k))
		copy(keyCopy, k)
		result = append(result, keyCopy)
		return true
	})
	return result
}

// 辅助函数：将iter.Seq[V]转换为[]V切片
func collectSeq[V any](seq func(func(V) bool)) []V {
	result := make([]V, 0)
	seq(func(v V) bool {
		result = append(result, v)
		return true
	})
	return result
}

// TestTrieMapEntryAPI tests the Entry API of TrieMap
func TestTrieMapEntryAPI(t *testing.T) {
	// 创建一个新的TrieMap，使用字符串作为键的元素类型，字符串作为值类型
	tm := New[string, string](func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// 测试Entry.Get方法 - 键不存在的情况
	entry := tm.Entry([]string{"a"})
	val, found := entry.Get()
	if found || val != "" {
		t.Errorf("Entry.Get: expected not found and empty value for nonexistent key, got found=%v, val=%s", found, val)
	}

	// 测试Entry.OrInsert方法 - 插入新值
	valPtr := entry.OrInsert("value_a")
	if *valPtr != "value_a" {
		t.Errorf("Entry.OrInsert: expected value 'value_a', got '%s'", *valPtr)
	}

	// 验证值已被正确插入
	actualVal, found := tm.Get([]string{"a"})
	if !found || actualVal != "value_a" {
		t.Errorf("Entry.OrInsert: expected to insert value, got found=%v, val=%s", found, actualVal)
	}

	// 测试Entry.OrInsert方法 - 键已存在
	valPtr = entry.OrInsert("new_value_a")
	if *valPtr != "value_a" {
		t.Errorf("Entry.OrInsert: expected to not override existing value, got '%s'", *valPtr)
	}

	// 测试Entry.AndModify方法
	entry.AndModify(func(v *string) {
		*v = "modified_value_a"
	})

	// 验证值已被修改
	actualVal, found = tm.Get([]string{"a"})
	if !found || actualVal != "modified_value_a" {
		t.Errorf("Entry.AndModify: expected modified value, got found=%v, val=%s", found, actualVal)
	}

	// 测试Entry.AndModify链式调用
	entry.AndModify(func(v *string) {
		*v += "_extra"
	}).OrInsert("ignored_value")

	// 验证链式调用的结果
	actualVal, found = tm.Get([]string{"a"})
	if !found || actualVal != "modified_value_a_extra" {
		t.Errorf("Entry.AndModify chain: expected chained result, got found=%v, val=%s", found, actualVal)
	}

	// 测试Entry.Insert方法 - 更新现有键
	oldVal, found := entry.Insert("updated_value_a")
	if !found || oldVal != "modified_value_a_extra" {
		t.Errorf("Entry.Insert: expected old value, got found=%v, oldVal=%s", found, oldVal)
	}

	// 验证值已被更新
	actualVal, found = tm.Get([]string{"a"})
	if !found || actualVal != "updated_value_a" {
		t.Errorf("Entry.Insert: expected updated value, got found=%v, val=%s", found, actualVal)
	}

	// 测试Entry.Insert方法 - 插入新键
	newEntry := tm.Entry([]string{"b"})
	oldVal, found = newEntry.Insert("value_b")
	if found || oldVal != "" {
		t.Errorf("Entry.Insert: expected not found and empty old value for new key, got found=%v, oldVal=%s", found, oldVal)
	}

	// 验证新值已被插入
	actualVal, found = tm.Get([]string{"b"})
	if !found || actualVal != "value_b" {
		t.Errorf("Entry.Insert: expected to insert new value, got found=%v, val=%s", found, actualVal)
	}

	// 测试Entry.OrInsertWith方法
	funcEntry := tm.Entry([]string{"c"})
	valPtr = funcEntry.OrInsertWith(func() string {
		return "value_c"
	})
	if *valPtr != "value_c" {
		t.Errorf("Entry.OrInsertWith: expected value 'value_c', got '%s'", *valPtr)
	}

	// 测试Entry.OrInsertWithKey方法
	keyFuncEntry := tm.Entry([]string{"d"})
	valPtr = keyFuncEntry.OrInsertWithKey(func(key []string) string {
		return "value_" + key[0]
	})
	if *valPtr != "value_d" {
		t.Errorf("Entry.OrInsertWithKey: expected value 'value_d', got '%s'", *valPtr)
	}

	// 测试Entry.Delete方法
	deleted := entry.Delete()
	if !deleted {
		t.Errorf("Entry.Delete: expected to delete existing key")
	}

	// 验证键已被删除
	_, found = tm.Get([]string{"a"})
	if found {
		t.Errorf("Entry.Delete: key still exists after deletion")
	}

	// 测试Entry.Delete方法 - 键不存在
	deleted = entry.Delete()
	if deleted {
		t.Errorf("Entry.Delete: expected not to delete nonexistent key")
	}

	// 测试Entry.Delete方法的另一个用例
	bEntry := tm.Entry([]string{"b"})
	deleted = bEntry.Delete()
	if !deleted {
		t.Errorf("Entry.Delete: expected to delete existing key")
	}

	// 验证键已被删除
	_, found = tm.Get([]string{"b"})
	if found {
		t.Errorf("Entry.Delete: key still exists after deletion")
	}
}
