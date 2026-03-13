// Package arraystack implements a generic stack (LIFO) data structure.
//
// ArrayStack provides last-in-first-out (LIFO) semantics with O(1) time complexity
// for push and pop operations. It's built on top of ArrayDeque for efficient memory usage.
//
// # Time Complexity
//
//   - Push: O(1) amortized
//   - Pop: O(1)
//   - Peek: O(1)
//   - Len, IsEmpty: O(1)
//   - Traversal: O(n)
//
// # Features
//
//   - LIFO (last-in-first-out) semantics
//   - O(1) push and pop operations
//   - Efficient memory usage with dynamic resizing
//   - Generic type support
//   - Iterator support for traversal
//
// # Usage
//
// ArrayStack is simple and efficient:
//
//	stack := arraystack.New[int]()
//
//	// Push elements (add to top)
//	stack.Push(1)
//	stack.Push(2)
//	stack.Push(3)
//
//	// Peek at top element without removing
//	top, _ := stack.Peek() // 3
//
//	// Pop elements (remove from top)
//	val, _ := stack.Pop() // 3
//	val, _ = stack.Pop() // 2
//
//	// Check if empty
//	if stack.IsEmpty() {
//	    fmt.Println("Stack is empty")
//	}
//
// # Use Cases
//
//   - Expression evaluation (postfix notation)
//   - Backtracking algorithms
//   - Function call stack simulation
//   - Undo mechanisms
//   - Depth-first search
//
// # Performance Notes
//
// ArrayStack provides O(1) operations for all basic stack operations.
// It's the preferred choice over ArrayDeque when you only need stack semantics.
package arraystack
