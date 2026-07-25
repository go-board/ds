package arraymap

import (
	"iter"

	"github.com/go-board/ds/bound"
	diter "github.com/go-board/ds/internal/iter"
	"github.com/go-board/ds/internal/kv"
)

func (m *ArrayMap[K, V]) iterPairAsc() iter.Seq[*kv.Pair[K, V]] {
	return func(yield func(*kv.Pair[K, V]) bool) {
		for i := range m.items {
			if !yield(&m.items[i]) {
				return
			}
		}
	}
}

func (m *ArrayMap[K, V]) iterPairDesc() iter.Seq[*kv.Pair[K, V]] {
	return func(yield func(*kv.Pair[K, V]) bool) {
		for i := len(m.items) - 1; i >= 0; i-- {
			if !yield(&m.items[i]) {
				return
			}
		}
	}
}

func (m *ArrayMap[K, V]) rangePairAsc(bounds bound.RangeBounds[K]) iter.Seq[*kv.Pair[K, V]] {
	start, end := m.rangeIndices(bounds)
	return func(yield func(*kv.Pair[K, V]) bool) {
		for i := start; i < end; i++ {
			if !yield(&m.items[i]) {
				return
			}
		}
	}
}

func (m *ArrayMap[K, V]) rangePairDesc(bounds bound.RangeBounds[K]) iter.Seq[*kv.Pair[K, V]] {
	start, end := m.rangeIndices(bounds)
	return func(yield func(*kv.Pair[K, V]) bool) {
		for i := end - 1; i >= start; i-- {
			if !yield(&m.items[i]) {
				return
			}
		}
	}
}

func (m *ArrayMap[K, V]) rangeIndices(bounds bound.RangeBounds[K]) (start, end int) {
	start = 0
	end = len(m.items)

	if v, ok := bounds.Start.Value(); ok {
		idx, found := m.search(v)
		if bounds.Start.IsIncluded() {
			start = idx
		} else {
			start = idx
			if found {
				start++
			}
		}
	}

	if v, ok := bounds.End.Value(); ok {
		idx, found := m.search(v)
		if bounds.End.IsIncluded() {
			end = idx
			if found {
				end++
			}
		} else {
			end = idx
		}
	}

	if start < 0 {
		start = 0
	}
	if end > len(m.items) {
		end = len(m.items)
	}
	if start > end {
		start = end
	}

	return start, end
}

// IterAsc returns an iterator over all key-value pairs in ascending key order.
func (m *ArrayMap[K, V]) IterAsc() iter.Seq2[K, V] {
	return diter.Split(m.iterPairAsc(), (*kv.Pair[K, V]).KV)
}

// IterMutAsc returns a mutable iterator over all key-value pairs in ascending key order.
func (m *ArrayMap[K, V]) IterMutAsc() iter.Seq2[K, *V] {
	return diter.Split(m.iterPairAsc(), (*kv.Pair[K, V]).KVMut)
}

// IterDesc returns an iterator over all key-value pairs in descending key order.
func (m *ArrayMap[K, V]) IterDesc() iter.Seq2[K, V] {
	return diter.Split(m.iterPairDesc(), (*kv.Pair[K, V]).KV)
}

// IterMutDesc returns a mutable iterator over all key-value pairs in descending key order.
func (m *ArrayMap[K, V]) IterMutDesc() iter.Seq2[K, *V] {
	return diter.Split(m.iterPairDesc(), (*kv.Pair[K, V]).KVMut)
}

// RangeAsc returns key-value pairs in ascending key order within bounds.
func (m *ArrayMap[K, V]) RangeAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, V] {
	return diter.Split(m.rangePairAsc(bounds), (*kv.Pair[K, V]).KV)
}

// RangeMutAsc returns mutable key-value pairs in ascending key order within bounds.
func (m *ArrayMap[K, V]) RangeMutAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, *V] {
	return diter.Split(m.rangePairAsc(bounds), (*kv.Pair[K, V]).KVMut)
}

// RangeDesc returns key-value pairs in descending key order within bounds.
func (m *ArrayMap[K, V]) RangeDesc(bounds bound.RangeBounds[K]) iter.Seq2[K, V] {
	return diter.Split(m.rangePairDesc(bounds), (*kv.Pair[K, V]).KV)
}

// RangeMutDesc returns mutable key-value pairs in descending key order within bounds.
func (m *ArrayMap[K, V]) RangeMutDesc(bounds bound.RangeBounds[K]) iter.Seq2[K, *V] {
	return diter.Split(m.rangePairDesc(bounds), (*kv.Pair[K, V]).KVMut)
}

// KeysAsc returns an iterator over all keys in ascending order.
func (m *ArrayMap[K, V]) KeysAsc() iter.Seq[K] { return diter.Keys(m.IterAsc()) }

// ValuesAsc returns an iterator over all values in ascending key order.
func (m *ArrayMap[K, V]) ValuesAsc() iter.Seq[V] { return diter.Values(m.IterAsc()) }

// ValuesMutAsc returns a mutable iterator over all values in ascending key order.
func (m *ArrayMap[K, V]) ValuesMutAsc() iter.Seq[*V] { return diter.Values(m.IterMutAsc()) }

// KeysDesc returns an iterator over all keys in descending order.
func (m *ArrayMap[K, V]) KeysDesc() iter.Seq[K] { return diter.Keys(m.IterDesc()) }

// ValuesDesc returns an iterator over all values in descending key order.
func (m *ArrayMap[K, V]) ValuesDesc() iter.Seq[V] { return diter.Values(m.IterDesc()) }

// ValuesMutDesc returns a mutable iterator over all values in descending key order.
func (m *ArrayMap[K, V]) ValuesMutDesc() iter.Seq[*V] { return diter.Values(m.IterMutDesc()) }

// Iter returns an iterator over all key-value pairs.
func (m *ArrayMap[K, V]) Iter() iter.Seq2[K, V] { return m.IterAsc() }

// IterMut returns a mutable iterator over all key-value pairs.
func (m *ArrayMap[K, V]) IterMut() iter.Seq2[K, *V] { return m.IterMutAsc() }

// Keys returns an iterator over all keys.
func (m *ArrayMap[K, V]) Keys() iter.Seq[K] { return m.KeysAsc() }

// Values returns an iterator over all values.
func (m *ArrayMap[K, V]) Values() iter.Seq[V] { return m.ValuesAsc() }

// ValuesMut returns a mutable iterator over all values.
func (m *ArrayMap[K, V]) ValuesMut() iter.Seq[*V] { return m.ValuesMutAsc() }
