# test_grpc
Test project with gRPC and REST HTTP gateway.


Test method
-----------

Begin

(docker terminal) docker-compose up --build
(browser)  http://localhost:8080/books
(test terminal) ./test/test.sh
(browser)  http://localhost:8080/books


View console output

(log terminal) docker logs -f gRPCserver


End

(docker terminal) Ctrl+z
(docker terminal) docker-compose down
(test terminal) sudo rm -R /tmp/data


Sample of testing orders
------------------------

(GET)
http://127.0.0.1:8080/books
curl -XGET "http://127.0.0.1:8080/books/3"

(POST)
curl -XPOST "http://127.0.0.1:8080/books" -d '{"author": "Some author", "title": "New book 1"}'

(DELETE)
curl -XDELETE "http://127.0.0.1:8080/books/91"

(PATCH) (partly update)
curl -XPATCH "http://127.0.0.1:8080/books/43" -d '{"author": "A100", "title": "Book 1"}'

(PUT) (full update)
curl -XPUT "http://127.0.0.1:8080/books/43" -d '{"author": "A1", "title": "Book 1"}'


Protobuf for gRPC compilation
-----------------------------
protoc -I. -I$GOPATH/src/github.com/dartburin --go_out=plugins=grpc:. api/proto/books.proto
protoc -I. -I$GOPATH/src/github.com/dartburin --grpc-gateway_out=logtostderr=true:. api/proto/books.proto

Start
-------
docker-compose up --build

