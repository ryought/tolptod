package wavelet

type Intersection struct {
	aL int
	aR int
	bL int
	bR int
	d  int  // Current depth
	c  byte // Last character of k-mer
}

type Queue []Intersection

func NewQueue() Queue {
	q := make([]Intersection, 0)
	return q
}

func (q *Queue) Len() int {
	return len(*q)
}

func (q *Queue) Pop() Intersection {
	i := (*q)[0]
	*q = (*q)[1:]
	return i
}

func (q *Queue) Push(i Intersection) {
	*q = append(*q, i)
}
