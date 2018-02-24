# golang-mongodb-microservice-example

This is a small example for a restful service using golang, mongodb and gorilla.

My lab consist in only one app server and one mongodb server:

    Mongodb container: alpine based with default port (27017).
    APP container:     alpine based with go program executable using port 8080.

First let's create the mongodb docker image.

```sh
docker build -f Dockerfile_mongo.dockerfile . --tag aldenso:mongodb-alpine
```

Run one container.

```sh
docker run -d --name mongodb -p 27017:27017 -v /vols4docker/mongodb:/data/db \
    aldenso:mongodb-alpine
```

Now let's build the app docker image.

Build the binary.

```sh
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o server .
```

Build the image.

```sh
docker build -f Dockerfile_goapi.dockerfile . --tag aldenso:goapi4mongo-alpine
```

export the mongodb connection info.

```sh
export MONGODB_IP="$(docker inspect -f '{{ .NetworkSettings.IPAddress }}' mongodb)"
```

Create the app container.

```sh
docker run -d --name goapi4mongo1 -p 8080:8080 -e MONGODB_IP=$MONGODB_IP \
    aldenso:goapi4mongo-alpine
```

Get the IP of the app container.

```sh
export APPSERVER_IP="$(docker inspect -f '{{ .NetworkSettings.IPAddress }}' goapi4mongo1)"
```

List todos:

```sh
curl -i http://$APPSERVER_IP:8080/api/todos
```

Show a single todo (replace {id} for the equivalent bson.ObjectIdHex):

```sh
curl -i http://$APPSERVER_IP:8080/api/todos/{id} # Replace the {id} with one created before.
```

Add todo:

```sh
curl -i -H "Content-Type: application/json" -X POST \
    -d '{"name": "Task_14", "completed": false}'  \
    http://$APPSERVER_IP:8080/api/todos
```

or

```sh
curl -i http://$APPSERVER_IP:8080/api/todos -X POST -d @add.json
```

where add.json file is something like:

```json
{
    "name":   "Task_14",
    "completed":   false
}
```

Update todo (replace {id} for the equivalent bson.ObjectIdHex):

```sh
curl -i -H "Content-Type: application/json" -X PUT \
    -d '{"name": "update_task", "completed": false}' \
    http://$APPSERVER_IP:8080/api/todos/{id}
```

or

```sh
curl -i http://$APPSERVER_IP:8080/api/todos/{id} -X PUT -d @update.json
```

where update.json file is something like:

```json
{
    "name":   "Task_X",
    "completed":   true
}
```

Delete todo:

```sh
curl -i http://$APPSERVER_IP:8080/api/todos/{id} -X DELETE
```

Search todo by name (replace {name} for the equivalent search pattern):

```sh
curl -i http://$APPSERVER_IP:8080/api/todos/search/byname/{name}
```

Search todo by status completed (replace {status} for true or false ):

```sh
curl -i http://$APPSERVER_IP:8080/api/todos/search/bystatus/{status}
```

Log samples:

```txt
2018/02/24 22:11:40 172.17.0.1:56552	GET	/api/todos	HTTP/1.1	200	4	1.621415ms
2018/02/24 22:11:53 172.17.0.1:56554	POST	/api/todos	HTTP/1.1	201	0	1.125878748s
2018/02/24 22:13:08 172.17.0.1:56562	GET	/api/todos/	HTTP/1.1	404	0	0
2018/02/24 22:13:10 172.17.0.1:56564	GET	/api/todos	HTTP/1.1	200	160	857.074µs
2018/02/24 22:13:44 172.17.0.1:56568	GET	/api/todos/5a91e328d7f8960001a99aeb	HTTP/1.1	200	132	1.702964ms
2018/02/24 22:14:40 172.17.0.1:56572	PUT	/api/todos/5a91e328d7f8960001a99aeb	HTTP/1.1	204	0	37.554697ms
2018/02/24 22:16:34 172.17.0.1:56580	GET	/api/todos	HTTP/1.1	200	159	898.811µs
2018/02/24 22:18:51 172.17.0.1:56592	GET	/api/todos/search/byname/Task	HTTP/1.1	200	159	901.338µs
2018/02/24 22:20:20 172.17.0.1:56600	GET	/api/todos/search/bystatus/false	HTTP/1.1	200	159	1.013693ms
2018/02/24 22:23:01 172.17.0.1:56618	POST	/api/todos	HTTP/1.1	201	0	907.716µs
2018/02/24 22:26:45 172.17.0.1:56634	GET	/api/todos	HTTP/1.1	200	314	992.902µs
2018/02/24 22:27:19 172.17.0.1:56640	DELETE	/api/todos/5a91e328d7f8960001a99aeb	HTTP/1.1	204	0	1.193428ms
2018/02/24 22:27:31 172.17.0.1:56642	DELETE	/api/todos/5a91e5c5d7f8960001a99aec	HTTP/1.1	204	0	3.093297ms
```

After playing with the api you can stop the containers.

```sh
docker stop goapi4mongo1 mongodb
```

Just for your information, check your new images and you'll confirm the small size of them.

```sh
docker images
```

```txt
REPOSITORY              TAG                  IMAGE ID            CREATED             SIZE
aldenso                 goapi4mongo-alpine   b9ba2da7dc62        About an hour ago   9.64MB
aldenso                 mongodb-alpine       0e53354d7bc1        About an hour ago   106MB
```
