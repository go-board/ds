// Package btreeset implements a B-tree based ordered set data structure.
//
// BTreeSet is an ordered collection type where elements are arranged according to a defined order,
// and each element appears at most once in the set.
//
// Features:
//   - Ordered storage with support for sequential element traversal
//   - Efficient insertion, deletion, and lookup operations with O(log n) time complexity
//   - Support for set operations (union, intersection, difference, etc.)
//   - Support for range queries
//
// Example usage:
//
//	// Create a new ordered set for integers
//	set := btreeset.NewOrdered[int]()
//
//	// Insert elements
//	set.Insert(5)
//	set.Insert(3)
//	set.Insert(7)
//
//	// Check if an element exists
//	exists := set.Contains(5) // true
//
//	// Iterate through elements (in order)
//	for val := range set.IterAsc() {
//	    fmt.Println(val) // Output: 3, 5, 7
//	}
//
//	// Set operations
//	otherSet := btreeset.NewOrdered[int]()
//	otherSet.Insert(5)
//	otherSet.Insert(10)
//
//	union := set.Union(otherSet)        // {3, 5, 7, 10}
//	intersection := set.Intersection(otherSet) // {5}
package btreeset

import (
	"cmp"

	"github.com/go-board/ds/btree"
)

// BTreeSet is an ordered collection type that uses BTree as its underlying storage.
//
// T represents the type of elements in the set, supporting any type with a comparison function or ordered types.
// Elements in the set are arranged according to the order defined by the comparator, with each element being unique.
//
// Note: BTreeSet is not thread-safe. Additional synchronization mechanisms are required for concurrent access.

type BTreeSet[T any] struct {
	btree      *btree.BTree[T] // Underlying B-tree storage structure
	comparator func(T, T) int  // Element comparison function, used to determine order
}

// New creates a new BTreeSet instance suitable for elements of any type.
//
// Parameters:
//
//	comparator: Function used to compare elements. Returns negative when a < b, zero when a == b,
//	            and positive when a > b.
//
// Returns:
//
//	A pointer to the newly created BTreeSet
//
// Time complexity: O(1)
func New[T any](comparator func(T, T) int) *BTreeSet[T] {
	return &BTreeSet[T]{
		btree:      btree.New(comparator),
		comparator: comparator,
	}
}

// NewOrdered creates a new ordered BTreeSet suitable for types implementing the cmp.Ordered interface
// (such as int, string, etc.).
//
// Returns:
//
//	A pointer to the newly created BTreeSet
//
// Time complexity: O(1)
func NewOrdered[T cmp.Ordered]() *BTreeSet[T] {
	return New(cmp.Compare[T])
}

// Insert adds an element to the set.
//
// Parameters:
//
//	value: The element value to insert
//
// Returns:
//
//	true if the element was not present and was successfully inserted;
//	false if the element was already present
//
// Time complexity: O(log n)
func (s *BTreeSet[T]) Insert(value T) bool {
	_, found := s.btree.Search(value)
	if !found {
		s.btree.Insert(value)
	}
	return !found
}

// Remove removes the specified element from the set.
//
// Parameters:
//
//	value: The element value to remove
//
// Returns:
//
//	true if the element was present and successfully removed;
//	false if the element was not present
//
// Time complexity: O(log n)
func (s *BTreeSet[T]) Remove(value T) bool {
	return s.btree.Remove(value)
}

// Contains checks if the set contains the specified element.
//
// Parameters:
//
//	value: The element value to check
//
// Returns:
//
//	true if the set contains the element; false otherwise
//
// Time complexity: O(log n)
func (s *BTreeSet[T]) Contains(value T) bool {
	_, found := s.btree.Search(value)
	return found
}

// Len returns the number of elements in the set.
//
// Returns:
//
//	The count of elements in the set
//
// Time complexity: O(1)
func (s *BTreeSet[T]) Len() int {
	return s.btree.Len()
}

