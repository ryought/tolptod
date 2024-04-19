package bitvec

import (
	// "bytes"
	"testing"
)

func TestCountZeros(t *testing.T) {
	tests := []struct {
		offset int
		count  int
	}{
		{0, 0}, // [0,0) is empty
		{1, 0}, // [0,1) = 1
		{2, 1}, // [0,2) = 01
		{3, 1}, // [0,3) = 101
		{4, 2},
		{5, 3}, // [0,5) = 00101
		{6, 3}, // [0,6) = 100101
		{7, 4}, // [0,7) = 0100101
		{8, 5}, // [0,8) = 00100101
	}

	for _, test := range tests {
		count := countZeros(0b_0010_0101, test.offset)
		if count != test.count {
			t.Error()
		}
	}
}

func TestCreateChunk(t *testing.T) {
	b := [64]byte{0, 1, 0, 1, 1, 1, 0, 1}
	chunk := createChunk(b)
	t.Logf("%08b", chunk)
	if chunk != 0b_1011_1010 {
		t.Error()
	}
}

func TestV2Bit(t *testing.T) {
	v := NewV2(16)
	v.Set(0, 1)
	v.Set(5, 1)
	v.Set(7, 1)
	v.Set(9, 1)
	v.Set(11, 1)
	v.Set(15, 1)
	v.Debug()
	t.Log("hoge")

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

	xs = []byte{
		1, 0, 0, 0, 0, 0, 0, 1, //
		0, 0, 0, 1, 0, 0, 0, 0, //
	}
	for i, x := range xs {
		if v.Get(i) != x {
			t.Error()
		}
	}

	v.UpdateRank()
	v.Debug()
}

func TestV2RankSmall(t *testing.T) {
	v := NewV2(33)
	v.Set(0, 1)
	v.Set(5, 1)
	v.Set(7, 1)
	v.Set(9, 1)
	v.Set(11, 1)
	v.Set(15, 1)
	v.Set(20, 1)
	v.Set(25, 1)
	v.Set(32, 1)
	v.UpdateRank()
	v.Debug()

	tests := []struct {
		i int
		r int
	}{
		{0, 0},
		{1, 0},
		{2, 1},
		{6, 4},
		{32, 24},
		{33, 24},
	}
	for _, test := range tests {
		r := v.Rank(test.i)
		t.Log(test.i, test.r, r)
		if r != test.r {
			t.Error()
		}
	}
}

func TestV2RankLarge(t *testing.T) {
	v := NewV2(555)
	v.UpdateRank()
	v.Debug()

	// Since all bits are zero, rank(i)=i for all i.
	for i := 0; i <= 555; i++ {
		r := v.Rank(i)
		t.Logf("rank(i=%d)=%d\n", i, r)
		if r != i {
			t.Error()
		}
	}
}

func TestV2RankN16(t *testing.T) {
	v := NewV2(16)
	v.Set(0, 1)
	v.Set(5, 1)
	v.Set(7, 1)
	v.Set(9, 1)
	v.Set(14, 1)
	v.UpdateRank()
	v.Debug()
	if v.Rank(0) != 0 {
		t.Error()
	}
	if v.Rank(6) != 4 {
		t.Error()
	}
	if v.Rank(16) != 11 {
		t.Error()
	}
}

func TestV2Build(t *testing.T) {
	n := 34
	x := []byte{
		0, 1, 0, 0, 1, 0, 1, 0, // 8
		1, 0, 1, 0, 0, 0, 0, 1, // 16
		0, 1, 0, 1, 1, 1, 0, 1, // 24
		1, 1, 1, 0, 0, 0, 1, 1, // 32
		0, 1, // 34
	}
	if len(x) != n {
		t.Error()
	}
	v := Build(n, func(i int) byte { return x[i] })
	v.Debug()

	t.Log("get")
	for i := 0; i < n; i++ {
		if x[i] != v.Get(i) {
			t.Errorf("Get: x[i=%d]=%d != %d", i, v.Get(i), x[i])
		}
	}

	t.Log("rank")
	for i := 0; i <= n; i++ {
		if naiveRank(x[:i]) != v.Rank(i) {
			t.Error("rank is different", i)
		}
	}
}

func naiveRank(b []byte) int {
	rank := 0
	for i := range b {
		if b[i] == 0 {
			rank += 1
		}
	}
	return rank
}
