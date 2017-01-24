package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func logger(inner http.Handler, name string) http.Handler {
	//name := runtime.FuncForPC(reflect.ValueOf(inner).Pointer()).Name()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if debug {
			formValues := ""
			if r.Method == "PUT" || r.Method == "POST" {
				r.ParseForm()
				for k, v := range r.Form {
					if len(v) != 0 {
						formValues += " " + k + "=" + v[0]
					}
				}
			}
			log.Println(name, r.RequestURI, r.RemoteAddr, r.Method, formValues)
		}
		w.Header().Set("X-fmp-json-Version", appVersionStr)
		//ctx := context.WithValue(r.Context(), serverKey, fmServer)

		inner.ServeHTTP(w, r)
	})
}

func getRequestData(w http.ResponseWriter, r *http.Request) (map[string]string, error) {
	if len(r.Header["Content-Type"]) == 0 {
		return nil, errors.New("Content-type not set")
	}
	data := make(map[string]string)
	switch r.Header["Content-Type"][0] {
	case "application/x-www-form-urlencoded":
		r.ParseForm()
		for k, v := range r.Form {
			if len(v) >= 1 {
				data[k] = v[0]
			}
		}
	case "application/json":
		var tmp map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&tmp); err != nil {
			return nil, errors.New("Not proper json")
		}
		for k, v := range tmp {
			switch t := v.(type) {
			case float64:
				data[k] = strconv.FormatFloat(t, 'f', -1, 64)
			case string:
				data[k] = t
			default:
				log.Printf("Err %T\n", v)
			}

		}
	}
	return data, nil
}
