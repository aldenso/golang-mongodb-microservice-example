package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Todo struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `json:"name"`
	Completed bool          `json:"completed"`
	Created   time.Time     `json:"createdon"`
}
