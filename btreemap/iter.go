package btreemap

import (
	"iter"

	"github.com/go-board/ds/bound"
	diter "github.com/go-board/ds/internal/iter"
)

// IterAsc returns an iterator over all key/value pairs in ascending key order.
func (m *BTreeMap[K, V]) IterAsc() iter.Seq2[K, V] {
	return diter.Split(m.btree.IterAsc(), (*node[K, V]).kv)
}

// IterMutAsc returns a mutable iterator over all key/value pairs in ascending key order.
func (m *BTreeMap[K, V]) IterMutAsc() iter.Seq2[K, *V] {
	return diter.Split(m.btree.IterAsc(), (*node[K, V]).kvMut)
}

// IterDesc returns an iterator over all key/value pairs in descending key order.
func (m *BTreeMap[K, V]) IterDesc() iter.Seq2[K, V] {
	return diter.Split(m.btree.IterDesc(), (*node[K, V]).kv)
}

// IterMutDesc returns a mutable iterator over all key/value pairs in descending key order.
func (m *BTreeMap[K, V]) IterMutDesc() iter.Seq2[K, *V] {
	return diter.Split(m.btree.IterDesc(), (*node[K, V]).kvMut)
}

// RangeAsc returns an iterator over key/value pairs in ascending key order within bounds.
func (m *BTreeMap[K, V]) RangeAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, V] {
	return diter.Split(m.rangeNode(bounds, false), (*node[K, V]).kv)
}

// RangeMutAsc returns a mutable iterator over key/value pairs in ascending key order within bounds.
func (m *BTreeMap[K, V]) RangeMutAsc(bounds bound.RangeBounds[K]) iter.Seq2[K, *V] {
	return diter.Split(m.rangeNode(bounds, false), (*node[K, V]).kvMut)
}

// RangeDesc returns an iterator over key/value pairs in descending key order within bounds.
func (m *BTreeMap[K, V]) RangeDesc(bounds bound.RangeBounds[K]) iter.Seq2[K, V] {
	return diter.Split(m.rangeNode(bounds, true), (*node[K, V]).kv)
}

// RangeMutDesc returns a mutable iterator over key/value pairs in descending key order within bounds.
func (m *BTreeMap[K, V]) RangeMutDesc(bounds bound.RangeBounds[K]) iter.Seq2[K, *V] {
	return diter.Split(m.rangeNode(bounds, true), (*node[K, V]).kvMut)
}

func (m *BTreeMap[K, V]) rangeNode(bounds bound.RangeBounds[K], desc bool) iter.Seq[*node[K, V]] {
	nb := mapBounds[K, V](bounds)
	if desc {
		return m.btree.RangeDesc(nb)
	}
	return m.btree.RangeAsc(nb)
}

func mapBounds[K, V any](b bound.RangeBounds[K]) bound.RangeBounds[*node[K, V]] {
	mapOne := func(src bound.Bound[K]) bound.Bound[*node[K, V]] {
		switch src.Kind() {
		case bound.Unbounded:
			return bound.NewUnbounded[*node[K, V]]()
		case bound.Included:
			v, _ := src.Value()
			return bound.NewIncluded(&node[K, V]{Key: v})
		default:
			v, _ := src.Value()
			return bound.NewExcluded(&node[K, V]{Key: v})
		}
	}
	return bound.NewRangeBounds(mapOne(b.Start), mapOne(b.End))
}

func (m *BTreeMap[K, V]) KeysAsc() iter.Seq[K]       { return diter.Keys(m.IterAsc()) }
func (m *BTreeMap[K, V]) ValuesAsc() iter.Seq[V]     { return diter.Values(m.IterAsc()) }
func (m *BTreeMap[K, V]) ValuesMutAsc() iter.Seq[*V] { return diter.Values(m.IterMutAsc()) }
func (m *BTreeMap[K, V]) KeysDesc() iter.Seq[K]      { return diter.Keys(m.IterDesc()) }
func (m *BTreeMap[K, V]) ValuesDesc() iter.Seq[V]    { return diter.Values(m.IterDesc()) }
func (m *BTreeMap[K, V]) ValuesMutDesc() iter.Seq[*V] {
	return diter.Values(m.IterMutDesc())
}

func (m *BTreeMap[K, V]) Extend(it iter.Seq2[K, V]) {
	for k, v := range it {
		m.Insert(k, v)
	}
}
