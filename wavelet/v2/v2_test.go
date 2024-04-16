package v2

import (
	"testing"
)

func TestV2(t *testing.T) {
	s := []byte("ATCGAGATG")
	printBits(s)

	w := NewV2(s, 2)
	// q := w.Access(3, 3)
	// t.Log(string(q))
	// w.Rank(8, []byte("GA"))

	r := w.Top(1, 1, 0, 1)
	t.Log(r)
}
