package wavelet

import "fmt"

// Create index [0,1,2,3,...,n-1]
func NewIndex(n int) []int {
	index := make([]int, n)
	for i := range index {
		index[i] = i
	}
	return index
}

// alphabet size
const N = 256

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

func printIndex(s []byte, index []int) {
	for i := range index {
		fmt.Printf("index[%d]=%d\t%s\n", i, index[i], s[index[i]:])
	}
}

// MSD Radix (stable) sort for suffixes
// https://en.wikipedia.org/wiki/Radix_sort
func RadixSortForDigit(s []byte, index []int, dest []int, k int) (count [N]int, pos [N]int) {
	n := len(s)

	// count character occurrences in s
	for _, i := range index {
		c := s[(i+k)%n]
		count[c] += 1
	}

	// cumurative count
	pos[0] = 0
	for c := 1; c < N; c++ {
		pos[c] = pos[c-1] + count[c-1]
	}

	for x, i := range index {
		c := s[(i+k)%n]
		dest[pos[c]] = index[x]
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
	queue := make([]Job, 0)
	index = NewIndex(n)
	dest := make([]int, n)
	// var dest []int
	queue = append(queue, Job{k: 0, i: 0, j: n})

	for len(queue) > 0 {
		fmt.Println("[queue]", queue)
		job := queue[0]
		queue = queue[1:]

		fmt.Printf("[radix] i=%d j=%d k=%d\n", job.i, job.j, job.k)
		fmt.Println("before", index)
		printIndex(s, index)
		count, pos := RadixSortForDigit(s, index[job.i:job.j], dest[job.i:job.j], job.k)
		fmt.Println("###count###")
		print(count)
		fmt.Println("###pos###")
		print(pos)
		copy(index[job.i:job.j], dest[job.i:job.j])
		fmt.Println("after", index)
		printIndex(s, index)

		// register new jobs
		for c := 0; c < N; c++ {
			if count[c] >= 2 && job.k < d {
				i := job.i
				if c > 0 {
					i = job.i + pos[c-1]
				}
				j := job.i + pos[c]
				job := Job{k: job.k + 1, i: i, j: j}
				queue = append(queue, job)
				fmt.Println("added job", len(queue), job)
			}
		}
		fmt.Println("queue", len(queue))
	}
	return
}
