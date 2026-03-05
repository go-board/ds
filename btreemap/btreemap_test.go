package btreemap

import (
	"github.com/go-board/ds/bound"
	"testing"
)

func intComparator(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// TestBTreeMapBasic validates insert, lookup, update, and delete behavior.
func TestBTreeMapBasic(t *testing.T) {
	m := New[int, string](intComparator)

	m.Insert(3, "three")
	m.Insert(1, "one")
	m.Insert(4, "four")

	if m.Len() != 3 {
		t.Errorf("Expected size 3, got %d", m.Len())
	}

	val, found := m.Get(1)
	if !found || val != "one" {
		t.Errorf("Expected 'one', got %v, found: %v", val, found)
	}

	val, found = m.Get(3)
	if !found || val != "three" {
		t.Errorf("Expected 'three', got %v, found: %v", val, found)
	}

	val, found = m.Get(4)
	if !found || val != "four" {
		t.Errorf("Expected 'four', got %v, found: %v", val, found)
	}

	_, found = m.Get(10)
	if found {
		t.Error("Expected not found for key 10")
	}

	if !m.ContainsKey(3) {
		t.Error("Expected ContainsKey(3) to return true")
	}

	if m.ContainsKey(10) {
		t.Error("Expected ContainsKey(10) to return false")
	}

	m.Insert(3, "THREE")
	val, found = m.Get(3)
	if !found || val != "THREE" {
		t.Errorf("Expected 'THREE' after update, got %v", val)
	}
	if m.Len() != 3 {
		t.Errorf("Expected size to remain 3 after update, got %d", m.Len())
	}

	m.Remove(1)
	if m.Len() != 2 {
		t.Errorf("Expected size 2 after deletion, got %d", m.Len())
	}

	_, found = m.Get(1)
	if found {
		t.Error("Expected not found for deleted key 1")
	}

	m.Remove(10)
	if m.Len() != 2 {
		t.Errorf("Expected size to remain 2 after deleting non-existent key, got %d", m.Len())
	}
}

// TestBTreeMapClear verifies that Clear removes all entries and resets size.
func TestBTreeMapClear(t *testing.T) {
	m := New[int, string](intComparator)

	for i := 0; i < 10; i++ {
		m.Insert(i, string(rune(i+'a')))
	}

	if m.Len() != 10 {
		t.Errorf("Expected size 10, got %d", m.Len())
	}

	m.Clear()

	if !m.IsEmpty() {
		t.Error("Map should be empty after Clear()")
	}

	if m.Len() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", m.Len())
	}
}

// collectKeys drains a key iterator into a slice for assertions.
func collectKeysAsc(keysIter func(func(int) bool)) []int {
	var keys []int
	keysIter(func(k int) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}

// collectValues drains a value iterator into a slice for assertions.
func collectValuesAsc(valuesIter func(func(string) bool)) []string {
	var values []string
	valuesIter(func(v string) bool {
		values = append(values, v)
		return true
	})
	return values
}

// TestBTreeMapEntryAndEntries checks entry materialization and ordering.
func TestBTreeMapEntryAndEntries(t *testing.T) {
	entry := node[int, string]{Key: 42, Value: "answer"}
	if entry.GetKey() != 42 || entry.GetValue() != "answer" {
		t.Errorf("node creation failed, expected (42, 'answer'), got (%v, %v)", entry.GetKey(), entry.GetValue())
	}

	m := New[int, string](intComparator)

	m.Insert(3, "three")
	m.Insert(1, "one")
	m.Insert(2, "two")

	var collectedEntries []node[int, string]
	for k, v := range m.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewUnbounded[int]())) {
		collectedEntries = append(collectedEntries, node[int, string]{Key: k, Value: v})
	}

	if len(collectedEntries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(collectedEntries))
	}

	expectedEntries := []node[int, string]{
		{Key: 1, Value: "one"},
		{Key: 2, Value: "two"},
		{Key: 3, Value: "three"},
	}

	for i, e := range collectedEntries {
		expected := expectedEntries[i]
		if e.Key != expected.Key || e.Value != expected.Value {
			t.Errorf("Expected entry (%v, %v) at position %d, got (%v, %v)",
				expected.Key, expected.Value, i, e.Key, e.Value)
		}

		if e.GetKey() != e.Key || e.GetValue() != e.Value {
			t.Errorf("Entry Getter methods failed for entry (%v, %v)", e.Key, e.Value)
		}
	}

	emptyMap := New[int, string](intComparator)
	emptyEntriesCount := 0
	for range emptyMap.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewUnbounded[int]())) {
		emptyEntriesCount++
	}
	if emptyEntriesCount != 0 {
		t.Errorf("Expected 0 entries for empty map, got %d", emptyEntriesCount)
	}
}

