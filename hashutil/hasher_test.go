package hashutil

import (
	"hash/maphash"
	"testing"
)

// helperHasher records the input to Hash calls for easy assertion
type helperHasher struct {
	hashed []int
}

func (h *helperHasher) Equal(x, y int) bool {
	return x == y
}

func (h *helperHasher) Hash(hash *maphash.Hash, v int) {
	h.hashed = append(h.hashed, v)
	var buf [8]byte
	for i := 0; i < 8; i++ {
		buf[i] = byte(v >> (8 * i))
	}
	hash.Write(buf[:])
}

func TestDefaultHasherHashConsistency(t *testing.T) {
	seed := maphash.MakeSeed()
	var h1, h2, h3 maphash.Hash
	h1.SetSeed(seed)
	h2.SetSeed(seed)
	h3.SetSeed(seed)

	hasher := Default[int]{}
	hasher.Hash(&h1, 42)
	hasher.Hash(&h2, 42)
	hasher.Hash(&h3, 7)

	if h1.Sum64() != h2.Sum64() {
		t.Fatal("Hash values for the same input should be consistent")
	}
	if h1.Sum64() == h3.Sum64() {
		t.Fatal("Hash values for different inputs should be different")
	}
	if !hasher.Equal(5, 5) || hasher.Equal(5, 4) {
		t.Fatal("Default.Equal behavior is abnormal")
	}
}

func TestSliceHasher(t *testing.T) {
	record := &helperHasher{}
	hasher := NewSliceHasher[[]int](record)

	if !hasher.Equal([]int{1, 2, 3}, []int{1, 2, 3}) {
		t.Fatal("Slices with the same content should be considered equal")
	}
	if hasher.Equal([]int{1, 2}, []int{2, 1}) {
		t.Fatal("Different order should be considered unequal")
	}

	var h maphash.Hash
	hasher.Hash(&h, []int{4, 5, 6})

	if len(record.hashed) != 3 {
		t.Fatalf("Hash should be called for each element, expected 3 times, got %d", len(record.hashed))
	}
	expected := map[int]int{4: 1, 5: 1, 6: 1}
	actual := map[int]int{}
	for _, v := range record.hashed {
		actual[v]++
	}
	if len(actual) != len(expected) {
		t.Fatalf("Recorded element set is incorrect: %#v", actual)
	}
	for k, c := range expected {
		if actual[k] != c {
			t.Fatalf("Element %d should be counted %d times, got %d", k, c, actual[k])
		}
	}
}

func TestMapHasher(t *testing.T) {
	record := &helperHasher{}
	hasher := NewMapHasher[map[string]int](record)

	a := map[string]int{"a": 1, "b": 2}
	b := map[string]int{"b": 2, "a": 1}
	c := map[string]int{"a": 1, "b": 3}

	if !hasher.Equal(a, b) {
		t.Fatal("Maps with the same content should be equal")
	}
	if hasher.Equal(a, c) {
		t.Fatal("Maps with different content should not be equal")
	}

	var h maphash.Hash
	hasher.Hash(&h, a)

	actual := map[int]int{}
	for _, v := range record.hashed {
		actual[v]++
	}
	if len(actual) != len(a) {
		t.Fatalf("Hash should iterate through each value in the map, got %#v", actual)
	}
	for _, expected := range a {
		if actual[expected] != 1 {
			t.Fatalf("Value %d should be hashed once, got %d times", expected, actual[expected])
		}
	}
}
