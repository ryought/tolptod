package bitvec

import "fmt"

type BitVecV2 struct {
	size   int
	chunks []byte
}

// Perform div a/b then ceil.
// Equivalent to int(math.Ceil(float(a) / float(b))) but not use float.
func div(a int, b int) int {
	return (a + (b - 1)) / b
}

func NewV2(size int) BitVecV2 {
	return BitVecV2{
		size:   size,
		chunks: make([]byte, div(size, 8)),
	}
}

func (bv BitVecV2) Size() int {
	return bv.size
}

func (bv BitVecV2) Debug() {
	fmt.Printf("size=%d\n", bv.size)
	fmt.Printf("i\t01234567\n")
	for i := range bv.chunks {
		fmt.Printf("%d\t%08b\n", i, bv.chunks[i])
	}
}

func (bv BitVecV2) Set(i int, x byte) {
	index := i >> 3
	offset := i & 0b111
	chunk := bv.chunks[index]
	if x == 0 {
		chunk = chunk & ^(1 << offset)
	} else {
		chunk = chunk | (1 << offset)
	}
	bv.chunks[index] = chunk
}

func (bv BitVecV2) Get(i int) byte {
	index := i >> 3
	offset := i & 0b111
	chunk := bv.chunks[index]
	return (chunk >> offset) & 1
}
