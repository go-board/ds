package hashset

import (
	"github.com/go-board/ds/hashmap"
	"github.com/go-board/ds/hashutil"
)

var nothing struct{} // Sentinel value for the HashMap-backed set implementation.

// HashSet is an unordered collection type that uses HashMap as its underlying storage.
//
// E is the type of elements in the set, supporting any type.
// H is the Hasher implementation used for hash computation and equality comparison.
// Elements in the set serve as keys in the HashMap, with empty struct values for space efficiency.
//
// Notes:
//  - HashSet does not guarantee element order, and iteration order may change with insertions and deletions
//  - HashSet is not thread-safe, requiring additional synchronization for multi-goroutine access
//  - For optimal performance, provide an efficient hash function implementation

type HashSet[E any, H hashutil.Hasher[E]] struct {
	table  *hashmap.HashMap[E, struct{}, H] // Backing hash table storage.
	hasher H                                // Hasher used for hashing and equality checks.
}

// New creates a new HashSet instance.
//
// Parameters:
//
//	hasher: Concrete implementation instance for element hash computation and equality comparison
//
// Return value:
//
//	Pointer to the newly created HashSet
//
// Time complexity: O(1)
func New[E any, H hashutil.Hasher[E]](hasher H) *HashSet[E, H] {
	return &HashSet[E, H]{
		table:  hashmap.New[E, struct{}](hasher),
		hasher: hasher,
	}
}

// NewComparable creates a new HashSet instance with a default hasher for comparable key types.
//
// Type Parameters:
//   - E: Element type, must be comparable
//
// Returns:
//   - Pointer to the newly created HashSet
//
// Example:
//
//	set := hashset.NewComparable[string]()
func NewComparable[E comparable]() *HashSet[E, hashutil.Default[E]] {
	return New(hashutil.Default[E]{})
}

// Insert adds an element to the set.
//
// Parameters:
//
//	value: The element value to insert
//
// Return value:
//
//	true if the element does not exist and was successfully inserted; false if the element already exists
//
// Time complexity: Average O(1), worst O(n)
func (hs *HashSet[E, H]) Insert(value E) bool {
	_, had := hs.table.Insert(value, nothing)
	return !had
}

// Remove removes the specified element from the set.
//
// Parameters:
//
//	value: The element value to remove
//
// Return value:
//
//	true if the element exists and was successfully removed; false if the element does not exist
//
// Time complexity: Average O(1), worst O(n)
func (hs *HashSet[E, H]) Remove(value E) bool {
	_, found := hs.table.Remove(value)
	return found
}

// Contains checks if the set contains the specified element.
//
// Parameters:
//
//	value: The element value to check
//
// Return value:
//
//	true if the set contains the element; false otherwise
//
// Time complexity: Average O(1), worst O(n)
func (hs *HashSet[E, H]) Contains(value E) bool {
	return hs.table.ContainsKey(value)
}

// Len returns the number of elements in the set.
//
// Return value:
//
//	The count of elements in the set
//
// Time complexity: O(1)
func (hs *HashSet[E, H]) Len() int {
	return hs.table.Len()
}

// IsEmpty checks if the set is empty.
//
// Return value:
//
//	true if the set is empty (contains no elements); false otherwise
//
// Time complexity: O(1)
func (hs *HashSet[E, H]) IsEmpty() bool {
	return hs.table.IsEmpty()
}

// Clear empties the set, removing all elements.
//
// After this operation, Len() will return 0.
//
// Time complexity: O(n), where n is the number of elements in the set
func (hs *HashSet[E, H]) Clear() {
	hs.table.Clear()
}

// Clone creates a deep copy of the set.
//
// Return value:
//
//	New HashSet instance containing all elements of the original set
//
// Note: This method performs a deep copy of the set structure, but elements themselves are copied by value (if elements are pointers, the pointers are copied).
//
// Time complexity: O(n), where n is the number of elements in the set
func (hs *HashSet[E, H]) Clone() *HashSet[E, H] {
	return &HashSet[E, H]{
		table:  hs.table.Clone(),
		hasher: hs.hasher,
	}
}

// Compact compacts the set, removing all deleted nodes.
//
// This operation frees memory and improves iteration and lookup efficiency, especially after numerous deletion operations.
//
// Time complexity: O(n), where n is the number of buckets in the set
func (hs *HashSet[E, H]) Compact() {
	hs.table.Compact()
}

// Entry retrieves the Entry state for an element, for flexible handling of insertion operations.
//
// Parameters:
//
//	value: The element value to retrieve the Entry state for
//
// Return value:
//
//	The Entry state for the corresponding element, which can be used for more complex insertion logic
//
// This method is similar to Rust's HashSet::entry method, allowing conditional insertion in a single lookup.
//
// Time complexity: Average O(1), worst O(n)
func (hs *HashSet[E, H]) Entry(value E) hashmap.Entry[E, struct{}, H] {
	return hs.table.Entry(value)
}

