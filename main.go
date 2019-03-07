package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// val is the value in the cache key/value mapping
type val struct {
	operation  Output
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
var routes map[string]bool

func init() {
	routes = map[string]bool{
		"add":      true,
		"subtract": true,
		"multiply": true,
		"divide":   true,
	}
}

func main() {
	http.HandleFunc("/", handler)

	// http.Handle("/", http.FileServer(http.Dir("./static")))
	// http.Handle("/welcome", http.FileServer(http.Dir("./static")))

	datastore.new()
	go datastore.setTicker()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("req.URL.Path ", request.URL.Path)

	path := strings.Split(request.URL.Path, "/")
	action := strings.Join(path, "")
	isComputation := routes[action]

	if len(path) > 1 && isComputation {
		checkInputs(action, response, request)
	}
}

func (v *val) resetExpiration() {
	v.expiration = time.Now()
}

func parseValues(req *http.Request) (string, string, string) {
	// err := req.ParseForm()
	// if err != nil {
	// 	log.Fatalln("Form fail: ", err)
	// }

	raw := req.URL.RawQuery
	x := req.FormValue("x")
	y := req.FormValue("y")
	return x, y, raw
}

func writeResponse(out Output, response http.ResponseWriter, errType bool) {
	errBytes := []byte("Error. Failed calculation, try again")

	if errType {
		response.Write(errBytes)
		return
	}

	jsonResult, err := json.Marshal(out)
	if err != nil {
		response.Write(errBytes)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(jsonResult)
}
