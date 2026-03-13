// Package arraystack implements tests for the arraystack package.
package arraystack

import (
	"testing"
)

// TestNew tests the New constructor
func TestNew(t *testing.T) {
	stack := New[int]()
	if stack == nil {
		t.Fatal("New should return a non-nil stack")
	}
	if !stack.IsEmpty() {
		t.Fatal("New stack should be empty")
	}
	if stack.Len() != 0 {
		t.Fatalf("New stack should have length 0, got %d", stack.Len())
	}
}

// TestPush tests the Push method
func TestPush(t *testing.T) {
	stack := New[int]()

	// Push a single element
	stack.Push(5)
	if stack.Len() != 1 {
		t.Fatalf("Stack length should be 1 after one push, got %d", stack.Len())
	}

	// Push multiple elements
	stack.Push(10)
	stack.Push(15)
	if stack.Len() != 3 {
		t.Fatalf("Stack length should be 3 after three pushes, got %d", stack.Len())
	}

	// Test with different types
	stringStack := New[string]()
	stringStack.Push("hello")
	stringStack.Push("world")
	if stringStack.Len() != 2 {
		t.Fatalf("String stack length should be 2, got %d", stringStack.Len())
	}
}

// TestPop tests the Pop method
func TestPop(t *testing.T) {
	stack := New[int]()

	// Test Pop on empty stack
	val, found := stack.Pop()
	if found {
		t.Fatal("Pop on empty stack should return false")
	}
	if val != 0 {
		t.Fatalf("Pop on empty stack should return zero value, got %d", val)
	}

	// Push elements and test Pop
	stack.Push(5)
	stack.Push(10)
	stack.Push(15)

	// Pop first element
	val, found = stack.Pop()
	if !found {
		t.Fatal("Pop should return true for non-empty stack")
	}
	if val != 15 {
		t.Fatalf("Pop should return 15, got %d", val)
	}
	if stack.Len() != 2 {
		t.Fatalf("Stack length should be 2 after one pop, got %d", stack.Len())
	}

	// Pop second element
	val, found = stack.Pop()
	if !found {
		t.Fatal("Pop should return true for non-empty stack")
	}
	if val != 10 {
		t.Fatalf("Pop should return 10, got %d", val)
	}
	if stack.Len() != 1 {
		t.Fatalf("Stack length should be 1 after two pops, got %d", stack.Len())
	}

	// Pop third element
	val, found = stack.Pop()
	if !found {
		t.Fatal("Pop should return true for non-empty stack")
	}
	if val != 5 {
		t.Fatalf("Pop should return 5, got %d", val)
	}
	if !stack.IsEmpty() {
		t.Fatal("Stack should be empty after popping all elements")
	}

	// Test Pop after stack is empty
	_, found = stack.Pop()
	if found {
		t.Fatal("Pop on empty stack should return false")
	}
}

// TestPeek tests the Peek method
func TestPeek(t *testing.T) {
	stack := New[int]()

	// Test Peek on empty stack
	val, found := stack.Peek()
	if found {
		t.Fatal("Peek on empty stack should return false")
	}
	if val != 0 {
		t.Fatalf("Peek on empty stack should return zero value, got %d", val)
	}

	// Push elements and test Peek
	stack.Push(5)
	stack.Push(10)

	// Test Peek doesn't modify the stack
	val, found = stack.Peek()
	if !found {
		t.Fatal("Peek should return true for non-empty stack")
	}
	if val != 10 {
		t.Fatalf("Peek should return 10, got %d", val)
	}
	if stack.Len() != 2 {
		t.Fatalf("Peek should not modify stack length, got %d", stack.Len())
	}

	// Peek after Pop
	_, _ = stack.Pop()
	val, found = stack.Peek()
	if !found {
		t.Fatal("Peek should return true for non-empty stack")
	}
	if val != 5 {
		t.Fatalf("Peek should return 5 after popping 10, got %d", val)
	}
}

