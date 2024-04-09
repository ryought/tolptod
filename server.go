package main

import (
	"encoding/json"
	// "flag"
	// "html/template"
	"index/suffixarray"
	"log"
	"net/http"
	"os"
)

type Ping struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

type Plot struct {
	Status int     `json:"status"`
	Result string  `json:"result"`
	Points []Point `json:"points"`
}

type Request struct {
	X       int  `json:"x"`
	Y       int  `json:"y"`
	XA      int  `json:"xA"`
	XB      int  `json:"xB"`
	YA      int  `json:"yA"`
	YB      int  `json:"yB"`
	K       int  `json:"k"`
	FreqLow int  `json:"freqLow"`
	FreqUp  int  `json:"freqUp"`
	Scale   int  `json:"scale"`
	Revcomp bool `json:"revcomp"`
}

func createGenerateHandler(xis []suffixarray.Index, yrs []Record) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("json")
		var req Request
		if err := json.Unmarshal([]byte(body), &req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println("/generate requested", req, req.K)
		points := FindMatch(xis[req.X], req.XA, req.XB, yrs[req.Y].Seq[req.YA:req.YB], req.Scale, req.K, req.FreqLow, req.FreqUp, req.Revcomp)
		plot := Plot{http.StatusOK, "ok", points}

		res, err := json.Marshal(plot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Write(res)
	}
}

func createInfoHandler(info Info) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("/info requested", info)
		res, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Write(res)
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: tolptod x.fa y.fa")
	}

	// parse fasta
	xrs, err := ParseFile(os.Args[1])
	if err != nil {
		log.Fatalf("ParseFile error: %s", err)
	}
	yrs, err := ParseFile(os.Args[2])
	if err != nil {
		log.Fatalf("ParseFile error: %s", err)
	}

	// dump seq infos
	for i, r := range xrs {
		log.Printf("x #%d %s (%d bp)\n", i, string(r.ID), len(r.Seq))
	}
	for i, r := range yrs {
		log.Printf("y #%d %s (%d bp)\n", i, string(r.ID), len(r.Seq))
	}

	// build
	log.Println("Building suffix array...")
	indexes := BuildIndexes(xrs)
	log.Println("Done")
	info := toInfo(xrs, yrs)

	http.HandleFunc("/", createInfoHandler(info))
	http.HandleFunc("/generate/", createGenerateHandler(indexes, yrs))

	log.Println("Server running on 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
