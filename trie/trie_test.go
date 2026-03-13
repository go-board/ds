package trie

import (
	"testing"
)

// TestTrieMapBasicOperations validates core insert/get/remove flows.
func TestTrieMapBasicOperations(t *testing.T) {
	tm := NewOrdered[string, string]()

	keys := [][]string{
		{"a"},
		{"a", "b"},
		{"a", "b", "c"},
		{"a", "d"},
		{"e"},
		{"e", "f"},
	}
	values := []string{"a", "ab", "abc", "ad", "e", "ef"}

	for i, key := range keys {
		tm.Insert(key, values[i])
	}

	if tm.Len() != len(keys) {
		t.Errorf("Len: expected %d, got %d", len(keys), tm.Len())
	}

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

	_, exists := tm.Get([]string{"nonexistent"})
	if exists {
		t.Errorf("Get: expected key [nonexistent] to not exist")
	}

	for i, key := range keys {
		value, deleted := tm.Remove(key)
		if !deleted {
			t.Errorf("Remove: expected key %v to be deleted", key)
		} else if value != values[i] {
			t.Errorf("Remove: expected value %s for key %v, got %s", values[i], key, value)
		}

		_, exists := tm.Get(key)
		if exists {
			t.Errorf("Get: expected key %v to be deleted", key)
		}

		if tm.Len() != len(keys)-i-1 {
			t.Errorf("Len after delete: expected %d, got %d", len(keys)-i-1, tm.Len())
		}
	}
}

