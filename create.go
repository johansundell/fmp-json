package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"postRecordHandler", "POST", "/pixfmp/{database}/{layout}/record/", postRecordHandler})
}

func postRecordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	data := make(map[string]string)
	for k, v := range r.Form {
		if len(v) != 0 {
			data[k] = v[0]
		}
	}
	username, password, _ := r.BasicAuth()
	fm := filemaker.NewServer(fmServer, username, password)
	req, err := fm.NewRow(vars["database"], vars["layout"], data)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	returnJson(w, req, vars["database"], vars["layout"], r)
}
