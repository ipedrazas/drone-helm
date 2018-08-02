#!/bin/bash

set -e

export GOOS=linux
export GOARCH=386
go get -u github.com/golang/dep/cmd/dep
dep ensure
go build
docker build -t drone-helm .
