package bitvec

import (
	"bytes"
	"testing"
)

func TestV2(t *testing.T) {
	v := NewV2(16)
	v.Set(0, 1)
	v.Set(5, 1)
	v.Set(7, 1)
	v.Set(9, 1)
	v.Set(11, 1)
	v.Set(15, 1)
	v.Debug()
	t.Log("hoge")

	if !bytes.Equal(v.chunks, []byte{0b10100001, 0b10001010}) {
		t.Error()
	}

	xs := []byte{
		1, 0, 0, 0, 0, 1, 0, 1, //
		0, 1, 0, 1, 0, 0, 0, 1, //
	}
	for i, x := range xs {
		if v.Get(i) != x {
			t.Error()
		}
	}

	v.Set(9, 0)
	v.Set(15, 0)
	v.Set(5, 0)
	v.Debug()

	if !bytes.Equal(v.chunks, []byte{0b10000001, 0b00001000}) {
		t.Error()
	}
}
