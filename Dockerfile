FROM golang:1.17.1-alpine

WORKDIR /go/src

COPY . .
RUN go mod download

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build -ldflags '-s -w' main.go