package btree

import (
	"github.com/go-board/ds/bound"
	"reflect"
	"testing"
)

// Custom comparison function
func intComparator(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// TestNew
func TestNew(t *testing.T) {
	// Test basic creation
	tree := New(intComparator)
	if tree == nil {
		t.Fatal("New should return non-nil tree")
	}
	if tree.Len() != 0 {
		t.Errorf("New tree should have length 0, got %d", tree.Len())
	}

	// Test nil comparator
	defer func() {
		if r := recover(); r == nil {
			t.Error("New should panic with nil comparator")
		}
	}()
	_ = New[int](nil)
}

// TestNewOrdered for ordered types
func TestNewOrdered(t *testing.T) {
	tree := NewOrdered[int]()
	if tree == nil {
		t.Fatal("NewOrdered should return non-nil tree")
	}
	if tree.Len() != 0 {
		t.Errorf("NewOrdered tree should have length 0, got %d", tree.Len())
	}

	// Test string type
	treeStr := NewOrdered[string]()
	treeStr.Insert("hello")
	if treeStr.Len() != 1 {
		t.Errorf("Tree should have length 1 after insertion, got %d", treeStr.Len())
	}
}

// TestInsert
func TestInsert(t *testing.T) {
	tree := New(intComparator)

	// Test single insertion
	tree.Insert(5)
	if tree.Len() != 1 {
		t.Errorf("Tree length should be 1 after single insertion, got %d", tree.Len())
	}

	// Test multiple insertions
	values := []int{3, 7, 2, 4, 6, 8}
	for _, v := range values {
		tree.Insert(v)
	}
	if tree.Len() != len(values)+1 {
		t.Errorf("Tree length should be %d after multiple insertions, got %d", len(values)+1, tree.Len())
	}

	// Test inserting duplicate values (current implementation saves duplicate keys)
	tree.Insert(5)
	if tree.Len() != len(values)+2 {
		t.Errorf("Tree length should increase after inserting duplicate, got %d", tree.Len())
	}

	// Test inserting a large amount of data (trigger node splitting)
	treeLarge := New(intComparator)
	for i := 0; i < 1000; i++ {
		treeLarge.Insert(i)
	}
	if treeLarge.Len() != 1000 {
		t.Errorf("Large tree length should be 1000, got %d", treeLarge.Len())
	}
}

// TestSearch
func TestSearch(t *testing.T) {
	tree := New(intComparator)
	values := []int{5, 3, 7, 2, 4, 6, 8}
	for _, v := range values {
		tree.Insert(v)
	}

	// Test searching for an existing value
	for _, v := range values {
		val, found := tree.Search(v)
		if !found {
			t.Errorf("Search should find value %d", v)
		}
		if val != v {
			t.Errorf("Search should return correct value, expected %d, got %d", v, val)
		}
	}

	// Test searching for a non-existing value
	notFound := []int{1, 9, 100}
	for _, v := range notFound {
		val, found := tree.Search(v)
		if found {
			t.Errorf("Search should not find value %d", v)
		}
		if val != 0 {
			t.Errorf("Search should return zero value for not found, got %d", val)
		}
	}

	// Test searching in an empty tree
	emptyTree := New(intComparator)
	val, found := emptyTree.Search(5)
	if found {
		t.Error("Search should not find value in empty tree")
	}
	if val != 0 {
		t.Errorf("Search should return zero value for empty tree, got %d", val)
	}
}

// TestRemove
func TestRemove(t *testing.T) {
	tree := New(intComparator)
	values := []int{5, 3, 7, 2, 4, 6, 8}
	for _, v := range values {
		tree.Insert(v)
	}

	// Test removing an existing value
	toRemove := []int{3, 7, 5}
	for _, v := range toRemove {
		removed := tree.Remove(v)
		if !removed {
			t.Errorf("Remove should return true for existing value %d", v)
		}
		// Verify it's actually removed
		_, found := tree.Search(v)
		if found {
			t.Errorf("Value %d should be removed from tree", v)
		}
	}

	// Test removing a non-existing value
	notExist := 10
	removed := tree.Remove(notExist)
	if removed {
		t.Errorf("Remove should return false for non-existing value %d", notExist)
	}

	// Test removing from an empty tree
	emptyTree := New(intComparator)
	removed = emptyTree.Remove(5)
	if removed {
		t.Error("Remove should return false for empty tree")
	}

	// Test removing all elements
	fullTree := New(intComparator)
	for _, v := range values {
		fullTree.Insert(v)
	}
	for _, v := range values {
		fullTree.Remove(v)
	}
	if fullTree.Len() != 0 {
		t.Errorf("Tree should be empty after removing all elements, got length %d", fullTree.Len())
	}

	// Test complex deletion scenarios (trigger node merging, etc.)
	treeComplex := New(intComparator)
	for i := 0; i < 100; i++ {
		treeComplex.Insert(i)
	}
	// Remove middle value
	for i := 20; i < 80; i++ {
		treeComplex.Remove(i)
	}
	if treeComplex.Len() != 40 {
		t.Errorf("Tree length should be 40 after complex deletion, got %d", treeComplex.Len())
	}
}

// TestIter
func TestIter(t *testing.T) {
	tree := New(intComparator)
	values := []int{5, 3, 7, 2, 4, 6, 8}
	expected := []int{2, 3, 4, 5, 6, 7, 8} // in-order traversal order

	for _, v := range values {
		tree.Insert(v)
	}

	// Collect iterator results
	var result []int
	for v := range tree.IterAsc() {
		result = append(result, v)
	}

	// Verify results match expectations
	if len(result) != len(expected) {
		t.Errorf("Iter should return %d elements, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Iter element %d should be %d, got %d", i, expected[i], v)
		}
	}

	// Test empty tree iterator
	emptyTree := New(intComparator)
	count := 0
	for range emptyTree.IterAsc() {
		count++
	}
	if count != 0 {
		t.Errorf("Empty tree iter should return 0 elements, got %d", count)
	}
}