// IsEmpty checks if the set is empty.
//
// Returns:
//
//	true if the set is empty (contains no elements); false otherwise
//
// Time complexity: O(1)
func (s *BTreeSet[T]) IsEmpty() bool {
	return s.btree.Len() == 0
}

// Clear empties the set, removing all elements.
//
// After this operation, Len() will return 0.
//
// Time complexity: O(1)
func (s *BTreeSet[T]) Clear() {
	s.btree = btree.New(s.comparator)
}

// Clone creates a deep copy of the set.
//
// Returns:
//
//	A new BTreeSet instance containing all elements from the original set
//
// Note: This method performs a deep copy of the set structure, but elements are copied by value
// (if elements are pointers, the pointers are copied).
//
// Time complexity: O(n)
func (s *BTreeSet[T]) Clone() *BTreeSet[T] {
	newSet := New(s.comparator)
	for elem := range s.btree.IterAsc() {
		newSet.btree.Insert(elem)
	}
	return newSet
}

// Extend adds all elements from another iterable collection to the current set.
//
// Parameters:
//
//	iter: Iterator providing elements to add
//
// For each element in the iterator, it is added only if it does not already exist in the current set.
//
// Time complexity: O(m log n), where m is the number of elements in the iterator and n is the current size of the set

// Union creates a new set containing all elements from both the current set and another set (union).
//
// Parameters:
//
//	other: The other set to union with the current set
//
// Returns:
//
//	A new set containing all distinct elements from both sets
//
// Mathematical definition: A ∪ B = {x | x ∈ A or x ∈ B}
//
// Time complexity: O(n + m), where n and m are the sizes of the two sets
func (s *BTreeSet[T]) Union(other *BTreeSet[T]) *BTreeSet[T] {
	result := s.Clone()
	result.Extend(other.IterAsc())
	return result
}

// Intersection creates a new set containing elements that are present in both the current set and another set (intersection).
//
// Parameters:
//
//	other: The other set to intersect with the current set
//
// Returns:
//
//	A new set containing elements that exist in both sets
//
// Mathematical definition: A ∩ B = {x | x ∈ A and x ∈ B}
//
// Time complexity: O(min(n, m) * log(max(n, m))), where n and m are the sizes of the two sets
func (s *BTreeSet[T]) Intersection(other *BTreeSet[T]) *BTreeSet[T] {
	result := New(s.comparator)

	// Optimization: iterate over the smaller set for better efficiency
	if s.Len() > other.Len() {
		s, other = other, s
	}

	for e := range s.IterAsc() {
		if other.Contains(e) {
			result.Insert(e)
		}
	}

	return result
}

// Difference creates a new set containing elements present in the current set but not in another set (difference).
//
// Parameters:
//
//	other: The other set to compute the difference against
//
// Returns:
//
//	A new set containing elements that belong to the current set but not to the other set
//
// Mathematical definition: A \ B = {x | x ∈ A and x ∉ B}
//
// Time complexity: O(n * log m), where n is the size of the current set and m is the size of the other set
func (s *BTreeSet[T]) Difference(other *BTreeSet[T]) *BTreeSet[T] {
	result := New(s.comparator)

	for e := range s.IterAsc() {
		if !other.Contains(e) {
			result.Insert(e)
		}
	}

	return result
}

// SymmetricDifference creates a new set containing elements that are present in either the current set or another set but not both (symmetric difference).
//
// Parameters:
//
//	other: The other set to compute the symmetric difference against
//
// Returns:
//
//	A new set containing elements that belong to either the current set or the other set but not both
//
// Mathematical definition: A △ B = (A \ B) ∪ (B \ A) = {x | x ∈ A XOR x ∈ B}
//
// Time complexity: O(n * log m + m * log n), where n and m are the sizes of the two sets
func (s *BTreeSet[T]) SymmetricDifference(other *BTreeSet[T]) *BTreeSet[T] {
	result := New(s.comparator)

	// Add elements that exist only in the current set
	for k := range s.IterAsc() {
		if !other.Contains(k) {
			result.Insert(k)
		}
	}

	// Add elements that exist only in the other set
	for k := range other.IterAsc() {
		if !s.Contains(k) {
			result.Insert(k)
		}
	}

	return result
}

