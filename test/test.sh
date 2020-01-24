#! /bin/bash

echo
curl -XPOST "http://127.0.0.1:8080/books" -d '{"author": "Some author 1", "title": "New book 1"}'
echo
curl -XPOST "http://127.0.0.1:8080/books" -d '{"author": "Some author 2", "title": "New book 2"}'
echo
curl -XPOST "http://127.0.0.1:8080/books" -d '{"author": "Some author 3", "title": "New book 3"}'
echo
curl -XPOST "http://127.0.0.1:8080/books" -d '{"author": "Some author 4", "title": "New book 4"}'
echo
curl -XPOST "http://127.0.0.1:8080/books" -d '{"author": "Some author 5", "title": "New book 5"}'
echo
curl -XPOST "http://127.0.0.1:8080/books" -d '{"author": "Some author 6", "title": "New book 6"}'
echo
curl -XPOST "http://127.0.0.1:8080/books" -d '{"author": "Some author 7", "title": "New book 7"}'
echo
curl -XDELETE "http://127.0.0.1:8080/books/4"
echo
curl -XPATCH "http://127.0.0.1:8080/books/5" -d '{"title": "Book 234264261"}'
echo
curl -XPATCH "http://127.0.0.1:8080/books/6" -d '{"author": "a"}'
echo
curl -XPUT "http://127.0.0.1:8080/books/7" -d '{"author": "A1", "title": "B1"}'
echo
curl -XGET "http://127.0.0.1:8080/books/3"
echo
curl -XGET "http://127.0.0.1:8080/books/5"
echo

