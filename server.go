package main

import (
	"encoding/json"
	"flag"

	// "html/template"
	_ "embed"
	"log"
	"net/http"

	"github.com/ryought/tolptod/fasta"
	"github.com/ryought/tolptod/suffixarray"
)

//go:embed app/dist/index.html
var html []byte

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

func createGenerateHandler(xis []suffixarray.Index, yrs []fasta.Record) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("json")
		var req Request
		if err := json.Unmarshal([]byte(body), &req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println("/generate requested", req, req.K)
		points := FindMatch(xis[req.X], req.XA, req.XB, yrs[req.Y].Seq[req.YA:req.YB], req.Scale, req.K, req.FreqLow, req.FreqUp, req.Revcomp)
		log.Println("matching done")
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/ requested")
	w.Write(html)
}

var addr = flag.String("b", ":8080", "address:port to bind. Default: localhost and port 8080.")
var help = flag.Bool("h", false, "show help")

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 || *help {
		log.Fatalf("Usage v0.2: tolptod -b localhost:8080 x.fa y.fa")
	}

	// parse fasta
	log.Println("Parsing", args[0])
	xrs, err := fasta.ParseFile(args[0])
	if err != nil {
		log.Fatalf("ParseFile error: %s", err)
	}
	log.Println("Parsing", args[1])
	yrs, err := fasta.ParseFile(args[1])
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

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/info/", createInfoHandler(info))
	http.HandleFunc("/generate/", createGenerateHandler(indexes, yrs))

	log.Printf("Server running on %s...", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
