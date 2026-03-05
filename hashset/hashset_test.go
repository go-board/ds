package hashset

import (
	"iter"
	"slices"
	"testing"

	"github.com/go-board/ds/hashutil"
)

func newIntHasher() hashutil.Default[int] {
	return hashutil.Default[int]{}
}

func TestHashSetBasicOperations(t *testing.T) {
	set := New[int](newIntHasher())

	if !set.IsEmpty() || set.Len() != 0 {
		t.Fatal("New set should be empty")
	}

	if !set.Insert(1) || set.Insert(1) {
		t.Fatal("Insert should differentiate between new and existing elements")
	}

	if !set.Contains(1) || set.Contains(2) {
		t.Fatal("Contains behavior unexpected")
	}

	if !set.Remove(1) || set.Remove(1) {
		t.Fatal("Remove behavior unexpected")
	}

	for i := 0; i < 5; i++ {
		set.Insert(i)
	}

	clone := set.Clone()
	set.Clear()
	if !clone.Contains(3) || !set.IsEmpty() {
		t.Fatal("Clone should get independent copy, Clear should empty original set")
	}

	clone.Compact() // Just ensuring it doesn't panic
}

func TestHashSetIterExtendAndEntry(t *testing.T) {
	set := New[int](newIntHasher())
	seq := func(yield func(int) bool) {
		for i := 0; i < 3; i++ {
			if !yield(i) {
				return
			}
		}
	}
	set.Extend(iter.Seq[int](seq))

	var collected []int
	for v := range set.Iter() {
		collected = append(collected, v)
	}
	slices.Sort(collected)
	if !slices.Equal(collected, []int{0, 1, 2}) {
		t.Fatalf("Iter result unexpected: %#v", collected)
	}

	entry := set.Entry(5)
	ptr := entry.OrInsert(struct{}{})
	if ptr == nil || !set.Contains(5) {
		t.Fatal("Entry OrInsert failed")
	}
	entry = set.Entry(5).AndModify(func(_ *struct{}) {})
	if _, ok := entry.Get(); !ok {
		t.Fatal("Entry Get should see just inserted element")
	}
}

func TestHashSetRelations(t *testing.T) {
	base := New[int](newIntHasher())
	for _, v := range []int{1, 2, 3} {
		base.Insert(v)
	}

	other := New[int](newIntHasher())
	for _, v := range []int{3, 4} {
		other.Insert(v)
	}

	union := base.Union(other)
	intersection := base.Intersection(other)
	diff := base.Difference(other)
	symmetric := base.SymmetricDifference(other)

	expect := func(hs *HashSet[int, hashutil.Default[int]], want []int) {
		var got []int
		for v := range hs.Iter() {
			got = append(got, v)
		}
		slices.Sort(got)
		if !slices.Equal(got, want) {
			t.Fatalf("Set contents expected %v, got %v", want, got)
		}
	}

	expect(union, []int{1, 2, 3, 4})
	expect(intersection, []int{3})
	expect(diff, []int{1, 2})
	expect(symmetric, []int{1, 2, 4})

	if !diff.IsSubset(union) || !union.IsSuperset(diff) {
		t.Fatal("Subset/superset judgment unexpected")
	}
	if union.IsDisjoint(other) || base.IsDisjoint(intersection) {
		t.Fatal("Disjoint judgment unexpected")
	}
}

// New test case: Test IsSubset function with all branches
func TestHashSetIsSubsetAllBranches(t *testing.T) {
	// Create test sets
	set1 := New[int](newIntHasher())
	set2 := New[int](newIntHasher())
	set3 := New[int](newIntHasher())

	// Insert elements
	for _, v := range []int{1, 2, 3, 4, 5} {
		set1.Insert(v)
	}

	for _, v := range []int{2, 3, 4} {
		set2.Insert(v)
	}

	for _, v := range []int{6, 7, 8} {
		set3.Insert(v)
	}

	// Test that set2 is a subset of set1
	if !set2.IsSubset(set1) {
		t.Error("set2 should be a subset of set1")
	}

	// Test that set1 is not a subset of set2
	if set1.IsSubset(set2) {
		t.Error("set1 should not be a subset of set2")
	}

	// Test case with disjoint sets
	if set3.IsSubset(set1) {
		t.Error("set3 should not be a subset of set1")
	}

	// Test that empty set is a subset of any set
	emptySet := New[int](newIntHasher())
	if !emptySet.IsSubset(set1) {
		t.Error("Empty set should be a subset of any set")
	}

	// Test that a set is a subset of itself
	if !set1.IsSubset(set1) {
		t.Error("Set should be a subset of itself")
	}
}

// New test case: Test IsDisjoint function with all branches
func TestHashSetIsDisjointAllBranches(t *testing.T) {
	// Create test sets
	set1 := New[int](newIntHasher())
	set2 := New[int](newIntHasher())
	set3 := New[int](newIntHasher())

	// Insert elements
	for _, v := range []int{1, 2, 3, 4, 5} {
		set1.Insert(v)
	}

	for _, v := range []int{6, 7, 8} {
		set2.Insert(v)
	}

	for _, v := range []int{4, 5, 6} {
		set3.Insert(v)
	}

	// Test disjoint sets
	if !set1.IsDisjoint(set2) {
		t.Error("set1 and set2 should be disjoint")
	}

	// Test intersecting sets
	if set1.IsDisjoint(set3) {
		t.Error("set1 and set3 should not be disjoint")
	}

	// Test that empty set is disjoint with any set
	emptySet := New[int](newIntHasher())
	if !emptySet.IsDisjoint(set1) {
		t.Error("Empty set should be disjoint with any set")
	}

	// Test two empty sets
	emptySet2 := New[int](newIntHasher())
	if !emptySet.IsDisjoint(emptySet2) {
		t.Error("Two empty sets should be disjoint")
	}
}

func TestHashSetEqual(t *testing.T) {
	a := New[int](newIntHasher())
	b := New[int](newIntHasher())
	for _, v := range []int{1, 2, 3} {
		a.Insert(v)
		b.Insert(v)
	}
	if !a.Equal(b) {
		t.Fatal("sets with same elements should be equal")
	}

	b.Insert(4)
	if a.Equal(b) {
		t.Fatal("sets with different lengths should not be equal")
	}

	c := New[int](newIntHasher())
	for _, v := range []int{1, 2, 4} {
		c.Insert(v)
	}
	if a.Equal(c) {
		t.Fatal("sets with same length but different elements should not be equal")
	}
}
