package skipmap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// Test basic insertion and retrieval functionality
func TestSkipMapBasicOperations(t *testing.T) {
	// Create a new SkipMap
	sm := NewOrdered[string, int]()

	// Test insertion
	oldValue, updated := sm.Insert("apple", 5)
	if oldValue != 0 || updated != false {
		t.Errorf("Insert(\"apple\", 5) should return 0, false, got %d, %v", oldValue, updated)
	}

	// Test update
	oldValue, updated = sm.Insert("apple", 10)
	if oldValue != 5 || updated != true {
		t.Errorf("Insert(\"apple\", 10) should return 5, true, got %d, %v", oldValue, updated)
	}

	// Test retrieval
	value, found := sm.Get("apple")
	if !found || value != 10 {
		t.Errorf("Get(\"apple\") should return 10, true, got %d, %v", value, found)
	}

	// Test retrieval of non-existent key
	value, found = sm.Get("orange")
	if found || value != 0 {
		t.Errorf("Get(\"orange\") should return 0, false, got %d, %v", value, found)
	}

	// Test key existence
	if !sm.ContainsKey("apple") {
		t.Error("ContainsKey(\"apple\") should return true")
	}
	if sm.ContainsKey("orange") {
		t.Error("ContainsKey(\"orange\") should return false")
	}

	// Test deletion
	oldValue, found = sm.Remove("apple")
	if !found || oldValue != 10 {
		t.Errorf("Remove(\"apple\") should return 10, true, got %d, %v", oldValue, found)
	}

	// Test deletion of non-existent key
	oldValue, found = sm.Remove("orange")
	if found || oldValue != 0 {
		t.Errorf("Remove(\"orange\") should return 0, false, got %d, %v", oldValue, found)
	}
}

// Test length and emptiness check functionality
func TestSkipMapLenAndEmpty(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Initial state should be empty
	if !sm.IsEmpty() {
		t.Error("New skip map should be empty")
	}
	if sm.Len() != 0 {
		t.Errorf("New skip map length should be 0, got %d", sm.Len())
	}

	// Length should increase after insertion
	sm.Insert("a", 1)
	if sm.IsEmpty() {
		t.Error("Skip map should not be empty after insertion")
	}
	if sm.Len() != 1 {
		t.Errorf("Skip map length should be 1 after insertion, got %d", sm.Len())
	}

	// Insert more elements
	sm.Insert("b", 2)
	sm.Insert("c", 3)
	if sm.Len() != 3 {
		t.Errorf("Skip map length should be 3 after insertions, got %d", sm.Len())
	}

	// Length should decrease after deletion
	sm.Remove("b")
	if sm.Len() != 2 {
		t.Errorf("Skip map length should be 2 after removal, got %d", sm.Len())
	}

	// Clear the map
	sm.Clear()
	if !sm.IsEmpty() {
		t.Error("Skip map should be empty after Clear()")
	}
	if sm.Len() != 0 {
		t.Errorf("Skip map length should be 0 after Clear(), got %d", sm.Len())
	}
}

// Test iterator functionality
func TestSkipMapIterators(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert some ordered key-value pairs
	data := []struct {
		Key   string
		Value int
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
		{"e", 5},
	}

	for _, item := range data {
		sm.Insert(item.Key, item.Value)
	}

	// Test Iter()
	iterData := make(map[string]int)
	for k, v := range sm.Iter() {
		iterData[k] = v
	}

	for _, item := range data {
		if val, ok := iterData[item.Key]; !ok || val != item.Value {
			t.Errorf("Iter() missing or wrong value for key %s: expected %d, got %v, %d", item.Key, item.Value, ok, val)
		}
	}

	// Test Keys()
	keys := make([]string, 0)
	for k := range sm.Keys() {
		keys = append(keys, k)
	}

	expectedKeys := []string{"a", "b", "c", "d", "e"}
	if len(keys) != len(expectedKeys) {
		t.Errorf("Keys() returned wrong number of keys: expected %d, got %d", len(expectedKeys), len(keys))
	} else {
		for i, k := range keys {
			if k != expectedKeys[i] {
				t.Errorf("Keys() returned wrong order at index %d: expected %s, got %s", i, expectedKeys[i], k)
			}
		}
	}

	// Test Values()
	values := make([]int, 0)
	for v := range sm.Values() {
		values = append(values, v)
	}

	expectedValues := []int{1, 2, 3, 4, 5}
	if len(values) != len(expectedValues) {
		t.Errorf("Values() returned wrong number of values: expected %d, got %d", len(expectedValues), len(values))
	} else {
		for i, v := range values {
			if v != expectedValues[i] {
				t.Errorf("Values() returned wrong order at index %d: expected %d, got %d", i, expectedValues[i], v)
			}
		}
	}

	// Test ValuesMut() and IterMut()
	for v := range sm.ValuesMut() {
		*v *= 2 // Multiply each value by 2
	}

	// Verify values have been modified
	for _, item := range data {
		if val, found := sm.Get(item.Key); !found || val != item.Value*2 {
			t.Errorf("ValuesMut() failed to modify value for key %s: expected %d, got %v, %d", item.Key, item.Value*2, found, val)
		}
	}

	// Modify back with IterMut
	for k, v := range sm.IterMut() {
		for _, item := range data {
			if k == item.Key {
				*v = item.Value // Restore original value
				break
			}
		}
	}

	// Verify values have been restored
	for _, item := range data {
		if val, found := sm.Get(item.Key); !found || val != item.Value {
			t.Errorf("IterMut() failed to restore value for key %s: expected %d, got %v, %d", item.Key, item.Value, found, val)
		}
	}
}