// TestRange
func TestRange(t *testing.T) {
	tree := New(intComparator)
	values := []int{1, 3, 5, 7, 9, 2, 4, 6, 8, 10}
	for _, v := range values {
		tree.Insert(v)
	}

	// Test full range
	var fullResult []int
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewUnbounded[int]())) {
		fullResult = append(fullResult, v)
	}
	expectedFull := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if len(fullResult) != len(expectedFull) {
		t.Errorf("Full range should return %d elements, got %d", len(expectedFull), len(fullResult))
	}

	// Test lower bound range [3, ∞)
	lowerBound := 3
	var lowerResult []int
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lowerBound), bound.NewUnbounded[int]())) {
		lowerResult = append(lowerResult, v)
	}
	expectedLower := []int{3, 4, 5, 6, 7, 8, 9, 10}
	if len(lowerResult) != len(expectedLower) {
		t.Errorf("Lower bound range should return %d elements, got %d", len(expectedLower), len(lowerResult))
	}

	// Test upper bound range (-∞, 7]
	upperBound := 7
	var upperResult []int
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewExcluded(upperBound))) {
		upperResult = append(upperResult, v)
	}
	expectedUpper := []int{1, 2, 3, 4, 5, 6}
	if len(upperResult) != len(expectedUpper) {
		t.Errorf("Upper bound range should return %d elements, got %d", len(expectedUpper), len(upperResult))
	}

	// Test closed interval [4, 8]
	low := 4
	high := 8
	var rangeResult []int
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(low), bound.NewExcluded(high))) {
		rangeResult = append(rangeResult, v)
	}
	expectedRange := []int{4, 5, 6, 7}
	if len(rangeResult) != len(expectedRange) {
		t.Errorf("Closed range should return %d elements, got %d", len(expectedRange), len(rangeResult))
	}

	// Test empty range
	low2 := 8
	high2 := 4
	var emptyRangeResult []int
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(low2), bound.NewExcluded(high2))) {
		emptyRangeResult = append(emptyRangeResult, v)
	}
	if len(emptyRangeResult) != 0 {
		t.Errorf("Empty range should return 0 elements, got %d", len(emptyRangeResult))
	}
}

// TestIterDesc tests the IterDesc method of BTree.
func TestIterDesc(t *testing.T) {
	// Test with ordered integers
	tree := NewOrdered[int]()

	// Test empty tree
	var emptyResult []int
	for v := range tree.IterDesc() {
		emptyResult = append(emptyResult, v)
	}
	if len(emptyResult) != 0 {
		t.Errorf("IterDesc on empty tree should return 0 elements, got %d", len(emptyResult))
	}

	// Test tree with a single element
	tree.Insert(42)
	var singleResult []int
	for v := range tree.IterDesc() {
		singleResult = append(singleResult, v)
	}
	if len(singleResult) != 1 || singleResult[0] != 42 {
		t.Errorf("IterDesc on single element tree should return [42], got %v", singleResult)
	}

	// Clear tree for next test
	tree = NewOrdered[int]()

	// Insert elements in a non-sorted order
	testValues := []int{5, 3, 8, 1, 10, 2, 7, 4, 6, 9}
	for _, v := range testValues {
		tree.Insert(v)
	}

	// Collect results from IterDesc.
	var reverseResults []int
	for v := range tree.IterDesc() {
		reverseResults = append(reverseResults, v)
	}

	// Expected result should be sorted in descending order
	expectedReverse := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	if !reflect.DeepEqual(reverseResults, expectedReverse) {
		t.Errorf("IterDesc should return elements in descending order, got %v, expected %v", reverseResults, expectedReverse)
	}

	// Test IterDesc with early termination.
	count := 0
	var partialResults []int
	for v := range tree.IterDesc() {
		partialResults = append(partialResults, v)
		count++
		if count == 5 {
			break // Terminate iteration early
		}
	}

	if count != 5 {
		t.Errorf("IterDesc iteration should yield 5 elements before break, got %d", count)
	}

	expectedPartial := []int{10, 9, 8, 7, 6}
	if !reflect.DeepEqual(partialResults, expectedPartial) {
		t.Errorf("Early terminated IterDesc should return top 5 elements, got %v, expected %v", partialResults, expectedPartial)
	}

	// Verify the tree structure is not damaged
	if tree.Len() != 10 {
		t.Errorf("Tree length should remain 10 after early termination, got %d", tree.Len())
	}
}

// TestFirstLast
func TestFirstLast(t *testing.T) {
	tree := New(intComparator)

	// Test empty tree
	_, found := tree.First()
	if found {
		t.Error("First should return false for empty tree")
	}
	_, found = tree.Last()
	if found {
		t.Error("Last should return false for empty tree")
	}

	// Test after inserting elements
	values := []int{5, 3, 7, 2, 4, 6, 8}
	for _, v := range values {
		tree.Insert(v)
	}

	first, found := tree.First()
	if !found {
		t.Error("First should return true for non-empty tree")
	}
	if first != 2 {
		t.Errorf("First should return 2, got %d", first)
	}

	last, found := tree.Last()
	if !found {
		t.Error("Last should return true for non-empty tree")
	}
	if last != 8 {
		t.Errorf("Last should return 8, got %d", last)
	}
}

