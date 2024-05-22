package gtf

import (
	"github.com/ryought/tolptod/fasta"
	"testing"
)

func TestGTFTree(t *testing.T) {
	seqs := []fasta.Seq{
		{Id: "chr1", Len: 100},
		{Id: "chr2", Len: 200},
	}
	features := []Feature{
		{"chr1", "a", "exon", 0, 50, "+", ""},
		{"chr1", "b", "exon", 50, 100, "+", ""},
		{"chr1", "c", "exon", 0, 1, "+", ""},
		{"chr2", "d", "exon", 20, 50, "+", ""},
		{"chr2", "e", "exon", 70, 90, "+", ""},
	}
	it := BuildIntervalTree(seqs, features)
	res := it.Find(1, 0, 100)
	t.Log("res", res)
	if len(it.Find(0, 0, 100)) != 3 {
		t.Error()
	}
	if len(it.Find(0, 200, 200)) != 0 {
		t.Error()
	}
	if len(it.Find(1, 0, 100)) != 2 {
		t.Error()
	}
}
