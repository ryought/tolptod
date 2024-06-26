package rand

import "math/rand"

func RandomDNA(n int) []byte {
	return RandomDNAWithSeed(n, 0)
}

func RandomDNAWithSeed(n int, seed int64) (ret []byte) {
	DNA := []byte{'A', 'C', 'G', 'T'}
	r := rand.New(rand.NewSource(seed))
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

// generate x[0,n) s.t. x[i] is random number in [0,m)
func RandomUint64(n int, m int64) (ret []int64) {
	r := rand.New(rand.NewSource(0))
	ret = make([]int64, n)
	for i := range ret {
		ret[i] = r.Int63n(m)
	}
	return
}
