# test_grpc
Test project with gRPC and REST HTTP gateway.


Protobuf for gRPC compilation
-----------------------------
protoc -I. -I$GOPATH/src/github.com/dartburin --go_out=plugins=grpc:. api/proto/books.proto
protoc -I. -I$GOPATH/src/github.com/dartburin --grpc-gateway_out=logtostderr=true:. api/proto/books.proto

Start
-------
docker-compose up --build


Sample of testing orders
------------------------

(GET)
http://127.0.0.1:8080/books

(POST)
curl -XPOST "http://127.0.0.1:8080/books" -d '{"author": "Some author", "title": "New book 1"}'

(DELETE)
curl -XDELETE "http://127.0.0.1:8080/books/91"

(PATCH) partly update
curl -XPATCH "http://127.0.0.1:8080/books/43" -d '{"author": "A100", "title": "Book 1"}'

(PUT) full update
curl -XPUT "http://127.0.0.1:8080/books/43" -d '{"author": "A1", "title": "Book 1"}'
