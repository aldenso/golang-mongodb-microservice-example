golang-mongodb-microservice-example
===================================

This is a small example for a restful service using golang, mongodb and gorilla

My lab consist in only one app server and one mongodb server

Mongodb_Address: "192.168.125.60:27017"

APP_Address: "192.168.125.1:8080"

List todos:

	# curl -i http://192.168.125.1:8080/api/todos

Show a single todo (replace {id} for the equivalent bson.ObjectIdHex):

	# curl -i http://192.168.125.1:8080/api/todos/{id}

Add todo:

	# curl -i http://192.168.125.1:8080/api/todos -X POST -d @add.json

where add.json file is something like:

	{
		"name":   "Task 14",
		"completed":   false
	}

Update todo (replace {id} for the equivalent bson.ObjectIdHex):

	# curl -i http://192.168.125.1:8080/api/todos/{id} -X PUT -d @update.json

where update.json file is something like:

	{
		"name":   "Task X",
		"completed":   true
	}

Delete todo:

	# curl -i http://192.168.125.1:8080/api/todos/{id} -X DELETE

Search todo (replace {name} for the equivalent search pattern):

	# curl -i http://192.168.125.1:8080/api/todos/search/{name}

Log examples:

	2016/05/13 02:01:39 192.168.125.1:15503        GET     /api/todos      HTTP/1.1        200     978     428.996µs
	2016/05/13 02:01:47 192.168.125.1:37320        GET     /api/todos/57351be5802abd1a283a9eca     HTTP/1.1        200     135     293.129µs
	2016/05/13 02:02:25 192.168.125.1:10409        GET     /api/todos/57351be5802abd1a283a9        HTTP/1.1        400     30      30.747µs
	2016/05/13 02:02:32 192.168.125.1:43340        GET     /api/todos/57351be5802abd1a283a9ecb     HTTP/1.1        404     28      543.703µs
	2016/05/13 02:03:58 192.168.125.1:61530        POST    /api/todos      HTTP/1.1        400     28      64.23µs
	2016/05/13 02:04:23 192.168.125.1:40092        POST    /api/todos      HTTP/1.1        201     0       595.388µs
	2016/05/13 02:04:49 192.168.125.1:13006        PUT     /api/todos/5735756fced374075efb7ef5     HTTP/1.1        204     0       1.204126ms
	2016/05/13 02:04:58 192.168.125.1:39523        PUT     /api/todos/5735756fced374075efb7        HTTP/1.1        400     30      21.706µs
	2016/05/13 02:05:21 192.168.125.1:17209        DELETE  /api/todos/5735756fced374075efb7ef5     HTTP/1.1        204     0       1.297254ms
	2016/05/13 02:05:25 192.168.125.1:39753        DELETE  /api/todos/5735756fced374075efb7ef5     HTTP/1.1        404     68      365.565µs
	2016/05/13 02:05:47 192.168.125.1:33098        GET     /api/todos/search/10    HTTP/1.1        200     166     706.37µs
	2016/05/13 02:05:53 192.168.125.1:43079        GET     /api/todos/search/a     HTTP/1.1        200     978     750.552µs
	2016/05/13 02:05:56 192.168.125.1:41361        GET     /api/todos/search/aq    HTTP/1.1        404     51      1.046498ms
