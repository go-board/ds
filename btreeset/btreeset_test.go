package btreeset

import (
	"github.com/go-board/ds/bound"
	"testing"
)

// Integer comparison function
func intComparator(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// Test basic functionality
func TestBTreeSetBasic(t *testing.T) {
	set := New(intComparator)

	// Test insertion and size
	if !set.Insert(1) {
		t.Errorf("Insert should return true for new element")
	}
	if set.Len() != 1 {
		t.Errorf("Size should be 1 after one insertion, got %d", set.Len())
	}

	// Test duplicate insertion
	if set.Insert(1) {
		t.Errorf("Insert should return false for duplicate element")
	}
	if set.Len() != 1 {
		t.Errorf("Size should remain 1 after duplicate insertion, got %d", set.Len())
	}

	// Test contains
	if !set.Contains(1) {
		t.Errorf("Set should contain element 1")
	}
	if set.Contains(2) {
		t.Errorf("Set should not contain element 2")
	}

	// Test removal
	if !set.Remove(1) {
		t.Errorf("Remove should return true for existing element")
	}
	if set.Len() != 0 {
		t.Errorf("Size should be 0 after removal, got %d", set.Len())
	}
	if !set.IsEmpty() {
		t.Errorf("Set should be empty after removal")
	}

	// Test removal of non-existing element
	if set.Remove(2) {
		t.Errorf("Remove should return false for non-existing element")
	}
}

// Test empty set
func TestBTreeSetEmpty(t *testing.T) {
	set := New(intComparator)

	if !set.IsEmpty() {
		t.Errorf("New set should be empty")
	}

	if set.Len() != 0 {
		t.Errorf("New set size should be 0, got %d", set.Len())
	}

	// Test empty set iterator
	iterCount := 0
	for range set.IterAsc() {
		iterCount++
	}
	if iterCount != 0 {
		t.Errorf("Empty set iterator should yield 0 elements, got %d", iterCount)
	}

	// Test removing non-existent element
	if set.Remove(1) {
		t.Errorf("Remove should return false for non-existing element in empty set")
	}
}

// Test Clear method
func TestBTreeSetClear(t *testing.T) {
	set := New(intComparator)

	// Add elements
	set.Insert(1)
	set.Insert(2)
	set.Insert(3)

	if set.Len() != 3 {
		t.Errorf("Size should be 3 before clear, got %d", set.Len())
	}

	// Clear the set
	set.Clear()

	if !set.IsEmpty() {
		t.Errorf("Set should be empty after clear")
	}

	if set.Len() != 0 {
		t.Errorf("Size should be 0 after clear, got %d", set.Len())
	}

	// Verify elements are removed
	if set.Contains(1) || set.Contains(2) || set.Contains(3) {
		t.Errorf("Elements should be removed after clear")
	}
}

// Test iterator
func TestBTreeSetIterAsc(t *testing.T) {
	set := New(intComparator)

	// Add elements (added out of order, but should iterate in ascending order)
	set.Insert(5)
	set.Insert(2)
	set.Insert(7)
	set.Insert(1)
	set.Insert(3)

	// Verify iteration order is ascending
	expectedOrder := []int{1, 2, 3, 5, 7}
	actualOrder := make([]int, 0, 5)

	for val := range set.IterAsc() {
		actualOrder = append(actualOrder, val)
	}

	if len(actualOrder) != len(expectedOrder) {
		t.Errorf("Expected %d elements, got %d", len(expectedOrder), len(actualOrder))
	}

	for i := range expectedOrder {
		if actualOrder[i] != expectedOrder[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expectedOrder[i], actualOrder[i])
		}
	}
}

