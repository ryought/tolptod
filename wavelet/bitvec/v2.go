package bitvec

import (
	"fmt"
	"math/bits"
)

const (
	L               = 256 // rankS is []byte, so maximum number is 256
	S               = 64  // if use uint64 to store bits, S=64
	U               = L / S
	USE_BATCH_BUILD = true
)

// BitVec2
// Use ~2n bits for ~n bits supporting rank operation
type BitVecV2 struct {
	size   int      // n
	chunks []uint64 // 1bit  x n          = n  (as 64bit x n/64)
	rankL  []int    // 64bit x n/L (=256) = n/4
	rankS  []byte   // 8bit  x n/S (=64)  = n/8
}

func NewV2(size int) BitVecV2 {
	return BitVecV2{
		size:   size,
		chunks: make([]uint64, size/S+1),
		rankL:  make([]int, size/L+1),
		rankS:  make([]byte, size/S+1),
	}
}

// Build BitVecV2 with size and accessor
func Build(size int, f func(i int) byte) BitVecV2 {
	bv := NewV2(size)

	if USE_BATCH_BUILD {
		// fill chunks for each S=64 bits
		for index := 0; index < size/S; index++ {
			offset := index * S
			bv.SetChunk(index, createChunk(apply64(f, offset)))
		}
		// fill the last chunk for each 1 bits..
		for i := (size / S) * S; i < size; i++ {
			bv.Set(i, f(i))
		}
	} else {
		for i := 0; i < size; i++ {
			bv.Set(i, f(i))
		}
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
			fmt.Printf("%d\t%064b\t%d\t%d\n", i, bv.chunks[i], bv.rankL[i/U], bv.rankS[i])
		} else {
			fmt.Printf("%d\t%064b\t-\t%d\n", i, bv.chunks[i], bv.rankS[i])
		}
	}
	fmt.Printf("********\n")
}

// Get i/64 and i%64
func address(i int) (int, int) {
	return i >> 6, i & 0b111111
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

// createChunk
// bits = [0,1,0,1,0,0,...]
// bits[i] is offset i
func createChunk(bits [S]byte) uint64 {
	var chunk uint64
	for i := 0; i < S; i++ {
		chunk |= uint64(bits[i]) << i
	}
	return chunk
}

func apply64(f func(i int) byte, offset int) [64]byte {
	var ret [64]byte
	for i := 0; i < 64; i++ {
		ret[i] = f(offset + i)
	}
	return ret
}

func (bv BitVecV2) SetChunk(index int, chunk uint64) {
	bv.chunks[index] = chunk
}

func (bv BitVecV2) Get(i int) byte {
	index, offset := address(i)
	chunk := bv.chunks[index]
	return byte((chunk >> offset) & 1)
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
		rank += S - bits.OnesCount64(chunk)
	}
}

// Count zero in first [0:offset) bits
//
// 0b_0010_0101 b
// ------|-<--> [0:offset)
// offset=4
//
// 0b_1111_0000 mask
func countZeros(b uint64, offset int) int {
	var mask uint64 = ^((1 << offset) - 1)
	return S - bits.OnesCount64(b|mask)
}

// Count zeros in [0:i)
func (bv BitVecV2) Rank(i int) int {
	index, offset := address(i)
	chunk := bv.chunks[index]
	rank := countZeros(chunk, offset)
	return bv.rankL[index/U] + int(bv.rankS[index]) + rank
}
