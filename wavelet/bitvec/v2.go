package bitvec

import (
	"fmt"
	"math/bits"
)

const (
	L    = 256
	S    = 8
	U    = L / BYTE
	BYTE = 8
)

// Use ~2n bits for ~n bits supporting rank operation
type BitVecV2 struct {
	size   int    // n
	chunks []byte // 1bit  x n           = n
	rankL  []int  // 64bit x n/L (L=256) = n/4
	rankS  []byte // 8bit  x n/S (S=8)   = n
}

func NewV2(size int) BitVecV2 {
	return BitVecV2{
		size:   size,
		chunks: make([]byte, size/BYTE+1),
		rankL:  make([]int, size/L+1),
		rankS:  make([]byte, size/S+1),
	}
}

// Build BitVecV2 with size and accessor
func Build(size int, f func(i int) byte) BitVecV2 {
	bv := NewV2(size)
	for i := 0; i < size; i++ {
		bv.Set(i, f(i))
	}
	bv.UpdateRank()
	return bv
}

func (bv BitVecV2) Size() int {
	return bv.size
}

func (bv BitVecV2) Debug() {
	fmt.Printf("********\n")
	fmt.Printf("size=%d\n", bv.size)
	for i := range bv.chunks {
		if i%U == 0 {
			fmt.Printf("%d\t%08b\t%d\t%d\n", i, bv.chunks[i], bv.rankL[i/U], bv.rankS[i])
		} else {
			fmt.Printf("%d\t%08b\t-\t%d\n", i, bv.chunks[i], bv.rankS[i])
		}
	}
	fmt.Printf("********\n")
}

// Get i/8 and i%8
func address(i int) (int, int) {
	return i >> 3, i & 0b111
}

func (bv BitVecV2) Set(i int, x byte) {
	index, offset := address(i)
	chunk := bv.chunks[index]
	if x == 0 {
		chunk = chunk & ^(1 << offset)
	} else {
		chunk = chunk | (1 << offset)
	}
	bv.chunks[index] = chunk
}

func (bv BitVecV2) Get(i int) byte {
	index, offset := address(i)
	chunk := bv.chunks[index]
	return (chunk >> offset) & 1
}

func (bv BitVecV2) UpdateRank() {
	rank := 0
	rankL := 0
	for k, chunk := range bv.chunks {
		if k%U == 0 {
			bv.rankL[k/U] = rank
			rankL = rank
		}
		bv.rankS[k] = byte(rank - rankL)
		rank += BYTE - bits.OnesCount8(chunk)
	}
}

// Count zero in first [0:offset) bits
//
// 0b_0010_0101 b
// ------|-<--> [0:offset)
// offset=4
//
// 0b_1111_0000 mask
func countZeros(b byte, offset int) int {
	mask := byte(^((1 << offset) - 1))
	return BYTE - bits.OnesCount8(b|mask)
}

// Count zeros in [0:i)
func (bv BitVecV2) Rank(i int) int {
	index, offset := address(i)
	chunk := bv.chunks[index]
	rank := countZeros(chunk, offset)
	return bv.rankL[index/U] + int(bv.rankS[index]) + rank
}