// Test set operations
func TestBTreeSetOperations(t *testing.T) {
	// Create two sets
	set1 := New(intComparator)
	set2 := New(intComparator)

	// Add elements
	for _, v := range []int{1, 2, 3, 4, 5} {
		set1.Insert(v)
	}
	for _, v := range []int{3, 4, 5, 6, 7} {
		set2.Insert(v)
	}

	// Test union
	union := set1.Union(set2)
	expectedUnion := []int{1, 2, 3, 4, 5, 6, 7}
	verifySetContent(t, union, expectedUnion, "Union")

	// Test intersection
	intersection := set1.Intersection(set2)
	expectedIntersection := []int{3, 4, 5}
	verifySetContent(t, intersection, expectedIntersection, "Intersection")

	// Test difference
	difference := set1.Difference(set2)
	expectedDifference := []int{1, 2}
	verifySetContent(t, difference, expectedDifference, "Difference")

	// Test symmetric difference
	symmetricDifference := set1.SymmetricDifference(set2)
	expectedSymDiff := []int{1, 2, 6, 7}
	verifySetContent(t, symmetricDifference, expectedSymDiff, "SymmetricDifference")

	// Test subset
	subset := New(intComparator)
	for _, v := range []int{3, 4} {
		subset.Insert(v)
	}
	if !subset.IsSubset(set1) {
		t.Errorf("{3,4} should be a subset of {1,2,3,4,5}")
	}
	if set1.IsSubset(subset) {
		t.Errorf("{1,2,3,4,5} should not be a subset of {3,4}")
	}

	// Test superset
	if !set1.IsSuperset(subset) {
		t.Errorf("{1,2,3,4,5} should be a superset of {3,4}")
	}
	if subset.IsSuperset(set1) {
		t.Errorf("{3,4} should not be a superset of {1,2,3,4,5}")
	}

	// Test disjoint
	disjoint := New(intComparator)
	for _, v := range []int{8, 9} {
		disjoint.Insert(v)
	}
	if !set1.IsDisjoint(disjoint) {
		t.Errorf("{1,2,3,4,5} and {8,9} should be disjoint")
	}
	if set1.IsDisjoint(set2) {
		t.Errorf("{1,2,3,4,5} and {3,4,5,6,7} should not be disjoint")
	}
}

// Test Extend method
func TestBTreeSetExtend(t *testing.T) {
	set := New(intComparator)

	// Create a sequence for extension
	slice := []int{5, 1, 3, 2, 4}

	// Extend the set
	set.Extend(iterFromSlice(slice))

	// Verify set content
	expected := []int{1, 2, 3, 4, 5}
	verifySetContent(t, set, expected, "Extend")
}

// Test Clone method
func TestBTreeSetClone(t *testing.T) {
	set1 := New(intComparator)
	for _, v := range []int{1, 2, 3, 4, 5} {
		set1.Insert(v)
	}

	// Clone the set
	set2 := set1.Clone()

	// Verify cloned set has the same content as original
	expected := []int{1, 2, 3, 4, 5}
	verifySetContent(t, set2, expected, "Clone")

	// Verify modifications to clone don't affect original set
	set2.Insert(6)
	if set1.Contains(6) {
		t.Errorf("Original set should not contain element added to clone")
	}
}

// Test First and Last methods
func TestBTreeSetFirstLast(t *testing.T) {
	set := New(intComparator)

	// Test empty set
	val, found := set.First()
	if found {
		t.Errorf("First should return false for empty set")
	}
	if val != 0 {
		t.Errorf("First should return zero value for empty set")
	}

	val, found = set.Last()
	if found {
		t.Errorf("Last should return false for empty set")
	}
	if val != 0 {
		t.Errorf("Last should return zero value for empty set")
	}

	// Add elements
	for _, v := range []int{5, 2, 7, 1, 3} {
		set.Insert(v)
	}

	// Test First method
	val, found = set.First()
	if !found {
		t.Errorf("First should return true for non-empty set")
	}
	if val != 1 {
		t.Errorf("First should return the smallest element, expected 1, got %d", val)
	}

	// Test Last method
	val, found = set.Last()
	if !found {
		t.Errorf("Last should return true for non-empty set")
	}
	if val != 7 {
		t.Errorf("Last should return the largest element, expected 7, got %d", val)
	}
}

// Test PopFirst and PopLast methods
func TestBTreeSetPopFirstLast(t *testing.T) {
	set := New(intComparator)

	// Add elements
	for _, v := range []int{5, 2, 7, 1, 3} {
		set.Insert(v)
	}

	originalSize := set.Len()

	// Test PopFirst method
	val, found := set.PopFirst()
	if !found {
		t.Errorf("PopFirst should return true for non-empty set")
	}
	if val != 1 {
		t.Errorf("PopFirst should return the smallest element, expected 1, got %d", val)
	}
	if set.Len() != originalSize-1 {
		t.Errorf("Size should decrease by 1 after PopFirst, expected %d, got %d", originalSize-1, set.Len())
	}
	if set.Contains(1) {
		t.Errorf("Element should be removed after PopFirst")
	}

	// Test PopLast method
	val, found = set.PopLast()
	if !found {
		t.Errorf("PopLast should return true for non-empty set")
	}
	if val != 7 {
		t.Errorf("PopLast should return the largest element, expected 7, got %d", val)
	}
	if set.Len() != originalSize-2 {
		t.Errorf("Size should decrease by 1 after PopLast, expected %d, got %d", originalSize-2, set.Len())
	}
	if set.Contains(7) {
		t.Errorf("Element should be removed after PopLast")
	}

	// Test behavior after clearing the set
	for set.Len() > 0 {
		set.PopFirst()
	}

	val, found = set.PopFirst()
	if found {
		t.Errorf("PopFirst should return false for empty set")
	}
	if val != 0 {
		t.Errorf("PopFirst should return zero value for empty set")
	}

	val, found = set.PopLast()
	if found {
		t.Errorf("PopLast should return false for empty set")
	}
	if val != 0 {
		t.Errorf("PopLast should return zero value for empty set")
	}
}