// TestBTreeMapRangeEntries validates bounded and unbounded range scans.
func TestBTreeMapRangeEntries(t *testing.T) {
	m := New[int, string](intComparator)

	for i := 1; i <= 10; i++ {
		m.Insert(i, "value-"+string(rune('0'+i)))
	}

	lowerBound := 3
	upperBound := 7
	var rangeEntries []node[int, string]
	for k, v := range m.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lowerBound), bound.NewExcluded(upperBound))) {
		rangeEntries = append(rangeEntries, node[int, string]{Key: k, Value: v})
	}

	expectedKeys := []int{3, 4, 5, 6}
	if len(rangeEntries) != len(expectedKeys) {
		t.Errorf("Expected %d entries in range, got %d", len(expectedKeys), len(rangeEntries))
	}

	for i, e := range rangeEntries {
		if e.Key != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, e.Key)
		}
		if e.Value != "value-"+string(rune('0'+e.Key)) {
			t.Errorf("Expected value 'value-%d' for key %d, got '%s'", e.Key, e.Key, e.Value)
		}
	}

	lowerOnly := 8
	var lowerRangeEntries []node[int, string]
	for k, v := range m.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lowerOnly), bound.NewUnbounded[int]())) {
		lowerRangeEntries = append(lowerRangeEntries, node[int, string]{Key: k, Value: v})
	}

	expectedLowerKeys := []int{8, 9, 10}
	if len(lowerRangeEntries) != len(expectedLowerKeys) {
		t.Errorf("Expected %d entries with lower bound, got %d", len(expectedLowerKeys), len(lowerRangeEntries))
	}

	upperOnly := 2
	var upperRangeEntries []node[int, string]
	for k, v := range m.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewExcluded(upperOnly))) {
		upperRangeEntries = append(upperRangeEntries, node[int, string]{Key: k, Value: v})
	}

	expectedUpperKeys := []int{1}
	if len(upperRangeEntries) != len(expectedUpperKeys) {
		t.Errorf("Expected %d entries with upper bound, got %d", len(expectedUpperKeys), len(upperRangeEntries))
	}
}

// TestBTreeMapEntryAPI exercises the Entry API contract for present/missing keys.
func TestBTreeMapEntryAPI(t *testing.T) {
	m := New[int, string](intComparator)

	entry1 := m.Entry(10)
	if m.ContainsKey(10) {
		t.Errorf("Expected key 10 to not exist")
	}

	valPtr1 := entry1.OrInsert("ten")
	if *valPtr1 != "ten" {
		t.Errorf("Expected OrInsert to return 'ten', got '%s'", *valPtr1)
	}

	value, found := m.Get(10)
	if !found || value != "ten" {
		t.Errorf("Expected value 'ten' to be inserted, got %v, found: %v", value, found)
	}

	entry2 := m.Entry(10)
	if !m.ContainsKey(10) {
		t.Errorf("Expected key 10 to exist")
	}

	valPtr2 := entry2.OrInsert("TEN")
	if *valPtr2 != "ten" {
		t.Errorf("Expected OrInsert to return existing value 'ten', got '%s'", *valPtr2)
	}

	value, found = m.Get(10)
	if !found || value != "ten" {
		t.Errorf("Expected value 'ten' to remain unchanged, got %v, found: %v", value, found)
	}

	counter := 0
	entry3 := m.Entry(20)
	valPtr3 := entry3.OrInsertWith(func() string {
		counter++
		return "twenty"
	})

	if *valPtr3 != "twenty" {
		t.Errorf("Expected OrInsertWith to return 'twenty', got '%s'", *valPtr3)
	}
	if counter != 1 {
		t.Errorf("Expected OrInsertWith to call function once, called %d times", counter)
	}

	entry4 := m.Entry(20)
	valPtr4 := entry4.OrInsertWith(func() string {
		counter++
		return "TWENTY"
	})

	if *valPtr4 != "twenty" {
		t.Errorf("Expected OrInsertWith to return existing value 'twenty', got '%s'", *valPtr4)
	}
	if counter != 1 {
		t.Errorf("Expected OrInsertWith to not call function for existing key, called %d times", counter)
	}

	val5, found5 := entry4.Get()
	if !found5 || val5 != "twenty" {
		t.Errorf("Expected Get to return 'twenty' and true, got '%v', found: %v", val5, found5)
	}

	entry5 := m.Entry(30)
	val6, found6 := entry5.Get()
	if found6 {
		t.Errorf("Expected Get to return false for non-existent key, got found: %v, value: '%v'", found6, val6)
	}

	old7, found7 := entry5.Insert("thirty")
	if found7 || old7 != "" {
		t.Errorf("Expected Insert new key to return zero/false, got old=%q found=%v", old7, found7)
	}

	old8, found8 := entry4.Insert("TWENTY")
	if !found8 || old8 != "twenty" {
		t.Errorf("Expected Insert existing key to return old value, got old=%q found=%v", old8, found8)
	}

	value, found = m.Get(20)
	if !found || value != "TWENTY" {
		t.Errorf("Expected value 'TWENTY' to be inserted, got %v, found: %v", value, found)
	}

	entry6 := m.Entry(10)
	entry6.Insert("TEN_MODIFIED")

	value, found = m.Get(10)
	if !found || value != "TEN_MODIFIED" {
		t.Errorf("Expected value to be modified to 'TEN_MODIFIED', got %v, found: %v", value, found)
	}
}