// Test range query functionality
func TestSkipMapRange(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert some key-value pairs
	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
		"f": 6,
		"g": 7,
	}

	for k, v := range data {
		sm.Insert(k, v)
	}

	// Test full range [nil, nil)
	fullRange := make(map[string]int)
	for k, v := range sm.Range(nil, nil) {
		fullRange[k] = v
	}

	if len(fullRange) != len(data) {
		t.Errorf("Range(nil, nil) returned wrong number of items: expected %d, got %d", len(data), len(fullRange))
	}

	// Test left-closed, right-open range ["b", "f")
	lowerB := "b"
	upperB := "f"
	expected := map[string]int{"b": 2, "c": 3, "d": 4, "e": 5}
	rangeData := make(map[string]int)
	for k, v := range sm.Range(&lowerB, &upperB) {
		rangeData[k] = v
	}

	if len(rangeData) != len(expected) {
		t.Errorf("Range(\"b\", \"f\") returned wrong number of items: expected %d, got %d", len(expected), len(rangeData))
	}

	for k, v := range expected {
		if val, ok := rangeData[k]; !ok || val != v {
			t.Errorf("Range(\"b\", \"f\") missing or wrong value for key %s: expected %d, got %v, %d", k, v, ok, val)
		}
	}

	// Test range with only lower bound ["d", nil)
	lowerD := "d"
	expectedLower := map[string]int{"d": 4, "e": 5, "f": 6, "g": 7}
	lowerRange := make(map[string]int)
	for k, v := range sm.Range(&lowerD, nil) {
		lowerRange[k] = v
	}

	if len(lowerRange) != len(expectedLower) {
		t.Errorf("Range(\"d\", nil) returned wrong number of items: expected %d, got %d", len(expectedLower), len(lowerRange))
	}

	// Test range with only upper bound [nil, "c")
	upperC := "c"
	expectedUpper := map[string]int{"a": 1, "b": 2}
	upperRange := make(map[string]int)
	for k, v := range sm.Range(nil, &upperC) {
		upperRange[k] = v
	}

	if len(upperRange) != len(expectedUpper) {
		t.Errorf("Range(nil, \"c\") returned wrong number of items: expected %d, got %d", len(expectedUpper), len(upperRange))
	}
}

