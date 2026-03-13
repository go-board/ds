// Package priorityqueue implements a generic priority queue data structure.
//
// PriorityQueue is a heap-based priority queue that supports both min-heap and max-heap modes.
// Elements are always retrieved in priority order (lowest for min-heap, highest for max-heap).
// All operations maintain the heap property in O(log n) time.
//
// # Time Complexity
//
//   - Push: O(log n)
//   - Pop: O(log n)
//   - Peek: O(1)
//   - Len, IsEmpty: O(1)
//   - Traversal: O(n) - not in priority order
//
// # Features
//
//   - Min-heap mode: retrieves smallest element first
//   - Max-heap mode: retrieves largest element first
//   - Efficient push and pop operations
//   - Peek at highest priority without removal
//   - Batch operations: PushBatch, Clear
//   - Generic type support with custom comparators
//
// # Usage
//
// Priority queues are used for scheduling, Dijkstra's algorithm, and more:
//
//	// Min-heap (smallest element first) - recommended for cmp.Ordered types
//	minHeap := priorityqueue.NewMinOrdered[int]()
//	minHeap.Push(5)
//	minHeap.Push(3)
//	minHeap.Push(7)
//
//	// Pop in priority order
//	for !minHeap.IsEmpty() {
//	    val, _ := minHeap.Pop()
//	    fmt.Println(val) // Output: 3, 5, 7
//	}
//
//	// Max-heap (largest element first)
//	maxHeap := priorityqueue.NewMaxOrdered[int]()
//	maxHeap.Push(5)
//	maxHeap.Push(3)
//	maxHeap.Push(7)
//
//	for !maxHeap.IsEmpty() {
//	    val, _ := maxHeap.Pop()
//	    fmt.Println(val) // Output: 7, 5, 3
//	}
//
//	// Custom comparator for complex types
//	type Task struct {
//	    Priority int
//	    Name     string
//	}
//	pq := priorityqueue.NewMin[Task](func(a, b Task) int {
//	    return a.Priority - b.Priority
//	})
//
// # Use Cases
//
//   - Task scheduling (highest priority first)
//   - Dijkstra's shortest path algorithm
//   - Huffman coding
//   - Event-driven simulations
//   - K-largest/k-smallest element problems
//
// # Implementation Notes
//
// PriorityQueue uses a binary heap implemented on a slice.
// The heap property is maintained after each push/pop operation.
package priorityqueue
