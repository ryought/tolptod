package main

import (
	// "time"
	// "html/template"
	"github.com/ryought/tolptod/fasta"
	"github.com/ryought/tolptod/suffixarray"
)

type Info struct {
	Xs []Seq `json:"xs"`
	Ys []Seq `json:"ys"`
}

type Seq struct {
	Id  string `json:"id"`
	Len int    `json:"len"`
}

func toSeqInfo(rs []fasta.Record) []Seq {
	is := make([]Seq, len(rs))
	for i, r := range rs {
		is[i].Id = string(r.ID)
		is[i].Len = len(r.Seq)
	}
	return is
}

func toInfo(xrs []fasta.Record, yrs []fasta.Record) Info {
	return Info{
		Xs: toSeqInfo(xrs),
		Ys: toSeqInfo(yrs),
	}
}

type Index struct {
	N        int
	Forward  suffixarray.Index
	Backward suffixarray.Index
}

func BuildIndexes(records []fasta.Record) []Index {
	indexes := make([]Index, len(records))
	for i, record := range records {
		// create suffix array
		indexes[i].N = len(record.Seq)
		indexes[i].Forward = *suffixarray.New(record.Seq)
		indexes[i].Backward = *suffixarray.New(fasta.RevComp(record.Seq))
	}
	return indexes
}

func FindMatch(x Index, xa int, xb int, y []byte, scale int, k int, freqLow int, freqUp int) ([]Point, []Point) {
	N := x.N
	nx := (xb-xa)/scale + 1
	ny := len(y)/scale + 1
	matF := NewMatrix(nx, ny)
	matB := NewMatrix(nx, ny)

	for j := 0; j < len(y)-k; j++ {
		kmer := y[j : j+k]

		_, posF := x.Forward.LookupWithin(kmer, xa, xb, freqUp+1)
		_, posB := x.Backward.LookupWithin(kmer, N-xb, N-xa, freqUp+1)
		n := len(posF) + len(posB)
		if freqLow <= n && n <= freqUp {
			// fill the cells
			// forward
			for _, i := range posF {
				if xa <= i && i < xb {
					matF.Set((i-xa)/scale, j/scale, true)
				}
			}
			// backward
			for _, i := range posB {
				if xa <= N-i && N-i < xb {
					matB.Set((N-i-xa)/scale, j/scale, true)
				}
			}
		}
	}

	return matF.Drain(), matB.Drain()
}
