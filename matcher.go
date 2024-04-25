package main

import (
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
	T          []byte
}

func NewMatcherSA(S []byte, T []byte) MatcherSA {
	m := MatcherSA{
		SAForward:  *suffixarray.New(S),
		SABackward: *suffixarray.New(fasta.RevComp(S)),
		T:          T,
	}
	return m
}

// between seqs[x][xL*W:xR*W] and seqs[y][yL*W:yR*W]
func (m MatcherSA) Match(W int, xL, xR, yL, yR int, K int, freqLow int, freqUp int) (Matrix, Matrix) {
	nx := xR - xL
	ny := yR - yL
	MF := NewMatrix(nx, ny)
	MB := NewMatrix(nx, ny)

	iL, iR := xL*W, xR*W
	jL, jR := yL*W, yR*W

	for j := jL; j < jR; j++ {
		kmer := m.T[j : j+K]
		_, iF := m.SAForward.LookupWithin(kmer, iL, iR, freqUp+1)
		_, iB := m.SABackward.LookupWithin(kmer, N-xb, N-xa, freqUp+1)
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
	m := MatcherWT{}
	return m
}
