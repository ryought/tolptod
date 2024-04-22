package wavelet

import (
	// "fmt"
	"github.com/ryought/tolptod/wavelet/bitvec"
	// "time"
)

type IntWavelet struct {
	bits    []bitvec.BitVecV2
	offsets []int
}

// depth
func (w IntWavelet) D() int {
	return len(w.offsets)
}

// length of s
func (w IntWavelet) N() int {
	return w.bits[0].Size()
}

// constructor
func NewIntWavelet(s []int64, D int) IntWavelet {
	if D < 1 || D > 64 {
		panic("W should be 1<=D<=64.")
	}
	bits := make([]bitvec.BitVecV2, D)
	offsets := make([]int, D)

	// X0
	x := make([]int, len(s))
	tmp := make([]int, len(s))
	for o := range x {
		x[o] = o
	}

	for d := 0; d < D; d++ {
		b := bitvec.Build(len(s), func(o int) byte {
			return byte(s[x[o]] >> d & 1)
		})
		offset := b.Rank(len(s))
		// fmt.Printf("bitvec in %d ms\n", time.Since(t0).Milliseconds())
		// fmt.Println(x, "X")
		// fmt.Println(b, "B", offset)

		// sort X0 to X1
		// t1 := time.Now()
		o0, o1 := 0, offset
		for o := range x {
			if b.Get(o) == 0 {
				tmp[o0] = x[o]
				o0 += 1
			} else {
				tmp[o1] = x[o]
				o1 += 1
			}
		}
		// fmt.Printf("sort in %d ms\n", time.Since(t1).Milliseconds())
		x, tmp = tmp, x
		// fmt.Println(x, "X'")

		bits[d] = b
		offsets[d] = offset
	}

	return IntWavelet{bits, offsets}
}

// Get s[i] for 0<=i<n
func (w IntWavelet) Access(i int) int64 {
	if i < 0 || i >= w.N() {
		panic("Access: index out of range")
	}

	// index in w.bits[d].
	// B[0][i]
	r := int64(0)
	o := i
	for d := 0; d < w.D(); d++ {
		b := w.bits[d].Get(o)
		r |= int64(b << d)
		if b == 1 {
			o = w.offsets[d] + o - w.bits[d].Rank(o)
		} else {
			o = w.bits[d].Rank(o)
		}
	}
	return r
}

// Get the occurrence of query in s[0:i) for 0<=i<=n.
func (w IntWavelet) Rank(i int, x int64) int {
	if i < 0 || i > w.N() {
		panic("i should be 0<=i<=N")
	}

	// subregion [oL:oR) in B[d] of bits represents query that have the same first d bit.
	// d = 0, [L=0:R=i) is valid region.
	oL, oR := 0, i

	for d := 0; d < w.D(); d++ {
		b := x >> d & 1
		if oL == oR {
			return 0
		}
		if b == 1 {
			oL = w.offsets[d] + oL - w.bits[d].Rank(oL)
			oR = w.offsets[d] + oR - w.bits[d].Rank(oR)
		} else {
			oL = w.bits[d].Rank(oL)
			oR = w.bits[d].Rank(oR)
		}
	}

	return oR - oL
}

// Find common int in S[aL:aR) and S[bL:bR).
func (w IntWavelet) Intersect(aL, aR, bL, bR int) (int, int) {
	if aL < 0 || aR > w.N() || aL > aR {
		panic("invalid search interval [aL:aR)")
	}
	if bL < 0 || bR > w.N() || bL > bR {
		panic("invalid search interval [bL:bR)")
	}

	q := NewStack()
	d := 0
	var c byte
	q.StackPush(Intersection{aL, aR, bL, bR, d, c})

	D := w.D()
	for q.Len() > 0 {
		is := q.StackPop()
		// fmt.Printf("poped [%d,%d) [%d,%d) d=%d\n", is.aL, is.aR, is.bL, is.bR, is.d)

		if is.d == D {
			// fmt.Println("found!")
			return is.aR - is.aL, is.bR - is.bL
		}

		// to left
		{
			aL := w.bits[is.d].Rank(is.aL)
			aR := w.bits[is.d].Rank(is.aR)
			bL := w.bits[is.d].Rank(is.bL)
			bR := w.bits[is.d].Rank(is.bR)
			d := is.d + 1
			c := is.c | (0 << d)
			// fmt.Printf("L [%d,%d) [%d,%d) d=%d\n", aL, aR, bL, bR, d)
			if aL < aR && bL < bR {
				q.StackPush(Intersection{aL, aR, bL, bR, d, c})
			}
		}

		// to right
		{
			aL := w.offsets[is.d] + is.aL - w.bits[is.d].Rank(is.aL)
			aR := w.offsets[is.d] + is.aR - w.bits[is.d].Rank(is.aR)
			bL := w.offsets[is.d] + is.bL - w.bits[is.d].Rank(is.bL)
			bR := w.offsets[is.d] + is.bR - w.bits[is.d].Rank(is.bR)
			d := is.d + 1
			c := is.c | (1 << d)
			// fmt.Printf("L [%d,%d) [%d,%d) d=%d\n", aL, aR, bL, bR, d)
			if aL < aR && bL < bR {
				q.StackPush(Intersection{aL, aR, bL, bR, d, c})
			}
		}
	}

	return 0, 0
}
