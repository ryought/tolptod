package wavelet

import (
	"bytes"
	"slices"
	"testing"
)

func TestCreateRank(t *testing.T) {
	b := []byte{0, 1, 1, 0, 1, 0, 0, 1, 0}
	R := []int{0, 1, 1, 1, 2, 2, 3, 4, 4, 5}
	r := createRank(b)
	t.Log(r)
	if !slices.Equal(r, R) {
		t.Errorf("rank is not correct")
	}
}

func TestAccess(t *testing.T) {
	s := []byte("ATCGAGATG")
	printBits(s)

	w := NewV2(s, 3)

	// access
	q := w.Access(8, 1)
	t.Log(string(q))

	tests := []struct {
		i int
		n int
	}{
		{0, 0},
		{1, 0},
		{2, 0},
		{0, 1},
		{8, 1},
		{3, 2},
		{0, 3},
		{6, 3},
	}
	for _, test := range tests {
		x := w.Access(test.i, test.n)
		t.Logf("i=%d n=%d %s\n", test.i, test.n, x)
		if !bytes.Equal(x, s[test.i:test.i+test.n]) {
			t.Errorf("access is not correct i=%d n=%d", test.i, test.n)
		}
	}

	// n := w.Rank(5, []byte("ATT"))
	// t.Log(n)

	// r := w.Top(1, 1, 0, 1)
	// t.Log(r)
}

func TestRank(t *testing.T) {
	s := []byte("ATCGAGATG")
	n := len(s)
	printBits(s)
	w := NewV2(s, 3)

	c := w.Rank(3, []byte("G"))
	t.Log(c)

	tests := []struct {
		i int
		q []byte
		c int
	}{
		// for all
		{n, []byte("GA"), 2},
		{n, []byte("A"), 3},
		{n, []byte("C"), 1},
		{n, []byte("G"), 3},
		{n, []byte("T"), 2},
		{n, []byte("Z"), 0},
		{n, []byte("CGA"), 1},
		{n, []byte("GGG"), 0},
		// subregion
		{3, []byte("A"), 1},
		{3, []byte("G"), 0},
		{3, []byte("GGG"), 0},
		{3, []byte("ATC"), 1},
		{3, []byte("TCG"), 1},
	}

	for _, test := range tests {
		c := w.Rank(test.i, test.q)
		t.Logf("i=%d q=%s c=%d\n", test.i, string(test.q), c)
		if c != test.c {
			t.Errorf("")
		}
	}
}

func TestTop(t *testing.T) {
	s := []byte("CATCGAGATGAGA")
	n := len(s)
	printBits(s)
	w := NewV2(s, 3)

	{
		q, c := w.Top(0, n, 1)
		t.Log(string(q), c)
		if !bytes.Equal(q, []byte("A")) || c != 5 {
			t.Errorf("")
		}
	}

	{
		q, c := w.Top(0, 4, 1)
		t.Log(string(q), c)
		if !bytes.Equal(q, []byte("C")) || c != 2 {
			t.Errorf("")
		}
	}

	{
		q, c := w.Top(0, n, 2)
		t.Log(string(q), c)
		if !bytes.Equal(q, []byte("GA")) || c != 4 {
			t.Errorf("")
		}
	}

	{
		q, c := w.Top(0, n, 3)
		t.Log(string(q), c)
		if !bytes.Equal(q, []byte("AGA")) || c != 2 {
			t.Errorf("")
		}
	}
}

func TestTopWithDoller(t *testing.T) {
	s := []byte("CATC$$$$GAC")
	n := len(s)
	printBits(s)
	w := NewV2(s, 3)
	{
		q, c := w.Top(0, n, 1)
		t.Log(string(q), c)
		if !bytes.Equal(q, []byte("C")) || c != 3 {
			t.Error()
		}
	}
}

func TestIntersectionEnd(t *testing.T) {
	{
		s := []byte("AAA$GGG$")
		w := NewV2(s, 3)
		a, b := w.Intersect(0, 4, 4, 8, 1)
		t.Log(a, b)
		if a != 0 && b != 0 {
			t.Error()
		}
	}

	{
		s := []byte("AAA$GGA$")
		w := NewV2(s, 3)
		a, b := w.Intersect(0, 4, 4, 8, 1)
		t.Log(a, b)
		// A is in common
		if a != 3 && b != 1 {
			t.Error()
		}

		a, b = w.Intersect(0, 4, 4, 8, 2)
		t.Log(a, b)
		if a != 0 && b != 0 {
			t.Error()
		}
	}
}

func TestIntersection(t *testing.T) {
	s := []byte("AAAAAGGGGG")
	w := NewV2(s, 3)

	tests := []struct {
		aL int
		aR int
		bL int
		bR int
		K  int
		a  int
		b  int
	}{
		{0, 5, 5, 10, 1, 0, 0},  // [0:5) [5:10) have no common char
		{0, 5, 0, 10, 1, 5, 5},  // [0:5) [0:10) have five "XX"s
		{0, 6, 4, 10, 2, 1, 1},  // XY in common
		{5, 9, 0, 9, 2, 4, 4},   // four YY in common
		{5, 10, 0, 10, 2, 4, 4}, // four GG in common
	}

	for _, test := range tests {
		a, b := w.Intersect(test.aL, test.aR, test.bL, test.bR, test.K)
		t.Log(test, a, b)
		if a != test.a || b != test.b {
			t.Error("not correct")
		}
	}

	{
		a, b := w.Intersect(5, 10, 0, 10, 2)
		t.Log(a, b)
	}
}