// TestClear tests the Clear method
func TestClear(t *testing.T) {
	stack := New[int]()

	// Clear empty stack
	stack.Clear()
	if !stack.IsEmpty() {
		t.Fatal("Clear on empty stack should leave it empty")
	}

	// Push elements and clear
	stack.Push(5)
	stack.Push(10)
	stack.Push(15)
	stack.Clear()

	if !stack.IsEmpty() {
		t.Fatal("Clear should empty the stack")
	}
	if stack.Len() != 0 {
		t.Fatalf("Stack length should be 0 after Clear, got %d", stack.Len())
	}

	// Verify we can use the stack after clearing
	stack.Push(20)
	if stack.Len() != 1 {
		t.Fatalf("Stack should accept new elements after Clear, length got %d", stack.Len())
	}
	val, found := stack.Peek()
	if !found || val != 20 {
		t.Fatalf("Stack should contain 20 after Clear and Push, got %d, found: %v", val, found)
	}
}

// TestIsEmpty tests the IsEmpty method
func TestIsEmpty(t *testing.T) {
	stack := New[int]()

	if !stack.IsEmpty() {
		t.Fatal("New stack should be empty")
	}

	stack.Push(5)
	if stack.IsEmpty() {
		t.Fatal("Stack with elements should not be empty")
	}

	_, _ = stack.Pop()
	if !stack.IsEmpty() {
		t.Fatal("Stack after popping all elements should be empty")
	}
}

// TestIter tests the Iter method (LIFO order)
func TestIter(t *testing.T) {
	stack := New[int]()

	// Test iteration on empty stack
	count := 0
	for range stack.Iter() {
		count++
	}
	if count != 0 {
		t.Fatalf("Iter on empty stack should yield 0 elements, got %d", count)
	}

	// Push elements
	stack.Push(5)
	stack.Push(10)
	stack.Push(15)

	// Collect elements using Iter (should be LIFO order: 15, 10, 5)
	var elements []int
	expected := []int{15, 10, 5}
	for val := range stack.Iter() {
		elements = append(elements, val)
	}

	if len(elements) != len(expected) {
		t.Fatalf("Iter should yield %d elements, got %d", len(expected), len(elements))
	}

	for i, val := range elements {
		if val != expected[i] {
			t.Fatalf("Iter element %d should be %d, got %d", i, expected[i], val)
		}
	}

	// Verify iteration doesn't modify the stack
	if stack.Len() != 3 {
		t.Fatalf("Iter should not modify stack length, got %d", stack.Len())
	}
}

// TestIterMut tests the IterMut method (LIFO order with mutable elements)
func TestIterMut(t *testing.T) {
	stack := New[int]()

	// Test iteration on empty stack
	count := 0
	for range stack.IterMut() {
		count++
	}
	if count != 0 {
		t.Fatalf("IterMut on empty stack should yield 0 elements, got %d", count)
	}

	// Push elements
	stack.Push(5)
	stack.Push(10)
	stack.Push(15)

	// Modify elements using IterMut (should be in LIFO order: 15, 10, 5)
	// Double each value
	expectedAfterModify := []int{10, 20, 30} // 5*2, 10*2, 15*2
	for valPtr := range stack.IterMut() {
		*valPtr *= 2
	}

	// Verify modification using Iter
	var actualValues []int
	for val := range stack.Iter() {
		actualValues = append(actualValues, val)
	}

	// Check values in LIFO order (30, 20, 10)
	if len(actualValues) != len(expectedAfterModify) {
		t.Fatalf("Stack should have %d elements after modification, got %d", len(expectedAfterModify), len(actualValues))
	}

	// Reverse the expected order to match LIFO iteration
	for i, val := range actualValues {
		expectedIndex := len(expectedAfterModify) - 1 - i
		if val != expectedAfterModify[expectedIndex] {
			t.Fatalf("Element at position %d should be %d, got %d", i, expectedAfterModify[expectedIndex], val)
		}
	}

	// Verify stack length is unchanged
	if stack.Len() != 3 {
		t.Fatalf("IterMut should not modify stack length, got %d", stack.Len())
	}
}

