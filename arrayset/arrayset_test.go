package arrayset

import (
	"slices"
	"testing"

	"github.com/go-board/ds/bound"
)

func TestArraySetBasicOperations(t *testing.T) {
	s := NewOrdered[int]()
	if !s.IsEmpty() || s.Len() != 0 {
		t.Fatal("new set should be empty")
	}

	if !s.Insert(2) || !s.Insert(1) || s.Insert(2) {
		t.Fatal("Insert behavior is unexpected")
	}

	if !s.Contains(1) || s.Contains(3) {
		t.Fatal("Contains behavior is unexpected")
	}

	if !s.Remove(1) || s.Remove(1) {
		t.Fatal("Remove behavior is unexpected")
	}
}

func TestArraySetOrderingAndRange(t *testing.T) {
	s := New[int](func(a, b int) int { return b - a }) // descending comparator
	for _, v := range []int{1, 3, 2, 4} {
		s.Insert(v)
	}

	if !slices.Equal(slices.Collect(s.IterAsc()), []int{4, 3, 2, 1}) {
		t.Fatalf("IterAsc order mismatch: %#v", slices.Collect(s.IterAsc()))
	}
	if !slices.Equal(slices.Collect(s.IterDesc()), []int{1, 2, 3, 4}) {
		t.Fatalf("IterDesc order mismatch: %#v", slices.Collect(s.IterDesc()))
	}

	bounds := bound.NewRangeBounds(bound.NewIncluded(3), bound.NewIncluded(2))
	if !slices.Equal(slices.Collect(s.RangeAsc(bounds)), []int{3, 2}) {
		t.Fatalf("RangeAsc mismatch: %#v", slices.Collect(s.RangeAsc(bounds)))
	}
}

func TestArraySetOrderOpsAndClone(t *testing.T) {
	s := NewOrdered[int]()
	s.Extend(slices.Values([]int{3, 1, 2}))

	if v, ok := s.First(); !ok || v != 1 {
		t.Fatalf("First = %d,%v", v, ok)
	}
	if v, ok := s.Last(); !ok || v != 3 {
		t.Fatalf("Last = %d,%v", v, ok)
	}
	if v, ok := s.PopFirst(); !ok || v != 1 {
		t.Fatalf("PopFirst = %d,%v", v, ok)
	}
	if v, ok := s.PopLast(); !ok || v != 3 {
		t.Fatalf("PopLast = %d,%v", v, ok)
	}

	clone := s.Clone()
	clone.Insert(9)
	if s.Contains(9) {
		t.Fatal("clone should be independent")
	}

	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("Clear should empty set")
	}
}

func TestArraySetSetAlgebra(t *testing.T) {
	a := NewOrdered[int]()
	b := NewOrdered[int]()
	a.Extend(slices.Values([]int{1, 2, 3}))
	b.Extend(slices.Values([]int{3, 4, 5}))

	if !slices.Equal(slices.Collect(a.Union(b).IterAsc()), []int{1, 2, 3, 4, 5}) {
		t.Fatal("Union mismatch")
	}
	if !slices.Equal(slices.Collect(a.Intersection(b).IterAsc()), []int{3}) {
		t.Fatal("Intersection mismatch")
	}
	if !slices.Equal(slices.Collect(a.Difference(b).IterAsc()), []int{1, 2}) {
		t.Fatal("Difference mismatch")
	}
	if !slices.Equal(slices.Collect(a.SymmetricDifference(b).IterAsc()), []int{1, 2, 4, 5}) {
		t.Fatal("SymmetricDifference mismatch")
	}

	if !a.Intersection(b).IsSubset(a) {
		t.Fatal("IsSubset mismatch")
	}
	if !a.Union(b).IsSuperset(a) {
		t.Fatal("IsSuperset mismatch")
	}
	if !a.IsDisjoint(NewOrdered[int]()) {
		t.Fatal("IsDisjoint with empty set should be true")
	}
}
