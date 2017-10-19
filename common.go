package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/johansundell/fmp-json/filemaker"
)

func returnJson(w http.ResponseWriter, data interface{}, database, layout string, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if mydebug {
		enc.SetIndent("", "\t")
	}
	switch r := data.(type) {
	case filemaker.Record:
		output := getFormatedData(r, database, layout, req)
		enc.Encode(output)
		return
	case filemaker.Records:
		output := make([]map[string]interface{}, 0)
		for _, v := range r {
			output = append(output, getFormatedData(v, database, layout, req))
		}
		enc.Encode(output)
		return
	default:
		log.Println("Error in encoding, unknown type", r)
	}
}

func getFormatedData(r filemaker.Record, database, layout string, req *http.Request) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range r {
		switch v.Type {
		case filemaker.FileMakerNumber:
			output[k], _ = strconv.ParseFloat(v.String(), 10)
		case filemaker.FileMakerDate, filemaker.FileMakerTimestamp:
			if v.String() == "" {
				output[k] = time.Time{}
			} else if date, err := time.Parse(v.Format, v.String()); err == nil {
				output[k] = date
			} else {
				log.Println(err, "recid", r["recid"])
				output[k] = time.Time{}
			}
		case filemaker.FileMakerContainer:
			if url, err := router.Get("getContainerHandler").URL("database", database, "layout", layout, "recid", r["recid"].String(), "field", k); err == nil {
				output[k] = req.Host + url.String()
			} else {
				log.Println(err)
			}
		default:
			output[k] = v.String()
		}
	}
	output["recid_url"] = getUrl(r, req, database, layout)
	return output
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fmp-json running ;)")
}

func getUrl(req filemaker.Record, r *http.Request, database, layout string) string {
	url, err := router.Get("getRecordHandler").URL("database", database, "layout", layout, "recid", req["recid"].String())
	if err != nil {
		return ""
	}
	return r.Host + url.String()
}

func redir(w http.ResponseWriter, r *http.Request) {
	if redirectPort != "" {
		host := r.Host
		if strings.Contains(host, ":") {
			host = host[:strings.Index(r.Host, ":")]
		}
		http.Redirect(w, r, "https://"+host+redirectPort+r.RequestURI, http.StatusMovedPermanently)
		return
	}
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func simpleFileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
