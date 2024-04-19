package bitvec

import (
	"fmt"
	"math/bits"
)

const (
	L               = 256
	S               = 8
	U               = L / BYTE
	BYTE            = 8
	USE_BATCH_BUILD = true
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

	if USE_BATCH_BUILD {
		// fill chunks for each 8 bits
		for index := 0; index < size/BYTE; index++ {
			i := index * 8
			bv.SetChunk(index, createChunk([8]byte{
				f(i),
				f(i + 1),
				f(i + 2),
				f(i + 3),
				f(i + 4),
				f(i + 5),
				f(i + 6),
				f(i + 7),
			}))
		}
		// fill the last chunk for each 1 bits..
		for i := (size / BYTE) * BYTE; i < size; i++ {
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

// createChunk
// bits = [0,1,0,1,0,0]
// bits[i] is offset i
func createChunk(bits [8]byte) byte {
	var chunk byte
	chunk |= bits[0] << 0
	chunk |= bits[1] << 1
	chunk |= bits[2] << 2
	chunk |= bits[3] << 3
	chunk |= bits[4] << 4
	chunk |= bits[5] << 5
	chunk |= bits[6] << 6
	chunk |= bits[7] << 7
	return chunk
}

func (bv BitVecV2) SetChunk(index int, chunk byte) {
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
