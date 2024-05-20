package main

import (
	"testing"
)

func TestIndex(t *testing.T) {
	s := []byte("ATCGGATCGT")
	// n := len(s)
	x := NewIndex(s)
	y := NewIndex(s)

	// bin=1
	c1 := Config{
		k:       4,
		bin:     1,
		freqLow: 0,
		freqUp:  -1,
		xL:      0,
		xR:      8,
		yL:      0,
		yR:      8,
	}
	mf1, mb1 := ComputeMatrix(x, y, c1)
	t.Log("f")
	mf1.Print()
	t.Log("b")
	mb1.Print()
}

func TestHoge(t *testing.T) {
	s := []byte("ATCGGATCGT")
	n := len(s)
	x := NewIndex(s)
	y := NewIndex(s)

	// bin=1
	c1 := Config{
		k:       4,
		bin:     1,
		freqLow: 0,
		freqUp:  1000,
		xL:      0,
		xR:      n,
		yL:      0,
		yR:      n,
	}
	mf1, mb1 := ComputeMatrix(x, y, c1)
	t.Log("f")
	mf1.Print()
	t.Log("b")
	mb1.Print()

	// bin=2
	c2 := Config{
		k:       4,
		bin:     2,
		freqLow: 0,
		freqUp:  1000,
		xL:      0,
		xR:      n,
		yL:      0,
		yR:      n,
	}
	mf2, mb2 := ComputeMatrix(x, y, c2)
	t.Log("f2")
	mf2.Print()
	t.Log("b2")
	mb2.Print()

	// bin=2 from cache1
	cache1 := Cache{matF: mf1, matB: mb1, config: c1}
	mf2s, mb2s := cache1.ComputeMatrix(c2)
	t.Log("f2s")
	mf2s.Print()
	t.Log("b2s")
	mb2s.Print()
	if !mf2.Equal(mf2s) {
		t.Error("error: mf2 != mf2s")
	}
	if !mb2.Equal(mb2s) {
		t.Error("error: mb2 != mb2s")
	}

	// bin=4
	c4 := Config{
		k:       4,
		bin:     4,
		freqLow: 0,
		freqUp:  -1,
		xL:      4,
		xR:      8,
		yL:      4,
		yR:      8,
	}
	t.Log("Compute...")
	mf4, mb4 := ComputeMatrix(x, y, c4)
	t.Log("f4")
	mf4.Print()
	t.Log("b4")
	mb4.Print()

	// bin=4 from cache
	// cache2 := NewCache(x, y, c2)
	// mf4s, mb4s := cache2.ComputeMatrix(c4)
	// t.Log("f4s")
	// mf4s.Print()
	// t.Log("b4s")
	// mb4s.Print()
}
