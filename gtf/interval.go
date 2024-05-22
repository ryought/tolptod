package gtf

import (
	"github.com/rdleal/intervalst/interval"
	"github.com/ryought/tolptod/fasta"
)

type GTFTree struct {
	tree []*interval.SearchTree[Feature, int]
}

func BuildIntervalTree(seqs []fasta.Seq, features []Feature) GTFTree {
	cmpFn := func(x, y int) int { return x - y }
	t := GTFTree{
		tree: make([]*interval.SearchTree[Feature, int], len(seqs)),
	}
	m := make(map[string]int, len(seqs))
	for i, seq := range seqs {
		t.tree[i] = interval.NewSearchTree[Feature](cmpFn)
		m[seq.Id] = i
	}
	for _, feature := range features {
		i, ok := m[feature.SeqName]
		if ok {
			t.tree[i].Insert(feature.Start, feature.End, feature)
		}
	}
	return t
}

func (t *GTFTree) Find(i int, start int, end int) []Feature {
	fs, _ := t.tree[i].AllIntersections(start, end)
	return fs
}
