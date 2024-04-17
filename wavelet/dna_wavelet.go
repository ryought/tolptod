package wavelet

type DNAWavelet struct {
	w Wavelet
}

func NewDNAWavelet(s []byte, K int) DNAWavelet {
	t := ToCompactDNA(s)
	W := 3
	return DNAWavelet{w: NewCustom(t, K, W)}
}

func (w DNAWavelet) Access(i int, K int) []byte {
	return FromCompactDNA(w.w.Access(i, K))
}

func (w DNAWavelet) Rank(i int, query []byte) int {
	return w.w.Rank(i, ToCompactDNA(query))
}

func (w DNAWavelet) Top(i int, j int, K int) ([]byte, int) {
	b, c := w.w.Top(i, j, K)
	return FromCompactDNA(b), c
}

func (w DNAWavelet) Intersect(aL, aR, bL, bR int, K int) (int, int) {
	return w.w.Intersect(aL, aR, bL, bR, K)
}
