package main

import (
	"log"
	"net/http"

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
