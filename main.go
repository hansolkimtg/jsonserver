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
	ID       string    `json:"id" datastore:"-" goon:"id"`
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

// jSONでget処理を行う
func handlerJSONget(w http.ResponseWriter, r *http.Request) {
	// リクエストからappengine, goonのコンテキストを生成する
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	// リクエストから「id」を取得する
	id := r.URL.Query().Get("id")

	// Employeeに取得した「id」を設定する
	emp := &Employee{
		ID: id,
	}

	// 設定した「id」を元に、DataStoreからデータを取得する
	if err := g.Get(emp); err != nil { // 取得できなかった場合
		// ログ・エラーメッセージを表示する
		log.Errorf(c, "could not get the record: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// 取得したデータをJSONに変換する
	jsonResp, err := json.Marshal(emp)
	if err != nil { // 変換できなかった場合
		// ログ・エラーメッセージを出力する
		log.Errorf(c, "could not convert to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// 取得データを画面に出力する
	fmt.Fprintln(w, fmt.Sprintf("emp: %s", jsonResp))
}

// jSONでpost処理を行う
func handlerJSONpost(w http.ResponseWriter, r *http.Request) {
	// リクエストからappengine, goonのコンテキストを生成する
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	// リクエストBodyのデータを取得する
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil { // 取得できなかった場合
		// ログを出力する
		log.Errorf(c, "err %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Employee型の変数を生成する
	var emp Employee

	// リクエストBodyのデータをEmployee型に変換する
	if err := json.Unmarshal(jsonBytes, &emp); err != nil { // 変換できなかった場合
		// ログ・エラーメッセージを出力する
		fmt.Fprintln(w, "1st error")
		log.Errorf(c, "1st error: %v", err)
	}

	// 情報ログを出力する
	log.Infof(c, "%v", time.Now())

	// Put処理を行う
	g.Put(&emp)
}

// GOONでpost処理を行う
func handlerGoonPost(w http.ResponseWriter, r *http.Request) {
	// リクエストからappengine, goonのコンテキストを生成する
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	// Employee型の変数を定義する
	emp := &Employee{
		ID:       "id",
		Name:     "mako",
		Role:     "Worker",
		HireDate: time.Now(),
		Account:  "kim",
	}

	// GOONのput処理を行う
	g.Put(emp)
}

// GOONでdelete処理を行う
func handlerGoonDelete(w http.ResponseWriter, r *http.Request) {
	// リクエストからappengine, goonのコンテキストを生成する
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	// キー付のEmployee型の変数を生成する
	employee := &Employee{
		ID: "id",
	}

	// 定義されたキーを元に、GOONのdelete処理を行う
	g.Delete(g.Key(employee))

	// デバッグのログを出力する
	log.Debugf(c, "emp: %#v", employee)
}

// GOONでput処理を行う
func handlerGoonPut(w http.ResponseWriter, r *http.Request) {
	// リクエストからgoonのコンテキストを生成する
	g := goon.NewGoon(r)

	// Employee型の変数を定義する
	emp := &Employee{
		ID:       "id",
		Name:     "kim",
		Role:     "Worker",
		HireDate: time.Now(),
		Account:  "kim",
	}

	// GOONのput処理を行う
	g.Put(emp)
}

// GOONでget処理を行う
func handlerGoonGet(w http.ResponseWriter, r *http.Request) {
	// リクエストからappengine, goonのコンテキストを生成する
	c := appengine.NewContext(r)
	g := goon.FromContext(c)

	// キー付のEmployee型の変数を生成する
	employee := &Employee{
		ID: "id",
	}

	// GOONのget処理を行う
	if err := g.Get(employee); err != nil { // 取得できなかった場合
		// ログ・エラーメッセージを出力する
		log.Errorf(c, "could not get the record: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// 取得したEmployeeを出力する
	fmt.Fprintln(w, fmt.Sprintf("emp: %v", employee))
}

// GOONを使わずにdelete処理を行う
func handlerDelete(w http.ResponseWriter, r *http.Request) {
	// リクエストからappengineのコンテキストを生成する
	c := appengine.NewContext(r)

	// DataStoreのキーを生成する
	key := datastore.NewKey(c, "employee", "abc", 0, nil)
	// 設定したキーを元に、対象データを削除する
	_ = datastore.Delete(c, key)
}

// GOONを使わずにput処理を行う
func handlerPut(w http.ResponseWriter, r *http.Request) {
	// リクエストからappengineのコンテキストを生成する
	c := appengine.NewContext(r)

	// Employee型の変数を定義する
	e2 := Employee{
		Name:     "Mary Citizen",
		Role:     "Worker",
		HireDate: time.Now(),
		Account:  "kim",
	}

	//　DataStoreのキーを生成する
	key := datastore.NewKey(c, "employee", "abc", 0, nil)
	// 設定したキーを元に、put処理を行う
	_, _ = datastore.Put(c, key, &e2)
}

// GOONを使わずにpost処理を行う
func handlerPost(w http.ResponseWriter, r *http.Request) {
	// リクエストからappengineのコンテキストを生成する
	c := appengine.NewContext(r)

	// Employee型の変数を定義する
	e1 := Employee{
		Name:     "Joe Citizen",
		Role:     "Manager",
		HireDate: time.Now(),
		Account:  "kim",
	}

	// DataStoreのキーを生成する
	key := datastore.NewKey(c, "employee", "abc", 0, nil)
	// 設定したキーを元に、put処理を行う
	_, _ = datastore.Put(c, key, &e1)
}

// GOONを使わずにget処理を行う
func handlerGet(w http.ResponseWriter, r *http.Request) {
	// リクエストからappengineのコンテキストを生成する
	c := appengine.NewContext(r)

	// DataStoreのキーを生成する
	employeeKey := datastore.NewKey(c, "employee", "abc", 0, nil)

	// Employee型の変数を生成する
	var employee Employee

	// 設定したキーを元に、get処理を行う
	err := datastore.Get(c, employeeKey, &employee)

	// 取得したデータを出力する
	fmt.Fprintln(w, fmt.Sprintf("emp: %v", employee))

	// デバッグのログを出力する
	log.Debugf(c, "emp: %#v", employee)

	// データを取得できなかった場合、メッセージを出力する
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
