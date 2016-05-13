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
		"/api",
		Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/api/todos",
		TodoIndex,
	},
	Route{
		"TodoShow",
		"GET",
		"/api/todos/{todoId}",
		TodoShow,
	},
	Route{
		"TodoDelete",
		"DELETE",
		"/api/todos/{todoId}",
		TodoDelete,
	},
	Route{
		"TodoAdd",
		"POST",
		"/api/todos",
		TodoAdd,
	},
	Route{
		"TodoUpdate",
		"PUT",
		"/api/todos/{todoId}",
		TodoUpdate,
	},
	Route{
		"TodoSearch",
		"GET",
		"/api/todos/search/{todoName:[A-Za-z0-9_]+}",
		TodoSearch,
	},
}
