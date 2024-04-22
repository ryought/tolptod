package suffixarray

import (
	"github.com/ryought/tolptod/fasta"
	"testing"
	// "github.com/ryought/tolptod/rand"
)

func TestLCP(t *testing.T) {
	s := []byte("mississippi")
	index := New(s)
	lcp := index.LCP()
	sa := index.SA()
	t.Log("lcp", lcp)
	t.Errorf("hoge")
	LCP := index.LCP()
	m, c := index.KmerMatches(LCP, 3)

	for i := range sa {
		t.Logf("%-20s\t%d\t%d", string(s[sa[i]:]), sa[i], lcp[i])
	}

	for i := range s {
		t.Logf("%-20s\t%d", string(s[i:]), m[i])
	}

	t.Log(c)
}

func TestLCPLarge(t *testing.T) {
	t.Log("parsing")
	records, _ := fasta.ParseFile("../chr1.fa")
	t.Log("indexing")
	index := New(records[0].Seq)
	t.Log("matching")
	LCP := index.LCP()
	_, c := index.KmerMatches(LCP, 20)
	t.Log(c)
}
