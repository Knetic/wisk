#!/bin/bash

export GOPATH="$(pwd)"

go get ./...
go build -o ./.output/${{=project.name=}} .
