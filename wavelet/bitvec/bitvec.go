package bitvec

import "fmt"

type BitVec struct {
	bits []bool
	rank []int
	// index []int
}

func New(bits []bool) BitVec {
	// rank
	rank := make([]int, len(bits))
	// index := make([]int, len(bits))
	count := 0
	for i := range bits {
		if bits[i] {
			// index[count] = i
			count += 1
		}
		rank[i] = count
	}

	return BitVec{
		bits,
		rank,
		// index,
	}
}

func bit2int(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func (bv BitVec) debug() {
	for i := range bv.bits {
		fmt.Printf("bits[i=%d]=%d\trank=%d\n", i, bit2int(bv.bits[i]), bv.rank[i])
	}
}

// Get bits[i]
func (bv BitVec) Access(i int) bool {
	return bv.bits[i]
}

// Count bit b (0 or 1) in bits[0:i] = bits[0],...,bits[i]
func (bv BitVec) Rank(b bool, i int) int {
	if b {
		return bv.rank[i]
	} else {
		return i + 1 - bv.rank[i]
	}
}

// Get position of i-th bit b
// func (bv BitVec) Select(b bool, i int) int {
// 	if b {
// 		return bv.index[i]
// 	} else {
// 		return i - bv.index[i]
// 	}
// }
