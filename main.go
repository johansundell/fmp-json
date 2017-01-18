package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		defaultHandler,
	},
}

var fmServer string

var router *mux.Router
var debug bool

const appVersionStr = "0.2"

type key int

const (
	serverKey = iota
	debugKey
)

func main() {
	var httpPortInterface, tlsPortInterface, sslKey, sslCert string
	var useSyslog bool
	flag.StringVar(&httpPortInterface, "http", ":8080", "HTTP port and interface the server will use, format interface:port")
	flag.StringVar(&tlsPortInterface, "tls", ":1443", "TLS port and interface the server will use, format interface:port")
	flag.StringVar(&fmServer, "server", "http://fmh-iwp12.no-ip.info", "The filemaker server to use as host")
	flag.BoolVar(&debug, "debug", true, "Debug requests")
	flag.BoolVar(&useSyslog, "usesyslog", false, "Use syslog")
	flag.StringVar(&sslCert, "ssl-cert", "server.crt", "Path to the ssl cert to use, if empty it will use http")
	flag.StringVar(&sslKey, "ssl-key", "server.key", "Path to the ssl key to use, if empty it will use http")
	flag.Parse()

	if useSyslog {
		logwriter, err := syslog.New(syslog.LOG_NOTICE, "fmp-json")
		if err != nil {
			panic(err)
		}
		log.SetOutput(logwriter)
	}

	router = NewRouter()
	srv := &http.Server{
		Handler: router,
		//Addr:    httpPortInterface,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 90 * time.Second,
		ReadTimeout:  90 * time.Second,
	}
	useTls := false
	if sslCert != "" && sslKey != "" && simpleFileExist(sslCert) && simpleFileExist(sslKey) {
		useTls = true
		srv.Addr = tlsPortInterface
		go func() {
			log.Fatal(srv.ListenAndServeTLS(sslCert, sslKey))
		}()
	}
	//log.Println("Could not find key or cert for ssl", sslKey, sslCert)
	if useTls {
		log.Fatal(http.ListenAndServe(httpPortInterface, http.HandlerFunc(redir)))
	} else {
		srv.Addr = httpPortInterface
		log.Fatal(srv.ListenAndServe())
	}
}

func returnJson(w http.ResponseWriter, data interface{}, database, layout string, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if debug {
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
		default:
			output[k] = v.String()
		}
	}
	output["recid_url"] = getUrl(r, req, database, layout)
	return output
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "TODO: Write some documentation here")
}

func getUrl(req filemaker.Record, r *http.Request, database, layout string) string {
	url, err := router.Get("getRecordHandler").URL("database", database, "layout", layout, "recid", req["recid"].String())
	if err != nil {
		return ""
	}
	return r.Host + url.String()
}

func redir(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func simpleFileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