// Test Entry API functionality
func TestSkipMapEntryAPI(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Test OrInsert - key doesn't exist
	valPtr1 := sm.Entry("apple").OrInsert(5)
	if *valPtr1 != 5 {
		t.Errorf("OrInsert for new key should return 5, got %d", *valPtr1)
	}

	// Verify value has been inserted
	val, found := sm.Get("apple")
	if !found || val != 5 {
		t.Errorf("Get after OrInsert should return 5, true, got %d, %v", val, found)
	}

	// Test OrInsert - key exists
	valPtr2 := sm.Entry("apple").OrInsert(10)
	if *valPtr2 != 5 || valPtr2 != valPtr1 {
		t.Errorf("OrInsert for existing key should return existing value 5, got %d", *valPtr2)
	}

	// Verify value hasn't been modified
	val, found = sm.Get("apple")
	if !found || val != 5 {
		t.Errorf("Value should remain 5 after OrInsert with existing key, got %d", val)
	}

	// Test OrInsertWith - key doesn't exist
	valPtr3 := sm.Entry("banana").OrInsertWith(func() int { return 20 })
	if *valPtr3 != 20 {
		t.Errorf("OrInsertWith for new key should return 20, got %d", *valPtr3)
	}

	// Verify value has been inserted
	val, found = sm.Get("banana")
	if !found || val != 20 {
		t.Errorf("Get after OrInsertWith should return 20, true, got %d, %v", val, found)
	}

	// Test OrInsertWithKey
	valPtr4 := sm.Entry("cherry").OrInsertWithKey(func(k string) int { return len(k) * 10 })
	if *valPtr4 != 60 {
		t.Errorf("OrInsertWithKey should return 60 (len(\"cherry\") * 10), got %d", *valPtr4)
	}

	// Test AndModify
	sm.Entry("apple").AndModify(func(v *int) { *v *= 2 })
	val, found = sm.Get("apple")
	if !found || val != 10 {
		t.Errorf("Value should be doubled to 10 after AndModify, got %d", val)
	}

	// Test chaining
	sm.Entry("date").AndModify(func(v *int) { *v = 100 }).OrInsert(50)
	val, found = sm.Get("date")
	if !found || val != 50 {
		t.Errorf("Value should be 50 after AndModify+OrInsert for new key, got %d", val)
	}

	// Test Get method
	v, ok := sm.Entry("apple").Get()
	if !ok || *v != 10 {
		t.Errorf("Entry.Get() should return 10, true for existing key, got %v, %v", v, ok)
	}

	v, ok = sm.Entry("grape").Get()
	if ok || v != nil {
		t.Errorf("Entry.Get() should return nil, false for non-existing key, got %v, %v", v, ok)
	}

	// Test Insert method
	oldVal, updated := sm.Entry("apple").Insert(100)
	if oldVal != 10 || !updated {
		t.Errorf("Entry.Insert() should return 10, true for existing key, got %d, %v", oldVal, updated)
	}

	oldVal, updated = sm.Entry("grape").Insert(200)
	if oldVal != 0 || updated {
		t.Errorf("Entry.Insert() should return 0, false for new key, got %d, %v", oldVal, updated)
	}
}

// Test ordering operations (First, Last, PopFirst, PopLast)
func TestSkipMapOrderOperations(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert some unordered key-value pairs
	data := map[string]int{
		"d": 4,
		"a": 1,
		"c": 3,
		"b": 2,
		"e": 5,
	}

	for k, v := range data {
		sm.Insert(k, v)
	}

	// Test First
	key, val, found := sm.First()
	if !found || key != "a" || val != 1 {
		t.Errorf("First() should return a, 1, true, got %s, %d, %v", key, val, found)
	}

	// Test Last
	key, val, found = sm.Last()
	if !found || key != "e" || val != 5 {
		t.Errorf("Last() should return e, 5, true, got %s, %d, %v", key, val, found)
	}

	// Test PopFirst
	key, val, found = sm.PopFirst()
	if !found || key != "a" || val != 1 {
		t.Errorf("PopFirst() should return a, 1, true, got %s, %d, %v", key, val, found)
	}

	// Verify element has been removed
	if sm.ContainsKey("a") {
		t.Error("Key 'a' should be removed after PopFirst()")
	}

	// Test First again, should return the next element
	key, val, found = sm.First()
	if !found || key != "b" || val != 2 {
		t.Errorf("First() after PopFirst() should return b, 2, true, got %s, %d, %v", key, val, found)
	}

	// Test PopLast
	key, val, found = sm.PopLast()
	if !found || key != "e" || val != 5 {
		t.Errorf("PopLast() should return e, 5, true, got %s, %d, %v", key, val, found)
	}

	// Verify element has been removed
	if sm.ContainsKey("e") {
		t.Error("Key 'e' should be removed after PopLast()")
	}

	// Test Last again, should return the previous element
	key, val, found = sm.Last()
	if !found || key != "d" || val != 4 {
		t.Errorf("Last() after PopLast() should return d, 4, true, got %s, %d, %v", key, val, found)
	}

	// Clear the map and test empty state
	sm.Clear()
	key, val, found = sm.First()
	if found {
		t.Errorf("First() on empty map should return false, got %v with %s, %d", found, key, val)
	}

	key, val, found = sm.Last()
	if found {
		t.Errorf("Last() on empty map should return false, got %v with %s, %d", found, key, val)
	}

	key, val, found = sm.PopFirst()
	if found {
		t.Errorf("PopFirst() on empty map should return false, got %v with %s, %d", found, key, val)
	}

	key, val, found = sm.PopLast()
	if found {
		t.Errorf("PopLast() on empty map should return false, got %v with %s, %d", found, key, val)
	}
}

