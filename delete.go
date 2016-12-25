package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"deleteRecordHandler", "DELETE", "/pixfmp/{database}/{layout}/record/{recid}/", deleteRecordHandler})
}

func deleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, password, _ := r.BasicAuth()
	fm := filemaker.NewServer(fmServer, username, password)
	err := fm.DeleteRow(vars["database"], vars["layout"], vars["recid"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	returnJson(w, true)
}
