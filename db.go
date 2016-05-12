package main

import "gopkg.in/mgo.v2"

func NewConnection() *mgo.Session {
	session, err := mgo.Dial("192.168.125.60")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}
