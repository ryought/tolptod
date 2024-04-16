package wavelet

import (
	"fmt"
	"github.com/ryought/tolptod/wavelet/bitvec"
	"math"
)

type Wavelet struct {
	vecs    []bitvec.BitVec
	offsets []int
}

func alphabetsize(xs []int) int {
	ret := 0
	for i := range xs {
		if xs[i] > ret {
			ret = xs[i] + 1
		}
	}
	return ret
}

func (w Wavelet) Debug() {
	for d := range w.vecs {
		fmt.Println("d", d)
		// w.vecs[d].Debug()
	}
}

func New(xs []int) Wavelet {
	// depth=ceil(lg(size))
	size := alphabetsize(xs)
	depth := int(math.Log2(float64(size))) + 1

	vecs := make([]bitvec.BitVec, depth)
	offsets := make([]int, depth)
	values := make([]int, len(xs))
	offset := 0

	for d := 0; d < depth; d++ {
		if d == 0 {
			copy(values, xs)
		} else {
			// sort values by bits
			// vecs[d-1]
		}

		bits := make([]bool, len(values))
		offset = 0
		for i := range values {
			x := xs[i] >> d & 1
			if x == 0 {
				offset += 1
			}
			// fmt.Printf("x %d %08b %d\n", xs[i], xs[i], x)
			bits[i] = bitvec.Int2Bit(x)
		}

		// fmt.Println("x", xs, bits)
		vecs[d] = bitvec.New(bits)
		offsets[d] = 0
	}

	return Wavelet{vecs, offsets}
}

func (w Wavelet) Intersect(i, j, k, l int) {
	fmt.Println("hoge")
}
