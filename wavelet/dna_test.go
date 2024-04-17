package wavelet

import (
	"bytes"
	"testing"
)

func TestDNAEncode(t *testing.T) {
	s := []byte("ACGATTCG$ATTCG GG#TTA")
	u := FromCompactDNA(ToCompactDNA(s))
	if !bytes.Equal(s, u) {
		t.Error()
	}

	// Lower character will be encoded as Capital
	s2 := []byte("ACgattcg$ATTCG GG#TTA")
	u2 := FromCompactDNA(ToCompactDNA(s2))
	if !bytes.Equal(u2, s) {
		t.Error()
	}
}
