package skipmap

import (
	"testing"
)

// TestSkipMapOrderedMapMethods tests the OrderedMap methods added to SkipMap
func TestSkipMapOrderedMapMethods(t *testing.T) {
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
		for k, vPtr := range m.RangeMut(&lower, &upper) {
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

	// Test IterBack
	t.Run("IterBack", func(t *testing.T) {
		expectedKeys := []string{"e", "d", "c", "b", "a"}
		var actualKeys []string
		for k, _ := range m.IterBack() {
			actualKeys = append(actualKeys, k)
		}

		if len(actualKeys) != len(expectedKeys) {
			t.Errorf("Expected %d keys from IterBack, got %d", len(expectedKeys), len(actualKeys))
			return
		}

		for i, key := range expectedKeys {
			if actualKeys[i] != key {
				t.Errorf("IterBack: expected key %s at index %d, got %s", key, i, actualKeys[i])
				return
			}
		}
	})

	// Test IterBackMut
	t.Run("IterBackMut", func(t *testing.T) {
		// Modify values using IterBackMut
		for k, vPtr := range m.IterBackMut() {
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

	// Test KeysBack
	t.Run("KeysBack", func(t *testing.T) {
		expectedKeys := []string{"e", "d", "c", "b", "a"}
		var actualKeys []string
		for k := range m.KeysBack() {
			actualKeys = append(actualKeys, k)
		}

		if len(actualKeys) != len(expectedKeys) {
			t.Errorf("Expected %d keys from KeysBack, got %d", len(expectedKeys), len(actualKeys))
			return
		}

		for i, key := range expectedKeys {
			if actualKeys[i] != key {
				t.Errorf("KeysBack: expected key %s at index %d, got %s", key, i, actualKeys[i])
				return
			}
		}
	})

	// Test ValuesBack
	t.Run("ValuesBack", func(t *testing.T) {
		expectedValues := []int{50, 40, 30, 20, 1}
		var actualValues []int
		for v := range m.ValuesBack() {
			actualValues = append(actualValues, v)
		}

		if len(actualValues) != len(expectedValues) {
			t.Errorf("Expected %d values from ValuesBack, got %d", len(expectedValues), len(actualValues))
			return
		}

		for i, val := range expectedValues {
			if actualValues[i] != val {
				t.Errorf("ValuesBack: expected value %d at index %d, got %d", val, i, actualValues[i])
				return
			}
		}
	})

	// Test ValuesBackMut
	t.Run("ValuesBackMut", func(t *testing.T) {
		// Modify values using ValuesBackMut
		for vPtr := range m.ValuesBackMut() {
			*vPtr *= 2
		}

		// Verify the values were modified
		expectedValues := map[string]int{
			"a": 2,
			"b": 40,
			"c": 60,
			"d": 80,
			"e": 100,
		}

		for k, expectedVal := range expectedValues {
			if actualVal, found := m.Get(k); !found || actualVal != expectedVal {
				t.Errorf("Expected %s:%d, got %d, found: %t", k, expectedVal, actualVal, found)
			}
		}
	})
}

// TestSkipMapEdgeCases tests edge cases for SkipMap methods
func TestSkipMapEdgeCases(t *testing.T) {
	// Test PopFirst with empty map
	t.Run("PopFirstEmpty", func(t *testing.T) {
		m := NewOrdered[string, int]()
		if k, v, found := m.PopFirst(); found {
			t.Errorf("Expected PopFirst to return false for empty map, got (k: %s, v: %d, found: %t)", k, v, found)
		}
	})

	// Test PopLast with empty map
	t.Run("PopLastEmpty", func(t *testing.T) {
		m := NewOrdered[string, int]()
		if k, v, found := m.PopLast(); found {
			t.Errorf("Expected PopLast to return false for empty map, got (k: %s, v: %d, found: %t)", k, v, found)
		}
	})

	// Test PopFirst with single element
	t.Run("PopFirstSingle", func(t *testing.T) {
		m := NewOrdered[string, int]()
		m.Insert("test", 42)
		if k, v, found := m.PopFirst(); !found || k != "test" || v != 42 {
			t.Errorf("Expected PopFirst to return (test, 42, true), got (k: %s, v: %d, found: %t)", k, v, found)
		}
		if !m.IsEmpty() {
			t.Error("Expected map to be empty after PopFirst on single element")
		}
	})

	// Test PopLast with single element
	t.Run("PopLastSingle", func(t *testing.T) {
		m := NewOrdered[string, int]()
		m.Insert("test", 42)
		if k, v, found := m.PopLast(); !found || k != "test" || v != 42 {
			t.Errorf("Expected PopLast to return (test, 42, true), got (k: %s, v: %d, found: %t)", k, v, found)
		}
		if !m.IsEmpty() {
			t.Error("Expected map to be empty after PopLast on single element")
		}
	})

	// Test Range with nil bounds
	t.Run("RangeNilBounds", func(t *testing.T) {
		m := NewOrdered[string, int]()
		m.Insert("a", 1)
		m.Insert("b", 2)
		m.Insert("c", 3)

		// Test with both bounds nil (should return all elements)
		count := 0
		for range m.Range(nil, nil) {
			count++
		}
		if count != 3 {
			t.Errorf("Expected Range(nil, nil) to return all 3 elements, got %d", count)
		}

		// Test with only lower bound nil
		upper := "b"
		count = 0
		for range m.Range(nil, &upper) {
			count++
		}
		if count != 1 {
			t.Errorf("Expected Range(nil, &upper) to return 1 element, got %d", count)
		}

		// Test with only upper bound nil
		lower := "b"
		count = 0
		for range m.Range(&lower, nil) {
			count++
		}
		if count != 2 {
			t.Errorf("Expected Range(&lower, nil) to return 2 elements, got %d", count)
		}
	})

	// Test RangeMut with nil bounds
	t.Run("RangeMutNilBounds", func(t *testing.T) {
		m := NewOrdered[string, int]()
		m.Insert("a", 1)
		m.Insert("b", 2)
		m.Insert("c", 3)

		// Test with both bounds nil (should return all elements)
		count := 0
		for range m.RangeMut(nil, nil) {
			count++
		}
		if count != 3 {
			t.Errorf("Expected RangeMut(nil, nil) to return all 3 elements, got %d", count)
		}
	})
}
