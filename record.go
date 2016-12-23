package main

import (
	"encoding/json"
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
	url, _ := router.Get("getRecordHandler").URL("database", vars["database"], "layout", vars["layout"], "recid", req["recid"])
	req["recid_url"] = r.Host + url.String()
	//fmt.Println(req)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(req)
}