// IsSubset checks if the current set is a subset of another set.
//
// Parameters:
//
//	other: The parent set to check against
//
// Returns:
//
//	true if all elements of the current set are present in the other set; false otherwise
//
// Mathematical definition: A ⊆ B ⇨ for all x ∈ A, x ∈ B
//
// Time complexity: O(n * log m), where n is the size of the current set and m is the size of the other set
func (s *BTreeSet[T]) IsSubset(other *BTreeSet[T]) bool {
	if s.Len() > other.Len() {
		return false
	}

	for e := range s.IterAsc() {
		if !other.Contains(e) {
			return false
		}
	}

	return true
}

// IsSuperset checks if the current set is a superset of another set.
//
// Parameters:
//
//	other: The subset to check against
//
// Returns:
//
//	true if all elements of the other set are present in the current set; false otherwise
//
// Mathematical definition: A ⊇ B ⇨ for all x ∈ B, x ∈ A
//
// Time complexity: O(m * log n), where m is the size of the other set and n is the size of the current set
func (s *BTreeSet[T]) IsSuperset(other *BTreeSet[T]) bool {
	return other.IsSubset(s)
}

// IsDisjoint checks if the current set and another set are disjoint (have no elements in common).
//
// Parameters:
//
//	other: The other set to check against
//
// Returns:
//
//	true if the two sets have no elements in common; false otherwise
//
// Mathematical definition: A and B are disjoint ⇨ A ∩ B = ∅
//
// Time complexity: O(min(n, m) * log(max(n, m))), where n and m are the sizes of the two sets
func (s *BTreeSet[T]) IsDisjoint(other *BTreeSet[T]) bool {
	// Optimization: check the smaller set
	if s.Len() > other.Len() {
		s, other = other, s
	}

	for e := range s.IterAsc() {
		if other.Contains(e) {
			return false
		}
	}

	return true
}

// Equal checks whether two sets contain exactly the same elements.
func (s *BTreeSet[T]) Equal(other *BTreeSet[T]) bool {
	if s.Len() != other.Len() {
		return false
	}
	return s.IsSubset(other)
}

// First returns the first (smallest) element in the set.
//
// Returns:
//
//	The smallest element in the set and a boolean indicating if the set is non-empty
//	If the set is empty, returns the zero value of the element type and false
//
// Time complexity: O(log n)
func (s *BTreeSet[T]) First() (T, bool) {
	val, found := s.btree.First()
	return val, found
}

// Last returns the last (largest) element in the set.
//
// Returns:
//
//	The largest element in the set and a boolean indicating if the set is non-empty
//	If the set is empty, returns the zero value of the element type and false
//
// Time complexity: O(log n)
func (s *BTreeSet[T]) Last() (T, bool) {
	val, found := s.btree.Last()
	return val, found
}

// PopFirst removes and returns the first (smallest) element from the set.
//
// Returns:
//
//	The removed smallest element and a boolean indicating if the operation was successful
//	If the set is empty, returns the zero value of the element type and false
//
// Time complexity: O(log n)
func (s *BTreeSet[T]) PopFirst() (T, bool) {
	val, found := s.btree.PopFirst()
	return val, found
}

// PopLast removes and returns the last (largest) element from the set.
//
// Returns:
//
//	The removed largest element and a boolean indicating if the operation was successful
//	If the set is empty, returns the zero value of the element type and false
//
// Time complexity: O(log n)
func (s *BTreeSet[T]) PopLast() (T, bool) {
	val, found := s.btree.PopLast()
	return val, found
}