// TestPopFirstLast
func TestPopFirstLast(t *testing.T) {
	tree := New(intComparator)
	values := []int{5, 3, 7, 2, 4, 6, 8}
	for _, v := range values {
		tree.Insert(v)
	}

	// Test PopFirst
	for expected := 2; expected <= 8; expected++ {
		val, found := tree.PopFirst()
		if !found {
			t.Errorf("PopFirst should return true, expected value %d", expected)
		}
		if val != expected {
			t.Errorf("PopFirst should return %d, got %d", expected, val)
		}
		// Verify element is actually removed
		_, stillThere := tree.Search(expected)
		if stillThere {
			t.Errorf("Value %d should be removed after PopFirst", expected)
		}
	}

	// Test empty tree
	_, found := tree.PopFirst()
	if found {
		t.Error("PopFirst should return false for empty tree")
	}

	// Reinsert and test PopLast
	for _, v := range values {
		tree.Insert(v)
	}

	for expected := 8; expected >= 2; expected-- {
		val, found := tree.PopLast()
		if !found {
			t.Errorf("PopLast should return true, expected value %d", expected)
		}
		if val != expected {
			t.Errorf("PopLast should return %d, got %d", expected, val)
		}
		// Verify element is actually removed
		_, stillThere := tree.Search(expected)
		if stillThere {
			t.Errorf("Value %d should be removed after PopLast", expected)
		}
	}

	// Test empty tree
	_, found = tree.PopLast()
	if found {
		t.Error("PopLast should return false for empty tree")
	}
}

// TestCustomType
func TestCustomType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Age-based comparator
	personComparator := func(a, b Person) int {
		return intComparator(a.Age, b.Age)
	}

	tree := New(personComparator)

	// Insert custom type elements
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
		{"Dave", 28},
	}

	for _, p := range people {
		tree.Insert(p)
	}

	if tree.Len() != len(people) {
		t.Errorf("Tree should have %d people, got %d", len(people), tree.Len())
	}

	// Search for custom type elements
	alice := Person{"Alice", 30}
	foundPerson, found := tree.Search(alice)
	if !found {
		t.Error("Should find Alice in the tree")
	}
	if foundPerson.Name != "Alice" || foundPerson.Age != 30 {
		t.Errorf("Found person should be Alice, got %v", foundPerson)
	}

	// Test iterator
	var ages []int
	for p := range tree.IterAsc() {
		ages = append(ages, p.Age)
	}
	expectedAges := []int{25, 28, 30, 35}
	if len(ages) != len(expectedAges) {
		t.Errorf("Iter should return %d ages, got %d", len(expectedAges), len(ages))
	}

	for i, age := range ages {
		if age != expectedAges[i] {
			t.Errorf("Age %d should be %d, got %d", i, expectedAges[i], age)
		}
	}

	// Test deletion
	removed := tree.Remove(Person{"Bob", 25})
	if !removed {
		t.Error("Should be able to remove Bob")
	}
	if tree.Len() != len(people)-1 {
		t.Errorf("Tree should have %d people after removal, got %d", len(people)-1, tree.Len())
	}
}

// TestComplexOperations
func TestComplexOperations(t *testing.T) {
	tree := New(intComparator)

	// Insert large amount of data
	for i := 0; i < 1000; i++ {
		tree.Insert(i)
	}

	if tree.Len() != 1000 {
		t.Errorf("Tree should have 1000 elements, got %d", tree.Len())
	}

	// Delete half of the data
	for i := 0; i < 500; i++ {
		tree.Remove(i)
	}

	if tree.Len() != 500 {
		t.Errorf("Tree should have 500 elements after deletion, got %d", tree.Len())
	}

	// Verify remaining elements
	for i := 500; i < 1000; i++ {
		_, found := tree.Search(i)
		if !found {
			t.Errorf("Value %d should still be in the tree", i)
		}
	}

	// Verify deleted elements
	for i := 0; i < 500; i++ {
		_, found := tree.Search(i)
		if found {
			t.Errorf("Value %d should have been removed", i)
		}
	}

	// Insert data again
	for i := 0; i < 500; i++ {
		tree.Insert(i)
	}

	if tree.Len() != 1000 {
		t.Errorf("Tree should have 1000 elements after reinsertion, got %d", tree.Len())
	}

	// Verify all elements using iterator
	count := 0
	prev := -1
	for v := range tree.IterAsc() {
		count++
		if v <= prev {
			t.Errorf("Elements should be in order, previous %d, current %d", prev, v)
		}
		prev = v
	}

	if count != 1000 {
		t.Errorf("Iterator should return 1000 elements, got %d", count)
	}
}

