package skipset

import (
	"fmt"
	"github.com/go-board/ds/bound"
	"math/rand"
	"sort"
	"testing"
	"time"
)

// Test basic insertion, deletion, and containment check functionality
func TestSkipSetBasicOperations(t *testing.T) {
	// Create a new SkipSet
	ss := NewOrdered[string]()

	// Test insertion
	added := ss.Insert("apple")
	if !added {
		t.Error("Insert(\"apple\") should return true")
	}

	// Test duplicate insertion
	added = ss.Insert("apple")
	if added {
		t.Error("Insert(\"apple\") again should return false")
	}

	// Test containment check
	if !ss.Contains("apple") {
		t.Error("Contains(\"apple\") should return true")
	}
	if ss.Contains("orange") {
		t.Error("Contains(\"orange\") should return false")
	}

	// Test deletion
	removed := ss.Remove("apple")
	if !removed {
		t.Error("Remove(\"apple\") should return true")
	}

	// Test deletion of non-existent element
	removed = ss.Remove("orange")
	if removed {
		t.Error("Remove(\"orange\") should return false")
	}

	// Verify element has been removed
	if ss.Contains("apple") {
		t.Error("Contains(\"apple\") should return false after removal")
	}
}

// Test length and empty check functionality
func TestSkipSetLenAndEmpty(t *testing.T) {
	ss := NewOrdered[string]()

	// Initial state should be empty
	if !ss.IsEmpty() {
		t.Error("New skip set should be empty")
	}
	if ss.Len() != 0 {
		t.Errorf("New skip set length should be 0, got %d", ss.Len())
	}

	// Length should increase after inserting elements
	ss.Insert("a")
	if ss.IsEmpty() {
		t.Error("Skip set should not be empty after insertion")
	}
	if ss.Len() != 1 {
		t.Errorf("Skip set length should be 1 after insertion, got %d", ss.Len())
	}

	// Insert more elements
	ss.Insert("b")
	ss.Insert("c")
	if ss.Len() != 3 {
		t.Errorf("Skip set length should be 3 after insertions, got %d", ss.Len())
	}

	// Length should decrease after removing elements
	ss.Remove("b")
	if ss.Len() != 2 {
		t.Errorf("Skip set length should be 2 after removal, got %d", ss.Len())
	}

	// Clear the set
	ss.Clear()
	if !ss.IsEmpty() {
		t.Error("Skip set should be empty after Clear()")
	}
	if ss.Len() != 0 {
		t.Errorf("Skip set length should be 0 after Clear(), got %d", ss.Len())
	}
}

// Test iterator functionality
func TestSkipSetIterator(t *testing.T) {
	ss := NewOrdered[string]()

	// Insert some ordered elements
	data := []string{"a", "b", "c", "d", "e"}

	for _, item := range data {
		ss.Insert(item)
	}

	// Test IterAsc()
	iterData := make([]string, 0)
	for val := range ss.IterAsc() {
		iterData = append(iterData, val)
	}

	// Verify iteration order
	sort.Strings(data) // Ensure data is sorted
	if len(iterData) != len(data) {
		t.Errorf("IterAsc() returned wrong number of items: expected %d, got %d", len(data), len(iterData))
	} else {
		for i, val := range iterData {
			if val != data[i] {
				t.Errorf("IterAsc() returned wrong value at index %d: expected %s, got %s", i, data[i], val)
			}
		}
	}
}

