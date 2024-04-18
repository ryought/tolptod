package wavelet

import (
	"bytes"
	"testing"
)

func TestDNAWavelet(t *testing.T) {
	K := 10
	//-----------0         10        20        30
	S := []byte("ATCTAGCTAGCTTAGCTAG$ATCGAGCTAGCTGGATC$")
	//           *                   * ATC
	//              *   *    *         TAGCT
	w := NewDNAWavelet(S, K)
	if !bytes.Equal(w.Access(0, K), S[0:K]) {
		t.Error()
	}
	if !bytes.Equal(w.Access(15, K), S[15:15+K]) {
		t.Error()
	}
	if w.Rank(30, []byte("ATC")) != 2 {
		t.Error()
	}

	kmer, count := w.Top(0, 20, 5)
	t.Log(string(kmer), count)
	if !bytes.Equal(kmer, []byte("TAGCT")) {
		t.Error()
	}
	if count != 3 {
		t.Error()
	}

	kmer, count = w.Top(0, len(S), 1)
	t.Log(string(kmer), count)
	if kmer[0] != 'T' || count != 10 {
		t.Error()
	}

	a, b := w.Intersect(0, 10, 10, 20, 5)
	t.Log(a, b)
	if a != 1 || b != 1 {
		t.Error()
	}
}
