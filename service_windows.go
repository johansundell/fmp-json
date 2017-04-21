// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"

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

var fmServer, redirectPort string

var router *mux.Router
var mydebug, displayDatabases, displayLayouts bool

const appVersionStr = "0.5"

type key int

type myservice struct{}

var elog debug.Log

type Settings struct {
	HttpInterface string `json:"http-interface"`
	HttpPort      string `json:"http-port"`
	FmServer      string `json:"fm-server"`
	Debug         bool   `json:"debug"`
}

func (s *Settings) isSane() bool {
	if s.FmServer == "" || s.HttpPort == "" {
		return false
	}
	return true
}

func (m *myservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	ex, err := os.Executable()
	if err != nil {
		elog.Info(1, err.Error())
		changes <- svc.Status{State: svc.StopPending}
		return
	}
	//
	settings := Settings{"", ":8081", "", false}
	dir, _ := filepath.Split(ex)
	filename := "settings.json"
	dat, err := ioutil.ReadFile(dir + filename)
	if err != nil {
		data, _ := json.Marshal(settings)
		ioutil.WriteFile(dir+filename, data, 0664)
		elog.Info(1, "Settings missing"+err.Error())
		changes <- svc.Status{State: svc.StopPending}
		return
	}
	if err := json.Unmarshal(dat, &settings); err != nil {
		elog.Info(1, err.Error())
		changes <- svc.Status{State: svc.StopPending}
		return
	}
	if !settings.isSane() {
		elog.Info(1, "Settings is not good")
		changes <- svc.Status{State: svc.StopPending}
		return
	}
	//elog.Info(1, "test "+dir)
	//changes <- svc.Status{State: svc.StopPending}
	//return

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	fasttick := time.Tick(500 * time.Millisecond)
	slowtick := time.Tick(2 * time.Second)
	tick := fasttick
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	var httpPortInterface string
	//var useSyslog bool
	httpPortInterface = os.Getenv("PIXFMT_HTTP_PORT")
	//tlsPortInterface = os.Getenv("PIXFMP_TLS_PORT")
	fmServer = os.Getenv("PIXFMP_FM_SERVER")
	//sslCert = os.Getenv("PIXFMP_TLS_CERT")
	//sslKey = os.Getenv("PIXFMP_TLS_KEY")
	redirectPort = os.Getenv("PIXFMP_HTTP_REDIRECT_TO")
	mydebug, _ = strconv.ParseBool(os.Getenv("PIXFMP_DEBUG"))

	// TEST
	httpPortInterface = settings.HttpPort
	fmServer = settings.FmServer

	router = NewRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    httpPortInterface,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 90 * time.Second,
		ReadTimeout:  90 * time.Second,
	}
	go func() {
		srv.ListenAndServe()
	}()
	/*
		srv.Addr = httpPortInterface
		log.Fatal(srv.ListenAndServe())
	*/

loop:
	for {
		select {
		case <-tick:
			//beep()
			//elog.Info(1, "beep")
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				elog.Info(1, "fmpjson shutdown ;)")
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				tick = slowtick
				srv.Close()
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				tick = fasttick
				go func() {
					srv.ListenAndServe()
				}()
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	srv.Close()
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &myservice{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}

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
