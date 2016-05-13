package main

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/api/", Index).Methods("GET")
	r.HandleFunc("/api/todos", TodoIndex).Methods("GET")
	r.HandleFunc("/api/todos/{todoId}", TodoShow).Methods("GET")
	r.HandleFunc("/api/todos", TodoAdd).Methods("POST")
	r.HandleFunc("/api/todos/{todoId}", TodoUpdate).Methods("PUT")
	r.HandleFunc("/api/todos/{todoId}", TodoDelete).Methods("DELETE")
	r.HandleFunc("/api/todos/search/{todoName}", TodoSearch).Methods("GET")
	return r
}
