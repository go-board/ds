package diter

import "iter"

func Keys[K, V any](it iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range it {
			if !yield(key) {
				return
			}
		}
	}
}

func Values[K, V any](it iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, value := range it {
			if !yield(value) {
				return
			}
		}
	}
}

func Split[E, K, V any](it iter.Seq[E], f func(E) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := range it {
			key, value := f(e)
			if !yield(key, value) {
				return
			}
		}
	}
}
