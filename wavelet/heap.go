package wavelet

import (
	"container/heap"
)

type PriorityQueue []*Search

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Size() > pq[j].Size()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*Search))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func NewPriorityQueue() PriorityQueue {
	h := PriorityQueue{}
	heap.Init(&h)
	return h
}

func (pq *PriorityQueue) HeapPush(search Search) {
	heap.Push(pq, &search)
}

func (pq *PriorityQueue) HeapPop() *Search {
	return heap.Pop(pq).(*Search)
}