// Test Clone functionality
func TestSkipMapClone(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert some key-value pairs
	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	for k, v := range data {
		sm.Insert(k, v)
	}

	// Create clone
	clone := sm.Clone()

	// Verify clone contains all original elements
	for k, v := range data {
		cloneVal, found := clone.Get(k)
		if !found || cloneVal != v {
			t.Errorf("Clone missing or wrong value for key %s: expected %d, got %v, %d", k, v, found, cloneVal)
		}
	}

	// Verify clone length
	if clone.Len() != sm.Len() {
		t.Errorf("Clone length should be %d, got %d", sm.Len(), clone.Len())
	}

	// Modifying clone should not affect original map
	clone.Insert("d", 4)
	clone.Insert("a", 100)

	// Verify original map hasn't been modified
	val, found := sm.Get("d")
	if found || val != 0 {
		t.Errorf("Original map should not have key 'd' after adding to clone, got %v, %d", found, val)
	}

	val, found = sm.Get("a")
	if !found || val != 1 {
		t.Errorf("Original map value for 'a' should remain 1 after modifying clone, got %v, %d", found, val)
	}
}

// Test Extend functionality
func TestSkipMapExtend(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert initial key-value pairs
	sm.Insert("a", 1)
	sm.Insert("b", 2)

	// Create iterator for extension
	extendData := []struct {
		K string
		V int
	}{
		{"c", 3},
		{"d", 4},
		{"a", 100}, // This key already exists, should update the value
	}

	extendIter := func(yield func(string, int) bool) {
		for _, item := range extendData {
			if !yield(item.K, item.V) {
				return
			}
		}
	}

	// Perform extension
	sm.Extend(extendIter)

	// Verify all key-value pairs
	expected := map[string]int{
		"a": 100,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	for k, v := range expected {
		val, found := sm.Get(k)
		if !found || val != v {
			t.Errorf("Get(%s) after Extend should return %d, true, got %d, %v", k, v, val, found)
		}
	}
}

// Test GetMut functionality
func TestSkipMapGetMut(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert some key-value pairs
	sm.Insert("apple", 5)
	sm.Insert("banana", 10)

	// Test GetMut - key exists
	valPtr, found := sm.GetMut("apple")
	if !found || *valPtr != 5 {
		t.Errorf("GetMut(\"apple\") should return pointer to 5, true, got %v, %v", valPtr, found)
	}

	// Modify value
	*valPtr = 15

	// Verify value has been modified
	val, found := sm.Get("apple")
	if !found || val != 15 {
		t.Errorf("Value should be updated to 15, got %d, %v", val, found)
	}

	// Test GetMut - key doesn't exist
	valPtr, found = sm.GetMut("cherry")
	if found || valPtr != nil {
		t.Errorf("GetMut(\"cherry\") should return nil, false, got %v, %v", valPtr, found)
	}
}

// Test GetKeyValue functionality
func TestSkipMapGetKeyValue(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert some key-value pairs
	sm.Insert("apple", 5)

	// Test GetKeyValue - key exists
	key, val, found := sm.GetKeyValue("apple")
	if !found || key != "apple" || val != 5 {
		t.Errorf("GetKeyValue(\"apple\") should return apple, 5, true, got %s, %d, %v", key, val, found)
	}

	// Test GetKeyValue - key doesn't exist
	key, val, found = sm.GetKeyValue("cherry")
	if found || key != "" || val != 0 {
		t.Errorf("GetKeyValue(\"cherry\") should return empty string, 0, false, got %s, %d, %v", key, val, found)
	}
}

// Test custom comparator
func TestSkipMapCustomComparator(t *testing.T) {
	// Create a SkipMap with custom comparator that sorts by key length
	customCmp := func(a, b string) int {
		if len(a) != len(b) {
			return len(a) - len(b)
		}
		return 0
	}

	sm := New[string, int](customCmp)

	// Insert keys of different lengths
	sm.Insert("a", 1)
	sm.Insert("ab", 2)
	sm.Insert("abc", 3)
	sm.Insert("abcd", 4)

	// Verify order is sorted by length
	keys := make([]string, 0)
	for k := range sm.Keys() {
		keys = append(keys, k)
	}

	expected := []string{"a", "ab", "abc", "abcd"}
	if len(keys) != len(expected) {
		t.Errorf("Keys length mismatch: expected %d, got %d", len(expected), len(keys))
	} else {
		for i, k := range keys {
			if k != expected[i] {
				t.Errorf("Key order mismatch at index %d: expected %s, got %s", i, expected[i], k)
			}
		}
	}

	// Test that keys with same length are considered equal
	oldVal, updated := sm.Insert("xy", 100)
	if oldVal != 2 || !updated {
		t.Errorf("Inserting key with same length should update, got %d, %v", oldVal, updated)
	}

	// Verify value has been updated
	val, found := sm.Get("xy")
	if !found || val != 100 {
		t.Errorf("Get(\"xy\") should return 100, true, got %d, %v", val, found)
	}
}

// Test performance and correctness under heavy load
func TestSkipMapLargeLoad(t *testing.T) {
	// Only run large test when test flag is enabled
	if testing.Short() {
		t.Skip("Skipping large load test in short mode")
	}

	const size = 10000
	sm := NewOrdered[int, int]()

	// Insert random data
	rand.Seed(time.Now().UnixNano())
	data := make(map[int]int)
	for i := 0; i < size; i++ {
		k := rand.Intn(size * 10)
		v := rand.Intn(size * 10)
		data[k] = v
		sm.Insert(k, v)
	}

	// Verify all data has been correctly inserted
	for k, v := range data {
		val, found := sm.Get(k)
		if !found || val != v {
			t.Errorf("Get(%d) should return %d, true, got %d, %v", k, v, val, found)
		}
	}

	// Verify length
	if sm.Len() != len(data) {
		t.Errorf("Length should be %d, got %d", len(data), sm.Len())
	}

	// Test iteration order (should be sorted by keys)
	keys := make([]int, 0, len(data))
	for k := range sm.Keys() {
		keys = append(keys, k)
	}

	// Verify keys are ordered
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("Keys not ordered: keys[%d]=%d < keys[%d]=%d", i, keys[i], i-1, keys[i-1])
		}
	}

	// Test deleting half of the data
	deleteCount := 0
	for k := range data {
		if deleteCount < len(data)/2 {
			sm.Remove(k)
			delete(data, k) // Also delete the key from data map
			deleteCount++
		} else {
			break
		}
	}

	// Verify data after deletion
	for k, v := range data {
		val, found := sm.Get(k)
		// All remaining keys should exist
		if !found || val != v {
			t.Errorf("Get(%d) should return %d, true, got %d, %v", k, v, val, found)
		}
	}
}

