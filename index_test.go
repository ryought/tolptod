package main

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/ryought/tolptod/rand"
)

var ctx context.Context = context.Background()

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
	mf1, mb1 := ComputeMatrix(ctx, x, y, c1)
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
	mf1, mb1 := ComputeMatrix(ctx, x, y, c1)
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
	mf2, mb2 := ComputeMatrix(ctx, x, y, c2)
	t.Log("f2")
	mf2.Print()
	t.Log("b2")
	mb2.Print()

	// bin=2 from cache1
	cache1 := Cache{matF: mf1, matB: mb1, config: c1}
	mf2s, mb2s := cache1.ComputeMatrix(ctx, c2)
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
	mf4, mb4 := ComputeMatrix(ctx, x, y, c4)
	t.Log("f4")
	mf4.Print()
	t.Log("b4")
	mb4.Print()

	// bin=4 from cache
	cache2 := NewCache(ctx, x, y, c2, nil)
	mf4s, mb4s := cache2.ComputeMatrix(ctx, c4)
	t.Log("f4s")
	mf4s.Print()
	t.Log("b4s")
	mb4s.Print()
	if !mf4.Equal(mf4s) {
		t.Error("error: mf4 != mf4s")
	}
	if !mb4.Equal(mb4s) {
		t.Error("error: mb4 != mb4s")
	}
}

func TestIndexPerformance(t *testing.T) {
	n := 1_000_000
	s := rand.RandomDNA(n)
	x := NewIndex(s)
	y := NewIndex(s)

	// bin=1
	ks := []int{5, 10, 15, 20}
	for _, k := range ks {
		start := time.Now()
		ComputeMatrix(ctx, x, y, Config{
			k:       k,
			bin:     1_000,
			freqLow: 0,
			freqUp:  -1,
			xL:      0,
			xR:      n,
			yL:      0,
			yR:      n,
		})
		t.Logf("n=%d\tk=%d\t%dus", n, k, time.Since(start).Microseconds())
	}
}

func TestListFn(t *testing.T) {
	// Unique
	xs := []int{0, 1, 5, 2, 1, 3, 2}
	t.Log(xs)
	ys := Unique(xs)
	expected := []int{0, 1, 2, 3, 5}
	t.Log(xs, ys)
	if !slices.Equal(ys, expected) {
		t.Error()
	}

	// Map
	{
		is := []int{0, 1, 5, 2, 1, 3, 2}
		js := Map(is, 10, 2, 0, 10, false)
		t.Log(is, js)
	}
	{
		is := []int{0, 1, 5, 2, 1, 3, 2}
		//         {      *  *     *  *}
		// x-L   = {-2,-1,3, 0, -1,1, 0}
		js := Map(is, 10, 2, 2, 8, false)
		t.Log(is, js)
	}
}
