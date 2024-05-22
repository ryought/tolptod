package gtf

import (
	"testing"
)

func TestGTF(t *testing.T) {
	r, err := ParseGTFFile("./test.gtf")
	t.Log("hoge", r, err)
}
