package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter() // this func is in router.go
	defer Session.Close() // related to Session in handlers.go
	log.Fatal(http.ListenAndServe(":8080", router))
}
