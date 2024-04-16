package wavelet

// reference
// https://miti-7.hatenablog.com/entry/2018/04/28/152259
// https://scrapbox.io/koki/Wavelet_Matrix
// https://www.slideshare.net/pfi/ss-15916040

import (
	"fmt"
)

type WaveletV2 struct {
	bits    [][]byte
	ranks   [][]int
	offsets []int
}

// i-th (0<=i<8) bit of s[k].
func ix(s []byte, k int, i int) byte {
	if k < len(s) && i < 8 {
		return s[k] >> i & 1
	} else {
		return 0
	}
}

func printBits(s []byte) {
	for i := range s {
		fmt.Printf("s[i=%d]=%c %08b\n", i, s[i], s[i])
	}
}

// max k size
func (w WaveletV2) K() int {
	return w.D() / 8
}

// depth
func (w WaveletV2) D() int {
	return len(w.offsets)
}

// length of s
func (w WaveletV2) N() int {
	return len(w.bits[0])
}

// constructor
func NewV2(s []byte, K int) WaveletV2 {
	D := K * 8
	bits := make([][]byte, D)
	ranks := make([][]int, D)
	offsets := make([]int, D)

	// X0
	x := make([]int, len(s))
	tmp := make([]int, len(s))
	for o := range x {
		x[o] = o
	}
	// B
	var b []byte

	for k := 0; k < K; k++ {
		for i := 0; i < 8; i++ {
			// depth
			d := k*8 + i
			fmt.Printf("k=%d i=%d d=%d\n", k, i, d)

			// B0
			b = make([]byte, len(s))
			offset := 0
			for o := range b {
				if ix(s, x[o]+k, i) == 0 {
					b[o] = 0
					offset += 1
				} else {
					b[o] = 1
				}
			}
			fmt.Println(x, "X")
			fmt.Println(b, "B", offset)

			// sort X0 to X1
			o0, o1 := 0, offset
			for o := range x {
				if b[o] == 0 {
					tmp[o0] = x[o]
					o0 += 1
				} else {
					tmp[o1] = x[o]
					o1 += 1
				}
			}
			copy(x, tmp)
			fmt.Println(x, "X'")

			bits[d] = b
			ranks[d] = createRank(b)
			fmt.Println(ranks[d], "rank")
			offsets[d] = offset
		}
	}

	return WaveletV2{bits, ranks, offsets}
}

// Create rank array.
// rank[i] is the occurrence count of 0s in b[0]...b[i] (both-closed slice)
// the occurrence of 1s in b[0]...b[i] is (i+1)-rank[i].
//
// For example,
// createRank([0 1 0 1 0 0])
// will be    [1 1 2 2 3 4]
func createRank(b []byte) []int {
	rank := make([]int, len(b))
	// current occurrence of zero in b
	count := 0
	for i := range b {
		if b[i] == 0 {
			count += 1
		}
		rank[i] = count
	}
	return rank
}

// Get s[i:i+K) subsequence.
func (w WaveletV2) Access(i int, K int) []byte {
	if i < 0 || i >= w.N() {
		panic("Access index out of range")
	}
	if K > w.K() {
		panic("k cannot be grater than Wavelet.K")
	}

	// subsequence to be returned
	s := make([]byte, K)

	// index in w.bits[d].
	o := i
	for k := 0; k < K; k++ {
		var c byte
		for i := 0; i < 8; i++ {
			d := k*8 + i
			// fmt.Printf("k=%d i=%d d=%d\n", k, i, d)
			b := w.bits[d][o]
			c = c | (b << i)
			// fmt.Printf("o=%d b=%d c=%d\n", o, b, c)
			// fmt.Println(w.bits[d], "B")
			// fmt.Println(w.ranks[d], "ranks")
			// fmt.Println(w.offsets[d], "offset")
			if b == 1 {
				// go to right (1).
				// number of zeros (offset) + number of ones before o.
				//
				//                 o        o=4
				// B[d-1]  0 1 0 1 1 1 0 0  offset=4
				// rank0   1 1 2 2 2 2 3 4
				// rank1   0 1 1 2 3 4 4 4  rank1[i] = (i+1) - rank0[i]
				// B[d]    0 0 0 0 1 1 1 1
				//         <-----> <--->
				//                     o    new o=6
				//         offset=4
				//                 rank1[o]=3
				//
				o = w.offsets[d] + o - w.ranks[d][o]
			} else {
				// go to left (0).
				//                 o        o=4
				// B[d-1]  0 1 0 1 0 1 1 0
				// rank0   1 1 2 2 3 3 3 4  rank0[4]=3
				// B[d]    0 0 0 0 1 1 1 1
				//         <--->
				//             o        new o=2
				o = w.ranks[d][o] - 1
			}
			// fmt.Println("newo", o)
		}
		s[k] = c
	}
	return s
}

// Get the occurrence of query in s[0:i) = s[0]...s[i-1]
func (w WaveletV2) Rank(i int, query []byte) int {
	if len(query) > w.K() {
		panic("query cannot be longer than Wavelet.K")
	}
	if !(i >= 0 && i <= w.N()) {
		panic("i should be 0<=i<=N")
	}

	// subregion [oL:oR) of bits represents query that have the same first d bit.
	oL, oR := 0, i

	for k := range query {
		for i := 0; i < 8; i++ {
			d := k*8 + i
			b := ix(query, k, i)
			fmt.Printf("k=%d i=%d b=%d d=%d\n", k, i, b, d)
			fmt.Printf("[%d, %d]\n", oL, oR)
			if oR <= oL {
				return 0
			}
			if b == 1 {
				oL = w.offsets[d] + oL - w.ranks[d][oL]
				oR = w.offsets[d] + oR - w.ranks[d][oR-1]
			} else {
				oL = w.ranks[d][oL] - 1
				oR = w.ranks[d][oR-1]
			}
		}
	}

	return oR - oL
}

type Search struct {
	d  int
	oL int
	oR int
}

func (s Search) Size() int {
	return s.oR - s.oL
}

// Get top-t frequent K-mers in s[i:j]
func (w WaveletV2) Top(i int, j int, t int, K int) int {
	if i < 0 || j >= w.N() || i > j {
		panic("invalid search interval")
	}
	if K > w.K() {
		panic("k cannot be grater than Wavelet.K")
	}

	h := New()

	// subregion [oL:oR] have d bits match.
	h.HeapPush(Search{oL: i, oR: j, d: 0})

	for h.Len() > 0 {
		s := h.HeapPop()
		fmt.Printf("searching [%d,%d] in %d\n", s.oL, s.oR, s.d)

		if s.d == K*8-1 {
			fmt.Println("found!")
			return s.Size()
		}

		// to left
		s0 := Search{
			oL: w.ranks[s.d][s.oL] - 1,
			oR: w.ranks[s.d][s.oR] - 1,
			d:  s.d + 1,
		}
		if s0.Size() > 0 {
			h.HeapPush(s0)
		}

		// to right
		s1 := Search{
			oL: w.offsets[s.d] + s.oL - w.ranks[s.d][s.oL],
			oR: w.offsets[s.d] + s.oR - w.ranks[s.d][s.oR],
			d:  s.d + 1,
		}
		if s1.Size() > 0 {
			h.HeapPush(s1)
		}
	}
	return 0
}

func (w WaveletV2) Intersect(i int, query []byte) int {
	return 0
}
