package priorityqueue

import "cmp"

// PriorityQueue is a heap-based priority queue
type PriorityQueue[T any] struct {
	data     []T              // slice that stores elements
	lessFunc func(a, b T) int // comparison function
}

// NewMin creates a new min-heap priority queue
//
// Parameters:
//   - cmp: comparison function that returns negative when a < b, zero when a == b, and positive when a > b
//
// Return value:
//   - Newly created min-heap priority queue instance
//
// Min-heap property: smallest element has highest priority
//
// Time complexity: O(1)
func NewMin[T any](cmp func(T, T) int) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data:     make([]T, 0),
		lessFunc: cmp, // min-heap uses the original comparison function directly
	}
}

// NewMinOrdered creates a new min-heap priority queue with ordered elements
//
// Type Parameters:
//   - T: Element type, must be ordered
//
// Return value:
//   - Newly created min-heap priority queue instance with ordered elements
//
// Min-heap property: smallest element has highest priority
//
// Time complexity: O(1)
func NewMinOrdered[T cmp.Ordered]() *PriorityQueue[T] {
	return NewMin(cmp.Compare[T])
}

// NewMax creates a new max-heap priority queue
//
// Parameters:
//   - cmp: comparison function that returns negative when a < b, zero when a == b, and positive when a > b
//
// Return value:
//   - Newly created max-heap priority queue instance
//
// Max-heap property: largest element has highest priority
//
// Time complexity: O(1)
func NewMax[T any](cmp func(T, T) int) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data: make([]T, 0),
		lessFunc: func(a, b T) int {
			return -cmp(a, b) // max-heap reverses the comparison result
		},
	}
}

// NewMaxOrdered creates a new max-heap priority queue with ordered elements
//
// Type Parameters:
//   - T: Element type, must be ordered
//
// Return value:
//   - Newly created max-heap priority queue instance with ordered elements
//
// Max-heap property: largest element has highest priority
//
// Time complexity: O(1)
func NewMaxOrdered[T cmp.Ordered]() *PriorityQueue[T] {
	return NewMax(cmp.Compare[T])
}

// Len returns the number of elements in the priority queue
//
// Returns:
//   - The number of elements in the queue
//
// Time complexity: O(1)
func (pq *PriorityQueue[T]) Len() int {
	return len(pq.data)
}

// IsEmpty checks if the priority queue is empty
//
// Returns:
//   - true if the queue is empty, false otherwise
//
// Time complexity: O(1)
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return len(pq.data) == 0
}

// Push adds an element to the priority queue
//
// Parameters:
//   - value: The value to add to the queue
//
// Time complexity: O(log n)
func (pq *PriorityQueue[T]) Push(value T) {
	// Add new element to the end of the heap
	pq.data = append(pq.data, value)
	// Perform swim operation to maintain heap property
	pq.swim(len(pq.data) - 1)
}

// Pop removes and returns the element with the highest priority (heap top)
//
// Returns:
//   - If the queue is not empty, returns the element with the highest priority and true
//   - If the queue is empty, returns the zero value of type T and false
//
// Time complexity: O(log n)
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	var zero T
	if pq.IsEmpty() {
		return zero, false
	}

	// Save the top element
	top := pq.data[0]
	// Move last element to the top
	lastIndex := len(pq.data) - 1
	pq.data[0] = pq.data[lastIndex]
	// Delete the last element
	pq.data = pq.data[:lastIndex]

	// If queue is not empty, perform sink operation
	if !pq.IsEmpty() {
		pq.sink(0)
	}

	return top, true
}

// Peek returns the element with the highest priority (heap top) but does not remove it
//
// Returns:
//   - If the queue is not empty, returns the element with the highest priority and true
//   - If the queue is empty, returns the zero value of type T and false
//
// Time complexity: O(1)
func (pq *PriorityQueue[T]) Peek() (T, bool) {
	var zero T
	if pq.IsEmpty() {
		return zero, false
	}
	return pq.data[0], true
}

// Clear clears the priority queue
//
// Time complexity: O(1)
func (pq *PriorityQueue[T]) Clear() {
	pq.data = make([]T, 0)
}

// Internal method: swim operation, maintains min-heap property
func (pq *PriorityQueue[T]) swim(index int) {
	parent := (index - 1) / 2
	// When current node is smaller than parent, swap them and continue swimming
	if index > 0 && pq.less(index, parent) {
		pq.swap(index, parent)
		pq.swim(parent)
	}
}

// Internal method: sink operation, maintains min-heap property
func (pq *PriorityQueue[T]) sink(index int) {
	smallest := index
	left := 2*index + 1
	right := 2*index + 2
	size := len(pq.data)

	// If left child exists and is smaller than current node
	if left < size && pq.less(left, smallest) {
		smallest = left
	}

	// If right child exists and is smaller than current smallest node
	if right < size && pq.less(right, smallest) {
		smallest = right
	}

	// If smallest node is not current node, swap and continue sinking
	if smallest != index {
		pq.swap(index, smallest)
		pq.sink(smallest)
	}
}

// Internal method: compare elements at two indices
func (pq *PriorityQueue[T]) less(i, j int) bool {
	// lessFunc returns negative value if element at i is smaller than element at j
	return pq.lessFunc(pq.data[i], pq.data[j]) < 0
}

// Internal method: swap elements at two indices
func (pq *PriorityQueue[T]) swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
}
