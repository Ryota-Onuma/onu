FROM golang:1.21.0-alpine3.17

WORKDIR /app/

ENV GOPATH=/go

RUN apk update && apk add git && apk add bash

COPY go.mod go.sum ./

RUN go mod download

RUN go get -u golang.org/x/tools/cmd/goyacc

COPY . .

RUN go install  golang.org/x/tools/cmd/goyacc

CMD go run main.go
