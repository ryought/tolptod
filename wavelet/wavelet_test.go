package wavelet

import (
	"fmt"
	"testing"
)

func TestWavelet(t *testing.T) {
	xs := []int{0, 3, 1, 2, 5, 0, 2, 3, 1, 4}
	w := New(xs)
	w.Debug()
}

func TestByte(t *testing.T) {
	for i := 0; i < 1<<8; i++ {
		fmt.Printf("i=%d\t%08b %q\n", i, i, byte(i))
	}
}
