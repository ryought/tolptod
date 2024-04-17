package wavelet

import (
	"bytes"
	"github.com/ryought/tolptod/rand"
	"index/suffixarray"
	"testing"
	"time"
)

func TestWaveletLarge(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	ns := []int{
		1_000, // 1KB
		10_000,
		100_000,
		1_000_000, // 1MB
		// 10_000_000,
	}

	for _, n := range ns {
		s := rand.RandomDNA(n)
		t.Log(n)

		t0 := time.Now()
		w := NewDNAWavelet(s, 100)
		t.Logf("wave %d ms", time.Since(t0).Milliseconds())

		t1 := time.Now()
		suffixarray.New(s)
		t.Logf("sais %d ms", time.Since(t1).Milliseconds())

		x := w.Access(10, 100)
		if !bytes.Equal(x, s[10:10+100]) {
			t.Error()
		}

	}
}

func BenchmarkWaveletLarge(b *testing.B) {
	s := rand.RandomDNA(1_000_000) // 1MB
	b.Log("building", len(s))
	b.StartTimer()
	w := New(s, 100)
	b.StopTimer()
	t := w.Access(0, 100)
	b.Log(t)
	b.Log(s[:100])
}
