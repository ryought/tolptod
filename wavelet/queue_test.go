package wavelet

import (
	"testing"
)

func TestQueue(t *testing.T) {
	h := NewQueue()

	// push
	h.Push(Intersection{d: 10})
	h.Push(Intersection{d: 20})
	h.Push(Intersection{d: 1})

	t.Log(h)
	if h.Len() != 3 {
		t.Error()
	}

	// pop
	if h.Pop().d != 10 {
		t.Error()
	}
	if h.Pop().d != 20 {
		t.Error()
	}
	if h.Pop().d != 1 {
		t.Error()
	}

	t.Log(h)
	if h.Len() != 0 {
		t.Error()
	}
}
