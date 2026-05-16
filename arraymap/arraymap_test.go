package arraymap

import (
	"iter"
	"slices"
	"testing"

	"github.com/go-board/ds/bound"
)

func TestArrayMapCustomComparatorAndOrdering(t *testing.T) {
	m := New[string, int](func(a, b string) int {
		if len(a) != len(b) {
			return len(a) - len(b)
		}
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})

	m.Insert("bbb", 3)
	m.Insert("a", 1)
	m.Insert("cc", 2)
	m.Insert("bb", 22)

	keys := slices.Collect(m.KeysAsc())
	if !slices.Equal(keys, []string{"a", "bb", "cc", "bbb"}) {
		t.Fatalf("unexpected key order: %#v", keys)
	}

	if v, ok := m.Get("bb"); !ok || v != 22 {
		t.Fatalf("Get bb = %v,%v", v, ok)
	}
}

func TestArrayMapBasicOperations(t *testing.T) {
	m := NewOrdered[int, int]()
	if !m.IsEmpty() || m.Len() != 0 {
		t.Fatal("new map should be empty")
	}

	old, existed := m.Insert(3, 30)
	if existed || old != 0 {
		t.Fatal("insert on missing key should return zero,false")
	}
	m.Insert(1, 10)
	m.Insert(2, 20)

	old, existed = m.Insert(2, 200)
	if !existed || old != 20 {
		t.Fatal("insert on existing key should return old,true")
	}

	if v, ok := m.Get(2); !ok || v != 200 {
		t.Fatalf("Get returned %v,%v", v, ok)
	}

	if ptr, ok := m.GetMut(1); !ok || ptr == nil {
		t.Fatal("GetMut should return writable pointer")
	} else {
		*ptr = 11
	}

	if k, v, ok := m.GetKeyValue(1); !ok || k != 1 || v != 11 {
		t.Fatalf("GetKeyValue returned %v,%v,%v", k, v, ok)
	}

	if !m.ContainsKey(3) || m.ContainsKey(9) {
		t.Fatal("ContainsKey behavior is unexpected")
	}

	removed, ok := m.Remove(2)
	if !ok || removed != 200 {
		t.Fatalf("Remove returned %d,%v", removed, ok)
	}

	if _, ok := m.Remove(2); ok {
		t.Fatal("second Remove should fail")
	}
}

func TestArrayMapOrderedMethods(t *testing.T) {
	m := NewOrdered[int, int]()
	for _, k := range []int{5, 1, 4, 2, 3} {
		m.Insert(k, k*10)
	}

	if k, v, ok := m.First(); !ok || k != 1 || v != 10 {
		t.Fatalf("First returned %d,%d,%v", k, v, ok)
	}
	if k, v, ok := m.Last(); !ok || k != 5 || v != 50 {
		t.Fatalf("Last returned %d,%d,%v", k, v, ok)
	}

	ascKeys := slices.Collect(m.KeysAsc())
	if !slices.Equal(ascKeys, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("KeysAsc unexpected: %#v", ascKeys)
	}

	descKeys := slices.Collect(m.KeysDesc())
	if !slices.Equal(descKeys, []int{5, 4, 3, 2, 1}) {
		t.Fatalf("KeysDesc unexpected: %#v", descKeys)
	}

	r := bound.NewRangeBounds(bound.NewIncluded(2), bound.NewExcluded(5))
	rangeKeys := make([]int, 0)
	for k := range m.RangeAsc(r) {
		rangeKeys = append(rangeKeys, k)
	}
	if !slices.Equal(rangeKeys, []int{2, 3, 4}) {
		t.Fatalf("RangeAsc unexpected: %#v", rangeKeys)
	}

	for k, v := range m.RangeMutDesc(r) {
		*v = k * 100
	}
	for _, k := range []int{2, 3, 4} {
		if v, _ := m.Get(k); v != k*100 {
			t.Fatalf("RangeMutDesc did not update key=%d", k)
		}
	}

	if k, v, ok := m.PopFirst(); !ok || k != 1 || v != 10 {
		t.Fatalf("PopFirst returned %d,%d,%v", k, v, ok)
	}
	if k, v, ok := m.PopLast(); !ok || k != 5 || v != 50 {
		t.Fatalf("PopLast returned %d,%d,%v", k, v, ok)
	}
}

func TestArrayMapEntryAPI(t *testing.T) {
	m := NewOrdered[string, int]()

	ptr := m.Entry("x").OrInsert(1)
	if *ptr != 1 {
		t.Fatalf("OrInsert returned %d", *ptr)
	}

	m.Entry("x").AndModify(func(v *int) { *v += 9 })
	if v, _ := m.Get("x"); v != 10 {
		t.Fatalf("AndModify failed, got %d", v)
	}

	calls := 0
	m.Entry("y").OrInsertWith(func() int {
		calls++
		return 2
	})
	m.Entry("y").OrInsertWith(func() int {
		calls++
		return 3
	})
	if calls != 1 {
		t.Fatalf("OrInsertWith callback calls = %d", calls)
	}

	m.Entry("z").OrInsertWithKey(func(k string) int { return len(k) })
	if v, _ := m.Entry("z").Get(); v != 1 {
		t.Fatalf("OrInsertWithKey returned %d", v)
	}

	old, existed := m.Entry("x").Insert(100)
	if !existed || old != 10 {
		t.Fatalf("Entry.Insert(existing) = %d,%v", old, existed)
	}
	old, existed = m.Entry("n").Insert(200)
	if existed || old != 0 {
		t.Fatalf("Entry.Insert(new) = %d,%v", old, existed)
	}

	if !m.Entry("n").Delete() || m.Entry("n").Delete() {
		t.Fatal("Delete return value is unexpected")
	}
}

func TestArrayMapCompatMethods(t *testing.T) {
	m := NewOrdered[int, int]()
	m.Extend(iter.Seq2[int, int](func(yield func(int, int) bool) {
		for i := 3; i >= 1; i-- {
			if !yield(i, i*10) {
				return
			}
		}
	}))

	if !slices.Equal(slices.Collect(m.Keys()), []int{1, 2, 3}) {
		t.Fatal("Keys should match ascending order")
	}
	for _, ptr := range slices.Collect(m.ValuesMut()) {
		*ptr += 1
	}
	for k, v := range m.IterMut() {
		*v += k
	}
	for k, v := range m.Iter() {
		if v != k*11+1 {
			t.Fatalf("unexpected Iter value for %d => %d", k, v)
		}
	}
}

func TestArrayMapCloneClearAndConstructors(t *testing.T) {
	src := map[int]int{2: 20, 1: 10}
	m := NewFromMap(src)
	if !slices.Equal(slices.Collect(m.KeysAsc()), []int{1, 2}) {
		t.Fatal("NewFromMap should build ordered keys")
	}

	clone := m.Clone()
	clone.Insert(3, 30)
	if m.ContainsKey(3) {
		t.Fatal("clone should be independent")
	}

	m.Clear()
	if !m.IsEmpty() || m.Len() != 0 {
		t.Fatal("Clear should empty map")
	}
}
