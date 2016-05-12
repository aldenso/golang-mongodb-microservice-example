package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		TodoIndex,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		TodoShow,
	},
	Route{
		"TodoDelete",
		"DELETE",
		"/todos/{todoId}",
		TodoDelete,
	},
	Route{
		"TodoAdd",
		"POST",
		"/todos",
		TodoAdd,
	},
	Route{
		"TodoUpdate",
		"PUT",
		"/todos/{todoId}",
		TodoUpdate,
	},
}
