package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestIsDNA(t *testing.T) {
	oks := []string{
		"",
		"ATCGANTT",
		"AGCTTC",
	}
	for _, ok := range oks {
		if !IsDNA([]byte(ok)) {
			t.Errorf("%s should be valid DNA", ok)
		}
	}

	ngs := []string{
		"ggg",
		"ATCGAXXTT",
	}
	for _, ng := range ngs {
		if IsDNA([]byte(ng)) {
			t.Errorf("%s should not be valid DNA", ng)
		}
	}
}

func TestParseFasta(t *testing.T) {
	const multiFasta = `
	>SequenceA hogehoge fugafuga
	GATCnnCGcTGAc
	>SequenceB ori motif
	CTAG
	AGCTTTAG
	>SequenceC CTCF binding motif
	CCGCGNGGNGGCAG
	`

	data := strings.NewReader(multiFasta)
	records, _ := Parse(data)

	expected := []Record{
		{ID: []byte("SequenceA"), Seq: []byte("GATCNNCGCTGAC")},
		{ID: []byte("SequenceB"), Seq: []byte("CTAGAGCTTTAG")},
		{ID: []byte("SequenceC"), Seq: []byte("CCGCGNGGNGGCAG")},
	}

	/*
		for i, record := range records {
			fmt.Println("record", i, string(record.Seq), string(record.ID))
		}
		for i, record := range expected {
			fmt.Println("expected", i, string(record.Seq), string(record.ID))
		}
	*/

	for i := range expected {
		if !bytes.Equal(expected[i].ID, records[i].ID) {
			t.Errorf("ID is different")
		}
		if !bytes.Equal(expected[i].Seq, records[i].Seq) {
			t.Errorf("Seq is different")
		}
	}
}

func TestRevComp(t *testing.T) {
	if !bytes.Equal(RevComp([]byte("ATCGGTA")), []byte("TACCGAT")) {
		t.Errorf("different")
	}
	if !bytes.Equal(RevComp([]byte("A")), []byte("T")) {
		t.Errorf("different")
	}
	if !bytes.Equal(RevComp([]byte("")), []byte("")) {
		t.Errorf("different")
	}
}
