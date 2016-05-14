package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s\t%s\t%s\t%d\t%d\t%d",
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
		r.Proto,
		http.StatusNotFound,
		0,
		0,
	)
	w.WriteHeader(http.StatusNotFound)
}

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/todos", TodoIndex).Methods("GET")
	r.HandleFunc("/api/todos/{todoId}", TodoShow).Methods("GET")
	r.HandleFunc("/api/todos", TodoAdd).Methods("POST")
	r.HandleFunc("/api/todos/{todoId}", TodoUpdate).Methods("PUT")
	r.HandleFunc("/api/todos/{todoId}", TodoDelete).Methods("DELETE")
	r.HandleFunc("/api/todos/search/{todoName}", TodoSearch).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(NotFound)
	return r
}
