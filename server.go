package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"index/suffixarray"
	"log"
	"net/http"
	"os"
	"time"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func search() {
	data := []byte("hogehogehoge")
	index := suffixarray.New(data)
	q := []byte("hoge")
	offsets := index.Lookup(q, -1)
	for i, offset := range offsets {
		fmt.Printf("found: %d %d\n", i, offset)
	}
}

func loadPage(title string) *Page {
	filename := title + ".txt"
	body, _ := os.ReadFile(filename)
	return &Page{Title: title, Body: body}
}

// http.HandlerFunc:
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request to /")
	fmt.Fprintf(w, "Hi there %s", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p := loadPage(title)
	t, _ := template.ParseFiles("template.html")
	t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func saveAndLoad() {
	fmt.Println("Serving...")
	p1 := &Page{Title: "Test", Body: []byte("sample page")}
	p1.save()
	p2 := loadPage("Test")
	fmt.Println(string(p2.Body))
}

type Ping struct {
	Status int    `json:"status"`
	Result string `json:"result"`
	Points []Point
}

type Point struct {
	X int
	Y int
}

func (p Point) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("[%d,%d]", p.X, p.Y)
	return []byte(s), nil
}

func pointJsonTest() {
	p := Point{X: 10, Y: 20}
	points := []Point{p}
	// err := json.Unmarshal([]byte("{\"x\":10,\"y\":30}"), &p)
	s, err := json.Marshal(points)
	fmt.Println(p, string(s), err)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if v != nil {
		for key, value := range v {
			fmt.Println("got", key, value)
		}
	}
	time.Sleep(5 * time.Second)
	ping := Ping{http.StatusOK, "ok", []Point{}}
	res, err := json.Marshal(ping)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	filename := os.Args[1]
	fmt.Println("opening", filename)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open file %q: %v", filename, err)
	}
	defer f.Close()
	records, err := Parse(f)
	for i, record := range records {
		fmt.Println("record", i, string(record.Seq), string(record.ID))
	}

	pointJsonTest()
	search()
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/ping/", pingHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