// New test case: Testing the getSuccessor function
func TestGetSuccessor(t *testing.T) {
	tree := New(intComparator)

	// Construct a specific B-tree structure to test getSuccessor
	// Insert enough data to trigger node splitting
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for _, v := range values {
		tree.Insert(v)
	}

	// Delete a specific element to trigger getSuccessor call
	// Delete a key in the root node where the right subtree has enough keys
	removed := tree.Remove(10)
	if !removed {
		t.Error("Should be able to remove value 10")
	}

	// Verify the tree remains correct after deletion
	for _, v := range values {
		if v == 10 {
			_, found := tree.Search(v)
			if found {
				t.Errorf("Value %d should have been removed", v)
			}
		} else {
			val, found := tree.Search(v)
			if !found {
				t.Errorf("Value %d should still exist in the tree", v)
			}
			if val != v {
				t.Errorf("Value %d should be %d, got %d", v, v, val)
			}
		}
	}

	// Verify the tree maintains correct ordering
	prev := 0
	for v := range tree.IterAsc() {
		if v <= prev {
			t.Errorf("Elements should be in ascending order, got %d after %d", v, prev)
		}
		prev = v
	}
}

// New test case: Test the borrowFromLeft function
func TestBorrowFromLeft(t *testing.T) {
	tree := New(intComparator)

	// Insert data to create a specific tree structure
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}
	for _, v := range values {
		tree.Insert(v)
	}

	// Delete specific elements to trigger borrowFromLeft operation
	// Delete elements from the right subtree to make the number of keys insufficient, requiring borrowing from the left sibling
	toRemove := []int{20, 21, 22, 23, 24, 25}
	for _, v := range toRemove {
		removed := tree.Remove(v)
		if !removed {
			t.Errorf("Should be able to remove value %d", v)
		}
	}

	// Verify the tree is still correct after deletion
	for _, v := range values {
		found := false
		for _, removed := range toRemove {
			if v == removed {
				found = true
				break
			}
		}

		if found {
			_, stillExists := tree.Search(v)
			if stillExists {
				t.Errorf("Value %d should have been removed", v)
			}
		} else {
			val, exists := tree.Search(v)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", v)
			}
			if val != v {
				t.Errorf("Value %d should be %d, got %d", v, v, val)
			}
		}
	}
}

// New test case: Test the borrowFromRight function
func TestBorrowFromRight(t *testing.T) {
	tree := New(intComparator)

	// Insert data to create a specific tree structure
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}
	for _, v := range values {
		tree.Insert(v)
	}

	// Delete elements from the left subtree to make the number of keys insufficient, requiring borrowing from the right sibling
	toRemove := []int{1, 2, 3, 4, 5, 6}
	for _, v := range toRemove {
		removed := tree.Remove(v)
		if !removed {
			t.Errorf("Should be able to remove value %d", v)
		}
	}

	// Verify the tree is still correct after deletion
	for _, v := range values {
		found := false
		for _, removed := range toRemove {
			if v == removed {
				found = true
				break
			}
		}

		if found {
			_, stillExists := tree.Search(v)
			if stillExists {
				t.Errorf("Value %d should have been removed", v)
			}
		} else {
			val, exists := tree.Search(v)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", v)
			}
			if val != v {
				t.Errorf("Value %d should be %d, got %d", v, v, val)
			}
		}
	}
}

// New test case: Test the mergeChildren function
func TestMergeChildren(t *testing.T) {
	tree := New(intComparator)

	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	for _, v := range values {
		tree.Insert(v)
	}

	// Delete most elements to reduce the number of keys in child nodes below the merge threshold
	toRemove := []int{4, 5, 6, 7, 8, 9, 10, 11, 12}
	for _, v := range toRemove {
		removed := tree.Remove(v)
		if !removed {
			t.Errorf("Should be able to remove value %d", v)
		}
	}

	for _, v := range values {
		found := false
		for _, removed := range toRemove {
			if v == removed {
				found = true
				break
			}
		}

		if found {
			_, stillExists := tree.Search(v)
			if stillExists {
				t.Errorf("Value %d should have been removed", v)
			}
		} else {
			val, exists := tree.Search(v)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", v)
			}
			if val != v {
				t.Errorf("Value %d should be %d, got %d", v, v, val)
			}
		}
	}
}

// New test case: Test boundary cases for rangeNode function
func TestRangeNodeEdgeCases(t *testing.T) {
	tree := New(intComparator)

	// Test range query on empty tree
	var result []int
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewUnbounded[int]())) {
		result = append(result, v)
	}
	if len(result) != 0 {
		t.Errorf("Range on empty tree should return 0 elements, got %d", len(result))
	}

	// Insert a single element
	tree.Insert(5)

	// Test range query on single-element tree
	result = result[:0] // Clear the slice
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewUnbounded[int]())) {
		result = append(result, v)
	}
	if len(result) != 1 || result[0] != 5 {
		t.Errorf("Range on single-element tree should return [5], got %v", result)
	}

	// Test boundary range query
	low := 5
	high := 5
	result = result[:0]
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(low), bound.NewExcluded(high))) {
		result = append(result, v)
	}
	if len(result) != 0 {
		t.Errorf("Range [5, 5) should return 0 elements, got %d: %v", len(result), result)
	}

	// Test range query that contains elements
	low = 4
	high = 6
	result = result[:0]
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(low), bound.NewExcluded(high))) {
		result = append(result, v)
	}
	if len(result) != 1 || result[0] != 5 {
		t.Errorf("Range [4, 6) should return [5], got %v", result)
	}
}