// TestClone tests the Clone method
func TestClone(t *testing.T) {
	// Clone empty stack
	stack1 := New[int]()
	stack2 := stack1.Clone()

	if stack2 == nil {
		t.Fatal("Clone should return a non-nil stack")
	}
	if !stack2.IsEmpty() {
		t.Fatal("Clone of empty stack should be empty")
	}

	// Clone non-empty stack
	stack1.Push(5)
	stack1.Push(10)
	stack1.Push(15)
	stack2 = stack1.Clone()

	if stack2 == nil {
		t.Fatal("Clone should return a non-nil stack")
	}
	if stack2.Len() != stack1.Len() {
		t.Fatalf("Clone should have same length as original, got %d vs %d", stack2.Len(), stack1.Len())
	}

	// Verify elements are the same
	for val1 := range stack1.Iter() {
		val2, found := stack2.Pop()
		if !found || val1 != val2 {
			t.Fatalf("Clone elements should match original, got %d vs %d (found: %v)", val2, val1, found)
		}
	}

	// Verify modifying clone doesn't affect original
	stack2.Push(20)
	if stack1.Len() != 3 {
		t.Fatalf("Modifying clone should not affect original, got %d", stack1.Len())
	}
	if stack2.Len() != 1 {
		t.Fatalf("Clone should be modified independently, got %d", stack2.Len())
	}
}

// TestExtend tests the Extend method with iter.Seq
func TestExtend(t *testing.T) {
	stack := New[int]()

	// Test with empty sequence
	emptySeq := func(yield func(int) bool) {
		// No elements to yield
	}
	stack.Extend(emptySeq)
	if !stack.IsEmpty() {
		t.Fatalf("Extend with empty sequence should not change the stack, but stack is not empty")
	}

	// Test with non-empty sequence
	values := []int{1, 2, 3, 4, 5}
	seq := func(yield func(int) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
	stack.Extend(seq)

	// Verify stack length
	if stack.Len() != len(values) {
		t.Fatalf("Stack length should be %d after Extend, got %d", len(values), stack.Len())
	}

	// Verify elements in correct order (LIFO)
	expectedValues := []int{5, 4, 3, 2, 1} // Top to bottom
	var actualValues []int
	for val := range stack.Iter() {
		actualValues = append(actualValues, val)
	}

	for i, val := range actualValues {
		if val != expectedValues[i] {
			t.Fatalf("Element at position %d should be %d, got %d", i, expectedValues[i], val)
		}
	}

	// Test extending with another sequence (should append to top)
	moreValues := []int{6, 7, 8}
	moreSeq := func(yield func(int) bool) {
		for _, v := range moreValues {
			if !yield(v) {
				return
			}
		}
	}
	stack.Extend(moreSeq)

	// Verify new length
	expectedLen := len(values) + len(moreValues)
	if stack.Len() != expectedLen {
		t.Fatalf("Stack length should be %d after second Extend, got %d", expectedLen, stack.Len())
	}

	// Verify new elements are at the top
	combinedExpected := []int{8, 7, 6, 5, 4, 3, 2, 1} // Top to bottom
	actualValues = nil
	for val := range stack.Iter() {
		actualValues = append(actualValues, val)
	}

	for i, val := range actualValues {
		if val != combinedExpected[i] {
			t.Fatalf("Element at position %d should be %d, got %d", i, combinedExpected[i], val)
		}
	}

	// Test with large sequence to trigger resizing
	stack.Clear()
	largeSeq := func(yield func(int) bool) {
		for i := 1; i <= 1000; i++ {
			if !yield(i) {
				return
			}
		}
	}
	stack.Extend(largeSeq)

	if stack.Len() != 1000 {
		t.Fatalf("Stack length should be 1000 after extending with large sequence, got %d", stack.Len())
	}

	// Verify top element is correct
	top, found := stack.Peek()
	if !found || top != 1000 {
		t.Fatalf("Top element should be 1000, got %d (found: %v)", top, found)
	}
}

// TestComplexOperations tests complex operations on the stack
func TestComplexOperations(t *testing.T) {
	stack := New[int]()

	// Push and Pop in various patterns
	stack.Push(1)
	stack.Push(2)
	_, _ = stack.Pop() // Remove 2
	stack.Push(3)
	stack.Push(4)
	_, _ = stack.Pop() // Remove 4
	_, _ = stack.Pop() // Remove 3
	stack.Push(5)

	// Verify final state
	if stack.Len() != 2 {
		t.Fatalf("Stack should have 2 elements after complex operations, got %d", stack.Len())
	}

	// Check elements (should be 5, 1)
	var elements []int
	for val := range stack.Iter() {
		elements = append(elements, val)
	}

	expected := []int{5, 1}
	if len(elements) != len(expected) {
		t.Fatalf("Stack should contain %d elements, got %d", len(expected), len(elements))
	}

	for i, val := range elements {
		if val != expected[i] {
			t.Fatalf("Element %d should be %d, got %d", i, expected[i], val)
		}
	}
}
