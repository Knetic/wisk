#!/bin/bash

export GOPATH="$(pwd)"
export GOOS=linux
export GOARCH=amd64

go get ./...
go build -o ./.output/${{=project.name=}} .
