package radix

import (
	"bytes"
	"fmt"
	"slices"
)

// Create index [0,1,2,3,...,n-1]
func NewIndex(n int) []int {
	index := make([]int, n)
	for i := range index {
		index[i] = i
	}
	return index
}

// alphabet size
// const N = 1 << 8
const N = 4

// workspace reset
func reset(v [N]int) {
	for c := range v {
		v[c] = 0
	}
}

func print(v [N]int) {
	cs := []byte{'A', 'C', 'G', 'T'}
	for _, c := range cs {
		fmt.Printf("v[%c]=%d\n", c, v[c])
	}
}

func printIndex(s []byte, index []int, k int) {
	for i := range index {
		fmt.Printf("index[%d]=%d\t%s\n", i, index[i], Slice(s, index[i], k))
	}
}

func Slice(s []byte, i int, n int) []byte {
	return s[i:min(i+n, len(s))]
}

// index is sorted
// s[index[i]:] < s[index[i+1]:]
func isSorted(s []byte, index []int, k int) bool {
	for i := 1; i < len(index); i++ {
		a := Slice(s, index[i-1], k)
		b := Slice(s, index[i], k)
		if bytes.Compare(a, b) > 0 {
			return false
		}
	}
	return true
}

// MSD Radix (stable) sort for suffixes
// https://en.wikipedia.org/wiki/Radix_sort
// count = [['A', 0, 10]]
// changed
func RadixSortForDigit(s []byte, index []int, dest []int, k int) (count [N]int, pos [N]int) {
	n := len(s)
	w := len(index)

	if w == 1 {
		return
	}

	// count character occurrences in s
	for _, i := range index {
		count[s[(i+k)%n]] += 1
	}

	// cumurative count
	i := 0
	for c := 0; c < N; c++ {
		pos[c] = i
		i += count[c]
	}

	for _, i := range index {
		c := s[(i+k)%n]
		dest[pos[c]] = i
		pos[c] += 1
	}

	return
}

type Job struct {
	k int
	i int
	j int
}

// Dispatch sort jobs
func RadixSort(s []byte, d int) (index []int) {
	n := len(s)
	queue := make([]Job, 0, 1000)
	index, dest := NewIndex(n), NewIndex(n)

	// first job
	queue = append(queue, Job{k: 0, i: 0, j: n})

	for len(queue) > 0 {
		job := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		// sort
		count, pos := RadixSortForDigit(s, index[job.i:job.j], dest[job.i:job.j], job.k)
		copy(index[job.i:job.j], dest[job.i:job.j])

		// register new jobs
		if job.k < d {
			for c := 0; c < N; c++ {
				if count[c] >= 2 {
					job := Job{
						k: job.k + 1,
						i: job.i + pos[c] - count[c],
						j: job.i + pos[c],
					}
					queue = append(queue, job)
				}
			}
		}
	}
	return
}

// use standard sort
func Sort(s []byte) (index []int) {
	index = NewIndex(len(s))
	slices.SortFunc(index, func(i, j int) int {
		return bytes.Compare(s[i:], s[j:])
	})
	return index
}
