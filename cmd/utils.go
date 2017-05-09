package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// JSON writes out pretty print JSON
func JSON(w http.ResponseWriter, v interface{}, s string, c int) {
	b, err := JSONIndent(v, s)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(c)
	w.Write(b)
}

// JSONError pretty prints any errors
func JSONError(w http.ResponseWriter, v interface{}, c int) {
	b, _ := json.MarshalIndent(v, "", "    ")
	w.WriteHeader(c)
	w.Write(b)
}

// JSONIndent pretty prints any JSON information
func JSONIndent(v interface{}, s string) (rj []byte, err error) {
	if len(s) != 0 && s == "false" {
		rj, err := json.Marshal(v)
		return rj, err
	}
	rj, err = json.MarshalIndent(v, "", "    ")
	return rj, err
}

func getURLVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func getFields(r *http.Request, f string) string {
	query := r.URL.Query()
	return query.Get(f)
}

func getVar(r *http.Request, v string) string {
	vars := getURLVars(r)
	return vars[v]
}

func getPrettyPrintValue(r *http.Request) string {
	return getFields(r, "pretty")
}

func getAccNumVar(r *http.Request) string {
	return getVar(r, "accNum")
}
