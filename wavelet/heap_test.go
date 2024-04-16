package wavelet

import (
	"testing"
)

func TestHeap(t *testing.T) {
	h := New()
	t.Log(h)

	h.HeapPush(Search{oL: 0, oR: 100, d: 1})
	h.HeapPush(Search{oL: 10, oR: 500, d: 2})
	h.HeapPush(Search{oL: 10, oR: 20, d: 0})

	for h.Len() > 0 {
		search := h.HeapPop()
		t.Log(search)
	}
}