// Test range query functionality
func TestSkipSetRangeAsc(t *testing.T) {
	ss := NewOrdered[string]()

	// Insert some elements
	data := []string{"a", "b", "c", "d", "e", "f", "g"}
	for _, item := range data {
		ss.Insert(item)
	}

	// Test full range [nil, nil)
	fullRange := make([]string, 0)
	for val := range ss.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[string](), bound.NewUnbounded[string]())) {
		fullRange = append(fullRange, val)
	}

	if len(fullRange) != len(data) {
		t.Errorf("RangeAsc(nil, nil) returned wrong number of items: expected %d, got %d", len(data), len(fullRange))
	}

	// Test left-closed right-open range ["b", "f")
	lowerB := "b"
	upperB := "f"
	expected := []string{"b", "c", "d", "e"}
	rangeData := make([]string, 0)
	for val := range ss.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lowerB), bound.NewExcluded(upperB))) {
		rangeData = append(rangeData, val)
	}

	if len(rangeData) != len(expected) {
		t.Errorf("RangeAsc(\"b\", \"f\") returned wrong number of items: expected %d, got %d", len(expected), len(rangeData))
	} else {
		for i, val := range rangeData {
			if val != expected[i] {
				t.Errorf("RangeAsc(\"b\", \"f\") returned wrong value at index %d: expected %s, got %s", i, expected[i], val)
			}
		}
	}

	// Test only lower bound ["d", nil)
	lowerD := "d"
	expectedLower := []string{"d", "e", "f", "g"}
	lowerRange := make([]string, 0)
	for val := range ss.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lowerD), bound.NewUnbounded[string]())) {
		lowerRange = append(lowerRange, val)
	}

	if len(lowerRange) != len(expectedLower) {
		t.Errorf("RangeAsc(\"d\", nil) returned wrong number of items: expected %d, got %d", len(expectedLower), len(lowerRange))
	}

	// Test only upper bound [nil, "c")
	upperC := "c"
	expectedUpper := []string{"a", "b"}
	upperRange := make([]string, 0)
	for val := range ss.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[string](), bound.NewExcluded(upperC))) {
		upperRange = append(upperRange, val)
	}

	if len(upperRange) != len(expectedUpper) {
		t.Errorf("RangeAsc(nil, \"c\") returned wrong number of items: expected %d, got %d", len(expectedUpper), len(upperRange))
	}
}

// Test order operations functionality (First, Last, PopFirst, PopLast)
func TestSkipSetOrderOperations(t *testing.T) {
	ss := NewOrdered[string]()

	// Insert some unordered elements
	data := []string{"d", "a", "c", "b", "e"}
	for _, item := range data {
		ss.Insert(item)
	}

	// Test First
	val, found := ss.First()
	if !found || val != "a" {
		t.Errorf("First() should return a, true, got %s, %v", val, found)
	}

	// Test Last
	val, found = ss.Last()
	if !found || val != "e" {
		t.Errorf("Last() should return e, true, got %s, %v", val, found)
	}

	// Test PopFirst
	val, found = ss.PopFirst()
	if !found || val != "a" {
		t.Errorf("PopFirst() should return a, true, got %s, %v", val, found)
	}

	// Verify element has been removed
	if ss.Contains("a") {
		t.Error("Element 'a' should be removed after PopFirst()")
	}

	// Test First again, should return the next element
	val, found = ss.First()
	if !found || val != "b" {
		t.Errorf("First() after PopFirst() should return b, true, got %s, %v", val, found)
	}

	// Test PopLast
	val, found = ss.PopLast()
	if !found || val != "e" {
		t.Errorf("PopLast() should return e, true, got %s, %v", val, found)
	}

	// Verify element has been removed
	if ss.Contains("e") {
		t.Error("Element 'e' should be removed after PopLast()")
	}

	// Test Last again, should return the previous element
	val, found = ss.Last()
	if !found || val != "d" {
		t.Errorf("Last() after PopLast() should return d, true, got %s, %v", val, found)
	}

	// Clear the set and test empty state
	ss.Clear()
	val, found = ss.First()
	if found {
		t.Errorf("First() on empty set should return false, got %v with %s", found, val)
	}

	val, found = ss.Last()
	if found {
		t.Errorf("Last() on empty set should return false, got %v with %s", found, val)
	}

	val, found = ss.PopFirst()
	if found {
		t.Errorf("PopFirst() on empty set should return false, got %v with %s", found, val)
	}

	val, found = ss.PopLast()
	if found {
		t.Errorf("PopLast() on empty set should return false, got %v with %s", found, val)
	}
}

