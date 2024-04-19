package wavelet

// reference
// Claude, F., Navarro, G. (2012). The Wavelet Matrix. In: Calderón-Benavides, L., González-Caro, C., Chávez, E., Ziviani, N. (eds) String Processing and Information Retrieval. SPIRE 2012. Lecture Notes in Computer Science, vol 7608. Springer, Berlin, Heidelberg. https://doi.org/10.1007/978-3-642-34109-0_18
// https://miti-7.hatenablog.com/entry/2018/04/28/152259
// https://scrapbox.io/koki/Wavelet_Matrix
// https://www.slideshare.net/pfi/ss-15916040

import (
	"fmt"
	"github.com/ryought/tolptod/wavelet/bitvec"
	// "time"
)

type Wavelet struct {
	bits     []bitvec.BitVecV2
	offsets  []int
	width    int
	terminal []byte
}

// i-th (0<=i<8) bit of s[k] (byte uint8).
func ix(s []byte, k int, i int) byte {
	if k < len(s) {
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
func (w Wavelet) K() int {
	return w.D() / w.W()
}

// bit width
func (w Wavelet) W() int {
	return w.width
}

// depth
func (w Wavelet) D() int {
	return len(w.offsets)
}

// length of s
func (w Wavelet) N() int {
	return w.bits[0].Size()
}

func (w Wavelet) IsTerminal(c byte) bool {
	for _, t := range w.terminal {
		if c == t {
			return true
		}
	}
	return false
}

// Constructor
func New(s []byte, K int) Wavelet {
	return NewCustom(s, K, 8, []byte{0, '$'})
}

// constructor
func NewCustom(s []byte, K int, W int, terminal []byte) Wavelet {
	if W < 1 || W > 8 {
		panic("W should be 1<=W<=8.")
	}
	D := K * W
	bits := make([]bitvec.BitVecV2, D)
	offsets := make([]int, D)

	// X0
	x := make([]int, len(s))
	tmp := make([]int, len(s))
	for o := range x {
		x[o] = o
	}

	for k := 0; k < K; k++ {
		// fmt.Printf("k=%d\n", k)
		for i := 0; i < W; i++ {
			// depth
			d := k*W + i

			// B0
			// t0 := time.Now()
			b := bitvec.Build(len(s), func(o int) byte {
				return ix(s, x[o]+k, i)
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
	}

	width := W
	return Wavelet{bits, offsets, width, terminal}
}

// Create rank array rank[0:n+1) of bytes b[0:n).
//
// rank[i] is the occurrence count of 0s in b[0:i)=b[0]...b[i-1].
//
// (the occurrence of 1s in b[0:i) is i-rank[i].)
//
// For example, createRank([0 1 0 1 0 0])=[0 1 1 2 2 3 4]
func createRank(b []byte) []int {
	n := len(b)
	rank := make([]int, n+1)
	// current occurrence of zero in b
	count := 0
	for i := range b {
		rank[i] = count
		if b[i] == 0 {
			count += 1
		}
	}
	rank[n] = count
	return rank
}

// Get s[i:i+K) subsequence for 0<=i<n.
func (w Wavelet) Access(i int, K int) []byte {
	if i < 0 || i >= w.N() {
		panic("Access: index out of range")
	}
	if K > w.K() {
		panic("Access: k cannot be grater than Wavelet.K")
	}

	// subsequence to be returned
	s := make([]byte, K)

	// index in w.bits[d].
	// B[0][i]
	o := i
	W := w.W()
	for k := 0; k < K; k++ {
		var c byte
		for i := 0; i < W; i++ {
			d := k*W + i
			b := w.bits[d].Get(o)
			c = c | (b << i)
			if b == 1 {
				// go to right (1).
				// new index =  (number of zeros (offset)) + (number of ones before o)
				//
				//                 o        o=4
				// B[d-1]  0 1 0 1 1 1 0 0  offset=4
				// rank0   1 1 2 2 2 2 3 4
				// rank1   0 1 1 2 3 4 4 4  rank1[i] = i - rank0[i]
				// B[d]    0 0 0 0 1 1 1 1
				//         <-----> <--->
				//                     o    new o=6
				//         offset=4
				//                 rank1[o]=3
				//
				o = w.offsets[d] + o - w.bits[d].Rank(o)
			} else {
				// go to left (0).
				//                 o        o=4
				// B[d-1]  0 1 0 1 0 1 1 0
				// rank0   1 1 2 2 3 3 3 4  rank0[4]=3
				// B[d]    0 0 0 0 1 1 1 1
				//         <--->
				//             o        new o=2
				o = w.bits[d].Rank(o)
			}
			// fmt.Println("newo", o)
		}
		s[k] = c
	}
	return s
}

// Get the occurrence of query in s[0:i) for 0<=i<=n.
func (w Wavelet) Rank(i int, query []byte) int {
	if len(query) > w.K() {
		panic("query cannot be longer than Wavelet.K")
	}
	if i < 0 || i > w.N() {
		panic("i should be 0<=i<=N")
	}

	// subregion [oL:oR) in B[d] of bits represents query that have the same first d bit.
	// d = 0, [L=0:R=i) is valid region.
	oL, oR := 0, i

	W := w.W()
	for k := range query {
		for i := 0; i < W; i++ {
			d := k*W + i
			b := ix(query, k, i)
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
	}

	return oR - oL
}

type Search struct {
	d  int
	oL int
	oR int
	b  []byte
}

func (s Search) Size() int {
	return s.oR - s.oL
}

func clone(b []byte) []byte {
	r := make([]byte, len(b))
	copy(r, b)
	return r
}

// Get top frequent K-mers in s[i:j)
func (w Wavelet) Top(i int, j int, K int) ([]byte, int) {
	if i < 0 || j > w.N() || i > j {
		panic("invalid search interval")
	}
	if K > w.K() {
		panic("k cannot be grater than Wavelet.K")
	}

	// initialize priority queue
	h := NewPriorityQueue()

	// subregion [oL:oR) have d bits match.
	h.HeapPush(Search{oL: i, oR: j, d: 0, b: make([]byte, 0)})

	W := w.W()
	for h.Len() > 0 {
		s := h.HeapPop()
		// fmt.Printf("searching [%d,%d) in d=%d \n", s.oL, s.oR, s.d)
		// fmt.Println("b", s.b)

		i := s.d % W
		k := s.d / W
		if i == 0 && k > 0 {
			if w.IsTerminal(s.b[k-1]) {
				// break this intersection
				continue
			}
		}
		if k == K {
			return s.b, s.Size()
		}

		// prepare a new byte to store k-mer bytes
		if i == 0 {
			s.b = append(s.b, 0)
		}

		// to left
		{
			oL := w.bits[s.d].Rank(s.oL)
			oR := w.bits[s.d].Rank(s.oR)
			if oL < oR {
				b := clone(s.b)
				b[k] = b[k] | (0 << i)
				// fmt.Printf("L i=%d k=%d %08b\n", i, k, b[k])
				// fmt.Printf("L [%d:%d)\n", oL, oR)
				h.HeapPush(Search{
					d:  s.d + 1,
					oL: oL,
					oR: oR,
					b:  b,
				})
			}
		}

		// to right
		{
			oL := w.offsets[s.d] + s.oL - w.bits[s.d].Rank(s.oL)
			oR := w.offsets[s.d] + s.oR - w.bits[s.d].Rank(s.oR)
			if oL < oR {
				b := clone(s.b)
				b[k] = b[k] | (1 << i)
				// fmt.Printf("R i=%d k=%d %08b\n", i, k, b[k])
				// fmt.Printf("R [%d:%d)\n", oL, oR)
				s1 := Search{
					d:  s.d + 1,
					oL: oL,
					oR: oR,
					b:  b,
				}
				h.HeapPush(s1)
			}
		}
	}

	// not found
	return []byte{}, 0
}

// Find common K-mer in S[aL:aR) and S[bL:bR).
func (w Wavelet) Intersect(aL, aR, bL, bR int, K int) (int, int) {
	if aL < 0 || aR > w.N() || aL > aR {
		panic("invalid search interval [aL:aR)")
	}
	if bL < 0 || bR > w.N() || bL > bR {
		panic("invalid search interval [bL:bR)")
	}
	if K > w.K() {
		panic("k cannot be grater than Wavelet.K")
	}

	q := NewQueue()
	d := 0
	var c byte
	q.HeapPush(Intersection{aL, aR, bL, bR, d, c})

	W := w.W()
	for q.Len() > 0 {
		is := q.HeapPop()
		// fmt.Printf("poped [%d,%d) [%d,%d) d=%d\n", is.aL, is.aR, is.bL, is.bR, is.d)

		i := is.d % W
		k := is.d / W
		if i == 0 && k > 0 {
			// fmt.Printf("current char %c\n", c)
			if w.IsTerminal(is.c) {
				// break this intersection
				continue
			} else {
				is.c = 0
			}
		}
		if k == K {
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
			c := is.c | (0 << i)
			// fmt.Printf("L [%d,%d) [%d,%d) d=%d\n", aL, aR, bL, bR, d)
			if aL < aR && bL < bR {
				q.HeapPush(Intersection{aL, aR, bL, bR, d, c})
			}
		}

		// to right
		{
			aL := w.offsets[is.d] + is.aL - w.bits[is.d].Rank(is.aL)
			aR := w.offsets[is.d] + is.aR - w.bits[is.d].Rank(is.aR)
			bL := w.offsets[is.d] + is.bL - w.bits[is.d].Rank(is.bL)
			bR := w.offsets[is.d] + is.bR - w.bits[is.d].Rank(is.bR)
			d := is.d + 1
			c := is.c | (1 << i)
			// fmt.Printf("L [%d,%d) [%d,%d) d=%d\n", aL, aR, bL, bR, d)
			if aL < aR && bL < bR {
				q.HeapPush(Intersection{aL, aR, bL, bR, d, c})
			}
		}
	}

	return 0, 0
}