// Test Range method
func TestBTreeSetRangeAsc(t *testing.T) {
	set := New(intComparator)

	// Add elements
	for i := 1; i <= 10; i++ {
		set.Insert(i)
	}

	// Test full range
	range1 := set.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewUnbounded[int]()))
	count1 := 0
	for range range1 {
		count1++
	}
	if count1 != 10 {
		t.Errorf("Range with no bounds should return all 10 elements, got %d", count1)
	}

	// Test with upper and lower bounds
	lower3 := 3
	upper7 := 7
	range2 := set.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower3), bound.NewExcluded(upper7)))
	values2 := make([]int, 0)
	for val := range range2 {
		values2 = append(values2, val)
	}
	expected2 := []int{3, 4, 5, 6}
	if len(values2) != len(expected2) {
		t.Errorf("Range with bounds [3,7) should return 4 elements, got %d", len(values2))
	}
	for i, v := range values2 {
		if v != expected2[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected2[i], v)
		}
	}

	// Test with only lower bound
	lower8 := 8
	range3 := set.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower8), bound.NewUnbounded[int]()))
	values3 := make([]int, 0)
	for val := range range3 {
		values3 = append(values3, val)
	}
	expected3 := []int{8, 9, 10}
	if len(values3) != len(expected3) {
		t.Errorf("Range with lower bound 8 should return 3 elements, got %d", len(values3))
	}
	for i, v := range values3 {
		if v != expected3[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected3[i], v)
		}
	}

	// Test with only upper bound
	upper4 := 4
	range4 := set.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewExcluded(upper4)))
	values4 := make([]int, 0)
	for val := range range4 {
		values4 = append(values4, val)
	}
	expected4 := []int{1, 2, 3}
	if len(values4) != len(expected4) {
		t.Errorf("Range with upper bound 4 should return 3 elements, got %d", len(values4))
	}
	for i, v := range values4 {
		if v != expected4[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected4[i], v)
		}
	}

	// Test with no matching elements
	lower11 := 11
	range5 := set.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower11), bound.NewUnbounded[int]()))
	count5 := 0
	for range range5 {
		count5++
	}
	if count5 != 0 {
		t.Errorf("Range with lower bound 11 should return 0 elements, got %d", count5)
	}
}

// Helper function: create iterator from slice
func iterFromSlice[T any](slice []T) func(func(T) bool) {
	return func(yield func(T) bool) {
		for _, item := range slice {
			if !yield(item) {
				return
			}
		}
	}
}

// Helper function: verify set content
func verifySetContent[T comparable](t *testing.T, set *BTreeSet[T], expected []T, operation string) {
	t.Helper()

	// Collect elements from the set
	actual := make([]T, 0, set.Len())
	for val := range set.IterAsc() {
		actual = append(actual, val)
	}

	// Verify element count
	if len(actual) != len(expected) {
		t.Errorf("%s: expected %d elements, got %d", operation, len(expected), len(actual))
		return
	}

	// Verify each element exists and is in correct order
	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("%s: at index %d, expected %v, got %v", operation, i, expected[i], actual[i])
		}
	}
}

func TestBTreeSetEqual(t *testing.T) {
	a := New(intComparator)
	b := New(intComparator)
	for _, v := range []int{1, 2, 3} {
		a.Insert(v)
		b.Insert(v)
	}
	if !a.Equal(b) {
		t.Fatal("sets with same elements should be equal")
	}

	b.Insert(4)
	if a.Equal(b) {
		t.Fatal("sets with different sizes should not be equal")
	}

	c := New(intComparator)
	for _, v := range []int{1, 2, 4} {
		c.Insert(v)
	}
	if a.Equal(c) {
		t.Fatal("sets with same size but different elements should not be equal")
	}
}
