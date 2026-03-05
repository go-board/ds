package bound

import "testing"

func intCmp(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func TestBoundKinds(t *testing.T) {
	ub := NewUnbounded[int]()
	if !ub.IsUnbounded() || ub.IsIncluded() || ub.IsExcluded() {
		t.Fatalf("unexpected unbounded flags: %+v", ub)
	}
	if _, ok := ub.Value(); ok {
		t.Fatalf("unbounded should not expose value")
	}

	in := NewIncluded(10)
	if !in.IsIncluded() || in.IsUnbounded() || in.IsExcluded() {
		t.Fatalf("unexpected included flags: %+v", in)
	}
	if v, ok := in.Value(); !ok || v != 10 {
		t.Fatalf("included value mismatch: %v %v", v, ok)
	}

	ex := NewExcluded(20)
	if !ex.IsExcluded() || ex.IsUnbounded() || ex.IsIncluded() {
		t.Fatalf("unexpected excluded flags: %+v", ex)
	}
}

func TestRangeBoundsContains(t *testing.T) {
	r := NewRangeBounds(NewIncluded(1), NewExcluded(5))
	cases := map[int]bool{
		0: false,
		1: true,
		4: true,
		5: false,
	}
	for v, want := range cases {
		if got := r.Contains(intCmp, v); got != want {
			t.Fatalf("contains(%d)=%v, want %v", v, got, want)
		}
	}

	r2 := NewRangeBounds(NewUnbounded[int](), NewIncluded(3))
	if !r2.Contains(intCmp, -100) || !r2.Contains(intCmp, 3) || r2.Contains(intCmp, 4) {
		t.Fatalf("unexpected contains result for unbounded lower range")
	}
}

func TestRangeBoundsIsValid(t *testing.T) {
	if !NewRangeBounds(NewIncluded(2), NewIncluded(2)).IsValid(intCmp) {
		t.Fatalf("[2,2] should be valid")
	}
	if NewRangeBounds(NewIncluded(2), NewExcluded(2)).IsValid(intCmp) {
		t.Fatalf("[2,2) should be invalid")
	}
	if NewRangeBounds(NewExcluded(3), NewIncluded(2)).IsValid(intCmp) {
		t.Fatalf("(3,2] should be invalid")
	}
	if !NewRangeBounds(NewUnbounded[int](), NewExcluded(0)).IsValid(intCmp) {
		t.Fatalf("unbounded side should be valid")
	}
}