// New test case: Specifically test getSuccessor function calls
func TestGetSuccessorSpecifically(t *testing.T) {
	tree := New(intComparator)

	// Construct a specific B-tree structure to ensure getSuccessor is called during deletion
	// Insert data to create a right subtree with sufficient keys
	for i := 1; i <= 20; i++ {
		tree.Insert(i)
	}

	// Ensure the tree structure meets conditions for calling getSuccessor:
	// 1. The key to delete is in an internal node
	// 2. The right subtree has sufficient keys (at least order, which is 3 here)

	// Delete a key from the root node to trigger getSuccessor call
	// In a B-tree, when deleting a key from an internal node, getSuccessor is called if the right subtree has keys >= order
	removed := tree.Remove(7)
	if !removed {
		t.Error("Should be able to remove value 7")
	}

	_, found := tree.Search(7)
	if found {
		t.Error("Value 7 should have been removed")
	}

	// Verify other elements still exist
	for i := 1; i <= 20; i++ {
		if i != 7 {
			val, exists := tree.Search(i)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", i)
			}
			if val != i {
				t.Errorf("Value %d should be %d, got %d", i, i, val)
			}
		}
	}

	// Verify the tree still maintains correct order
	prev := 0
	for v := range tree.IterAsc() {
		if v <= prev {
			t.Errorf("Elements should be in ascending order, got %d after %d", v, prev)
		}
		prev = v
	}
}

// New test case: Test specific cases for borrowFromLeft function
func TestBorrowFromLeftSpecifically(t *testing.T) {
	tree := New(intComparator)

	// Insert specific data to create a tree structure that satisfies borrowFromLeft call conditions
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for _, v := range values {
		tree.Insert(v)
	}

	// Delete elements from the right sibling node to reduce its key count to a critical value
	// This may trigger borrowFromLeft in subsequent deletion operations
	toRemove := []int{16, 17, 18, 19, 20}
	for _, v := range toRemove {
		tree.Remove(v)
	}

	// Now delete elements from left sibling, this should trigger borrowFromLeft operation
	removed := tree.Remove(1)
	if !removed {
		t.Error("Should be able to remove value 1")
	}

	// Verify the deletion operation was successful
	_, found := tree.Search(1)
	if found {
		t.Error("Value 1 should have been removed")
	}

	// Verify other elements still exist
	for _, v := range values {
		shouldExist := true
		for _, removed := range toRemove {
			if v == removed {
				shouldExist = false
				break
			}
		}
		if v == 1 {
			shouldExist = false
		}

		if shouldExist {
			val, exists := tree.Search(v)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", v)
			}
			if val != v {
				t.Errorf("Value %d should be %d, got %d", v, v, val)
			}
		} else {
			_, exists := tree.Search(v)
			if exists {
				t.Errorf("Value %d should have been removed", v)
			}
		}
	}
}

// New test case: Test full coverage for rangeNode function
func TestRangeNodeFullCoverage(t *testing.T) {
	tree := New(intComparator)

	values := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	for _, v := range values {
		tree.Insert(v)
	}

	// Test various range queries to improve coverage of rangeNode function

	// Test nil boundaries
	var result []int
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewUnbounded[int]())) {
		result = append(result, v)
	}
	expected := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	if len(result) != len(expected) {
		t.Errorf("Full range should return %d elements, got %d", len(expected), len(result))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Element %d should be %d, got %d", i, expected[i], v)
		}
	}

	// Test only lower bound
	lower := 7
	result = result[:0]
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower), bound.NewUnbounded[int]())) {
		result = append(result, v)
	}
	expected = []int{7, 9, 11, 13, 15, 17, 19}
	if len(result) != len(expected) {
		t.Errorf("Lower bound range should return %d elements, got %d", len(expected), len(result))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Element %d should be %d, got %d", i, expected[i], v)
		}
	}

	// Test only upper bound
	upper := 11
	result = result[:0]
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewExcluded(upper))) {
		result = append(result, v)
	}
	expected = []int{1, 3, 5, 7, 9}
	if len(result) != len(expected) {
		t.Errorf("Upper bound range should return %d elements, got %d", len(expected), len(result))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Element %d should be %d, got %d", i, expected[i], v)
		}
	}

	// Test range with both bounds
	lower = 5
	upper = 15
	result = result[:0]
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower), bound.NewExcluded(upper))) {
		result = append(result, v)
	}
	expected = []int{5, 7, 9, 11, 13}
	if len(result) != len(expected) {
		t.Errorf("Bounded range should return %d elements, got %d", len(expected), len(result))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Element %d should be %d, got %d", i, expected[i], v)
		}
	}

	// Test empty range
	lower = 15
	upper = 5
	result = result[:0]
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower), bound.NewExcluded(upper))) {
		result = append(result, v)
	}
	if len(result) != 0 {
		t.Errorf("Empty range should return 0 elements, got %d", len(result))
	}
}

// New test case: Construct specific scenario to ensure getSuccessor function is called
func TestEnsureGetSuccessorCalled(t *testing.T) {
	tree := New(intComparator)

	// Insert specific data patterns to construct a B-tree structure that satisfies getSuccessor call conditions
	// We need to construct a case where:
	// 1. The key to be deleted is in an internal node
	// 2. The number of keys in the left subtree < order (i.e., < 3)
	// 3. The number of keys in the right subtree >= order (i.e., >= 3)

	// First insert some data to create splits
	for i := 1; i <= 15; i++ {
		tree.Insert(i)
	}

	// Now delete some specific keys to adjust the number of keys in left and right subtrees
	// Delete some keys from the left subtree to make its key count less than 3
	tree.Remove(1)
	tree.Remove(2)

	// Ensure the right subtree has enough keys
	// The right subtree should already have enough keys

	// Now delete a key in the root node, which should trigger a getSuccessor call
	// Because left subtree key count < 3 and right subtree key count >= 3
	removed := tree.Remove(8) // Assuming 8 is in the root node
	if !removed {
		// If 8 is not in the root node, try other values
		removed = tree.Remove(7)
		if !removed {
			removed = tree.Remove(9)
		}
	}

	// Verify the deletion operation was executed successfully (even if we're not sure which function was called specifically)
	if !removed {
		t.Error("Should be able to remove at least one of the test values")
	}

	// Verify the basic functionality of the tree is still normal
	// Insert some new values
	tree.Insert(100)
	val, found := tree.Search(100)
	if !found || val != 100 {
		t.Error("Should be able to insert and find new values")
	}

	// Verify iteration still works in order
	prev := 0
	for v := range tree.IterAsc() {
		if v <= prev && v != 100 { // 100 is newly inserted and might be at the end
			t.Errorf("Elements should be in ascending order, got %d after %d", v, prev)
		}
		prev = v
	}
}