// Test using custom type as key
func TestSkipMapCustomKeyType(t *testing.T) {
	// Define a custom type
	type Person struct {
		Name string
		Age  int
	}

	// Define comparison function
	personCmp := func(a, b Person) int {
		if a.Age != b.Age {
			return a.Age - b.Age
		}
		return 0
	}

	sm := New[Person, string](personCmp)

	// Insert data
	people := []struct {
		Key   Person
		Value string
	}{
		{Person{"Alice", 30}, "Engineer"},
		{Person{"Bob", 25}, "Designer"},
		{Person{"Charlie", 35}, "Manager"},
	}

	for _, p := range people {
		sm.Insert(p.Key, p.Value)
	}

	// Verify data
	for _, p := range people {
		val, found := sm.Get(p.Key)
		if !found || val != p.Value {
			t.Errorf("Get(%v) should return %s, true, got %s, %v", p.Key, p.Value, val, found)
		}
	}

	// Test that people with same age are considered equal
	oldVal, updated := sm.Insert(Person{"David", 25}, "Developer")
	if oldVal != "Designer" || !updated {
		t.Errorf("Inserting person with same age should update, got %s, %v", oldVal, updated)
	}

	// Verify update
	val, found := sm.Get(Person{"David", 25})
	if !found || val != "Developer" {
		t.Errorf("Get updated person should return Developer, true, got %s, %v", val, found)
	}

	// Test iteration order (should be sorted by age)
	ages := make([]int, 0)
	for k := range sm.Keys() {
		ages = append(ages, k.Age)
	}

	expectedAges := []int{25, 30, 35}
	for i, age := range ages {
		if age != expectedAges[i] {
			t.Errorf("Age order mismatch at index %d: expected %d, got %d", i, expectedAges[i], age)
		}
	}
}

