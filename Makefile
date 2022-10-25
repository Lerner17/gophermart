SHELL=/bin/bash

all: prep build perm

prep:
	go mod tidy

build:
	go build -o bin/gophermart ./cmd/gophermart/main.go

perm:
	chmod -R +x bin