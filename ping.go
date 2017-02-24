package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"pingServer", "GET", "/pixfmp/server/ping/", pingServer})
}

func pingServer(w http.ResponseWriter, r *http.Request) {
	// UGLY copy ;)
	username, password, _ := r.BasicAuth()
	fm := filemaker.NewServer(fmServer, username, password)
	req, err := fm.ListDatabases()
	ping := false
	if err != nil || len(req) == 0 {
		log.Println(err)
		ping = false
	} else {
		ping = true
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if debug {
		enc.SetIndent("", "\t")
	}
	enc.Encode(ping)
	return
	/*output := make([]map[string]string, 0)
	for k, v := range req {
		log.Println(k, v["DATABASE_NAME"])
		row := make(map[string]string)
		row["database"] = v["DATABASE_NAME"].String()
		if url, err := router.GetRoute("listLayouts").URL("database", v["DATABASE_NAME"].String()); err == nil {
			row["layouts"] = r.Host + url.String()
		}
		output = append(output, row)
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if debug {
		enc.SetIndent("", "\t")
	}
	enc.Encode(output)
	*/
}
