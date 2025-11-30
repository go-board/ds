# go-board/ds

<div align="center">

A comprehensive, high-performance generic data structures library for Go 1.24+

</div>

## Overview

**go-board/ds** is a feature-rich data structures library for Go that leverages the power of Go 1.24's generics to provide type-safe, efficient implementations of common and advanced data structures. This library is designed to be both performant and developer-friendly, with a consistent API design across all implementations.

## Features

- **Full Generic Support**: Leverage Go 1.24 generics for type-safe data structures
- **Comprehensive Data Structure Collection**: Linked lists, deques, maps, sets, priority queues, and more
- **Performance-Optimized Implementations**: Each data structure is carefully optimized for its specific use cases
- **Consistent API Design**: Familiar interfaces across all data structures
- **Flexible Hashing Strategies**: Customizable hashers for any type
- **Complete Iterator Support**: Multiple iteration patterns for easy data processing

## Installation

```bash
go get github.com/go-board/ds
```

## Quick Start

Here's a simple example demonstrating some of the library's capabilities:

```go
package main

import (
	"fmt"

	"github.com/go-board/ds"
)

func main() {
	// Create a new hash map
	m := ds.NewHashMap[string, int](ds.DefaultHasher[string]{})
	m.Insert("apple", 5)
	m.Insert("banana", 3)
	m.Insert("cherry", 7)

	// Retrieve a value
	if val, found := m.Get("banana"); found {
		fmt.Printf("Banana count: %d\n", val)
	}

	// Create a linked list
	list := ds.NewLinkedList[int]()
	list.PushBack(10)
	list.PushBack(20)
	list.PushFront(5)

	// Iterate through the list
	fmt.Println("Linked list elements:")
	for val := range list.Iter() {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
}
```

## Data Structure Capabilities

The following table shows the core capabilities supported by each data structure:
- ✅ Indicates the data structure supports this function
- ❌ Indicates the data structure does not support this function

| Function | LinkedList | ArrayDeque | HashMap | HashSet | BTree | BTreeMap | BTreeSet | SkipMap | SkipSet | TrieMap | PriorityQueue |
|---------|------------|------------|---------|---------|-------|----------|----------|---------|---------|---------|---------------|
| **Basic Operations** | | | | | | | | | | |
| Clear   | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Len     | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| IsEmpty | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Compact | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ |
| **Access & Modification** | | | | | | | | | | |
| Get     | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ✅ | ❌ | ✅ | ❌ |
| GetMut  | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ✅ | ❌ | ✅ | ❌ |
| Insert  | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| Remove  | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| **List/Queue Operations** | | | | | | | | | | |
| PushFront | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| PopFront  | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| PushBack  | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| PopBack   | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Front     | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Back      | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Set Operations** | | | | | | | | | | |
| Contains  | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ✅ | ❌ | ❌ |
| ContainsKey | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ✅ | ❌ | ❌ |
| **Ordered Data Operations** | | | | | | | | | | |
| First    | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| Last     | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| Range    | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| PopFirst | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| PopLast  | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| **Iterators** | | | | | | | | | | |
| Iter     | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| IterMut  | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ❌ | ✅ | ✅ | ❌ |
| **Collection Operations** | | | | | | | | | | |
| Extend   | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| Difference | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| Union    | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| Intersection | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| IsSubset | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| IsSuperset | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| SymmetricDifference | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| Equal    | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| **Priority Queue Operations** | | | | | | | | | | |
| Push    | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Pop     | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Peek    | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |

## Core Data Structures

### Linked List

```go
// Create a new linked list
list := ds.NewLinkedList[int]()

// Add elements
list.PushBack(1)
list.PushFront(0)
list.PushBack(2)

// Remove elements
value, _ := list.PopFront()  // 0
value, _ = list.PopBack()    // 2

// Iterate
for val := range list.Iter() {
    fmt.Println(val)  // 1
}
```

### Hash Map

```go
// Create a new hash map with default hasher
m := ds.NewHashMap[string, int](ds.DefaultHasher[string]{})

// Insert key-value pairs
m.Insert("one", 1)
m.Insert("two", 2)

// Retrieve values
value, found := m.Get("one")

// Update values
oldValue, updated := m.Insert("one", 10)

// Iterate through all key-value pairs
for key, val := range m.Iter() {
    fmt.Printf("%s: %d\n", key, val)
}
```

### Priority Queue

```go
// Create a min-heap (sorts by age ascending)
minHeap := ds.NewMinPriorityQueue[Person](func(a, b Person) int {
    return a.Age - b.Age
})

// Add elements
minHeap.Push(Person{Name: "Alice", Age: 30})
minHeap.Push(Person{Name: "Bob", Age: 25})

// Get highest priority element
person, _ := minHeap.Pop()  // Bob

// Peek at the highest priority element without removing
person, _ = minHeap.Peek()  // Alice
```

### BTreeMap

```go
// Create an ordered BTreeMap for ordered types
m := ds.NewOrderedBTreeMap[string, int]()

// Insert key-value pairs
m.Insert("apple", 5)
m.Insert("banana", 3)
m.Insert("cherry", 7)

// Get first and last elements
first, _ := m.First()      // apple: 5
last, _ := m.Last()        // cherry: 7

// Range query
m.Range(func(key string, value int) bool {
    fmt.Printf("%s: %d\n", key, value)

### TrieMap

```go
// Create a new TrieMap
m := ds.NewTrieMap[string, int]()

// Insert key-value pairs
m.Insert("apple", 5)
m.Insert("app", 3)
m.Insert("banana", 7)

// Get values
value, found := m.Get("apple")  // 5, true

// Check if key exists
value, found = m.Get("ap")      // 0, false (key not found)

// Remove a key
m.Remove("app")

// Iterate through all key-value pairs
for key, val := range m.Iter() {
    fmt.Printf("%s: %d\n", key, val)
}
```
```

## Custom Hashing

For complex types or custom hashing behavior, implement your own `Hasher`:

```go
// Custom struct
type Point struct {
    X, Y int
}

// Custom hasher
type PointHasher struct{}

func (PointHasher) Equal(p1, p2 Point) bool {
    return p1.X == p2.X && p1.Y == p2.Y
}

func (PointHasher) Hash(h *maphash.Hash, p Point) {
    maphash.WriteInt(h, p.X)
    maphash.WriteInt(h, p.Y)
}

// Usage
pointMap := ds.NewHashMap[Point, string](PointHasher{})
pointMap.Insert(Point{X: 1, Y: 2}, "Start point")
```

## Requirements

- Go 1.24 or higher (for full generics support)
- No external dependencies

## Testing

Run the test suite:

```bash
go test ./... -v
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Ensure your code follows Go's standard formatting
2. Add appropriate test cases for any new functionality
3. Update documentation as needed
4. Maintain consistency with existing API design

## License

This project is licensed under the MIT License. See the LICENSE file for details.