FROM golang:1.20.0-alpine3.17

WORKDIR /go/src

RUN apk update && apk add git && apk add bash

COPY . .

RUN go mod tidy

CMD go run main.go
