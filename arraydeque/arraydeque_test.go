package arraydeque

import (
	"testing"
)

func TestNew(t *testing.T) {
	// Test creating a new empty deque
	dq := New[int]()
	if dq == nil {
		t.Fatal("New should return a non-nil pointer")
	}
	if !dq.IsEmpty() {
		t.Fatal("New should return an empty deque")
	}
	if dq.Len() != 0 {
		t.Fatalf("New should return a deque with length 0, got %d", dq.Len())
	}
}

func TestPushBack(t *testing.T) {
	dq := New[int]()

	// Test pushing one element
	dq.PushBack(1)
	if dq.Len() != 1 {
		t.Fatalf("PushBack should increase length to 1, got %d", dq.Len())
	}
	if val, ok := dq.Back(); !ok || val != 1 {
		t.Fatalf("PushBack should set back element to 1, got %v, %v", val, ok)
	}

	// Test pushing multiple elements
	dq.PushBack(2)
	dq.PushBack(3)
	if dq.Len() != 3 {
		t.Fatalf("PushBack should increase length to 3, got %d", dq.Len())
	}
	if val, ok := dq.Back(); !ok || val != 3 {
		t.Fatalf("PushBack should set back element to 3, got %v, %v", val, ok)
	}
}

func TestPushFront(t *testing.T) {
	dq := New[int]()

	// Test pushing one element to front
	dq.PushFront(1)
	if dq.Len() != 1 {
		t.Fatalf("PushFront should increase length to 1, got %d", dq.Len())
	}
	if val, ok := dq.Front(); !ok || val != 1 {
		t.Fatalf("PushFront should set front element to 1, got %v, %v", val, ok)
	}

	// Test pushing multiple elements to front
	dq.PushFront(2)
	dq.PushFront(3)
	if dq.Len() != 3 {
		t.Fatalf("PushFront should increase length to 3, got %d", dq.Len())
	}
	if val, ok := dq.Front(); !ok || val != 3 {
		t.Fatalf("PushFront should set front element to 3, got %v, %v", val, ok)
	}
}

func TestPopBack(t *testing.T) {
	dq := New[int]()

	// Test popping from empty deque
	if _, ok := dq.PopBack(); ok {
		t.Fatal("PopBack should return false for empty deque")
	}

	// Test popping after pushing
	dq.PushBack(1)
	dq.PushBack(2)
	dq.PushBack(3)

	if val, ok := dq.PopBack(); !ok || val != 3 {
		t.Fatalf("PopBack should return 3, got %v, %v", val, ok)
	}
	if dq.Len() != 2 {
		t.Fatalf("PopBack should decrease length to 2, got %d", dq.Len())
	}

	if val, ok := dq.PopBack(); !ok || val != 2 {
		t.Fatalf("PopBack should return 2, got %v, %v", val, ok)
	}
	if dq.Len() != 1 {
		t.Fatalf("PopBack should decrease length to 1, got %d", dq.Len())
	}

	if val, ok := dq.PopBack(); !ok || val != 1 {
		t.Fatalf("PopBack should return 1, got %v, %v", val, ok)
	}
	if dq.Len() != 0 {
		t.Fatalf("PopBack should decrease length to 0, got %d", dq.Len())
	}
}

func TestPopFront(t *testing.T) {
	dq := New[int]()

	// Test popping from empty deque
	if _, ok := dq.PopFront(); ok {
		t.Fatal("PopFront should return false for empty deque")
	}

	// Test popping after pushing to front
	dq.PushFront(1)
	dq.PushFront(2)
	dq.PushFront(3)

	if val, ok := dq.PopFront(); !ok || val != 3 {
		t.Fatalf("PopFront should return 3, got %v, %v", val, ok)
	}
	if dq.Len() != 2 {
		t.Fatalf("PopFront should decrease length to 2, got %d", dq.Len())
	}

	if val, ok := dq.PopFront(); !ok || val != 2 {
		t.Fatalf("PopFront should return 2, got %v, %v", val, ok)
	}
	if dq.Len() != 1 {
		t.Fatalf("PopFront should decrease length to 1, got %d", dq.Len())
	}

	if val, ok := dq.PopFront(); !ok || val != 1 {
		t.Fatalf("PopFront should return 1, got %v, %v", val, ok)
	}
	if dq.Len() != 0 {
		t.Fatalf("PopFront should decrease length to 0, got %d", dq.Len())
	}
}