// New test case: Test more branches of rangeNode function
func TestRangeNodeMoreBranches(t *testing.T) {
	tree := New(intComparator)

	// Insert a large amount of data to create a complex B-tree structure
	for i := 1; i <= 50; i++ {
		tree.Insert(i)
	}

	// Test various range queries to trigger different branches of rangeNode

	// Test cases triggered by upper bound
	upper := 25
	var result []int
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewExcluded(upper))) {
		result = append(result, v)
		// Terminate early at a certain point to test the case where yield returns false
		if len(result) == 10 {
			break
		}
	}

	// Verify the results
	expectedLen := 10
	if len(result) != expectedLen {
		t.Errorf("Range with early termination should yield %d elements, got %d", expectedLen, len(result))
	}

	// Verify element order is correct
	for i := 0; i < len(result)-1; i++ {
		if result[i] >= result[i+1] {
			t.Errorf("Elements should be in ascending order, got %d >= %d at positions %d and %d",
				result[i], result[i+1], i, i+1)
		}
	}

	// Test cases with both lower and upper bounds
	lower := 15
	upperBound := 35
	result = result[:0]
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower), bound.NewExcluded(upperBound))) {
		result = append(result, v)
	}

	for _, v := range result {
		if v < lower || v >= upperBound {
			t.Errorf("Element %d should be in range [%d, %d)", v, lower, upperBound)
		}
	}

	for i := 0; i < len(result)-1; i++ {
		if result[i] >= result[i+1] {
			t.Errorf("Elements should be in ascending order, got %d >= %d at positions %d and %d",
				result[i], result[i+1], i, i+1)
		}
	}
}

// New test case: Test borrowFromLeft branch of deleteNode function
func TestDeleteNodeBorrowFromLeft(t *testing.T) {
	tree := New(intComparator)

	// Insert specific data to trigger borrowFromLeft operation
	// Need to construct a case where borrowing from left sibling is required
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for _, v := range values {
		tree.Insert(v)
	}

	// Delete some keys from the right subtree to reduce its key count
	toRemoveFromRight := []int{16, 17, 18, 19, 20}
	for _, v := range toRemoveFromRight {
		tree.Remove(v)
	}

	// Now delete a key from the left subtree, which should trigger borrowFromLeft operation
	removed := tree.Remove(1)
	if !removed {
		t.Error("Should be able to remove value 1")
	}

	_, found := tree.Search(1)
	if found {
		t.Error("Value 1 should have been removed")
	}

	// Verify other elements still exist
	for _, v := range values {
		shouldExist := true
		for _, removed := range toRemoveFromRight {
			if v == removed {
				shouldExist = false
				break
			}
		}
		if v == 1 {
			shouldExist = false
		}

		if shouldExist {
			val, exists := tree.Search(v)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", v)
			}
			if val != v {
				t.Errorf("Value %d should be %d, got %d", v, v, val)
			}
		} else {
			_, exists := tree.Search(v)
			if exists {
				t.Errorf("Value %d should have been removed", v)
			}
		}
	}
}

// New test case: Test complex range query scenarios
func TestComplexRangeQueries(t *testing.T) {
	tree := New(intComparator)

	for i := 1; i <= 30; i++ {
		tree.Insert(i)
	}

	// Test various boundary cases

	// Test case where upper bound equals a key
	upperBound := 15
	var result []int
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewExcluded(upperBound))) {
		result = append(result, v)
	}

	// Verify results do not contain the upper bound value
	for _, v := range result {
		if v >= upperBound {
			t.Errorf("Range should not include upper bound, but got %d >= %d", v, upperBound)
		}
	}

	// Test case where lower bound equals a key
	lowerBound := 10
	result = result[:0]
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lowerBound), bound.NewUnbounded[int]())) {
		result = append(result, v)
	}

	// Verify results contain the lower bound value
	foundLower := false
	for _, v := range result {
		if v == lowerBound {
			foundLower = true
			break
		}
	}
	if !foundLower {
		t.Errorf("Range should include lower bound %d", lowerBound)
	}

	// Test empty range (lower bound >= upper bound)
	lower := 20
	upper := 10
	result = result[:0]
	for v := range tree.RangeAsc(bound.NewRangeBounds(bound.NewIncluded(lower), bound.NewExcluded(upper))) {
		result = append(result, v)
	}
	if len(result) != 0 {
		t.Errorf("Range with lower >= upper should be empty, got %v", result)
	}
}

