package main

import (
	"encoding/json"
	// "flag"
	"fmt"
	// "html/template"
	"log"
	"net/http"
	"os"
	"time"
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
	X       int `json:"x"`
	Y       int `json:"y"`
	XA      int `json:"xA"`
	XB      int `json:"xB"`
	YA      int `json:"yA"`
	YB      int `json:"yB"`
	K       int `json:"k"`
	FreqLow int `json:"freqLow"`
	FreqUp  int `json:"freqUp"`
	Scale   int `json:"scale"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if v != nil {
		for key, value := range v {
			fmt.Println("got", key, value)
		}
	}
	time.Sleep(5 * time.Second)
	ping := Ping{http.StatusOK, "ok"}
	res, err := json.Marshal(ping)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Write(res)
}

func createGenerateHandler(xrs []Record, yrs []Record) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("json")
		var req Request
		if err := json.Unmarshal([]byte(body), &req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Println("requeted..", req, req.K)
		points := FindMatch(xrs[req.X].Seq[req.XA:req.XB], yrs[req.Y].Seq[req.YA:req.YB], req.Scale, req.K)
		// points := []Point{
		// 	{X: 0, Y: 0},
		// 	{X: 1, Y: 1},
		// 	{X: 2, Y: 2},
		// }
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
		fmt.Println("info", info)
		res, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("res", string(res))
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Write(res)
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: tolptod x.fa y.fa")
	}
	xrs, err := ParseFile(os.Args[1])
	if err != nil {
		log.Fatalf("ParseFile error: %s", err)
	}
	yrs, err := ParseFile(os.Args[2])
	if err != nil {
		log.Fatalf("ParseFile error: %s", err)
	}
	PrintRecords(xrs)
	PrintRecords(yrs)

	info := toInfo(xrs, yrs)

	// pointJsonTest()
	// search()
	// http.HandleFunc("/", handler)
	// http.HandleFunc("/view/", viewHandler)
	// http.HandleFunc("/edit/", editHandler)
	// http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/", createInfoHandler(info))
	http.HandleFunc("/generate/", createGenerateHandler(xrs, yrs))
	http.HandleFunc("/ping/", pingHandler)

	fmt.Println("running..")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
