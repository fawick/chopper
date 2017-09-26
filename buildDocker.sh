#!/bin/bash
#Builds binary inside alpine
go-bindata-assetfs static_source/*/* static_source/
docker run -it -v "$GOPATH":/work -e "GOPATH=/work" -e "CGO_ENABLED=1" -e "GOOS=linux" -w /work/src/github.com/boundlessgeo/chopper golang:1.8 go build -ldflags -v
docker build . -t quay.io/boundlessgeo/chopper:0.0.18
docker push quay.io/boundlessgeo/chopper:0.0.18
