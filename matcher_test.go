package main

import (
	"github.com/ryought/tolptod/fasta"
	"testing"
)

func TestBits(t *testing.T) {
	tests := []struct {
		x    int64
		bits int
	}{
		{0b000, 0},
		{0b001, 1},
		{0b010, 2},
		{0b011, 2},
		{0b011, 2},
		{0b100, 3},
		{0b101, 3},
		{0b110, 3},
		{0b111, 3},
	}
	for _, test := range tests {
		t.Log(test.x, Bits(test.x), test.bits)
		if Bits(test.x) != test.bits {
			t.Error()
		}
	}
}

func TestMatcherSA(t *testing.T) {
	S := []byte("ATGGATCGG")
	N := len(S)
	// m := NewMatcherSA(S, fasta.RevComp(S))
	m := NewMatcherSA(S, S)
	K := 4

	// True data
	FTrue := NewMatrix(N, N)
	BTrue := NewMatrix(N, N)
	FTrue.Set(0, 0, true)
	FTrue.Set(1, 1, true)
	FTrue.Set(2, 2, true)
	FTrue.Set(3, 3, true)
	FTrue.Set(4, 4, true)
	FTrue.Set(5, 5, true)
	BTrue.Set(3, 3, true)

	// W=1
	F, B := m.Match(1, 0, N, 0, N, K, 0, 100)
	t.Logf("S =%s", S)
	t.Logf("S'=%s", fasta.RevComp(S))
	// t.Errorf("hoge", hoge)
	t.Log("F")
	F.Print()
	t.Log("B")
	B.Print()
	t.Log("hoge")

	t.Log(F.X, FTrue.X)
	if !F.Equal(FTrue) {
		t.Error("F not true")
	}
	if !B.Equal(BTrue) {
		t.Error("B not true")
	}

	// W=2
	F, B = m.Match(2, 0, N, 0, N, K, 0, 100)
	t.Log("F")
	F.Print()
	t.Log("B")
	B.Print()
}

func TestMatcherWT(t *testing.T) {
	S := []byte("ATGGATCGG")
	N := len(S)
	K := 4
	m := NewMatcherWT(S, S, K)

	// W=1
	F, B := m.Match(1, 0, N, 0, N)
	t.Logf("S =%s", S)
	t.Logf("S'=%s", fasta.RevComp(S))
	t.Log("F")
	F.Print()
	t.Log("B")
	B.Print()
	t.Log("hoge")
}