// Test clone functionality
func TestSkipSetClone(t *testing.T) {
	ss := NewOrdered[string]()

	// Insert some elements
	data := []string{"a", "b", "c", "d", "e"}
	for _, item := range data {
		ss.Insert(item)
	}

	// create clone
	clone := ss.Clone()

	// verify clone contains all original elements
	for _, item := range data {
		if !clone.Contains(item) {
			t.Errorf("Clone missing element %s", item)
		}
	}

	// verify clone length
	if clone.Len() != ss.Len() {
		t.Errorf("Clone length should be %d, got %d", ss.Len(), clone.Len())
	}

	// modifying clone should not affect original set
	clone.Insert("f")
	clone.Remove("a")

	// verify original set was not modified
	if !ss.Contains("a") {
		t.Error("Original set should still contain 'a' after removing from clone")
	}
	if ss.Contains("f") {
		t.Error("Original set should not contain 'f' after adding to clone")
	}
}

// Test extend functionality
func TestSkipSetExtend(t *testing.T) {
	ss := NewOrdered[string]()

	// Insert initial elements
	ss.Insert("a")
	ss.Insert("b")

	// create iterator for extend
	extendData := []string{"c", "d", "a", "e"}
	extendIter := func(yield func(string) bool) {
		for _, item := range extendData {
			if !yield(item) {
				return
			}
		}
	}

	// perform extend
	ss.Extend(extendIter)

	// verify all elements were added
	expected := []string{"a", "b", "c", "d", "e"}
	for _, item := range expected {
		if !ss.Contains(item) {
			t.Errorf("Set should contain %s after Extend", item)
		}
	}

	// Verify length
	if ss.Len() != len(expected) {
		t.Errorf("Set length should be %d after Extend, got %d", len(expected), ss.Len())
	}
}

// Test set operation: union
func TestSkipSetUnion(t *testing.T) {
	// Create two sets
	set1 := NewOrdered[string]()
	set2 := NewOrdered[string]()

	// Insert elements
	data1 := []string{"a", "b", "c", "d"}
	data2 := []string{"c", "d", "e", "f"}

	for _, item := range data1 {
		set1.Insert(item)
	}

	for _, item := range data2 {
		set2.Insert(item)
	}

	// Calculate union
	union := set1.Union(set2)

	// Verify union contains all elements
	expected := []string{"a", "b", "c", "d", "e", "f"}
	for _, item := range expected {
		if !union.Contains(item) {
			t.Errorf("Union should contain %s", item)
		}
	}

	// verify length
	if union.Len() != len(expected) {
		t.Errorf("Union length should be %d, got %d", len(expected), union.Len())
	}

	// verify original set was not modified
	for _, item := range data1 {
		if !set1.Contains(item) {
			t.Error("Original set1 should not be modified")
		}
	}
	for _, item := range data2 {
		if !set2.Contains(item) {
			t.Error("Original set2 should not be modified")
		}
	}
}

// Test set operation: intersection
func TestSkipSetIntersection(t *testing.T) {
	// Create two sets
	set1 := NewOrdered[string]()
	set2 := NewOrdered[string]()

	// Insert elements
	data1 := []string{"a", "b", "c", "d"}
	data2 := []string{"c", "d", "e", "f"}

	for _, item := range data1 {
		set1.Insert(item)
	}

	for _, item := range data2 {
		set2.Insert(item)
	}

	// Calculate intersection
	intersection := set1.Intersection(set2)

	// Verify intersection contains only common elements
	expected := []string{"c", "d"}
	for _, item := range expected {
		if !intersection.Contains(item) {
			t.Errorf("Intersection should contain %s", item)
		}
	}

	// verify does not contain other elements
	nonExpected := []string{"a", "b", "e", "f"}
	for _, item := range nonExpected {
		if intersection.Contains(item) {
			t.Errorf("Intersection should not contain %s", item)
		}
	}

	// verify length
	if intersection.Len() != len(expected) {
		t.Errorf("Intersection length should be %d, got %d", len(expected), intersection.Len())
	}
}

