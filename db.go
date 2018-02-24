package main

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

// NewConnection create connection to DB
func NewConnection() *mgo.Session {
	MONGODB := os.Getenv("MONGODB_IP")
	if MONGODB == "" {
		log.Fatal("You need to export MONGODB_IP environment variable")
	}
	session, err := mgo.Dial(MONGODB)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}
