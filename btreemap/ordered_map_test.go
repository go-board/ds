package btreemap

import (
	"testing"

	"github.com/go-board/ds/bound"
)

// TestBTreeMapOrderedMapMethods tests the OrderedMap methods added to BTreeMap
func TestBTreeMapOrderedMapMethods(t *testing.T) {
	m := NewOrdered[string, int]()

	// Insert some test data
	m.Insert("a", 1)
	m.Insert("b", 2)
	m.Insert("c", 3)
	m.Insert("d", 4)
	m.Insert("e", 5)

	// Test RangeMut
	t.Run("RangeMut", func(t *testing.T) {
		// Test with range ["b", "d")
		lower := "b"
		upper := "d"
		count := 0
		for k, vPtr := range m.RangeMutAsc(bound.NewRangeBounds(bound.NewIncluded(lower), bound.NewExcluded(upper))) {
			count++
			if k == "b" {
				*vPtr = 20
			}
			if k == "c" {
				*vPtr = 30
			}
		}
		if count != 2 {
			t.Errorf("Expected 2 elements in RangeMut, got %d", count)
		}

		// Verify the values were modified
		if val, found := m.Get("b"); !found || val != 20 {
			t.Errorf("Expected b:20, got %d, found: %t", val, found)
		}
		if val, found := m.Get("c"); !found || val != 30 {
			t.Errorf("Expected c:30, got %d, found: %t", val, found)
		}
	})

	// Test IterDesc
	t.Run("IterDesc", func(t *testing.T) {
		expectedKeys := []string{"e", "d", "c", "b", "a"}
		var actualKeys []string
		for k := range m.IterDesc() {
			actualKeys = append(actualKeys, k)
		}

		if len(actualKeys) != len(expectedKeys) {
			t.Errorf("Expected %d keys from IterDesc, got %d", len(expectedKeys), len(actualKeys))
			return
		}

		for i, key := range expectedKeys {
			if actualKeys[i] != key {
				t.Errorf("IterDesc: expected key %s at index %d, got %s", key, i, actualKeys[i])
				return
			}
		}
	})

	// Test IterMutDesc
	t.Run("IterMutDesc", func(t *testing.T) {
		// Modify values using IterMutDesc
		for k, vPtr := range m.IterMutDesc() {
			if k == "e" {
				*vPtr = 50
			}
			if k == "d" {
				*vPtr = 40
			}
		}

		// Verify the values were modified
		if val, found := m.Get("e"); !found || val != 50 {
			t.Errorf("Expected e:50, got %d, found: %t", val, found)
		}
		if val, found := m.Get("d"); !found || val != 40 {
			t.Errorf("Expected d:40, got %d, found: %t", val, found)
		}
	})
}

// TestBTreeMapMissingMethods tests the methods that were not covered by existing tests
func TestBTreeMapMissingMethods(t *testing.T) {
	m := NewOrdered[string, int]()

	// Test GetKeyValue
	t.Run("GetKeyValue", func(t *testing.T) {
		// Test with non-existent key
		if k, v, found := m.GetKeyValue("nonexistent"); found {
			t.Errorf("Expected GetKeyValue to return false for non-existent key, got (k: %s, v: %d, found: %t)", k, v, found)
		}

		// Test with existing key
		m.Insert("test", 42)
		if k, v, found := m.GetKeyValue("test"); !found || k != "test" || v != 42 {
			t.Errorf("Expected GetKeyValue to return (test, 42, true), got (k: %s, v: %d, found: %t)", k, v, found)
		}
	})

	// Test ContainsKey
	t.Run("ContainsKey", func(t *testing.T) {
		// Test with non-existent key
		if m.ContainsKey("nonexistent") {
			t.Error("Expected ContainsKey to return false for non-existent key")
		}

		// Test with existing key
		if !m.ContainsKey("test") {
			t.Error("Expected ContainsKey to return true for existing key")
		}
	})

	// Test IsEmpty
	t.Run("IsEmpty", func(t *testing.T) {
		// Test with empty map
		emptyMap := NewOrdered[string, int]()
		if !emptyMap.IsEmpty() {
			t.Error("Expected IsEmpty to return true for empty map")
		}

		// Test with non-empty map
		if m.IsEmpty() {
			t.Error("Expected IsEmpty to return false for non-empty map")
		}
	})
}

