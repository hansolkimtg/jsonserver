package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine/log"

	"github.com/mjibson/goon"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Employee employee
type Employee struct {
	ID string `json:"id" datastore:"-" goon:"id"`
	//`datastore:"-" goon:"id"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
	HireDate time.Time `json:"hiredate"` // 2014-12-31 08:04:18 +0900 JST
	Account  string    `json:"account"`
}

func main() {
	http.HandleFunc("/jsonpost", handlerJSONpost)
	http.HandleFunc("/jsonget", handlerJSONget)
	http.HandleFunc("/get", handlerGet)
	http.HandleFunc("/post", handlerPost)
	http.HandleFunc("/put", handlerPut)
	http.HandleFunc("/delete", handlerDelete)
	http.HandleFunc("/goon-post", handlerGoonPost)
	http.HandleFunc("/goon-get", handlerGoonGet)
	http.HandleFunc("/goon-put", handlerGoonPut)
	http.HandleFunc("/goon-delete", handlerGoonDelete)
	appengine.Main()
}

func handlerJSONget(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	// get id from Request
	id := r.URL.Query().Get("id")

	emp := &Employee{
		ID: id,
	}

	// get entity from id
	if err := g.Get(emp); err != nil {
		log.Errorf(c, "could not get the record: %v", err)
		http.Error(w, "An error occurred.", http.StatusInternalServerError)
	}

	// convert to JSON
	jsonResp, err := json.Marshal(emp)
	if err != nil {
		log.Errorf(c, "could not convert to JSON: %v", err)
		http.Error(w, "An error occurred.", http.StatusInternalServerError)
	}

	// print out
	fmt.Fprintln(w, fmt.Sprintf("emp: %s", jsonResp))
}

func handlerJSONpost(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	//fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)

	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "err %v", err.Error())
	}

	log.Infof(c, "json: %s", string(jsonBytes))

	var emp Employee

	if err := json.Unmarshal(jsonBytes, &emp); err != nil {
		fmt.Fprintln(w, "1st error")
		log.Errorf(c, "1st error: %v", err)
	}

	log.Infof(c, "%v", time.Now())

	g.Put(&emp)
}

func handlerGoonPost(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	emp := &Employee{
		ID:       "id",
		Name:     "mako",
		Role:     "Worker",
		HireDate: time.Now(),
		Account:  "kim",
	}

	g.Put(emp)
}

func handlerGoonDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	employee := &Employee{
		ID: "id",
	}

	g.Delete(g.Key(employee))

	log.Debugf(c, "emp: %#v", employee)
}
func handlerGoonPut(w http.ResponseWriter, r *http.Request) {
	g := goon.NewGoon(r)

	emp := &Employee{
		ID:       "id",
		Name:     "kim",
		Role:     "Worker",
		HireDate: time.Now(),
		Account:  "kim",
	}

	g.Put(emp)
}

func handlerGoonGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	employee := &Employee{
		ID: "id",
	}

	if err := g.Get(employee); err != nil {
		fmt.Fprintln(w, "no!")
		log.Errorf(c, "could not get the record: %v", err)
		http.Error(w, "An error occurred.", http.StatusInternalServerError)

	}
	fmt.Fprintln(w, fmt.Sprintf("emp: %v", employee))

}

func handlerDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "employee", "abc", 0, nil)
	_ = datastore.Delete(c, key)
}

func handlerPut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	e2 := Employee{
		Name:     "Mary Citizen",
		Role:     "Worker",
		HireDate: time.Now(),
		Account:  "kim",
	}

	key := datastore.NewKey(c, "employee", "abc", 0, nil)
	_, _ = datastore.Put(c, key, &e2)
}

func handlerPost(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	e1 := Employee{
		Name:     "Joe Citizen",
		Role:     "Manager",
		HireDate: time.Now(),
		Account:  "kim",
	}

	key := datastore.NewKey(c, "employee", "abc", 0, nil)
	_, _ = datastore.Put(c, key, &e1)
}

func handlerGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	employeeKey := datastore.NewKey(c, "employee", "abc", 0, nil)
	//employeeKey := datastore.NewIncompleteKey(c, "employee", nil)

	var employee Employee

	err := datastore.Get(c, employeeKey, &employee)

	fmt.Fprintln(w, fmt.Sprintf("emp: %v", employee))

	log.Debugf(c, "emp: %#v", employee)

	if err != nil {
		fmt.Fprintln(w, "no!")
	}

}
