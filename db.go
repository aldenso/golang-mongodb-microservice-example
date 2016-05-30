package main

import "gopkg.in/mgo.v2"

// NewConnection create connection to DB
func NewConnection() *mgo.Session {
	session, err := mgo.Dial("192.168.125.60")
	//session, err := mgo.Dial("172.17.0.1")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}
