package wavelet

import (
	"github.com/ryought/tolptod/rand"
	"testing"
)

func BenchmarkRadixLarge(b *testing.B) {
	s := rand.RandomDNA(1_000_000) // 1MB
	b.Log("building", len(s))
	b.StartTimer()
	b.StopTimer()
}
