package wavelet

import (
	// "bytes"
	// "strings"
	// "fmt"
	// "github.com/ryought/tolptod/fasta"
	"testing"
	"time"
)

func TestRadix(t *testing.T) {
	s := []byte("ATCGGATCGATTCG$")
	// workspace
	// var count [N]int
	// var pos [N]int
	// fmt.Println("index[1]", dest)

	index := RadixSort(s, 5)
	// index := Sort(s)
	// fmt.Println(index)
	printIndex(s, index, 10)

	if !isSorted(s, index, 5) {
		t.Error("not sorted")
	}
}

func TestRadixLarge(t *testing.T) {
	s := RandomDNA(1_000_000)
	s = append(s, '$')
	t0 := time.Now()
	index := RadixSort(s, 40)
	t.Logf("radix %d ms", time.Since(t0).Milliseconds())
	if !isSorted(s, index, 40) {
		t.Error("not sorted")
	}

	t1 := time.Now()
	index = Sort(s)
	t.Logf("sort %d ms", time.Since(t1).Milliseconds())
	if !isSorted(s, index, 40) {
		t.Error("not sorted")
	}
	// printIndex(s, index, 40)
}

func BenchmarkRadixLarge(b *testing.B) {
	// records, _ := fasta.ParseFile("../chr1.fa")
	// s := records[0].Seq
	s := RandomDNA(1_000_000) // 1MB
	b.Log("building", len(s))
	b.StartTimer()
	RadixSort(s, 40)
	// Sort(s)
	b.StopTimer()
}
