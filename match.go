package main

import (
	"encoding/json"
	// "flag"
	"fmt"
	// "time"
	// "html/template"
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

func toSeqInfo(rs []Record) []Seq {
	is := make([]Seq, len(rs))
	for i, r := range rs {
		is[i].Id = string(r.ID)
		is[i].Len = len(r.Seq)
	}
	return is
}

func toInfo(xrs []Record, yrs []Record) Info {
	return Info{
		Xs: toSeqInfo(xrs),
		Ys: toSeqInfo(yrs),
	}
}

type Point struct {
	X int
	Y int
}

func (p Point) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("[%d,%d]", p.X, p.Y)
	return []byte(s), nil
}

func pointJsonTest() {
	p := Point{X: 10, Y: 20}
	points := []Point{p}
	// err := json.Unmarshal([]byte("{\"x\":10,\"y\":30}"), &p)
	s, err := json.Marshal(points)
	fmt.Println(p, string(s), err)
}

type Matrix struct {
	X int
	Y int
	m []bool
}

func NewMatrix(X int, Y int) Matrix {
	m := make([]bool, X*Y)
	return Matrix{X, Y, m}
}

func (m Matrix) Set(x int, y int, v bool) {
	if x < 0 || x >= m.X {
		fmt.Println("x out of range", x, m.X)
	}
	if y < 0 || y >= m.Y {
		fmt.Println("y out of range", y, m.Y)
	}
	m.m[x*m.Y+y] = v
}

func (m Matrix) Get(x int, y int) bool {
	return m.m[x*m.Y+y]
}
func (m Matrix) Drain() []Point {
	points := make([]Point, 0)
	for X := 0; X < m.X; X++ {
		for Y := 0; Y < m.Y; Y++ {
			if m.Get(X, Y) {
				points = append(points, Point{X, Y})
			}
		}
	}
	return points
}

func BuildIndexes(records []Record) []suffixarray.Index {
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
			kmer = RevComp(kmer)
		}

		// start := time.Now()
		offsets := x.LookupWithin(kmer, xa, xb, freqUp+1)
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
