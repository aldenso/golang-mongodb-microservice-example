package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes { //this routes come from routes.go
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name) //related to logger.go

		router. //Check routes.go
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
