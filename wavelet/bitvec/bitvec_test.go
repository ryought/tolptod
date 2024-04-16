package bitvec

import (
	"testing"
)

func TestHoge(t *testing.T) {
	bits := []bool{true, false, true, false, false, true, false, true, true}
	bv := New(bits)
	bv.Debug()

	tests := []struct {
		b bool
		i int
		r int
	}{
		{true, 0, 1},
		{true, 4, 2},
		{true, 5, 3},
		{true, 8, 5},
		{false, 0, 0},
		{false, 1, 1},
		{false, 8, 4},
	}
	for _, test := range tests {
		if bv.Rank(test.b, test.i) != test.r {
			t.Errorf("not true rank(%t,%d)!=%d", test.b, test.i, test.r)
		}
	}
}
