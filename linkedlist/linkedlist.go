// Package linkedlist implements a generic doubly linked list data structure.
//
// LinkedList is a doubly linked list that supports efficient addition and removal of elements at both ends,
// as well as bidirectional traversal. It is suitable for scenarios with frequent insertion and deletion operations.
//
// Features:
//   - Efficient addition and removal of elements at both ends
//   - Support for bidirectional traversal
//   - Support for batch operations
//   - Support for retaining elements that meet certain conditions
//
// Example usage:
//
//	// Create an integer linked list
//	list := linkedlist.New[int]()
//
//	// Add elements at the end
//	list.PushBack(1)
//	list.PushBack(2)
//
//	// Add element at the beginning
//	list.PushFront(0)
//
//	// Access elements
//	front, _ := list.Front() // 0
//	back, _ := list.Back()   // 2
//
//	// Traverse elements
//	for val := range list.Iter() {
//	    fmt.Println(val) // Output: 0, 1, 2
//	}
//
//	// Remove elements
//	val, _ := list.PopFront() // 0
//	val, _ := list.PopBack()  // 2
package linkedlist

// node represents a node structure in the doubly linked list
type node[T any] struct {
	value T        // value stored in the node
	prev  *node[T] // pointer to the previous node
	next  *node[T] // pointer to the next node
}

// LinkedList is a generic list implemented based on a doubly linked list
type LinkedList[T any] struct {
	head *node[T] // pointer to the first actual node
	tail *node[T] // pointer to the last actual node
	size int      // number of elements in the list
}

// New creates a new empty linked list
//
// Type Parameters:
//   - T: The type of elements in the list
//
// Returns:
//   - A new LinkedList instance
//
// Time complexity: O(1)
func New[T any]() *LinkedList[T] {
	// Empty list has head and tail as nil
	return &LinkedList[T]{}
}

// Append appends all elements from another linked list to the end of the current list
// After appending, the other list will be empty
//
// Parameters:
//   - other: The linked list to append
//
// Time complexity: O(1)
func (l *LinkedList[T]) Append(other *LinkedList[T]) {
	// If other list is empty, do nothing
	if other.size == 0 {
		return
	}

	if l.size == 0 {
		// If current list is empty, directly take over other list's nodes
		l.head = other.head
		l.tail = other.tail
		l.size = other.size
	} else {
		// Connect current list's tail to other list's head
		l.tail.next = other.head
		other.head.prev = l.tail

		// Update current list's tail and size
		l.tail = other.tail
		l.size += other.size
	}

	// Clear other list
	other.Clear()
}

// Back returns the last element in the linked list
//
// Returns:
//   - If the list is not empty, returns the last element and true
//   - If the list is empty, returns zero value and false
//
// Time complexity: O(1)
func (l *LinkedList[T]) Back() (T, bool) {
	if l.size == 0 {
		var zero T
		return zero, false
	}
	return l.tail.value, true
}

// BackMut returns a mutable reference to the last element in the linked list and a success flag
//
// Returns:
//   - If the list is not empty, returns a pointer to the last element and true
//   - If the list is empty, returns nil and false
//
// Time complexity: O(1)
func (l *LinkedList[T]) BackMut() (*T, bool) {
	if l.size == 0 {
		return nil, false
	}
	return &l.tail.value, true
}

// Clear removes all elements from the linked list, resetting it to its initial empty state
//
// Time complexity: O(1)
func (l *LinkedList[T]) Clear() {
	// Reset head and tail pointers to nil
	l.head = nil
	l.tail = nil

	// Reset size counter
	l.size = 0
}

// Front returns the first element in the linked list
//
// Returns:
//   - If the list is not empty, returns the first element and true
//   - If the list is empty, returns zero value and false
//
// Time complexity: O(1)
func (l *LinkedList[T]) Front() (T, bool) {
	if l.size == 0 {
		var zero T
		return zero, false
	}
	return l.head.value, true
}

// FrontMut returns a mutable reference to the first element in the linked list and a success flag
//
// Returns:
//   - If the list is not empty, returns a pointer to the first element and true
//   - If the list is empty, returns nil and false
//
// Time complexity: O(1)
func (l *LinkedList[T]) FrontMut() (*T, bool) {
	if l.size == 0 {
		return nil, false
	}
	return &l.head.value, true
}

// iterator-related methods have been moved to `iter.go`

// IsEmpty checks if the linked list is empty
//
// Returns:
//   - true if the list is empty, false otherwise
//
// Time complexity: O(1)
func (l *LinkedList[T]) IsEmpty() bool {
	return l.size == 0
}

// Len returns the number of elements in the linked list
//
// Returns:
//   - The number of elements in the list
//
// Time complexity: O(1)
func (l *LinkedList[T]) Len() int {
	return l.size
}

// PopBack removes and returns the last element in the linked list
//
// Returns:
//   - If the list is not empty, returns the last element and true
//   - If the list is empty, returns zero value and false
//
// Time complexity: O(1)
func (l *LinkedList[T]) PopBack() (T, bool) {
	if l.size == 0 {
		var zero T
		return zero, false
	}

	value := l.tail.value

	if l.size == 1 {
		l.head = nil
		l.tail = nil
	} else {
		l.tail = l.tail.prev
		l.tail.next = nil
	}

	// decrement size count
	l.size--

	return value, true
}

// PopFront removes and returns the first element in the linked list
//
// Returns:
//   - If the list is not empty, returns the first element and true
//   - If the list is empty, returns zero value and false
//
// Time complexity: O(1)
func (l *LinkedList[T]) PopFront() (T, bool) {
	if l.size == 0 {
		var zero T
		return zero, false
	}

	value := l.head.value

	if l.size == 1 {
		l.head = nil
		l.tail = nil
	} else {
		l.head = l.head.next
		l.head.prev = nil
	}

	l.size--

	return value, true
}

// PushBack appends a new element to the end of the linked list
//
// Parameters:
//   - value: The value to add to the list
//
// Time complexity: O(1)
func (l *LinkedList[T]) PushBack(value T) {
	newNode := &node[T]{value: value}

	if l.size == 0 {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		newNode.prev = l.tail
		l.tail = newNode
	}

	l.size++
}

// PushFront appends a new element to the beginning of the linked list
//
// Parameters:
//   - value: The value to add to the list
//
// Time complexity: O(1)
func (l *LinkedList[T]) PushFront(value T) {
	newNode := &node[T]{value: value}

	if l.size == 0 {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode.next = l.head
		l.head.prev = newNode
		l.head = newNode
	}

	l.size++
}

// Retain retains only the elements that satisfy the given predicate function
//
// Parameters:
//   - f: A function that takes an element of type T and returns a boolean value
//
// Time complexity: O(n)
func (l *LinkedList[T]) Retain(f func(T) bool) {
	current := l.head
	for current != nil {
		next := current.next
		if !f(current.value) {
			if current.prev == nil {
				l.head = next
			} else {
				current.prev.next = next
			}

			if current.next == nil {
				l.tail = current.prev
			} else {
				current.next.prev = current.prev
			}

			l.size--
		}
		current = next
	}
}
