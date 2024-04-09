package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type Record struct {
	ID  []byte
	Seq []byte
}

func IsDNA(b []byte) bool {
	for _, c := range b {
		if c != 'A' && c != 'C' && c != 'G' && c != 'T' && c != 'N' {
			return false
		}
	}
	return true
}

func Parse(f io.Reader) ([]Record, error) {
	records := make([]Record, 0)
	r := bufio.NewReader(f)

	var id, seq []byte

	for {
		line, _, err := r.ReadLine()

		// end of file
		if err == io.EOF {
			record := Record{ID: id, Seq: seq}
			records = append(records, record)
			break
		}

		// read error
		if err != nil {
			return nil, err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == '>' {
			// save previous id and seq
			if len(id) > 0 {
				if !IsDNA(seq) {
					return nil, fmt.Errorf("Contains non-DNA character")
				}
				record := Record{ID: id, Seq: seq}
				records = append(records, record)
			}

			// store new id
			seg := bytes.Fields(line[1:])[0]
			id = make([]byte, len(seg))
			copy(id, seg)
			seq = nil
		} else {
			line = bytes.ToUpper(line)
			seq = append(seq, line...)
		}
	}

	return records, nil
}

func ParseFile(filename string) ([]Record, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func PrintRecords(records []Record) {
	for i, record := range records {
		fmt.Println("record", i, string(record.ID), len(record.Seq))
	}
}
