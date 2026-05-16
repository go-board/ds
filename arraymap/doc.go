// Package arraymap implements an ordered map backed by a sorted slice.
//
// ArrayMap keeps pairs sorted by key using a user-provided comparator.
// Lookups use binary search, while insertions and deletions may shift elements.
// It is suitable for small to medium ordered maps with compact memory layout.
package arraymap
