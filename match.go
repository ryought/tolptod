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

func BuildIndexes(records []fasta.Record) []suffixarray.Index {
	indexes := make([]suffixarray.Index, len(records))
	for i, record := range records {
		// create suffix array
		indexes[i] = *suffixarray.New(record.Seq)
	}
	return indexes
}

func FindMatch(x suffixarray.Index, xa int, xb int, y []byte, scale int, k int, freqLow int, freqUp int, revcomp bool) []Point {
	nx := (xb-xa)/scale + 1
	ny := len(y)/scale + 1
	m := NewMatrix(nx, ny)

	for j := 0; j < len(y)-k; j++ {
		kmer := y[j : j+k]
		if revcomp {
			kmer = fasta.RevComp(kmer)
		}

		// start := time.Now()
		_, offsets := x.LookupWithin(kmer, xa, xb, freqUp+1)
		// elapsed := time.Now().Sub(start)
		// fmt.Println("j", j, elapsed.Milliseconds())

		// this k-mer satisfies the freq restriction.
		if freqLow <= len(offsets) && len(offsets) <= freqUp {
			// fill the cells
			for _, i := range offsets {
				if xa <= i && i < xb {
					m.Set((i-xa)/scale, j/scale, true)
				}
			}
		}
	}
	points := m.Drain()

	return points
}