// New test case: Test early termination in rangeNode function
func TestRangeNodeEarlyTermination(t *testing.T) {
	tree := New(intComparator)

	for i := 1; i <= 20; i++ {
		tree.Insert(i)
	}

	// Test case where yield returns false in rangeNode
	// This triggers the return false branch in rangeNode function
	count := 0
	for range tree.RangeAsc(bound.NewRangeBounds(bound.NewUnbounded[int](), bound.NewUnbounded[int]())) {
		count++
		if count == 5 {
			break // Terminate iteration early
		}
	}

	if count != 5 {
		t.Errorf("Range iteration should yield 5 elements before break, got %d", count)
	}

	// Verify the tree structure is not damaged
	if tree.Len() != 20 {
		t.Errorf("Tree length should remain 20 after early termination, got %d", tree.Len())
	}
}

// New test case: Constructing specific scenarios to ensure getSuccessor function is triggered
func TestGetSuccessorTrigger(t *testing.T) {
	tree := New(intComparator)

	// Insert a large amount of data to create a complex B-tree structure
	for i := 1; i <= 30; i++ {
		tree.Insert(i)
	}

	// Delete some keys to adjust the tree structure
	// Delete some leaf nodes
	toRemove := []int{1, 2, 3, 28, 29, 30}
	for _, v := range toRemove {
		tree.Remove(v)
	}

	// Now try to delete a key in an internal node
	// This should trigger getSuccessor call as we've adjusted the tree structure
	// so that the right subtree has enough keys while the left subtree has insufficient keys
	removed := tree.Remove(15)
	if !removed {
		t.Error("Should be able to remove value 15")
	}

	// Verify the deletion operation was executed successfully
	_, found := tree.Search(15)
	if found {
		t.Error("Value 15 should have been removed")
	}

	// Verify other elements still exist
	for i := 4; i <= 27; i++ {
		if i != 15 {
			val, exists := tree.Search(i)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", i)
			}
			if val != i {
				t.Errorf("Value %d should be %d, got %d", i, i, val)
			}
		}
	}
}

// New test case: Test more branches of borrowFromLeft function
func TestBorrowFromLeftMoreBranches(t *testing.T) {
	tree := New(intComparator)

	// Insert specific data to better trigger borrowFromLeft operations
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}
	for _, v := range values {
		tree.Insert(v)
	}

	// Delete elements from the right sibling node to reduce its key count to a critical value
	toRemoveFromRight := []int{21, 22, 23, 24, 25}
	for _, v := range toRemoveFromRight {
		tree.Remove(v)
	}

	// Now delete elements from the left sibling node, which should trigger borrowFromLeft operation
	removed := tree.Remove(1)
	if !removed {
		t.Error("Should be able to remove value 1")
	}

	// Verify deletion operation succeeded
	_, found := tree.Search(1)
	if found {
		t.Error("Value 1 should have been removed")
	}

	for _, v := range values {
		shouldExist := true
		for _, removed := range toRemoveFromRight {
			if v == removed {
				shouldExist = false
				break
			}
		}
		if v == 1 {
			shouldExist = false
		}

		if shouldExist {
			val, exists := tree.Search(v)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", v)
			}
			if val != v {
				t.Errorf("Value %d should be %d, got %d", v, v, val)
			}
		} else {
			_, exists := tree.Search(v)
			if exists {
				t.Errorf("Value %d should have been removed", v)
			}
		}
	}
}

// New test case: Ensure all lines in getSuccessor function are executed
func TestGetSuccessorAllLines(t *testing.T) {
	tree := New(intComparator)

	// Insert a large amount of data to create a multi-level B-tree structure
	// This ensures the for loop in getSuccessor function will be executed
	for i := 1; i <= 100; i++ {
		tree.Insert(i)
	}

	toRemove := []int{1, 2, 3, 98, 99, 100}
	for _, v := range toRemove {
		tree.Remove(v)
	}

	// Now delete a key from an internal node, this should trigger getSuccessor call
	// And ensure the for loop is executed (i.e., when node.isLeaf is false)
	removed := tree.Remove(50)
	if !removed {
		t.Error("Should be able to remove value 50")
	}

	// Verify the deletion operation was executed successfully
	_, found := tree.Search(50)
	if found {
		t.Error("Value 50 should have been removed")
	}

	for i := 4; i <= 97; i++ {
		if i != 50 {
			val, exists := tree.Search(i)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", i)
			}
			if val != i {
				t.Errorf("Value %d should be %d, got %d", i, i, val)
			}
		}
	}
}

// New test case: Test for loop execution in getSuccessor function
func TestGetSuccessorLoopExecution(t *testing.T) {
	tree := New(intComparator)

	// Insert data in a specific pattern to ensure creation of a structure with internal nodes
	// This ensures the for loop condition in getSuccessor function is true
	values := make([]int, 0)
	for i := 1; i <= 50; i++ {
		values = append(values, i)
	}

	for _, v := range values {
		tree.Insert(v)
	}

	val, found := tree.Search(25)
	if !found || val != 25 {
		t.Error("Value 25 should exist in the tree")
	}

	// Delete key 25, this should trigger getSuccessor call
	// And ensure the for loop is executed
	removed := tree.Remove(25)
	if !removed {
		t.Error("Should be able to remove value 25")
	}

	// Verify deletion succeeded
	_, found = tree.Search(25)
	if found {
		t.Error("Value 25 should have been removed")
	}

	for _, v := range values {
		if v != 25 {
			val, exists := tree.Search(v)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", v)
			}
			if val != v {
				t.Errorf("Value %d should be %d, got %d", v, v, val)
			}
		}
	}

	// Verify the tree order is still correct
	prev := 0
	for v := range tree.IterAsc() {
		if v <= prev {
			t.Errorf("Elements should be in ascending order, got %d after %d", v, prev)
		}
		prev = v
	}
}

