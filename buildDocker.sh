#!/bin/bash
#Builds binary inside alpine
go-bindata-assetfs static_source/*/* static_source/
docker run --rm -v "$GOPATH":/work -e "GOPATH=/work" -e "CGO_ENABLED=1" -e "GOOS=linux" -w /work/src/github.com/tingold/squirrelchopper golang:1.8 go build -ldflags -v
#docker run --rm -v "$GOPATH":/work -e "GOPATH=/work" -e "CGO_ENABLED=1" -e "GOOS=linux" -w /work/src/github.com/tingold/squirrelchopper -e "CC=/usr/local/musl/bin/musl-gcc" tingold/alpine-cgo-musl go build --ldflags '-linkmode external -extldflags "-static"' -v
docker build . -t quay.io/ruptivegeo/chopper:0.0.1
docker push quay.io/ruptivegeo/chopper:0.0.1