// Test set operation: difference
func TestSkipSetDifference(t *testing.T) {
	// Create two sets
	set1 := NewOrdered[string]()
	set2 := NewOrdered[string]()

	// Insert elements
	data1 := []string{"a", "b", "c", "d"}
	data2 := []string{"c", "d", "e", "f"}

	for _, item := range data1 {
		set1.Insert(item)
	}

	for _, item := range data2 {
		set2.Insert(item)
	}

	// Calculate difference: set1 - set2
	difference := set1.Difference(set2)

	// Verify difference contains only elements in set1 but not in set2
	expected := []string{"a", "b"}
	for _, item := range expected {
		if !difference.Contains(item) {
			t.Errorf("Difference should contain %s", item)
		}
	}

	// verify does not contain other elements
	nonExpected := []string{"c", "d", "e", "f"}
	for _, item := range nonExpected {
		if difference.Contains(item) {
			t.Errorf("Difference should not contain %s", item)
		}
	}

	// verify length
	if difference.Len() != len(expected) {
		t.Errorf("Difference length should be %d, got %d", len(expected), difference.Len())
	}
}

// Test set operation: symmetric difference
func TestSkipSetSymmetricDifference(t *testing.T) {
	// Create two sets
	set1 := NewOrdered[string]()
	set2 := NewOrdered[string]()

	// Insert elements
	data1 := []string{"a", "b", "c", "d"}
	data2 := []string{"c", "d", "e", "f"}

	for _, item := range data1 {
		set1.Insert(item)
	}

	for _, item := range data2 {
		set2.Insert(item)
	}

	// Calculate symmetric difference
	symDiff := set1.SymmetricDifference(set2)

	// Verify symmetric difference contains elements that appear in exactly one of the sets
	expected := []string{"a", "b", "e", "f"}
	for _, item := range expected {
		if !symDiff.Contains(item) {
			t.Errorf("SymmetricDifference should contain %s", item)
		}
	}

	// Verify not containing common elements
	nonExpected := []string{"c", "d"}
	for _, item := range nonExpected {
		if symDiff.Contains(item) {
			t.Errorf("SymmetricDifference should not contain %s", item)
		}
	}

	// Verify length
	if symDiff.Len() != len(expected) {
		t.Errorf("SymmetricDifference length should be %d, got %d", len(expected), symDiff.Len())
	}
}

// Test set relations: subset, superset, and disjoint
func TestSkipSetSetRelations(t *testing.T) {
	// Create several sets for testing relations
	set1 := NewOrdered[string]()
	set2 := NewOrdered[string]()
	set3 := NewOrdered[string]()
	set4 := NewOrdered[string]()

	// Insert elements
	for _, item := range []string{"a", "b", "c", "d", "e"} {
		set1.Insert(item)
	}

	for _, item := range []string{"b", "c", "d"} {
		set2.Insert(item)
	}

	for _, item := range []string{"a", "b", "c", "d", "e"} {
		set3.Insert(item)
	}

	for _, item := range []string{"x", "y", "z"} {
		set4.Insert(item)
	}

	// Test subset relation
	if !set2.IsSubset(set1) {
		t.Error("set2 should be a subset of set1")
	}

	if !set1.IsSubset(set3) {
		t.Error("set1 should be a subset of set3 (equal sets)")
	}

	if set1.IsSubset(set2) {
		t.Error("set1 should not be a subset of set2")
	}

	if set1.IsSubset(set4) {
		t.Error("set1 should not be a subset of set4")
	}

	// Test superset relation
	if !set1.IsSuperset(set2) {
		t.Error("set1 should be a superset of set2")
	}

	if !set3.IsSuperset(set1) {
		t.Error("set3 should be a superset of set1 (equal sets)")
	}

	if set2.IsSuperset(set1) {
		t.Error("set2 should not be a superset of set1")
	}

	if set4.IsSuperset(set1) {
		t.Error("set4 should not be a superset of set1")
	}

	// test equality relationship
	if !set1.Equal(set3) {
		t.Error("set1 should be equal to set3")
	}

	if set1.Equal(set2) {
		t.Error("set1 should not be equal to set2")
	}

	// test disjoint relationship
	if !set1.IsDisjoint(set4) {
		t.Error("set1 and set4 should be disjoint")
	}

	if set1.IsDisjoint(set2) {
		t.Error("set1 and set2 should not be disjoint")
	}
}

