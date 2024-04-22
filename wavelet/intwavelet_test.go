package wavelet

import (
	"github.com/ryought/tolptod/rand"
	"testing"
)

func TestIntWavelet(t *testing.T) {
	//           0  1  2  3  4  5  6  7  8  9  10 11 12
	x := []int64{0, 1, 1, 0, 2, 0, 2, 1, 4, 3, 5, 5, 9}
	w := NewIntWavelet(x, 4)
	for i := range x {
		t.Log(i, w.Access(i))
	}
	t.Log(w.Rank(5, 1))
	t.Log(w.Intersect(0, 4, 9, 13))
}

func BenchmarkIntWavelet(b *testing.B) {
	n := 100_000_000
	x := rand.RandomUint64(n, 10_000)
	w := NewIntWavelet(x, 15)
	b.StartTimer()
	w.Intersect(0, n/2, n/2, n)
	b.StopTimer()
}
