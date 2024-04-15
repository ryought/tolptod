package suffixarray

import (
	"testing"
	// "github.com/ryought/tolptod/rand"
)

func TestLCP(t *testing.T) {
	s := []byte("ATCGGATCGATTCG$")
	index := New(s)
	lcp := index.LCP()
	sa := index.SA()
	t.Log("lcp", lcp)
	t.Errorf("hoge")

	for i := range sa {
		t.Log(string(s[sa[i]:]), sa[i], lcp[i])
	}
}
