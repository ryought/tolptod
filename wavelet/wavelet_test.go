package wavelet

import (
	"testing"
)

func TestV2(t *testing.T) {
	s := []byte("ATCGAGATG")
	printBits(s)

	w := NewV2(s, 3)
	q := w.Access(3, 3)
	t.Log(string(q))
	n := w.Rank(5, []byte("ATT"))
	t.Log(n)

	// r := w.Top(1, 1, 0, 1)
	// t.Log(r)
}