// Test custom comparator
func TestSkipSetCustomComparator(t *testing.T) {
	// Create a SkipSet with a custom comparator, sorting by string length
	customCmp := func(a, b string) int {
		if len(a) != len(b) {
			return len(a) - len(b)
		}
		return 0
	}

	ss := New[string](customCmp)

	// Insert strings of different lengths
	ss.Insert("a")
	ss.Insert("ab")
	ss.Insert("abc")
	ss.Insert("abcd")

	// Verify ordering is by length
	values := make([]string, 0)
	for val := range ss.IterAsc() {
		values = append(values, val)
	}

	expected := []string{"a", "ab", "abc", "abcd"}
	if len(values) != len(expected) {
		t.Errorf("Values length mismatch: expected %d, got %d", len(expected), len(values))
	} else {
		for i, val := range values {
			if len(val) != len(expected[i]) {
				t.Errorf("Value order mismatch at index %d: expected length %d, got length %d (value: %s)",
					i, len(expected[i]), len(val), val)
			}
		}
	}

	// Test that strings of the same length are considered equal
	added := ss.Insert("xy")
	if added {
		t.Error("Inserting string with same length should not be added")
	}

	// Verify set size remains the same
	if ss.Len() != 4 {
		t.Errorf("Set length should remain 4, got %d", ss.Len())
	}
}

// Test using custom type as elements
func TestSkipSetCustomElementType(t *testing.T) {
	// Define a custom type
	type Person struct {
		Name string
		Age  int
	}

	// define comparator function
	personCmp := func(a, b Person) int {
		if a.Age != b.Age {
			return a.Age - b.Age
		}
		return 0
	}

	ss := New[Person](personCmp)

	// Insert data
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}

	for _, p := range people {
		ss.Insert(p)
	}

	// Verify data
	for _, p := range people {
		if !ss.Contains(p) {
			t.Errorf("Set should contain %v", p)
		}
	}

	// Test that people with the same age are considered equal
	added := ss.Insert(Person{"David", 25})
	if added {
		t.Error("Inserting person with same age should not be added")
	}

	// Verify set size remains the same
	if ss.Len() != 3 {
		t.Errorf("Set length should remain 3, got %d", ss.Len())
	}

	// Test iteration order (should be by age)
	ages := make([]int, 0)
	for p := range ss.IterAsc() {
		ages = append(ages, p.Age)
	}

	expectedAges := []int{25, 30, 35}
	for i, age := range ages {
		if age != expectedAges[i] {
			t.Errorf("Age order mismatch at index %d: expected %d, got %d", i, expectedAges[i], age)
		}
	}
}

