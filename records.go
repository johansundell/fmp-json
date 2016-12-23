package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"getRecordsHandler", "GET", "/pixfmp/{database}/{layout}/records/", getRecordsHandler})
}

func getRecordsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, password, _ := r.BasicAuth()
	params := make([]filemaker.SearchParam, 0)
	if strings.Contains(r.RequestURI, "?") {
		part := r.RequestURI[strings.Index(r.RequestURI, "?")+1:]
		var err error
		part, err = url.QueryUnescape(part)
		if err != nil {
			log.Println(err)
		}
		parts := strings.FieldsFunc(part, func(r rune) bool { return r == '&' })
		sp := filemaker.SearchParam{}
		for _, v := range parts {
			switch {
			case strings.Contains(v, "="):
				sp.Op = filemaker.Equal
			case strings.Contains(v, ">"):
				sp.Op = filemaker.MoreThan
			case strings.Contains(v, "<"):
				sp.Op = filemaker.LessThan
			}
			vals := strings.FieldsFunc(v, func(r rune) bool { return r == '=' || r == '>' || r == '<' })
			if len(vals) == 2 {
				sp.Name = vals[0]
				sp.Value = vals[1]
				params = append(params, sp)
			}
		}

	}
	//fmt.Println(params)
	//fmt.Println(vars)
	//fmt.Println(username, password)
	fm := filemaker.NewServer(fmServer, username, password)
	req, err := fm.Get(vars["database"], vars["layout"], params)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 404)
		return
	}
	for _, v := range req {
		url, _ := router.Get("getRecordHandler").URL("database", vars["database"], "layout", vars["layout"], "recid", v["recid"])
		v["recid_url"] = r.Host + url.String()
	}
	//fmt.Println(req)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(req)
}
