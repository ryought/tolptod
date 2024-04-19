package wavelet

import (
	"github.com/ryought/tolptod/rand"
	"testing"
)

func BenchmarkWavelet1MB(b *testing.B) {
	n := 10_000_000 // 10MB
	K := 30
	s := rand.RandomDNA(n)
	b.Log("build")
	w := NewDNAWavelet(s, K)
	b.Log("query")
	x, y := w.Intersect(0, n/2, n/2, n, K)
	b.Log(x, y)
}