// TestBTreeMapForEach verifies in-order iteration over all entries.
func TestBTreeMapForEach(t *testing.T) {
	m := New[int, string](intComparator)

	m.Insert(3, "three")
	m.Insert(1, "one")
	m.Insert(2, "two")

	var keys []int
	var values []string

	for k, v := range m.IterAsc() {
		keys = append(keys, k)
		values = append(values, v)
	}

	expectedKeys := []int{1, 2, 3}
	expectedValues := []string{"one", "two", "three"}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys in ForEach, got %d", len(expectedKeys), len(keys))
	}

	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, k)
		}
	}

	if len(values) != len(expectedValues) {
		t.Errorf("Expected %d values in ForEach, got %d", len(expectedValues), len(values))
	}

	for i, v := range values {
		if v != expectedValues[i] {
			t.Errorf("Expected value %s at position %d, got %s", expectedValues[i], i, v)
		}
	}
}

// TestBTreeMapRange verifies full scan and early termination behavior.
func TestBTreeMapRangeAsc(t *testing.T) {
	m := New[int, string](intComparator)

	for i := 1; i <= 5; i++ {
		m.Insert(i, string(rune(i+'a'-1)))
	}

	var keys []int
	var values []string

	for k, v := range m.IterAsc() {
		keys = append(keys, k)
		values = append(values, v)
	}

	expectedKeys := []int{1, 2, 3, 4, 5}
	expectedValues := []string{"a", "b", "c", "d", "e"}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys in Range iteration, got %d", len(expectedKeys), len(keys))
	}

	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, k)
		}
	}

	if len(values) != len(expectedValues) {
		t.Errorf("Expected %d values in Range iteration, got %d", len(expectedValues), len(values))
	}

	for i, v := range values {
		if v != expectedValues[i] {
			t.Errorf("Expected value %s at position %d, got %s", expectedValues[i], i, v)
		}
	}

	keys = nil
	values = nil
	for k, v := range m.IterAsc() {
		keys = append(keys, k)
		values = append(values, v)
		if k >= 3 {
			break // Stop iteration early and verify yielded prefix.
		}
	}

	expectedKeys = []int{1, 2, 3}
	expectedValues = []string{"a", "b", "c"}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys in Range with early termination, got %d", len(expectedKeys), len(keys))
	}

	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, k)
		}
	}
}

// TestBTreeMapLargeDataSet provides a coarse regression check on larger input sizes.
func TestBTreeMapLargeDataSet(t *testing.T) {
	m := New[int, int](intComparator)

	for i := 0; i < 1000; i++ {
		m.Insert(i, i*10)
	}

	if m.Len() != 1000 {
		t.Errorf("Expected size 1000, got %d", m.Len())
	}

	for i := 0; i < 1000; i += 100 {
		val, found := m.Get(i)
		if !found || val != i*10 {
			t.Errorf("Expected %d for key %d, got %v, found: %v", i*10, i, val, found)
		}
	}

	for i := 0; i < 1000; i += 2 {
		m.Remove(i)
	}

	if m.Len() != 500 {
		t.Errorf("Expected size 500 after deleting even keys, got %d", m.Len())
	}

	for i := 0; i < 1000; i++ {
		_, found := m.Get(i)
		if i%2 == 0 && found {
			t.Errorf("Expected key %d to be deleted", i)
		} else if i%2 == 1 && !found {
			t.Errorf("Expected key %d to exist", i)
		}
	}
}

// collectStringKeys drains a string-key iterator into a slice for assertions.
func collectStringKeysAsc(keysIter func(func(string) bool)) []string {
	var keys []string
	keysIter(func(k string) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}

func TestBTreeMapGetComparator(t *testing.T) {
	m1 := NewOrdered[string, int]()
	cmp1 := m1.GetComparator()
	if cmp1 == nil {
		t.Fatal("GetComparator should return a non-nil comparator")
	}
	if cmp1("a", "b") >= 0 {
		t.Fatal("ordered comparator should return negative for a < b")
	}

	customCmp := func(a, b string) int {
		if len(a) != len(b) {
			return len(a) - len(b)
		}
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	m2 := New[string, int](customCmp)
	cmp2 := m2.GetComparator()
	if cmp2 == nil {
		t.Fatal("GetComparator should return custom comparator")
	}
	if cmp2("a", "bb") >= 0 {
		t.Fatal("custom comparator should compare by length first")
	}
	if cmp2("xx", "xx") != 0 {
		t.Fatal("custom comparator should return zero for equal keys")
	}
}
