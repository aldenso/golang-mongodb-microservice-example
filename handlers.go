package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//Establish the main session, this comes from db.go
var Session = NewConnection()

// function to help in responses
func JsonResponse(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

// function to help in error responses
func JsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome! to my first TODO API")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	collection.Find(bson.M{}).All(&todos)
	response, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		panic(err)
	}
	JsonResponse(w, response, http.StatusOK)
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	vars := mux.Vars(r)
	if bson.IsObjectIdHex(vars["todoId"]) != true {
		JsonError(w, "bad entry for id", http.StatusBadRequest)
		return
	}
	todoId := bson.ObjectIdHex(vars["todoId"])
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	collection.Find(bson.M{"_id": todoId}).One(&todo)
	if todo.Id == "" {
		JsonError(w, "todo not found", http.StatusNotFound)
	} else {
		response, err := json.MarshalIndent(todo, "", "    ")
		if err != nil {
			panic(err)
		}
		JsonResponse(w, response, http.StatusOK)
	}
}

func TodoAdd(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	if todo.Name == "" {
		JsonError(w, "Incorrect body", http.StatusBadRequest)
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
		JsonError(w, "Failed to insert todo", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.URL.Path+"/"+string(todo.Id.Hex()))
	w.WriteHeader(http.StatusCreated)
}

func TodoUpdate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	vars := mux.Vars(r)
	if bson.IsObjectIdHex(vars["todoId"]) != true {
		JsonError(w, "bad entry for id", http.StatusBadRequest)
		return
	}
	json.NewDecoder(r.Body).Decode(&todo)
	if todo.Name == "" {
		JsonError(w, "Incorrect body", http.StatusBadRequest)
		return
	}
	todoId := bson.ObjectIdHex(vars["todoId"])
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Update(bson.M{"_id": todoId},
		bson.M{"$set": bson.M{"name": todo.Name, "completed": todo.Completed}})
	if err != nil {
		JsonError(w, "Could not find Todo "+string(todoId.Hex())+" to update", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func TodoDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := bson.ObjectIdHex(vars["todoId"])
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Remove(bson.M{"_id": todoId})
	if err != nil {
		JsonError(w, "Could not find Todo "+string(todoId.Hex())+" to delete", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// Search Todo by Name
func TodoSearch(w http.ResponseWriter, r *http.Request) {
	var todo []Todo
	vars := mux.Vars(r)
	todoName := vars["todoName"]
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Find(bson.M{"name": bson.M{"$regex": todoName}}).All(&todo)
	if err != nil {
		JsonError(w, "Failed to search todo name", http.StatusInternalServerError)
		return
	}
	response, err := json.MarshalIndent(todo, "", "    ")
	if err != nil {
		panic(err)
	}
	if string(response) == "null" {
		JsonError(w, "Could not find any Todo containing "+todoName, http.StatusNotFound)
		return
	}
	JsonResponse(w, response, http.StatusOK)
}
