package skipmap

import (
	"iter"

	"github.com/go-board/ds/bound"
	diter "github.com/go-board/ds/internal/iter"
)

// KeysAsc returns all keys in ascending order.
func (m *SkipMap[K, V]) KeysAsc() iter.Seq[K] { return diter.Keys(m.IterAsc()) }

// ValuesAsc returns all values in ascending key order.
func (m *SkipMap[K, V]) ValuesAsc() iter.Seq[V] { return diter.Values(m.IterAsc()) }

// ValuesMutAsc returns mutable values in ascending key order.
func (m *SkipMap[K, V]) ValuesMutAsc() iter.Seq[*V] { return diter.Values(m.IterMutAsc()) }

// IterAsc returns all key/value pairs in ascending key order.
func (m *SkipMap[K, V]) IterAsc() iter.Seq2[K, V] {
	all := bound.NewRangeBounds(bound.NewUnbounded[K](), bound.NewUnbounded[K]())
	return diter.Split(m.rangeNodeAsc(all), (*node[K, V]).kv)
}

// IterMutAsc returns mutable key/value pairs in ascending key order.
func (m *SkipMap[K, V]) IterMutAsc() iter.Seq2[K, *V] {
	all := bound.NewRangeBounds(bound.NewUnbounded[K](), bound.NewUnbounded[K]())
	return diter.Split(m.rangeNodeAsc(all), (*node[K, V]).kvMut)
}

// RangeAsc returns key/value pairs in ascending key order within bounds.
func (m *SkipMap[K, V]) RangeAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, V] {
	return diter.Split(m.rangeNodeAsc(bounds), (*node[K, V]).kv)
}

// RangeMutAsc returns mutable key/value pairs in ascending key order within bounds.
func (m *SkipMap[K, V]) RangeMutAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, *V] {
	return diter.Split(m.rangeNodeAsc(bounds), (*node[K, V]).kvMut)
}

// IterDesc returns all key/value pairs in descending key order.
func (m *SkipMap[K, V]) IterDesc() iter.Seq2[K, V] {
	return diter.Split(m.iterNodeDesc(), (*node[K, V]).kv)
}

// IterMutDesc returns mutable key/value pairs in descending key order.
func (m *SkipMap[K, V]) IterMutDesc() iter.Seq2[K, *V] {
	return diter.Split(m.iterNodeDesc(), (*node[K, V]).kvMut)
}

// RangeDesc returns key/value pairs in descending key order within bounds.
func (m *SkipMap[K, V]) RangeDesc(bounds bound.RangeBounds[K]) iter.Seq2[K, V] {
	return diter.Split(m.rangeNodeDesc(bounds), (*node[K, V]).kv)
}

// RangeMutDesc returns mutable key/value pairs in descending key order within bounds.
func (m *SkipMap[K, V]) RangeMutDesc(bounds bound.RangeBounds[K]) iter.Seq2[K, *V] {
	return diter.Split(m.rangeNodeDesc(bounds), (*node[K, V]).kvMut)
}

// KeysDesc returns all keys in descending order.
func (m *SkipMap[K, V]) KeysDesc() iter.Seq[K] { return diter.Keys(m.IterDesc()) }

// ValuesDesc returns all values in descending key order.
func (m *SkipMap[K, V]) ValuesDesc() iter.Seq[V] { return diter.Values(m.IterDesc()) }

// ValuesMutDesc returns mutable values in descending key order.
func (m *SkipMap[K, V]) ValuesMutDesc() iter.Seq[*V] { return diter.Values(m.IterMutDesc()) }

func (m *SkipMap[K, V]) iterNodeDesc() iter.Seq[*node[K, V]] {
	return func(yield func(*node[K, V]) bool) {
		if m == nil || m.head == nil {
			return
		}
		last := m.head.next[0]
		if last == nil {
			return
		}
		for last.next[0] != nil {
			last = last.next[0]
		}

		predecessor := func(key K) *node[K, V] {
			current := m.head
			for i := m.level; i >= 0; i-- {
				for current.next[i] != nil && m.comparator(current.next[i].Key, key) < 0 {
					current = current.next[i]
				}
			}
			if current == m.head {
				if current.next[0] != nil && m.comparator(current.next[0].Key, key) >= 0 {
					return nil
				}
			}
			return current
		}

		cur := last
		for cur != nil {
			if !yield(cur) {
				return
			}
			prev := predecessor(cur.Key)
			if prev == nil || prev == m.head {
				return
			}
			cur = prev
		}
	}
}

func (m *SkipMap[K, V]) rangeNodeAsc(bounds bound.RangeBounds[K]) iter.Seq[*node[K, V]] {
	return func(yield func(*node[K, V]) bool) {
		current := m.head.next[0]
		if v, ok := bounds.Start.Value(); ok {
			for current != nil {
				cmp := m.comparator(current.Key, v)
				if bounds.Start.IsIncluded() {
					if cmp >= 0 {
						break
					}
				} else {
					if cmp > 0 {
						break
					}
				}
				current = current.next[0]
			}
		}

		for current != nil {
			if !bounds.Contains(m.comparator, current.Key) {
				if v, ok := bounds.End.Value(); ok && m.comparator(current.Key, v) > 0 {
					break
				}
				current = current.next[0]
				continue
			}
			if !yield(current) {
				return
			}
			current = current.next[0]
		}
	}
}

func (m *SkipMap[K, V]) rangeNodeDesc(bounds bound.RangeBounds[K]) iter.Seq[*node[K, V]] {
	return func(yield func(*node[K, V]) bool) {
		for n := range m.iterNodeDesc() {
			if !bounds.Contains(m.comparator, n.Key) {
				continue
			}
			if !yield(n) {
				return
			}
		}
	}
}

// Extend inserts all key/value pairs from it into the map.
func (m *SkipMap[K, V]) Extend(it iter.Seq2[K, V]) {
	for k, v := range it {
		m.Insert(k, v)
	}
}
