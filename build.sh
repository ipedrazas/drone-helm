#!/bin/bash

set -e

export GOOS=linux
export GOARCH=386
go build
docker build -t drone-helm .