// Test mutable iterator modification functionality
func TestSkipMapMutIterators(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert some data
	sm.Insert("a", 1)
	sm.Insert("b", 2)
	sm.Insert("c", 3)

	// Test ValuesMut modification
	for v := range sm.ValuesMut() {
		*v *= 10
	}

	// Verify modification
	expected := map[string]int{"a": 10, "b": 20, "c": 30}
	for k, expVal := range expected {
		val, found := sm.Get(k)
		if !found || val != expVal {
			t.Errorf("ValuesMut: Get(%s) should return %d, true, got %d, %v", k, expVal, val, found)
		}
	}

	// Test IterMut modification
	for _, v := range sm.IterMut() {
		*v += 1
	}

	// Verify modification
	expected = map[string]int{"a": 11, "b": 21, "c": 31}
	for k, expVal := range expected {
		val, found := sm.Get(k)
		if !found || val != expVal {
			t.Errorf("IterMut: Get(%s) should return %d, true, got %d, %v", k, expVal, val, found)
		}
	}
}

// Test modifying map during iteration
func TestSkipMapModifyDuringIteration(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert some data
	for i := 0; i < 10; i++ {
		sm.Insert(fmt.Sprintf("key%d", i), i)
	}

	// Delete some keys during iteration
	toDelete := []string{"key2", "key5", "key8"}
	for k := range sm.Keys() {
		for _, delKey := range toDelete {
			if k == delKey {
				sm.Remove(k)
				break
			}
		}
	}

	// Verify deletion
	for _, delKey := range toDelete {
		_, found := sm.Get(delKey)
		if found {
			t.Errorf("Key %s should be deleted", delKey)
		}
	}

	// Verify other keys still exist
	expectedKeys := []string{"key0", "key1", "key3", "key4", "key6", "key7", "key9"}
	for _, key := range expectedKeys {
		_, found := sm.Get(key)
		if !found {
			t.Errorf("Key %s should still exist", key)
		}
	}
}

// New test case: Test GetComparator function
func TestSkipMapGetComparator(t *testing.T) {
	// Test comparator for ordered type
	sm1 := NewOrdered[string, int]()
	comparator1 := sm1.GetComparator()
	if comparator1 == nil {
		t.Error("GetComparator should return a non-nil function for ordered type")
	}

	// Test custom comparator
	customCmp := func(a, b string) int {
		if len(a) != len(b) {
			return len(a) - len(b)
		}
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}
	sm2 := New[string, int](customCmp)
	comparator2 := sm2.GetComparator()
	if comparator2 == nil {
		t.Error("GetComparator should return a non-nil function for custom comparator")
	}

	// Verify comparator functionality
	result := comparator2("a", "b")
	if result >= 0 {
		t.Error("Custom comparator should return negative value for 'a' < 'b'")
	}

	// Test comparator consistency
	result1 := comparator2("test", "test")
	if result1 != 0 {
		t.Error("Custom comparator should return 0 for equal strings")
	}

	result2 := comparator2("zzz", "aaa")
	if result2 <= 0 {
		t.Error("Custom comparator should return positive value for 'zzz' > 'aaa'")
	}
}

// New test case: Test getNode function in various scenarios
func TestSkipMapGetNode(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Test getNode for empty map
	node, found := sm.getNode("nonexistent")
	if found || node != nil {
		t.Error("getNode should return nil, false for empty map")
	}

	// Test after inserting elements
	sm.Insert("key1", 100)
	sm.Insert("key2", 200)

	// Test existing key
	node, found = sm.getNode("key1")
	if !found || node == nil || node.Value != 100 {
		t.Error("getNode should return valid node for existing key")
	}

	// Test non-existent key
	node, found = sm.getNode("nonexistent")
	if found || node != nil {
		t.Error("getNode should return nil, false for non-existing key")
	}
}

