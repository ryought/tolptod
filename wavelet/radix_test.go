package wavelet

import (
	// "bytes"
	// "strings"
	"fmt"
	"testing"
)

func TestRadix(t *testing.T) {
	s := []byte("ATCGGATCGATTCG$")
	// workspace
	// var count [N]int
	// var pos [N]int
	// fmt.Println("index[1]", dest)

	index := RadixSort(s, 3)
	fmt.Println(index)
}
