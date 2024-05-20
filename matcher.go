package main

import (
	"fmt"
	"github.com/ryought/tolptod/fasta"
	"github.com/ryought/tolptod/suffixarray"
	"github.com/ryought/tolptod/wavelet"
	"time"
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
	nx := CeilDiv(xR-xL, W)
	ny := CeilDiv(yR-yL, W)
	MF := NewMatrix(nx, ny)
	MB := NewMatrix(nx, ny)

	for y := yL; y < min(yR, Y-K+1); y++ {
		kmer := m.T[y : y+K]
		_, xF := m.SAForward.LookupWithin(kmer, xL, xR, freqUp+1)
		_, xB := m.SABackward.LookupWithin(kmer, X-xR-K, X-xL-K, freqUp+1)
		n := len(xF) + len(xB)
		// fmt.Println(y, string(kmer), "F", xF, "B", xB)
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
	// Forward: S$T
	indexF := suffixarray.New(fasta.Join(S, T))
	LCPF := indexF.LCP()
	kmersF, maxKmerF := indexF.KmerMatches(LCPF, K)
	fmt.Println("maxKmerF", maxKmerF)
	wF := wavelet.NewIntWavelet(kmersF)

	// Backward: S$rev(T)
	indexB := suffixarray.New(fasta.Join(S, fasta.RevComp(T)))
	LCPB := indexB.LCP()
	kmersB, maxKmerB := indexB.KmerMatches(LCPB, K)
	fmt.Println("maxKmerB", maxKmerB)
	wB := wavelet.NewIntWavelet(kmersB)

	m := MatcherWT{
		K:        K,
		Forward:  wF,
		Backward: wB,
	}
	return m
}

// between S[xL:xR] and T[yL:yR]
func (m MatcherWT) Match(W int, xL, xR, yL, yR int) (Matrix, Matrix) {
	nx := CeilDiv(xR-xL, W)
	ny := CeilDiv(yR-yL, W)
	MF := NewMatrix(nx, ny)
	MB := NewMatrix(nx, ny)
	fmt.Printf("matrix %d %d\n", nx, ny)
	var duration int64
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			aL, aR := xL+i*W, min(xL+(i+1)*W, xR)
			bL, bR := yL+j*W, min(yL+(j+1)*W, yR)

			start := time.Now()
			_, cx, cy := m.Forward.Intersect(aL, aR, bL, bR)
			t := time.Since(start).Nanoseconds()
			duration += t

			fmt.Println(i, j, t)

			if cx > 0 && cy > 0 {
				MF.Set(i, j, true)
			}
		}
	}
	fmt.Printf("average query time %d ns for W=%d\n", duration/int64(nx*ny), W)
	return MF, MB
}
