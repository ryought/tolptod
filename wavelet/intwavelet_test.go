package wavelet

import (
	"github.com/ryought/tolptod/rand"
	"testing"
	"time"
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

func TestIntWaveletBySize(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	ns := []int{
		1_000,
		10_000,
		100_000,
		1_000_000,
		10_000_000,
		100_000_000,
	}

	for _, n := range ns {
		x := rand.RandomUint64(n, 2<<15)
		t0 := time.Now()
		w := NewIntWavelet(x, 15)
		t.Logf("%d\t%d us", n, time.Since(t0).Microseconds())

		t1 := time.Now()
		for i := 0; i < 100; i++ {
			w.Intersect(0, n/2, n/2, n)
		}
		t.Logf("%d\t%d us", n, time.Since(t1).Microseconds())
		// t.Log(n)
	}
}

func BenchmarkIntWavelet(b *testing.B) {
	n := 10_000_000
	x := rand.RandomUint64(n, 10_000)
	w := NewIntWavelet(x, 15)
	b.StartTimer()
	w.Intersect(0, n/2, n/2, n)
	b.StopTimer()
}
