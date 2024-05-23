package main

import (
	"encoding/json"
	"flag"

	// "html/template"
	_ "embed"
	"log"
	"net/http"
	"strings"

	"github.com/ryought/tolptod/fasta"
	"github.com/ryought/tolptod/gtf"
)

//go:embed app/dist/index.html
var html []byte

// versions injected when build using goreleaser
var (
	Version  = "unset"
	Revision = "unset"
)

type Ping struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

type Plot struct {
	Status   int     `json:"status"`
	Result   string  `json:"result"`
	Forward  []Point `json:"forward"`
	Backward []Point `json:"backward"`
}

type Features struct {
	Status int           `json:"status"`
	Result string        `json:"result"`
	X      []gtf.Feature `json:"x"`
	Y      []gtf.Feature `json:"y"`
}

type Request struct {
	X        int  `json:"x"`
	Y        int  `json:"y"`
	XA       int  `json:"xA"`
	XB       int  `json:"xB"`
	YA       int  `json:"yA"`
	YB       int  `json:"yB"`
	K        int  `json:"k"`
	FreqLow  int  `json:"freqLow"`
	FreqUp   int  `json:"freqUp"`
	Scale    int  `json:"scale"`
	UseCache bool `json:"useCache"`
}

func createGenerateHandler(index IndexV2, cache *Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body := r.FormValue("json")
		var req Request
		if err := json.Unmarshal([]byte(body), &req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var forward, backward Matrix
		config := Config{
			k:       req.K,
			bin:     req.Scale,
			freqLow: req.FreqLow,
			freqUp:  req.FreqUp,
			xL:      req.XA,
			xR:      req.XB,
			yL:      req.YA,
			yR:      req.YB,
		}
		if req.UseCache {
			log.Println("/generate requested with cache", req, req.K)
			p := 0 // progress percentage
			forward, backward = cache.ComputeMatrixWithProgress(
				ctx, config,
				func(y, yL, yR int) {
					newp := 100 * (y - yL) / (yR - yL)
					if newp > p {
						p = newp
						log.Printf("progress %d%% (y=%d in [%d, %d])\n", p, y, yL, yR)
					}
				},
			)
			log.Println("matching done")
		} else {
			log.Println("/generate requested", req, req.K)
			p := 0 // progress percentage
			forward, backward = ComputeMatrixWithProgress(
				ctx, index.xindex[req.X], index.yindex[req.Y], config,
				func(y, yL, yR int) {
					newp := 100 * (y - yL) / (yR - yL)
					if newp > p {
						p = newp
						log.Printf("progress %d%% (y=%d in [%d, %d])\n", p, y, yL, yR)
					}
				},
			)
			log.Println("matching done")
		}
		plot := Plot{http.StatusOK, "ok", forward.Drain(), backward.Drain()}

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

func createCacheHandler(index IndexV2, cache *Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body := r.FormValue("json")
		var req Request
		if err := json.Unmarshal([]byte(body), &req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println("/cache requested", req.X, req.Y, req.K, req.Scale)
		p := 0 // progress percentage
		*cache = NewCache(
			ctx,
			index.xindex[req.X],
			index.yindex[req.Y],
			Config{
				k:       req.K,
				bin:     req.Scale,
				freqLow: req.FreqLow,
				freqUp:  req.FreqUp,
			},
			func(y, yL, yR int) {
				newp := 100 * (y - yL) / (yR - yL)
				if newp > p {
					p = newp
					log.Printf("progress %d%% (y=%d in [%d, %d])\n", p, y, yL, yR)
				}
			},
		)
		log.Println("cache done")
		res, err := json.Marshal("ok")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Write(res)
	}
}

func createInfoHandler(info fasta.Info) http.HandlerFunc {
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

func createFeaturesHandler(xf, yf gtf.GTFTree) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("json")
		var req Request
		if err := json.Unmarshal([]byte(body), &req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println("/features requested", req.X, req.Y, req.K, req.Scale)
		xFs := xf.Find(req.X, req.XA, req.XB)
		yFs := yf.Find(req.Y, req.YA, req.YB)
		res, err := json.Marshal(Features{
			Status: http.StatusOK,
			Result: "ok",
			X:      xFs,
			Y:      yFs,
		})
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

func createGFFTree(seqs []fasta.Seq, f string) gtf.GTFTree {
	if f == "" {
		return gtf.BuildIntervalTree(seqs, []gtf.Feature{})
	}
	if strings.HasSuffix(f, "bed") {
		log.Println("Parsing", f, "as BED")
		features, err := gtf.ParseBEDFile(f)
		if err != nil {
			log.Fatalf("ParseBED error: %s", err)
		}
		return gtf.BuildIntervalTree(seqs, features)
	} else {
		log.Println("Parsing", f, "as GFF/GTF")
		features, err := gtf.ParseGTFFile(f)
		if err != nil {
			log.Fatalf("ParseGTF error: %s", err)
		}
		return gtf.BuildIntervalTree(seqs, features)
	}
}

var addr = flag.String("b", ":8080", "address:port to bind. Default: localhost and port 8080.")
var help = flag.Bool("h", false, "show help")
var version = flag.Bool("v", false, "show version")
var xGFF = flag.String("x", "", "annotation GFF/GTF file for x.fa")
var yGFF = flag.String("y", "", "annotation GFF/GTF file for y.fa")

func main() {
	flag.Parse()
	args := flag.Args()

	if *version {
		log.Printf("%s (%s)", Version, Revision)
		return
	}
	if len(args) < 2 || *help {
		log.Fatalf("Usage %s: tolptod -b localhost:8080 x.fa y.fa", Version)
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
	info := fasta.ToInfo(xrs, yrs)

	// dump seq infos
	for i, r := range xrs {
		log.Printf("x #%d %s (%d bp)\n", i, string(r.ID), len(r.Seq))
	}
	for i, r := range yrs {
		log.Printf("y #%d %s (%d bp)\n", i, string(r.ID), len(r.Seq))
	}

	// GFF parser
	xf := createGFFTree(info.Xs, *xGFF)
	yf := createGFFTree(info.Ys, *yGFF)

	// build
	log.Println("Building suffix array...")
	indexes := NewIndexV2FromRecords(xrs, yrs)
	log.Println("Done")
	var cache Cache

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/info/", createInfoHandler(info))
	http.HandleFunc("/cache/", createCacheHandler(indexes, &cache))
	http.HandleFunc("/generate/", createGenerateHandler(indexes, &cache))
	http.HandleFunc("/features/", createFeaturesHandler(xf, yf))

	log.Printf("Server running on %s...", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