// New test case: Test edge cases for Entry's OrInsertWith and OrInsertWithKey functions
func TestSkipMapEntryOrInsertWithEdgeCases(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Test OrInsertWith - key doesn't exist, using function that returns fixed value
	callCount := 0
	valPtr1 := sm.Entry("apple").OrInsertWith(func() int {
		callCount++
		return 42
	})

	if *valPtr1 != 42 {
		t.Errorf("OrInsertWith should return 42, got %d", *valPtr1)
	}

	if callCount != 1 {
		t.Errorf("Function should be called once, called %d times", callCount)
	}

	// Test OrInsertWith - key exists, function should not be called
	callCount = 0
	valPtr2 := sm.Entry("apple").OrInsertWith(func() int {
		callCount++
		return 100
	})

	if *valPtr2 != 42 {
		t.Errorf("OrInsertWith for existing key should return existing value 42, got %d", *valPtr2)
	}

	if callCount != 0 {
		t.Errorf("Function should not be called for existing key, called %d times", callCount)
	}

	// Test OrInsertWithKey - key doesn't exist, using key-based function
	valPtr3 := sm.Entry("banana").OrInsertWithKey(func(k string) int {
		return len(k) * 10
	})

	if *valPtr3 != 60 {
		t.Errorf("OrInsertWithKey should return 60 (len(\"banana\") * 10), got %d", *valPtr3)
	}

	// Test OrInsertWithKey - key exists, function should not be called
	valPtr4 := sm.Entry("banana").OrInsertWithKey(func(k string) int {
		return 1000
	})

	if *valPtr4 != 60 {
		t.Errorf("OrInsertWithKey for existing key should return existing value 60, got %d", *valPtr4)
	}

	// Test OrInsertWithKey with complex key-based logic
	valPtr5 := sm.Entry("cherry").OrInsertWithKey(func(k string) int {
		// Calculate value based on key length and characters
		result := len(k)
		for _, c := range k {
			result += int(c)
		}
		return result
	})

	expected := len("cherry")
	for _, c := range "cherry" {
		expected += int(c)
	}

	if *valPtr5 != expected {
		t.Errorf("OrInsertWithKey should return %d, got %d", expected, *valPtr5)
	}
}

// New test case: Test edge cases for iterators
func TestSkipMapIteratorsEdgeCases(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Test iterators on empty map
	count := 0
	for range sm.Keys() {
		count++
	}
	if count != 0 {
		t.Errorf("Keys() on empty map should yield 0 items, got %d", count)
	}

	count = 0
	for range sm.Values() {
		count++
	}
	if count != 0 {
		t.Errorf("Values() on empty map should yield 0 items, got %d", count)
	}

	count = 0
	for range sm.ValuesMut() {
		count++
	}
	if count != 0 {
		t.Errorf("ValuesMut() on empty map should yield 0 items, got %d", count)
	}

	count = 0
	for range sm.Iter() {
		count++
	}
	if count != 0 {
		t.Errorf("Iter() on empty map should yield 0 items, got %d", count)
	}

	count = 0
	for range sm.IterMut() {
		count++
	}
	if count != 0 {
		t.Errorf("IterMut() on empty map should yield 0 items, got %d", count)
	}

	// Insert a single element for testing
	sm.Insert("single", 1)

	// Test iterators on single-element map
	keys := make([]string, 0)
	for k := range sm.Keys() {
		keys = append(keys, k)
	}
	if len(keys) != 1 || keys[0] != "single" {
		t.Errorf("Keys() on single-element map should yield ['single'], got %v", keys)
	}

	values := make([]int, 0)
	for v := range sm.Values() {
		values = append(values, v)
	}
	if len(values) != 1 || values[0] != 1 {
		t.Errorf("Values() on single-element map should yield [1], got %v", values)
	}

	// Test ValuesMut on single element
	for v := range sm.ValuesMut() {
		*v = 42
	}

	val, found := sm.Get("single")
	if !found || val != 42 {
		t.Errorf("ValuesMut should modify value to 42, got %d", val)
	}

	// Test IterMut on single element
	for k, v := range sm.IterMut() {
		if k == "single" {
			*v = 100
		}
	}

	val, found = sm.Get("single")
	if !found || val != 100 {
		t.Errorf("IterMut should modify value to 100, got %d", val)
	}
}

// New test case: Test edge cases for PopLast
func TestSkipMapPopLastEdgeCases(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Test PopLast on empty map
	key, val, found := sm.PopLast()
	if found {
		t.Errorf("PopLast() on empty map should return false, got %v with %s, %d", found, key, val)
	}

	// Test PopLast on single-element map
	sm.Insert("single", 42)
	key, val, found = sm.PopLast()
	if !found || key != "single" || val != 42 {
		t.Errorf("PopLast() on single-element map should return single, 42, true, got %s, %d, %v", key, val, found)
	}

	// Verify map is now empty
	if !sm.IsEmpty() {
		t.Error("Map should be empty after PopLast on single-element map")
	}

	// Test empty map again
	key, val, found = sm.PopLast()
	if found {
		t.Errorf("PopLast() on empty map should return false, got %v with %s, %d", found, key, val)
	}
}

// New test case: Test edge cases for New function
func TestSkipMapNewEdgeCases(t *testing.T) {
	// Test that nil comparator panics
	defer func() {
		if r := recover(); r == nil {
			t.Error("New with nil comparator should panic")
		}
	}()
	_ = New[string, int](nil)
}