func TestFrontAndBack(t *testing.T) {
	dq := New[int]()

	// Test on empty deque
	if _, ok := dq.Front(); ok {
		t.Fatal("Front should return false for empty deque")
	}
	if _, ok := dq.Back(); ok {
		t.Fatal("Back should return false for empty deque")
	}

	// Test after pushing to back
	dq.PushBack(1)
	if val, ok := dq.Front(); !ok || val != 1 {
		t.Fatalf("Front should return 1, got %v, %v", val, ok)
	}
	if val, ok := dq.Back(); !ok || val != 1 {
		t.Fatalf("Back should return 1, got %v, %v", val, ok)
	}

	// Test after pushing to front
	dq.PushFront(2)
	if val, ok := dq.Front(); !ok || val != 2 {
		t.Fatalf("Front should return 2, got %v, %v", val, ok)
	}
	if val, ok := dq.Back(); !ok || val != 1 {
		t.Fatalf("Back should return 1, got %v, %v", val, ok)
	}

	// Test after pushing more elements
	dq.PushBack(3)
	if val, ok := dq.Front(); !ok || val != 2 {
		t.Fatalf("Front should return 2, got %v, %v", val, ok)
	}
	if val, ok := dq.Back(); !ok || val != 3 {
		t.Fatalf("Back should return 3, got %v, %v", val, ok)
	}
}

// TestGetAndGetMut function removed as Get and GetMut methods are no longer supported in queue implementation

func TestLenAndIsEmpty(t *testing.T) {
	dq := New[int]()

	// Test on empty deque
	if dq.Len() != 0 {
		t.Fatalf("Len should return 0 for empty deque, got %d", dq.Len())
	}
	if !dq.IsEmpty() {
		t.Fatal("IsEmpty should return true for empty deque")
	}

	// Test after pushing one element
	dq.PushBack(1)
	if dq.Len() != 1 {
		t.Fatalf("Len should return 1 after pushing one element, got %d", dq.Len())
	}
	if dq.IsEmpty() {
		t.Fatal("IsEmpty should return false after pushing one element")
	}

	// Test after pushing multiple elements
	dq.PushBack(2)
	dq.PushBack(3)
	if dq.Len() != 3 {
		t.Fatalf("Len should return 3 after pushing three elements, got %d", dq.Len())
	}
	if dq.IsEmpty() {
		t.Fatal("IsEmpty should return false after pushing three elements")
	}

	// Test after popping all elements
	dq.PopFront()
	dq.PopFront()
	dq.PopFront()
	if dq.Len() != 0 {
		t.Fatalf("Len should return 0 after popping all elements, got %d", dq.Len())
	}
	if !dq.IsEmpty() {
		t.Fatal("IsEmpty should return true after popping all elements")
	}
}

func TestClear(t *testing.T) {
	dq := New[int]()
	dq.PushBack(1)
	dq.PushBack(2)
	dq.PushBack(3)

	// Test clearing non-empty deque
	dq.Clear()
	if dq.Len() != 0 {
		t.Fatalf("Clear should set length to 0, got %d", dq.Len())
	}
	if !dq.IsEmpty() {
		t.Fatal("Clear should make deque empty")
	}
	if _, ok := dq.Front(); ok {
		t.Fatal("Clear should remove all elements, Front should return false")
	}
	if _, ok := dq.Back(); ok {
		t.Fatal("Clear should remove all elements, Back should return false")
	}

	// Test clearing empty deque (should not panic)
	dq.Clear()
	if dq.Len() != 0 {
		t.Fatalf("Clear on empty deque should keep length 0, got %d", dq.Len())
	}
}

func TestClone(t *testing.T) {
	dq := New[int]()
	dq.PushBack(1)
	dq.PushBack(2)
	dq.PushBack(3)

	// Test cloning non-empty deque
	clone := dq.Clone()
	if clone == nil {
		t.Fatal("Clone should return a non-nil pointer")
	}
	if clone == dq {
		t.Fatal("Clone should return a different pointer")
	}
	if clone.Len() != dq.Len() {
		t.Fatalf("Clone should have same length, got %d, expected %d", clone.Len(), dq.Len())
	}

	// Test that clone is a deep copy
	if val, ok := clone.PopBack(); !ok || val != 3 {
		t.Fatalf("Clone should have same elements, got %v, %v", val, ok)
	}
	// Original deque should still have all elements
	if dq.Len() != 3 {
		t.Fatalf("Original deque should not be modified by clone operations, got length %d", dq.Len())
	}

	// Test cloning empty deque
	empty := New[int]()
	emptyClone := empty.Clone()
	if emptyClone == nil {
		t.Fatal("Clone should return a non-nil pointer for empty deque")
	}
	if !emptyClone.IsEmpty() {
		t.Fatal("Clone of empty deque should be empty")
	}
}

