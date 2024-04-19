package wavelet

import (
	"bytes"
	"github.com/ryought/tolptod/rand"
	"github.com/ryought/tolptod/suffixarray"
	"testing"
	"time"
)

func TestWaveletVsSuffixArray(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	ns := []int{
		// 1_000, // 1KB
		// 10_000,
		// 100_000,
		1_000_000, // 1MB
		// 10_000_000, // 10MB
	}

	for _, n := range ns {
		s := rand.RandomDNA(n)
		t.Log(n)

		t0 := time.Now()
		K := 40
		w := NewDNAWavelet(s, K)
		t.Logf("wave %d ms", time.Since(t0).Milliseconds())

		t1 := time.Now()
		sa := suffixarray.New(s)
		t.Logf("sais %d ms", time.Since(t1).Milliseconds())

		x := w.Access(10, K)
		if !bytes.Equal(x, s[10:10+K]) {
			t.Error()
		}

		t2 := time.Now()
		a, b := w.Intersect(0, n/2, n/2, n, K)
		t.Log(a, b)
		t.Logf("wave intersect %d ms", time.Since(t2).Milliseconds())

		t3 := time.Now()
		a, b = sa.Intersect(0, n/2, n/2, n, K)
		t.Log(a, b)
		t.Logf("sais intersect %d ms", time.Since(t3).Milliseconds())
	}
}
