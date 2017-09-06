#!/bin/bash
#Builds binary inside alpine
go-bindata-assetfs static_source/*/* static_source/
docker run -it -v "$GOPATH":/work -e "GOPATH=/work" -e "CGO_ENABLED=1" -e "GOOS=linux" -w /work/src/github.com/ruptivespatial/chopper golang:1.8 go build -ldflags -v
#docker run --rm -v "$GOPATH":/work -e "GOPATH=/work" -e "CGO_ENABLED=1" -e "GOOS=linux" -w /work/src/github.com/tingold/squirrelchopper -e "CC=/usr/local/musl/bin/musl-gcc" tingold/alpine-cgo-musl go build --ldflags '-linkmode external -extldflags "-static"' -v
docker build . -t tingold/sc:0.0.18
docker push tingold/sc:0.0.18
