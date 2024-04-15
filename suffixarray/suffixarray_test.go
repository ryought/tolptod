package suffixarray

import (
	"github.com/ryought/tolptod/fasta"
	"testing"
)

func TestSA(t *testing.T) {
	s := "ATCGGATCG"
	index := New([]byte(s))
	sa := index.SA()
	t.Log(sa)

	for i := range sa {
		t.Log(sa[i], s[sa[i]:])
	}

	// t.Errorf("hoge")
}

func BenchmarkSALarge(b *testing.B) {
	records, _ := fasta.ParseFile("../chr1.fa")
	b.StartTimer()
	b.Log("building", len(records[0].Seq))
	index := New(records[0].Seq)
	b.StopTimer()
	b.Log(index.sa.len())
}
