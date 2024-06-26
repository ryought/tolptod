package main

import (
	"encoding/json"
	"fmt"
)

// 2D boolean matrix
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
		panic("x out of range")
	}
	if y < 0 || y >= m.Y {
		panic("y out of range")
	}
	m.m[y*m.X+x] = v
}

func (m Matrix) Get(x int, y int) bool {
	return m.m[y*m.X+x]
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
func (m Matrix) Print() {
	for Y := 0; Y < m.Y; Y++ {
		for X := 0; X < m.X; X++ {
			if m.Get(X, Y) {
				fmt.Printf("1 ")
			} else {
				fmt.Printf("0 ")
			}
		}
		fmt.Printf("\n")
	}
}

func (m Matrix) Equal(m2 Matrix) bool {
	// size should be the same
	if !(m.X == m2.X && m.Y == m2.Y) {
		return false
	}
	for X := 0; X < m.X; X++ {
		for Y := 0; Y < m.Y; Y++ {
			if m.Get(X, Y) != m2.Get(X, Y) {
				return false
			}
		}
	}
	return true
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
