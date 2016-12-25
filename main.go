package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		defaultHandler,
	},
}

var hostname, basePath string
var router *mux.Router
var fmServer = "http://fmh-iwp12.no-ip.info"

func main() {
	flag.StringVar(&hostname, "hostname", ":8080", "The hostname that will server the files")
	flag.StringVar(&fmServer, "server", "http://fmh-iwp12.no-ip.info", "The filemaker server to use as host")
	flag.Parse()

	router = NewRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    hostname,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 90 * time.Second,
		ReadTimeout:  90 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func returnJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(data)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TODO: Write some documentation here")
}
