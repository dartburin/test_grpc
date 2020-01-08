# test_grpc
Test project with gRPC and REST proxy.


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
curl -XPOST "http://127.0.0.1:8080/books" -d '{"Author": "Some author", "Title": "New book 1"}'

(DELETE)
curl -XDELETE "http://127.0.0.1:8080/books/91"

(PATCH)
curl -XPATCH "http://127.0.0.1:8080/books/43" -d '{"Author": "A100", "Title": "Book 1"}'

(PUT)
curl -XPUT "http://127.0.0.1:8080/books/43" -d '{"Author": "A1", "Title": "Book 1"}'
