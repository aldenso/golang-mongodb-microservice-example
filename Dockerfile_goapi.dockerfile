# Build binary with the following command
# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o server .
FROM alpine

WORKDIR /src

COPY server /src/

ENTRYPOINT ["./server"]