// New test case: Test normal cases for New function
func TestSkipMapNewNormalCases(t *testing.T) {
	// Test creating map with custom comparator
	customCmp := func(a, b string) int {
		// Compare by string length, then by lex order when lengths are equal
		if len(a) != len(b) {
			return len(a) - len(b)
		}
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}

	sm := New[string, int](customCmp)
	if sm == nil {
		t.Error("New should create a valid SkipMap")
	}

	// Test insertion and retrieval
	sm.Insert("a", 1)
	sm.Insert("bb", 2)
	sm.Insert("ccc", 3)

	// Verify sorting by length
	keys := make([]string, 0)
	for k := range sm.Keys() {
		keys = append(keys, k)
	}

	expected := []string{"a", "bb", "ccc"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("Keys should be ordered by length, expected %s, got %s", expected[i], k)
		}
	}
}

// New test case: Test early termination of iterators
func TestSkipMapIteratorsEarlyTermination(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert multiple elements
	for i := 1; i <= 10; i++ {
		sm.Insert(fmt.Sprintf("key%d", i), i)
	}

	// Test early termination of Keys iterator
	count := 0
	for range sm.Keys() {
		count++
		if count == 5 {
			break
		}
	}
	if count != 5 {
		t.Errorf("Keys iterator should yield 5 elements before break, got %d", count)
	}

	// Test early termination of Values iterator
	count = 0
	for range sm.Values() {
		count++
		if count == 3 {
			break
		}
	}
	if count != 3 {
		t.Errorf("Values iterator should yield 3 elements before break, got %d", count)
	}

	// Test early termination of ValuesMut iterator
	count = 0
	for range sm.ValuesMut() {
		count++
		if count == 4 {
			break
		}
	}
	if count != 4 {
		t.Errorf("ValuesMut iterator should yield 4 elements before break, got %d", count)
	}

	// Test early termination of Iter iterator
	count = 0
	for range sm.Iter() {
		count++
		if count == 6 {
			break
		}
	}
	if count != 6 {
		t.Errorf("Iter iterator should yield 6 elements before break, got %d", count)
	}

	// Test early termination of IterMut iterator
	count = 0
	for range sm.IterMut() {
		count++
		if count == 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("IterMut iterator should yield 2 elements before break, got %d", count)
	}

	// Verify map structure is not corrupted
	if sm.Len() != 10 {
		t.Errorf("Map length should remain 10 after early termination, got %d", sm.Len())
	}
}

// New test case: Test more branches of Range function
func TestSkipMapRangeMoreBranches(t *testing.T) {
	sm := NewOrdered[string, int]()

	// Insert data
	for i := 1; i <= 20; i++ {
		sm.Insert(fmt.Sprintf("key%d", i), i)
	}

	// Test early termination of range query
	lower := "key5"
	upper := "key8"
	count := 0
	for range sm.Range(&lower, &upper) {
		count++
		if count == 3 {
			break
		}
	}
	if count != 3 {
		t.Errorf("Range iterator should yield 3 elements before break, got %d", count)
	}

	// Test case with only lower bound
	lower = "key10"
	count = 0
	for range sm.Range(&lower, nil) {
		count++
		if count == 5 {
			break
		}
	}
	if count != 5 {
		t.Errorf("Range with lower bound iterator should yield 5 elements before break, got %d", count)
	}

	// Test case with only upper bound
	upper = "key4"
	count = 0
	for range sm.Range(nil, &upper) {
		count++
		if count == 3 {
			break
		}
	}
	if count != 3 {
		t.Errorf("Range with upper bound iterator should yield 3 elements before break, got %d", count)
	}

	// Verify map structure is not corrupted
	if sm.Len() != 20 {
		t.Errorf("Map length should remain 20 after range queries, got %d", sm.Len())
	}

	// Test edge case: lowerBound > upperBound
	lower = "key15"
	upper = "key10"
	count = 0
	for range sm.Range(&lower, &upper) {
		count++
	}
	if count != 0 {
		t.Errorf("Range with lower > upper should yield 0 elements, got %d", count)
	}

	// Test edge case: lowerBound equals a key
	lower = "key10"
	upper = "key12"
	values := make([]int, 0)
	for k, v := range sm.Range(&lower, &upper) {
		if k == "key10" {
			values = append(values, v)
		}
	}
	if len(values) != 1 || values[0] != 10 {
		t.Errorf("Range should include key10 with value 10, got %v", values)
	}
}
