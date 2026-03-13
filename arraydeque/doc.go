// Package arraydeque implements a generic double-ended queue (deque) data structure.
//
// ArrayDeque is a linear collection that supports adding and removing elements from both ends
// in O(1) amortized time. It uses a circular buffer implementation for efficient memory usage.
//
// # Time Complexity
//
//   - PushBack, Back, Len, Cap, IsEmpty: O(1)
//   - PushFront, PopFront: O(n) worst case (shifts all elements)
//   - PopBack: O(1) amortized
//   - Get, Set: O(1)
//   - Traversal: O(n)
//
// # Features
//
//   - Double-ended operations: push/pop from both front and back
//   - Efficient circular buffer implementation
//   - Dynamic capacity management with growth strategy
//   - Memory efficient (no memory leaks)
//   - Generic type support
//   - Iterator support
//
// # Usage
//
// ArrayDeque is ideal for implementing queues and double-ended queues:
//
//	dq := arraydeque.New[int]()
//
//	// Add elements to the back
//	dq.PushBack(1)
//	dq.PushBack(2)
//	dq.PushBack(3)
//
//	// Add element to the front
//	dq.PushFront(0)
//
//	// Access elements without removing
//	front, _ := dq.Front() // 0
//	back, _ := dq.Back()   // 3
//
//	// Remove elements
//	val, _ := dq.PopFront() // 0, front element
//	val, _ = dq.PopBack()   // 3, back element
//
// # Use Cases
//
//   - Implementing queues (FIFO)
//   - Implementing stacks (LIFO) - but prefer ArrayStack
//   - Sliding window algorithms
//   - Breadth-first search implementations
//   - Ring buffers
//
// # Performance Notes
//
// PushFront and PopFront are O(n) because they require shifting all elements.
// If you only need front operations, consider using a linked list instead.
package arraydeque