// TestBTreeMapEntryMethods tests the Entry methods that were not covered by existing tests
func TestBTreeMapEntryMethods(t *testing.T) {
	m := NewOrdered[string, int]()

	// Test OrInsert
	t.Run("EntryOrInsert", func(t *testing.T) {
		// Test with new key
		entry := m.Entry("new")
		valPtr := entry.OrInsert(100)
		if *valPtr != 100 {
			t.Errorf("Expected OrInsert to return pointer to 100 for new key, got pointer to %d", *valPtr)
		}

		// Test with existing key (should not update value)
		entry = m.Entry("new")
		valPtr = entry.OrInsert(200)
		if *valPtr != 100 {
			t.Errorf("Expected OrInsert to return pointer to existing value 100 for existing key, got pointer to %d", *valPtr)
		}
	})

	// Test OrInsertWith
	t.Run("EntryOrInsertWith", func(t *testing.T) {
		// Test with new key
		entry := m.Entry("orinsertwith")
		valPtr := entry.OrInsertWith(func() int { return 300 })
		if *valPtr != 300 {
			t.Errorf("Expected OrInsertWith to return pointer to 300 for new key, got pointer to %d", *valPtr)
		}

		// Test with existing key (should not update value)
		entry = m.Entry("orinsertwith")
		valPtr = entry.OrInsertWith(func() int { return 400 })
		if *valPtr != 300 {
			t.Errorf("Expected OrInsertWith to return pointer to existing value 300 for existing key, got pointer to %d", *valPtr)
		}
	})

	// Test OrInsertWithKey
	t.Run("EntryOrInsertWithKey", func(t *testing.T) {
		// Test with new key
		entry := m.Entry("orinsertwithkey")
		valPtr := entry.OrInsertWithKey(func(key string) int { return len(key) })
		if *valPtr != 15 {
			t.Errorf("Expected OrInsertWithKey to return pointer to 15 for new key, got pointer to %d", *valPtr)
		}

		// Test with existing key (should not update value)
		entry = m.Entry("orinsertwithkey")
		valPtr = entry.OrInsertWithKey(func(key string) int { return len(key) * 2 })
		if *valPtr != 15 {
			t.Errorf("Expected OrInsertWithKey to return pointer to existing value 15 for existing key, got pointer to %d", *valPtr)
		}
	})

	// Test Insert
	t.Run("EntryInsert", func(t *testing.T) {
		entry := m.Entry("entryinsert")
		_, _ = entry.Insert(500)
		if val, found := m.Get("entryinsert"); !found || val != 500 {
			t.Errorf("Expected Insert to set value to 500, got %d, found: %t", val, found)
		}
	})

	// Test AndModify
	t.Run("EntryAndModify", func(t *testing.T) {
		entry := m.Entry("entryinsert")
		entry.AndModify(func(v *int) {
			*v *= 2
		})
		if val, found := m.Get("entryinsert"); !found || val != 1000 {
			t.Errorf("Expected AndModify to double the value to 1000, got %d, found: %t", val, found)
		}

		// Test with non-existent key
		entry = m.Entry("nonexistent")
		entry.AndModify(func(v *int) {
			*v = 999
		})
		if val, found := m.Get("nonexistent"); found {
			t.Errorf("Expected AndModify to do nothing for non-existent key, but found value %d", val)
		}
	})

	// Test Get
	t.Run("EntryGet", func(t *testing.T) {
		// Insert test key first
		m.Insert("test", 42)

		// Test with existing key
		entry := m.Entry("test")
		if val, found := entry.Get(); !found || val != 42 {
			t.Errorf("Expected Get to return (42, true), got (%d, %t)", val, found)
		}

		// Test with non-existent key
		entry = m.Entry("nonexistent")
		if _, found := entry.Get(); found {
			t.Error("Expected Get to return false for non-existent key")
		}
	})
}
