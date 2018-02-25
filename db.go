package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

// NewConnection create connection to DB
func NewConnection() *mgo.Session {
	MONGODB := os.Getenv("MONGODB_IP")
	if MONGODB == "" {
		log.Fatal("You need to export MONGODB_IP environment variable")
	}
	//session, err := mgo.Dial(MONGODB)
	session, err := mgo.DialWithTimeout(MONGODB, (60 * time.Second))
	if err != nil {
		log.Fatal(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}
