#!/bin/sh

cd $GOPATH/src/github.com/jimmy-go/vovo

if [ "$1" == "bench" ]; then
    go test -race -bench=.
    exit;
fi

if [ "$1" == "html" ]; then
    go tool cover -html=coverage.out
    exit;
fi

go test -cover -coverprofile=coverage.out
