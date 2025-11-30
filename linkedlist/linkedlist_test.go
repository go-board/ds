package linkedlist

import (
	"slices"
	"testing"
)

func collect[T any](seq []T, value T) []T {
	return append(seq, value)
}

func TestLinkedListPushPopAndPeek(t *testing.T) {
	list := New[int]()
	if !list.IsEmpty() || list.Len() != 0 {
		t.Fatal("New linked list should be empty")
	}

	list.PushBack(1)
	list.PushFront(0)
	list.PushBack(2)

	if list.Len() != 3 {
		t.Fatalf("Length should be 3, got %d", list.Len())
	}

	front, ok := list.Front()
	if !ok || front != 0 {
		t.Fatalf("Front should return 0, got %v %v", front, ok)
	}
	if v, ok := list.FrontMut(); !ok || *v != 0 {
		t.Fatalf("FrontMut returned unexpected values: %v %v", v, ok)
	} else {
		*v = -1
	}

	back, ok := list.Back()
	if !ok || back != 2 {
		t.Fatalf("Back should return 2, got %v %v", back, ok)
	}
	if v, ok := list.BackMut(); !ok || *v != 2 {
		t.Fatalf("BackMut returned unexpected values: %v %v", v, ok)
	} else {
		*v = 3
	}

	if val, ok := list.PopFront(); !ok || val != -1 {
		t.Fatalf("PopFront should get -1, returned %v %v", val, ok)
	}
	if val, ok := list.PopBack(); !ok || val != 3 {
		t.Fatalf("PopBack should get 3, returned %v %v", val, ok)
	}
	if val, ok := list.PopBack(); !ok || val != 1 {
		t.Fatalf("PopBack should get 1, returned %v %v", val, ok)
	}
	if _, ok := list.PopBack(); ok {
		t.Fatal("PopBack on empty list should return false")
	}

	list.PushFront(10)
	list.Clear()
	if !list.IsEmpty() || list.Len() != 0 {
		t.Fatal("List should be empty after Clear")
	}
}

func TestLinkedListIterators(t *testing.T) {
	list := New[int]()
	for i := 0; i < 5; i++ {
		list.PushBack(i)
	}

	var iterValues []int
	for v := range list.Iter() {
		iterValues = collect(iterValues, v)
	}
	if !slices.Equal(iterValues, []int{0, 1, 2, 3, 4}) {
		t.Fatalf("Iter order unexpected: %#v", iterValues)
	}

	for v := range list.IterMut() {
		*v *= 10
	}
	var updated []int
	for v := range list.Iter() {
		updated = append(updated, v)
	}
	if !slices.Equal(updated, []int{0, 10, 20, 30, 40}) {
		t.Fatalf("IterMut modification didn't take effect: %#v", updated)
	}
}

func TestLinkedListAppendAndRetain(t *testing.T) {
	a := New[int]()
	b := New[int]()
	for i := 0; i < 3; i++ {
		a.PushBack(i)
		b.PushBack(i + 3)
	}

	a.Append(b)
	if a.Len() != 6 {
		t.Fatalf("Length after append should be 6, got %d", a.Len())
	}
	if !b.IsEmpty() {
		t.Fatal("Appended list should be cleared")
	}

	a.Retain(func(v int) bool { return v%2 == 0 })
	var values []int
	for v := range a.Iter() {
		values = append(values, v)
	}
	if !slices.Equal(values, []int{0, 2, 4}) {
		t.Fatalf("Retain result unexpected: %#v", values)
	}
}

func TestLinkedListEdgeCases(t *testing.T) {
	list := New[int]()
	if _, ok := list.Front(); ok {
		t.Fatal("Front on empty list should not return an element")
	}
	if _, ok := list.Back(); ok {
		t.Fatal("Back on empty list should not return an element")
	}

	list.PushBack(1)
	if val, ok := list.PopFront(); !ok || val != 1 {
		t.Fatalf("PopFront on single element should return 1, got %v %v", val, ok)
	}

	list.PushFront(2)
	if val, ok := list.PopBack(); !ok || val != 2 {
		t.Fatalf("PopBack on single element should return 2, got %v %v", val, ok)
	}

	other := New[int]()
	other.PushBack(10)
	list.Append(other)
	if list.Len() != 1 || !other.IsEmpty() {
		t.Fatal("Empty list appending non-empty list should directly take over nodes")
	}

	empty := New[int]()
	list.Append(empty) // other is empty, early return coverage

	list.PushBack(11)
	list.Retain(func(v int) bool { return v != 10 })
	if front, _ := list.Front(); front != 11 {
		t.Fatalf("Retain should remove head node, got %d", front)
	}
}

