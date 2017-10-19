package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"listLayouts", "GET", "/pixfmp/{database}/", listLayouts})
}

func listLayouts(w http.ResponseWriter, r *http.Request) {
	if !displayLayouts {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	vars := mux.Vars(r)
	username, password, _ := r.BasicAuth()
	fm := filemaker.NewServer(fmServer, username, password)
	req, _, err := fm.Listlayouts(vars["database"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	log.Println(req)
	output := make([]map[string]string, 0)
	for k, v := range req {
		log.Println(k, v["LAYOUT_NAME"])
		if v["LAYOUT_NAME"].String() != "" {
			row := make(map[string]string)
			row["database"] = v["LAYOUT_NAME"].String()
			if url, err := router.GetRoute("getRecordsHandler").URL("database", vars["database"], "layout", v["LAYOUT_NAME"].String()); err == nil {
				row["records"] = r.Host + url.String()
			}
			output = append(output, row)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if mydebug {
		enc.SetIndent("", "\t")
	}
	enc.Encode(output)
}
