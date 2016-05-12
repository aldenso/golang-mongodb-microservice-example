package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//Establish the main session, this comes from db.go
var Session = NewConnection()

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome! to my first TODO API")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	collection.Find(bson.M{}).All(&todos)
	w.Header().Set("Content-Type", "application/json")
	j, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Write(j)
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	vars := mux.Vars(r)
	if bson.IsObjectIdHex(vars["todoId"]) != true {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad entry for id!"))
		return
	}
	todoId := bson.ObjectIdHex(vars["todoId"])
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	collection.Find(bson.M{"_id": todoId}).One(&todo)
	if todo.Id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("todo not found!"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		j, err := json.MarshalIndent(todo, "", "    ")
		if err != nil {
			panic(err)
		}
		w.Write(j)
	}
}

func TodoAdd(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	if todo.Name == "" || !todo.Completed {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect Body!"))
		return
	}
	obj_id := bson.NewObjectId()
	todo.Id = obj_id
	todo.Created = time.Now()
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Insert(todo)
	if err != nil {
		log.Println("Failed insert book: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.URL.Path+"/"+string(todo.Id))
	w.WriteHeader(http.StatusCreated)
}

func TodoUpdate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	vars := mux.Vars(r)
	if bson.IsObjectIdHex(vars["todoId"]) != true {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad entry for id!"))
		return
	}
	json.NewDecoder(r.Body).Decode(&todo)
	if todo.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect Body!"))
		return
	}
	todoId := bson.ObjectIdHex(vars["todoId"])
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Update(bson.M{"_id": todoId},
		bson.M{"$set": bson.M{"name": todo.Name, "completed": todo.Completed}})
	if err != nil {
		log.Printf("Could not find Todo %s to update", todoId)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func TodoDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := bson.ObjectIdHex(vars["todoId"])
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Remove(bson.M{"_id": todoId})
	if err != nil {
		log.Printf("Could not find Todo %s to delete", todoId)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
