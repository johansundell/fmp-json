// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main

import (
	"flag"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
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

var fmServer, redirectPort string

var router *mux.Router
var mydebug, displayDatabases, displayLayouts bool

type key int

func main() {
	var httpPortInterface, tlsPortInterface, sslKey, sslCert string
	var useSyslog bool
	httpPortInterface = os.Getenv("PIXFMT_HTTP_PORT")
	tlsPortInterface = os.Getenv("PIXFMP_TLS_PORT")
	fmServer = os.Getenv("PIXFMP_FM_SERVER")
	sslCert = os.Getenv("PIXFMP_TLS_CERT")
	sslKey = os.Getenv("PIXFMP_TLS_KEY")
	redirectPort = os.Getenv("PIXFMP_HTTP_REDIRECT_TO")
	mydebug, _ = strconv.ParseBool(os.Getenv("PIXFMP_DEBUG"))
	flag.StringVar(&httpPortInterface, "http", httpPortInterface, "HTTP port and interface the server will use, format interface:port")
	flag.StringVar(&tlsPortInterface, "tls", tlsPortInterface, "TLS port and interface the server will use, format interface:port")
	flag.StringVar(&redirectPort, "redirect-to", tlsPortInterface, "When using TLS, redirect all request using http to this port")
	flag.StringVar(&fmServer, "server", fmServer, "The filemaker server to use as host")
	flag.BoolVar(&mydebug, "debug", mydebug, "Debug requests")
	flag.BoolVar(&useSyslog, "usesyslog", false, "Use syslog")
	flag.StringVar(&sslCert, "ssl-cert", sslCert, "Path to the ssl cert to use, if empty it will use http")
	flag.StringVar(&sslKey, "ssl-key", sslKey, "Path to the ssl key to use, if empty it will use http")
	flag.BoolVar(&displayDatabases, "list-databases", true, "Display all XML enabled databases")
	flag.BoolVar(&displayLayouts, "list-layouts", true, "Display all XML enabled layouts")
	flag.Parse()

	if fmServer == "" {
		log.Fatalln("Filemaker server not set, exiting...")
	}

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
	/*router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})*/
	//log.Println("Could not find key or cert for ssl", sslKey, sslCert)
	if useTls {
		log.Fatal(http.ListenAndServe(httpPortInterface, http.HandlerFunc(redir)))
	} else {
		srv.Addr = httpPortInterface
		log.Fatal(srv.ListenAndServe())
	}
}
