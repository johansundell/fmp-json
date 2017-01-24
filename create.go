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
	data, err := getRequestData(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 404)
		return
	}
	//log.Println(data)
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
