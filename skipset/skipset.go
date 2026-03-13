package skipset

import (
	"cmp"

	"github.com/go-board/ds/skipmap"
)

// empty struct used as the value type in SkipMap, occupies no memory
var nothing struct{}

// SkipSet implements a sorted set based on skip lists.
// A skip list is a data structure that allows fast lookup, functioning as a multi-layered linked list.
// Each layer is sorted, and upper layers are subsets of lower layers, used to accelerate the search process.
// SkipSet implements set functionality by internally maintaining a SkipMap where keys are set elements
// and values are empty structs.

type SkipSet[E any] struct {
	mapImpl *skipmap.SkipMap[E, struct{}] // SkipMap used internally
}

// New creates a new empty SkipSet using the specified key comparison function.
// Parameters:
//   - comparator: Function for comparing keys, must not be nil. Should return a negative number when a < b, 0 when a == b, and a positive number when a > b.
//
// Returns:
//   - Pointer to the newly created SkipSet
func New[E any](comparator func(E, E) int) *SkipSet[E] {
	return &SkipSet[E]{
		mapImpl: skipmap.New[E, struct{}](comparator),
	}
}

// NewOrdered creates a new empty SkipSet for element types that support ordered comparison.
// This is a convenience function that uses [cmp.Compare] as the comparator.
// Type Parameters:
//   - E: Element type, must implement the [cmp.Ordered] interface
//
// Returns:
//   - Pointer to the newly created SkipSet
func NewOrdered[E cmp.Ordered]() *SkipSet[E] {
	return &SkipSet[E]{
		mapImpl: skipmap.NewOrdered[E, struct{}](),
	}
}

// Insert adds an element to the set.
// Parameters:
//   - key: The element to add
//
// Returns:
//   - true if the element was newly added (didn't exist before)
//   - false if the element already existed
func (ss *SkipSet[E]) Insert(key E) bool {
	_, updated := ss.mapImpl.Insert(key, nothing)
	return !updated // If updated is false, it means it's newly inserted
}

// Remove removes the specified element from the set.
// Parameters:
//   - key: The element to remove
//
// Returns:
//   - true if the element existed and was removed
//   - false if the element didn't exist
func (ss *SkipSet[E]) Remove(key E) bool {
	_, found := ss.mapImpl.Remove(key)
	return found
}

// Contains checks if the set contains the specified element.
// Parameters:
//   - key: The element to check
//
// Returns:
//   - true if the element exists, false otherwise
func (ss *SkipSet[E]) Contains(key E) bool {
	return ss.mapImpl.ContainsKey(key)
}

// Len returns the number of elements in the set.
// Returns:
//   - The number of elements in the set
func (ss *SkipSet[E]) Len() int {
	return ss.mapImpl.Len()
}

// IsEmpty checks if the set is empty (contains no elements).
// Returns:
//   - true if the set is empty, false otherwise
func (ss *SkipSet[E]) IsEmpty() bool {
	return ss.mapImpl.IsEmpty()
}

// Clear removes all elements from the set, making it empty.
func (ss *SkipSet[E]) Clear() {
	ss.mapImpl.Clear()
}

// Clone creates a deep copy of the set.
// Returns:
//   - A new SkipSet containing all the same elements
func (ss *SkipSet[E]) Clone() *SkipSet[E] {
	clone := &SkipSet[E]{
		mapImpl: ss.mapImpl.Clone(),
	}
	return clone
}

// First returns the first (smallest) element in the set.
// Returns:
//   - The element and true if the set is not empty
//   - Zero value and false if the set is empty
func (ss *SkipSet[E]) First() (E, bool) {
	key, _, found := ss.mapImpl.First()
	return key, found
}

// Last returns the last (largest) element in the set.
// Returns:
//   - The element and true if the set is not empty
//   - Zero value and false if the set is empty
func (ss *SkipSet[E]) Last() (E, bool) {
	key, _, found := ss.mapImpl.Last()
	return key, found
}

// PopFirst removes and returns the first (smallest) element in the set.
// Returns:
//   - The removed element and true if the set is not empty
//   - Zero value and false if the set is empty
func (ss *SkipSet[E]) PopFirst() (E, bool) {
	key, _, found := ss.mapImpl.PopFirst()
	return key, found
}

// PopLast removes and returns the last (largest) element in the set.
// Returns:
//   - The removed element and true if the set is not empty
//   - Zero value and false if the set is empty
func (ss *SkipSet[E]) PopLast() (E, bool) {
	key, _, found := ss.mapImpl.PopLast()
	return key, found
}

