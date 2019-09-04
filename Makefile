SHELL := $(shell which bash)
ENV = /usr/bin/env

VERSION=$(shell git describe)

.SHELLFLAGS = -c
.SILENT: ;
.ONESHELL: ;
.NOTPARALLEL: ;
.EXPORT_ALL_VARIABLES: ;

.PHONY: dep build test test-coverage

dep:
	go mod download

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -o ./dcc cmd/dcc/*.go

test:
	go test -race -v ./... -cover

test-coverage:
	go test -race -v -cover -coverprofile=coverage/out ./...
	go tool cover -html=coverage/out -o coverage/coverage.html
	rm coverage/out
	go get -u github.com/jstemmer/go-junit-report
	go test ./... -v 2>&1 | go-junit-report > report.xml
