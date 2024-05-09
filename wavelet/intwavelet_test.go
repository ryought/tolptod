package wavelet

import (
	"github.com/ryought/tolptod/rand"
	"testing"
	"time"
)

func TestBits(t *testing.T) {
	tests := []struct {
		x    int64
		bits int
	}{
		{0b000, 0},
		{0b001, 1},
		{0b010, 2},
		{0b011, 2},
		{0b011, 2},
		{0b100, 3},
		{0b101, 3},
		{0b110, 3},
		{0b111, 3},
	}
	for _, test := range tests {
		t.Log(test.x, Bits(test.x), test.bits)
		if Bits(test.x) != test.bits {
			t.Error()
		}
	}
}

func TestIntWavelet(t *testing.T) {
	//           0  1  2  3  4  5  6  7  8  9  10 11 12
	x := []int64{0, 1, 1, 0, 2, 0, 2, 1, 4, 3, 5, 5, 9}
	w := NewIntWavelet(x)
	for i := range x {
		t.Logf("x[i=%d]\t%d", i, w.Access(i))
	}
	t.Log(w.Rank(5, 0))
	t.Log(w.Rank(5, 1))
	t.Log(w.Rank(5, 2))
	t.Log(w.Rank(5, 3))
	t.Log(w.Intersect(0, 4, 9, 13))
	t.Log(w.Intersect(0, 4, 4, 8))
	t.Log(w.Intersect(0, 4, 4, 7))
	t.Log(w.Intersect(3, 5, 5, 7))
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
		w := NewIntWavelet(x)
		t.Logf("%d\t%d us", n, time.Since(t0).Microseconds())

		t1 := time.Now()
		N := 100
		for i := 0; i < N; i++ {
			w.Intersect(0, n/2, n/2, n)
		}
		t.Logf("%d\t%d us", n, time.Since(t1).Microseconds()/int64(N))
		// t.Log(n)
	}
}

func BenchmarkIntWavelet(b *testing.B) {
	n := 10_000_000
	x := rand.RandomUint64(n, 10_000)
	w := NewIntWavelet(x)

	// b.StartTimer()
	b.ResetTimer()
	w.Intersect(0, n/2, n/2, n)
	// b.StopTimer()
}
