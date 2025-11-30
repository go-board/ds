package hashmap

import (
	"encoding/binary"
	"hash/maphash"
	"iter"
	"slices"
	"testing"

	"github.com/go-board/ds/hashutil"
)

type collisionHasher struct{}

func (collisionHasher) Equal(a, b int) bool { return a == b }

func (collisionHasher) Hash(h *maphash.Hash, _ int) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], 1)
	h.Write(buf[:])
}

func TestHashMapBasicOperations(t *testing.T) {
	hm := New[int, int](hashutil.Default[int]{})

	if !hm.IsEmpty() || hm.Len() != 0 {
		t.Fatal("New HashMap should be empty")
	}

	old, updated := hm.Insert(1, 10)
	if updated || old != 0 {
		t.Fatal("First insertion should not return an old value")
	}
	old, updated = hm.Insert(1, 20)
	if !updated || old != 10 {
		t.Fatal("Subsequent insertion should return the old value")
	}

	if v, ok := hm.Get(1); !ok || v != 20 {
		t.Fatalf("Get should return 20, got %v %v", v, ok)
	}

	if key, val, ok := hm.GetKeyValue(1); !ok || key != 1 || val != 20 {
		t.Fatal("GetKeyValue returned unexpectedly")
	}

	if ptr, ok := hm.GetMut(1); !ok || ptr == nil {
		t.Fatal("GetMut should return a writable pointer")
	} else {
		*ptr = 30
	}

	if !hm.ContainsKey(1) || hm.ContainsKey(2) {
		t.Fatal("ContainsKey behavior is unexpected")
	}

	if val, ok := hm.Remove(1); !ok || val != 30 {
		t.Fatal("Remove should delete key 1 and return the latest value")
	}
	if _, ok := hm.Remove(1); ok {
		t.Fatal("Duplicate removal should return false")
	}
}

func TestHashMapCompactExtendAndIterators(t *testing.T) {
	hm := New[int, int](collisionHasher{})

	hm.Insert(1, 10)
	hm.Insert(2, 20)
	hm.Remove(1)
	hm.Insert(1, 30) // Reuse soft deleted node

	hm.Compact()

	seqData := []struct {
		k int
		v int
	}{{3, 30}, {4, 40}}
	hm.Extend(iter.Seq2[int, int](func(yield func(int, int) bool) {
		for _, item := range seqData {
			if !yield(item.k, item.v) {
				return
			}
		}
	}))

	keys := make([]int, 0)
	for key := range hm.Keys() {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	if !slices.Equal(keys, []int{1, 2, 3, 4}) {
		t.Fatalf("Keys result is unexpected: %#v", keys)
	}

	values := make([]int, 0)
	for value := range hm.Values() {
		values = append(values, value)
	}
	slices.Sort(values)
	if len(values) != 4 {
		t.Fatalf("Values count is unexpected: %#v", values)
	}

	pairs := make(map[int]int)
	for key, value := range hm.Iter() {
		pairs[key] = value
	}
	if len(pairs) != 4 {
		t.Fatalf("Iter should iterate all key-value pairs: %#v", pairs)
	}

	for key, ptr := range hm.IterMut() {
		*ptr = key * 100
	}
	for key, value := range hm.Iter() {
		if value != key*100 {
			t.Fatalf("IterMut modification did not take effect: key=%d val=%d", key, value)
		}
	}
}

func TestHashMapEntryAPI(t *testing.T) {
	hm := New[int, int](hashutil.Default[int]{})

	entry := hm.Entry(1)
	ptr := entry.OrInsert(10)
	if ptr == nil || hm.Len() != 1 {
		t.Fatal("OrInsert should insert a new key")
	}

	entry = hm.Entry(1)
	entry.AndModify(func(v *int) { *v += 5 })
	if value, _ := hm.Get(1); value != 15 {
		t.Fatalf("AndModify should modify the value, got %d", value)
	}

	hm.Entry(2).OrInsertWith(func() int { return 20 })
	hm.Entry(3).OrInsertWithKey(func(k int) int { return k * 10 })

	if value, _ := hm.Entry(2).Get(); value != 20 {
		t.Fatal("Get should return the newly inserted value")
	}
	if value, _ := hm.Entry(3).Get(); value != 30 {
		t.Fatal("OrInsertWithKey result is unexpected")
	}

	ptr = hm.Entry(2).Insert(200)
	if ptr == nil || *ptr != 200 {
		t.Fatal("Insert should update and return a writable pointer")
	}

	entry = hm.Entry(99)
	if value, ok := entry.Get(); ok || value != 0 {
		t.Fatal("Get for non-existent key should return zero value and false")
	}
}

func TestHashMapEdgeCases(t *testing.T) {
	hm := New[int, int](collisionHasher{})
	hm.Insert(1, 10)
	hm.Insert(2, 20)

	if val, ok := hm.Get(3); ok || val != 0 {
		t.Fatal("Get should return zero value when bucket exists but key is missing")
	}

	hm.Remove(1)
	hm.Remove(2)
	if val, ok := hm.Get(1); ok || val != 0 {
		t.Fatal("Deleted key should not return value again")
	}
	if k, v, ok := hm.GetKeyValue(1); ok || k != 1 || v != 0 {
		t.Fatalf("GetKeyValue for deleted key should return original key and zero value, got %v %v %v", k, v, ok)
	}
	if ptr, ok := hm.GetMut(1); ok || ptr != nil {
		t.Fatal("GetMut for deleted key should return nil")
	}

	// Trigger Compact
	hm.Insert(3, 30)
	if hm.Len() != 1 {
		t.Fatalf("Element count should be 1 after Compact, got %d", hm.Len())
	}

	// Extend triggers Compact branch
	hm.Remove(3)
	hm.Extend(iter.Seq2[int, int](func(yield func(int, int) bool) {
		for i := 4; i <= 5; i++ {
			if !yield(i, i*10) {
				return
			}
		}
	}))

	// Covering both new and old branches of Entry.Insert
	ptr := hm.Entry(4).Insert(400)
	if ptr == nil || *ptr != 400 {
		t.Fatal("Entry.Insert should update existing key")
	}
	ptr = hm.Entry(6).Insert(600)
	if ptr == nil || *ptr != 600 {
		t.Fatal("Entry.Insert should insert new key")
	}

	if val := hm.Entry(4).OrInsert(444); *val != 400 {
		t.Fatal("OrInsert should return original value for existing key")
	}
	callCount := 0
	hm.Entry(7).OrInsertWith(func() int {
		callCount++
		return 700
	})
	if callCount != 1 {
		t.Fatalf("OrInsertWith should only be called when key is missing, count=%d", callCount)
	}
	hm.Entry(7).OrInsertWith(func() int {
		t.Fatal("Callback should not be called for existing key")
		return 0
	})

	hm.Entry(8).OrInsertWithKey(func(k int) int { return k * 100 })
	hm.Entry(8).OrInsertWithKey(func(k int) int {
		t.Fatalf("OrInsertWithKey callback should not be called for existing key: %d", k)
		return 0
	})

	// Early exit paths for Keys/Values/Iter/IterMut
	hm.Keys()(func(int) bool { return false })
	hm.Values()(func(int) bool { return false })
	hm.Iter()(func(int, int) bool { return false })
	hm.IterMut()(func(int, *int) bool { return false })
}
