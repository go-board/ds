package priorityqueue

import (
	"testing"
)

// TestEmptyMinHeap tests empty min heap
func TestEmptyMinHeap(t *testing.T) {
	pq := NewMinOrdered[int]()

	// Test empty queue properties
	if !pq.IsEmpty() {
		t.Error("New queue should be empty")
	}

	if pq.Len() != 0 {
		t.Errorf("New queue length should be 0, got %d", pq.Len())
	}

	// Test popping from empty queue
	val, found := pq.Pop()
	if found || val != 0 {
		t.Error("Pop from empty queue should return false and zero value")
	}

	// Test peeking from empty queue
	val, found = pq.Peek()
	if found || val != 0 {
		t.Error("Peek from empty queue should return false and zero value")
	}
}

// TestMinHeapOperations tests min heap basic operations
func TestMinHeapOperations(t *testing.T) {
	pq := NewMinOrdered[int]()

	// Insert elements
	pq.Push(5)
	pq.Push(3)
	pq.Push(7)
	pq.Push(1)
	pq.Push(9)

	// Verify queue length
	if pq.Len() != 5 {
		t.Errorf("Queue length should be 5, got %d", pq.Len())
	}

	// Verify queue is not empty
	if pq.IsEmpty() {
		t.Error("Queue should not be empty after push operations")
	}

	// Verify the top element (minimum value)
	val, found := pq.Peek()
	if !found || val != 1 {
		t.Errorf("Min heap peek should return 1, got %d, found: %v", val, found)
	}

	// Pop elements, verify order
	expected := []int{1, 3, 5, 7, 9}
	for i, exp := range expected {
		val, found := pq.Pop()
		if !found || val != exp {
			t.Errorf("Pop #%d expected %d, got %d", i, exp, val)
		}
	}

	// Verify queue is empty
	if !pq.IsEmpty() {
		t.Error("Queue should be empty after popping all elements")
	}
}

// TestMaxHeapOperations tests max heap basic operations
func TestMaxHeapOperations(t *testing.T) {
	pq := NewMaxOrdered[int]()

	// Insert elements
	pq.Push(5)
	pq.Push(3)
	pq.Push(7)
	pq.Push(1)
	pq.Push(9)

	// Verify queue length
	if pq.Len() != 5 {
		t.Errorf("Queue length should be 5, got %d", pq.Len())
	}

	// Verify the top element (maximum value)
	val, found := pq.Peek()
	if !found || val != 9 {
		t.Errorf("Max heap peek should return 9, got %d, found: %v", val, found)
	}

	// Pop elements, verify order
	expected := []int{9, 7, 5, 3, 1}
	for i, exp := range expected {
		val, found := pq.Pop()
		if !found || val != exp {
			t.Errorf("Pop #%d expected %d, got %d", i, exp, val)
		}
	}
}

// TestStringPriorityQueue tests string type priority queue
func TestStringPriorityQueue(t *testing.T) {
	pq := NewMinOrdered[string]()

	// Insert strings
	strings := []string{"banana", "apple", "cherry", "date", "elderberry"}
	for _, s := range strings {
		pq.Push(s)
	}

	// Verify pop order
	expected := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for i, exp := range expected {
		val, found := pq.Pop()
		if !found || val != exp {
			t.Errorf("String pop #%d expected %s, got %s", i, exp, val)
		}
	}
}

// TestClearOperation tests Clear method
func TestClearOperation(t *testing.T) {
	pq := NewMinOrdered[int]()

	// Insert some elements
	for i := range 10 {
		pq.Push(i)
	}

	// Clear the queue
	pq.Clear()

	// Verify queue is cleared
	if !pq.IsEmpty() {
		t.Error("Queue should be empty after Clear")
	}

	if pq.Len() != 0 {
		t.Errorf("Queue length should be 0 after Clear, got %d", pq.Len())
	}
}

// TestMixedOperations tests mixed operations
func TestMixedOperations(t *testing.T) {
	pq := NewMinOrdered[int]()

	// Test mixed push and pop operations
	pq.Push(10)
	pq.Push(5)
	pq.Push(15)

	val, _ := pq.Pop() // Should return 5
	if val != 5 {
		t.Errorf("Expected 5, got %d", val)
	}

	pq.Push(3)
	pq.Push(7)

	val, _ = pq.Pop() // Should return 3
	if val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}

	val, _ = pq.Pop() // Should return 7
	if val != 7 {
		t.Errorf("Expected 7, got %d", val)
	}

	val, _ = pq.Pop() // Should return 10
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	val, _ = pq.Pop() // Should return 15
	if val != 15 {
		t.Errorf("Expected 15, got %d", val)
	}
}

// TestLargeData tests large amount of data
func TestLargeData(t *testing.T) {
	pq := NewMinOrdered[int]()
	const size = 1000

	// Insert large number of elements
	for i := size - 1; i >= 0; i-- {
		pq.Push(i)
	}

	// Verify queue length
	if pq.Len() != size {
		t.Errorf("Queue length should be %d, got %d", size, pq.Len())
	}

	// Verify pop order
	for i := range size {
		val, found := pq.Pop()
		if !found || val != i {
			t.Errorf("Pop #%d expected %d, got %d", i, i, val)
			break
		}
	}
}
