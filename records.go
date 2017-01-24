package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"getRecordsHandler", "GET", "/pixfmp/{database}/{layout}/records/", getRecordsHandler})
	routes = append(routes, Route{"getRecordsPageHandler", "GET", "/pixfmp/{database}/{layout}/records/{start}/{stop}/", getRecordsPageHandler})
}

func getRecordsPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	start, err := strconv.Atoi(vars["start"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 501)
		return
	}
	stop, err := strconv.Atoi(vars["stop"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 501)
		return
	}
	getRecords(w, r, start, stop)
}

func getRecordsHandler(w http.ResponseWriter, r *http.Request) {
	getRecords(w, r, -1, -1)
}

func getRecords(w http.ResponseWriter, r *http.Request, start, stop int) {
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

	fm := filemaker.NewServer(fmServer, username, password)
	req, err := fm.Get(vars["database"], vars["layout"], params, start, stop)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 404)
		return
	}
	returnJson(w, req, vars["database"], vars["layout"], r)
}
