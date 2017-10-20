package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"execScript", "GET", "/pixfmp/{database}/{layout}/script/{script}/", execScript})
}

func execScript(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, password, _ := r.BasicAuth()
	param := strings.TrimSpace(r.URL.RawQuery)
	fm := filemaker.NewServer(fmServer, username, password)
	req, found, err := fm.RunScript(vars["database"], vars["layout"], vars["script"], param)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("found-count", fmt.Sprintf("%d", found))
	returnJson(w, req, vars["database"], vars["layout"], r)
}
