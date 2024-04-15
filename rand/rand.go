package rand

import "math/rand"

func RandomDNA(n int) (ret []byte) {
	DNA := []byte{'A', 'C', 'G', 'T'}
	r := rand.New(rand.NewSource(0))
	ret = make([]byte, n)
	for i := range ret {
		ret[i] = DNA[r.Intn(4)]
	}
	return
}

func RandomByte(n int) (ret []byte) {
	r := rand.New(rand.NewSource(0))
	ret = make([]byte, n)
	for i := range ret {
		ret[i] = byte(r.Intn(4))
	}
	return
}