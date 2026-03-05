package bound

// Kind describes the boundary type.
type Kind uint8

const (
	// Unbounded means no boundary on this side.
	Unbounded Kind = iota
	// Included means the boundary value is included in the range.
	Included
	// Excluded means the boundary value is excluded from the range.
	Excluded
)

// Bound represents one side of a range.
//
// For Unbounded bounds, the value is ignored.
type Bound[T any] struct {
	kind  Kind
	value T
}

// NewUnbounded creates an unbounded boundary.
func NewUnbounded[T any]() Bound[T] {
	return Bound[T]{kind: Unbounded}
}

// NewIncluded creates an inclusive boundary.
func NewIncluded[T any](value T) Bound[T] {
	return Bound[T]{kind: Included, value: value}
}

// NewExcluded creates an exclusive boundary.
func NewExcluded[T any](value T) Bound[T] {
	return Bound[T]{kind: Excluded, value: value}
}

// Kind returns the boundary kind.
func (b Bound[T]) Kind() Kind {
	return b.kind
}

// Value returns the boundary value when bounded.
func (b Bound[T]) Value() (T, bool) {
	if b.kind == Unbounded {
		var zero T
		return zero, false
	}
	return b.value, true
}

// IsUnbounded reports whether the boundary is unbounded.
func (b Bound[T]) IsUnbounded() bool { return b.kind == Unbounded }

// IsIncluded reports whether the boundary is inclusive.
func (b Bound[T]) IsIncluded() bool { return b.kind == Included }

// IsExcluded reports whether the boundary is exclusive.
func (b Bound[T]) IsExcluded() bool { return b.kind == Excluded }

// RangeBounds describes a full range by start and end boundaries.
type RangeBounds[T any] struct {
	Start Bound[T]
	End   Bound[T]
}

// NewRangeBounds creates a range bounds instance.
func NewRangeBounds[T any](start, end Bound[T]) RangeBounds[T] {
	return RangeBounds[T]{Start: start, End: end}
}

// Contains reports whether value is inside the range.
//
// Comparator contract:
//   - cmp(a, b) < 0 means a < b
//   - cmp(a, b) = 0 means a == b
//   - cmp(a, b) > 0 means a > b
func (r RangeBounds[T]) Contains(cmp func(T, T) int, value T) bool {
	if v, ok := r.Start.Value(); ok {
		c := cmp(value, v)
		if r.Start.IsIncluded() {
			if c < 0 {
				return false
			}
		} else {
			if c <= 0 {
				return false
			}
		}
	}

	if v, ok := r.End.Value(); ok {
		c := cmp(value, v)
		if r.End.IsIncluded() {
			if c > 0 {
				return false
			}
		} else {
			if c >= 0 {
				return false
			}
		}
	}

	return true
}

// IsValid reports whether the range itself is internally consistent.
func (r RangeBounds[T]) IsValid(cmp func(T, T) int) bool {
	start, okStart := r.Start.Value()
	end, okEnd := r.End.Value()
	if !okStart || !okEnd {
		return true
	}

	c := cmp(start, end)
	if c < 0 {
		return true
	}
	if c > 0 {
		return false
	}

	return r.Start.IsIncluded() && r.End.IsIncluded()
}