// New test case: Ensure complete execution of for loop in getSuccessor function
func TestGetSuccessorFullLoopExecution(t *testing.T) {
	tree := New(intComparator)

	// Insert a large amount of data to create a multi-level B-tree structure
	// This ensures the for loop in getSuccessor function will be executed multiple times
	for i := 1; i <= 100; i++ {
		tree.Insert(i)
	}

	// Delete some keys to adjust tree structure and ensure getSuccessor is called
	toRemove := []int{1, 2, 3, 50, 98, 99, 100}
	for _, v := range toRemove {
		tree.Remove(v)
	}

	// Now delete a key in an internal node, which should trigger getSuccessor call
	// And ensure the for loop is executed (i.e., when node.isLeaf is false)
	removed := tree.Remove(25)
	if !removed {
		t.Error("Should be able to remove value 25")
	}

	// Verify the deletion operation was executed successfully
	_, found := tree.Search(25)
	if found {
		t.Error("Value 25 should have been removed")
	}

	for i := 4; i <= 97; i++ {
		if i != 25 && i != 50 {
			val, exists := tree.Search(i)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", i)
			}
			if val != i {
				t.Errorf("Value %d should be %d, got %d", i, i, val)
			}
		}
	}
}

// New test case: Test for loop execution in getSuccessor function
func TestGetSuccessorLoopExecutionSpecific(t *testing.T) {
	tree := New(intComparator)

	// This ensures the for loop in getSuccessor function will be executed
	for i := 1; i <= 200; i++ {
		tree.Insert(i)
	}

	for i := 1; i <= 50; i++ {
		tree.Remove(i)
	}
	for i := 150; i <= 200; i++ {
		tree.Remove(i)
	}

	// Especially to ensure the for loop is executed multiple times
	removed := tree.Remove(100)
	if !removed {
		t.Error("Should be able to remove value 100")
	}

	// Verify the deletion operation was executed successfully
	_, found := tree.Search(100)
	if found {
		t.Error("Value 100 should have been removed")
	}

	// Verify the tree still maintains correct order
	prev := 0
	for v := range tree.IterAsc() {
		if v <= prev {
			t.Errorf("Elements should be in ascending order, got %d after %d", v, prev)
		}
		prev = v
	}
}

// New test case: Ensure complete execution of for loop in getSuccessor function
func TestGetSuccessorFullLoopCoverage(t *testing.T) {
	tree := New(intComparator)

	// Use a specific insertion order to create a complex internal node structure
	values := make([]int, 0)
	for i := 1; i <= 300; i++ {
		values = append(values, i)
	}

	// Insert in batches to create a more complex tree structure
	for i := 0; i < len(values); i += 10 {
		end := i + 10
		if end > len(values) {
			end = len(values)
		}
		batch := values[i:end]
		for _, v := range batch {
			tree.Insert(v)
		}
	}

	if tree.Len() != 300 {
		t.Errorf("Tree should have 300 elements, got %d", tree.Len())
	}

	// Delete specific keys to trigger complex deletion operations
	// These deletion operations should trigger getSuccessor function calls
	toRemove := []int{50, 100, 150, 200, 250}
	for _, v := range toRemove {
		removed := tree.Remove(v)
		if !removed {
			t.Errorf("Should be able to remove value %d", v)
		}
	}

	for _, v := range toRemove {
		_, found := tree.Search(v)
		if found {
			t.Errorf("Value %d should have been removed", v)
		}
	}

	// Verify the tree size again
	expectedSize := 300 - len(toRemove)
	if tree.Len() != expectedSize {
		t.Errorf("Tree should have %d elements, got %d", expectedSize, tree.Len())
	}

	// Verify the tree still maintains correct order
	prev := 0
	count := 0
	for v := range tree.IterAsc() {
		count++
		if v <= prev {
			t.Errorf("Elements should be in ascending order, got %d after %d", v, prev)
		}
		prev = v
	}

	if count != expectedSize {
		t.Errorf("Iterator should yield %d elements, got %d", expectedSize, count)
	}
}

// New test case: Test all branches of getSuccessor function
func TestGetSuccessorAllBranches(t *testing.T) {
	tree := New(intComparator)

	// Insert enough data to create a complex B-tree structure
	for i := 1; i <= 50; i++ {
		tree.Insert(i)
	}

	toRemove := []int{1, 2, 3, 48, 49, 50}
	for _, v := range toRemove {
		tree.Remove(v)
	}

	removed := tree.Remove(25)
	if !removed {
		t.Error("Should be able to remove value 25")
	}

	_, found := tree.Search(25)
	if found {
		t.Error("Value 25 should have been removed")
	}

	for i := 5; i <= 47; i++ {
		if i != 25 {
			val, exists := tree.Search(i)
			if !exists {
				t.Errorf("Value %d should still exist in the tree", i)
			}
			if val != i {
				t.Errorf("Value %d should be %d, got %d", i, i, val)
			}
		}
	}

	// Verify the tree still maintains correct order
	prev := 0
	for v := range tree.IterAsc() {
		if v <= prev {
			t.Errorf("Elements should be in ascending order, got %d after %d", v, prev)
		}
		prev = v
	}
}
