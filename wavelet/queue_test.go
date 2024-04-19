package wavelet

import (
	"testing"
)

func TestQueue(t *testing.T) {
	h := NewQueue()

	// push
	h.HeapPush(Intersection{d: 10, aL: 0, aR: 10, bL: 10, bR: 20})
	h.HeapPush(Intersection{d: 30, aL: 10, aR: 100, bL: 10, bR: 200})
	h.HeapPush(Intersection{d: 20, aL: 0, aR: 10, bL: 20, bR: 20})

	t.Log(h)
	if h.Len() != 3 {
		t.Error()
	}

	// pop
	is := h.HeapPop()
	t.Log(is)
	if is.d != 30 {
		t.Error()
	}
	is = h.HeapPop()
	t.Log(is)
	if is.d != 10 {
		t.Error()
	}
	is = h.HeapPop()
	t.Log(is)
	if is.d != 20 {
		t.Error()
	}

	t.Log(h)
	if h.Len() != 0 {
		t.Error()
	}
}
