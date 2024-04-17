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
