package wavelet

import (
	// "fmt"
	"github.com/ryought/tolptod/wavelet/bitvec"
	"slices"
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

// ceil(log2(x))
func Bits(x int64) int {
	var y int64
	i := 0
	for i < 64 {
		if x <= y {
			break
		}
		y = (y << 1) + 1
		i += 1
	}
	return i
}

// constructor
func NewIntWavelet(s []int64) IntWavelet {
	// max int
	S := slices.Max(s)
	D := Bits(S) + 1
	if slices.Min(s) < 0 {
		panic("negative int is not supported")
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
		var f func(int) byte
		if d == 0 {
			f = func(o int) byte {
				if s[x[o]] > 0 {
					return 1
				} else {
					return 0
				}
			}
		} else {
			f = func(o int) byte {
				return byte(s[x[o]] >> (d - 1) & 1)
			}
		}
		b := bitvec.Build(len(s), f)
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
		if d == 0 {
			if b == 0 {
				return 0
			}
		} else {
			r |= int64(b << (d - 1))
		}
		if b == 1 {
			o = w.offsets[d] + o - w.bits[d].Rank(o)
		} else {
			o = w.bits[d].Rank(o)
		}
	}
	return r
}

// Get the occurrence of x in s[0:i) for 0<=i<=n.
func (w IntWavelet) Rank(i int, x int64) int {
	if i < 0 || i > w.N() {
		panic("i should be 0<=i<=N")
	}

	// subregion [oL:oR) in B[d] of bits represents query that have the same first d bit.
	// d = 0, [L=0:R=i) is valid region.
	oL, oR := 0, i

	for d := 0; d < w.D(); d++ {
		// x to b
		var b int64
		if d == 0 {
			if x > 0 {
				b = 1
			} else {
				b = 0
			}
		} else {
			b = x >> (d - 1) & 1
		}

		// descend wavelet tree to left/right according to b
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
func (w IntWavelet) Intersect(aL, aR, bL, bR int) (int64, int, int) {
	if aL < 0 || aR > w.N() || aL > aR {
		panic("invalid search interval [aL:aR)")
	}
	if bL < 0 || bR > w.N() || bL > bR {
		panic("invalid search interval [bL:bR)")
	}

	q := NewStack()
	d := 0
	var c int64
	q.StackPush(Intersection{aL, aR, bL, bR, d, c})

	D := w.D()
	for q.Len() > 0 {
		is := q.StackPop()
		// fmt.Printf("poped [%d,%d) [%d,%d) d=%d\n", is.aL, is.aR, is.bL, is.bR, is.d)

		if is.d == D {
			// fmt.Println("found!")
			return is.c, is.aR - is.aL, is.bR - is.bL
		}

		if is.d == 0 {
			// right only
			isR := Intersection{
				aL: w.offsets[is.d] + is.aL - w.bits[is.d].Rank(is.aL),
				aR: w.offsets[is.d] + is.aR - w.bits[is.d].Rank(is.aR),
				bL: w.offsets[is.d] + is.bL - w.bits[is.d].Rank(is.bL),
				bR: w.offsets[is.d] + is.bR - w.bits[is.d].Rank(is.bR),
				d:  is.d + 1,
				c:  is.c,
			}
			if isR.IsOpen() {
				q.StackPush(isR)
			}
		} else {
			// to left
			isL := Intersection{
				aL: w.bits[is.d].Rank(is.aL),
				aR: w.bits[is.d].Rank(is.aR),
				bL: w.bits[is.d].Rank(is.bL),
				bR: w.bits[is.d].Rank(is.bR),
				d:  is.d + 1,
				c:  is.c | (0 << (is.d - 1)),
			}

			// to right
			isR := Intersection{
				aL: w.offsets[is.d] + is.aL - w.bits[is.d].Rank(is.aL),
				aR: w.offsets[is.d] + is.aR - w.bits[is.d].Rank(is.aR),
				bL: w.offsets[is.d] + is.bL - w.bits[is.d].Rank(is.bL),
				bR: w.offsets[is.d] + is.bR - w.bits[is.d].Rank(is.bR),
				d:  is.d + 1,
				c:  is.c | (1 << (is.d - 1)),
			}

			if isL.IsOpen() && isR.IsOpen() {
				if isL.Priority() > isR.Priority() {
					q.StackPush(isR)
					q.StackPush(isL)
				} else {
					q.StackPush(isL)
					q.StackPush(isR)
				}
			} else if isL.IsOpen() {
				q.StackPush(isL)
			} else if isR.IsOpen() {
				q.StackPush(isR)
			}
		}
	}

	return 0, 0, 0
}
