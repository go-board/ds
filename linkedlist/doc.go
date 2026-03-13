// Package linkedlist implements a generic doubly linked list data structure.
//
// LinkedList is a linear collection where each element contains a value and pointers
// to the next and previous elements. This enables efficient insertion and deletion
// at both ends and anywhere in the list.
//
// # Time Complexity
//
//   - PushBack, PushFront: O(1)
//   - PopBack, PopFront: O(1)
//   - InsertAfter, InsertBefore: O(1) given iterator/position
//   - Get, Set: O(n) - requires traversal
//   - Remove: O(1) given iterator
//   - Traversal: O(n)
//
// # Features
//
//   - O(1) insertion/deletion at both ends
//   - Bidirectional traversal (forward and backward)
//   - Bidirectional iterators for flexible navigation
//   - Batch operations: PushBackBatch, PushFrontBatch
//   - Filtering: Retain, RemoveIf
//   - Generic type support
//
// # Usage
//
// LinkedList excels at frequent insertions and deletions:
//
//	list := linkedlist.New[int]()
//
//	// Add elements to both ends
//	list.PushBack(1)
//	list.PushBack(2)
//	list.PushFront(0)
//
//	// Access without removal
//	front, _ := list.Front() // 0
//	back, _ := list.Back()   // 2
//
//	// Remove from both ends
//	val, _ := list.PopFront() // 0
//	val, _ = list.PopBack()   // 2
//
//	// Iterate forward
//	for iter := list.Iter(); iter.Next(); {
//	    fmt.Println(iter.Value())
//	}
//
//	// Iterate backward
//	for iter := list.IterRev(); iter.Prev(); {
//	    fmt.Println(iter.Value())
//	}
//
// # Use Cases
//
//   - When frequent insertions/deletions in the middle are needed
//   - Implementing queues (with O(1) enqueue/dequeue)
//   - When you need bidirectional traversal
//   - Implementing other data structures (e.g., stacks, deques)
//
// # Performance Notes
//
// LinkedList provides O(1) insertion/deletion but O(n) random access.
// If you need fast random access, consider ArrayDeque instead.
package linkedlist
