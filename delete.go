package main

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(true)
}
