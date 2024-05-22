package fasta

type Info struct {
	Xs []Seq `json:"xs"`
	Ys []Seq `json:"ys"`
}

type Seq struct {
	Id  string `json:"id"`
	Len int    `json:"len"`
}

func ToSeqInfo(rs []Record) []Seq {
	is := make([]Seq, len(rs))
	for i, r := range rs {
		is[i].Id = string(r.ID)
		is[i].Len = len(r.Seq)
	}
	return is
}

func ToInfo(xrs []Record, yrs []Record) Info {
	return Info{
		Xs: ToSeqInfo(xrs),
		Ys: ToSeqInfo(yrs),
	}
}
