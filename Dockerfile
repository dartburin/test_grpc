FROM golang:latest

ARG A_DB_USER=postgres
ARG A_DB_PASS=postgres
ARG A_DB_BASE=books
ARG A_DB_HOST=servdb
ARG A_DB_PORT=5432
ARG A_GRPC_PORT=1111


ENV APP_NAME serverGRPC

ENV DB_USER=${A_DB_USER} 
ENV DB_PASS=${A_DB_PASS} 
ENV DB_BASE=${A_DB_BASE} 
ENV DB_HOST=${A_DB_HOST} 
ENV DB_PORT=${A_DB_PORT} 
ENV GRPC_PORT=${A_GRPC_PORT} 

#RUN go version

RUN mkdir -p ${GOPATH}/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

COPY go.mod ./
COPY go.sum ./

COPY ./cmd ./
COPY ./cmd/server/ ./cmd/
COPY ./cmd/server/* ./cmd/server/

COPY ./api ./
COPY ./api/proto/ ./api/
COPY ./api/proto/* ./api/proto/

COPY ./third_party ./

COPY ./internal/ ./
COPY ./internal/api ./internal/
COPY ./internal/api/server ./internal/api/
COPY ./internal/api/server/* ./internal/api/server/
COPY ./internal/db ./internal/
COPY ./internal/db/* ./internal/db/

#RUN ls -l
#RUN pwd

RUN go mod download 
#RUN protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. cmd/grpc/books/books.proto
#RUN protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. cmd/grpc/books/books.proto
RUN go build cmd/server/main.go
RUN ls -l

EXPOSE ${GRPC_PORT}
EXPOSE ${DB_PORT}


CMD ./main -dbhost=${DB_HOST} -dbbase=${DB_BASE} -dbuser=${DB_USER} -dbpass=${DB_PASS} -dbport=${DB_PORT} -grpcport=${GRPC_PORT}

