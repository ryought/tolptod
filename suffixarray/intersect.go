package suffixarray

func (x *Index) Intersect(aL, aR, bL, bR int, K int) (int, int) {
	for a := aL; a < min(aR, len(x.data)-K); a++ {
		kmer := x.data[a : a+K]
		na, nb := 0, 0

		matches := x.LookupAll(kmer)
		count := matches.len()

		for i := 0; i < count; i++ {
			var pos int
			if matches.int32 != nil {
				pos = int(matches.int32[i])
			} else {
				pos = int(matches.int64[i])
			}
			if aL <= pos && pos < aR {
				na += 1
			}
			if bL <= pos && pos < bR {
				nb += 1
			}
		}

		if na > 0 && nb > 0 {
			return na, nb
		}
	}

	// not found
	return 0, 0
}