// Union computes the union of the current set and another set, returning a new set containing all unique elements.
// Parameters:
//   - other: Another [SkipSet]
//
// Returns:
//   - New [SkipSet] containing all elements from the current set and the other set
func (ss *SkipSet[E]) Union(other *SkipSet[E]) *SkipSet[E] {
	// Create a new set
	result := ss.Clone()

	// Add all elements from the other set
	for key := range other.IterAsc() {
		result.Insert(key)
	}

	return result
}

// Intersection computes the intersection of the current set and another set, returning a new set containing only elements present in both sets.
// Parameters:
//   - other: Another [SkipSet]
//
// Returns:
//   - New [SkipSet] containing elements that exist in both the current set and the other set
func (ss *SkipSet[E]) Intersection(other *SkipSet[E]) *SkipSet[E] {
	// Create a new set
	result := New(ss.mapImpl.GetComparator())

	// Iterate over the smaller set for optimization
	var small, large *SkipSet[E]
	if ss.Len() <= other.Len() {
		small, large = ss, other
	} else {
		small, large = other, ss
	}

	// Check if each element in the smaller set exists in the larger set
	for key := range small.IterAsc() {
		if large.Contains(key) {
			result.Insert(key)
		}
	}

	return result
}

// Difference computes the difference between the current set and another set, returning a new set containing elements in the current set but not in the other set.
// Parameters:
//   - other: Another [SkipSet]
//
// Returns:
//   - New [SkipSet] containing elements that exist in the current set but not in the other set
func (ss *SkipSet[E]) Difference(other *SkipSet[E]) *SkipSet[E] {
	// Create a new set
	result := New(ss.mapImpl.GetComparator())

	// Check if each element in the current set is not in the other set
	for key := range ss.IterAsc() {
		if !other.Contains(key) {
			result.Insert(key)
		}
	}

	return result
}

// SymmetricDifference computes the symmetric difference between the current set and another set, returning a new set containing elements that are in either set but not in both.
// Parameters:
//   - other: Another [SkipSet]
//
// Returns:
//   - New [SkipSet] containing elements that are in either set but not in both sets simultaneously
func (ss *SkipSet[E]) SymmetricDifference(other *SkipSet[E]) *SkipSet[E] {
	// Create two difference sets
	diff1 := ss.Difference(other)
	diff2 := other.Difference(ss)

	// Return the union of both differences
	for key := range diff2.IterAsc() {
		diff1.Insert(key)
	}

	return diff1
}

// IsSubset checks if the current set is a subset of another set.
// Parameters:
//   - other: Another [SkipSet]
//
// Returns:
//   - true if all elements of the current set are in the other set
//   - false otherwise
func (ss *SkipSet[E]) IsSubset(other *SkipSet[E]) bool {
	// If the current set has more elements than the other, it can't be a subset
	if ss.Len() > other.Len() {
		return false
	}

	// Check if each element in the current set is in the other set
	for key := range ss.IterAsc() {
		if !other.Contains(key) {
			return false
		}
	}

	return true
}

// IsSuperset checks if the current set is a superset of another set.
// Parameters:
//   - other: Another [SkipSet]
//
// Returns:
//   - true if all elements of the other set are in the current set
//   - false otherwise
func (ss *SkipSet[E]) IsSuperset(other *SkipSet[E]) bool {
	// Call the other set's IsSubset method
	return other.IsSubset(ss)
}

// IsDisjoint checks if the current set and another set are disjoint (have no elements in common).
// Parameters:
//   - other: Another [SkipSet]
//
// Returns:
//   - true if the two sets have no common elements
//   - false otherwise
func (ss *SkipSet[E]) IsDisjoint(other *SkipSet[E]) bool {
	// Iterate over the smaller set for optimization
	var small, large *SkipSet[E]
	if ss.Len() <= other.Len() {
		small, large = ss, other
	} else {
		small, large = other, ss
	}

	// Check if each element in the smaller set exists in the larger set
	for key := range small.IterAsc() {
		if large.Contains(key) {
			return false
		}
	}

	return true
}

// Equal checks if the current set is equal to another set (contains the same elements).
// Parameters:
//   - other: Another [SkipSet]
//
// Returns:
//   - true if the two sets contain the same elements
//   - false otherwise
func (ss *SkipSet[E]) Equal(other *SkipSet[E]) bool {
	// If the element counts differ, sets can't be equal
	if ss.Len() != other.Len() {
		return false
	}

	// Check if each element in the current set is in the other set
	return ss.IsSubset(other)
}