// Test performance and correctness under high load scenarios
func TestSkipSetLargeLoad(t *testing.T) {
	// Only run large tests when the test flag is enabled
	if testing.Short() {
		t.Skip("Skipping large load test in short mode")
	}

	const size = 10000
	ss := NewOrdered[int]()

	// Insert random data
	rand.Seed(time.Now().UnixNano())
	data := make(map[int]bool)
	for i := 0; i < size; i++ {
		k := rand.Intn(size * 10)
		data[k] = true
		ss.Insert(k)
	}

	// Verify all data was inserted correctly
	for k := range data {
		if !ss.Contains(k) {
			t.Errorf("Set should contain %d", k)
		}
	}

	// Verify length
	if ss.Len() != len(data) {
		t.Errorf("Length should be %d, got %d", len(data), ss.Len())
	}

	// Test iteration order (should be sorted by elements)
	elements := make([]int, 0, len(data))
	for val := range ss.IterAsc() {
		elements = append(elements, val)
	}

	// Verify elements are ordered
	for i := 1; i < len(elements); i++ {
		if elements[i] < elements[i-1] {
			t.Errorf("Elements not ordered: elements[%d]=%d < elements[%d]=%d", i, elements[i], i-1, elements[i-1])
		}
	}

	// Test deleting half of the data
	deleteCount := 0
	deleted := make(map[int]bool)
	for k := range data {
		if deleteCount < len(data)/2 {
			ss.Remove(k)
			deleted[k] = true
			deleteCount++
		} else {
			break
		}
	}

	// Verify data after deletion
	for k := range data {
		if deleted[k] {
			// This element should be deleted
			if ss.Contains(k) {
				t.Errorf("Set should not contain deleted element %d", k)
			}
		} else {
			// This element should still be present
			if !ss.Contains(k) {
				t.Errorf("Set should contain %d", k)
			}
		}
	}
}

// Test modifying the set during iteration
func TestSkipSetModifyDuringIteration(t *testing.T) {
	ss := NewOrdered[string]()

	// Insert some data
	for i := 0; i < 10; i++ {
		ss.Insert(fmt.Sprintf("item%d", i))
	}

	// Delete some elements during iteration
	toDelete := []string{"item2", "item5", "item8"}
	for val := range ss.IterAsc() {
		for _, delVal := range toDelete {
			if val == delVal {
				ss.Remove(val)
				break
			}
		}
	}

	// Verify deletion
	for _, delVal := range toDelete {
		if ss.Contains(delVal) {
			t.Errorf("Element %s should be deleted", delVal)
		}
	}

	// Verify other elements still exist
	expectedItems := []string{"item0", "item1", "item3", "item4", "item6", "item7", "item9"}
	for _, item := range expectedItems {
		if !ss.Contains(item) {
			t.Errorf("Element %s should still exist", item)
		}
	}
}

// New test case: Testing all branches of the Range function
func TestSkipSetRangeAllBranches(t *testing.T) {
	ss := NewOrdered[string]()

	// Insert some elements
	data := []string{"a", "b", "c", "d", "e", "f", "g"}
	for _, item := range data {
		ss.Insert(item)
	}

	// Test edge case: lowerBound > upperBound
	lower := "f"
	upper := "b"
	rangeData := make([]string, 0)
	for val := range ss.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower), bound.NewExcluded(upper))) {
		rangeData = append(rangeData, val)
	}
	if len(rangeData) != 0 {
		t.Errorf("Range with lower > upper should return empty, got %v", rangeData)
	}

	// Test edge case: lowerBound equals an element
	lower = "c"
	upper = "f"
	rangeData = make([]string, 0)
	for val := range ss.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower), bound.NewExcluded(upper))) {
		rangeData = append(rangeData, val)
	}
	expected := []string{"c", "d", "e"}
	if len(rangeData) != len(expected) {
		t.Errorf("Range returned wrong number of items: expected %d, got %d", len(expected), len(rangeData))
	}
	for i, val := range rangeData {
		if val != expected[i] {
			t.Errorf("Range returned wrong value at index %d: expected %s, got %s", i, expected[i], val)
		}
	}

	// Test iterator early termination
	count := 0
	for range ss.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[string](), bound.NewUnbounded[string]())) {
		count++
		if count == 3 {
			break
		}
	}
	if count != 3 {
		t.Errorf("Range iterator should yield 3 elements before break, got %d", count)
	}
}

