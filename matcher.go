package main

import (
	"fmt"
	"github.com/ryought/tolptod/fasta"
	"github.com/ryought/tolptod/suffixarray"
	"github.com/ryought/tolptod/wavelet"
)

type MatcherConfig struct {
	k          int
	freqLow    int
	freqUp     int
	useWavelet bool
}

type MatcherSA struct {
	SAForward  suffixarray.Index
	SABackward suffixarray.Index
	S          []byte
	T          []byte
}

func NewMatcherSA(S []byte, T []byte) MatcherSA {
	m := MatcherSA{
		SAForward:  *suffixarray.New(S),
		SABackward: *suffixarray.New(fasta.RevComp(S)),
		S:          S,
		T:          T,
	}
	return m
}

func FloorDiv(x int, y int) int {
	// if x > 0 {
	// 	return 1 + (x-1)/y
	// } else {
	// 	return x / y
	// }
	return 1 + (x-1)/y
}

// between S[xL:xR] and T[yL:yR]
func (m MatcherSA) Match(W int, xL, xR, yL, yR int, K int, freqLow int, freqUp int) (Matrix, Matrix) {
	X := len(m.S)
	Y := len(m.T)
	if W <= 0 {
		panic("W should be >0")
	}
	if K <= 0 {
		panic("K should be >0")
	}
	if xL < 0 || xR > X {
		panic("[xL:xR] out of range")
	}
	if yL < 0 || yR > Y {
		panic("[yL:yR] out of range")
	}
	nx := FloorDiv(xR-xL, W)
	ny := FloorDiv(yR-yL, W)
	MF := NewMatrix(nx, ny)
	MB := NewMatrix(nx, ny)

	for y := yL; y < min(yR, Y-K+1); y++ {
		kmer := m.T[y : y+K]
		_, xF := m.SAForward.LookupWithin(kmer, xL, xR, freqUp+1)
		_, xB := m.SABackward.LookupWithin(kmer, X-xR-K, X-xL-K, freqUp+1)
		n := len(xF) + len(xB)
		fmt.Println(y, string(kmer), "F", xF, "B", xB)
		if freqLow <= n && n <= freqUp {
			// forward
			for _, x := range xF {
				MF.Set((x-xL)/W, (y-yL)/W, true)
			}
			// backward
			for _, x := range xB {
				MB.Set((X-x-K-xL)/W, (y-yL)/W, true)
			}
		}
	}
	return MF, MB
}

type MatcherWT struct {
	K        int
	Forward  wavelet.IntWavelet
	Backward wavelet.IntWavelet
}

func NewMatcherWT(S []byte, T []byte, K int) MatcherWT {
	// S#T$ and S#rev(T)$
	index := suffixarray.New(fasta.Join(S, T))
	LCP := index.LCP()
	kmers, _ := index.KmerMatches(LCP, K)
	w := wavelet.NewIntWavelet(kmers, 15)

	m := MatcherWT{
		K:       K,
		Forward: w,
	}
	return m
}
