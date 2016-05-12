golang-mongodb-microservice-example
===================================

This is a small example for a restful service using golang, mongodb and gorilla

Still in progress, needs a lot of improvements...

My lab has and app server and a mongodb server

mongodb_Address: "192.168.125.60:27017"

app_Address: "192.168.125.1:8080"

List todos:

	# curl -i http://192.168.125.1:8080/todos

Show a single todo (replace {id} for the equivalent bson.ObjectIdHex):

	# curl -i http://192.168.125.1:8080/todos/{id}

Add todo:

	# curl -i http://192.168.125.1:8080/todos -X POST -d @add.json

where add.json file is something like:

	{
		"name":   "Task 14",
		"completed":   false
	}

Update todo (replace {id} for the equivalent bson.ObjectIdHex):

	# curl -i http://192.168.125.1:8080/todos/{id} -X PUT -d @add.json

where add.json file is something like:

	{
		"name":   "Task X",
		"completed":   true
	}

Delete todo:

	# curl -i http://192.168.125.1:8080/todos/{id} -X DELETE
