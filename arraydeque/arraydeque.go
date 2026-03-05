// Package arraydeque implements a generic double-ended queue data structure based on a linear slice.
//
// ArrayDeque is a data structure that allows efficient addition and removal of elements at both ends,
// internally using a linear slice as the underlying storage.
// All operations have the following time complexity:
// - PushBack, Back, Len, Cap, IsEmpty: O(1) time complexity
// - PushFront, PopFront: O(n) time complexity due to slice shifting
// - PopBack: O(1) time complexity (amortized)
//
// Features:
//   - Support for addition and removal of elements at both ends
//   - Support for capacity management
//   - Support for iterator access
//   - Avoid memory leaks
//
// Example usage:
//
//	// Create an integer double-ended queue
//	dq := New[int]()
//
//	// Add elements from the end
//	dq.PushBack(1)
//	dq.PushBack(2)
//
//	// Add element from the front
//	dq.PushFront(0)
//
//	// Access elements
//	front, _ := dq.Front() // 0
//	back, _ := dq.Back()   // 2
//
//	// Traverse elements
//	for val := range dq.Iter() {
//	    fmt.Println(val) // Output: 0, 1, 2
//	}
//
//	// Remove elements
//	val, _ := dq.PopFront() // 0
//	val, _ := dq.PopBack()  // 2
package arraydeque

// ArrayDeque is a generic double-ended queue implementation based on a linear slice
// Using a slice as the underlying storage with direct indexing
type ArrayDeque[T any] struct {
	// data is the underlying storage slice
	data []T
}

// New creates a new empty double-ended queue
//
// Type parameters:
//   - T: The type of elements in the queue
//
// Returns:
//   - A new ArrayDeque instance
//
// Time complexity: O(1)
func New[T any]() *ArrayDeque[T] {
	return &ArrayDeque[T]{
		data: make([]T, 0, 8),
	}
}

// PushBack adds an element to the end of the queue
//
// Parameters:
//   - value: The element value to add
//
// Time complexity: O(1) (amortized, O(n) when growing)
func (d *ArrayDeque[T]) PushBack(value T) {
	d.data = append(d.data, value)
}

// PushFront adds an element to the front of the queue
//
// Parameters:
//   - value: The element value to add
//
// Time complexity: O(n) (requires shifting all elements)
func (d *ArrayDeque[T]) PushFront(value T) {
	d.data = append([]T{value}, d.data...)
}

// PopBack removes and returns the element from the end of the queue
//
// Returns:
//   - If the queue is not empty, returns the removed element and true
//   - If the queue is empty, returns a zero value and false
//
// Time complexity: O(1)
// Note: This method sets the position of the removed element to zero value to avoid memory leaks
func (d *ArrayDeque[T]) PopBack() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	lastIndex := len(d.data) - 1
	value := d.data[lastIndex]
	d.data = d.data[:lastIndex]
	return value, true
}

// PopFront removes and returns the element from the front of the queue
//
// Returns:
//   - If the queue is not empty, returns the removed element and true
//   - If the queue is empty, returns a zero value and false
//
// Time complexity: O(n) (requires shifting all elements)
// Note: This method sets the position of the removed element to zero value to avoid memory leaks
func (d *ArrayDeque[T]) PopFront() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	value := d.data[0]
	d.data = d.data[1:]
	return value, true
}

// Back returns the element at the end of the queue (without removing it)
//
// Returns:
//   - If the queue is not empty, returns the back element and true
//   - If the queue is empty, returns a zero value and false
//
// Time complexity: O(1)
func (d *ArrayDeque[T]) Back() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	return d.data[len(d.data)-1], true
}

// Front returns the element at the front of the queue (without removing it)
//
// Returns:
//   - If the queue is not empty, returns the front element and true
//   - If the queue is empty, returns a zero value and false
//
// Time complexity: O(1)
func (d *ArrayDeque[T]) Front() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	return d.data[0], true
}

// Len returns the number of elements in the queue
//
// Returns:
//   - The number of elements in the queue
//
// Time complexity: O(1)
func (d *ArrayDeque[T]) Len() int {
	return len(d.data)
}

// Capacity returns the capacity of the queue
//
// Returns:
//   - The capacity of the queue
//
// Time complexity: O(1)
func (d *ArrayDeque[T]) Capacity() int {
	return cap(d.data)
}

// IsEmpty checks if the queue is empty
//
// Returns:
//   - true if the queue is empty, false otherwise
//
// Time complexity: O(1)
func (d *ArrayDeque[T]) IsEmpty() bool {
	return len(d.data) == 0
}

// Clear empties the queue, removing all elements
//
// Time complexity: O(n)
// Note: This method sets all elements to zero value to avoid memory leaks
func (d *ArrayDeque[T]) Clear() {
	clear(d.data)
	d.data = d.data[:0:min(len(d.data), 8)]
}

// ShrinkToFit reduces the capacity of the queue to match its length
//
// Time complexity: O(n)
func (d *ArrayDeque[T]) ShrinkToFit() {
	if len(d.data) == 0 {
		d.data = make([]T, 0, 8) // Keep minimum capacity at 8 instead of setting to nil.
		return
	}
	newData := make([]T, len(d.data))
	copy(newData, d.data)
	d.data = newData
}

// Reserve ensures that the queue has enough capacity to hold the specified number of elements
//
// Parameters:
//   - capacity: The capacity to ensure
//
// Time complexity: O(n) (only when expansion is needed)
func (d *ArrayDeque[T]) Reserve(capacity int) {
	if capacity <= cap(d.data) {
		return
	}
	newData := make([]T, len(d.data), capacity)
	copy(newData, d.data)
	d.data = newData
}

// Clone creates a deep copy of the queue
//
// Returns:
//   - A new ArrayDeque instance containing all elements of the original queue
//
// Time complexity: O(n)
func (d *ArrayDeque[T]) Clone() *ArrayDeque[T] {
	newD := &ArrayDeque[T]{
		data: make([]T, len(d.data), cap(d.data)),
	}
	copy(newD.data, d.data)
	return newD
}
