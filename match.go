package main

import (
	"encoding/json"
	// "flag"
	"fmt"
	// "html/template"
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

func FindMatch() []Point {
	// create suffix array
	points := make([]Point, 0)
	return points
}
