SHELL := $(shell which bash)
ENV = /usr/bin/env

VERSION=$(shell git describe --tags --always)

.SHELLFLAGS = -c
.SILENT: ;
.ONESHELL: ;
.NOTPARALLEL: ;
.EXPORT_ALL_VARIABLES: ;

.PHONY: dep build build-image test coverage

dep:
	go mod download

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -o ./dcc cmd/dcc/*.go

build-image:
	docker build -t krijoh/dcc:${VERSION} -f Dockerfile .

test:
	go test -race -v ./... -cover

coverage:
	go test -race -coverprofile=profile.out -covermode=atomic ./...
	if [ -f profile.out ]; then
			cat profile.out >> coverage.txt
			rm profile.out
	fi