package gtf

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func TestGTF(t *testing.T) {
	gtf := strings.Join([]string{
		"chr1\tHoge\tgene\t182\t882\t.\t+\t.\tID=hoge;name=fuga",
		"chr2\tHige\texon\t888\t999\t.\t-\t.\tID=hige;name=higee",
	}, "\n")
	res, err := ParseGTF(strings.NewReader(gtf))
	t.Log(res, err)
	expected := []Feature{
		{"chr1", "Hoge", "gene", 182, 882, "+", "ID=hoge;name=fuga"},
		{"chr2", "Hige", "exon", 888, 999, "-", "ID=hige;name=higee"},
	}
	if !slices.Equal(res, expected) {
		t.Error()
	}
}

func TestGTFFile(t *testing.T) {
	r, err := ParseGTFFile("./test.gtf")
	t.Log("hoge", r, len(r), err)
	if err != nil {
		t.Error()
	}
	if len(r) != 9 {
		t.Error()
	}
}

func TestGTFJSON(t *testing.T) {
	f := Feature{
		SeqName:    "hoge",
		Source:     "hoge",
		Type:       "fuga",
		Start:      0,
		End:        100,
		Strand:     "+",
		Attributes: "fuga",
	}
	res, _ := json.Marshal(f)
	t.Log(string(res))
	ex := `{"seqname":"hoge","source":"hoge","type":"fuga","start":0,"end":100,"strand":"+","attributes":"fuga"}`
	if string(res) != ex {
		t.Error()
	}
}