func TestReserveAndShrinkToFit(t *testing.T) {
	dq := New[int]()

	// Test Reserve with larger capacity
	dq.Reserve(100)
	if dq.Capacity() < 100 {
		t.Fatalf("Reserve should set capacity to at least 100, got %d", dq.Capacity())
	}

	// Test Reserve with smaller capacity (should not change)
	currentCap := dq.Capacity()
	dq.Reserve(50)
	if dq.Capacity() != currentCap {
		t.Fatalf("Reserve with smaller capacity should not change capacity, got %d, expected %d", dq.Capacity(), currentCap)
	}

	// Test ShrinkToFit
	dq.PushBack(1)
	dq.PushBack(2)
	dq.ShrinkToFit()
	if dq.Capacity() < dq.Len() {
		t.Fatalf("ShrinkToFit should keep capacity at least length, got cap %d, len %d", dq.Capacity(), dq.Len())
	}

	// Test ShrinkToFit on empty deque
	dq.Clear()
	dq.ShrinkToFit()
	if dq.Capacity() == 0 {
		t.Fatal("ShrinkToFit on empty deque should not set capacity to 0")
	}
}

func TestIter(t *testing.T) {
	dq := New[int]()
	dq.PushBack(1)
	dq.PushBack(2)
	dq.PushBack(3)

	// Test Iter
	var result []int
	for val := range dq.Iter() {
		result = append(result, val)
	}
	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Fatalf("Iter should yield %d elements, got %d", len(expected), len(result))
	}
	for i, val := range result {
		if val != expected[i] {
			t.Fatalf("Iter should yield %v, got %v at index %d", expected[i], val, i)
		}
	}

	// Test Iter on empty deque
	empty := New[int]()
	count := 0
	for range empty.Iter() {
		count++
	}
	if count != 0 {
		t.Fatalf("Iter on empty deque should yield 0 elements, got %d", count)
	}
}

func TestIterMut(t *testing.T) {
	dq := New[int]()
	dq.PushBack(1)
	dq.PushBack(2)
	dq.PushBack(3)

	// Test IterMut
	for ptr := range dq.IterMut() {
		*ptr *= 2
	}

	var result []int
	for val := range dq.Iter() {
		result = append(result, val)
	}
	expected := []int{2, 4, 6}
	if len(result) != len(expected) {
		t.Fatalf("IterMut should allow modifying all elements, got %d elements, expected %d", len(result), len(expected))
	}
	for i, val := range result {
		if val != expected[i] {
			t.Fatalf("IterMut should modify elements correctly, got %v, expected %v at index %d", val, expected[i], i)
		}
	}
}

func TestExtend(t *testing.T) {
	dq := New[int]()
	dq.PushBack(1)
	dq.PushBack(2)

	// Test Extend with slice iterator
	toExtend := []int{3, 4, 5}
	dq.Extend(func(yield func(int) bool) {
		for _, val := range toExtend {
			if !yield(val) {
				return
			}
		}
	})

	if dq.Len() != 5 {
		t.Fatalf("Extend should increase length to 5, got %d", dq.Len())
	}

	var result []int
	for val := range dq.Iter() {
		result = append(result, val)
	}
	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Fatalf("Extend should add all elements, got %d elements, expected %d", len(result), len(expected))
	}
	for i, val := range result {
		if val != expected[i] {
			t.Fatalf("Extend should add elements in correct order, got %v, expected %v at index %d", val, expected[i], i)
		}
	}
}

func TestEdgeCases(t *testing.T) {
	dq := New[int]()

	// Test pushing and popping alternately
	for i := 0; i < 100; i++ {
		dq.PushBack(i)
		dq.PopFront()
	}
	if !dq.IsEmpty() {
		t.Fatal("Alternate push and pop should result in empty deque")
	}

	// Test pushing to front and popping from back alternately
	for i := 0; i < 100; i++ {
		dq.PushFront(i)
		dq.PopBack()
	}
	if !dq.IsEmpty() {
		t.Fatal("Alternate push front and pop back should result in empty deque")
	}

	// Test large number of elements to trigger multiple grow operations
	for i := 0; i < 1000; i++ {
		dq.PushBack(i)
	}
	if dq.Len() != 1000 {
		t.Fatalf("Should handle 1000 elements, got %d", dq.Len())
	}
	for i := 0; i < 1000; i++ {
		if val, ok := dq.PopFront(); !ok || val != i {
			t.Fatalf("Should pop elements in order, got %v, %v at index %d", val, ok, i)
		}
	}
	if !dq.IsEmpty() {
		t.Fatal("Should be empty after popping all elements")
	}
}

