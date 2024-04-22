package suffixarray

import "fmt"

// create LCP array
func (x *Index) LCP() []int64 {
	sa := x.SA()
	s := x.data

	// sa[r] = i means that suffix S[i:] is r-th in suffix array.
	// rank[i] = r
	rank := make([]int64, len(s))
	for r, i := range sa {
		rank[i] = int64(r)
	}

	lcp := make([]int64, len(s))
	// calculate for LCP[i] in decreasing order of suffix length.
	var l int64
	n := int64(len(s))
	for i_, r := range rank {
		i := int64(i_)
		// previous element in suffix array table
		if r == 0 {
			lcp[r] = -1
		} else {
			i0 := sa[r-1]
			for i0+l < n && i+l < n && s[i0+l] == s[i+l] {
				l += 1
			}
			lcp[r] = l
			l = max(l-1, 0)
		}
	}

	return lcp
}

func (x *Index) KmerMatches(LCP []int64, K int) ([]int64, int64) {
	sa := x.SA()
	var kmerId int64
	inRun := false
	matches := make([]int64, len(sa))

	for i := 0; i < len(sa); i++ {
		if LCP[i] >= int64(K) {
			if !inRun {
				kmerId += 1
				inRun = true
				matches[sa[i-1]] = kmerId
			}
			matches[sa[i]] = kmerId
		} else {
			inRun = false
		}
	}
	return matches, kmerId
}
