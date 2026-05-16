package arrayset

import (
	"cmp"
	"iter"

	"github.com/go-board/ds/arraymap"
)

var nothing struct{}

// ArraySet is an ordered set backed by ArrayMap.
type ArraySet[E any] struct {
	mapImpl *arraymap.ArrayMap[E, struct{}]
}

// New creates an empty ArraySet using comparator.
func New[E any](comparator func(E, E) int) *ArraySet[E] {
	return &ArraySet[E]{mapImpl: arraymap.New[E, struct{}](comparator)}
}

// NewOrdered creates an empty ArraySet for ordered element types.
func NewOrdered[E cmp.Ordered]() *ArraySet[E] {
	return &ArraySet[E]{mapImpl: arraymap.NewOrdered[E, struct{}]()}
}

// Insert adds value to the set.
func (s *ArraySet[E]) Insert(value E) bool {
	_, existed := s.mapImpl.Insert(value, nothing)
	return !existed
}

// Remove deletes value from the set.
func (s *ArraySet[E]) Remove(value E) bool {
	_, found := s.mapImpl.Remove(value)
	return found
}

// Contains reports whether value exists in the set.
func (s *ArraySet[E]) Contains(value E) bool {
	return s.mapImpl.ContainsKey(value)
}

// Len returns the number of elements.
func (s *ArraySet[E]) Len() int {
	return s.mapImpl.Len()
}

// IsEmpty reports whether the set has no elements.
func (s *ArraySet[E]) IsEmpty() bool {
	return s.mapImpl.IsEmpty()
}

// Clear removes all elements.
func (s *ArraySet[E]) Clear() {
	s.mapImpl.Clear()
}

// Clone creates a shallow copy of the set.
func (s *ArraySet[E]) Clone() *ArraySet[E] {
	return &ArraySet[E]{mapImpl: s.mapImpl.Clone()}
}

// First returns the first element in ascending order.
func (s *ArraySet[E]) First() (E, bool) {
	k, _, ok := s.mapImpl.First()
	return k, ok
}

// Last returns the last element in ascending order.
func (s *ArraySet[E]) Last() (E, bool) {
	k, _, ok := s.mapImpl.Last()
	return k, ok
}

// PopFirst removes and returns the first element in ascending order.
func (s *ArraySet[E]) PopFirst() (E, bool) {
	k, _, ok := s.mapImpl.PopFirst()
	return k, ok
}

// PopLast removes and returns the last element in ascending order.
func (s *ArraySet[E]) PopLast() (E, bool) {
	k, _, ok := s.mapImpl.PopLast()
	return k, ok
}

// GetComparator returns the comparator used by the set.
func (s *ArraySet[E]) GetComparator() func(E, E) int {
	return s.mapImpl.GetComparator()
}

// Union returns a new set containing elements from both sets.
func (s *ArraySet[E]) Union(other *ArraySet[E]) *ArraySet[E] {
	result := s.Clone()
	result.Extend(other.IterAsc())
	return result
}

// Intersection returns a new set containing elements present in both sets.
func (s *ArraySet[E]) Intersection(other *ArraySet[E]) *ArraySet[E] {
	result := New(s.GetComparator())
	nextA, stopA := iter.Pull(s.IterAsc())
	defer stopA()
	nextB, stopB := iter.Pull(other.IterAsc())
	defer stopB()

	a, okA := nextA()
	b, okB := nextB()
	cmpFn := s.GetComparator()

	for okA && okB {
		switch c := cmpFn(a, b); {
		case c == 0:
			result.Insert(a)
			a, okA = nextA()
			b, okB = nextB()
		case c < 0:
			a, okA = nextA()
		default:
			b, okB = nextB()
		}
	}

	return result
}

// Difference returns a new set containing elements in s but not in other.
func (s *ArraySet[E]) Difference(other *ArraySet[E]) *ArraySet[E] {
	result := New(s.GetComparator())
	nextA, stopA := iter.Pull(s.IterAsc())
	defer stopA()
	nextB, stopB := iter.Pull(other.IterAsc())
	defer stopB()

	a, okA := nextA()
	b, okB := nextB()
	cmpFn := s.GetComparator()

	for okA {
		if !okB {
			result.Insert(a)
			a, okA = nextA()
			continue
		}

		switch c := cmpFn(a, b); {
		case c == 0:
			a, okA = nextA()
			b, okB = nextB()
		case c < 0:
			result.Insert(a)
			a, okA = nextA()
		default:
			b, okB = nextB()
		}
	}

	return result
}

// SymmetricDifference returns a new set containing elements in either set but not both.
func (s *ArraySet[E]) SymmetricDifference(other *ArraySet[E]) *ArraySet[E] {
	result := New(s.GetComparator())
	nextA, stopA := iter.Pull(s.IterAsc())
	defer stopA()
	nextB, stopB := iter.Pull(other.IterAsc())
	defer stopB()

	a, okA := nextA()
	b, okB := nextB()
	cmpFn := s.GetComparator()

	for okA && okB {
		switch c := cmpFn(a, b); {
		case c == 0:
			a, okA = nextA()
			b, okB = nextB()
		case c < 0:
			result.Insert(a)
			a, okA = nextA()
		default:
			result.Insert(b)
			b, okB = nextB()
		}
	}
	for okA {
		result.Insert(a)
		a, okA = nextA()
	}
	for okB {
		result.Insert(b)
		b, okB = nextB()
	}

	return result
}

// IsSubset reports whether s is a subset of other.
func (s *ArraySet[E]) IsSubset(other *ArraySet[E]) bool {
	if s.Len() > other.Len() {
		return false
	}
	nextA, stopA := iter.Pull(s.IterAsc())
	defer stopA()
	nextB, stopB := iter.Pull(other.IterAsc())
	defer stopB()

	a, okA := nextA()
	b, okB := nextB()
	cmpFn := s.GetComparator()

	for okA {
		if !okB {
			return false
		}

		switch c := cmpFn(a, b); {
		case c == 0:
			a, okA = nextA()
			b, okB = nextB()
		case c < 0:
			return false
		default:
			b, okB = nextB()
		}
	}

	return true
}

// IsSuperset reports whether s is a superset of other.
func (s *ArraySet[E]) IsSuperset(other *ArraySet[E]) bool {
	return other.IsSubset(s)
}

// IsDisjoint reports whether s and other have no common elements.
func (s *ArraySet[E]) IsDisjoint(other *ArraySet[E]) bool {
	nextA, stopA := iter.Pull(s.IterAsc())
	defer stopA()
	nextB, stopB := iter.Pull(other.IterAsc())
	defer stopB()

	a, okA := nextA()
	b, okB := nextB()
	cmpFn := s.GetComparator()

	for okA && okB {
		switch c := cmpFn(a, b); {
		case c == 0:
			return false
		case c < 0:
			a, okA = nextA()
		default:
			b, okB = nextB()
		}
	}

	return true
}
