package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
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

var httpPortInterface, tlsPortInterface, fmServer, sslKey, sslCert string
var router *mux.Router
var debug, useSyslog bool

const appVersionStr = "0.1"

func main() {
	flag.StringVar(&httpPortInterface, "http", ":8080", "HTTP port and interface the server will use, format interface:port")
	flag.StringVar(&tlsPortInterface, "tls", ":1443", "TLS port and interface the server will use, format interface:port")
	flag.StringVar(&fmServer, "server", "http://fmh-iwp12.no-ip.info", "The filemaker server to use as host")
	flag.BoolVar(&debug, "debug", false, "Debug requests")
	flag.BoolVar(&useSyslog, "usesyslog", false, "Use syslog")
	flag.StringVar(&sslCert, "ssl-cert", "server.crt", "Path to the ssl cert to use, if empty it will use http")
	flag.StringVar(&sslKey, "ssl-key", "server.key", "Path to the ssl key to use, if empty it will use http")
	flag.Parse()

	if useSyslog {
		logwriter, err := syslog.New(syslog.LOG_NOTICE, "pixext")
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

func returnJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if debug {
		enc.SetIndent("", "\t")
	}
	enc.Encode(data)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TODO: Write some documentation here")
}

func setUrl(req filemaker.Record, r *http.Request, database, layout string) {
	url, _ := router.Get("getRecordHandler").URL("database", database, "layout", layout, "recid", req["recid"])
	req["recid_url"] = r.Host + url.String()
}

func redir(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
	//http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusTemporaryRedirect)
}

func simpleFileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
