package main

import (
	"github.com/ryought/tolptod/fasta"
	"github.com/ryought/tolptod/suffixarray"
)

type Config struct {
	k       int
	bin     int
	freqLow int
	freqUp  int
	xL      int
	xR      int
	yL      int
	yR      int
}

type Index struct {
	N        int
	Forward  suffixarray.Index
	Backward suffixarray.Index
}

func NewIndex(s []byte) Index {
	return Index{
		N:        len(s),
		Forward:  *suffixarray.New(s),
		Backward: *suffixarray.New(fasta.RevComp(s)),
	}
}

type IndexV2 struct {
	xindex []Index
	yindex []Index
}

type Cache struct {
	matF   Matrix
	matB   Matrix
	config Config
}

func NewIndexV2(xs [][]byte, ys [][]byte) IndexV2 {
	xindex := make([]Index, len(xs))
	for i, x := range xs {
		xindex[i] = NewIndex(x)
	}
	yindex := make([]Index, len(ys))
	for i, y := range ys {
		yindex[i] = NewIndex(y)
	}
	return IndexV2{xindex, yindex}
}

func NewIndexV2FromRecords(xs []fasta.Record, ys []fasta.Record) IndexV2 {
	xindex := make([]Index, len(xs))
	for i, x := range xs {
		xindex[i] = NewIndex(x.Seq)
	}
	yindex := make([]Index, len(ys))
	for i, y := range ys {
		yindex[i] = NewIndex(y.Seq)
	}
	return IndexV2{xindex, yindex}
}

func ComputeMatrix(xindex, yindex Index, config Config) (Matrix, Matrix) {
	X := xindex.N
	Y := yindex.N
	W := config.bin
	K := config.k
	nx := CeilDiv(config.xR-config.xL, W)
	ny := CeilDiv(config.yR-config.yL, W)
	matF := NewMatrix(nx, ny)
	matB := NewMatrix(nx, ny)

	for y := config.yL; y < min(config.yR, Y-K+1); y++ {
		kmer := yindex.Forward.Bytes()[y : y+K]
		xF := xindex.Forward.LookupAll(kmer)
		xB := xindex.Backward.LookupAll(kmer)

		// count for match in the region
		n := 0
		for i := 0; i < xF.Len(); i++ {
			x := xF.Get(i)
			if config.xL <= x && x < config.xR {
				n += 1
			}
		}
		for i := 0; i < xB.Len(); i++ {
			x := X - xB.Get(i) - K
			if config.xL <= x && x < config.xR {
				n += 1
			}
		}

		// fill the table
		if config.freqLow <= n && (config.freqUp == -1 || n <= config.freqUp) {
			for i := 0; i < xF.Len(); i++ {
				x := xF.Get(i)
				if config.xL <= x && x < config.xR {
					matF.Set((x-config.xL)/W, (y-config.yL)/W, true)
				}
			}
			for i := 0; i < xB.Len(); i++ {
				x := X - xB.Get(i) - K
				if config.xL <= x && x < config.xR {
					matB.Set((x-config.xL)/W, (y-config.yL)/W, true)
				}
			}
		}
	}
	return matF, matB
}

func NewCache(xindex, yindex Index, config Config) Cache {
	c := Config{
		k:       config.k,
		bin:     config.bin,
		freqLow: config.freqLow,
		freqUp:  config.freqUp,
		xL:      0,
		xR:      xindex.N,
		yL:      0,
		yR:      yindex.N,
	}
	matF, matB := ComputeMatrix(xindex, yindex, c)
	return Cache{
		matF:   matF,
		matB:   matB,
		config: c,
	}
}

func (c Cache) ComputeMatrix(config Config) (Matrix, Matrix) {
	W0 := c.config.bin
	W := config.bin
	R := W / W0
	if W%W0 != 0 {
		panic("query bin size should be a multiple of cache bin size")
	}
	if W < W0 {
		panic("query bin size should be larger than the cache bin size")
	}
	if config.xL%W != 0 || config.xR%W != 0 || config.yR%W != 0 || config.yL%W != 0 {
		panic("region left/right should be a multiple of bin size")
	}
	nx := CeilDiv(config.xR-config.xL, W)
	ny := CeilDiv(config.yR-config.yL, W)
	matF := NewMatrix(nx, ny)
	matB := NewMatrix(nx, ny)
	xL0 := config.xL / W0
	xR0 := config.xR / W0
	yL0 := config.yL / W0
	yR0 := config.yR / W0
	for x := xL0; x < xR0; x++ {
		for y := yL0; y < yR0; y++ {
			if c.matF.Get(x, y) == true {
				matF.Set((x-xL0)/R, (y-yL0)/R, true)
			}
			if c.matB.Get(x, y) == true {
				matB.Set((x-xL0)/R, (y-yL0)/R, true)
			}
		}
	}
	return matF, matB
}