// TestTrieMapPrefixMatching validates prefix-based key and value lookups.
func TestTrieMapPrefixMatching(t *testing.T) {
	tm := NewOrdered[string, string]()

	tm.Insert([]string{"a"}, "a")
	tm.Insert([]string{"a", "b"}, "ab")
	tm.Insert([]string{"a", "b", "c"}, "abc")
	tm.Insert([]string{"a", "d"}, "ad")
	tm.Insert([]string{"e"}, "e")
	tm.Insert([]string{"e", "f"}, "ef")

	// Query values for a broad prefix and compare as a set.
	valuesA := tm.ValuesByPrefix([]string{"a"})
	valuesASlice := collectSeq(valuesA)
	expectedValuesA := []string{"a", "ab", "abc", "ad"}
	if len(valuesASlice) != len(expectedValuesA) {
		t.Errorf("ValuesByPrefix [a]: expected %d values, got %d", len(expectedValuesA), len(valuesASlice))
	} else {
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

	// Query keys for a narrower prefix and compare as a set.
	keysAB := tm.KeysByPrefix([]string{"a", "b"})
	keysABSlice := collectSeq2D(keysAB)
	expectedKeysAB := [][]string{
		{"a", "b"},
		{"a", "b", "c"},
	}
	if len(keysABSlice) != len(expectedKeysAB) {
		t.Errorf("KeysByPrefix [a b]: expected %d keys, got %d", len(expectedKeysAB), len(keysABSlice))
	} else {
		keyMap := make(map[string]bool)
		for _, key := range keysABSlice {
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

	// Non-existent prefixes should produce empty iterators.
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

// TestTrieMapIterators verifies full-map, prefix, and mutable iterator behavior.
func TestTrieMapIterators(t *testing.T) {
	tm := NewOrdered[string, string]()

	tm.Insert([]string{"a"}, "a")
	tm.Insert([]string{"a", "b"}, "ab")
	tm.Insert([]string{"a", "b", "c"}, "abc")
	tm.Insert([]string{"a", "d"}, "ad")
	tm.Insert([]string{"e"}, "e")
	tm.Insert([]string{"e", "f"}, "ef")

	keys := tm.Keys()
	keysSlice := collectSeq2D(keys)
	if len(keysSlice) != tm.Len() {
		t.Errorf("Keys: expected %d keys, got %d", tm.Len(), len(keysSlice))
	}

	values := tm.Values()
	valuesSlice := collectSeq(values)
	if len(valuesSlice) != tm.Len() {
		t.Errorf("Values: expected %d values, got %d", tm.Len(), len(valuesSlice))
	}

	keyValues := make(map[string]string)
	for key, value := range tm.Iter() {
		keyStr := ""
		for _, k := range key {
			keyStr += k
		}
		keyValues[keyStr] = value
	}

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

	prefix := []string{"a"}
	prefixKeyValues := make(map[string]string)
	for key, value := range tm.IterByPrefix(prefix) {
		keyStr := ""
		for _, k := range key {
			keyStr += k
		}
		prefixKeyValues[keyStr] = value
	}

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

	for _, valuePtr := range tm.IterMut() {
		*valuePtr += "_mutated"
	}

	for keyStr, expectedValue := range expectedKeyValues {
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

// TestTrieMapDelete covers leaf, internal, and missing-key delete paths.
func TestTrieMapDelete(t *testing.T) {
	tm := NewOrdered[string, string]()

	tm.Insert([]string{"a"}, "a")
	tm.Insert([]string{"a", "b"}, "ab")
	tm.Insert([]string{"a", "b", "c"}, "abc")
	tm.Insert([]string{"a", "d"}, "ad")

	_, deleted := tm.Remove([]string{"a", "b", "c"})
	if !deleted {
		t.Errorf("Remove: expected to delete leaf node [a b c]")
	}

	_, exists := tm.Get([]string{"a", "b", "c"})
	if exists {
		t.Errorf("Delete: leaf node [a b c] still exists after deletion")
	}

	_, exists = tm.Get([]string{"a", "b"})
	if !exists {
		t.Errorf("Delete: parent node [a b] no longer exists after deleting leaf")
	}

	_, deleted = tm.Remove([]string{"a", "b"})
	if !deleted {
		t.Errorf("Remove: expected to delete intermediate node [a b]")
	}

	_, exists = tm.Get([]string{"a", "b"})
	if exists {
		t.Errorf("Delete: intermediate node [a b] still exists after deletion")
	}

	_, deleted = tm.Remove([]string{"a"})
	if !deleted {
		t.Errorf("Remove: expected to delete root child node [a]")
	}

	_, exists = tm.Get([]string{"a"})
	if exists {
		t.Errorf("Delete: root child node [a] still exists after deletion")
	}

	_, exists = tm.Get([]string{"a", "d"})
	if !exists {
		t.Errorf("Delete: node [a d] no longer exists after deleting [a]")
	}

	_, deleted = tm.Remove([]string{"nonexistent"})
	if deleted {
		t.Errorf("Remove: expected to return false for nonexistent key")
	}
}

// TestTrieMapClone verifies clone completeness and deep-copy semantics.
func TestTrieMapClone(t *testing.T) {
	tm := NewOrdered[string, string]()

	tm.Insert([]string{"a"}, "a")
	tm.Insert([]string{"a", "b"}, "ab")
	tm.Insert([]string{"a", "b", "c"}, "abc")
	tm.Insert([]string{"a", "d"}, "ad")

	clone := tm.Clone()

	if clone.Len() != tm.Len() {
		t.Errorf("Clone: expected len %d, got %d", tm.Len(), clone.Len())
	}

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

	clone.Insert([]string{"a"}, "modified_a")
	clone.Insert([]string{"new"}, "new_value")

	value, exists := tm.Get([]string{"a"})
	if !exists || value != "a" {
		t.Errorf("Clone: original TrieMap was modified, expected value 'a', got '%v'", value)
	}

	_, exists = tm.Get([]string{"new"})
	if exists {
		t.Errorf("Clone: original TrieMap contains new key inserted in clone")
	}
}

// collectSeq2D drains a [][]K iterator into a slice for assertions.
func collectSeq2D[K any](seq func(func([]K) bool)) [][]K {
	result := make([][]K, 0)
	seq(func(k []K) bool {
		keyCopy := make([]K, len(k))
		copy(keyCopy, k)
		result = append(result, keyCopy)
		return true
	})
	return result
}

// collectSeq drains a value iterator into a slice for assertions.
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
	tm := NewOrdered[string, string]()

	entry := tm.Entry([]string{"a"})
	val, found := entry.Get()
	if found || val != "" {
		t.Errorf("Entry.Get: expected not found and empty value for nonexistent key, got found=%v, val=%s", found, val)
	}

	valPtr := entry.OrInsert("value_a")
	if *valPtr != "value_a" {
		t.Errorf("Entry.OrInsert: expected value 'value_a', got '%s'", *valPtr)
	}

	actualVal, found := tm.Get([]string{"a"})
	if !found || actualVal != "value_a" {
		t.Errorf("Entry.OrInsert: expected to insert value, got found=%v, val=%s", found, actualVal)
	}

	valPtr = entry.OrInsert("new_value_a")
	if *valPtr != "value_a" {
		t.Errorf("Entry.OrInsert: expected to not override existing value, got '%s'", *valPtr)
	}

	entry.AndModify(func(v *string) {
		*v = "modified_value_a"
	})

	actualVal, found = tm.Get([]string{"a"})
	if !found || actualVal != "modified_value_a" {
		t.Errorf("Entry.AndModify: expected modified value, got found=%v, val=%s", found, actualVal)
	}

	entry.AndModify(func(v *string) {
		*v += "_extra"
	}).OrInsert("ignored_value")

	actualVal, found = tm.Get([]string{"a"})
	if !found || actualVal != "modified_value_a_extra" {
		t.Errorf("Entry.AndModify chain: expected chained result, got found=%v, val=%s", found, actualVal)
	}

	oldVal, found := entry.Insert("updated_value_a")
	if !found || oldVal != "modified_value_a_extra" {
		t.Errorf("Entry.Insert: expected old value, got found=%v, oldVal=%s", found, oldVal)
	}

	actualVal, found = tm.Get([]string{"a"})
	if !found || actualVal != "updated_value_a" {
		t.Errorf("Entry.Insert: expected updated value, got found=%v, val=%s", found, actualVal)
	}

	newEntry := tm.Entry([]string{"b"})
	oldVal, found = newEntry.Insert("value_b")
	if found || oldVal != "" {
		t.Errorf("Entry.Insert: expected not found and empty old value for new key, got found=%v, oldVal=%s", found, oldVal)
	}

	actualVal, found = tm.Get([]string{"b"})
	if !found || actualVal != "value_b" {
		t.Errorf("Entry.Insert: expected to insert new value, got found=%v, val=%s", found, actualVal)
	}

	funcEntry := tm.Entry([]string{"c"})
	valPtr = funcEntry.OrInsertWith(func() string {
		return "value_c"
	})
	if *valPtr != "value_c" {
		t.Errorf("Entry.OrInsertWith: expected value 'value_c', got '%s'", *valPtr)
	}

	keyFuncEntry := tm.Entry([]string{"d"})
	valPtr = keyFuncEntry.OrInsertWithKey(func(key []string) string {
		return "value_" + key[0]
	})
	if *valPtr != "value_d" {
		t.Errorf("Entry.OrInsertWithKey: expected value 'value_d', got '%s'", *valPtr)
	}

	deleted := entry.Delete()
	if !deleted {
		t.Errorf("Entry.Delete: expected to delete existing key")
	}

	_, found = tm.Get([]string{"a"})
	if found {
		t.Errorf("Entry.Delete: key still exists after deletion")
	}

	deleted = entry.Delete()
	if deleted {
		t.Errorf("Entry.Delete: expected not to delete nonexistent key")
	}

	bEntry := tm.Entry([]string{"b"})
	deleted = bEntry.Delete()
	if !deleted {
		t.Errorf("Entry.Delete: expected to delete existing key")
	}

	_, found = tm.Get([]string{"b"})
	if found {
		t.Errorf("Entry.Delete: key still exists after deletion")
	}
}