// New test case: Test BackMut and FrontMut functions for edge cases
func TestBackMutAndFrontMutEdgeCases(t *testing.T) {
	list := New[int]()

	// Test BackMut on empty list
	ptr, ok := list.BackMut()
	if ok || ptr != nil {
		t.Errorf("BackMut on empty list should return nil, false, got %v, %v", ptr, ok)
	}

	// Test FrontMut on empty list
	ptr, ok = list.FrontMut()
	if ok || ptr != nil {
		t.Errorf("FrontMut on empty list should return nil, false, got %v, %v", ptr, ok)
	}

	// Test BackMut after adding elements
	list.PushBack(42)
	ptr, ok = list.BackMut()
	if !ok || ptr == nil {
		t.Errorf("BackMut on non-empty list should return pointer, true, got %v, %v", ptr, ok)
	} else if *ptr != 42 {
		t.Errorf("BackMut should return pointer to 42, got %d", *ptr)
	}

	// Modify the value obtained through BackMut
	*ptr = 100
	val, ok := list.Back()
	if !ok || val != 100 {
		t.Errorf("Back should return 100 after modification, got %d", val)
	}

	// Test FrontMut after adding elements
	list.PushFront(10)
	ptr, ok = list.FrontMut()
	if !ok || ptr == nil {
		t.Errorf("FrontMut on non-empty list should return pointer, true, got %v, %v", ptr, ok)
	} else if *ptr != 10 {
		t.Errorf("FrontMut should return pointer to 10, got %d", *ptr)
	}

	// Modify the value obtained through FrontMut
	*ptr = 20
	val, ok = list.Front()
	if !ok || val != 20 {
		t.Errorf("Front should return 20 after modification, got %d", val)
	}
}

// New test case: Test iterator edge cases
func TestIteratorsEdgeCases(t *testing.T) {
	list := New[int]()

	// Test Iter on empty list
	count := 0
	for range list.Iter() {
		count++
	}
	if count != 0 {
		t.Errorf("Iter on empty list should yield 0 elements, got %d", count)
	}

	// Test IterMut on empty list
	count = 0
	for range list.IterMut() {
		count++
	}
	if count != 0 {
		t.Errorf("IterMut on empty list should yield 0 elements, got %d", count)
	}

	// Test Iter on single-element list
	list.PushBack(1)
	values := make([]int, 0)
	for v := range list.Iter() {
		values = append(values, v)
	}
	if len(values) != 1 || values[0] != 1 {
		t.Errorf("Iter on single-element list should yield [1], got %v", values)
	}

	// Test IterMut on single-element list
	values = values[:0] // Clear the slice
	for v := range list.IterMut() {
		values = append(values, *v)
	}
	if len(values) != 1 || values[0] != 1 {
		t.Errorf("IterMut on single-element list should yield [1], got %v", values)
	}

	// Test IterMut for value modification
	for v := range list.IterMut() {
		*v = 42
	}
	val, ok := list.Front()
	if !ok || val != 42 {
		t.Errorf("Front should return 42 after IterMut modification, got %d", val)
	}
}

// New test case: Test PopFront edge cases
func TestPopFrontEdgeCases(t *testing.T) {
	list := New[int]()

	// Test PopFront on empty list
	val, ok := list.PopFront()
	if ok {
		t.Errorf("PopFront on empty list should return false, got %v with value %d", ok, val)
	}

	// Test PopFront on single-element list
	list.PushBack(100)
	val, ok = list.PopFront()
	if !ok || val != 100 {
		t.Errorf("PopFront on single-element list should return 100, true, got %d, %v", val, ok)
	}

	// Verify list is now empty
	if !list.IsEmpty() || list.Len() != 0 {
		t.Error("List should be empty after PopFront on single-element list")
	}

	// Test PopFront on empty list again
	val, ok = list.PopFront()
	if ok {
		t.Errorf("PopFront on empty list should return false, got %v with value %d", ok, val)
	}
}

// New test case: Test Append function with more cases
func TestAppendMoreCases(t *testing.T) {
	// Test appending empty list to non-empty list
	list1 := New[int]()
	list1.PushBack(1)
	list1.PushBack(2)

	list2 := New[int]()

	// Save original state
	originalLen := list1.Len()

	// Append empty list to non-empty list
	list1.Append(list2)

	// Verify list1 hasn't changed
	if list1.Len() != originalLen {
		t.Errorf("Appending empty list should not change the original list length, expected %d, got %d", originalLen, list1.Len())
	}

	val1, _ := list1.Front()
	val2, _ := list1.Back()
	if val1 != 1 || val2 != 2 {
		t.Errorf("Appending empty list should not change the original list elements, expected 1,2 got %d,%d", val1, val2)
	}

	// Verify list2 is still empty
	if !list2.IsEmpty() {
		t.Error("Appended empty list should still be empty")
	}
}

// New test case: Test iterator early termination
func TestIteratorsEarlyTermination(t *testing.T) {
	list := New[int]()

	// Add multiple elements
	for i := 0; i < 5; i++ {
		list.PushBack(i)
	}

	// Test early termination with Iter
	count := 0
	for range list.Iter() {
		count++
		if count == 3 {
			break // Early termination
		}
	}
	if count != 3 {
		t.Errorf("Iter should yield 3 elements before break, got %d", count)
	}

	// Test early termination with IterMut
	count = 0
	for range list.IterMut() {
		count++
		if count == 2 {
			break // Early termination
		}
	}
	if count != 2 {
		t.Errorf("IterMut should yield 2 elements before break, got %d", count)
	}

	// Verify list structure is not corrupted
	if list.Len() != 5 {
		t.Errorf("List length should remain 5 after early termination, got %d", list.Len())
	}
}
