package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//Session Establish the main session, this comes from db.go
var Session = NewConnection()

var null = "null"

//JSONResponse function to help in responses
func JSONResponse(w http.ResponseWriter, r *http.Request, start time.Time, response []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	log.Printf("%s\t%s\t%s\t%s\t%d\t%d\t%s",
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
		r.Proto,
		code,
		len(response),
		time.Since(start),
	)
	if string(response) != "" {
		w.Write(response)
	}
}

//JSONError function to help in error responses
func JSONError(w http.ResponseWriter, r *http.Request, start time.Time, message string, code int) {
	j := map[string]string{"message": message}
	response, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	log.Printf("%s\t%s\t%s\t%s\t%d\t%d\t%s",
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
		r.Proto,
		code,
		len(response),
		time.Since(start),
	)
	w.Write(response)
}

//TodoIndex handler to route index
func TodoIndex(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var todos []Todo
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	collection.Find(bson.M{}).All(&todos)
	response, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		panic(err)
	}
	JSONResponse(w, r, start, response, http.StatusOK)
}

//TodoShow handler to show all todos
func TodoShow(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var todo Todo
	vars := mux.Vars(r)
	if bson.IsObjectIdHex(vars["todoID"]) != true {
		JSONError(w, r, start, "bad entry for id", http.StatusBadRequest)
		return
	}
	todoID := bson.ObjectIdHex(vars["todoID"])
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	collection.Find(bson.M{"_id": todoID}).One(&todo)
	if todo.ID == "" {
		JSONError(w, r, start, "todo not found", http.StatusNotFound)
	} else {
		response, err := json.MarshalIndent(todo, "", "    ")
		if err != nil {
			panic(err)
		}
		JSONResponse(w, r, start, response, http.StatusOK)
	}
}

// TodoAdd handler to add new todo
func TodoAdd(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	if todo.Name == "" {
		JSONError(w, r, start, "Incorrect body", http.StatusBadRequest)
		return
	}
	objID := bson.NewObjectId()
	todo.ID = objID
	todo.Created = time.Now()
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Insert(todo)
	if err != nil {
		JSONError(w, r, start, "Failed to insert todo", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", r.URL.Path+"/"+string(todo.ID.Hex()))
	JSONResponse(w, r, start, []byte{}, http.StatusCreated)
}

//TodoUpdate handler to update a previous todo
func TodoUpdate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var todo Todo
	vars := mux.Vars(r)
	if bson.IsObjectIdHex(vars["todoID"]) != true {
		JSONError(w, r, start, "bad entry for id", http.StatusBadRequest)
		return
	}
	json.NewDecoder(r.Body).Decode(&todo)
	if todo.Name == "" {
		JSONError(w, r, start, "Incorrect body", http.StatusBadRequest)
		return
	}
	todoID := bson.ObjectIdHex(vars["todoID"])
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Update(bson.M{"_id": todoID},
		bson.M{"$set": bson.M{"name": todo.Name, "completed": todo.Completed}})
	if err != nil {
		JSONError(w, r, start, "Could not find Todo "+string(todoID.Hex())+" to update", http.StatusNotFound)
		return
	}
	JSONResponse(w, r, start, []byte{}, http.StatusNoContent)
}

//TodoDelete handler to delete a todo
func TodoDelete(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	vars := mux.Vars(r)
	todoID := bson.ObjectIdHex(vars["todoID"])
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Remove(bson.M{"_id": todoID})
	if err != nil {
		JSONError(w, r, start, "Could not find Todo "+string(todoID.Hex())+" to delete", http.StatusNotFound)
		return
	}
	JSONResponse(w, r, start, []byte{}, http.StatusNoContent)
}

//TodoSearchName handler Todo by Name
func TodoSearchName(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var todo []Todo
	vars := mux.Vars(r)
	todoName := vars["todoName"]
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	err := collection.Find(bson.M{"name": bson.M{"$regex": todoName}}).All(&todo)
	if err != nil {
		JSONError(w, r, start, "Failed to search todo name", http.StatusInternalServerError)
		return
	}
	response, err := json.MarshalIndent(todo, "", "    ")
	if err != nil {
		panic(err)
	}
	if string(response) == null {
		JSONError(w, r, start, "Could not find any Todo containing "+todoName, http.StatusNotFound)
		return
	}
	JSONResponse(w, r, start, response, http.StatusOK)
}

//TodoSearchStatus search todo by status (completed, not completed)
func TodoSearchStatus(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var todo []Todo
	vars := mux.Vars(r)
	todoStatus := vars["status"]
	session := Session.Copy()
	defer session.Close()
	collection := session.DB("prod").C("todos")
	if todoStatus == "true" {
		err := collection.Find(bson.M{"completed": bson.M{"$eq": true}}).All(&todo)
		if err != nil {
			JSONError(w, r, start, "Failed to search todo name", http.StatusInternalServerError)
			return
		}
		response, err := json.MarshalIndent(todo, "", "    ")
		if err != nil {
			panic(err)
		}
		if string(response) == null {
			JSONError(w, r, start, "Could not find any Todo containing status "+todoStatus, http.StatusNotFound)
			return
		}
		JSONResponse(w, r, start, response, http.StatusOK)
	} else if todoStatus == "false" {
		err := collection.Find(bson.M{"completed": bson.M{"$eq": false}}).All(&todo)
		if err != nil {
			JSONError(w, r, start, "Failed to search todo name", http.StatusInternalServerError)
			return
		}
		response, err := json.MarshalIndent(todo, "", "    ")
		if err != nil {
			panic(err)
		}
		if string(response) == null {
			JSONError(w, r, start, "Could not find any Todo containing status "+todoStatus, http.StatusNotFound)
			return
		}
		JSONResponse(w, r, start, response, http.StatusOK)
	} else {
		JSONError(w, r, start, "bad request, must be true or false, not "+todoStatus, http.StatusNotFound)
	}
}
