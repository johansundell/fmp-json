package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"putRecordHandler", "PUT", "/pixfmp/{database}/{layout}/record/{recid}/", putRecordHandler})
	routes = append(routes, Route{"replaceRecordHandler", "PUT", "/pixfmp/{database}/{layout}/records/{mainkey}/{mainvalue}/", replaceRecordHandler})
}

func replaceRecordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	useForce := false
	if r.URL.Query().Get("force") == "true" {
		useForce = true
	}
	username, password, _ := r.BasicAuth()
	mainkey := vars["mainkey"]
	mainvalue := vars["mainvalue"]
	data, err := getRequestData(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 404)
		return
	}
	log.Println(mainkey, mainvalue, data)

	params := make([]filemaker.SearchParam, 0)
	params = append(params, filemaker.SearchParam{Name: mainkey, Value: mainvalue, Op: filemaker.Equal})
	fm := filemaker.NewServer(fmServer, username, password)
	req, _, err := fm.Get(vars["database"], vars["layout"], params, -1, -1)
	if err != nil {
		if !useForce {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			// Not found so lets create it
			data[mainkey] = mainvalue
			record, err := fm.NewRow(vars["database"], vars["layout"], data)
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(404), 404)
				return
			}
			returnJson(w, record, vars["database"], vars["layout"], r)
		}
		return
	}
	log.Println(req)
	records := make(filemaker.Records, 0)
	for _, row := range req {
		reqid := row["recid"].String()
		_ = reqid
		/*delete(row, "recid")
		log.Println(row)*/
		result, err := fm.EditRow(vars["database"], vars["layout"], reqid, data)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		records = append(records, result)
	}
	returnJson(w, records, vars["database"], vars["layout"], r)
	//log.Println(vars)
	//http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

func putRecordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data, err := getRequestData(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 404)
		return
	}
	//log.Println(data)
	username, password, _ := r.BasicAuth()
	fm := filemaker.NewServer(fmServer, username, password)
	req, err := fm.EditRow(vars["database"], vars["layout"], vars["recid"], data)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	returnJson(w, req, vars["database"], vars["layout"], r)
}
