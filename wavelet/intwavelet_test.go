package wavelet

import (
	"github.com/ryought/tolptod/rand"
	mathrand "math/rand"
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

func NaiveOccurrence(s []int64, x int64) int {
	count := 0
	for _, v := range s {
		if v == x {
			count += 1
		}
	}
	return count
}

func CreateTestWavelet(t *testing.T, x []int64) {
	N := len(x)
	w := NewIntWavelet(x)

	// access
	for i := 0; i < N; i++ {
		got := w.Access(i)
		expected := x[i]
		if got != expected {
			t.Error("notcorrect", i)
		}
	}

	// rank
	for i := 0; i <= N; i++ {
		for a := 0; a < 10; a++ {
			got := w.Rank(i, int64(a))
			expected := NaiveOccurrence(x[:i], int64(a))
			if got != expected {
				t.Error("notcorrect", i, a)
			}
		}
	}

	// intersect
	for aL := 0; aL <= N; aL++ {
		for aR := aL; aR <= N; aR++ {
			for bL := 0; bL <= N; bL++ {
				for bR := bL; bR <= N; bR++ {
					c, a, b := w.Intersect(aL, aR, bL, bR)
					ea := NaiveOccurrence(x[aL:aR], c)
					eb := NaiveOccurrence(x[bL:bR], c)
					if c > 0 {
						if ea != a || eb != b {
							t.Error("incorrect", aL, aR, bL, bR)
						}
					}
				}
			}
		}
	}

}

func TestIntWaveletRandom(t *testing.T) {
	x := rand.RandomUint64(50, 10)
	CreateTestWavelet(t, x)
}

func TestIntWavelet(t *testing.T) {
	//           0  1  2  3  4  5  6  7  8  9  10 11 12
	x := []int64{0, 1, 1, 0, 2, 0, 2, 1, 4, 3, 5, 5, 9}
	N := len(x)
	w := NewIntWavelet(x)
	for i := range x {
		t.Logf("x[i=%d]\t%d", i, w.Access(i))
	}

	for i := 0; i < N; i++ {
		got := w.Access(i)
		expected := x[i]
		t.Logf("Access(i=%d)=%d %d", i, got, expected)
		if got != expected {
			t.Error("notcorrect", i)
		}
	}

	for i := 0; i <= N; i++ {
		for a := 0; a < 10; a++ {
			got := w.Rank(i, int64(a))
			expected := NaiveOccurrence(x[:i], int64(a))
			t.Logf("Rank(i=%d,x=%d)=%d %d", i, a, got, expected)
			if got != expected {
				t.Error("notcorrect", i, a)
			}
		}
	}

	for aL := 0; aL <= N; aL++ {
		for aR := aL; aR <= N; aR++ {
			for bL := 0; bL <= N; bL++ {
				for bR := bL; bR <= N; bR++ {
					c, a, b := w.Intersect(aL, aR, bL, bR)
					ea := NaiveOccurrence(x[aL:aR], c)
					eb := NaiveOccurrence(x[bL:bR], c)
					t.Logf("Intersect(%d,%d,%d,%d)\t%d\t(%d %d)\t(%d %d)", aL, aR, bL, bR, c, a, b, ea, eb)
					if c > 0 {
						if ea != a || eb != b {
							t.Error("incorrect", aL, aR, bL, bR)
						}
					}
				}
			}
		}
	}

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

func TestIntWaveletPerformanceCheck(t *testing.T) {
	sigma := int64(10_000) // alphabet size
	n := 1_000_000         // 1MB
	x := rand.RandomUint64(n, sigma)
	w := NewIntWavelet(x)
	bands := []int{1, 10, 100, 1_000, 10_000, 100_000}
	r := mathrand.New(mathrand.NewSource(0))
	for _, band := range bands {
		var duration int64
		N := 1000
		count := 0
		for trial := 0; trial < N; trial++ {
			i := r.Intn(n - band)
			j := r.Intn(n - band)

			start := time.Now()
			_, ca, cb := w.Intersect(i, i+band, j, j+band)
			duration += time.Since(start).Nanoseconds()
			count += ca
			count += cb
		}
		t.Logf("%10d\t%d\t%d", band, duration/int64(N), count)
	}
	t.Logf("band size\ttime ns\tn matches")
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