func TestCap(t *testing.T) {
	dq := New[int]()

	// Test initial capacity
	initialCap := dq.Capacity()
	if initialCap == 0 {
		t.Fatal("Initial capacity should be greater than 0")
	}

	// Test capacity after pushing elements
	for i := 0; i < initialCap+1; i++ {
		dq.PushBack(i)
	}
	if dq.Capacity() <= initialCap {
		t.Fatalf("Capacity should increase after exceeding initial capacity, got %d, expected > %d", dq.Capacity(), initialCap)
	}
}

func TestIterBack(t *testing.T) {
	dq := New[int]()
	dq.PushBack(1)
	dq.PushBack(2)
	dq.PushBack(3)

	// Test IterBack (reverse order: 3, 2, 1)
	var result []int
	for val := range dq.IterBack() {
		result = append(result, val)
	}
	expected := []int{3, 2, 1}
	if len(result) != len(expected) {
		t.Fatalf("IterBack should yield %d elements, got %d", len(expected), len(result))
	}
	for i, val := range result {
		if val != expected[i] {
			t.Fatalf("IterBack should yield %v, got %v at index %d", expected[i], val, i)
		}
	}

	// Test IterBack on empty deque
	empty := New[int]()
	count := 0
	for range empty.IterBack() {
		count++
	}
	if count != 0 {
		t.Fatalf("IterBack on empty deque should yield 0 elements, got %d", count)
	}

	// Test with mixed push operations to verify circular buffer handling
	dq2 := New[int]()
	for i := 1; i <= 10; i++ {
		dq2.PushBack(i)
	}
	// Remove first 5 elements to create a circular situation
	for i := 0; i < 5; i++ {
		dq2.PopFront()
	}
	// Add more elements
	for i := 11; i <= 15; i++ {
		dq2.PushBack(i)
	}

	// IterBack should still work correctly with circular buffer
	var result2 []int
	for val := range dq2.IterBack() {
		result2 = append(result2, val)
	}
	expected2 := []int{15, 14, 13, 12, 11, 10, 9, 8, 7, 6}
	if len(result2) != len(expected2) {
		t.Fatalf("IterBack with circular buffer should yield %d elements, got %d", len(expected2), len(result2))
	}
	for i, val := range result2 {
		if val != expected2[i] {
			t.Fatalf("IterBack with circular buffer should yield %v, got %v at index %d", expected2[i], val, i)
		}
	}
}

func TestIterBackMut(t *testing.T) {
	dq := New[int]()
	dq.PushBack(1)
	dq.PushBack(2)
	dq.PushBack(3)

	// Test IterBackMut - modify elements in reverse order
	multiplier := 2
	for ptr := range dq.IterBackMut() {
		*ptr *= multiplier
		multiplier++ // Use different multipliers for different positions
	}

	// Verify modified elements using IterBack (should be 3*2=6, 2*3=6, 1*4=4)
	var result []int
	for val := range dq.IterBack() {
		result = append(result, val)
	}
	expected := []int{6, 6, 4}
	if len(result) != len(expected) {
		t.Fatalf("IterBackMut should allow modifying all elements, got %d elements, expected %d", len(result), len(expected))
	}
	for i, val := range result {
		if val != expected[i] {
			t.Fatalf("IterBackMut should modify elements correctly in reverse order, got %v, expected %v at index %d", val, expected[i], i)
		}
	}

	// Verify using regular Iter (should be 4, 6, 6)
	var resultRegular []int
	for val := range dq.Iter() {
		resultRegular = append(resultRegular, val)
	}
	expectedRegular := []int{4, 6, 6}
	if len(resultRegular) != len(expectedRegular) {
		t.Fatalf("After IterBackMut, Iter should show %d elements, got %d", len(expectedRegular), len(resultRegular))
	}
	for i, val := range resultRegular {
		if val != expectedRegular[i] {
			t.Fatalf("After IterBackMut, Iter should show %v, got %v at index %d", expectedRegular[i], val, i)
		}
	}

	// Test IterBackMut on empty deque (should not panic)
	empty := New[int]()
	count := 0
	for range empty.IterBackMut() {
		count++
	}
	if count != 0 {
		t.Fatalf("IterBackMut on empty deque should yield 0 elements, got %d", count)
	}
}