// Union creates a new set containing all elements from the current set and another set (union).
//
// Parameters:
//
//	other: Another set to union with the current set
//
// Return value:
//
//	New set containing all distinct elements from both sets
//
// Mathematical definition: A ∪ B = {x | x ∈ A or x ∈ B}
//
// Time complexity: Average O(n + m), where n and m are the sizes of the two sets
func (hs *HashSet[E, H]) Union(other *HashSet[E, H]) *HashSet[E, H] {
	result := hs.Clone()
	result.Extend(other.Iter())
	return result
}

// Intersection creates a new set containing elements that exist in both the current set and another set (intersection).
//
// Parameters:
//
//	other: Another set to intersect with the current set
//
// Return value:
//
//	New set containing elements that exist in both sets
//
// Mathematical definition: A ∩ B = {x | x ∈ A and x ∈ B}
//
// Time complexity: Average O(min(n, m)), where n and m are the sizes of the two sets
func (hs *HashSet[E, H]) Intersection(other *HashSet[E, H]) *HashSet[E, H] {
	result := hs.Clone()
	result.Clear()

	// Iterate over the smaller set to reduce lookup cost.
	if hs.Len() > other.Len() {
		hs, other = other, hs
	}

	for e := range hs.Iter() {
		if other.Contains(e) {
			result.Insert(e)
		}
	}

	return result
}

// Difference creates a new set containing elements that exist in the current set but not in another set (difference).
//
// Parameters:
//
//	other: Another set to compute difference with the current set
//
// Return value:
//
//	New set containing elements belonging to the current set but not to the other set
//
// Mathematical definition: A \ B = {x | x ∈ A and x ∉ B}
//
// Time complexity: Average O(n), where n is the size of the current set
func (hs *HashSet[E, H]) Difference(other *HashSet[E, H]) *HashSet[E, H] {
	result := hs.Clone()
	result.Clear()

	for e := range hs.Iter() {
		if !other.Contains(e) {
			result.Insert(e)
		}
	}

	return result
}

// SymmetricDifference creates a new set containing elements that exist only in the current set or only in another set (symmetric difference).
//
// Parameters:
//
//	other: Another set to compute symmetric difference with the current set
//
// Return value:
//
//	New set containing elements that belong only to the current set or only to the other set
//
// Mathematical definition: A △ B = (A \ B) ∪ (B \ A) = {x | x ∈ A XOR x ∈ B}
//
// Time complexity: Average O(n + m), where n and m are the sizes of the two sets
func (hs *HashSet[E, H]) SymmetricDifference(other *HashSet[E, H]) *HashSet[E, H] {
	result := hs.Clone()
	result.Clear()

	// Add elements that exist only in the current set.
	for k := range hs.Iter() {
		if !other.Contains(k) {
			result.Insert(k)
		}
	}

	// Add elements that exist only in the other set.
	for k := range other.Iter() {
		if !hs.Contains(k) {
			result.Insert(k)
		}
	}

	return result
}

// IsSubset checks if the current set is a subset of another set.
//
// Parameters:
//
//	other: Parent set to check against
//
// Return value:
//
//	true if all elements of the current set exist in the other set; false otherwise
//
// Mathematical definition: A ⊆ B ⇨ for all x ∈ A, x ∈ B
//
// Time complexity: Average O(n), where n is the size of the current set
func (hs *HashSet[E, H]) IsSubset(other *HashSet[E, H]) bool {
	if hs.Len() > other.Len() {
		return false
	}

	for e := range hs.Iter() {
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
//	other: Subset to check against
//
// Return value:
//
//	true if all elements of the other set exist in the current set; false otherwise
//
// Mathematical definition: A ⊇ B ⇨ for all x ∈ B, x ∈ A
//
// Time complexity: Average O(m), where m is the size of the other set
func (hs *HashSet[E, H]) IsSuperset(other *HashSet[E, H]) bool {
	return other.IsSubset(hs)
}

// IsDisjoint checks if the current set and another set are disjoint (have no common elements).
//
// Parameters:
//
//	other: Another set to check against
//
// Return value:
//
//	true if the two sets have no common elements; false otherwise
//
// Mathematical definition: A and B are disjoint ⇨ A ∩ B = ∅
//
// Time complexity: Average O(min(n, m)), where n and m are the sizes of the two sets
func (hs *HashSet[E, H]) IsDisjoint(other *HashSet[E, H]) bool {
	// Check the smaller set first for better performance.
	if hs.Len() > other.Len() {
		hs, other = other, hs
	}

	for e := range hs.Iter() {
		if other.Contains(e) {
			return false
		}
	}

	return true
}

// Equal checks if two sets contain exactly the same elements.
func (hs *HashSet[E, H]) Equal(other *HashSet[E, H]) bool {
	if hs.Len() != other.Len() {
		return false
	}
	return hs.IsSubset(other)
}
