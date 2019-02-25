package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// cache is the in memory key/value store
type cache struct {
	mu      sync.Mutex
	keyVals map[string]*val
}

// val is the value in the cache key/value mapping
type val struct {
	operation  *Output
	expiration time.Time
}

// Output is the JSON output of entire query returned to user
type Output struct {
	Action string
	X      float64
	Y      float64
	Answer float64
	Cached bool
}

var datastore cache

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/add/", add)
	http.HandleFunc("/subtract/", subtract)
	http.HandleFunc("/multiply/", multiply)
	http.HandleFunc("/divide/", divide)

	datastore.new()
	go datastore.setTicker()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func add(response http.ResponseWriter, request *http.Request) {
	compute("add", response, request)
}

func subtract(response http.ResponseWriter, request *http.Request) {
	compute("subtract", response, request)
}

func multiply(response http.ResponseWriter, request *http.Request) {
	compute("multiply", response, request)
}

func divide(response http.ResponseWriter, request *http.Request) {
	compute("divide", response, request)
}

func (v *val) resetExpiration() {
	v.expiration = time.Now()
}

func parseValues(req *http.Request) (string, string, string) {
	err := req.ParseForm()
	if err != nil {
		log.Fatalln("Form fail: ", err)
	}

	raw := req.URL.RawQuery
	x := req.FormValue("x")
	y := req.FormValue("y")
	return x, y, raw
}

func writeResponse(out *Output, response http.ResponseWriter) []byte {
	jsonResult, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonResult)
	return jsonResult
}