// New test case: Testing all branches of the Intersection function
func TestSkipSetIntersectionAllBranches(t *testing.T) {
	// Create two sets
	set1 := NewOrdered[string]()
	set2 := NewOrdered[string]()
	set3 := NewOrdered[string]()

	// Insert elements
	data1 := []string{"a", "b", "c", "d", "e"}
	data2 := []string{"c", "d", "e", "f", "g"}
	data3 := []string{"x", "y", "z"}

	for _, item := range data1 {
		set1.Insert(item)
	}

	for _, item := range data2 {
		set2.Insert(item)
	}

	for _, item := range data3 {
		set3.Insert(item)
	}

	// Test case with intersection
	intersection := set1.Intersection(set2)
	expected := []string{"c", "d", "e"}
	result := make([]string, 0)
	for val := range intersection.IterAsc() {
		result = append(result, val)
	}
	if len(result) != len(expected) {
		t.Errorf("Intersection length should be %d, got %d", len(expected), len(result))
	}
	for i, val := range result {
		if val != expected[i] {
			t.Errorf("Intersection value at index %d should be %s, got %s", i, expected[i], val)
		}
	}

	// Test case with no intersection
	intersection = set1.Intersection(set3)
	if intersection.Len() != 0 {
		t.Errorf("Intersection of disjoint sets should be empty, got length %d", intersection.Len())
	}

	// Test intersection with empty set
	emptySet := NewOrdered[string]()
	intersection = set1.Intersection(emptySet)
	if intersection.Len() != 0 {
		t.Errorf("Intersection with empty set should be empty, got length %d", intersection.Len())
	}
}

// New test case: Testing all branches of the IsSubset function
func TestSkipSetIsSubsetAllBranches(t *testing.T) {
	// Create test sets
	set1 := NewOrdered[string]()
	set2 := NewOrdered[string]()
	set3 := NewOrdered[string]()

	// Insert elements
	for _, item := range []string{"a", "b", "c", "d", "e"} {
		set1.Insert(item)
	}

	for _, item := range []string{"b", "c", "d"} {
		set2.Insert(item)
	}

	for _, item := range []string{"f", "g", "h"} {
		set3.Insert(item)
	}

	// Test that set2 is a subset of set1
	if !set2.IsSubset(set1) {
		t.Error("set2 should be a subset of set1")
	}

	// test set1 is not subset of set2
	if set1.IsSubset(set2) {
		t.Error("set1 should not be a subset of set2")
	}

	// Test with disjoint sets
	if set3.IsSubset(set1) {
		t.Error("set3 should not be a subset of set1")
	}

	// Test that empty set is a subset of any set
	emptySet := NewOrdered[string]()
	if !emptySet.IsSubset(set1) {
		t.Error("Empty set should be a subset of any set")
	}

	// test set is subset of itself
	if !set1.IsSubset(set1) {
		t.Error("Set should be a subset of itself")
	}
}

// New test case: Testing all branches of the IsDisjoint function
func TestSkipSetIsDisjointAllBranches(t *testing.T) {
	// Create test sets
	set1 := NewOrdered[string]()
	set2 := NewOrdered[string]()
	set3 := NewOrdered[string]()

	// insert elements
	for _, item := range []string{"a", "b", "c", "d", "e"} {
		set1.Insert(item)
	}

	for _, item := range []string{"f", "g", "h"} {
		set2.Insert(item)
	}

	for _, item := range []string{"d", "e", "f"} {
		set3.Insert(item)
	}

	// Test disjoint sets
	if !set1.IsDisjoint(set2) {
		t.Error("set1 and set2 should be disjoint")
	}

	// Test intersecting sets
	if set1.IsDisjoint(set3) {
		t.Error("set1 and set3 should not be disjoint")
	}

	// Test that empty set is disjoint with any set
	emptySet := NewOrdered[string]()
	if !emptySet.IsDisjoint(set1) {
		t.Error("Empty set should be disjoint with any set")
	}

	// Test two empty sets
	emptySet2 := NewOrdered[string]()
	if !emptySet.IsDisjoint(emptySet2) {
		t.Error("Two empty sets should be disjoint")
	}
}
