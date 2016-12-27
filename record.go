package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"getRecordHandler", "GET", "/pixfmp/{database}/{layout}/record/{recid}/", getRecordHandler})
}

func getRecordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, password, _ := r.BasicAuth()
	fm := filemaker.NewServer(fmServer, username, password)
	req, err := fm.GetRow(vars["database"], vars["layout"], vars["recid"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 404)
		return
	}
	setUrl(req, r, vars["database"], vars["layout"])

	returnJson(w, req)
}
