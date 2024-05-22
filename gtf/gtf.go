package gtf

import (
	// "bufio"
	// "bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Feature struct {
	SeqName    string `json:"seqname"`
	Source     string `json:"source"`
	Type       string `json:"type"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
	Strand     string `json:"strand"`
	Attributes string `json:"attributes"`
}

func ParseGTF(f io.Reader) ([]Feature, error) {
	records := make([]Feature, 0)
	r := csv.NewReader(f)
	r.Comma = '\t'
	r.Comment = '#'
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(row) != 9 {
			return nil, fmt.Errorf("invalid GTF")
		}
		start, err := strconv.Atoi(row[3])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(row[4])
		if err != nil {
			return nil, err
		}
		record := Feature{
			SeqName:    row[0],
			Source:     row[1],
			Type:       row[2],
			Start:      start,
			End:        end,
			Strand:     row[6],
			Attributes: row[8],
		}
		records = append(records, record)
	}

	return records, nil
}

func ParseGTFFile(filename string) ([]Feature, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseGTF(f)
}
