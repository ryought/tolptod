package wavelet

import (
	"container/heap"
)

type Intersection struct {
	aL int
	aR int
	bL int
	bR int
	d  int  // Current depth
	c  byte // Last character of k-mer
}

func (is Intersection) Priority() int {
	return min(is.aR-is.aL, is.bR-is.bL)
}

type Queue []*Intersection

func NewQueue() Queue {
	q := Queue{}
	// q := make([]*Intersection, 0, 1024)
	heap.Init(&q)
	return q
}

func (q *Queue) HeapPush(is Intersection) {
	heap.Push(q, &is)
}

func (q *Queue) HeapPop() *Intersection {
	return heap.Pop(q).(*Intersection)
}

func (q Queue) Len() int {
	return len(q)
}

func (q Queue) Less(i, j int) bool {
	return q[i].Priority() > q[j].Priority()
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *Queue) Push(x any) {
	*q = append(*q, x.(*Intersection))
}

func (q *Queue) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*q = old[0 : n-1]
	return item
}

func NewStack() Queue {
	q := make(Queue, 0, 1024)
	return q
}

func (q *Queue) StackPop() *Intersection {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*q = old[0 : n-1]
	return item
}

func (q *Queue) StackPush(is Intersection) {
	q.Push(&is)
}
