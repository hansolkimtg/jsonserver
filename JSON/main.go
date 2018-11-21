package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/appengine"
)

// Employee employee
type Employee struct {
	ID       string `datastore:"-" goon:"id"`
	Name     string
	Role     string
	HireDate time.Time
	Account  string
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
	appengine.Main()
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)

	length, _ := strconv.Atoi(r.Header.Get("Content-Length"))

	body := make([]byte, length)
	length, _ = r.Body.Read(body)

	var jsonBody map[string]Employee

	_ = json.Unmarshal(body[:length], &jsonBody)

	fmt.Printf("%v\n", jsonBody)

	// v := r.URL.Query()
	// if v == nil {
	// 	return
	// }

	// for key, vs := range v{
	// 	fmt.Fprintf(w, "%s = %s\n", key, vs[0])
	// }
}
