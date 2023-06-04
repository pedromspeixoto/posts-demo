#!/usr/bin/env bash

go get -u github.com/swaggo/swag/cmd/swag@v1.8.3
swag init --parseDependency --parseInternal --parseDepth 5 -g cmd/main.go

export GO111MODULE=on
export GOOS=linux
export GOARCH=amd64

go mod tidy
go build -a -o bin/posts-api cmd/main.go